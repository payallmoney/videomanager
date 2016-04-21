package worker

import (
	"github.com/go-martini/martini"
//	"strings"
//	"encoding/json"
	"github.com/martini-contrib/sessions"
	"gopkg.in/mgo.v2"
	"github.com/martini-contrib/render"
	"net/http"
	"github.com/payallmoney/videomanager/src/util"
	"log"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

func Router(m martini.Router) {
	m.Any("/login", Login)
	m.Any("/logout", Logout)
}

func Login(session sessions.Session, db *mgo.Database, r render.Render , req *http.Request, writer http.ResponseWriter) string {
	writer.Header().Set("Content-Type", "text/javascript")
	//params := util.JsonBody(req)
	//params = req.PostForm
	userid := req.FormValue("userid")
	password := req.FormValue("password")
	log.Println("userid", userid)
	log.Println("password", password)
	if userid != "" {
		result := bson.M{}
		err := db.C("auth_worker").Find(bson.M{"_id": userid}).One(&result)
		if err == nil {
			values := result
			if values["password"] == password {
				session.Set("worker_userid", values["_id"])
				session.Set("worker_username", values["name"])
				values["password"] = nil;
				values["password_rp"] = nil;
				fmt.Println("登录成功!")
				//r.JSON(200, util.Jsonp(util.JsonRet{true, "登录成功", values},req))
				return util.Jsonp(util.JsonRet{true, "登录成功", values},req)
			}
		}
	}
	return util.Jsonp(util.JsonRet{false, "登录失败!用户名或密码错误!", nil},req)
}

func Logout(session sessions.Session, r render.Render) {
	session.Delete("worker_userid")
	r.Redirect("/")
}

func UserInfo(session sessions.Session, r render.Render) {
	r.JSON(200, util.JsonRet{true, "", bson.M{"userid":session.Get("user_userid"), "name":session.Get("user_username")}})
}

func Auth(session sessions.Session, c martini.Context, r render.Render, req *http.Request) {
	v := session.Get("worker_userid")
	log.Println(!noAuth(req))
	if v == nil && !noAuth(req) {
		r.JSON(401,util.Jsonp( util.JsonRet{false, "没有权限!", nil},req))
	}else {
		c.Next();
	}
}

func isJson(req *http.Request) bool {
	log.Println(req.Header.Get("accept"));
	accept :=req.Header.Get("accept")
	if accept == "" || len(accept) <16 {
		return false;
	}else{
		return req.Header.Get("accept")[:16] == "application/json"
	}
}
func noAuth(req *http.Request) bool {
	noauth := bson.M{
		"/worker":true,
		"/worker/":true,
		"/worker/login":true,
	}
	url := req.URL.String()
	if noauth[url] == nil {
		flag := false
		for key , val := range(noauth){
			passurl := key+"?"
			log.Println(url[:len(passurl)],passurl)
			if url[:len(passurl)] == passurl && val.(bool){
				flag = true
				break
			}
		}
		if flag {
			return flag
		}else{
			return false
		}
	}else {
		return noauth[url].(bool)
	}
}