package action

import (
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/syndesisio/syndesis-operator/pkg/apis/syndesis/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	replaceResourcesIfPresent = true
)

type InstallationAction interface {

	CanExecute(syndesis *v1alpha1.Syndesis) bool

	Execute(syndesis *v1alpha1.Syndesis) error

}

type updateFunction func(runtime.Object)

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

func setNamespaceAndOwnerReference(resource interface{}, syndesis *v1alpha1.Syndesis) {
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

func createOrReplace(res runtime.Object) error {
	return createOrReplaceForce(res, false)
}

func createOrReplaceForce(res runtime.Object, force bool) error {
	if err := sdk.Create(res); err != nil && k8serrors.IsAlreadyExists(err) {
		if force || canResourceBeReplaced(res) {
			err = sdk.Delete(res, sdk.WithDeleteOptions(&metav1.DeleteOptions{}))
			if err != nil {
				return err
			}
			return sdk.Create(res)
		} else {
			return nil
		}
	} else {
		return err
	}
}

func updateOnLatestRevision(res runtime.Object, change updateFunction) error {
	change(res)
	err := sdk.Update(res)
	if err != nil && k8serrors.IsConflict(err) {
		attempts := 1
		for attempts <= 5 && err != nil && k8serrors.IsConflict(err) {
			err = sdk.Get(res)
			if err != nil {
				return err
			}

			change(res)
			err = sdk.Update(res)
			attempts++
		}
	}
	return err
}

func canResourceBeReplaced(res runtime.Object) bool {
	if !replaceResourcesIfPresent {
		return false
	}

	if _, blacklisted := res.(*corev1.PersistentVolumeClaim); blacklisted {
		return false
	}
	return true
}