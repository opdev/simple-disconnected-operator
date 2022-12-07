/*
Copyright 2022 The OpDev Maintainers.

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

// DisconnectedFriendlyAppSpec defines the desired state of DisconnectedFriendlyApp
type DisconnectedFriendlyAppSpec struct {
	// BusyBoxReplicas is the number of replicas of the BusyBox app associated
	// with this instance of DisconnectedFriendlyApp.
	BusyBoxReplicas *int32 `json:"busyBoxReplicas,omitempty"`
	// SleeperReplicas is the number of replicas of the Sleeper app associated
	// with this instance of DisconnectedFriendlyApp.
	SleeperReplicas *int32 `json:"sleeperReplicas,omitempty"`
}

// DisconnectedFriendlyAppStatus defines the observed state of DisconnectedFriendlyApp
type DisconnectedFriendlyAppStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=dfa

// DisconnectedFriendlyApp is the Schema for the disconnectedfriendlyapps API
type DisconnectedFriendlyApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DisconnectedFriendlyAppSpec   `json:"spec,omitempty"`
	Status DisconnectedFriendlyAppStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DisconnectedFriendlyAppList contains a list of DisconnectedFriendlyApp
type DisconnectedFriendlyAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DisconnectedFriendlyApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DisconnectedFriendlyApp{}, &DisconnectedFriendlyAppList{})
}

var oneReplica int32 = 1

func (a *DisconnectedFriendlyApp) BusyBoxReplicas() *int32 {
	if a.Spec.BusyBoxReplicas == nil {
		return &oneReplica

	}

	return a.Spec.BusyBoxReplicas
}

func (a *DisconnectedFriendlyApp) SleeperReplicas() *int32 {
	if a.Spec.SleeperReplicas == nil {
		return &oneReplica
	}

	return a.Spec.SleeperReplicas
}
