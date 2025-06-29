/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	healthcheckconfigv1alpha1 "github.com/gardener/gardener/extensions/pkg/apis/config/v1alpha1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	componentbaseconfigv1alpha1 "k8s.io/component-base/config/v1alpha1"
)

// STACKITClusterSpec defines the desired state of STACKITCluster.
type STACKITClusterSpec struct {
	// ClientConnection specifies the kubeconfig file and client connection
	// settings for the proxy server to use when communicating with the apiserver.
	// +optional
	ClientConnection *componentbaseconfigv1alpha1.ClientConnectionConfiguration `json:"clientConnection,omitempty"`
	// ETCD is the etcd configuration.
	ETCD ETCD `json:"etcd"`
	// HealthCheckConfig is the config for the health check controller
	// +optional
	HealthCheckConfig *healthcheckconfigv1alpha1.HealthCheckConfig `json:"healthCheckConfig,omitempty"`
	// StackitRegion sets the STACKIT region.
	StackitRegion string `json:"stackitRegion,omitempty"`
	// StackitAPIEndpoints is the config for STACKIT API Endpoints
	// +optional
	StackitAPIEndpoints *StackitAPIEndpoints `json:"stackitAPIEndpoints,omitempty"`

	// RegistryCaches optionally configures a container registry cache(s) that will be
	// configured on every shoot machine at boot time (and reconciled while its running).
	// +optional
	RegistryCaches []RegistryCacheConfiguration `json:"registryCaches,omitempty"`

	STACKITToken *STACKITToken `json:"stackitToken,,omitempty"`
}

type STACKITToken struct {
	// Name of the secret.
	Name string `json:"name,omitempty"`

	// Key is the name of the key in the secret.
	Key string `json:"key,omitempty"`
}

// ETCD is an etcd configuration.
type ETCD struct {
	// ETCDStorage is the etcd storage configuration.
	Storage ETCDStorage `json:"storage"`
	// ETCDBackup is the etcd backup configuration.
	Backup ETCDBackup `json:"backup"`
}

// ETCDStorage is an etcd storage configuration.
type ETCDStorage struct {
	// ClassName is the name of the storage class used in etcd-main volume claims.
	// +optional
	ClassName *string `json:"className,omitempty"`
	// Capacity is the storage capacity used in etcd-main volume claims.
	// +optional
	Capacity *resource.Quantity `json:"capacity,omitempty"`
}

// ETCDBackup is an etcd backup configuration.
type ETCDBackup struct {
	// Schedule is the etcd backup schedule.
	// +optional
	Schedule *string `json:"schedule,omitempty"`
}

// StackitAPIEndpoints contains all STACKIT API Endpoints.
type StackitAPIEndpoints struct {
	// LoadBalancer is the Endpoint of the LoadBalancer API.
	// +optional
	Loadbalancer *string `json:"loadbalancer,omitempty"`
	// Token is the Token endpoint.
	// +optional
	Token *string `json:"token,omitempty"`
}

// RegistryCacheConfiguration configures a single registry cache.
type RegistryCacheConfiguration struct {
	// Server is the URL of the upstream registry.
	Server string `json:"server"`
	// Cache is the URL of the cache registry.
	Cache string `json:"cache"`
	// CABundle optionally specifies a CA Bundle to trust when connecting to the cache registry.
	CABundle []byte `json:"caBundle,omitempty"`
	// Capabilities optionally specifies what operations the cache registry is capable of.
	Capabilities []string `json:"capabilities,omitempty"`
}

// STACKITClusterStatus defines the observed state of STACKITCluster.
type STACKITClusterStatus struct {
	// Ready indicate STACKIT cluster is ready.
	Ready bool `json:"ready,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// STACKITCluster is the Schema for the stackitclusters API.
type STACKITCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   STACKITClusterSpec   `json:"spec,omitempty"`
	Status STACKITClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// STACKITClusterList contains a list of STACKITCluster.
type STACKITClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []STACKITCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&STACKITCluster{}, &STACKITClusterList{})
}
