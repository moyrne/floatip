package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"

	"github.com/moyrne/floatip/internal/adapter/tencent"
	"github.com/moyrne/floatip/internal/config"
	"github.com/moyrne/floatip/internal/syncer"
)

func main() {
	cfg := loadConfig()

	adapter, err := tencent.NewTencentAdapter(cfg.SecretID, cfg.SecretKey, cfg.Region)
	if err != nil {
		log.Fatalf("failed to create tencent adapter: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	s := syncer.New(adapter, cfg)
	go s.Run(ctx)

	<-sigCh
	log.Println("shutting down...")
	cancel()
}

func loadConfig() config.Config {
	cfg := config.Default()

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/floatip")
	v.SetEnvPrefix("FLOATIP")
	v.AutomaticEnv()

	v.SetDefault("sync_interval", cfg.SyncInterval)
	v.SetDefault("ip_echo_url", cfg.IPEchoURL)
	v.SetDefault("template_id", cfg.TemplateID)
	v.SetDefault("secret_id", cfg.SecretID)
	v.SetDefault("secret_key", cfg.SecretKey)
	v.SetDefault("region", cfg.Region)
	v.SetDefault("rule_protocol", cfg.RuleProtocol)
	v.SetDefault("rule_port", cfg.RulePort)
	v.SetDefault("rule_action", cfg.RuleAction)
	v.SetDefault("rule_description", cfg.RuleDescription)

	_ = v.ReadInConfig()
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}

	return cfg
}
