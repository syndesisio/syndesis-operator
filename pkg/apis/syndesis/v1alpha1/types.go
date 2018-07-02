package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type SyndesisList struct {
	metav1.TypeMeta  `json:",inline"`
	metav1.ListMeta  `json:"metadata"`
	Items []Syndesis `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Syndesis struct {
	metav1.TypeMeta       `json:",inline"`
	metav1.ObjectMeta     `json:"metadata"`
	Spec   SyndesisSpec   `json:"spec"`
	Status SyndesisStatus `json:"status,omitempty"`
}

type SyndesisSpec struct {
	RouteHostName    string      `json:"routeHostname, omitempty"`
	DemoData         bool        `json:"demoData, omitempty"`
	IntegrationLimit int         `json:"integrationLimit, omitempty"`
	Components       *Components `json:"components, omitempty"`
}

type SyndesisStatus struct {
}

type Components struct {
	Db         DbResources `json:"db, omitempty"`
	Prometheus Resources   `json:"prometheus, omitempty"`
	Server     Resources   `json:"server, omitempty"`
	Meta       Resources   `json:"meta, omitempty"`
}

type DbResources struct {
	Resources v1.ResourceRequirements `json:"resources, omitempty"`
	User      string                  `json:"user, omitempty"`
}

type Resources struct {
	Resources v1.ResourceRequirements `json:"resources, omitempty"`
}
