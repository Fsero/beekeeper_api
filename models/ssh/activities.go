package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	elastic "gopkg.in/olivere/elastic.v5"
)

type ByUnixTimeActivities []attackerActivityDoc

func (a ByUnixTimeActivities) Len() int      { return len(a) }
func (a ByUnixTimeActivities) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByUnixTimeActivities) Less(i, j int) bool {

	return a[i].Timestamp.Before(a[j].Timestamp)

}

func getActivities(searchResult *elastic.SearchResult) (activities []attackerActivityDoc, err error) {
	// searchResult is of type SearchResult and returns hits, suggestions,
	// and all kinds of other information from Elasticsearch.
	//fmt.Printf("Query took %d milliseconds\n", searchResult.TookInMillis)

	// Each is a convenience function that iterates over hits in a search result.
	// It makes sure you don't need to check for nil values in the response.
	// However, it ignores errors in serialization. If you want full control
	// over the process, see below.
	// var ttyp attackerActivityDoc
	// for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
	// 	t := item.(attackerActivityDoc)
	// 	fmt.Printf("Activity by %+v \n", t)
	// }

	// TotalHits is another convenience function that works even when something goes wrong.
	//fmt.Printf("Found a total of %d activities\n", searchResult.TotalHits())

	// Here's how you iterate through the search results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		//	fmt.Printf("Found a total of %d activities \n", searchResult.Hits.TotalHits)

		// Iterate through results
		for _, hit := range searchResult.Hits.Hits {
			// hit.Index contains the name of the index

			// Deserialize hit.Source into a Tweet (could also be just a map[string]interface{}).
			var t attackerActivityDoc
			err := json.Unmarshal(*hit.Source, &t)
			activities = append(activities, t)

			if err != nil {
				// Deserialization failed
			}

			// Work with tweet
			//		fmt.Printf("Activity by %s: %+v\n", t)
		}
		return activities, nil

	} else {
		// No hits
		fmt.Print("Found no activities\n")
	}

	return nil, errors.New("ObjectId Not Exist")
}

func GetActivities(ContainerId string, probeName string, timestamp string) (activities []attackerActivityDoc, err error) {
	var e ElasticOutputClient
	e.url = beego.AppConfig.String("elasticsearchurl")
	err = e.Init()
	if err != nil {
		logrus.Fatalf("Unable to initializa ES client %s", err)
	}
	// Search with a term query
	termQuery := elastic.NewTermQuery("containerid", ContainerId)
	probeQuery := elastic.NewTermQuery("probe_name", probeName)

	query := elastic.NewBoolQuery()

	d1 := elastic.NewRangeQuery("@timestamp")
	i, err := strconv.Atoi(timestamp)
	if err != nil {
		logrus.Fatal(err)
	}
	d1.From(i)
	d1.To(i + 600000)
	query = query.Must(termQuery).Must(probeQuery).Must(d1)
	searchResult, err := e.Search("ssh_activities", 30, query)
	activities, err = getActivities(searchResult)
	return activities, err
}
