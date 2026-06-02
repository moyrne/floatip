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
	CreateRule(ctx context.Context, templateID string, rule Rule) error
	UpdateRule(ctx context.Context, templateID string, rule Rule) error
}
