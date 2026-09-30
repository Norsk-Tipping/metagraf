/*
Copyright 2020 The metaGraf Authors

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

// Package v1alpha1 contains the subset of the ArgoCD Application
// (argoproj.io/v1alpha1) API types that metagraf generates. The field names
// and JSON tags mirror github.com/argoproj/argo-cd/pkg/apis/application/v1alpha1
// so generated manifests stay identical without depending on the argo-cd module.
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	Group           = "argoproj.io"
	Version         = "v1alpha1"
	ApplicationKind = "Application"
)

var SchemeGroupVersion = schema.GroupVersion{Group: Group, Version: Version}

// ApplicationResource is the GroupVersionResource for ArgoCD Applications.
var ApplicationResource = SchemeGroupVersion.WithResource("applications")

type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ApplicationSpec `json:"spec"`
}

type ApplicationSpec struct {
	Source      ApplicationSource      `json:"source"`
	Destination ApplicationDestination `json:"destination"`
	Project     string                 `json:"project"`
	SyncPolicy  *SyncPolicy            `json:"syncPolicy,omitempty"`
	Info        []Info                 `json:"info,omitempty"`
}

type ApplicationSource struct {
	RepoURL        string                      `json:"repoURL"`
	Path           string                      `json:"path,omitempty"`
	TargetRevision string                      `json:"targetRevision,omitempty"`
	Directory      *ApplicationSourceDirectory `json:"directory,omitempty"`
}

type ApplicationSourceDirectory struct {
	Recurse bool                     `json:"recurse,omitempty"`
	Jsonnet ApplicationSourceJsonnet `json:"jsonnet,omitempty"`
	Exclude string                   `json:"exclude,omitempty"`
}

type ApplicationSourceJsonnet struct {
	ExtVars []JsonnetVar `json:"extVars,omitempty"`
	TLAs    []JsonnetVar `json:"tlas,omitempty"`
	Libs    []string     `json:"libs,omitempty"`
}

type JsonnetVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Code  bool   `json:"code,omitempty"`
}

type ApplicationDestination struct {
	Server    string `json:"server,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
}

type SyncPolicy struct {
	Automated   *SyncPolicyAutomated `json:"automated,omitempty"`
	SyncOptions []string             `json:"syncOptions,omitempty"`
	Retry       *RetryStrategy       `json:"retry,omitempty"`
}

type SyncPolicyAutomated struct {
	Prune      bool `json:"prune,omitempty"`
	SelfHeal   bool `json:"selfHeal,omitempty"`
	AllowEmpty bool `json:"allowEmpty,omitempty"`
}

type RetryStrategy struct {
	Limit   int64    `json:"limit,omitempty"`
	Backoff *Backoff `json:"backoff,omitempty"`
}

type Backoff struct {
	Duration    string `json:"duration,omitempty"`
	Factor      *int64 `json:"factor,omitempty"`
	MaxDuration string `json:"maxDuration,omitempty"`
}

type Info struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
