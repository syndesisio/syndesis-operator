package template

import (
	"github.com/openshift/api/template/v1"
	"github.com/syndesisio/syndesis-operator/pkg/apis/syndesis/v1alpha1"
	"github.com/syndesisio/syndesis-operator/pkg/openshift/template"
	"github.com/syndesisio/syndesis-operator/pkg/syndesis/configuration"
	"github.com/syndesisio/syndesis-operator/pkg/util"
	"k8s.io/apimachinery/pkg/runtime"
)


type InstallParams struct {
	OAuthClientSecret	string
}

func GetInstallResources(syndesis *v1alpha1.Syndesis, params InstallParams) ([]runtime.RawExtension, error) {
	res, err := util.LoadKubernetesResourceFromAsset("template.yaml")
	if err != nil {
		return nil, err
	}

	templ := res.(*v1.Template)
	processor, err := template.NewTemplateProcessor(syndesis.Namespace)
	if err != nil {
		return nil, err
	}

	config := configuration.GetEnvVars(syndesis)
	config[string(configuration.EnvOpenshiftOauthClientSecret)] = params.OAuthClientSecret

	return processor.Process(templ, config)
}
