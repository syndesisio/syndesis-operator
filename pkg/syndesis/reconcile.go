package syndesis

import (
	api "github.com/syndesisio/syndesis-operator/pkg/apis/syndesis/v1alpha1"
	"github.com/syndesisio/syndesis-operator/pkg/util"
	"github.com/openshift/api/template/v1"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"github.com/syndesisio/syndesis-operator/pkg/openshift/template"
)

// Reconcile the state of the Syndesis infrastructure elements
// For a new installation: Instantiate the Syndesis Template with the parameters set in the resource
// For an update: execute the Update Pod for the desired version.
func Reconcile(syndesis *api.Syndesis) error {
	if exists, err := Exists(syndesis); err == nil && !exists {
		return Create(syndesis)
	} else {
		// TODO update
		return err
	}
}

func Exists(syndesis *api.Syndesis) (bool, error) {
	svc := corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind: "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: syndesis.Namespace,
			Name: "syndesis-server",
		},
	}
	err := sdk.Get(&svc)
	if err != nil && errors.IsNotFound(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func Create(syndesis *api.Syndesis) error {

	serviceAccountRes, err := util.LoadKubernetesResourceFromAsset("oauth-client-sa.yaml")
	if err != nil {
		return err
	}
	customizeKubernetesResource(serviceAccountRes, syndesis)
	err = sdk.Create(serviceAccountRes)
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
		if err != nil {
			return err
		}
	}

	return nil
}

func customizeKubernetesResource(resource interface{}, syndesis *api.Syndesis) {
	if kObj, ok := resource.(metav1.Object); ok {
		kObj.SetNamespace(syndesis.Namespace)
	}
}
