package syncer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/moyrne/floatip/internal/adapter"
	"github.com/moyrne/floatip/internal/config"
)

type mockAdapter struct {
	rules []adapter.Rule
}

func (m *mockAdapter) ListRules(ctx context.Context, templateID string) ([]adapter.Rule, error) {
	return m.rules, nil
}

func (m *mockAdapter) CreateRule(ctx context.Context, templateID string, rule adapter.Rule) error {
	m.rules = append(m.rules, rule)
	return nil
}

func (m *mockAdapter) UpdateRule(ctx context.Context, templateID string, rule adapter.Rule) error {
	for i, r := range m.rules {
		if r.ID == rule.ID {
			m.rules[i] = rule
			break
		}
	}
	return nil
}

func testSyncer(t *testing.T, rules []adapter.Rule, ip string) *Syncer {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(ip + "\n"))
	}))
	t.Cleanup(ts.Close)

	mock := &mockAdapter{rules: rules}
	cfg := config.Default()
	cfg.IPEchoURL = ts.URL
	return New(mock, cfg)
}

func TestSync_NoManagedRule_CreatesRule(t *testing.T) {
	s := testSyncer(t, []adapter.Rule{}, "1.2.3.4")
	s.sync(context.Background())

	mockRules := s.adapter.(*mockAdapter).rules
	if len(mockRules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(mockRules))
	}
	if mockRules[0].CidrBlock != "1.2.3.4" {
		t.Fatalf("expected IP 1.2.3.4, got %s", mockRules[0].CidrBlock)
	}
	if mockRules[0].Description != "floatip-sync" {
		t.Fatalf("expected description floatip-sync, got %s", mockRules[0].Description)
	}
}

func TestSync_ManagedRuleIPMatches_Skips(t *testing.T) {
	s := testSyncer(t, []adapter.Rule{
		{ID: "1", Protocol: "tcp", Port: "22", CidrBlock: "1.2.3.4", Action: "accept", Description: "floatip-sync"},
	}, "1.2.3.4")
	s.sync(context.Background())

	mockRules := s.adapter.(*mockAdapter).rules
	if len(mockRules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(mockRules))
	}
	if mockRules[0].CidrBlock != "1.2.3.4" {
		t.Fatalf("expected unchanged IP, got %s", mockRules[0].CidrBlock)
	}
}

func TestSync_ManagedRuleIPDiffers_UpdatesRule(t *testing.T) {
	s := testSyncer(t, []adapter.Rule{
		{ID: "1", Protocol: "tcp", Port: "22", CidrBlock: "5.6.7.8", Action: "accept", Description: "floatip-sync"},
	}, "1.2.3.4")
	s.sync(context.Background())

	mockRules := s.adapter.(*mockAdapter).rules
	if len(mockRules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(mockRules))
	}
	if mockRules[0].CidrBlock != "1.2.3.4" {
		t.Fatalf("expected updated IP 1.2.3.4, got %s", mockRules[0].CidrBlock)
	}
}

func TestSync_UnmanagedRuleWithSameIP_StillCreates(t *testing.T) {
	s := testSyncer(t, []adapter.Rule{
		{ID: "1", Protocol: "tcp", Port: "22", CidrBlock: "1.2.3.4", Action: "accept", Description: "other-rule"},
	}, "1.2.3.4")
	s.sync(context.Background())

	mockRules := s.adapter.(*mockAdapter).rules
	if len(mockRules) != 2 {
		t.Fatalf("expected 2 rules (existing + new), got %d", len(mockRules))
	}
}

func TestSync_DetectionFailure_DoesNothing(t *testing.T) {
	mock := &mockAdapter{rules: []adapter.Rule{}}
	cfg := config.Default()
	cfg.IPEchoURL = ""
	s := New(mock, cfg)
	s.sync(context.Background())

	if len(mock.rules) != 0 {
		t.Fatalf("expected 0 rules on detection failure, got %d", len(mock.rules))
	}
}
