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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// STACKITMachineSpec defines the desired state of STACKITMachine.
type STACKITMachineSpec struct {
	// FloatingPoolName contains the FloatingPoolName name in which LoadBalancer FIPs should be created.
	FloatingPoolName string `json:"floatingPoolName"`
	// FloatingPoolSubnetName contains the fixed name of subnet or matching name pattern for subnet
	// in the Floating IP Pool where the router should be attached to.
	// +optional
	FloatingPoolSubnetName *string `json:"floatingPoolSubnetName,omitempty"`
	// Networks is the OpenStack specific network configuration
	Networks Networks `json:"networks"`
}

// Networks holds information about the Kubernetes and infrastructure networks.
type Networks struct {
	// Router indicates whether to use an existing router or create a new one.
	// +optional
	Router *Router `json:"router,omitempty"`
	// Worker is a CIDRs of a worker subnet (private) to create (used for the VMs).
	// Deprecated - use `workers` instead.
	Worker string `json:"worker"`
	// Workers is a CIDRs of a worker subnet (private) to create (used for the VMs).
	Workers string `json:"workers"`
	// ID is the ID of an existing private network.
	// +optional
	ID *string `json:"id,omitempty"`
	// SubnetID is the ID of an existing subnet.
	// +optional
	SubnetID *string `json:"subnetId,omitempty"`
	// ShareNetwork holds information about the share network (used for shared file systems like NFS)
	// +optional
	ShareNetwork *ShareNetwork `json:"shareNetwork,omitempty"`
	// DNSServers overrides the default dns configuration from cloud profile
	// +optional
	DNSServers *[]string `json:"dnsServers,omitempty"`
}

// Router indicates whether to use an existing router or create a new one.
type Router struct {
	// ID is the router id of an existing OpenStack router.
	ID string `json:"id"`
}

// ShareNetwork holds information about the share network (used for shared file systems like NFS)
type ShareNetwork struct {
	// Enabled is the switch to enable the creation of a share network
	Enabled bool `json:"enabled"`
}

// NodeStatus contains information about Node related resources.
type NodeStatus struct {
	// KeyName is the name of the SSH key.
	KeyName string `json:"keyName"`
}

// NetworkStatus contains information about a generated Network or resources created in an existing Network.
type NetworkStatus struct {
	// ID is the Network id.
	ID string `json:"id"`
	// Name is the Network name.
	Name string `json:"name"`
	// FloatingPool contains information about the floating pool.
	FloatingPool FloatingPoolStatus `json:"floatingPool"`
	// Router contains information about the Router and related resources.
	Router RouterStatus `json:"router"`
	// Subnets is a list of subnets that have been created.
	Subnets []Subnet `json:"subnets"`
	// ShareNetwork contains information about a created/provided ShareNetwork
	// +optional
	ShareNetwork *ShareNetworkStatus `json:"shareNetwork,omitempty"`
}

// RouterStatus contains information about a generated Router or resources attached to an existing Router.
type RouterStatus struct {
	// ID is the Router id.
	ID string `json:"id"`
	// IP is the router ip.
	IP string `json:"ip"`
}

// FloatingPoolStatus contains information about the floating pool.
type FloatingPoolStatus struct {
	// ID is the floating pool id.
	ID string `json:"id"`
	// Name is the floating pool name.
	Name string `json:"name"`
}

// ShareNetworkStatus contains information about a generated ShareNetwork
type ShareNetworkStatus struct {
	// ID is the Network id.
	ID string `json:"id"`
	// Name is the Network name.
	Name string `json:"name"`
}

// Purpose is a purpose of a resource.
type Purpose string

const (
	// PurposeNodes is a Purpose for node resources.
	PurposeNodes Purpose = "nodes"
)

// Subnet is an OpenStack subnet related to a Network.
type Subnet struct {
	// Purpose is a logical description of the subnet.
	Purpose Purpose `json:"purpose"`
	// ID is the subnet id.
	ID string `json:"id"`
}

// SecurityGroup is an OpenStack security group related to a Network.
type SecurityGroup struct {
	// Purpose is a logical description of the security group.
	Purpose Purpose `json:"purpose"`
	// ID is the security group id.
	ID string `json:"id"`
	// Name is the security group name.
	Name string `json:"name"`
}

// STACKITMachineStatus defines the observed state of STACKITMachine.
type STACKITMachineStatus struct {
	// Networks contains information about the created Networks and some related resources.
	Networks NetworkStatus `json:"networks"`
	// Node contains information about Node related resources.
	Node NodeStatus `json:"node"`
	// SecurityGroups is a list of security groups that have been created.
	SecurityGroups []SecurityGroup `json:"securityGroups"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// STACKITMachine is the Schema for the stackitmachines API.
type STACKITMachine struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   STACKITMachineSpec   `json:"spec,omitempty"`
	Status STACKITMachineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// STACKITMachineList contains a list of STACKITMachine.
type STACKITMachineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []STACKITMachine `json:"items"`
}

func init() {
	SchemeBuilder.Register(&STACKITMachine{}, &STACKITMachineList{})
}
