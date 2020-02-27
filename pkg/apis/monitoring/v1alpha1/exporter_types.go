//
// Copyright 2020 IBM Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ExporterSpec defines the desired state of Exporter
type ExporterSpec struct {
	Certs            Certs            `json:"certs,omitempty"`
	Collectd         Collectd         `json:"collectd,omitempty"`
	NodeExporter     NodeExporter     `json:"nodeExporter,omitempty"`
	KubeStateMetrics KubeStateMetrics `json:"kubeStateMetrics,omitempty"`
	ImagePolicy      v1.PullPolicy    `json:"imagePolicy,omitempty"`
	ImagePullSecrets []string         `json:"imagePullSecrets,omitempty"`
	RouterImage      string           `json:"routerImage,omitempty"`
}

// Certs defines certifications used by all exporters
type Certs struct {
	// All certificates for monitoring stack should be signed by this CA
	CASecret string `json:"caSecret"`
	// Exorters' tls cert. Define the secret name. It will not be recreated when existing
	// It can be created by either this operator or prometheus operator. Make sure secret name defined by both operator same
	ExporterSecret string `json:"exporterSecret"`
	// The clusterissuer name. It is used when provider is certmanager
	Issuer string `json:"issuer,omitempty"`
}

// Collectd defines desired state of Collectd exporter
type Collectd struct {
	Enable bool `json:"enable"`
	//SCC privileged should have been added to the service account already by administrator
	ServiceAccount string                  `json:"serviceAccount,omitempty"`
	MetricsPort    int32                   `json:"metricsPort,omitempty"`
	CollectorPort  int32                   `json:"collectorPort,omitempty"`
	Image          string                  `json:"image,omitempty"`
	RouterResource v1.ResourceRequirements `json:"routerResource,omitempty"`
	Resource       v1.ResourceRequirements `json:"resource,omitempty"`
}

// NodeExporter defines desired state of NodeExporter exporter
type NodeExporter struct {
	Enable bool `json:"enable"`
	////SCC privileged should have been added to the service account already by administrator
	ServiceAccount string                  `json:"serviceAccount,omitempty"`
	HostPort       int32                   `json:"hostPort,omitempty"`
	ServicePort    int32                   `json:"servicePort,omitempty"`
	HealthyPort    int32                   `json:"healtyPort,omitempty"`
	RouterResource v1.ResourceRequirements `json:"routerResource,omitempty"`
	Resource       v1.ResourceRequirements `json:"resource,omitempty"`
	Image          string                  `json:"image,omitempty"`
}

// KubeStateMetrics defines desired state of kube-state-metrics
type KubeStateMetrics struct {
	Enable bool `json:"enable"`
	//SCC privileged should have been added to the service account already by administrator
	ServiceAccount string                  `json:"serviceAccount,omitempty"`
	Port           int32                   `json:"port,omitempty"`
	Image          string                  `json:"image,omitempty"`
	RouterResource v1.ResourceRequirements `json:"routerResource,omitempty"`
	Resource       v1.ResourceRequirements `json:"resource,omitempty"`
}

// ExporterStatus defines the observed state of Exporter
type ExporterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Exporter is the Schema for the exporters API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=exporters,scope=Namespaced
type Exporter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ExporterSpec   `json:"spec,omitempty"`
	Status            ExporterStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ExporterList contains a list of Exporter
type ExporterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Exporter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Exporter{}, &ExporterList{})
}
