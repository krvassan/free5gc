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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// Free5gcSpec defines the desired state of Free5gc.
type Free5gcSpec struct {
	AMF  Free5gcService `json:"amf,omitempty"`
	SMF  Free5gcService `json:"smf,omitempty"`
	AUSF Free5gcService `json:"ausf,omitempty"`
	PCF  Free5gcService `json:"pcf,omitempty"`
	NRF  Free5gcService `json:"nrf,omitempty"`
	UPF  Free5gcService `json:"upf,omitempty"`
	UDM  Free5gcService `json:"udm,omitempty"`
	UDR  Free5gcService `json:"udr,omitempty"`
	NSSF Free5gcService `json:"nssf,omitempty"`
}

// Free5gcService définit un service Free5gc
type Free5gcService struct {
	Image    string        `json:"image,omitempty"`
	Replicas int           `json:"replicas,omitempty"`
	Network  NetworkConfig `json:"network,omitempty"`
}

// Définit la configuration réseau
type NetworkConfig struct {
	Port    int    `json:"port,omitempty"`
	Address string `json:"address,omitempty"`
}

// Free5gcStatus defines the observed state of Free5gc.
type Free5gcStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Free5gc is the Schema for the free5gcs API.
type Free5gc struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Free5gcSpec   `json:"spec,omitempty"`
	Status Free5gcStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// Free5gcList contains a list of Free5gc.
type Free5gcList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Free5gc `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Free5gc{}, &Free5gcList{})
}
