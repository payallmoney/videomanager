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
	m.Get("/logout", Logout)

	m.Any("/video/upload", videoupload)
	m.Any("/video/list/:id", clientlist)
	m.Any("/video/list", videolist)
	m.Any("/client/reg/:id", reg)
	m.Any("/client/active/:id", active)
	m.Any("/client/version/:id", videoversion)
	m.Any("/uploadpage", uploadpage)
	//m.Any("/upload", videouploadpage)
	m.Any("/video/del/:id", del)
	m.Any("/video/changename/:id/:name", changename)

	//客户端
	m.Any("/clients", clients)
	m.Any("/client/add/:id", client_add)
	m.Any("/client/videoadd/:id/:videoid", client_video_add)
	m.Any("/client/videochange/:id/:idx/:videoid", client_video_change)
	m.Any("/client/videodel/:id/:idx", client_video_del)
	m.Any("/client/del/:id", client_del)

	m.Any("/program/version", programVersion)
	m.Any("/program/upload", programUpload)
	m.Any("/program/delete/:version", programDel)
	m.Any("/program/list", programList)
	m.Any("/program/reset", reset)
	//worker

	m.Any("/workers", workers)
	m.Any("/worker/add", worker_add)
	m.Any("/worker/del", worker_del)
	m.Any("/worker/udp", worker_udp)
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
				session.Set("admin_userid", values["_id"])
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
	r.Redirect("/admin")
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
