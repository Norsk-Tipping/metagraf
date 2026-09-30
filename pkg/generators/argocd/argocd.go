package argocd

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

import (
	"context"
	gojson "encoding/json"
	"fmt"
	"os"

	"github.com/laetho/metagraf/internal/pkg/k8sclient"
	"github.com/laetho/metagraf/internal/pkg/params"
	argoapp "github.com/laetho/metagraf/pkg/apis/argocd/v1alpha1"
	"github.com/laetho/metagraf/pkg/metagraf"
	"gopkg.in/yaml.v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	log "k8s.io/klog"
)

// Generator for the Application type
type ApplicationGenerator struct {
	MetaGraf   metagraf.MetaGraf
	Properties metagraf.MGProperties
	Options    ApplicationOptions
	Resource   argoapp.Application
}

type ApplicationOption func(*ApplicationOptions)

type ApplicationOptions struct {
	Namespace                       string
	ApplicationProject              string
	ApplicationDestinationNamespace string
	ApplicationRepoURL              string
	ApplicationRepoPath             string
	ApplicationTargetRevision       string
	// Kubernetes API Server endpoint.
	ApplicationDestinationServer      string
	ApplicationSourceDirectoryRecurse bool
	SyncPolicyRetry                   bool
	SyncPolicyRetryLimit              int64
	AutomatedSyncPolicy               bool
	AutomatedSyncPolicyPrune          bool
	AutomatedSyncPolicySelfHeal       bool
}

// Instance of ApplicationOptions that can be used for propagating flags.
var AppOpts ApplicationOptions

func init() {
	// Initialize AppOpts with defaults
	AppOpts.ApplicationDestinationServer = "https://kubernetes.default.svc"
	AppOpts.AutomatedSyncPolicy = true
	AppOpts.AutomatedSyncPolicyPrune = true
	AppOpts.AutomatedSyncPolicySelfHeal = true
	AppOpts.SyncPolicyRetry = false
	AppOpts.SyncPolicyRetryLimit = 5

}

// Creates a new ApplicationOptions based on defaults from AppOpts and runs
// the functional ApplicationOption methods against it from options.
func NewApplicationOptions(options ...ApplicationOption) ApplicationOptions {
	opts := AppOpts
	o := &opts
	for _, opt := range options {
		opt(o)
	}
	return *o
}

// Factory for creating a new ApplicationGenerator
func NewApplicationGenerator(mg metagraf.MetaGraf, prop metagraf.MGProperties, options ApplicationOptions) ApplicationGenerator {
	g := ApplicationGenerator{
		MetaGraf:   mg,
		Properties: prop,
		Options:    options,
	}
	return g
}

func (g *ApplicationOptions) SetOptions(options ...ApplicationOption) {
	for _, opt := range options {
		opt(g)
	}
}

// Sets the Kubernetes namespace for the generated Application resource.
func ApplicationNamespace(ns string) ApplicationOption {
	return func(o *ApplicationOptions) {
		o.Namespace = ns
	}
}

// Sets the Kubernetes namespace for the generated Application resource.
func ApplicationTargetNamespace(ns string) ApplicationOption {
	return func(o *ApplicationOptions) {
		o.Namespace = ns
	}
}

func ApplicationDestinationNamepace(ns string) ApplicationOption {
	return func(g *ApplicationOptions) {
		g.ApplicationDestinationNamespace = ns
	}
}

func (g ApplicationGenerator) applicationSyncPolicy() *argoapp.SyncPolicy {

	sp := &argoapp.SyncPolicy{}

	if g.Options.AutomatedSyncPolicy {
		sp.Automated = &argoapp.SyncPolicyAutomated{
			Prune:    g.Options.AutomatedSyncPolicyPrune,
			SelfHeal: g.Options.AutomatedSyncPolicySelfHeal,
		}
	}

	if g.Options.SyncPolicyRetry {
		rp := &argoapp.RetryStrategy{
			Limit:   g.Options.SyncPolicyRetryLimit,
			Backoff: nil,
		}
		sp.Retry = rp
	}
	return sp
}

func GetArgoCDSourceDirectory() *argoapp.ApplicationSourceDirectory {
	asd := argoapp.ApplicationSourceDirectory{
		Recurse: params.ArgoCDApplicationSourceDirectoryRecurse,
		Jsonnet: argoapp.ApplicationSourceJsonnet{},
	}
	return &asd
}

func (g *ApplicationGenerator) Application(name string) argoapp.Application {

	var meta []argoapp.Info

	obj := argoapp.Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       argoapp.ApplicationKind,
			APIVersion: argoapp.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   g.Options.Namespace,
			Labels:      g.MetaGraf.Metadata.Labels,
			Annotations: g.MetaGraf.Metadata.Annotations,
		},
		Spec: argoapp.ApplicationSpec{
			Destination: argoapp.ApplicationDestination{
				Server: g.Options.ApplicationDestinationServer,
				// todo revert to g.Options.ApplicationDestinationNamespace
				Namespace: g.Options.Namespace,
			},
			Source: argoapp.ApplicationSource{
				RepoURL:        g.Options.ApplicationRepoURL,
				Path:           g.Options.ApplicationRepoPath,
				TargetRevision: g.Options.ApplicationTargetRevision,
			},

			Project:    g.Options.ApplicationProject,
			SyncPolicy: g.applicationSyncPolicy(),
			Info:       meta,
		},
	}
	return obj
}

func OutputApplication(obj argoapp.Application, format string) {
	// Round-trip through a map so keys are sorted the same way for json and yaml.
	b, err := gojson.Marshal(obj)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	jsonMap := make(map[string]interface{})
	err = gojson.Unmarshal(b, &jsonMap)
	if err != nil {
		panic(err)
	}

	if format == "json" {
		oj, err := gojson.MarshalIndent(jsonMap, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(oj))
		return
	} else {
		oy, err := yaml.Marshal(jsonMap)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(oy))
	}
}

func StoreApplication(obj argoapp.Application) {

	log.V(2).Infof("ResourceVersion: %v Length: %v", obj.ResourceVersion, len(obj.ResourceVersion))
	log.V(2).Infof("Namespace: %v", params.NameSpace)

	content, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&obj)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	u := &unstructured.Unstructured{Object: content}

	client := k8sclient.GetArgoCDApplicationClient(params.NameSpace)
	if len(obj.ResourceVersion) > 0 {
		// update
		result, err := client.Update(context.TODO(), u, metav1.UpdateOptions{})
		if err != nil {
			log.Info(err)
			return
		}
		log.Infof("Updated ArgoCD Application: %v(%v)", result.GetName(), obj.Name)
	} else {
		result, err := client.Create(context.TODO(), u, metav1.CreateOptions{})
		if err != nil {
			log.Info(err)
			return
		}
		log.Infof("Created ArgoCD Application: %v(%v)", result.GetName(), obj.Name)
	}
}
