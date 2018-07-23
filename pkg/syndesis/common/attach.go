package common

import (
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/syndesisio/syndesis-operator/pkg/apis/syndesis/v1alpha1"
	"github.com/syndesisio/syndesis-operator/pkg/util"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func AttachSyndesisToResource(syndesis *v1alpha1.Syndesis) error {

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
				SetNamespaceAndOwnerReference(res, syndesis)
				if err := sdk.Update(res); err != nil {
					return err
				}
			}

		}
	}

	return nil
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
