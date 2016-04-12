package user

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

	//注册
	m.Any("/register", Register)
	//登录
	m.Post("/login", Login)
	//注销
	m.Get("/logout", Logout)
	//用户信息
	m.Get("/user/userinfo", UserInfo)

	//视频
	m.Any("/video/upload", videoupload)

	m.Any("/video/list", videolist)
	m.Any("/video/del/:id", del)
	m.Any("/video/changename/:id/:name", changename)
	//客户端

	m.Any("/client/version/:id", videoversion)

	m.Any("/clients", clients)
	m.Any("/client/add/:id", client_add)
	m.Any("/client/videoadd/:id/:videoid", client_video_add)
	m.Any("/client/videochange/:id/:idx/:videoid", client_video_change)
	m.Any("/client/videodel/:id/:idx", client_video_del)
	m.Any("/client/del/:id", client_del)

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
		err := db.C("auth_user").Find(bson.M{"_id": userid}).One(&result)
		if err == nil {
			values := result
			if values["password"] == password {
				session.Set("user_userid", values["_id"])
				session.Set("user_username", values["name"])
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
	session.Delete("user_userid")
	r.Redirect("/")
}

func UserInfo(session sessions.Session, r render.Render){
	r.JSON(200, JsonRet{true, "", bson.M{"userid":session.Get("user_userid"),"name":session.Get("user_username")}})
}

func Auth(session sessions.Session, c martini.Context, r render.Render, req *http.Request) {
	v := session.Get("user_userid")
	if v == nil && !noAuth(req) {
		if isJson(req) {
			r.JSON(401, JsonRet{false, "登录失败!用户名错误!", nil})
		}else {
			r.Redirect("/")
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
		"/":true,
		"/register":true,
		"/login":true    }
	url := req.URL.String()
	if noauth[url] == nil {
		return false
	}else {
		return noauth[url].(bool)
	}
}



func Register(session sessions.Session, db *mgo.Database, r render.Render, req *http.Request , writer http.ResponseWriter) {
	params := util.JsonBody(req)
	if params == nil {
		r.JSON(200, JsonRet{false, "用户名不能为空!", nil})
		return
	}
	userid := params["userid"]
	name := params["name"]
	password := params["password"]
	fmt.Println("userid", userid)
	if userid == "" {
		r.JSON(200, JsonRet{false, "用户名不能为空!", nil})
	} else {
		var result  bson.M
		err := db.C("auth_user").Find(bson.M{"_id": userid}).One(&result)
		util.CheckErr(err)
		if result != nil {
			r.JSON(200, JsonRet{false, "注册失败!用户名已经存在!", nil})
		}else {
			err = db.C("auth_user").Insert(bson.M{"_id": userid, "password":password, "name":name})
			if err == nil {
				r.JSON(200, JsonRet{true, "注册成功!", nil})
			}else {
				r.JSON(200, JsonRet{false, "注册失败!注册出错,请与系统管理员联系!", err})
			}
		}
	}
}