package admin

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/payallmoney/videomanager/src/util"
	"log"
)

func users(r render.Render, db *mgo.Database, params martini.Params, req *http.Request, w http.ResponseWriter) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	result := []bson.M{}
	db.C("auth_worker").Find(bson.M{}).Sort("_id").All(&result)
	r.JSON(200, result)
}

func user_del(r render.Render, db *mgo.Database, params martini.Params, req *http.Request, w http.ResponseWriter) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	postdata := util.JsonBody(req)
	db.C("auth_worker").Remove(bson.M{"_id":postdata["account"]})
	r.JSON(200, util.JsonRet{true,"删除成功!",nil})
}


func user_add(r render.Render, db *mgo.Database, params martini.Params, req *http.Request, w http.ResponseWriter) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	postdata := util.JsonBody(req)
	log.Println(postdata);

	var result bson.M
	db.C("auth_worker").Find(bson.M{"_id":postdata["account"]}).One(&result)
	if result != nil{
		db.C("auth_worker").Update(bson.M{"_id":postdata["_id"]},postdata)
	}else{
		postdata["_id"] = postdata["account"]
		db.C("auth_worker").Insert(postdata)
	}
	r.JSON(200, util.JsonRet{true,"保存成功!",postdata})
}

