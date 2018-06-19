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


const (
	CMD_REGISTER = "REGISTER"
	CMD_INVITE   = "INVITE"
	CMD_ANSWER   = "ANSWER"
	CMD_INDICATE = "INDICATE"
	CMD_DELINDICATE = "DELINDICATE"
	CMD_BYE = "BYE"
	CMD_TICK = "TICK"
)

func NewGateway(domain string, port int) *Gateway {
	return &Gateway{domain, port, nil}
}

func (g *Gateway) GetSession(id uint64) *link.Session {
	return g.server.GetSession(id)
}

func (g *Gateway) start() {
	json := codec.Json()
	json.RegisterName("PtpvMessage", PtpvMessage{})
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
			handleReqMessage(session, msg)
		}

	}
}

func handleReqMessage(session *link.Session, msg *PtpvMessage) {
	switch msg.Cmd{
	case CMD_REGISTER:
		logs.Info("REGISTER [%s]",msg.From)
		ret, c := contactRegist(msg.From, msg.Body)
		if ret == CONTACT_STATE_UNKOWN {
			rsp := PtpvMessage{Cmd:"ACK", From:msg.From, To:msg.To, Channel:"", Body:"400:unknow account"}
			session.Send(rsp)
		} else {
			c.Session = session.ID()
			c.State = CONTACT_STATE_REACHABLE
			c.UpdateContactDb()
			AddActiveContact(c)
			rsp := PtpvMessage{Cmd:"ACK", From:msg.From, To:msg.To, Channel:"", Body:"200:ok"}
			session.Send(rsp)
			logs.Info("REGISER success")
		}
	case CMD_INVITE:
		logs.Info("INVITE [%s -> %s]", msg.From, msg.To)
		cf, exist := GetActiveContact(msg.From)
		if exist {
			switch cf.State {
			case CONTACT_STATE_REACHABLE:
				cf.State = CONTACT_STATE_CALLING
				cf.UpdateContactDb()
				ct, exist := GetActiveContact(msg.To)
				if exist {
					switch ct.State {
					case CONTACT_STATE_REACHABLE:
						st := gateway.GetSession(ct.Session)
						if st == nil {
							logs.Info("get st nil")
						} else {
							err := st.Send(*msg)
							if err != nil {
								logs.Error(err)
							}
						}

					case CONTACT_STATE_CALLING:
						logs.Error("INVITE callee calling")
						rsp := PtpvMessage{Cmd:"ACK", From:msg.From, To:msg.To, Channel:"", Body:"404:invite contact calling"}
						session.Send(rsp)
					case CONTACT_STATE_UNKOWN:
						logs.Error("INVITE callee not registed")
						rsp := PtpvMessage{Cmd:"ACK", From:msg.From, To:msg.To, Channel:"", Body:"403:invite contact registed"}
						session.Send(rsp)

					default:
						logs.Error("unkonw contact state")

					}
				} else {
					logs.Error("INVITE caller not registed")
					rsp := PtpvMessage{Cmd:"ACK", From:msg.From, To:msg.To, Channel:"", Body:"400:invite contact not registed"}
					session.Send(rsp)
				}
			case CONTACT_STATE_CALLING:
				logs.Info("INVITE caller calling")
				rsp := PtpvMessage{Cmd:"ACK", From:msg.From, To:msg.To, Channel:"", Body:"401:you state calling"}
				session.Send(rsp)
			case CONTACT_STATE_UNKOWN:
				logs.Info("INVITE caller not registed")
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
	case CMD_ANSWER:
		logs.Info("ANSWER [%s -> %s]", msg.From, msg.To)
		ct, exist := GetActiveContact(msg.To)
		if exist && ct.State != CONTACT_STATE_UNKOWN {
			cf, exist := GetActiveContact(msg.From)
			if exist && cf.State == CONTACT_STATE_REACHABLE {
				cf.State = CONTACT_STATE_CALLING
				cf.UpdateContactDb()
			} else {
				logs.Error("ANSWER caller not exist")
			}
			st := gateway.GetSession(ct.Session)
			err := st.Send(msg)
			if err != nil {
				logs.Error(err)
			}
		}
	case CMD_INDICATE:
		logs.Info("INDICATE [%s -> %s]", msg.From, msg.To)
		ct, exist := GetActiveContact(msg.To)
		if exist && ct.State != CONTACT_STATE_UNKOWN {
			st := gateway.GetSession(ct.Session)
			err := st.Send(msg)
			if err != nil {
				logs.Error(err)
			}
		} else {
			logs.Info("INDICATE callee not exist")
		}
	case CMD_DELINDICATE:
		logs.Info("DELINDICATE [%s -> %s]",msg.From,msg.To)
		ct, exist := GetActiveContact(msg.To)
		if exist && ct.State != CONTACT_STATE_UNKOWN {
			st := gateway.GetSession(ct.Session)
			err := st.Send(msg)
			if err != nil {
				logs.Error(err)
			}

		} else {
			logs.Info("DELINDICATE callee not exist")
		}
	case CMD_BYE:
		logs.Info("BYE [%s -> %s]",msg.From,msg.To)
		cf, exist := GetActiveContact(msg.From)
		if exist && cf.State != CONTACT_STATE_UNKOWN {
			cf.State = CONTACT_STATE_REACHABLE
			cf.UpdateContactDb()
		} else {
			logs.Info("BYE caller not exist")
		}
	case CMD_TICK:
		logs.Info("TICK session %d",session.ID())
	}

}

