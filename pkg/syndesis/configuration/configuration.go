package configuration

import (
	"github.com/syndesisio/syndesis-operator/pkg/apis/syndesis/v1alpha1"
	"strconv"
)

type SyndesisEnvVar string

const (
	EnvRouteHostname 					SyndesisEnvVar = "ROUTE_HOSTNAME"
	//EnvOpenshiftMaster 					SyndesisEnvVar = "OPENSHIFT_MASTER"
	//EnvOpenshiftConsoleUrl				SyndesisEnvVar = "OPENSHIFT_CONSOLE_URL"
	EnvOpenshiftProject					SyndesisEnvVar = "OPENSHIFT_PROJECT"
	EnvOpenshiftOauthClientSecret		SyndesisEnvVar = "OPENSHIFT_OAUTH_CLIENT_SECRET"
	EnvPostgresqlMemoryLimit			SyndesisEnvVar = "POSTGRESQL_MEMORY_LIMIT"
	EnvPostgresqlImageStreamNamespace	SyndesisEnvVar = "POSTGRESQL_IMAGE_STREAM_NAMESPACE"
	EnvPostgresqlUser					SyndesisEnvVar = "POSTGRESQL_USER"
	//EnvPostgresqlPassword				SyndesisEnvVar = "POSTGRESQL_PASSWORD"
	EnvPostgresqlDatabase				SyndesisEnvVar = "POSTGRESQL_DATABASE"
	EnvPostgresqlVolumeCapacity			SyndesisEnvVar = "POSTGRESQL_VOLUME_CAPACITY"
	//EnvPostgresqlSampledbPassword		SyndesisEnvVar = "POSTGRESQL_SAMPLEDB_PASSWORD"
	//EnvTestSupportEnabled				SyndesisEnvVar = "TEST_SUPPORT_ENABLED"
	//EnvOauthCookieSecret				SyndesisEnvVar = "OAUTH_COOKIE_SECRET"
	//EnvSyndesisEncryptKey				SyndesisEnvVar = "SYNDESIS_ENCRYPT_KEY"
	EnvPrometheusVolumeCapacity			SyndesisEnvVar = "PROMETHEUS_VOLUME_CAPACITY"
	EnvPrometheusMemoryLimit			SyndesisEnvVar = "PROMETHEUS_MEMORY_LIMIT"
	EnvMetaVolumeCapacity				SyndesisEnvVar = "META_VOLUME_CAPACITY"
	EnvMetaMemoryLimit					SyndesisEnvVar = "META_MEMORY_LIMIT"
	EnvServerMemoryLimit				SyndesisEnvVar = "SERVER_MEMORY_LIMIT"
	//EnvClientStateAuthenticationKey		SyndesisEnvVar = "CLIENT_STATE_AUTHENTICATION_KEY"
	//EnvClientStateEncryptionKey			SyndesisEnvVar = "CLIENT_STATE_ENCRYPTION_KEY"
	EnvImageStreamNamespace				SyndesisEnvVar = "IMAGE_STREAM_NAMESPACE"
	EnvControllersIntegrationEnabled	SyndesisEnvVar = "CONTROLLERS_INTEGRATION_ENABLED"
	EnvSyndesisRegistry					SyndesisEnvVar = "SYNDESIS_REGISTRY"
	EnvDemoDataEnabled					SyndesisEnvVar = "DEMO_DATA_ENABLED"
	EnvMaxIntegrationsPerUser			SyndesisEnvVar = "MAX_INTEGRATIONS_PER_USER"
)

type SyndesisEnvVarConfig struct {
	Var		SyndesisEnvVar
	Value	string
}

type SyndesisEnvVarExtractor func(syndesis *v1alpha1.Syndesis) *SyndesisEnvVarConfig


var (
	extractors = []SyndesisEnvVarExtractor {
		envOpenshiftProject,
		envRouteHostname,
		envSyndesisRegistry,
		envDemoDataEnabled,
		envMaxIntegrationsPerUser,
		envControllersIntegrationsEnabled,
		envImageStreamNamespace,

		envPostgresqlMemoryLimit,
		envPostgresqlImageStreamNamespace,
		envPostgresqlUser,
		envPostgresqlDatabase,
		envPostgresqlVolumeCapacity,

		envPrometheusMemoryLimit,
		envPrometheusVolumeCapacity,

		envServerMemoryLimit,

		envMetaMemoryLimit,
		envMetaVolumeCapacity,
	}
)

func GetEnvVars(syndesis *v1alpha1.Syndesis) map[string]string {
	configs := make(map[string]string)
	for _, extractor := range extractors {
		conf := extractor(syndesis)
		if conf != nil {
			configs[string(conf.Var)] = conf.Value
		}
	}
	return configs
}


/*
 * List of specific extractors into Environment variables
 */

func envRouteHostname(syndesis *v1alpha1.Syndesis) *SyndesisEnvVarConfig {
	if routeHost := syndesis.Spec.RouteHostName; routeHost != "" {
		return &SyndesisEnvVarConfig{
			Var: EnvRouteHostname,
			Value: routeHost,
		}
	}
	return nil
}

func envOpenshiftProject(syndesis *v1alpha1.Syndesis) *SyndesisEnvVarConfig {
	return &SyndesisEnvVarConfig{
		Var: EnvOpenshiftProject,
		Value: syndesis.Namespace,
	}
}

func envSyndesisRegistry(syndesis *v1alpha1.Syndesis) *SyndesisEnvVarConfig {
	if registry := syndesis.Spec.Registry; registry != "" {
		return &SyndesisEnvVarConfig{
			Var:   EnvSyndesisRegistry,
			Value: registry,
		}
	}
	return nil
}

