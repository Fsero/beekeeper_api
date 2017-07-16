package models

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"bitbucket.org/fseros/beekeeper_api/helpers"

	"regexp"

	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
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

func GetIncident(IncidentId string) (incident *Incident, err error) {
	return nil, nil
}

func GetAllIncidents() map[string]*Incident {
	var e ElasticOutputClient
	Incidents = make(map[string]*Incident)
	e.url = beego.AppConfig.String("elasticsearchurl")
	err := e.Init()
	if err != nil {
		logrus.Fatalf("Unable to initialize ES client %s", err)
	}
	searchResult, err := e.client.Search().
		Index("alerts-*").         // search in index "twitter"
		Sort("@timestamp", false). // sort by "user" field, ascending
		From(0).Size(30).          // take documents 0-9
		Pretty(true).              // pretty print request and response JSON
		Do(context.Background())   // execute
	if err != nil {
		// Handle error
		panic(err)
	}

	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	//fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over the process, see below.
	// var ttyp alertDoc
	// for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
	// 	t := item.(alertDoc)
	// 	//fmt.Printf("alertDoc by %+v \n", t)
	// }

	// TotalHits is another convenience function that works even when something goes wrong.
	fmt.Printf("Found a total of %d alertDoc \n", searchResult.TotalHits())

	// Here's how you iterate through the search results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d attempts \n", searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t alertDoc
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}
			if t.Message == "" {
				//fmt.Printf("empty alert, why? %+v", t)
				continue
			}
			//fmt.Printf("Found alert!,  %+v \n", t)
			slices := strings.Split(t.Raw, "falco: ")

			input := slices[1]
			//input := strings.Replace(slices[1], `\"\"`, `''`, -1)
			//fmt.Println(input)
			//input = strings.Replace(input, `"`, `|`, -1)
			//input = strings.Replace(input, `|`, `"`, -1)

			var alert Alert
			if err = json.Unmarshal([]byte(input), &alert); err != nil {

				fmt.Printf("something went wrong %s %s\n", input, err)

			}

			containerIDRegexp := regexp.MustCompile(`(id=(\w+))`)
			matches := containerIDRegexp.FindStringSubmatch(t.Message)
			if len(matches) > 0 {
				containerID := matches[2]
				stringSlice := strings.Split(t.Source, "/")
				probeName := fmt.Sprintf("%s.superprivyhosting.com", stringSlice[3])
				//fmt.Printf("getting info for %s %+v \n", probeName, stringSlice)

				probes, err := helpers.GetProbe(probeName)
				probe := probes[0]

				offenders, _ := GetOffenders(containerID, probe.FQDN, fmt.Sprintf("%d000", alert.Timestamp.Unix()))
				activities, _ := GetActivities(containerID, probe.FQDN, fmt.Sprintf("%d000", alert.Timestamp.Unix()))

				sort.Sort(ByUnixTimeActivities(activities))
				sort.Sort(ByUnixTimeOffenders(offenders))
				fmt.Printf("Alert %+v \n", alert)
				//started_at and finished_at should be obtained from offenders and activities and not the alert.

				//fmt.Println("%s %s %+v %+v", containerID, fmt.Sprintf("%d000", t.Timestamp.Unix()), offenders, activities)
				var lastSeen time.Time
				if len(offenders) > 0 {
					lastSeen = offenders[len(offenders)-1].Timestamp
				} else {
					lastSeen = time.Time{}
				}

				if err == nil {
					inc := Incident{Activities: activities, Triggered: t.Message, ID: fmt.Sprintf("%s-%d000", probe.Provider, alert.Timestamp.Unix()), StartedAt: alert.Timestamp, FinishedAt: lastSeen, Provider: Provider{Provider: probe.Provider, Country: probe.Country}, Offenders: offenders}
					//fmt.Println("%+v", inc)
					Incidents[inc.ID] = &inc
				}
			}
			// Work with tweet
			//fmt.Printf("Attempt by %s: %+v\n", t)
		}
	} else {
		// No hits
		//fmt.Print("Found no activities\n")
	}

	fmt.Printf("%d incidents found \n", len(Incidents))

	return Incidents
}
