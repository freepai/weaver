package config

import "github.com/spf13/viper"

var (
	WatchUrlKey = "watch.url"
	WatchRefKey = "watch.ref"
	WatchPathKey = "watch.path"
	InternalKey = "internal"

	SSHKeyPrivateFileKey = "sshKey.key"
	SSHKeyPrivatePasswordKey = "sshKey.password"
	SSHKeyPublicFileKey = "sshKey.public"
)

func InitCommon(cfg *viper.Viper) {
	cfg.SetDefault(InternalKey, "@every 10s")
}