package routers

import (
	"github.com/astaxie/beego"
)

func init() {

	beego.GlobalControllerRouter["bitbucket.org/fseros/beekeeper_api/controllers/ssh:FeedController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/beekeeper_api/controllers/ssh:FeedController"],
		beego.ControllerComments{
			Method: "Post",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			Params: nil})

	beego.GlobalControllerRouter["bitbucket.org/fseros/beekeeper_api/controllers/ssh:FeedController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/beekeeper_api/controllers/ssh:FeedController"],
		beego.ControllerComments{
			Method: "Get",
			Router: `/:objectId`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["bitbucket.org/fseros/beekeeper_api/controllers/ssh:FeedController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/beekeeper_api/controllers/ssh:FeedController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

	beego.GlobalControllerRouter["bitbucket.org/fseros/beekeeper_api/controllers/ssh:FeedController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/beekeeper_api/controllers/ssh:FeedController"],
		beego.ControllerComments{
			Method: "Put",
			Router: `/:objectId`,
			AllowHTTPMethods: []string{"put"},
			Params: nil})

	beego.GlobalControllerRouter["bitbucket.org/fseros/beekeeper_api/controllers/ssh:FeedController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/beekeeper_api/controllers/ssh:FeedController"],
		beego.ControllerComments{
			Method: "Delete",
			Router: `/:objectId`,
			AllowHTTPMethods: []string{"delete"},
			Params: nil})

	beego.GlobalControllerRouter["bitbucket.org/fseros/beekeeper_api/controllers/ssh:IncidentController"] = append(beego.GlobalControllerRouter["bitbucket.org/fseros/beekeeper_api/controllers/ssh:IncidentController"],
		beego.ControllerComments{
			Method: "GetAll",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			Params: nil})

}
