/*
Copyright 2026.

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

// WebAppSpec defines the desired state of WebApp
type WebAppSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	//+kubebuilder:validation:MinLength=1
	Image string `json:"image"`

	//+kubebuilder:validation:Minimum=0
	Replicas *int32 `json:"replicas,omitempty"`

	//+kubebuilder:validation:Minimum=1
	//+kubebuilder:validation:Maximum=65535
	Port int32 `json:"port"`
}

// WebAppStatus defines the observed state of WebApp.
type WebAppStatus struct {
	Phase             string `json:"phase,omitempty"`
	Message           string `json:"message,omitempty"`
	AvailableReplicas int32  `json:"availableReplicas,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// WebApp is the Schema for the webapps API
type WebApp struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of WebApp
	// +required
	Spec WebAppSpec `json:"spec"`

	// status defines the observed state of WebApp
	// +optional
	Status WebAppStatus `json:"status,omitzero"`
}

// +kubebuilder:object:root=true

// WebAppList contains a list of WebApp
type WebAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []WebApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WebApp{}, &WebAppList{})
}
