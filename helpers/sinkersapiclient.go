package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
)

// "ProbeID": "kTuoo-GmR9VEFo7C4_f7t2YTkpv22OeWbE1dIg==",
// "fqdn": "srv02.superprivyhosting.com",
// "ipv4": "146.185.164.76",
// "ipv6": "",
// "provider": "Digital Ocean",
// "geolongitude": "4.954100",
// "geolatitude": "52.296398",
// "country": "Netherlands",
// "sshprivateKey": "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlFcEFJQkFBS0NBUUVBcVZHeFVtcStWa3Nrc052aUVneVpCZWx3cFFreW9kT2dqcU52UmEzUGxIU0QzSDZ6ClBFa2JFdlp4U2h1TUR4YTJMMUNvd3AyUFRtZERObXhLcDBFSGJleER2eXpqeWJoNUVQdEtDVzN1bVpIZHl6ckkKdjJwMEdMVHJBMm9uWjJ3bWdma3R3b0lwTEVRQ2JWd0NiMzlWUmxjNFpuZE1wK3g3MHFCTEpLZWN2ajRjNDZQYQpyS1NGR3BWbEE0OWp5UitQaEFKeDdIRzdmZW44aVVOVEJSUFNxbzFzdTlQanVZYnhJNmdTUzVrOU1FeUloZ0xoCk5Ub0VFSHNPQW9waGIrMGJZUi9PdmNWdGNiZk9iZEYzcGVrcUh3SlN3Y21WMGNuM2pKejl3a2pkM004d3JSelMKYkx1UTRUMmZEOVlwTXArWURsTnV4SDYxYkhaTjFhT3Nmc29RbXdJREFRQUJBb0lCQUdOdkJUcUlVMFRzRmgwQwoxQzJUVmw0aGJEU1BSVHZCd3oxZy8xeWxLUTFlcTMxV3JyMk5sU2U2c0djNzdER0VQZk5sWStYK2o0VVVvV1VaCmpYSFJmNkp2S0kzaHQ5Zmp6TDFMUlh4cUliL3Y2SmVMNnc1MlhyMlBxYUEwdS9WQmp3K3ZITTlvanZOZThTbjEKNmJ3K3cvNXVCRUl3ZDBUNlhQRWhqMTkwUzFsd3ZZRmhaWHpRNnFjT1Fyd1ZadUJUUVRsZUFFbWIxZlVkRHo0OAp4QUUwL1R2bVJiczdlM1l6Z1lhU0k4RWJuWnFZNm1XSFN3MzJ5dHdhemVjTUNtSzFNMXdRQjJjK1lwSGhTdExvCjFxWlFjSGZwRnVQZ1FYYTB5WGQrL0R4OWQwMFByUEg5aHIrWEpPUWg3T2NPVk41K2pWTmtjTkV4SDViQ2hWOG4KQzZJdzVERUNnWUVBM3J4Ym8yb0gwUEdlWEViMjUvVlpIeTB2bVJvM0VzL3d6ZFFpYWgwQXA0aDRCL0NlUTRtbQpSWStXUHhKNEFDaWdnU3BsTEorTXhPZXczVXFPRVNSbThvUysxcWI2OWM2YTdUd0VxZWFGQ2cxVmN1ZjBEL25qCjJnSTd2NnVJMHJLcG0ybjdFd0VYbkNEaGZvL2FncUd6YzBWa2lhWitMNytIRkZIazltTlJWQ01DZ1lFQXdwc2IKWkVIQVdtcTNjajMyTGFqVTRlZUZYV1BmSml6enhiZDlaRmpZbVNhNnkzYm4zU3ZFWTdOak4vQ2paaFVMZmhSTApBL24yd1RMZ09HaXZobFdTR3YyVEJxSmVrUDNFc2R4RTlna0xMTWNscDhxZTgyd0k4U2NPSWZNUFIyREd5b1lYCkhzYm9ZTG1WUnZrUGVzUTVNTG1SemZPK3JmSHYrYlZCdFpZLy9Ta0NnWUVBb1N4TloxZHZabnVnYXdlUzNOQ0YKOS9CYmsvOExRUnFsRmx2ZHQwbGJVdCtHYzhCaTFWNUNxZTAzL0ZYaDdjTjRPVjh6TFBJYkM5VFgvNWxXYWdNYwpWM3RGR05CbG92OG96bWZ5dS9xcDVGYzNzTmsxbTJYb3diV0NCTFVjWWRLVXRuZ2ZEV1pwN2psQTByTkhtK1ZrCmxCSHZxVWVINGdkR3VLWjE5dkJ4Umw4Q2dZRUFoSkd2UGtRQWFsZktjanU5aVdzNjRrMmFyMzBLbGZJSGVvZysKRm03ajFxam9sUlNDYlV1VWRLck9pMXdWbzhQd1dVb3Z0QnpET09lVWtUalhZYmJIV2pXbHc5NDJkNlU0S2tXNApnTGEyY3lHVENGUGlwa2JSYko1RFpXTXo1RmNMOVFrVmxQVEJkcXJXQTB4RmZFZFNBbHhYOUNuNG1ueDNFdStrClBMU0hFTWtDZ1lCcmswUjhNeEFMQVh5MVZUNHA1dUI3MTZmeXdBYkNRZGVvMGlXMitzVDhMQVM3YkM5eWVJMmgKZS9zcDdOVXBTc3hOZnFBZ0d4Uk9YZkxPa2JJa2psaEZTLzhuMHZiY0tMR3g1UW1mbmMwaTJtRGh3VTNoZjNZVgpXUXNiemI1U01SVVI0dWZ4U054OHNVZW9PbUNXbmtDUm5oWGtzYmQ1dXRxeWd3eEJ4SGg2ZGc9PQotLS0tLUVORCBSU0EgUFJJVkFURSBLRVktLS0tLQo=",
// "sshpublicKey": "c3NoLXJzYSBBQUFBQjNOemFDMXljMkVBQUFBREFRQUJBQUFCQVFDcFViRlNhcjVXU3lTdzIrSVNESmtGNlhDbENUS2gwNkNPbzI5RnJjK1VkSVBjZnJNOFNSc1M5bkZLRzR3UEZyWXZVS2pDblk5T1owTTJiRXFuUVFkdDdFTy9MT1BKdUhrUSswb0piZTZaa2QzTE9zaS9hblFZdE9zRGFpZG5iQ2FCK1MzQ2dpa3NSQUp0WEFKdmYxVkdWemhtZDB5bjdIdlNvRXNrcDV5K1BoempvOXFzcElVYWxXVURqMlBKSDQrRUFuSHNjYnQ5NmZ5SlExTUZFOUtxald5NzArTzVodkVqcUJKTG1UMHdUSWlHQXVFMU9nUVFldzRDaW1GdjdSdGhIODY5eFcxeHQ4NXQwWGVsNlNvZkFsTEJ5WlhSeWZlTW5QM0NTTjNjenpDdEhOSnN1NURoUFo4UDFpa3luNWdPVTI3RWZyVnNkazNWbzZ4K3loQ2IgYW5zaWJsZS1nZW5lcmF0ZWQgb24gc3J2MDIuc3VwZXJwcml2eWhvc3RpbmcuY29tCg==",
// "tracespath": "/var/log/traces",
// "enabled": true,
// "created_at": "2017-06-16T11:45:10.735914298Z",
// "updated_at": "2017-06-16T11:45:10.735914388Z"

