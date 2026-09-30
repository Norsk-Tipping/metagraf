package modules

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteResource(t *testing.T) {
	notFound := errors.New("not found")

	tests := []struct {
		name          string
		getErr        error
		delErr        error
		expectDelCall bool
	}{
		{name: "deletes existing resource", expectDelCall: true},
		{name: "skips missing resource", getErr: notFound, expectDelCall: false},
		{name: "reports failed delete", delErr: errors.New("forbidden"), expectDelCall: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			delCalled := false
			deleteResource("ConfigMap", "test",
				func() error { return tt.getErr },
				func() error {
					delCalled = true
					return tt.delErr
				},
			)
			assert.Equal(t, tt.expectDelCall, delCalled)
		})
	}
}
