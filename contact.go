package main

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/logs"
	"sync"
)

const (
	CONTACT_STATE_UNKOWN = iota
	CONTACT_STATE_REACHABLE
	CONTACT_STATE_CALLING
)

type Contact struct {
	Id        int
	Account   string
	Password  string
	SessionId uint64
	State     int
}

func (c *Contact) UpdateContactDb() {
	o := orm.NewOrm()
	_, err := o.Update(c)
	if err != nil {
		logs.Error("contact db update error")
	}
}

func contactRegist(a string, p string) (int, Contact) {
	o := orm.NewOrm()
	contact := Contact{Account:a, Password:p}
	err := o.Read(&contact)
	if err != nil {
		logs.Error(err)
		return CONTACT_STATE_UNKOWN, contact
	}else {
		return CONTACT_STATE_REACHABLE, contact
	}
}

var activeContacts = struct {
	sync.RWMutex
	cm map[string]*Contact
}{cm:make(map[string]*Contact)}

func AddActiveContact(c *Contact) {
	activeContacts.Lock()
	defer activeContacts.Unlock()
	activeContacts.cm[c.Account] = c
}

func GetActiveContact(acc string) (*Contact, bool) {
	activeContacts.RLock()
	defer activeContacts.RUnlock()
	c, ok := activeContacts.cm[acc]
	return c, ok
}

type Conversation struct {

}

func init() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "yaoguoju:Admin@123@/ptpv?charset=utf8")
	orm.RegisterModel(new(Contact))
}

func insertTestContacts() {
	o := orm.NewOrm()
	o.Using("default")
	contact1 := new(Contact)
	contact1.Account = "1000";
	contact1.Password = "1234";
	logs.Info(o.Insert(contact1))
}







