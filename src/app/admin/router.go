package admin

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
	"gopkg.in/mgo.v2"
	"github.com/martini-contrib/render"
	"net/http"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"github.com/payallmoney/videomanager/src/util"
)

type JsonRet struct {
	Success bool
	Msg    string
	Data   interface{}
}

func Router(m martini.Router) {
	m.Any("/", index)
	m.Any("", index)
	m.Post("/login", Login)
	m.Any("/video/list", videolist)
	m.Any("/video/clients", clientlist)
	m.Any("/program/version", programVersion)
	m.Any("/program/upload", programUpload)
	m.Any("/program/delete/:version", programDel)
	m.Any("/program/list", programList)
	m.Any("/program/reset", reset)
	m.Any("/status/:id", client_status)
}

func Login(session sessions.Session, db *mgo.Database, r render.Render, req *http.Request, writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "text/javascript")
	params := util.JsonBody(req)
	userid := params["userid"]
	password := params["password"]
	log.Println("userid", userid)
	log.Println("password", password)
	if userid != "" {
		result := bson.M{}
		err := db.C("auth_admin").Find(bson.M{"_id": userid}).One(&result)
		if err == nil {
			values := result
			if values["password"] == password {
				session.Set("admin_userid", values["id"])
				session.Set("admin_username", values["name"])
				values["password"] = nil;
				fmt.Println("登录成功!")
				r.JSON(200, JsonRet{true, "登录成功", values})
				return
			}
		}
	}
	r.JSON(200, JsonRet{false,  "登录失败!用户名或密码错误!", nil})
}

func Logout(session sessions.Session, r render.Render) {
	session.Delete("admin_userid")
	r.HTML(200, "login", "登出成功")
}

func Auth(session sessions.Session, c martini.Context, r render.Render, req *http.Request) {
	v := session.Get("admin_userid")
	if v == nil && !noAuth(req) {
		if isJson(req) {
			r.JSON(401, JsonRet{false, "登录失败!用户名错误!", nil})
		}else {
			r.Redirect("/admin")
		}
	}else {
		c.Next();
	}
}

func isJson(req *http.Request) bool {
	return req.Header.Get("accept")[:16] == "application/json"
}
func noAuth(req *http.Request) bool {
	noauth := bson.M{
		"/admin":true,
		"/admin/login":true    }
	url := req.URL.String()
	if noauth[url] == nil {
		return false
	}else {
		return noauth[url].(bool)
	}
}
