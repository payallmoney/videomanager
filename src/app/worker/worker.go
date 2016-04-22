package worker

import (
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/payallmoney/videomanager/src/util"
	"net/http"
)

func findUser(params martini.Params, db *mgo.Database, req *http.Request) string {
	userid := params["userid"]
	var result bson.M
	err := db.C("auth_user").Find(bson.M{"_id": userid}).One(&result)
	if err == nil && result != nil {
		result["password"] = nil;
		result["password_rp"] = nil;
		var videos  []bson.M
		db.C("video_client").Find(bson.M{"user":userid}).Sort("_id").All(&videos)
		result["clients"] = videos
		return util.Jsonp(util.JsonRet{true, "", result}, req)
	}else {
		return util.Jsonp(util.JsonRet{false, "用户不存在!", nil}, req)
	}
}

func findClient(params martini.Params, db *mgo.Database, req *http.Request) string {
	id := params["clientid"]
	userid := params["userid"]
	var result bson.M
	err := db.C("video_client").Find(bson.M{"_id": id}).One(&result)
	if err ==nil && result !=nil  {
		if(result["user"] == nil){
			db.C("video_client").Update(bson.M{"_id": params["id"]}, bson.M{"$set":bson.M{"user":userid}})
			return util.Jsonp(util.JsonRet{true, "绑定成功!", result}, req)
		}else{
			return util.Jsonp(util.JsonRet{false, "无法绑定:设备已经与其他用户绑定!", nil}, req)
		}
	}else{
		return util.Jsonp(util.JsonRet{false, "无法绑定:设备未注册!", nil}, req)
	}
}
