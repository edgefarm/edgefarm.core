package framework

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (f *Framework) SetSecret(nameSpace string, secretName string, key string, value string) {
	_, err := f.ClientSet.CoreV1().Secrets(nameSpace).Create(f.Context, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
		},
		StringData: map[string]string{
			key: value,
		},
		Type: corev1.SecretTypeOpaque,
	}, metav1.CreateOptions{})
	ExpectNoError(err)
}
