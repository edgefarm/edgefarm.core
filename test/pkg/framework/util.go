package framework

import (
	corev1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"

	"k8s.io/utils/pointer"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func (f *Framework) GetDefaultSecurityContext() *corev1.SecurityContext {
	return &corev1.SecurityContext{
		RunAsUser: pointer.Int64(12345),
	}
}

// ApplyOrUpdate applies or updates the given object.
func (f *Framework) ApplyOrUpdate(obj ctrlclient.Object) error {
	err := f.CtrlClient.Create(f.Context, obj)
	if err != nil {
		if !kerrors.IsAlreadyExists(err) {
			return err
		}
		err = f.CtrlClient.Update(f.Context, obj)
		if err != nil {
			return err
		}
	}
	return nil
}
