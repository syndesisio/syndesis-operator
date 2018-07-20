package action

import (
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"github.com/syndesisio/syndesis-operator/pkg/apis/syndesis/v1alpha1"
	"github.com/syndesisio/syndesis-operator/pkg/syndesis/configuration"
)

// Merge into the Syndesis custom resource the configuration extracted from the installation
// This state does not belong to the standard workflow, it should be triggered manually (e.g. the operator may put a syndesis resource into this state).
type Merge struct {}


func (a *Merge) CanExecute(syndesis *v1alpha1.Syndesis) bool {
	return syndesisInstallationStatusIs(syndesis, v1alpha1.SyndesisInstallationStatusMerging)
}

func (a *Merge) Execute(syndesis *v1alpha1.Syndesis) error {

	logrus.Info("Merging Syndesis legacy configuration into resource ", syndesis.Name)

	config, err := configuration.GetSyndesisEnvVarsFromOpenshiftNamespace(syndesis.Namespace)
	if err != nil {
		return nil
	}

	target := syndesis.DeepCopy()
	configuration.SetConfigurationFromEnvVars(config, target)

	target.Status.InstallationStatus = v1alpha1.SyndesisInstallationStatusAttaching
	target.Status.Reason = v1alpha1.SyndesisStatusReasonMissing
	target.Status.Description = ""

	logrus.Info("Syndesis resource ", syndesis.Name, " merged with legacy installation")
	return sdk.Update(target)
}
