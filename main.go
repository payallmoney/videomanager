package main

import (
	"github.com/go-martini/martini"
	"flag"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"net/http"
	"github.com/larspensjo/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/payallmoney/videomanager/src/auth"
	"github.com/payallmoney/videomanager/src/util"
	"fmt"
	"github.com/payallmoney/videomanager/src/app/admin"
	"github.com/payallmoney/videomanager/src/app/user"
	"log"
	"reflect"
)



func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	m := getMartini()
	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("goodshare_session", store))

	m.Use(render.Renderer(
		render.Options{
		Directory: "templates",
	}))
	//配置文件
	configFile := flag.String("configfile", "config.ini", "配置文件")
	inicfg, err := config.ReadDefault(*configFile)
	if err != nil {
		panic(err)
	}
	m.Map(inicfg)
	//数据库
	db := util.GetDB(inicfg)
	m.Map(db)
	//缓存
	cache := make(map[string]interface{})
	m.Map(cache)

	m.Any("/client/reg/:id", reg)
	m.Any("/client/active/:id", active)
	m.Any("/client/list/:id", clientlist)
	m.Any("/client/status/:id", client_status)

	//静态内容
	m.Use(martini.Static("static"))


	//需要权限的内容
	m.Group("/admin", admin.Router,admin.Auth)
	m.Group("", user.Router,user.Auth)
	m.Run();
	//m.RunOnAddr(":3333")
}

func index(db *mgo.Database, r render.Render, req *http.Request, inicfg *config.Config) {
	ret := make(map[string]interface{})
	r.HTML(200, "index", ret)
}

func getinit(session sessions.Session, db *mgo.Database, r render.Render, req *http.Request , writer http.ResponseWriter) string {
	writer.Header().Set("Content-Type", "text/javascript")
	ret := bson.M{};
	callback := req.FormValue("callback")
	cats := []interface{}{}
	err := db.C("good_cats").Find(bson.M{}).All(&cats)

	if err == nil {
		ret["cats"] = cats;
		return util.Jsonp(auth.JsonRet{"login", 200, "初始化数据成功", ret}, callback)
	}else {
		return util.Jsonp(auth.JsonRet{"login", 401, "初始化数据成功失败!请稍后再试!", nil}, callback)
	}
}


func getMartini() *martini.ClassicMartini {
	//与martini.Classic相比,去掉了martini.Logger() handler  去掉讨厌的http请求日志
	base := martini.New()
	router := martini.NewRouter()
	base.Use(martini.Recovery())
	base.Use(martini.Static("static"))
	base.MapTo(router, (*martini.Routes)(nil))
	base.Action(router.Handle)
	//	return &martini.ClassicMartini{base, router}
	return martini.Classic()
}


func reg(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db.C("video_client").Insert(bson.M{"_id": params["id"]})
	r.JSON(200, "注册成功")
}

func active(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db.C("video_client").Update(bson.M{"_id": params["id"]}, bson.M{"$set":bson.M{"active":true}})
	db.C("video_client").Insert(bson.M{"_id": params["id"]})
	r.JSON(200, "激活成功")
}

func clientlist(r render.Render, db *mgo.Database, params martini.Params, req *http.Request, w http.ResponseWriter) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")

	result := bson.M{}
	//	db.C("video_client_list").Find(bson.M{"_id": params["id"]}).One(&result)
	//	r.JSON(200, result["videolist"])
	fmt.Println("id==="+params["id"])
	db.C("video_client").Find(bson.M{"_id": params["id"]}).One(&result)
	//	list ,_:= bson.Marshal( result["videolist"])
	list := result["videolist"];
	var ret  []string
	if (list != nil) {
		videolist := reflect.ValueOf(list)
		ret = make([]string,videolist.Len())
		for i := 0; i < videolist.Len(); i ++ {
			row := videolist.Index(i).Elem()
			ret[i] = row.MapIndex(reflect.ValueOf("src")).Elem().String();
		}
	}
	r.JSON(200, ret)
}

func client_status(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database) {
	ret := bson.M{}
	result := bson.M{}
	db.C("video_client").Find(bson.M{"_id": params["id"]}).One(&result)
	//ret := map[string]interface{}{"version":1,"files":[]string{"/uploadjs/tasklib.js"}}
	if(util.IsZero(result)){
		ret["status"]="未注册"
	}else{
		if result["user"] == nil{
			ret["status"]="未绑定"
		}else if result["active"] == nil{
			ret["status"]="未激活"
		}else{
			ret["status"]="已绑定"
		}
	}
	r.JSON(200, ret)
}