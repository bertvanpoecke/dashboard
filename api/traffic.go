package api

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	gofeed "github.com/mmcdole/gofeed"
	"github.com/sirupsen/logrus"
)

const (
	debug        = true
	fetchTimeout = 30 * time.Second

	trafficXmlURL        = "http://www.verkeerscentrum.be/rss/4-INC%7CLOS%7CINF%7CPEVT.xml"
	trafficSummaryString = "Totale filelengte op snelwegen in Vlaanderen: "
	trendString          = " Trend: "
)

type Incident struct {
	Road        string `json:"road,omitempty"`
	Direction   string `json:"direction,omitempty"`
	Location    string `json:"location,omitempty"`
	Description string `json:"description,omitempty"`
	Link        string `json:"link,omitempty"`
}

type TrafficInfo struct {
	Timestamp        time.Time   `json:"timestamp,omitempty"`
	TotalJamLengthKm string      `json:"jamlength,omitempty"`
	Trend            string      `json:"trend,omitempty"`
	Incidents        []*Incident `json:"incidents,omitempty"`
}

type TrafficInfoGetter struct {
	TrafficInfo chan *TrafficInfo
	Active      bool
}

func NewTrafficInfoGetter() *TrafficInfoGetter {
	return &TrafficInfoGetter{
		TrafficInfo: make(chan *TrafficInfo),
		Active:      false,
	}
}

func (t *TrafficInfoGetter) Start() {
	t.Active = true
	go t.fetchTraffic()
	logrus.Info("Traffic Info Getter started")
}

func (t *TrafficInfoGetter) Stop() {
	t.Active = false
	close(t.TrafficInfo)
	logrus.Info("Traffic Info Getter stopped")
}

func (t *TrafficInfoGetter) fetchTraffic() {
	for {
		if !t.Active {
			return
		}

		parser := gofeed.NewParser()
		feed, err := parser.ParseURL(trafficXmlURL)
		if err != nil {
			logrus.Warningf("Unable to parse traffic feed (%v)", err)
		}

		totalJamLengthKm := ""
		jamTrend := ""
		incidents := []*Incident{}
		for i, item := range feed.Items {
			if i == 0 && strings.HasPrefix(item.Title, trafficSummaryString) {
				totalJamLengthKm, jamTrend, err = parseGeneralInfo(item)
				if err != nil {
					logrus.Warningf("%v", err)
				}
				continue
			}

			incident, err := parseItem(item)
			if err != nil {
				logrus.Warningf("%v", err)
				continue
			}
			incidents = append(incidents, &incident)
		}

		trafficInfo := &TrafficInfo{
			Timestamp:        time.Now(),
			TotalJamLengthKm: totalJamLengthKm,
			Trend:            jamTrend,
			Incidents:        incidents,
		}

		if debug {
			trafficstring, err := json.Marshal(trafficInfo)
			if err != nil {
				logrus.Debugf("%v", err)
			}
			logrus.Info(string(trafficstring))
		} else {
			t.TrafficInfo <- trafficInfo
		}

		time.Sleep(fetchTimeout)
	}
}

func replaceArrow(s string) string {
	return strings.ReplaceAll(s, "->", "\u2192")
}

func parseItem(item *gofeed.Item) (Incident, error) {
	incident := Incident{}
	chuncks := strings.Split(item.Title, ":")
	if len(chuncks) != 3 && len(chuncks) != 4 {
		return incident, fmt.Errorf("Unable to parse incident title. Incident skipped.")
	}
	incident.Road = strings.TrimSpace(chuncks[0])
	incident.Direction = replaceArrow(strings.TrimSpace(chuncks[1]))
	incident.Link = strings.TrimSpace(item.Link)
	if len(chuncks) == 3 {
		incident.Description = strings.TrimSpace(chuncks[2])
	} else { // length = 4
		incident.Location = strings.TrimSpace(chuncks[2])
		incident.Description = strings.TrimSpace(chuncks[3])
	}
	return incident, nil
}

func parseGeneralInfo(item *gofeed.Item) (string, string, error) {
	if strings.HasPrefix(item.Title, trafficSummaryString) {
		chuncks := strings.Split(item.Title, ";")
		if len(chuncks) != 2 {
			return "", "", fmt.Errorf("Unable to parse general info.")
		}
		jamLength := strings.TrimSpace(strings.TrimPrefix(chuncks[0], trafficSummaryString))
		jamTrend := strings.TrimSpace(strings.TrimPrefix(chuncks[1], trendString))
		return jamLength, jamTrend, nil
	}
	return "", "", fmt.Errorf("No general info item.")
}
