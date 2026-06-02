package syncer

import (
	"context"
	"log"
	"time"

	"github.com/moyrne/floatip/internal/adapter"
	"github.com/moyrne/floatip/internal/config"
	"github.com/moyrne/floatip/internal/detector"
)

type Syncer struct {
	adapter adapter.FirewallAdapter
	cfg     config.Config
}

func New(adapter adapter.FirewallAdapter, cfg config.Config) *Syncer {
	return &Syncer{adapter: adapter, cfg: cfg}
}

func (s *Syncer) Run(ctx context.Context) {
	ticker := time.NewTicker(s.cfg.SyncInterval)
	defer ticker.Stop()

	s.sync(ctx)

	for {
		select {
		case <-ticker.C:
			s.sync(ctx)
		case <-ctx.Done():
			log.Println("syncer stopped")
			return
		}
	}
}

func (s *Syncer) sync(ctx context.Context) {
	log.Printf("sync: %+v", s.cfg.InstanceIDs)
	start := time.Now()
	defer func() { log.Printf("sync cost: %v", time.Since(start)) }()

	ip, err := detector.DetectPublicIP(s.cfg.IPEchoURL)
	if err != nil {
		log.Printf("failed to detect public IP: %v", err)
		return
	}

	if rule := s.syncRule(ctx, ip); rule != nil {
		s.ensureInstanceSync(ctx, rule)
	}
}

// syncRule ensures the template has a rule matching the detected IP.
// Returns the rule if the operation succeeds, nil on error.
func (s *Syncer) syncRule(ctx context.Context, ip string) *adapter.Rule {
	managedRule, err := s.adapter.GetRuleByDescription(ctx, s.cfg.TemplateID, s.cfg.RuleDescription)
	if err != nil {
		log.Printf("failed to find managed rule: %v", err)
		return nil
	}

	if managedRule != nil && managedRule.CidrBlock == ip {
		return managedRule
	}

	if managedRule != nil {
		managedRule.CidrBlock = ip
		if err := s.adapter.UpdateRule(ctx, s.cfg.TemplateID, *managedRule); err != nil {
			log.Printf("failed to update rule %s: %v", managedRule.ID, err)
			return nil
		}
		log.Printf("updated rule %s with IP: %s", managedRule.ID, ip)
		return managedRule
	}

	rule := adapter.Rule{
		Protocol:    s.cfg.RuleProtocol,
		Port:        s.cfg.RulePort,
		CidrBlock:   ip,
		Action:      s.cfg.RuleAction,
		Description: s.cfg.RuleDescription,
	}
	if err := s.adapter.CreateRule(ctx, s.cfg.TemplateID, rule); err != nil {
		log.Printf("failed to create firewall rule: %v", err)
		return nil
	}
	log.Printf("created firewall rule with description %q for IP: %s", s.cfg.RuleDescription, ip)
	return &rule
}

func ruleEqual(a, b *adapter.Rule) bool {
	if a == nil || b == nil {
		return false
	}
	return a.Protocol == b.Protocol &&
		a.Port == b.Port &&
		a.CidrBlock == b.CidrBlock &&
		a.Action == b.Action &&
		a.Description == b.Description
}

func (s *Syncer) ensureInstanceSync(ctx context.Context, rule *adapter.Rule) {
	if len(s.cfg.InstanceIDs) == 0 {
		return
	}

	var drift []string
	for _, instanceID := range s.cfg.InstanceIDs {
		instanceRule, err := s.adapter.GetInstanceRuleByDescription(ctx, instanceID, rule.Description)
		if err != nil {
			log.Printf("failed to check instance %s: %v", instanceID, err)
			drift = append(drift, instanceID)
			continue
		}
		if !ruleEqual(instanceRule, rule) {
			log.Printf("instance %s firewall rule missing or out of date", instanceID)
			drift = append(drift, instanceID)
		}
	}

	if len(drift) == 0 {
		log.Printf("all instances firewall rules are in sync with template")
		return
	}

	log.Printf("instance firewall drift detected for %v, re-applying template", drift)
	s.applyTemplate(ctx, drift)
}

func (s *Syncer) applyTemplate(ctx context.Context, instanceIDs []string) {
	if len(instanceIDs) == 0 {
		return
	}
	if err := s.adapter.ApplyTemplate(ctx, s.cfg.TemplateID, instanceIDs); err != nil {
		log.Printf("failed to apply firewall template: %v", err)
		return
	}
	log.Printf("applied firewall template %s to instances: %v", s.cfg.TemplateID, instanceIDs)
}
