package main

import (
	"fmt"
	"github.com/freepai/weaver/config"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	log "github.com/phachon/go-logger"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"github.com/google/uuid"
	"os"
)

var logger = log.NewLogger()

func gitClone(cfg *viper.Viper) (*git.Worktree, error) {
	// Clone the given repository to the given directory
	logger.Infof("git clone %s", cfg.GetString(config.WatchUrlKey))

	// Auth
	privateKeyFile := cfg.GetString(config.SSHKeyPrivateFileKey)
	password := cfg.GetString(config.SSHKeyPrivatePasswordKey)

	publicKeys, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, password)
	if err != nil {
		logger.Warningf("generate publickeys failed: %s\n", err.Error())
		return nil, err
	}

	directory := fmt.Sprintf("/tmp/.weaver/%s/", uuid.New())

	repo, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:      cfg.GetString(config.WatchUrlKey),
		RemoteName: cfg.GetString(config.WatchRefKey),
		Auth: publicKeys,
		Progress: os.Stdout,
	})

	if err!=nil {
		return nil, err
	}

	w, err := repo.Worktree()
	if err!=nil {
		return nil, err
	}

	return w, nil
}

func gitPull(w *git.Worktree) error {
	err := w.Pull(&git.PullOptions{

	})

	if err!=nil {
		return err
	}

	logger.Info("git pull ok")
	return nil
}

func syncToKube(w *git.Worktree, path string) {

}

func main() {
	cfg := config.InitViper()

	w, err := gitClone(cfg)
	if err!=nil {
		logger.Emergencyf("Error in git clone: %v", err)
	}

	// 定时
	c := cron.New()

	c.AddFunc(cfg.GetString(config.InternalKey), func() {
		logger.Info("cron executed")

		err := gitPull(w)
		if err!=nil {
			syncToKube(w, cfg.GetString(config.WatchPathKey))
		}
	})

	c.Run()
}