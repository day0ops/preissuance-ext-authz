// Package logger provides a small wrapper around zap for constructing the
// service's structured logger.
package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// New returns a configured zap logger. When debug is true it returns a
// development logger (human-readable, debug level, stack traces on warn);
// otherwise it returns a production logger (JSON, info level).
func New(debug bool) (*zap.Logger, error) {
	if debug {
		log, err := zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("build development logger: %w", err)
		}
		return log, nil
	}

	log, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("build production logger: %w", err)
	}
	return log, nil
}
