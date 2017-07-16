package controllers

import (
	"fmt"
	"log"
	"sort"

	"bitbucket.org/fseros/beekeeper_api/models/ssh"

	"time"

	"github.com/Sirupsen/logrus"
	"github.com/astaxie/beego"
)

// Operations about object
type IncidentController struct {
	beego.Controller
}

// @Title GetIncidents
// @Description get all objects
// @Success 200 {object} []models.Incident
// @Failure 500 "Unable to get incidents"
// @router / [get]
func (o *IncidentController) GetIncidents() {
	var start_date, end_date string
	var obs map[string]*models.Incident
	var err error
	var endDate, startDate time.Time
	var endDateErr, startDateErr error
	var pageSize int
	o.Ctx.Input.Bind(&start_date, "from")
	o.Ctx.Input.Bind(&end_date, "to")
	o.Ctx.Input.Bind(&pageSize, "size")

	fmt.Printf("%s %s \n", start_date, end_date)

	if pageSize == 0 {
		pageSize = 30
	}
	if start_date != "" {
		timeString := time.RFC3339
		if end_date == "" {
			endDate = time.Now()
		} else {
			endDate, endDateErr = time.Parse(timeString, end_date)
		}
		startDate, startDateErr = time.Parse(timeString, start_date)

		if startDateErr != nil || endDateErr != nil {

			log.Println(startDateErr, endDateErr)
			return
		}
		logrus.Infof("searching for %d incidents from %s to %s", pageSize, start_date, end_date)
		obs, err = models.GetIncident(startDate, endDate, pageSize)

	} else {
		obs, err = models.GetAllIncidents(pageSize)
	}
	if err == nil {
		var keys []string
		for k := range obs {
			keys = append(keys, k)
		}
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
		var values []models.Incident
		values = make([]models.Incident, 0)
		for _, k := range keys {
			values = append(values, *(obs[k]))
		}
		o.Data["json"] = values
		o.ServeJSON()
	} else {
		o.Data["json"] = fmt.Sprintf("{ 'msg': '%s' }", err.Error())
		o.Ctx.Output.SetStatus(500)
		o.ServeJSON()
		o.StopRun()
	}
}
