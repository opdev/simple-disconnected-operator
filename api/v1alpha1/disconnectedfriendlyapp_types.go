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
	"github.com/opdev/simple-disconnected-operator/image"
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
	// ObservedGeneration is the last generation that the controller has acted upon.
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
	// BusyBoxImage is the image that the controller will use for the BusyBox deployment.
	// This is inferred from the environment variable DFA_BUSYBOX_IMAGE or defaulted.
	BusyBoxImage string `json:"busyBoxImage,omitempty"`
	// SleeperImage is the image that the controller will use for the BusyBox deployment.
	// This is inferred from the environment variable DFA_SLEEPER_IMAGE or defaulted.
	SleeperImage string `json:"sleeperImage,omitempty"`
	// DefaultBusyBoxImage is the default value hardcoded for BusyBox in this operator.
	DefaultBusyBoxImage string `json:"defaultBusyBoxImage,omitempty"`
	// DefaultSleeperImage is the default value hardcoded for Sleeper in this operator.
	DefaultSleeperImage string `json:"defaultSleeperImage,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=dfa
// +kubebuilder:printcolumn:name="Resolved BusyBox Image",type=string,JSONPath=`.status.busyBoxImage`
// +kubebuilder:printcolumn:name="Default BusyBox Image",type=string,JSONPath=`.status.defaultBusyBoxImage`
// +kubebuilder:printcolumn:name="Resolved Sleeper Image",type=string,JSONPath=`.status.sleeperImage`
// +kubebuilder:printcolumn:name="Default SleeperImage",type=string,JSONPath=`.status.defaultSleeperImage`

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

// PopulateStatus updates the status based on other values from the same resource.
func (a *DisconnectedFriendlyApp) PopulateStatus() {
	a.Status.ObservedGeneration = a.ObjectMeta.Generation
	a.Status.BusyBoxImage = image.BusyBox()
	a.Status.SleeperImage = image.Sleeper()
	a.Status.DefaultBusyBoxImage = image.DefaultBusyBoxImage
	a.Status.DefaultSleeperImage = image.DefaultSleeperImage
}
