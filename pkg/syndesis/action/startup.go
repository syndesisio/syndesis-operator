package action

import (
	"github.com/openshift/api/apps/v1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"github.com/syndesisio/syndesis-operator/pkg/apis/syndesis/v1alpha1"
	"github.com/syndesisio/syndesis-operator/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Waits for all pods to startup, then mark Syndesis as "Running"

type Startup struct {}


func (a *Startup) CanExecute(syndesis *v1alpha1.Syndesis) bool {
	return syndesisInstallationStatusIs(syndesis, v1alpha1.SyndesisInstallationStatusStarting)
}

func (a *Startup) Execute(syndesis *v1alpha1.Syndesis) error {

	options := sdk.WithListOptions(&metav1.ListOptions{
		LabelSelector: "syndesis.io/app=syndesis,syndesis.io/type=infrastructure",
	})
	list := metav1.List{
		TypeMeta: metav1.TypeMeta{
			Kind: "DeploymentConfig",
			APIVersion: "apps.openshift.io/v1",
		},
	}
	if err := sdk.List(syndesis.Namespace, &list, options); err != nil {
		return err
	}

	ready := true
	for _, o := range list.Items {
		if deplObj, err := util.LoadKubernetesResource(o.Raw); err != nil {
			return err
		} else if depl, ok := deplObj.(*v1.DeploymentConfig); ok {
			if depl.Spec.Replicas != depl.Status.ReadyReplicas {
				ready = false
			}
		}
	}

	if ready {
		target := syndesis.DeepCopy()
		target.Status.InstallationStatus = v1alpha1.SyndesisInstallationStatusInstalled
		target.Status.Reason = v1alpha1.SyndesisStatusReasonMissing
		logrus.Info("Syndesis resource ", syndesis.Name, " started up")
		return sdk.Update(target)
	} else {
		logrus.Info("Waiting for Syndesis resource ", syndesis.Name, " to startup")
		return nil
	}
}

