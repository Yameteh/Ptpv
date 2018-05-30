package main

import (
	"github.com/funny/link"
	"github.com/funny/link/codec"
	"github.com/astaxie/beego/logs"
	"fmt"
)

type Gateway struct {
	Domain string
	Port   int
	server *link.Server
}

func NewGateway(domain string, port int) *Gateway {
	return &Gateway{domain, port,nil}
}

func (g *Gateway) GetSession(id uint64) *link.Session{
	return g.server.GetSession(id)
}

func (g *Gateway) start() {
	json := codec.Json()
	json.RegisterName("PtpvMessage",PtpvMessage{})
	address := fmt.Sprintf("%s:%d", g.Domain, g.Port)
	var err error
	g.server, err = link.Listen("tcp", address, json, 0, link.HandlerFunc(handleSessionLoop))
	if err != nil {
		logs.Error(err)
		return
	}
	g.server.Serve()
}

func handleSessionLoop(session *link.Session) {
	for {
		req, err := session.Receive()
		if err != nil {
			logs.Error(err)
			return
		}

		if msg, ok := req.(*PtpvMessage); ok {
			logs.Info(msg)
			handleReqMessage(session, msg)
		}

	}
}

func handleReqMessage(session *link.Session, msg *PtpvMessage) {
	logs.Info("handle req message")
	switch msg.Cmd{
	case "REGISTER":
		ret, c := contactRegist(msg.From, msg.Body)
		if ret == CONTACT_STATE_UNKOWN {
			rsp := PtpvMessage{Cmd:"ACK", From:msg.From, To:msg.To, Channel:"", Body:"400:unknow account"}
			session.Send(rsp)
		} else {
			c.SessionId = session.ID()
			c.State = CONTACT_STATE_REACHABLE
			c.UpdateContactDb()
			AddActiveContact(&c)
			rsp := PtpvMessage{Cmd:"ACK", From:msg.From, To:msg.To, Channel:"", Body:"200:ok"}
			session.Send(rsp)
		}
	case "INVITE":
		cf, exist := GetActiveContact(msg.From)
		if exist {
			switch cf.State {
			case CONTACT_STATE_REACHABLE:
				ct, exist := GetActiveContact(msg.To)
				if exist {
					switch ct.State {
					case CONTACT_STATE_REACHABLE:
						st := gateway.GetSession(ct.SessionId)
						st.Send(msg)
					case CONTACT_STATE_CALLING:
						logs.Error("invite contact calling")
						rsp := PtpvMessage{Cmd:"ACK", From:msg.From, To:msg.To, Channel:"", Body:"404:invite contact calling"}
						session.Send(rsp)
					case CONTACT_STATE_UNKOWN:
						logs.Error("invite contact not registed")
						rsp := PtpvMessage{Cmd:"ACK", From:msg.From, To:msg.To, Channel:"", Body:"403:invite contact registed"}
						session.Send(rsp)

					default:
						logs.Error("unkonw contact state")

					}
				}else {
					logs.Error("invite account not registed")
					rsp := PtpvMessage{Cmd:"ACK", From:msg.From, To:msg.To, Channel:"", Body:"400:invite contact not registed"}
					session.Send(rsp)
				}
			case CONTACT_STATE_CALLING:
				logs.Info("invite when calling")
				rsp := PtpvMessage{Cmd:"ACK", From:msg.From, To:msg.To, Channel:"", Body:"401:you state calling"}
				session.Send(rsp)
			case CONTACT_STATE_UNKOWN:
				logs.Info("invite when not registed")
				rsp := PtpvMessage{Cmd:"ACK", From:msg.From, To:msg.To, Channel:"", Body:"402:you not registed"}
				session.Send(rsp)
			default:
				logs.Error("unknow contact state")
			}
		} else {
			logs.Error("you not reigsted")
			rsp := PtpvMessage{Cmd:"ACK", From:msg.From, To:msg.To, Channel:"", Body:"402:you not registed"}
			session.Send(rsp)
		}
	case "ANSWER":
		ct, exist := GetActiveContact(msg.To)
		if exist && ct.State != CONTACT_STATE_UNKOWN {
			cf, exist := GetActiveContact(msg.From)
			if exist && cf.State == CONTACT_STATE_REACHABLE{
				cf.State = CONTACT_STATE_CALLING
				cf.UpdateContactDb()
			}else {
				logs.Error("answer From not exist")
			}
			st := gateway.GetSession(ct.SessionId)
			st.Send(msg)
		}
	case "INDICATE":
		ct, exist := GetActiveContact(msg.To)
		if exist && ct.State != CONTACT_STATE_UNKOWN {
			st := gateway.GetSession(ct.SessionId)
			st.Send(msg)
		}else {
			logs.Info("indicate To not exist")
		}
	case "DELINDICATE":
		ct, exist := GetActiveContact(msg.To)
		if exist && ct.State != CONTACT_STATE_UNKOWN {
			st := gateway.GetSession(ct.SessionId)
			st.Send(msg)
		}else {
			logs.Info("remove indicate to not exist")
		}
	case "BYE":
		cf, exist := GetActiveContact(msg.From)
		if exist && cf.State != CONTACT_STATE_UNKOWN {
			cf.State = CONTACT_STATE_REACHABLE
			cf.UpdateContactDb()
		}else {
			logs.Info("bye From not exist")
		}
	}
}

