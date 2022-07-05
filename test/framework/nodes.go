package framework

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetNodes returns a list of nodes in the cluster
// use options to filter nodes e.g. by label
func (f *Framework) GetNodes(options metav1.ListOptions) (*corev1.NodeList, error) {
	nodes, err := f.ClientSet.CoreV1().Nodes().List(f.Context, options)
	return nodes, err
}

func (f *Framework) SetNodeLabel(node *corev1.Node, key string, value string) {
	node.ObjectMeta.Labels[key] = value
	_, err := f.ClientSet.CoreV1().Nodes().Update(f.Context, node, metav1.UpdateOptions{})
	ExpectNoError(err)
}

func (f *Framework) RemoveNodeLabel(node *corev1.Node, key string) {

	node, err := f.ClientSet.CoreV1().Nodes().Get(f.Context, node.Name, metav1.GetOptions{})
	ExpectNoError(err)

	delete(node.ObjectMeta.Labels, key)

	_, err = f.ClientSet.CoreV1().Nodes().Update(f.Context, node, metav1.UpdateOptions{})
	ExpectNoError(err)
}
