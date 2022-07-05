package framework

import (
	"bytes"
	"io"

	corev1 "k8s.io/api/core/v1"
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
