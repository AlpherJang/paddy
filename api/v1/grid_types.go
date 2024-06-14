/*
Copyright 2024.

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
// GridPhase is a label for the condition of a canary at the current time
type GridPhase string

const GridKind = "Grid"

const (
	InitializingGridPhase GridPhase = "Initializing"
	InitializedGridPhase  GridPhase = "Initialized"
	WaitingGridPhase      GridPhase = "Waiting"
	ProgressingGridPhase  GridPhase = "Progressing"
	FinalisingGridPhase   GridPhase = "Finalising"
	SucceededGridPhase    GridPhase = "Succeeded"
	FailedGridPhase       GridPhase = "Failed"
	TerminatingGridPhase  GridPhase = "Terminating"
	TerminatedGridPhase   GridPhase = "Terminated"
)

type GridService struct {
	// Name of Kubernetes service generated by paddy, default same with target
	// +optional
	Name string `json:"name,omitempty"`
	// Port of service
	// +required
	Port       int32  `json:"port"`
	PortName   string `json:"portName"`
	TargetPort string `json:"targetPort"`
}

type TargetObject struct {
	// Name of the referent
	// +required
	Name string `json:"name"`
	// Kind of referent
	// +required
	Kind string `json:"kind"`
	// APIVersion
	// +required
	APIVersion string `json:"apiVersion"`
}

type AutoScaler struct {
	// Name of scaler, default same with target
	// +optional
	Name string `json:"name,omitempty"`
	// Kind of scaler
	// +required
	Kind string `json:"kind"`
	// APIVersion of scaler
	// +required
	APIVersion string `json:"apiVersion"`
}

type HookType string

const (
	// StartRolloutHook execute webhook before start rollout
	StartRolloutHook HookType = "start-rollout"
	// FinishedRolloutHook execute webhook after the rollout finished
	FinishedRolloutHook HookType = "finished-rollout"
	// RollbackHook rollback canary analysis if webhook returns HTTP 200
	RollbackHook HookType = "rollback"
	// StartPromoteHook execute webhook before promote
	StartPromoteHook HookType = "start-promote"
	// PromotedFinishHook execute webhook after promoted
	PromotedFinishHook HookType = "promoted-finish"
	// PromotingHook execute webhook between promoting, after change traffic weight
	PromotingHook HookType = "promoting"
)

type WebHook struct {
	// Address of Webhook
	// +required
	Address string `json:"address"`
	// Name of WebHook
	Name string `json:"name"`
	// Timeout of WebHook
	Timeout string `json:"timeout,omitempty"`
	// MetaData of Webhook, key-value pair
	// +optional
	MetaData map[string]string `json:"metaData,omitempty"`
	// Number of retry
	// +optional
	Retries int      `json:"retries,omitempty"`
	Type    HookType `json:"type"`
}

// GridSpec defines the desired state of Grid
type GridSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// WorkInstance defined which paddy instance to work with this rollout box
	// +optional
	WorkInstance string `json:"workInstance"`
	// Namespace of resource in Kubernetes,includes deployment\service\scalers
	// +optional
	Namespace string `json:"namespace"`
	// Services of kubernetes
	// +optional
	Services []GridService `json:"services"`
	// AutoScaler of service
	// +optional
	AutoScaler AutoScaler   `json:"autoScaler"`
	TargetRef  TargetObject `json:"targetRef"`
	// WebHooks of Grid
	// +optional
	WebHooks []WebHook `json:"webHooks,omitempty"`
}

// GridStatus defines the observed state of Grid
type GridStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
	Phase      GridPhase          `json:"phase"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Grid is the Schema for the grids API
type Grid struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GridSpec   `json:"spec,omitempty"`
	Status GridStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GridList contains a list of Grid
type GridList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Grid `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Grid{}, &GridList{})
}

func (g *Grid) GetLogName() string {
	return g.Spec.Namespace + g.Spec.TargetRef.Name
}

func (g *Grid) GetPrimaryName() string {
	return g.Spec.TargetRef.Name + "-primary"
}
