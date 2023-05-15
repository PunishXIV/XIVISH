package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/punishxiv/xivish"
)

func main() {
	s := xivish.NewScraper()
	s.AttachLogger(xivish.Logger{
		DisableAll: false,
		Queue: xivish.LoggerQueueCfg{
			Enabled:  true,
			Interval: 250 * time.Millisecond,
		},
	})

	err := s.AttachQueue(1, 1000)
	if err != nil {
		log.Fatal(err)
	}

	c := s.AttachCharacterCollector(func(character xivish.Character) {
		charJson, err := json.MarshalIndent(character, "", "  ")
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("character: %s", charJson)
	})

	_ = c.AddCharacterToQueue("na", 35383842)

	err = s.StartQueue()
	if err != nil {
		log.Fatal(err)
	}

	s.Close()
}
