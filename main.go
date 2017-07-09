package main

import (
	_ "bitbucket.org/fseros/beekeeper_api/routers"

	"github.com/astaxie/beego"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	//models.GetAllIncidents()
	beego.Run()

}
