package syncer

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/moyrne/floatip/internal/adapter"
	"github.com/moyrne/floatip/internal/config"
)

type mockAdapter struct {
	rules              []adapter.Rule
	instanceRules      map[string][]adapter.Rule
	applyCallCount     int
	instanceRulesErr   error
}

func (m *mockAdapter) ListRules(_ context.Context, _ string) ([]adapter.Rule, error) {
	return m.rules, nil
}

func (m *mockAdapter) ListInstanceRules(_ context.Context, instanceID string) ([]adapter.Rule, error) {
	if m.instanceRulesErr != nil {
		return nil, m.instanceRulesErr
	}
	if m.instanceRules == nil {
		return m.rules, nil
	}
	rules, ok := m.instanceRules[instanceID]
	if !ok {
		return m.rules, nil
	}
	return rules, nil
}

func (m *mockAdapter) GetRuleByDescription(_ context.Context, _ string, description string) (*adapter.Rule, error) {
	for _, r := range m.rules {
		if r.Description == description {
			return &r, nil
		}
	}
	return nil, nil
}

func (m *mockAdapter) GetInstanceRuleByDescription(_ context.Context, instanceID string, description string) (*adapter.Rule, error) {
	if m.instanceRulesErr != nil {
		return nil, m.instanceRulesErr
	}
	var rules []adapter.Rule
	if m.instanceRules != nil {
		rules = m.instanceRules[instanceID]
	} else {
		rules = m.rules
	}
	for _, r := range rules {
		if r.Description == description {
			return &r, nil
		}
	}
	return nil, nil
}

func (m *mockAdapter) CreateRule(_ context.Context, _ string, rule adapter.Rule) error {
	m.rules = append(m.rules, rule)
	return nil
}

func (m *mockAdapter) UpdateRule(_ context.Context, _ string, rule adapter.Rule) error {
	for i, r := range m.rules {
		if r.ID == rule.ID {
			m.rules[i] = rule
			break
		}
	}
	return nil
}

func (m *mockAdapter) ApplyTemplate(_ context.Context, _ string, _ []string) error {
	m.applyCallCount++
	return nil
}

func testSyncer(t *testing.T, rules []adapter.Rule, ip string, instanceRules map[string][]adapter.Rule) (*Syncer, *mockAdapter) {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(ip + "\n"))
	}))
	t.Cleanup(ts.Close)

	mock := &mockAdapter{rules: rules, instanceRules: instanceRules}
	cfg := config.Default()
	cfg.IPEchoURL = ts.URL
	cfg.InstanceIDs = []string{"ins-test"}
	return New(mock, cfg), mock
}

func TestSync_NoManagedRule_CreatesRule(t *testing.T) {
	s, mock := testSyncer(t, []adapter.Rule{}, "1.2.3.4", map[string][]adapter.Rule{})
	s.sync(context.Background())

	if len(mock.rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(mock.rules))
	}
	if mock.rules[0].CidrBlock != "1.2.3.4" {
		t.Fatalf("expected IP 1.2.3.4, got %s", mock.rules[0].CidrBlock)
	}
	if mock.rules[0].Description != "floatip-sync" {
		t.Fatalf("expected description floatip-sync, got %s", mock.rules[0].Description)
	}
	if mock.applyCallCount != 1 {
		t.Fatalf("expected 1 apply call, got %d", mock.applyCallCount)
	}
}

func TestSync_ManagedRuleIPMatches_Skips(t *testing.T) {
	rule := adapter.Rule{ID: "1", Protocol: "tcp", Port: "22", CidrBlock: "1.2.3.4", Action: "accept", Description: "floatip-sync"}
	s, mock := testSyncer(t, []adapter.Rule{rule}, "1.2.3.4", nil)
	s.sync(context.Background())

	if len(mock.rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(mock.rules))
	}
	if mock.rules[0].CidrBlock != "1.2.3.4" {
		t.Fatalf("expected unchanged IP, got %s", mock.rules[0].CidrBlock)
	}
	if mock.applyCallCount != 0 {
		t.Fatalf("expected 0 apply calls, got %d", mock.applyCallCount)
	}
}

func TestSync_ManagedRuleIPDiffers_UpdatesRule(t *testing.T) {
	templateRule := adapter.Rule{ID: "1", Protocol: "tcp", Port: "22", CidrBlock: "5.6.7.8", Action: "accept", Description: "floatip-sync"}
	instanceRule := adapter.Rule{Protocol: "tcp", Port: "22", CidrBlock: "5.6.7.8", Action: "accept", Description: "floatip-sync"}
	s, mock := testSyncer(t, []adapter.Rule{templateRule}, "1.2.3.4", map[string][]adapter.Rule{"ins-test": {instanceRule}})
	s.sync(context.Background())

	if len(mock.rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(mock.rules))
	}
	if mock.rules[0].CidrBlock != "1.2.3.4" {
		t.Fatalf("expected updated IP 1.2.3.4, got %s", mock.rules[0].CidrBlock)
	}
	if mock.applyCallCount != 1 {
		t.Fatalf("expected 1 apply call, got %d", mock.applyCallCount)
	}
}

