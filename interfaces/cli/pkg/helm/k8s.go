package helm

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

func IsNginxInstalled() (bool, error) {
	ctx := context.Background()
	config := ctrl.GetConfigOrDie()
	clientset := kubernetes.NewForConfigOrDie(config)

	ns, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, n := range ns.Items {
		list, err := clientset.CoreV1().Services(n.GetName()).List(ctx, metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/component=controller,app.kubernetes.io/name=ingress-nginx",
		})
		if err != nil {
			return false, err
		}

		if len(list.Items) > 0 {
			return true, nil
		}
	}

	return false, nil
}
