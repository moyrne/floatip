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
	region string
}

func NewTencentAdapter(secretID, secretKey, region string) (*TencentAdapter, error) {
	cred := common.NewCredential(secretID, secretKey)
	cpf := profile.NewClientProfile()
	client, err := lighthouse.NewClient(cred, region, cpf)
	if err != nil {
		return nil, fmt.Errorf("failed to create lighthouse client: %w", err)
	}
	return &TencentAdapter{client: client, region: region}, nil
}

const pageSize = 100

func templateRuleToRule(r *lighthouse.FirewallTemplateRuleInfo) adapter.Rule {
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
	return rule
}

func (a *TencentAdapter) ListRules(ctx context.Context, templateID string) ([]adapter.Rule, error) {
	var rules []adapter.Rule
	var offset int64 = 0
	limit := int64(pageSize)

	for {
		req := lighthouse.NewDescribeFirewallTemplateRulesRequest()
		req.TemplateId = &templateID
		req.Offset = &offset
		req.Limit = &limit

		resp, err := a.client.DescribeFirewallTemplateRulesWithContext(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to list firewall template rules: %w", err)
		}

		for _, r := range resp.Response.TemplateRuleSet {
			rules = append(rules, templateRuleToRule(r))
		}

		if len(resp.Response.TemplateRuleSet) < pageSize {
			break
		}
		offset += int64(len(resp.Response.TemplateRuleSet))
	}
	return rules, nil
}

func firewallRuleToRule(info *lighthouse.FirewallRuleInfo) adapter.Rule {
	rule := adapter.Rule{
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
	return rule
}

func (a *TencentAdapter) ListInstanceRules(ctx context.Context, instanceID string) ([]adapter.Rule, error) {
	var rules []adapter.Rule
	var offset int64 = 0
	limit := int64(pageSize)

	for {
		req := lighthouse.NewDescribeFirewallRulesRequest()
		req.InstanceId = &instanceID
		req.Offset = &offset
		req.Limit = &limit

		resp, err := a.client.DescribeFirewallRulesWithContext(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("failed to list instance firewall rules: %w", err)
		}

		for _, info := range resp.Response.FirewallRuleSet {
			rules = append(rules, firewallRuleToRule(info))
		}

		if len(resp.Response.FirewallRuleSet) < pageSize {
			break
		}
		offset += int64(len(resp.Response.FirewallRuleSet))
	}
	return rules, nil
}

func (a *TencentAdapter) GetRuleByDescription(ctx context.Context, templateID string, description string) (*adapter.Rule, error) {
	rules, err := a.ListRules(ctx, templateID)
	if err != nil {
		return nil, err
	}
	for _, r := range rules {
		if r.Description == description {
			return &r, nil
		}
	}
	return nil, nil
}

func (a *TencentAdapter) GetInstanceRuleByDescription(ctx context.Context, instanceID string, description string) (*adapter.Rule, error) {
	rules, err := a.ListInstanceRules(ctx, instanceID)
	if err != nil {
		return nil, err
	}
	for _, r := range rules {
		if r.Description == description {
			return &r, nil
		}
	}
	return nil, nil
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

func (a *TencentAdapter) ApplyTemplate(ctx context.Context, templateID string, instanceIDs []string) error {
	req := lighthouse.NewApplyFirewallTemplateRequest()
	req.TemplateId = &templateID
	req.ApplyInstances = make([]*lighthouse.InstanceIdentifier, len(instanceIDs))
	for i, id := range instanceIDs {
		req.ApplyInstances[i] = &lighthouse.InstanceIdentifier{InstanceId: &id, Region: &a.region}
	}

	_, err := a.client.ApplyFirewallTemplateWithContext(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to apply firewall template: %w", err)
	}
	return nil
}
