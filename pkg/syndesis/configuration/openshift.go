package configuration

import (
	"errors"
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

const (
	SyndesisGlobalConfigSecret			= "syndesis-global-config"
	SyndesisGlobalConfigVersionProperty	= "syndesis"
	SyndesisGlobalConfigParamsProperty	= "params"
)

func GetSyndesisEnvVarsFromOpenshiftNamespace(namespace string) (map[string]string, error) {
	secret := v1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind: "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name: SyndesisGlobalConfigSecret,
		},
	}
	if err := sdk.Get(&secret); err != nil {
		return nil, err
	}

	if envBlob, present := secret.Data[SyndesisGlobalConfigParamsProperty]; present {
		return parseConfigurationBlob(envBlob), nil
	} else {
		return nil, errors.New("no configuration found")
	}

}

func parseConfigurationBlob(blob []byte) map[string]string {
	strs := strings.Split(string(blob), "\n")
	configs := make(map[string]string, 0)
	for _, conf := range strs {
		conf := strings.Trim(conf, " \r\t")
		if conf == "" {
			continue
		}
		kv := strings.SplitAfterN(conf, "=", 2)
		if len(kv) == 2 {
			key := strings.TrimRight(kv[0], "=")
			value := kv[1]
			configs[key] = value
		}
	}
	return configs
}