func TestSync_UnmanagedRuleWithSameIP_StillCreates(t *testing.T) {
	s, mock := testSyncer(t, []adapter.Rule{
		{ID: "1", Protocol: "tcp", Port: "22", CidrBlock: "1.2.3.4", Action: "accept", Description: "other-rule"},
	}, "1.2.3.4", map[string][]adapter.Rule{})
	s.sync(context.Background())

	if len(mock.rules) != 2 {
		t.Fatalf("expected 2 rules (existing + new), got %d", len(mock.rules))
	}
	if mock.applyCallCount != 1 {
		t.Fatalf("expected 1 apply call, got %d", mock.applyCallCount)
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

func TestEnsureInstanceSync_AllMatch(t *testing.T) {
	rule := &adapter.Rule{Protocol: "tcp", Port: "22", CidrBlock: "1.2.3.4", Action: "accept", Description: "floatip-sync"}
	mock := &mockAdapter{rules: []adapter.Rule{*rule}}
	cfg := config.Default()
	cfg.InstanceIDs = []string{"ins-1"}
	s := New(mock, cfg)

	s.ensureInstanceSync(context.Background(), rule)

	if mock.applyCallCount != 0 {
		t.Fatalf("expected 0 apply calls when all match, got %d", mock.applyCallCount)
	}
}

func TestEnsureInstanceSync_InstanceDiverges(t *testing.T) {
	rule := &adapter.Rule{Protocol: "tcp", Port: "22", CidrBlock: "1.2.3.4", Action: "accept", Description: "floatip-sync"}
	instanceRule := adapter.Rule{Protocol: "tcp", Port: "22", CidrBlock: "9.9.9.9", Action: "accept", Description: "floatip-sync"}
	mock := &mockAdapter{
		rules:         []adapter.Rule{*rule},
		instanceRules: map[string][]adapter.Rule{"ins-1": {instanceRule}},
	}
	cfg := config.Default()
	cfg.InstanceIDs = []string{"ins-1"}
	s := New(mock, cfg)

	s.ensureInstanceSync(context.Background(), rule)

	if mock.applyCallCount != 1 {
		t.Fatalf("expected 1 apply call when instance diverges, got %d", mock.applyCallCount)
	}
}

func TestEnsureInstanceSync_MultipleInstancesOneDiverges(t *testing.T) {
	rule := &adapter.Rule{Protocol: "tcp", Port: "22", CidrBlock: "1.2.3.4", Action: "accept", Description: "floatip-sync"}
	instanceRule2 := adapter.Rule{Protocol: "tcp", Port: "22", CidrBlock: "9.9.9.9", Action: "accept", Description: "floatip-sync"}
	mock := &mockAdapter{
		rules: []adapter.Rule{*rule},
		instanceRules: map[string][]adapter.Rule{
			"ins-1": {*rule},
			"ins-2": {instanceRule2},
		},
	}
	cfg := config.Default()
	cfg.InstanceIDs = []string{"ins-1", "ins-2"}
	s := New(mock, cfg)

	s.ensureInstanceSync(context.Background(), rule)

	if mock.applyCallCount != 1 {
		t.Fatalf("expected 1 apply call when one instance diverges, got %d", mock.applyCallCount)
	}
}

func TestEnsureInstanceSync_NoInstances(t *testing.T) {
	mock := &mockAdapter{}
	cfg := config.Default()
	cfg.InstanceIDs = nil
	s := New(mock, cfg)

	s.ensureInstanceSync(context.Background(), &adapter.Rule{})

	if mock.applyCallCount != 0 {
		t.Fatalf("expected 0 apply calls with no instances, got %d", mock.applyCallCount)
	}
}

func TestEnsureInstanceSync_ListFailureNeedsApply(t *testing.T) {
	rule := &adapter.Rule{Protocol: "tcp", Port: "22", CidrBlock: "1.2.3.4", Action: "accept", Description: "floatip-sync"}
	mock := &mockAdapter{
		rules:            []adapter.Rule{*rule},
		instanceRulesErr: errors.New("list failed"),
	}
	cfg := config.Default()
	cfg.InstanceIDs = []string{"ins-1"}
	s := New(mock, cfg)

	s.ensureInstanceSync(context.Background(), rule)

	if mock.applyCallCount != 1 {
		t.Fatalf("expected 1 apply call when GetInstanceRuleByDescription fails, got %d", mock.applyCallCount)
	}
}
