// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"sma_easy_helper/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/api/v1",
		beego.NSNamespace("/clusters",
			beego.NSInclude(
				&controllers.ClusterController{},
				),
		),
		beego.NSNamespace("/operations",
			beego.NSInclude(
				&controllers.OperationController{},
				),
		),
		beego.NSNamespace("/files",
			beego.NSInclude(
				&controllers.FileController{},
				),
		),
		beego.NSNamespace("/init",
			beego.NSInclude(
				&controllers.BaseController{},
				),
		),
		//beego.NSNamespace("/logs",
		//	beego.NSInclude(
		//		&controllers.FileController{},
		//	),
		//),
	)
	beego.AddNamespace(ns)
	// Register this router for the SSH connection
	beego.Router("/api/v1/ws", &controllers.SSHWebSocketController{})
}
