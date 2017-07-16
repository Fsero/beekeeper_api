package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"bitbucket.org/fseros/beekeeper_api/helpers"

	"regexp"

	"github.com/Sirupsen/logrus"
	"gopkg.in/olivere/elastic.v5"
)

var (
	Incidents map[string]*Incident
)

type Provider struct {
	Provider string
	Country  string
}

type Incident struct {
	ID        string
	Triggered string
	Provider
	StartedAt  time.Time
	FinishedAt time.Time
	Offenders  []attackerLoginAttemptDoc
	Activities []attackerActivityDoc
}

// Jul 15 20:59:00 srv01 falco: {"output":"20:48:49.599471717: Alert Shell spawned in a container other than entrypoint (user=root ssh (id=94cf593573b3) ssh (id=94cf593573b3) shell=bash parent=sshd cmdline=bash -c /usr/lib/openssh/sftp-server)","priority":"Alert","rule":"Run shell in container","time":"2017-07-15T20:48:49.599471717Z"}

type Alert struct {
	Output    string    `json:"output"`
	Priority  string    `json:"priority"`
	Rule      string    `json:"rule"`
	Timestamp time.Time `json:"time"`
}

func GetIncident(from time.Time, to time.Time) (map[string]*Incident, error) {
	return nil, nil
}

func getIncidents(searchResult *elastic.SearchResult) (map[string]*Incident, error) {
	Incidents = make(map[string]*Incident, 0)
	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var t alertDoc
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				logrus.Fatalf("[getIncidents] Unable to unmarshall json %s %s", *hit.Source, err)
			}
			if t.Message == "" {
				logrus.Debugf("[getIncident] empty alert?")
				continue
			}
			slices := strings.Split(t.Raw, "falco: ")
			input := slices[1]
			var alert Alert
			if err = json.Unmarshal([]byte(input), &alert); err != nil {
				logrus.Errorf("[getIncidents] Unable to unmarshall json %s %s", input, err)
				return nil, errors.New("unable to get alerts")

			}
			containerIDRegexp := regexp.MustCompile(`(id=(\w+))`)
			matches := containerIDRegexp.FindStringSubmatch(t.Message)
			if len(matches) > 0 {
				containerID := matches[2]
				stringSlice := strings.Split(t.Source, "/")
				probeName := fmt.Sprintf("%s.superprivyhosting.com", stringSlice[3])
				probes, err := helpers.GetProbe(probeName)
				var probe helpers.Probe
				if len(probes) > 0 {
					probe = probes[0]
				} else {
					logrus.Errorf("[GetIncidents] Unable to get probe info")
					return nil, errors.New("unable to get probe info")
				}
				offenders, _ := GetOffenders(containerID, probe.FQDN, fmt.Sprintf("%d000", alert.Timestamp.Unix()))
				activities, _ := GetActivities(containerID, probe.FQDN, fmt.Sprintf("%d000", alert.Timestamp.Unix()))
				sort.Sort(ByUnixTimeActivities(activities))
				sort.Sort(ByUnixTimeOffenders(offenders))
				logrus.Debugf("[getIncidents] new alert found! %+v", alert)
				var lastSeen time.Time
				if len(offenders) > 0 {
					lastSeen = offenders[len(offenders)-1].Timestamp
				} else {
					lastSeen = time.Time{}
				}

				if err == nil {
					inc := Incident{Activities: activities, Triggered: t.Message, ID: fmt.Sprintf("%s-%d000", probe.Provider, alert.Timestamp.Unix()), StartedAt: alert.Timestamp, FinishedAt: lastSeen, Provider: Provider{Provider: probe.Provider, Country: probe.Country}, Offenders: offenders}
					logrus.Debugf("[getIncidents] adding incident %+v", inc)
					Incidents[inc.ID] = &inc
				}
			}
		}
	}
	return Incidents, nil
}

func GetAllIncidents() (map[string]*Incident, error) {
	var e ElasticOutputClient
	searchResult, err := e.Search("alerts-*", 30, nil)
	Incidents, err := getIncidents(searchResult)
	logrus.Debugf("[GetAllIncidents] %d incidents found \n", len(Incidents))
	return Incidents, err
}