type Probe struct {
	ProbeID       string    `json:"ProbeID"`
	FQDN          string    `json:"fqdn"`
	IPv4          string    `json:"ipv4"`
	IPv6          string    `json:"ipv6"`
	Provider      string    `json:"provider"`
	Geolongitude  string    `json:"geolongitude"`
	Geolatitude   string    `json:"geolatitude"`
	Country       string    `json:"country"`
	SSHprivatekey string    `json:"sshprivatekey"`
	SSHpublickey  string    `json:"sshpublickey"`
	Tracespath    string    `json:"tracespath"`
	Enabled       bool      `json:"enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     time.Time `json:"deleted_at"`
}

func GetProbe(ProbeName string) ([]Probe, error) {

	sinkersurl := beego.AppConfig.String("sinkersapiurl")
	url := fmt.Sprintf("%s/v1/probe/name/%s", sinkersurl, ProbeName)
	fmt.Printf("querying %s for info about %s", url, ProbeName)
	// Build the request
	logrus.Infof("[GetProbe] requesting Probe info %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logrus.Fatal("NewRequest: ", err)
		return nil, err
	}

	// For control over HTTP client headers,
	// redirect policy, and other settings,
	// create a Client
	// A Client is an HTTP client
	client := &http.Client{}

	// Send the request via a client
	// Do sends an HTTP request and
	// returns an HTTP response
	resp, err := client.Do(req)
	if err != nil {
		logrus.Fatal("Do: ", err)
		return nil, err
	}

	// Callers should close resp.Body
	// when done reading from it
	// Defer the closing of the body
	defer resp.Body.Close()

	// Fill the record with the data from the JSON
	if resp.StatusCode != 200 {
		logrus.Fatalf("We were unable to retrieve provider info, request %s answered %s code", url, resp.Status)
		return nil, fmt.Errorf("unable to get Probe info from sinkers API")
	}

	probes := make([]Probe, 0)

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&probes); err != nil {
		logrus.Fatalf("Unable to get Probe info. sinker API is maybe down? error: %s ", err)
	}

	return probes, nil

}