func envDemoDataEnabled(syndesis *v1alpha1.Syndesis) *SyndesisEnvVarConfig {
	if demodata := syndesis.Spec.DemoData; demodata != nil {
		return &SyndesisEnvVarConfig{
			Var:   EnvDemoDataEnabled,
			Value: strconv.FormatBool(*demodata),
		}
	}
	return nil
}


func envMaxIntegrationsPerUser(syndesis *v1alpha1.Syndesis) *SyndesisEnvVarConfig {
	if integrations := syndesis.Spec.IntegrationLimit; integrations != nil {
		return &SyndesisEnvVarConfig{
			Var: EnvMaxIntegrationsPerUser,
			Value: strconv.Itoa(*integrations),
		}
	}
	return nil
}

func envControllersIntegrationsEnabled(syndesis *v1alpha1.Syndesis) *SyndesisEnvVarConfig {
	if deploy := syndesis.Spec.DeployIntegrations; deploy != nil {
		return &SyndesisEnvVarConfig{
			Var: EnvControllersIntegrationEnabled,
			Value: strconv.FormatBool(*deploy),
		}
	}
	return nil
}

func envImageStreamNamespace(syndesis *v1alpha1.Syndesis) *SyndesisEnvVarConfig {
	if namespace := syndesis.Spec.ImageStreamNamespace; namespace != "" {
		return &SyndesisEnvVarConfig{
			Var: EnvImageStreamNamespace,
			Value: namespace,
		}
	}
	return nil
}

/*
 * Postgres
 */

func envPostgresqlMemoryLimit(syndesis *v1alpha1.Syndesis) *SyndesisEnvVarConfig {
	if limits := syndesis.Spec.Components.Db.Resources.Limits.Memory(); limits != nil && limits.Value() > 0 {
		return &SyndesisEnvVarConfig{
			Var:   EnvPostgresqlMemoryLimit,
			Value: limits.String(),
		}
	}
	return nil
}

func envPostgresqlImageStreamNamespace(syndesis *v1alpha1.Syndesis) *SyndesisEnvVarConfig {
	if ns := syndesis.Spec.Components.Db.ImageStreamNamespace; ns != "" {
		return &SyndesisEnvVarConfig{
			Var:   EnvPostgresqlImageStreamNamespace,
			Value: ns,
		}
	}
	return nil
}

func envPostgresqlUser(syndesis *v1alpha1.Syndesis) *SyndesisEnvVarConfig {
	if user := syndesis.Spec.Components.Db.User; user != "" {
		return &SyndesisEnvVarConfig{
			Var:   EnvPostgresqlUser,
			Value: user,
		}
	}
	return nil
}

func envPostgresqlDatabase(syndesis *v1alpha1.Syndesis) *SyndesisEnvVarConfig {
	if database := syndesis.Spec.Components.Db.Database; database != "" {
		return &SyndesisEnvVarConfig{
			Var:   EnvPostgresqlDatabase,
			Value: database,
		}
	}
	return nil
}

func envPostgresqlVolumeCapacity(syndesis *v1alpha1.Syndesis) *SyndesisEnvVarConfig {
	if capacity := syndesis.Spec.Components.Db.Resources.VolumeCapacity; capacity != "" {
		return &SyndesisEnvVarConfig{
			Var:   EnvPostgresqlVolumeCapacity,
			Value: capacity,
		}
	}
	return nil
}

/*
 * Prometheus
 */

func envPrometheusMemoryLimit(syndesis *v1alpha1.Syndesis) *SyndesisEnvVarConfig {
	if limits := syndesis.Spec.Components.Prometheus.Resources.Limits.Memory(); limits != nil && limits.Value() > 0 {
		return &SyndesisEnvVarConfig{
			Var:   EnvPrometheusMemoryLimit,
			Value: limits.String(),
		}
	}
	return nil
}

func envPrometheusVolumeCapacity(syndesis *v1alpha1.Syndesis) *SyndesisEnvVarConfig {
	if capacity := syndesis.Spec.Components.Prometheus.Resources.VolumeCapacity; capacity != "" {
		return &SyndesisEnvVarConfig{
			Var:   EnvPrometheusVolumeCapacity,
			Value: capacity,
		}
	}
	return nil
}

/*
 * Prometheus
 */

func envServerMemoryLimit(syndesis *v1alpha1.Syndesis) *SyndesisEnvVarConfig {
	if limits := syndesis.Spec.Components.Server.Resources.Limits.Memory(); limits != nil && limits.Value() > 0 {
		return &SyndesisEnvVarConfig{
			Var:   EnvServerMemoryLimit,
			Value: limits.String(),
		}
	}
	return nil
}

/*
 * Meta
 */

func envMetaMemoryLimit(syndesis *v1alpha1.Syndesis) *SyndesisEnvVarConfig {
	if limits := syndesis.Spec.Components.Meta.Resources.Limits.Memory(); limits != nil && limits.Value() > 0 {
		return &SyndesisEnvVarConfig{
			Var:   EnvMetaMemoryLimit,
			Value: limits.String(),
		}
	}
	return nil
}

func envMetaVolumeCapacity(syndesis *v1alpha1.Syndesis) *SyndesisEnvVarConfig {
	if capacity := syndesis.Spec.Components.Meta.Resources.VolumeCapacity; capacity != "" {
		return &SyndesisEnvVarConfig{
			Var:   EnvMetaVolumeCapacity,
			Value: capacity,
		}
	}
	return nil
}

