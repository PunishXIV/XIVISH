package xivish

import (
	"fmt"
	"regexp"
	"strconv"

	"strings"

	"github.com/gocolly/colly/v2"
)

var (
	characterURL       = "https://%s.finalfantasyxiv.com/lodestone/character/%d"
	characterSelectors = map[string]string{
		"contents":    "div#character",
		"pathName":    "a.frame__chara__link",
		"idRegExp":    `(\d+)\/?$`,
		"avatar":      "div.frame__chara__face > img",
		"name":        "div.frame__chara__box > p.frame__chara__name",
		"title":       "div.frame__chara__box > p.frame__chara__title",
		"bio":         "div.character__selfintroduction",
		"world":       "div.frame__chara__box > p.frame__chara__world",
		"worldRegExp": `(\w+) ?\[(\w+)\]`,
	}
)

type Character struct {
	Id         int
	Avatar     string
	Name       string
	Title      string
	Bio        string
	Server     string
	DataCenter string
}

type CharacterScraper struct {
	scraper *Scraper
}

// AttachCharacterCollector attaches a character collector to the scraper
func (s *Scraper) AttachCharacterCollector(cb func(character Character)) CharacterScraper {
	s.createBaseCollector()

	s.collector.OnHTML(characterSelectors["contents"], func(e *colly.HTMLElement) {
		pathName := e.ChildAttr(characterSelectors["pathName"], "href")
		id := regexp.MustCompile(characterSelectors["idRegExp"]).FindStringSubmatch(pathName)[1]
		idAsInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			s.logger.LogError("parsing character id", err)
		}

		avatar := e.ChildAttr(characterSelectors["avatar"], "src")
		name := e.ChildText(characterSelectors["name"])
		title := e.ChildText(characterSelectors["title"])
		bio, _ := e.DOM.Find(characterSelectors["bio"]).Html()
		bio = strings.Replace(bio, "<br/>", "\n", -1)

		world := e.ChildText(characterSelectors["world"])
		worldRegexp := regexp.MustCompile(characterSelectors["worldRegExp"]).FindStringSubmatch(world)
		server := worldRegexp[1]
		dataCenter := worldRegexp[2]

		character := Character{
			Id:         int(idAsInt),
			Avatar:     avatar,
			Name:       name,
			Title:      title,
			Bio:        bio,
			Server:     server,
			DataCenter: dataCenter,
		}

		cb(character)
	})

	return CharacterScraper{scraper: s}
}

// AddCharacterToQueue adds a character to the queue
func (s CharacterScraper) AddCharacterToQueue(region string, id int) error {
	return s.scraper.addUrlToQueue(fmt.Sprintf(characterURL, region, id))
}
