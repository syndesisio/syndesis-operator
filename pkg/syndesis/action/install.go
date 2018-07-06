package action

import (
	"errors"
	"github.com/openshift/api/template/v1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"github.com/syndesisio/syndesis-operator/pkg/apis/syndesis/v1alpha1"
	"github.com/syndesisio/syndesis-operator/pkg/openshift/serviceaccount"
	"github.com/syndesisio/syndesis-operator/pkg/openshift/template"
	"github.com/syndesisio/syndesis-operator/pkg/util"
	coreV1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

// Install syndesis into the namespace, taking resources from the template

type Install struct {}


func (a *Install) CanExecute(syndesis *v1alpha1.Syndesis) bool {
	return syndesisInstallationStatusIs(syndesis, v1alpha1.SyndesisInstallationStatusInstalling)
}

func (a *Install) Execute(syndesis *v1alpha1.Syndesis) error {

	logrus.Info("Installing Syndesis resource ", syndesis.Name)

	serviceAccountRes, err := util.LoadKubernetesResourceFromAsset("oauth-client-sa.yaml")
	if err != nil {
		return err
	}

	var saName string
	if sa, ok := serviceAccountRes.(*coreV1.ServiceAccount); ok {
		saName = sa.Name
	} else {
		return errors.New("Cannot determine service account name")
	}

	customizeKubernetesResource(serviceAccountRes, syndesis)
	err = sdk.Create(serviceAccountRes)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}

	token, err := serviceaccount.GetServiceAccountToken(saName, syndesis.Namespace)
	if err != nil {
		return err
	}

	res, err := util.LoadKubernetesResourceFromAsset("template.yaml")
	if err != nil {
		return err
	}

	templ := res.(*v1.Template)
	processor, err := template.NewTemplateProcessor(syndesis.Namespace)
	if err != nil {
		return err
	}

	params := make(map[string]string)
	params["ROUTE_HOSTNAME"] = syndesis.Spec.RouteHostName
	params["OPENSHIFT_PROJECT"] = syndesis.Namespace
	params["OPENSHIFT_OAUTH_CLIENT_SECRET"] = token

	list, err := processor.Process(templ, params)
	if err != nil {
		return err
	}

	for _, obj := range list {
		res, err := util.LoadKubernetesResource(obj.Raw)
		if err != nil {
			return err
		}

		customizeKubernetesResource(res, syndesis)

		err = sdk.Create(res)
		if err != nil && !k8serrors.IsAlreadyExists(err) {
			return err
		}
	}

	// Installation completed, set the next state
	target := syndesis.DeepCopy()
	target.Status.InstallationStatus = v1alpha1.SyndesisInstallationStatusStarting
	target.Status.Reason = v1alpha1.SyndesisStatusReasonMissing

	logrus.Info("Syndesis resource ", syndesis.Name, " installed")

	return sdk.Update(target)
}

