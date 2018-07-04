package openshift

import (
	"github.com/operator-framework/operator-sdk/pkg/k8sclient"

	templateclientset "github.com/openshift/origin/pkg/template/generated/internalclientset"
	templateclient "github.com/openshift/origin/pkg/template/client/internalversion"
	"github.com/openshift/origin/pkg/template/apis/template"
	v1template "github.com/openshift/api/template/v1"
	templateconversion "github.com/openshift/origin/pkg/template/apis/template/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type TemplateProcessor struct {
	namespace	string
}

func NewTemplateProcessor(namespace string) *TemplateProcessor {
	return &TemplateProcessor{
		namespace: namespace,
	}
}

func (t *TemplateProcessor) Process(sourceTemplate *v1template.Template, parameters map[string]string) ([]runtime.RawExtension, error) {

	config := k8sclient.GetKubeConfig()
	templateClient, err := templateclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	client := templateClient.Template()
	processor := templateclient.NewTemplateProcessorClient(client.RESTClient(), t.namespace)

	templ := template.Template{}
	err = templateconversion.Convert_v1_Template_To_template_Template(sourceTemplate, &templ, nil)
	if err != nil {
		return nil, err
	}

	t.fillParameters(&templ, parameters)

	procTempl, err := processor.Process(&templ)
	if err != nil {
		return nil, err
	}

	processed := v1template.Template{}
	err = templateconversion.Convert_template_Template_To_v1_Template(procTempl, &processed, nil)
	if err != nil {
		return nil, err
	}

	return processed.Objects, nil
}

func (t *TemplateProcessor) fillParameters(template *template.Template, parameters map[string]string) {
	for i, param := range template.Parameters {
		if value, ok := parameters[param.Name]; ok {
			template.Parameters[i].Value = value
		}
	}
}