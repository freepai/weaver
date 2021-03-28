package config

import (
	"flag"
	log "github.com/phachon/go-logger"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

var logger = log.NewLogger()

func InitViper() *viper.Viper {
	cfg := viper.New()

	// 解析命令行
	pflag.StringP("config", "f", "", "config file path")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	// 初始化默认值
	InitCommon(cfg)

	// 从文件中读取配置
	if configFile, _ := pflag.CommandLine.GetString("config"); configFile!="" {
		logger.Infof("configFile: %s", configFile)

		cfg.SetConfigFile(configFile)
	} else {
		cfg.SetConfigName("config")        // name of config file (without extension)
		cfg.SetConfigType("yaml")          // REQUIRED if the config file does not have the extension in the name
		cfg.AddConfigPath("/etc/weaver/")  // path to look for the config file in
		cfg.AddConfigPath("$HOME/.weaver") // call multiple times to add many search paths
		cfg.AddConfigPath(".")             // optionally look for config in the working directory
	}

	if err := cfg.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logger.Warningf("Warn in load config: %v", err)
		} else {
			logger.Emergencyf("Error in load config: %v", err)
		}
	}

	// 从环境变量中获取配置
	cfg.AutomaticEnv()
	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 从命令行参数中获取配置
	cfg.BindPFlags(pflag.CommandLine)

	return cfg
}
