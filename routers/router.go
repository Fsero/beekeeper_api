// @APIVersion 1.0.0
// @Title beekeper API
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"bitbucket.org/fseros/beekeeper_api/controllers/ssh"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/ssh/incident",
			beego.NSInclude(
				&controllers.IncidentController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
