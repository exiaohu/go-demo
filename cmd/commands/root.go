package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/exiaohu/go-demo/config"
	"github.com/exiaohu/go-demo/pkg/logger"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "playground",
		Short: "A playground application",
		Long:  `A playground application built with Go, demonstrating best practices.`,
	}
)

// Execute 执行根命令
func Execute(gitCommit, buildTime string) {
	// 注入版本信息
	rootCmd.AddCommand(newVersionCmd(gitCommit, buildTime))
	rootCmd.AddCommand(newServerCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
}

func initConfig() {
	// 这里可以放置初始化配置的逻辑，如果需要的话
}

// loadConfigHelper 辅助加载配置和日志
func loadConfigHelper() (*config.Config, error) {
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if err := logger.Initialize(cfg.Debug); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	return cfg, nil
}
