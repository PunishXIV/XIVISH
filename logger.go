package xivish

import (
	"fmt"
	"time"
)

type LoggerQueueCfg struct {
	Enabled  bool
	Interval time.Duration
	stop     chan bool
}

type Logger struct {
	DisableAll bool
	Queue      LoggerQueueCfg
}

// AttachLogger attaches a logger to the scraper
func (s *Scraper) AttachLogger(logger Logger) {
	s.logger = logger
}

// log logs a message with a prefix and status
func (l *Logger) log(status string, prefix string, msg any) {
	if !l.DisableAll {
		fmt.Printf("[%s] %s: %v\n", status, prefix, msg)
	}
}

// LogInfo logs an info message
func (l *Logger) LogInfo(prefix string, msg any) {
	l.log("INFO", prefix, msg)
}

// LogError logs an error message
func (l *Logger) LogError(prefix string, msg any) {
	l.log("ERROR", prefix, msg)
}

// DisableAll disables all logging
func (l *Logger) Disable() {
	l.DisableAll = true

	l.Close()
}

// Close stops the logging
func (l *Logger) Close() {
	if l.Queue.stop != nil {
		l.Queue.stop <- true
	}
}

// SetQueueState sets whether the queue progress should be logged and the interval
func (l *Logger) SetQueueState(enabled bool, interval *time.Duration) {
	l.Queue.Enabled = enabled
	if interval != nil {
		l.Queue.Interval = *interval
	}

	if !enabled && l.Queue.stop != nil {
		l.Queue.stop <- true
	}
}
