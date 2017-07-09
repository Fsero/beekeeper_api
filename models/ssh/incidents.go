package models

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"regexp"

	"github.com/Sirupsen/logrus"
)

var (
	Incidents map[string]*Incident
)

type Provider struct {
	Name    string
	Country string
}

type Incident struct {
	ID         string
	StartedAt  time.Time
	FinishedAt time.Time
	Offenders  []attackerLoginAttemptDoc
	Activities []attackerActivityDoc
	Provider
}

func GetIncident(IncidentId string) (incident *Incident, err error) {
	return nil, nil
}

func GetAllIncidents() map[string]*Incident {
	var e ElasticOutputClient
	Incidents = make(map[string]*Incident)
	e.url = "http://main01.superprivyhosting.com:9200"
	err := e.Init()
	if err != nil {
		logrus.Fatalf("Unable to initializa ES client %s", err)
	}
	// Search with a term query
	//termQuery := elastic.NewTermQuery("containerid", ContainerId)
	//successQuery := elastic.NewTermQuery("successful", true)

	//dateQuery := elastic.NewTermQuery("@timestamp", timestamp)

	// d1 := elastic.NewRangeQuery("@timestamp")
	// i, err := strconv.Atoi(timestamp)
	// if err != nil {
	// 	logrus.Fatal(err)
	// }
	// d1.From(i - 600000)

	// d1.To(i + 600000)
	// //d1 = d1.Lte(timestamp)
	// d2 := elastic.NewBoolQuery()
	// d2.Must(successQuery)
	// d2.Must(d1)

	//d2 := elastic.NewRangeQuery("%s+10m")

	searchResult, err := e.client.Search().
		Index("alerts-*"). // search in index "twitter"
		//Query(termQuery).
		//Query(d2).
		//Query(d2).
		//Query(dateQuery).         // specify the query
		Sort("@timestamp", false). // sort by "user" field, ascending
		From(0).Size(100).         // take documents 0-9
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
	//fmt.Printf("Found a total of %d alertDoc \n", searchResult.TotalHits())

	// Here's how you iterate through the search results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		//	fmt.Printf("Found a total of %d attempts \n", searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t alertDoc
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				// Deserialization failed
			}

			containerIDRegexp := regexp.MustCompile(`(id=(\w+))`)
			matches := containerIDRegexp.FindStringSubmatch(t.Message)
			if len(matches) > 0 {
				containerID := matches[2]
				offenders, _ := GetOffenders(containerID, fmt.Sprintf("%d000", t.Timestamp.Unix()))
				activities, _ := GetActivities(containerID, fmt.Sprintf("%d000", t.Timestamp.Unix()))
				//fmt.Println("%s %s %+v %+v", containerID, fmt.Sprintf("%d000", t.Timestamp.Unix()), offenders, activities)
				inc := Incident{Activities: activities, ID: fmt.Sprintf("%d000", t.Timestamp.Unix()), StartedAt: t.Timestamp, Provider: Provider{}, Offenders: offenders}
				//fmt.Println("%+v", inc)
				Incidents[inc.ID] = &inc
			}
			// Work with tweet
			//fmt.Printf("Attempt by %s: %+v\n", t)
		}
	} else {
		// No hits
		//fmt.Print("Found no activities\n")
	}

	return Incidents
}
