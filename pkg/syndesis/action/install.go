package action

import (
	"github.com/openshift/api/template/v1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"github.com/syndesisio/syndesis-operator/pkg/apis/syndesis/v1alpha1"
	"github.com/syndesisio/syndesis-operator/pkg/openshift/serviceaccount"
	"github.com/syndesisio/syndesis-operator/pkg/openshift/template"
	"github.com/syndesisio/syndesis-operator/pkg/syndesis/configuration"
	"github.com/syndesisio/syndesis-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Install syndesis into the namespace, taking resources from the template

const (
	replaceResourcesIfPresent = true
)

type Install struct {}


func (a *Install) CanExecute(syndesis *v1alpha1.Syndesis) bool {
	return syndesisInstallationStatusIs(syndesis, v1alpha1.SyndesisInstallationStatusInstalling)
}

func (a *Install) Execute(syndesis *v1alpha1.Syndesis) error {

	logrus.Info("Installing Syndesis resource ", syndesis.Name)

	sa := newSyndesisServiceAccount()
	setNamespaceAndOwnerReference(sa, syndesis)
	// We don't replace the service account if already present, to let Kubernetes generate its tokens
	err := sdk.Create(sa)
	if err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}

	token, err := serviceaccount.GetServiceAccountToken(sa.Name, syndesis.Namespace)
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

	params := configuration.GetEnvVars(syndesis)
	params[string(configuration.EnvOpenshiftOauthClientSecret)] = token

	list, err := processor.Process(templ, params)
	if err != nil {
		return err
	}

	for _, obj := range list {
		res, err := util.LoadKubernetesResource(obj.Raw)
		if err != nil {
			return err
		}

		setNamespaceAndOwnerReference(res, syndesis)

		err = createOrReplace(res)
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

func createOrReplace(res runtime.Object) error {
	if err := sdk.Create(res); err != nil && k8serrors.IsAlreadyExists(err) {
		if canResourceBeReplaced(res) {
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

func canResourceBeReplaced(res runtime.Object) bool {
	if !replaceResourcesIfPresent {
		return false
	}

	if _, blacklisted := res.(*corev1.PersistentVolumeClaim); blacklisted {
		return false
	}
	return true
}

func newSyndesisServiceAccount() *corev1.ServiceAccount {
	sa := corev1.ServiceAccount{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind: "ServiceAccount",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "syndesis-oauth-client",
			Labels: map[string]string{
				"app": "syndesis",
			},
			Annotations: map[string]string {
				"serviceaccounts.openshift.io/oauth-redirecturi.local": "https://localhost:4200",
				"serviceaccounts.openshift.io/oauth-redirecturi.route": "https://",
				"serviceaccounts.openshift.io/oauth-redirectreference.route": `{"kind": "OAuthRedirectReference", "apiVersion": "v1", "reference": {"kind": "Route","name": "syndesis"}}`,
			},
		},
	}

	return &sa
}
