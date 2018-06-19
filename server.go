package main

import (
	"github.com/astaxie/beego/logs"

)


func main()  {
	initLog()
	logs.Info("ptpv server start")


}

var gateway *Gateway


func initLog()  {
	logs.SetLogger("console")
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)

	gateway = NewGateway("0.0.0.0",8787)
	gateway.start()
}




