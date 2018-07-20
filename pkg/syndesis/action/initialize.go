package action

import (
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"github.com/syndesisio/syndesis-operator/pkg/apis/syndesis/v1alpha1"
	"github.com/syndesisio/syndesis-operator/pkg/syndesis/configuration"
)

// Initializes a Syndesis resource with no status and starts the installation process
type Initialize struct {}


func (a *Initialize) CanExecute(syndesis *v1alpha1.Syndesis) bool {
	return syndesisInstallationStatusIs(syndesis,
		v1alpha1.SyndesisInstallationStatusMissing,
		v1alpha1.SyndesisInstallationStatusNotInstalled)
}

func (a *Initialize) Execute(syndesis *v1alpha1.Syndesis) error {

	list := v1alpha1.NewSyndesisList()
	err := sdk.List(syndesis.Namespace, list)
	if err != nil {
		return err
	}

	syndesisAlreadyInstalled, err := isSyndesisAlreadyInstalled(syndesis)
	if err != nil {
		return err
	}

	target := syndesis.DeepCopy()

	if len(list.Items) > 1 {
		// We want one instance per namespace at most
		target.Status.InstallationStatus = v1alpha1.SyndesisInstallationStatusNotInstalled
		target.Status.Reason = v1alpha1.SyndesisStatusReasonDuplicate
		target.Status.Description = "Cannot install two Syndesis resources in the same namespace"
		logrus.Error("Cannot initialize Syndesis resource ", syndesis.Name, ": duplicate")
	} else if syndesisAlreadyInstalled {
		// One Syndesis CR and Syndesis already installed: integrating the installation in the resource
		target.Status.InstallationStatus = v1alpha1.SyndesisInstallationStatusAttaching
		target.Status.Reason = v1alpha1.SyndesisStatusReasonMissing
		target.Status.Description = "Existing Syndesis installation detected: attaching it to Syndesis resource"
		logrus.Info("Existing Syndesis installation detected: attaching it to Syndesis resource ", syndesis.Name)
	} else {
		syndesisVersion, err := configuration.GetSyndesisVersionFromOperatorTemplate()
		if err != nil {
			return err
		}

		target.Status.InstallationStatus = v1alpha1.SyndesisInstallationStatusInstalling
		target.Status.Reason = v1alpha1.SyndesisStatusReasonMissing
		target.Status.Description = ""
		target.Status.Version = syndesisVersion
		logrus.Info("Syndesis resource ", syndesis.Name, " initialized: installing version ", syndesisVersion)
	}

	return sdk.Update(target)
}

func isSyndesisAlreadyInstalled(syndesis *v1alpha1.Syndesis) (bool, error) {
	// Detects if the configmap with the version is present
	v, err := configuration.GetSyndesisVersionFromNamespace(syndesis.Namespace)
	return err == nil && v != "", err
}
