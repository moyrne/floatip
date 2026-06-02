package config

import "time"

type Config struct {
	SyncInterval    time.Duration `mapstructure:"sync_interval"`
	IPEchoURL       string        `mapstructure:"ip_echo_url"`
	TemplateID      string        `mapstructure:"template_id"`
	SecretID        string        `mapstructure:"secret_id"`
	SecretKey       string        `mapstructure:"secret_key"`
	Region          string        `mapstructure:"region"`
	RuleProtocol    string        `mapstructure:"rule_protocol"`
	RulePort        string        `mapstructure:"rule_port"`
	RuleAction      string        `mapstructure:"rule_action"`
	RuleDescription string        `mapstructure:"rule_description"`
	InstanceIDs     []string      `mapstructure:"instance_ids"`
}

func Default() Config {
	return Config{
		SyncInterval:    5 * time.Minute,
		IPEchoURL:       "https://checkip.amazonaws.com",
		RuleProtocol:    "ALL",
		RulePort:        "ALL",
		RuleAction:      "accept",
		RuleDescription: "floatip-sync",
	}
}
