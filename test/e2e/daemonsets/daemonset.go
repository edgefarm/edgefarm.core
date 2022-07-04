package daemonsets

import (
	"fmt"

	"github.com/edgefarm/edgefarm.core/test/framework"
	"github.com/loft-sh/vcluster/pkg/util/random"
	"github.com/onsi/ginkgo/v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	testingContainerName  = "nginx"
	testingContainerImage = "nginxinc/nginx-unprivileged"
	initialNsLabelKey     = "testing-ns-label"
	initialNsLabelValue   = "testing-ns-label-value"
)

var _ = ginkgo.Describe("Daemonsets", func() {
	var (
		f         *framework.Framework
		iteration int
		ns        string
	)
	ginkgo.JustBeforeEach(func() {
		// use default framework
		f = framework.DefaultFramework
		iteration++
		ns = fmt.Sprintf("e2e-pods-%d-%s", iteration, random.RandomString(5))

		// create test namespace
		_, err := f.ClientSet.CoreV1().Namespaces().Create(f.Context, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{
			Name:   ns,
			Labels: map[string]string{initialNsLabelKey: initialNsLabelValue},
		}}, metav1.CreateOptions{})
		framework.ExpectNoError(err)
	})
	ginkgo.AfterEach(func() {
		// delete test namespace
		err := f.DeleteTestNamespace(ns, false)
		framework.ExpectNoError(err)
	})

	ginkgo.It("Test daemonsets ", func() {

		ds := &appsv1.DaemonSet{
			TypeMeta: metav1.TypeMeta{
				Kind:       "DaemonSet",
				APIVersion: "apps/v1",
			},

			ObjectMeta: metav1.ObjectMeta{
				Name:      "nginx-all-nodes",
				Namespace: ns,
			},
			Spec: appsv1.DaemonSetSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: map[string]string{
						"app": "nginx",
					},
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"app": "nginx",
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Name:            testingContainerName,
								SecurityContext: f.GetDefaultSecurityContext(),
								Image:           testingContainerImage,
							},
						},
						Tolerations: []corev1.Toleration{
							{
								Key:      "edgefarm.applications",
								Operator: "Exists",
								Effect:   "NoExecute",
							},
						},
					},
				},
			},
		}
		err := f.ApplyOrUpdate(ds)
		framework.ExpectNoError(err)

	})

})
