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
	isync "github.com/pmaojo/n8n-ops/internal/sync"
	"github.com/sirupsen/logrus"
)

// gitAutoCommitter defines the subset of git.GitStatusChecker used by the watch service.
type gitAutoCommitter interface {
	AutoCommitWorkflows(ctx context.Context, message string) (string, error)
}

// workflowSyncer abstracts the sync.Service dependency.
type workflowSyncer interface {
	Sync(ctx context.Context, opts isync.Options) error
}

// Service watches workflows for changes using a polling mechanism.
type Service struct {
	Client            client.Client
	CredentialManager *credentials.CredentialManager
	Logger            logrus.FieldLogger
	Environment       string
	Git               gitAutoCommitter
	Syncer            workflowSyncer
}

// Options configures the watch behaviour.
type Options struct {
	Interval   time.Duration
	AutoCommit bool
	AutoSync   bool
}

// NewService constructs a new watcher service.
func NewService(cli client.Client, cm *credentials.CredentialManager, gitChecker gitAutoCommitter, syncer workflowSyncer, logger logrus.FieldLogger, env string) *Service {
	return &Service{
		Client:            cli,
		CredentialManager: cm,
		Logger:            logger,
		Environment:       env,
		Git:               gitChecker,
		Syncer:            syncer,
	}
}

// Watch starts monitoring workflow changes until the context is canceled.
func (s *Service) Watch(ctx context.Context, opts Options) error {
	if err := s.Client.HealthCheck(ctx); err != nil {
		return fmt.Errorf("failed to connect to n8n: %w", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(sigChan)
		close(sigChan)
	}()

	ticker := time.NewTicker(opts.Interval)
	defer ticker.Stop()

	state := make(map[string]time.Time)
	if err := s.initialSync(ctx, state); err != nil {
		s.Logger.WithError(err).Warn("initial sync failed")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-sigChan:
			return nil
		case <-ticker.C:
			if err := s.checkChanges(ctx, state, opts); err != nil {
				s.Logger.WithError(err).Error("failed to check changes")
			}
		}
	}
}

func (s *Service) initialSync(ctx context.Context, state map[string]time.Time) error {
	wfs, err := s.Client.GetWorkflows(ctx)
	if err != nil {
		return err
	}
	for _, wf := range wfs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		state[wf.ID] = wf.UpdatedAt
	}
	return nil
}

func (s *Service) checkChanges(ctx context.Context, state map[string]time.Time, opts Options) error {
	wfs, err := s.Client.GetWorkflows(ctx)
	if err != nil {
		return err
	}

	changed := false
	for _, wf := range wfs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		last, ok := state[wf.ID]
		if !ok {
			s.Logger.Infof("new workflow detected: %s", wf.Name)
			state[wf.ID] = wf.UpdatedAt
			changed = true
			continue
		}
		if wf.UpdatedAt.After(last) {
			s.Logger.Infof("workflow updated: %s", wf.Name)
			state[wf.ID] = wf.UpdatedAt
			changed = true
		}
	}

	if !changed {
		return nil
	}

	if opts.AutoSync && s.Syncer != nil {
		if err := s.Syncer.Sync(ctx, isync.Options{}); err != nil {
			s.Logger.WithError(err).Warn("auto-sync failed")
		}
	}

	if opts.AutoCommit && s.Git != nil {
		if msg, err := s.Git.AutoCommitWorkflows(ctx, ""); err != nil {
			s.Logger.WithError(err).Warn("auto-commit failed")
		} else if msg != "" {
			s.Logger.Info(msg)
		}
	}

	return nil
}
