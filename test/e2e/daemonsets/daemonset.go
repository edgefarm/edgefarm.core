package daemonsets

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/edgefarm/edgefarm.core/test/framework"
	"github.com/loft-sh/vcluster/pkg/util/random"
	"github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	numNodes              = 4 // total nodes in cluster, must match cluster configuration!
	testingContainerName  = "busybox"
	testingContainerImage = "busybox:latest"
	nodeLabelKey          = "tagged"
	edgeLabelKey          = "node-role.kubernetes.io/edge"
)

var _ = ginkgo.Describe("Daemonsets", func() {
	var (
		f         *framework.Framework
		iteration int
		ns        string
		nodes     *corev1.NodeList
	)
	ginkgo.JustBeforeEach(func() {
		// use default framework
		f = framework.DefaultFramework
		iteration++
		ns = fmt.Sprintf("e2e-ds-%d-%s", iteration, random.RandomString(5))

		framework.ExpectNoError(f.CreateTestNamespace(ns))
		var err error
		nodes, err = f.GetNodes(metav1.ListOptions{})
		framework.ExpectNoError(err)

		Expect(len(nodes.Items)).To(BeNumerically("==", numNodes))
	})
	ginkgo.AfterEach(func() {
		// delete test namespace
		err := f.DeleteTestNamespace(ns, false)
		framework.ExpectNoError(err)

		// remove test labels from nodes
		for _, n := range nodes.Items {
			framework.ExpectNoError(f.RemoveNodeLabel(&n, nodeLabelKey))
		}
	})

	ginkgo.It("Daemonset with no labelselector starts pods on all nodes", func() {
		for _, n := range nodes.Items {
			framework.ExpectNoError(f.SetNodeLabel(&n, nodeLabelKey, "dontcare"))
		}

		framework.ExpectNoError(createDaemonSet(ns, map[string]string{}))

		err := wait.PollImmediate(time.Second, time.Hour, func() (bool, error) {
			return podsAreAppliedToAllSelectedNodes(ns, nodes, nodeLabelKey)
		})
		framework.ExpectNoError(err)
	})

	ginkgo.It("Daemonset with labelselector starts pod on one node", func() {
		for _, n := range nodes.Items {
			_, ok := n.ObjectMeta.Labels[edgeLabelKey]
			if ok {
				framework.ExpectNoError(f.SetNodeLabel(&n, nodeLabelKey, "dontcare"))
				break
			}
		}

		framework.ExpectNoError(createDaemonSet(ns, map[string]string{
			nodeLabelKey: "dontcare",
		}))

		err := wait.PollImmediate(time.Second, time.Hour, func() (bool, error) {
			return podsAreAppliedToAllSelectedNodes(ns, nodes, nodeLabelKey)
		})
		framework.ExpectNoError(err)
	})

	ginkgo.It("Daemonset with labelselector starts pods on two nodes", func() {
		i := 0
		for _, n := range nodes.Items {
			_, ok := n.ObjectMeta.Labels[edgeLabelKey]
			if ok {
				framework.ExpectNoError(f.SetNodeLabel(&n, nodeLabelKey, "dontcare"))
				i++
				if i == 2 {
					break
				}
			}
		}

		framework.ExpectNoError(createDaemonSet(ns, map[string]string{
			nodeLabelKey: "dontcare",
		}))

		err := wait.PollImmediate(time.Second, time.Hour, func() (bool, error) {
			return podsAreAppliedToAllSelectedNodes(ns, nodes, nodeLabelKey)
		})
		framework.ExpectNoError(err)
	})

	ginkgo.It("Daemonset Pod Logs are accessible", func() {
		for _, n := range nodes.Items {
			_, ok := n.ObjectMeta.Labels[edgeLabelKey]
			if ok {
				framework.ExpectNoError(f.SetNodeLabel(&n, nodeLabelKey, "dontcare"))
				break
			}
		}

		framework.ExpectNoError(createDaemonSet(ns, map[string]string{
			nodeLabelKey: "dontcare",
		}))

		err := wait.PollImmediate(time.Second, time.Hour, func() (bool, error) {
			return podsAreAppliedToAllSelectedNodes(ns, nodes, nodeLabelKey)
		})
		framework.ExpectNoError(err)

		pods, err := f.ClientSet.CoreV1().Pods(ns).List(f.Context, metav1.ListOptions{})
		framework.ExpectNoError(err)

		err = wait.PollImmediate(time.Second, time.Hour, func() (bool, error) {
			s, err := f.GetPodLog(pods.Items[0], 10)
			framework.ExpectNoError(err)

			if strings.Contains(s, "Hello") {
				return true, nil
			}
			return false, nil
		})
		framework.ExpectNoError(err)
	})

})

func createDaemonSet(nameSpace string, nodeSelector map[string]string) error {
	f := framework.DefaultFramework
	ds := &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DaemonSet",
			APIVersion: "apps/v1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      "busybox-test",
			Namespace: nameSpace,
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "busybox",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "busybox",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            testingContainerName,
							SecurityContext: f.GetDefaultSecurityContext(),
							Image:           testingContainerImage,
							Command: []string{
								"sh",
								"-c",
								"echo Hello && sleep 10 && echo Later && sleep 5 && echo MuchLater",
							},
						},
					},
					Tolerations: []corev1.Toleration{
						{
							Key:      "edgefarm.applications",
							Operator: "Exists",
							Effect:   "NoExecute",
						},
					},
					NodeSelector: nodeSelector,
				},
			},
		},
	}
	err := f.ApplyOrUpdate(ds)
	return err
}

func podsAreAppliedToAllSelectedNodes(nameSpace string, nodes *corev1.NodeList, labelKey string) (bool, error) {
	f := framework.DefaultFramework

	// get all running pods in namespace
	pods, err := f.ClientSet.CoreV1().Pods(nameSpace).List(f.Context, metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	// get of nodes where the pods are running
	podNodes := make([]string, 0)
	for _, p := range pods.Items {
		if p.Status.Phase == corev1.PodRunning {
			podNodes = append(podNodes, p.Spec.NodeName)
		}
	}

	// get list of tagged nodes
	taggedNodes := make([]string, 0)
	for _, n := range nodes.Items {
		_, ok := n.ObjectMeta.Labels[labelKey]
		if ok {
			taggedNodes = append(taggedNodes, n.Name)
		}
	}

	// check if the two lists are identical
	sort.Strings(podNodes)
	sort.Strings(taggedNodes)

	//fmt.Printf("podNodes: %v, taggedNodes: %v\n", podNodes, taggedNodes)

	if len(taggedNodes) == 0 {
		return false, fmt.Errorf("no tagged nodes")
	}

	if len(podNodes) > len(taggedNodes) {
		return false, fmt.Errorf("too many pods started")
	}

	if len(podNodes) == len(taggedNodes) {
		for i := 0; i < len(podNodes); i++ {
			if podNodes[i] != taggedNodes[i] {
				return false, fmt.Errorf("pod started on wrong node")
			}
		}
		return true, nil
	}

	return false, nil
}
