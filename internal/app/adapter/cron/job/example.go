package job

import "log/slog"

// ExampleJob example job
type ExampleJob struct {
	logger *slog.Logger
}

// NewExampleJob build example job
func NewExampleJob(logger *slog.Logger) *ExampleJob {
	return &ExampleJob{
		logger: logger,
	}
}

// Run execute job
func (s ExampleJob) Run() {
	s.logger.Info("example Job executed successfully")
}
