package main

import (
	"github.com/astaxie/beego/logs"

)


func main()  {
	initLog()
	logs.Info("ptpv server start")


}

var gateway *Gateway

const json11  = "{\"body\":\"1234\",\"channel\":\"\",\"cmd\":\"REGISTER\",\"from\":\"1000\",\"to\":\"\"}";

func initLog()  {
	logs.SetLogger("console")
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)

	//a := PtpvMessage{}
	//a.Cmd = "fadf"
	//a.From = "faef"
	//a.Channel = "afd"
	//a.To = "daff"
	//
	//value,err := json.Marshal(a)
	//logs.Info(string(value))
	//
	//b := TestMessage{}
	//b.A = "afe"
	//b.B = "adfaef"
	//
	//valueb,err := json.Marshal(b)
	//logs.Info(string(valueb))
	//
	//c := PtpvMessage{}
	//err = json.Unmarshal([]byte(json11),&c)
	//if err != nil {
	//	logs.Error(err)
	//}
	//logs.Info(c.From)
	gateway = NewGateway("0.0.0.0",8787)
	gateway.start()
}




