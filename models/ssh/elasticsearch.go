package models

import (
	"time"

	"github.com/Sirupsen/logrus"

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
