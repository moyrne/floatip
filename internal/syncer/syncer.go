package syncer

import (
	"context"
	"log"
	"slices"
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
	ip, err := detector.DetectPublicIP(s.cfg.IPEchoURL)
	if err != nil {
		log.Printf("failed to detect public IP: %v", err)
		return
	}

	rules, err := s.adapter.ListRules(ctx, s.cfg.TemplateID)
	if err != nil {
		log.Printf("failed to list firewall rules: %v", err)
		return
	}

	idx := slices.IndexFunc(rules, func(r adapter.Rule) bool {
		return r.Description == s.cfg.RuleDescription
	})
	var managedRule *adapter.Rule
	if idx >= 0 {
		managedRule = &rules[idx]
	}

	if managedRule != nil && managedRule.CidrBlock == ip {
		return
	}

	if managedRule != nil {
		managedRule.CidrBlock = ip
		if err := s.adapter.UpdateRule(ctx, s.cfg.TemplateID, *managedRule); err != nil {
			log.Printf("failed to update rule %s: %v", managedRule.ID, err)
			return
		}
		log.Printf("updated rule %s with IP: %s", managedRule.ID, ip)
		return
	}

	if err := s.adapter.CreateRule(ctx, s.cfg.TemplateID, adapter.Rule{
		Protocol:    s.cfg.RuleProtocol,
		Port:        s.cfg.RulePort,
		CidrBlock:   ip,
		Action:      s.cfg.RuleAction,
		Description: s.cfg.RuleDescription,
	}); err != nil {
		log.Printf("failed to create firewall rule: %v", err)
		return
	}
	log.Printf("created firewall rule with description %q for IP: %s", s.cfg.RuleDescription, ip)
}
