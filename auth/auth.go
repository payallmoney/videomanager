package auth

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
	"github.com/martini-contrib/render"
	"net/http"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	util "github.com/payallmoney/videomanager/util"
	//	"encoding/json"
)

type JsonRet struct {
	Type   string
	Status int
	Msg    string
	Data   interface{}
}

func Login(session sessions.Session, db *mgo.Database, r render.Render, req *http.Request , writer http.ResponseWriter) string {
	writer.Header().Set("Content-Type", "text/javascript")
	userid := req.FormValue("userid")
	callback := req.FormValue("callback")
	password := req.FormValue("password")
	fmt.Println("userid", userid)
	if userid == "" {
		return util.Jsonp(JsonRet{"login", 401, "请登录", "abc"}, callback)
	} else {
		result := bson.M{}
		err := db.C("auth_user").Find(bson.M{"id": userid}).One(&result)
		fmt.Println(password)
		fmt.Println(result)

		if err == nil {
			values := result
			if values["password"] == password {
				session.Set("userid", values["id"])
				session.Set("username", values["name"])
				values["password"] = nil;
				items := [] interface{}{}
				db.C("good_shares").Find(bson.M{"acc": userid}).All(&items)
				values["items"] = items;
				fmt.Println("登录成功!")
				return util.Jsonp(JsonRet{"login", 200, "登录成功", values}, callback)

			}else {
				return util.Jsonp(JsonRet{"login", 401, "登录失败!密码错误!", nil}, callback)
			}
		}else {
			return util.Jsonp(JsonRet{"login", 401, "登录失败!用户名错误!", nil}, callback)
		}
	}
}

func Refresh(session sessions.Session, db *mgo.Database, r render.Render, req *http.Request , writer http.ResponseWriter) string {
	writer.Header().Set("Content-Type", "text/javascript")
	userid := session.Get("userid")
	callback := req.FormValue("callback")
	if userid == "" {
		return util.Jsonp(JsonRet{"login", 401, "请登录", "abc"}, callback)
	} else {
		result := bson.M{}
		err := db.C("auth_user").Find(bson.M{"id": userid}).One(&result)
		if err == nil {
			values := result
			values["password"] = nil;
			items := [] interface{}{}
			db.C("good_shares").Find(bson.M{"acc": userid}).All(&items)
			values["items"] = items;
			fmt.Println("登录成功!")
			return util.Jsonp(JsonRet{"login", 200, "刷新数据成功", values}, callback)
		}else {
			return util.Jsonp(JsonRet{"login", 401, "刷新数据失败!请稍后再试!", nil}, callback)
		}
	}
}

func Register(session sessions.Session, db *mgo.Database, r render.Render, req *http.Request , writer http.ResponseWriter) string {
	writer.Header().Set("Content-Type", "text/javascript")
	userid := req.FormValue("userid")
	name := req.FormValue("name")
	callback := req.FormValue("callback")
	password := req.FormValue("password")
	fmt.Println("userid", userid)
	if userid == "" {
		return util.Jsonp(JsonRet{"register", 400, "用户名不能为空", nil}, callback)
	} else {
		result := bson.M{}
		err := db.C("auth_user").Find(bson.M{"id": userid}).One(&result)
		if err == nil {
			return util.Jsonp(JsonRet{"login", 400, "注册失败!用户名已经存在!", nil}, callback)
		}else {
			err = db.C("auth_user").Insert(bson.M{"id": userid, "password":password, "name":name, "points":10, "todayadd":0, "creditpoint":0, "creditprecent":100})
			if err == nil {
				return util.Jsonp(JsonRet{"login", 200, "注册成功", nil}, callback)
			}else {
				return util.Jsonp(JsonRet{"login", 400, "注册失败!请稍后再试", nil}, callback)
			}
		}
	}
}

func Logout(session sessions.Session, r render.Render) {
	v := session.Get("userid")
	fmt.Println(v)
	session.Delete("userid")
	r.HTML(200, "login", "登出成功")
}

func Auth(session sessions.Session, c martini.Context, r render.Render) {
	v := session.Get("userid")
	fmt.Println(v)
	if v == nil {
		r.Redirect("/login")
	}else {
		c.Next();
	}
}


func Share(session sessions.Session, db *mgo.Database, r render.Render, req *http.Request , params martini.Params, writer http.ResponseWriter) string {
	writer.Header().Set("Content-Type", "text/javascript")
	userid := session.Get("userid")
	name := req.FormValue("name")
	imgsrc := req.FormValue("imgsrc")
	callback := req.FormValue("callback")
	remark := req.FormValue("remark")
	sharepoint := req.FormValue("sharepoint")
	cat := req.FormValue("cat")


	fmt.Println("name==" + name)

	fmt.Println(userid)
	fmt.Println("imgsrc==" + imgsrc)
	fmt.Println("remark==" + remark)
	fmt.Println("sharepoint==" + sharepoint)
	fmt.Println("cat==" + cat)

	if name == "" {
		return util.Jsonp(JsonRet{"register", 400, "名称不能为空", nil}, callback)
	}else if imgsrc == "" {
		return util.Jsonp(JsonRet{"register", 400, "图片不能为空", nil}, callback)
	}else if remark == "" {
		return util.Jsonp(JsonRet{"register", 400, "详细信息不能为空", nil}, callback)
	}else if cat == "" || len(cat) <= 0 {
		return util.Jsonp(JsonRet{"register", 400, "分类不能为空", nil}, callback)
	}  else {
		err := db.C("good_shares").Insert(
		bson.M{"acc": userid, "name":name, "imgsrc":imgsrc,
		"sharepoint":sharepoint, "remark":remark, "state":"正在出借",
		"cat":cat})
		if err == nil {
			return util.Jsonp(JsonRet{"login", 200, "分享成功", nil}, callback)
		}else {
			return util.Jsonp(JsonRet{"login", 400, "分享失败!请稍后再试", nil}, callback)
		}
	}
}

