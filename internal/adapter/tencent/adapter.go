package tencent

import (
	"context"
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"

	"github.com/moyrne/floatip/internal/adapter"
)

type TencentAdapter struct {
	client *lighthouse.Client
}

func NewTencentAdapter(secretID, secretKey, region string) (*TencentAdapter, error) {
	cred := common.NewCredential(secretID, secretKey)
	cpf := profile.NewClientProfile()
	client, err := lighthouse.NewClient(cred, region, cpf)
	if err != nil {
		return nil, fmt.Errorf("failed to create lighthouse client: %w", err)
	}
	return &TencentAdapter{client: client}, nil
}

func (a *TencentAdapter) ListRules(ctx context.Context, templateID string) ([]adapter.Rule, error) {
	req := lighthouse.NewDescribeFirewallTemplateRulesRequest()
	req.TemplateId = &templateID

	resp, err := a.client.DescribeFirewallTemplateRulesWithContext(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list firewall template rules: %w", err)
	}

	rules := make([]adapter.Rule, 0, len(resp.Response.TemplateRuleSet))
	for _, r := range resp.Response.TemplateRuleSet {
		info := r.FirewallRuleInfo
		rule := adapter.Rule{
			ID:       *r.TemplateRuleId,
			Protocol: *info.Protocol,
			Action:   *info.Action,
		}
		if info.CidrBlock != nil {
			rule.CidrBlock = *info.CidrBlock
		}
		if info.Port != nil {
			rule.Port = *info.Port
		}
		if info.FirewallRuleDescription != nil {
			rule.Description = *info.FirewallRuleDescription
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func (a *TencentAdapter) CreateRule(ctx context.Context, templateID string, rule adapter.Rule) error {
	fr := &lighthouse.FirewallRule{
		Protocol:                &rule.Protocol,
		Port:                    &rule.Port,
		CidrBlock:               &rule.CidrBlock,
		Action:                  &rule.Action,
		FirewallRuleDescription: &rule.Description,
	}
	req := lighthouse.NewCreateFirewallTemplateRulesRequest()
	req.TemplateId = &templateID
	req.TemplateRules = []*lighthouse.FirewallRule{fr}

	_, err := a.client.CreateFirewallTemplateRulesWithContext(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create firewall template rule: %w", err)
	}
	return nil
}

func (a *TencentAdapter) UpdateRule(ctx context.Context, templateID string, rule adapter.Rule) error {
	fr := &lighthouse.FirewallRule{
		Protocol:                &rule.Protocol,
		Port:                    &rule.Port,
		CidrBlock:               &rule.CidrBlock,
		Action:                  &rule.Action,
		FirewallRuleDescription: &rule.Description,
	}
	req := lighthouse.NewReplaceFirewallTemplateRuleRequest()
	req.TemplateId = &templateID
	req.TemplateRuleId = &rule.ID
	req.TemplateRule = fr

	_, err := a.client.ReplaceFirewallTemplateRuleWithContext(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to update firewall template rule: %w", err)
	}
	return nil
}
