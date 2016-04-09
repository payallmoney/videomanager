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
	"github.com/payallmoney/videomanager/auth"
	util "github.com/payallmoney/videomanager/util"
	"fmt"
	"path/filepath"
	"os"
	"github.com/satori/go.uuid"
	"strings"
	"io"
	"encoding/json"
	"github.com/payallmoney/videomanager/app/video"
	"log"
)



func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	m := getMartini()
	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("goodshare_session", store))
	m.Use(martini.Static("static"))
	m.Use(render.Renderer())
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
	m.Any("/login", auth.Login)
	m.Any("/refresh", auth.Refresh)
	m.Any("/register", auth.Register)
	m.Get("/logout", auth.Logout)

	m.Any("/share", auth.Share,auth.Auth)
	m.Get("/init", getinit)
	m.Get("/items", items)
	m.Any("/img/upload", imgupload)




	m.Post("/img/delete", imgdelete)
	m.Get("/", index)
	//静态内容
	//	m.Use(martini.Static("static"))
	//需要权限的内容
	m.Group("/video", video.Router,auth.AdminAuth)
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

func imgupload(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	util.CheckErr(err)
	err = req.ParseMultipartForm(100000)
	util.CheckErr(err)
	//get a ref to the parsed multipart form
	m := req.MultipartForm

	//get the *fileheaders
	fmt.Println(json.Marshal(m))
	files := m.File["file_data"]

	filenames := []string{}
	for i, _ := range files {
		//for each fileheader, get a handle to the actual file
		file, _ := files[i].Open()
		defer file.Close()
		//create destination file making sure the path is writeable.
		ext := filepath.Ext(files[i].Filename)
		newfilename := uuid.NewV4().String() + ext
		if _, err := os.Stat(newfilename); err == nil {
			newfilename = uuid.NewV4().String()+ext
		}
		fmt.Println(dir + "/static/upload/" + newfilename)
		dst, _ := os.Create(dir + "/static/upload/" + newfilename)
		filenames = append(filenames, "/upload/"+newfilename)
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	ret := map[string]string{"urls":strings.Join(filenames, ",")}
	r.JSON(200, ret)
}



func imgdelete(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	req.ParseForm()
	ret := map[string]string{"msg":"删除成功"}
	r.JSON(200, ret)
}






func items(session sessions.Session, db *mgo.Database, r render.Render, req *http.Request , writer http.ResponseWriter)string{
	writer.Header().Set("Content-Type", "text/javascript")
	callback := req.FormValue("callback")
	items := []interface{}{}
	err := db.C("good_shares").Find(bson.M{}).All(&items)

	if err == nil {
		return util.Jsonp(auth.JsonRet{"login", 200, "初始化数据成功", items}, callback)
	}else {
		return util.Jsonp(auth.JsonRet{"login", 401, "初始化数据成功失败!请稍后再试!", nil}, callback)
	}
}
