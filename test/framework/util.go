package framework

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/utils/pointer"
)

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

func (f *Framework) DeleteTestNamespace(ns string, waitUntilDeleted bool) error {
	err := f.ClientSet.CoreV1().Namespaces().Delete(f.Context, ns, metav1.DeleteOptions{})
	if err != nil {
		if kerrors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if !waitUntilDeleted {
		return nil
	}
	return wait.PollImmediate(time.Second, PollTimeout, func() (bool, error) {
		_, err = f.ClientSet.CoreV1().Namespaces().Get(f.Context, ns, metav1.GetOptions{})
		if kerrors.IsNotFound(err) {
			return true, nil
		}
		return false, err
	})
}

func (f *Framework) GetDefaultSecurityContext() *corev1.SecurityContext {
	return &corev1.SecurityContext{
		RunAsUser: pointer.Int64(12345),
	}
}
