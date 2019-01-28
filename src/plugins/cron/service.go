package cron

import (
	"fmt"
	"github.com/agui2200/roadrunner"
	rr "github.com/agui2200/roadrunner/cmd/rr/cmd"
	"github.com/agui2200/roadrunner/service/rpc"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"os"
	"os/exec"
	"path"
)

// 定时的任务计划

const ID = "cron"

var c = cron.New()
var runningTask []string
var config Config
var logger = rr.Logger

const taskFormat = "%s %s %s"

type Service struct {
}

func (s *Service) Init(r *rpc.Service, cfg *Config) (ok bool, err error) {
	if cfg.WorkDir == "" {
		return false, errors.New("require php script work dir")
	}
	if r != nil {
		err = r.Register(ID, &Command{})
		if err != nil {
			return false, err
		}
	}
	config = *cfg
	c.Start()
	return true, nil
}

func (s *Service) Serve() error {
	return nil
}

func (s *Service) Stop() {
	c.Stop()
}

type Command struct {
}

func (cd *Command) AddFunc(spec, filename string) error {
	filename = path.Join(config.WorkDir, filename)
	finfo, err := os.Stat(filename)
	if err != nil {
		return err
	}
	if finfo == nil {
		return errors.New(fmt.Sprintf("[%s] not exists", filename))
	}
	cmd := exec.Command("php", filename)
	err = c.AddFunc(spec, func() {
		logger.Infof("[%s] running ...", filename)
		w, err := roadrunner.NewPipeFactory().SpawnWorker(cmd)
		if err != nil {
			logger.Error(err)
			return
		}
		go func() {
			err := w.Wait()
			if err != nil {
				logger.Error(err)
			}
		}()
	})
	if err != nil {
		return err
	}
	runningTask = append(runningTask, fmt.Sprintf(taskFormat, spec, cmd.Path, cmd.Args))
	return nil
}

func (c *Command) runningTasks() []string {
	return runningTask
}
