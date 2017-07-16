package models

import (
	"context"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"

	elastic "gopkg.in/olivere/elastic.v5"
)

type attackerLoginAttemptDoc struct {
	//ID          string           `json:"id"`
	Timestamp   time.Time        `json:"@timestamp"`
	ContainerID string           `json:"containerid"`
	IP          string           `json:"ip"`
	Country     string           `json:"country"`
	Location    elastic.GeoPoint `json:"location"`
	User        string           `json:"user"`
	Password    string           `json:"password"`
	Successful  bool             `json:"successful"`
	//ProbeIP       string           `json:"probe_ip"`
	//ProbeName     string           `json:"probe_name"`
	//ProbeProvider string           `json:"probe_provider"`
	//ProbeLocation elastic.GeoPoint `json:"probe_provider_location"`
}

type attackerActivityDoc struct {
	//ID          string    `json:"id"`
	Timestamp   time.Time `json:"@timestamp"`
	ContainerID string    `json:"containerid"`
	PID         string    `json:"pid"`
	User        string    `json:"user"`
	//SourceFile  string    `json:"source"`
	Activity string `json:"activity"`
	//ProbeIP       string           `json:"probe_ip"`
	//ProbeName     string           `json:"probe_name"`
	//ProbeProvider string           `json:"probe_provider"`
	//ProbeLocation elastic.GeoPoint `json:"probe_provider_location"`
}

type alertDoc struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"@timestamp"`
	Message   string    `json:"msg"`
	Raw       string    `json:"message"`
	Rule      string    `json:"rule"`
	Source    string    `json:"source"`
}

type ElasticOutputClient struct {
	client   *elastic.Client
	url      string
	sniff    bool
	n        int
	bulkSize int
}

func (e *ElasticOutputClient) Init() error {
	logrus.Debug("initialize elasticsearch output")
	client, err := elastic.NewClient(elastic.SetURL(e.url), elastic.SetSniff(e.sniff))
	if err != nil {
		return err
	}
	e.client = client
	return nil
}

func (e *ElasticOutputClient) Search(index string, pageSize int, query *elastic.BoolQuery) (*elastic.SearchResult, error) {
	var searchResult *elastic.SearchResult
	var err error
	e.url = beego.AppConfig.String("elasticsearchurl")
	err = e.Init()
	if err != nil {
		logrus.Fatalf("Unable to initialize ES client %s", err)
	}
	if query == nil {
		searchResult, err = e.client.Search().
			Index(index).              // search in index "twitter"
			Sort("@timestamp", false). // sort by "user" field, ascending
			From(0).Size(pageSize).    // take documents 0-9
			Pretty(true).              // pretty print request and response JSON
			Do(context.Background())   // execute
	} else {
		searchResult, err = e.client.Search().
			Index(index). // search in index "twitter"
			Query(query).
			Sort("@timestamp", false). // sort by "user" field, ascending
			From(0).Size(pageSize).    // take documents 0-9
			Pretty(true).              // pretty print request and response JSON
			Do(context.Background())   // execute
	}

	if err != nil {
		logrus.Errorf("[Search] unable to search in Elasticsearch %s", err)
		return nil, err
	}
	return searchResult, err
}
