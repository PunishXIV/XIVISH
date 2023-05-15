package xivish

import (
	"fmt"
	"time"

	"github.com/gocolly/colly/v2/queue"
)

// SetupQueue sets up the queue for the scraper
func (s *Scraper) AttachQueue(threads int, maxSize int) error {
	q, err := queue.New(
		threads,
		&queue.InMemoryQueueStorage{MaxSize: maxSize},
	)
	if err != nil {
		s.logger.LogError("queue", err)
		return err
	}

	s.queue = q
	s.logger.LogInfo("queue created", fmt.Sprintf("%d threads and max size of %d", threads, maxSize))
	return nil
}

// addUrlToQueue adds an url to the queue
func (s Scraper) addUrlToQueue(url string) error {
	if s.queue == nil {
		s.logger.LogError("scraper", "no queue available")
		return fmt.Errorf("no queue available")
	}
	err := s.queue.AddURL(url)
	if err != nil {
		s.logger.LogError("queue", err)
		return err
	}
	return nil
}

// StartQueue initiates scraping with the queue
func (s Scraper) StartQueue() error {
	if s.collector == nil {
		s.logger.LogError("scraper", "no collector available")
		return fmt.Errorf("no collector available")
	}
	if s.queue == nil {
		s.logger.LogError("scraper", "no queue available")
		return fmt.Errorf("no queue available")
	}

	size, err := s.queue.Size()
	if err != nil {
		s.logger.LogError("queue size", err)
		return err
	}
	s.logger.LogInfo("queue status", fmt.Sprintf("started with %d items", size))

	s.startQueueLogger()

	err = s.queue.Run(s.collector)
	if err != nil {
		s.logger.LogError("queue", err)
		return err
	}

	s.logger.LogInfo("queue status", "finished")
	return nil
}

// startQueueLogger starts the queue logger in the background
func (s *Scraper) startQueueLogger() {
	if (s.logger.Queue == LoggerQueueCfg{}) || !s.logger.Queue.Enabled {
		return
	}

	s.logger.Queue.stop = make(chan bool)

	go func() {
	loop:
		for {
			select {
			case <-s.logger.Queue.stop:
				break loop
			default:
				time.Sleep(s.logger.Queue.Interval)
				s.logQueueProgress()
			}
		}
	}()
}

// logQueueProgress logs the queue progress
func (s Scraper) logQueueProgress() {
	size, err := s.queue.Size()
	if err != nil {
		s.logger.LogError("queue size", err)
	}

	if size > 0 {
		s.logger.LogInfo("queue size", size)
	}
}
