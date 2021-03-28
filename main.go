package main

import (
	"errors"
	"fmt"
	"github.com/freepai/weaver/config"
	"github.com/freepai/weaver/vo"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	log "github.com/phachon/go-logger"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	gitCfg "github.com/go-git/go-git/v5/config"
	"github.com/google/uuid"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/exec"
)

var logger = log.NewLogger()

func gitClone(cfg *viper.Viper, auth transport.AuthMethod) (*git.Worktree, error) {
	// Clone the given repository to the given directory
	logger.Infof("git clone %s", cfg.GetString(config.WatchUrlKey))

	directory := fmt.Sprintf("/tmp/.weaver/%s/", uuid.New())
	url := cfg.GetString(config.WatchUrlKey)

	repo, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL:      url,
		RemoteName: cfg.GetString(config.WatchRefKey),
		Auth: auth,
		Progress: os.Stdout,
	})

	if err!=nil {
		return nil, err
	}

	remote := &gitCfg.RemoteConfig{
		Name: "origin",
		URLs: []string{url},
	}
	_, err = repo.CreateRemote(remote)
	if err!=nil {
		return nil, err
	}

	w, err := repo.Worktree()
	if err!=nil {
		return nil, err
	}

	return w, nil
}

func gitPull(w *git.Worktree, auth transport.AuthMethod) error {
	logger.Info("git pull start")

	err := w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth: auth,
	})

	if err!=nil {
		return errors.New(fmt.Sprintf("Warn in git pull:%s", err.Error()))
	}

	logger.Info("git pull ok")

	return nil
}

func syncToKube(w *git.Worktree, path string) {
	logger.Info("start syncToKube")

	yamlFile, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", w.Filesystem.Root(), path))
	if err != nil {
		logger.Errorf("cmd.Run() failed with %s\n", err)
	}

	logger.Infof("path content:\n%s\n", string(yamlFile))

	app := vo.Application{}
	err = yaml.Unmarshal(yamlFile, &app)
	if err != nil {
		logger.Emergencyf("Unmarshal: %v", err)
	}

	logger.Infof("app: %v", app)

	for _, file := range app.Resources {
		resFile := fmt.Sprintf("%s/%s", w.Filesystem.Root(), file)
		logger.Infof("resFile: %s", resFile)

		kubeApply(resFile)
	}
}

func kubeApply(resFile string) {
	args := []string{"apply", "-f", resFile}
	cmd := exec.Command("kubectl", args...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Infof("combined out:\n%s\n", string(out))
		logger.Errorf("cmd.Run() failed with %s\n", err)
	} else {
		logger.Infof("combined out:\n%s\n", string(out))
	}
}

func main() {
	cfg := config.InitViper()

	// Auth
	privateKeyFile := cfg.GetString(config.SSHKeyPrivateFileKey)
	password := cfg.GetString(config.SSHKeyPrivatePasswordKey)

	auth, err := ssh.NewPublicKeysFromFile("git", privateKeyFile, password)
	if err != nil {
		logger.Emergencyf("generate publickeys failed: %s\n", err.Error())
	}

	w, err := gitClone(cfg, auth)
	if err!=nil {
		logger.Emergencyf("Error in git clone: %v", err)
	}

	logger.Info("git clone ok!")

	// 定时
	c := cron.New()

	cron := cfg.GetString(config.InternalKey)
	c.AddFunc(cron, func() {
		logger.Info("cron executed")

		err := gitPull(w, auth)
		if err!=nil {
			logger.Warningf("Warn in gitPull: %v", err)
			return
		}

		syncToKube(w, cfg.GetString(config.WatchPathKey))
	})

	c.Run()
}