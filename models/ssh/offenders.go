package models

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
	elastic "gopkg.in/olivere/elastic.v5"
)

type ByUnixTimeOffenders []attackerLoginAttemptDoc

func (a ByUnixTimeOffenders) Len() int      { return len(a) }
func (a ByUnixTimeOffenders) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByUnixTimeOffenders) Less(i, j int) bool {

	return a[i].Timestamp.Before(a[j].Timestamp)

}

func getOffenders(searchResult *elastic.SearchResult) (attempts []attackerLoginAttemptDoc, err error) {
	if searchResult.Hits.TotalHits > 0 {
		for _, hit := range searchResult.Hits.Hits {
			var t attackerLoginAttemptDoc
			err := json.Unmarshal(*hit.Source, &t)
			if err != nil {
				return nil, errors.New("No login attempts found, failed to unmarshal JSON")
			}
			attempts = append(attempts, t)
		}
		return attempts, nil

	}
	return nil, errors.New("No login attempts found")
}
func GetOffenders(ContainerId string, probeName string, timestamp string) (attempts []attackerLoginAttemptDoc, err error) {
	var e ElasticOutputClient
	e.url = beego.AppConfig.String("elasticsearchurl")
	err = e.Init()
	if err != nil {
		logrus.Fatalf("Unable to initialize ES client %s", err)
	}
	termQuery := elastic.NewTermQuery("containerid", ContainerId)
	probeQuery := elastic.NewTermQuery("probe_name", probeName)
	successQuery := elastic.NewTermQuery("successful", true)
	query := elastic.NewBoolQuery()
	d1 := elastic.NewRangeQuery("@timestamp")
	i, err := strconv.Atoi(timestamp)
	if err != nil {
		logrus.Fatal(err)
	}
	d1.From(i - 600000)
	d1.To(i + 600000)
	query = query.Must(termQuery).Must(probeQuery).Must(successQuery).Must(d1)
	searchResult, err := e.Search("ssh_login_attempts", 30, query)
	attempts, err = getOffenders(searchResult)
	return attempts, err
}
