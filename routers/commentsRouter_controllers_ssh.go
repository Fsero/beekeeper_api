package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["bitbucket.org/fseros/beekeeper_api/controllers/ssh:IncidentController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/beekeeper_api/controllers/ssh:IncidentController"],
		beego.ControllerComments{
			Method: "GetIncidents",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
