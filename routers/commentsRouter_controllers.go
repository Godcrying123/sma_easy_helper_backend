package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["sma_easy_helper/controllers:ClusterController"] = append(beego.GlobalControllerRouter["sma_easy_helper/controllers:ClusterController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["sma_easy_helper/controllers:ClusterController"] = append(beego.GlobalControllerRouter["sma_easy_helper/controllers:ClusterController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:cname`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["sma_easy_helper/controllers:OperationController"] = append(beego.GlobalControllerRouter["sma_easy_helper/controllers:OperationController"],
        beego.ControllerComments{
            Method: "GetAll",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["sma_easy_helper/controllers:OperationController"] = append(beego.GlobalControllerRouter["sma_easy_helper/controllers:OperationController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/:oname`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
