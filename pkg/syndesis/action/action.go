package action

import (
	"github.com/syndesisio/syndesis-operator/pkg/apis/syndesis/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type InstallationAction interface {

	CanExecute(syndesis *v1alpha1.Syndesis) bool

	Execute(syndesis *v1alpha1.Syndesis) error

}

func syndesisInstallationStatusIs(syndesis *v1alpha1.Syndesis, statuses ...v1alpha1.SyndesisInstallationStatus) bool {
	if syndesis == nil {
		return false
	}

	currentStatus := syndesis.Status.InstallationStatus
	for _, status := range statuses {
		if currentStatus == status {
			return true
		}
	}
	return false
}

func customizeKubernetesResource(resource interface{}, syndesis *v1alpha1.Syndesis) {
	if kObj, ok := resource.(metav1.Object); ok {
		kObj.SetNamespace(syndesis.Namespace)

		kObj.SetOwnerReferences([]metav1.OwnerReference{
			*metav1.NewControllerRef(syndesis, schema.GroupVersionKind{
				Group:   v1alpha1.SchemeGroupVersion.Group,
				Version: v1alpha1.SchemeGroupVersion.Version,
				Kind:    syndesis.Kind,
			}),
		})
	}
}
