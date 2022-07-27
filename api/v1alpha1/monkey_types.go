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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MonkeySpec defines the desired state of Monkey
type MonkeySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// noop defines whether to log only
	// +optional
	Noop bool `json:"noop,omitempty"`

	// interval defines interval to requeue Chaos experiment to kill a random pod with matching selector
	// +optional
	Interval string `json:"interval,omitempty"`

	// Namespace defines namespace to search for pods to delete
	Namespace string `json:"namespace,omitempty"`

	Selector metav1.LabelSelector `json:"selector,omitempty"`
}

// MonkeyStatus defines the observed state of Monkey
type MonkeyStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Monkey is the Schema for the monkeys API
type Monkey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MonkeySpec   `json:"spec,omitempty"`
	Status MonkeyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MonkeyList contains a list of Monkey
type MonkeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Monkey `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Monkey{}, &MonkeyList{})
}
