package services

import (
	"context"
	"sync"
	"time"

	"github.com/lielamurs/aggregator/internal/config"
	"github.com/sirupsen/logrus"
)

type SubmissionProcessor struct {
	submissionService SubmissionService
	config            config.SubmissionProcessorConfig
	logger            *logrus.Logger
	ctx               context.Context
	cancel            context.CancelFunc
	wg                sync.WaitGroup
	running           bool
	mu                sync.RWMutex
}

func NewSubmissionProcessor(
	submissionService SubmissionService,
	config config.SubmissionProcessorConfig,
	logger *logrus.Logger,
) *SubmissionProcessor {
	return &SubmissionProcessor{
		submissionService: submissionService,
		config:            config,
		logger:            logger,
	}
}

func (p *SubmissionProcessor) Start() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return nil
	}

	p.logger.WithField("interval_seconds", p.config.IntervalSeconds).Info("Starting submission processor")

	p.ctx, p.cancel = context.WithCancel(context.Background())
	p.running = true

	p.wg.Add(1)
	go p.run()

	return nil
}

func (p *SubmissionProcessor) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.running {
		return nil
	}

	p.logger.Info("Stopping submission processor")

	p.cancel()
	p.wg.Wait()
	p.running = false

	p.logger.Info("Submission processor stopped")
	return nil
}

func (p *SubmissionProcessor) run() {
	defer p.wg.Done()

	logger := p.logger.WithField("component", "submission_processor")
	interval := time.Duration(p.config.IntervalSeconds) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	logger.Info("Submission processor started")

	p.processSubmissions(logger)

	for {
		select {
		case <-p.ctx.Done():
			logger.Info("Submission processor context cancelled")
			return
		case <-ticker.C:
			p.processSubmissions(logger)
		}
	}
}

func (p *SubmissionProcessor) processSubmissions(logger *logrus.Entry) {
	startTime := time.Now()
	logger.Debug("Starting submission processing cycle")

	if err := p.submissionService.ProcessSubmissions(p.ctx); err != nil {
		logger.WithError(err).Error("Submission processing cycle failed")
	} else {
		duration := time.Since(startTime)
		logger.WithField("duration", duration).Debug("Submission processing cycle completed")
	}
}
