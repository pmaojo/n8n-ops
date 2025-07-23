package watch

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pmaojo/n8n-ops/internal/client"
	"github.com/pmaojo/n8n-ops/internal/credentials"
	"github.com/sirupsen/logrus"
)

// Service watches workflows for changes using a polling mechanism.
type Service struct {
	Client            client.Client
	CredentialManager *credentials.CredentialManager
	Logger            logrus.FieldLogger
	Environment       string
}

// Options configures the watch behaviour.
type Options struct {
	Interval   time.Duration
	AutoCommit bool
	AutoSync   bool
}

// NewService constructs a new watcher service.
func NewService(cli client.Client, cm *credentials.CredentialManager, logger logrus.FieldLogger, env string) *Service {
	return &Service{Client: cli, CredentialManager: cm, Logger: logger, Environment: env}
}

// Watch starts monitoring workflow changes until the context is canceled.
func (s *Service) Watch(ctx context.Context, opts Options) error {
	if err := s.Client.HealthCheck(ctx); err != nil {
		return fmt.Errorf("failed to connect to n8n: %w", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	state := make(map[string]time.Time)
	if err := s.initialSync(state); err != nil {
		s.Logger.WithError(err).Warn("initial sync failed")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-sigChan:
			return nil
		case <-ticker.C:
			if err := s.checkChanges(state, opts); err != nil {
				s.Logger.WithError(err).Error("failed to check changes")
			}
		}
	}
}

func (s *Service) initialSync(state map[string]time.Time) error {
	wfs, err := s.Client.GetWorkflows(context.Background())
	if err != nil {
		return err
	}
	for _, wf := range wfs {
		state[wf.ID] = time.Now()
	}
	return nil
}

func (s *Service) checkChanges(state map[string]time.Time, opts Options) error {
	wfs, err := s.Client.GetWorkflows(context.Background())
	if err != nil {
		return err
	}

	for _, wf := range wfs {
		last, ok := state[wf.ID]
		if !ok {
			s.Logger.Infof("new workflow detected: %s", wf.Name)
			state[wf.ID] = time.Now()
			continue
		}
		if time.Since(last) > time.Minute {
			s.Logger.Infof("workflow updated: %s", wf.Name)
			state[wf.ID] = time.Now()
		}
	}

	return nil
}
