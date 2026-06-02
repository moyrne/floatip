package adapter

import "context"

type Rule struct {
	ID          string
	Protocol    string
	Port        string
	CidrBlock   string
	Action      string
	Description string
}

type FirewallAdapter interface {
	ListRules(ctx context.Context, templateID string) ([]Rule, error)
	ListInstanceRules(ctx context.Context, instanceID string) ([]Rule, error)
	GetRuleByDescription(ctx context.Context, templateID string, description string) (*Rule, error)
	GetInstanceRuleByDescription(ctx context.Context, instanceID string, description string) (*Rule, error)
	CreateRule(ctx context.Context, templateID string, rule Rule) error
	UpdateRule(ctx context.Context, templateID string, rule Rule) error
	ApplyTemplate(ctx context.Context, templateID string, instanceIDs []string) error
}
