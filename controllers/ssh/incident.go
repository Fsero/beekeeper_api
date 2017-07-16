package controllers

import (
	"fmt"
	"sort"

	"bitbucket.org/fseros/beekeeper_api/models/ssh"

	"time"

	"github.com/astaxie/beego"
)

// Operations about object
type IncidentController struct {
	beego.Controller
}

// @Title Create
// @Description create object
// @Param	body		body 	models.Object	true		"The object content"
// @Success 200 {string} models.Object.Id
// @Failure 403 body is empty
// @router / [post]
// func (o *ObjectController) Post() {
// 	var ob models.Object
// 	json.Unmarshal(o.Ctx.Input.RequestBody, &ob)
// 	objectid := models.AddOne(ob)
// 	o.Data["json"] = map[string]string{"ObjectId": objectid}
// 	o.ServeJSON()
// }

// @Title Get
// @Description find object by objectid
// @Param	IncidentId		path 	string	true		"the objectid you want to get"
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router /:from [get]

func (o *IncidentController) Get() {
	timestamp := o.Ctx.Input.Param(":timestamp")
	if timestamp != "" {
		ob, err := models.GetIncident(time.Now(), time.Now())
		if err != nil {
			o.Data["json"] = err.Error()
		} else {
			o.Data["json"] = ob
		}
	}
	o.ServeJSON()
}

// @Title GetAll
// @Description get all objects
// @Success 200 {object} models.Object
// @Failure 403 :objectId is empty
// @router / [get]
func (o *IncidentController) GetAll() {
	obs, err := models.GetAllIncidents()
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
