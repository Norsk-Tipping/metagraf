package argocd

import (
	"encoding/json"
	"testing"

	"github.com/laetho/metagraf/pkg/metagraf"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
)

func testGenerator() ApplicationGenerator {
	mg := metagraf.MetaGraf{}
	mg.Metadata.Labels = map[string]string{"team": "x"}

	opts := NewApplicationOptions(func(o *ApplicationOptions) {
		o.Namespace = "myns"
		o.ApplicationProject = "proj"
		o.ApplicationRepoURL = "https://git.example/r.git"
		o.ApplicationRepoPath = "deploy"
		o.ApplicationTargetRevision = "HEAD"
		o.AutomatedSyncPolicySelfHeal = false
		o.SyncPolicyRetry = true
		o.SyncPolicyRetryLimit = 3
	})
	return NewApplicationGenerator(mg, metagraf.MGProperties{}, opts)
}

// The generated manifest must keep the ArgoCD Application wire format.
func TestApplicationJSON(t *testing.T) {
	g := testGenerator()
	app := g.Application("serviceav1")

	b, err := json.Marshal(app)
	assert.NoError(t, err)

	expected := `{
		"apiVersion": "argoproj.io/v1alpha1",
		"kind": "Application",
		"metadata": {
			"creationTimestamp": null,
			"labels": {"team": "x"},
			"name": "serviceav1",
			"namespace": "myns"
		},
		"spec": {
			"destination": {"namespace": "myns", "server": "https://kubernetes.default.svc"},
			"project": "proj",
			"source": {"path": "deploy", "repoURL": "https://git.example/r.git", "targetRevision": "HEAD"},
			"syncPolicy": {"automated": {"prune": true}, "retry": {"limit": 3}}
		}
	}`
	assert.JSONEq(t, expected, string(b))
}

func TestApplicationToUnstructured(t *testing.T) {
	g := testGenerator()
	app := g.Application("serviceav1")

	content, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&app)
	assert.NoError(t, err)
	assert.Equal(t, "Application", content["kind"])
	assert.Equal(t, "argoproj.io/v1alpha1", content["apiVersion"])

	spec := content["spec"].(map[string]interface{})
	source := spec["source"].(map[string]interface{})
	assert.Equal(t, "https://git.example/r.git", source["repoURL"])
}
