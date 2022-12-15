/*
Copyright 2022.

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

// Container defines container related properties.
type Container struct {
	Image string `json:"image"`
	Port  int    `json:"port"`
}

// Service defines service related properties.
type Service struct {
	Port int `json:"port"`
}

// SchedulingConfig defines scheduling related properties.
type Scheduling struct {
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=23
	StartTime int `json:"startTime"`
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=23
	EndTime int `json:"endTime"`
	// +kubebuilder:validation:Minimum=0
	Replica int `json:"replica"`
}

// TDSetSpec defines the desired state of TDSet
type TDSetSpec struct {
	// +kubebuilder:validation:Required
	Container Container `json:"container"`
	// +kubebuilder:validation:Optional
	Service Service `json:"service,omitempty"`
	// +kubebuilder:validation:Required
	SchedulingConfig []*Scheduling `json:"schedulingConfig"`
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Minimum=1
	DefaultReplica int32 `json:"defaultReplica"`
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=1440
	IntervalMint int32 `json:"intervalMint"`
}

// TDSetStatus defines the observed state of TDSet
type TDSetStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// TDSet is the Schema for the tdsets API
type TDSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TDSetSpec   `json:"spec,omitempty"`
	Status TDSetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TDSetList contains a list of TDSet
type TDSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TDSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TDSet{}, &TDSetList{})
}
