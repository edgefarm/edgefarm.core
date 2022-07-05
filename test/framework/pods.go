package framework

import (
	"bytes"
	"io"
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
