package framework

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

// GetPodLog returns the last log lines from the log
// Multiple calls to this func will return the same lines again, unless more than tailLines are
// available
func (f *Framework) GetPodLog(pod corev1.Pod, tailLines int64) (string, error) {
	req := f.ClientSet.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
		TailLines: &tailLines,
	})
	stream, err := req.Stream(f.Context)
	if err != nil {
		return "", err
	}
	defer stream.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, stream)
	if err != nil {
		return "", err
	}
	str := buf.String()
	return str, nil
}

func (f *Framework) WaitForPodRunning(podName string, ns string) error {
	return wait.PollImmediate(time.Second, PollTimeout, func() (bool, error) {
		pod, err := f.ClientSet.CoreV1().Pods(ns).Get(f.Context, podName, metav1.GetOptions{})
		if err != nil {
			if kerrors.IsNotFound(err) {
				return false, nil
			}
			return false, err
		}
		if pod.Status.Phase != corev1.PodRunning {
			return false, nil
		}
		return true, nil
	})
}

// GetRunningPodsNodeNames returns a list of node names for running pods matching podNamePrefix
func (f *Framework) GetRunningPodsNodeNames(nameSpace string, podNamePrefix string) []string {
	podNodes := make([]string, 0)

	// get all running pods in namespace
	pods, err := f.ClientSet.CoreV1().Pods(nameSpace).List(f.Context, metav1.ListOptions{})
	if err != nil {
		return podNodes // empty list
	}

	// get nodes where the pods are running
	for _, p := range pods.Items {
		if strings.HasPrefix(p.Name, podNamePrefix) && p.Status.Phase == corev1.PodRunning {
			podNodes = append(podNodes, p.Spec.NodeName)
		}
	}
	return podNodes
}

// GetRunnningPodsNames returns a list of pod names for running pods matching podNamePrefix
func (f *Framework) GetRunningPodsNames(nameSpace string, podNamePrefix string) []string {
	podNames := make([]string, 0)

	// get all running pods in namespace
	pods, err := f.ClientSet.CoreV1().Pods(nameSpace).List(f.Context, metav1.ListOptions{})
	if err != nil {
		return podNames // empty list
	}

	for _, p := range pods.Items {
		if strings.HasPrefix(p.Name, podNamePrefix) && p.Status.Phase == corev1.PodRunning {
			podNames = append(podNames, p.Name)
		}
	}
	return podNames
}

// GetPodImage returns the image of the given pod/container
func (f *Framework) GetPodImage(nameSpace string, podName string, containerName string) (string, error) {

	pod, err := f.GetPodByName(nameSpace, podName)
	if err != nil {
		return "", err
	}

	for _, c := range pod.Spec.Containers {
		if c.Name == containerName {
			return c.Image, nil
		}
	}
	return "", fmt.Errorf("containerName not found")
}

// GetPodByName returns the pod with the given name
func (f *Framework) GetPodByName(nameSpace string, podName string) (*corev1.Pod, error) {
	return f.ClientSet.CoreV1().Pods(nameSpace).Get(f.Context, podName, metav1.GetOptions{})
}

func (f *Framework) WaitForNoPodsInNamespace(nameSpace string, timeout time.Duration) error {
	err := wait.PollImmediate(time.Second, timeout, func() (bool, error) {

		// get all running pods in namespace
		pods, err := f.ClientSet.CoreV1().Pods(nameSpace).List(f.Context, metav1.ListOptions{})
		if err != nil {
			return false, err
		}
		if len(pods.Items) == 0 {
			return true, nil
		}
		return false, nil
	})
	return err
}
