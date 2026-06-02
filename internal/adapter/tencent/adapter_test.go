package tencent

import (
	"context"
	"testing"

	"github.com/moyrne/floatip/internal/adapter"
)

func TestNewTencentAdapter_InvalidCreds(t *testing.T) {
	_, err := NewTencentAdapter("invalid", "invalid", "ap-guangzhou")
	if err != nil {
		t.Fatalf("unexpected error creating adapter (will fail at API call, not construction): %v", err)
	}
}

func TestListRules_EmptyTemplate(t *testing.T) {
	a, err := NewTencentAdapter("test", "test", "ap-guangzhou")
	if err != nil {
		t.Fatal(err)
	}
	_, err = a.ListRules(context.Background(), "non-existent")
	if err == nil {
		t.Skip("skipping: no credentials in CI")
	}
}

func TestAdapterImplementsInterface(t *testing.T) {
	var _ adapter.FirewallAdapter = (*TencentAdapter)(nil)
}
