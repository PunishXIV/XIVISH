package xivish

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
)

var (
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36"
)

type Scraper struct {
	collector *colly.Collector
	queue     *queue.Queue
	logger    Logger
}

// NewScraper creates a new scraper
func NewScraper() Scraper {
	return Scraper{
		collector: nil,
		queue:     nil,
		logger:    Logger{DisableAll: true},
	}
}

// createBaseCollector creates a collector with some defaults
func (s *Scraper) createBaseCollector() {
	s.attachCollector(colly.NewCollector(
		colly.UserAgent(userAgent),
	))

	_ = s.SetLimiter(1*time.Second, 1)
}

// SetLimiter sets the rate limiter for the scraper
func (s Scraper) SetLimiter(duration time.Duration, parallel int) error {
	s.logger.LogInfo("limiter", fmt.Sprintf("limiting to %d request per %s", parallel, duration))
	err := s.collector.Limit(&colly.LimitRule{DomainGlob: "*", Delay: duration, Parallelism: parallel})

	if err != nil {
		s.logger.LogError("limiter", err)
		return err
	}
	return nil
}

// attachCollector attaches a collector to the scraper
func (s *Scraper) attachCollector(c *colly.Collector) {
	s.collector = c
}

// Close closes the scraper
func (s Scraper) Close() {
	s.collector.Wait()
	s.logger.LogInfo("scraper status", "closed")
	s.logger.Close()
}

// Restart stops the queue and starts it again after a delay
func (s Scraper) Restart(delay time.Duration) error {
	s.queue.Stop()
	time.Sleep(delay)
	err := s.queue.Run(s.collector)
	if err != nil {
		s.logger.LogError("queue", err)
		return err
	}
	return nil
}
