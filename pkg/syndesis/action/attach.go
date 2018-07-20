package action

import (
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"github.com/syndesisio/syndesis-operator/pkg/apis/syndesis/v1alpha1"
	"github.com/syndesisio/syndesis-operator/pkg/syndesis/configuration"
	"github.com/syndesisio/syndesis-operator/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Attach a detached Syndesis installation to the Syndesis custom resource
type Attach struct {}


func (a *Attach) CanExecute(syndesis *v1alpha1.Syndesis) bool {
	return syndesisInstallationStatusIs(syndesis, v1alpha1.SyndesisInstallationStatusAttaching)
}

func (a *Attach) Execute(syndesis *v1alpha1.Syndesis) error {

	// Checking that there's only one installation to avoid stealing resources
	if anotherInstallation, err := isAnotherActiveInstallationPresent(syndesis); err != nil {
		return err
	} else if anotherInstallation {
		target := syndesis.DeepCopy()

		target.Status.InstallationStatus = v1alpha1.SyndesisInstallationStatusNotInstalled
		target.Status.Reason = v1alpha1.SyndesisStatusReasonDuplicate
		target.Status.Description = "Cannot merge because another Syndesis resources is present in the same namespace"

		logrus.Error("Cannot merge resource ", syndesis.Name, " because another Syndesis resources is present in the same namespace")
		return sdk.Update(target)
	}

	logrus.Info("Attaching Syndesis installation to resource ", syndesis.Name)

	for _, selector := range getAllManagerSelectors() {
		for _, metaType := range getAllManagedResourceTypes() {

			options := sdk.WithListOptions(&selector)
			list := metav1.List{
				TypeMeta: metaType,
			}
			if err := sdk.List(syndesis.Namespace, &list, options); err != nil {
				return err
			}

			for _, obj := range list.Items {
				res, err := util.LoadKubernetesResource(obj.Raw)
				if err != nil {
					return err
				}
				setNamespaceAndOwnerReference(res, syndesis)
				if err := sdk.Update(res); err != nil {
					return err
				}
			}

		}
	}

	syndesisVersion, err := configuration.GetSyndesisVersionFromNamespace(syndesis.Namespace)
	if err != nil {
		return err
	}

	target := syndesis.DeepCopy()
	target.Status.InstallationStatus = v1alpha1.SyndesisInstallationStatusStarting
	target.Status.Reason = v1alpha1.SyndesisStatusReasonMissing
	target.Status.Description = ""
	target.Status.Version = syndesisVersion

	logrus.Info("Syndesis installation attached to resource ", syndesis.Name)
	return sdk.Update(target)
}

func getAllManagerSelectors() []metav1.ListOptions {
	return []metav1.ListOptions {
		{
			LabelSelector: "syndesis.io/app=syndesis,syndesis.io/type=infrastructure",
		},
		{
			LabelSelector: "app=syndesis,syndesis.io/app=todo",
		},
		{
			LabelSelector: "syndesis.io/app=syndesis,syndesis.io/component=syndesis-db",
		},
	}
}

func getAllManagedResourceTypes() []metav1.TypeMeta {
	return []metav1.TypeMeta{
		{
			APIVersion: "v1",
			Kind: "ConfigMap",
		},
		{
			APIVersion: "v1",
			Kind: "PersistentVolumeClaim",
		},
		{
			APIVersion: "v1",
			Kind: "Secret",
		},
		{
			APIVersion: "v1",
			Kind: "Service",
		},
		{
			APIVersion: "v1",
			Kind: "ServiceAccount",
		},
		{
			APIVersion: "batch/v1beta1",
			Kind: "CronJob",
		},
		{
			APIVersion: "apps.openshift.io/v1",
			Kind: "DeploymentConfig",
		},
		{
			APIVersion: "authorization.openshift.io/v1",
			Kind: "RoleBinding",
		},
		{
			APIVersion: "build.openshift.io/v1",
			Kind: "BuildConfig",
		},
		{
			APIVersion: "image.openshift.io/v1",
			Kind: "ImageStream",
		},
		{
			APIVersion: "route.openshift.io/v1",
			Kind: "Route",
		},
		{
			APIVersion: "template.openshift.io/v1",
			Kind: "Template",
		},
	}
}

func isAnotherActiveInstallationPresent(syndesis *v1alpha1.Syndesis) (bool, error) {
	lst := v1alpha1.NewSyndesisList()
	err := sdk.List(syndesis.Namespace, lst)
	if err != nil {
		return false, err
	}

	for _, that := range lst.Items {
		if that.Name != syndesis.Name &&
			that.Status.InstallationStatus != v1alpha1.SyndesisInstallationStatusNotInstalled &&
			that.Status.InstallationStatus != v1alpha1.SyndesisStatusReasonMissing {
				return true, nil
		}
	}

	return false, nil
}