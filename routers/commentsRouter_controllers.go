package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["sma_easy_helper/controllers:BaseController"] = append(beego.GlobalControllerRouter["sma_easy_helper/controllers:BaseController"],
        beego.ControllerComments{
            Method: "InfoClusterRead",
            Router: `/cluster/read`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["sma_easy_helper/controllers:BaseController"] = append(beego.GlobalControllerRouter["sma_easy_helper/controllers:BaseController"],
        beego.ControllerComments{
            Method: "InfoClusterWrite",
            Router: `/cluster/write`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["sma_easy_helper/controllers:BaseController"] = append(beego.GlobalControllerRouter["sma_easy_helper/controllers:BaseController"],
        beego.ControllerComments{
            Method: "InfoOperationRead",
            Router: `/operation/read`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["sma_easy_helper/controllers:BaseController"] = append(beego.GlobalControllerRouter["sma_easy_helper/controllers:BaseController"],
        beego.ControllerComments{
            Method: "InfoOperationWrite",
            Router: `/operation/write`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

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

    beego.GlobalControllerRouter["sma_easy_helper/controllers:FileController"] = append(beego.GlobalControllerRouter["sma_easy_helper/controllers:FileController"],
        beego.ControllerComments{
            Method: "List",
            Router: `/list`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["sma_easy_helper/controllers:FileController"] = append(beego.GlobalControllerRouter["sma_easy_helper/controllers:FileController"],
        beego.ControllerComments{
            Method: "Read",
            Router: `/read`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["sma_easy_helper/controllers:FileController"] = append(beego.GlobalControllerRouter["sma_easy_helper/controllers:FileController"],
        beego.ControllerComments{
            Method: "Save",
            Router: `/write`,
            AllowHTTPMethods: []string{"post"},
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
