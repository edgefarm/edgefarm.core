package framework

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

const edgeLabelKey = "node-role.kubernetes.io/edge"

// GetNodes returns a list of nodes in the cluster
// use options to filter nodes e.g. by label
func (f *Framework) GetNodes(options metav1.ListOptions) (*corev1.NodeList, error) {
	nodes, err := f.ClientSet.CoreV1().Nodes().List(f.Context, options)
	return nodes, err
}

func (f *Framework) SetNodeLabel(node *corev1.Node, key string, value string) error {
	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		n, err := f.ClientSet.CoreV1().Nodes().Get(f.Context, node.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		n.ObjectMeta.Labels[key] = value

		_, err = f.ClientSet.CoreV1().Nodes().Update(f.Context, n, metav1.UpdateOptions{})
		return err
	})
	return err
}

func (f *Framework) RemoveNodeLabel(node *corev1.Node, key string) error {

	err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		n, err := f.ClientSet.CoreV1().Nodes().Get(f.Context, node.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		delete(n.ObjectMeta.Labels, key)

		_, err = f.ClientSet.CoreV1().Nodes().Update(f.Context, n, metav1.UpdateOptions{})
		return err
	})
	return err
}

func (f *Framework) NodeIsReady(n *corev1.Node) bool {
	for _, c := range n.Status.Conditions {
		if c.Type == "Ready" {
			return c.Status == corev1.ConditionTrue
		}
	}
	return false
}

func (f *Framework) LabelReadyEdgeNodes(nameSpace string, numNodes int, labelKey string) error {
	nodes, err := f.GetNodes(metav1.ListOptions{})
	if err != nil {
		return err
	}
	i := 0
	for _, n := range nodes.Items {
		if f.NodeIsReady(&n) {
			_, ok := n.ObjectMeta.Labels[edgeLabelKey]
			if ok {
				err := f.SetNodeLabel(&n, labelKey, "")
				if err != nil {
					return err
				}
				i++
				if i == numNodes {
					break
				}
			}
		}
	}
	if i < numNodes {
		return fmt.Errorf("cannot tag requested number of nodes")
	}
	return nil
}

func (f *Framework) RemoveNodeLabels(labelKey string) error {
	nodes, err := f.GetNodes(metav1.ListOptions{})
	if err != nil {
		return err
	}
	// remove test labels from nodes
	for _, n := range nodes.Items {
		err := f.RemoveNodeLabel(&n, labelKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Framework) GetTaggedNodes(labelKey string) ([]string, error) {
	taggedNodes := make([]string, 0)
	nods, err := f.GetNodes(metav1.ListOptions{})
	if err != nil {
		return taggedNodes, err
	}

	for _, n := range nods.Items {
		_, ok := n.ObjectMeta.Labels[labelKey]
		if ok {
			taggedNodes = append(taggedNodes, n.Name)
		}
	}
	return taggedNodes, nil
}
