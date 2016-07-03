package user

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
	"github.com/martini-contrib/render"
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/payallmoney/videomanager/src/util"
	"fmt"
	"path/filepath"
	"os"
	"github.com/satori/go.uuid"
	"io"
	"log"
	"runtime"
	"strings"
	"time"
)

func clients(r render.Render, db *mgo.Database, params martini.Params, req *http.Request, w http.ResponseWriter,session sessions.Session,) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	result := []bson.M{}
	db.C("video_client").Find(bson.M{"user": session.Get("user_userid")}).Sort("_id").All(&result)
	for _, value := range result {
		list := bson.M{}
		db.C("video_list").Find(bson.M{"_id":value["id"]}).One(&list)
		value["list"] = list["videolist"];
		status := bson.M{}
		log.Println("id===="+value["_id"].(string))
		db.C("client_status_log").Find(bson.M{"client_id":value["_id"].(string)}).Sort("-reportTime").Limit(1).One(&status)
		if(status["reportTime"]!=nil){
			reportTime := status["reportTime"].(time.Time)
			log.Println(reportTime.Format("2006-01-02 15:04:05.000"))
			log.Println(time.Now().Add(-time.Minute*3).Format("2006-01-02 15:04:05.000"))
			if reportTime.Before(time.Now().Add(-time.Minute*3)){
				value["status"] = "离线";
			}else{
				value["status"] = "在线";
			}
		}else{
			log.Println("....nil...")
			value["status"] = "离线";
		}
	}
	r.JSON(200, result)
}



func videolist(r render.Render, db *mgo.Database, params martini.Params, req *http.Request, w http.ResponseWriter,session sessions.Session,) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	result := []bson.M{}
	db.C("video_list").Find(bson.M{"user": session.Get("user_userid")}).All(&result)
	r.JSON(200, result)
}


func videoversion(r render.Render, db *mgo.Database, params martini.Params, req *http.Request, w http.ResponseWriter) {
	result := bson.M{}
	db.C("video_client").Find(bson.M{"_id": params["id"]}).One(&result)
	//ret := map[string]interface{}{"version":1,"files":[]string{"/uploadjs/tasklib.js"}}
	r.JSON(200, result)
}



func uploadpage(r render.Render, w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.HTML(200, "upload", nil)
}

func videoupload(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database,session sessions.Session, ) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	util.CheckErr(err)
	err = req.ParseMultipartForm(100000)
	util.CheckErr(err)
	//get a ref to the parsed multipart form
	m := req.MultipartForm

	//get the *fileheaders
	//	fmt.Println(json.Marshal(m))
	files := m.File["file_data"]

	filenames := []string{}

	for i, _ := range files {
		//for each fileheader, get a handle to the actual file
		file, _ := files[i].Open()
		defer file.Close()
		//create destination file making sure the path is writeable.
		ext := filepath.Ext(files[i].Filename)
		newname := uuid.NewV4().String()
		newfilename := newname + ext
		if _, err := os.Stat(newfilename); err == nil {
			newfilename = uuid.NewV4().String() + ext
		}
		fmt.Println(dir + "/static/uploadvideo/" + newfilename)
		realpath :=dir + "/static/uploadvideo/" + newfilename
		dst, _ := os.Create(realpath)

		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id := bson.NewObjectId()
		obj := bson.M{"_id":id,
			"user": session.Get("user_userid"),
			"name":files[i].Filename[0:len(files[i].Filename)-len(ext)],
			"hash":nil,
			"status":"正在转换",
			"src":"/uploadvideo/" + newfilename}
		db.C("video_list").Insert(obj)
		filenames = append(filenames, util.Js(obj))
		//进行转换
		go convertVideo( id,newfilename,db)
	}
	r.JSON(200, filenames)
}
func convertVideo(id bson.ObjectId,filename string, db *mgo.Database) bson.M{
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	util.CheckErr(err)
	path := dir+"/static/uploadvideo/"
	if runtime.GOOS == "windows" {
		path = strings.Replace(path,"/","\\",-1);
	}
	var strs string

	if runtime.GOOS == "windows" {
		cmd:=fmt.Sprintf("cd /d %s && HandBrakeCLI -i %s -o %s -q 14 -e x264",path,filename,filename+".mp4")
		strs,err  = util.Cmd(cmd)
	}else {
		cmd:=fmt.Sprintf("cd %s && HandBrakeCLI -i %s -o %s -q 14 -e x264",path,filename,filename+".mp4")
		log.Println("=================================================")
		log.Println(cmd)
		log.Println("=================================================")
		strs,err  = util.Cmd(cmd)
	}

	if err !=nil{
		db.C("video_list").Update(bson.M{"_id":id},bson.M{"_id":bson.M{"status":"转换失败","convertinfo":strs}})
		return bson.M{"status":"转换失败","convertinfo":fmt.Sprintf("%v\r\n%v",err,strs)}
	}else{
		hash , _:=util.ComputeMd5(path+filename+".mp4")
		db.C("video_list").Update(bson.M{"_id":id},bson.M{"$set":bson.M{"status":"等待审核","convertinfo":strs,"hash":hash}})

		return bson.M{"status":"转换成功","convertinfo":string(strs)}
	}
}

func convert(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database, ) {
	var result  bson.M
	db.C("video_list").Find(bson.M{"_id": bson.ObjectIdHex(params["id"])}).One(&result);
	if result !=nil  {
		r.JSON(200, bson.M{"status":"转换出错!","msg":"找不到文件!"})
	}else{
		r.JSON(200, convertVideo(bson.ObjectIdHex(params["id"]),result["src"].(string)[14:],db))
	}
}


func del(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database, ) {
	db.C("video_list").Remove(bson.M{"_id": bson.ObjectIdHex(params["id"])});
	r.JSON(200, nil)
}

func changename(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database) {
	db.C("video_list").Update(bson.M{"_id": bson.ObjectIdHex(params["id"])}, bson.M{"$set":bson.M{"name":params["name"]}})
}

func client_add(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database,session sessions.Session,) {
	var result  bson.M
	db.C("video_client").Find(bson.M{"_id": params["id"]}).One(&result);
	if result !=nil  {
		log.Println(result["user"])
		if(result["user"] == nil){
			db.C("video_client").Update(bson.M{"_id": params["id"]}, bson.M{"$set":bson.M{"user":session.Get("user_userid")}})
			r.JSON(200,bson.M{"success":true,"msg":"绑定成功!"})
		}else{
			r.JSON(200,bson.M{"success":false,"msg":"无法绑定:设备已经与其他用户绑定!"})
		}
	}else{
		r.JSON(200,bson.M{"success":false,"msg":"无法绑定:设备未注册!"})
	}
}
func client_del(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database) {
	db.C("video_client").Remove(bson.M{"_id": params["id"]})
}
func client_unbind(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database) {
	db.C("video_client").Update(bson.M{"_id": params["id"]}, bson.M{"$set":bson.M{"user":nil}})
}




func client_video_add(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database) {
	result := bson.M{}
	db.C("video_list").Find(bson.M{"_id": bson.ObjectIdHex(params["videoid"])}).One(&result);
	str, _ := result["src"].(string)
	str = str+".mp4"
	hash, _ := result["hash"].(string)
	db.C("video_client").Update(bson.M{"_id": params["id"]}, bson.M{"$push":bson.M{"videolist":bson.M{"_id":params["videoid"], "src":str,"hash":hash}}})
}
func client_video_del(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database) {
	//	db.C("video_client").Remove(bson.M{"_id": params["id"] ,"videolist":params["idx"]});
	db.C("video_client").Update(bson.M{"_id": params["id"]}, bson.M{"$unset" : bson.M{"videolist." + params["idx"] : 1 }});
	db.C("video_client").Update(bson.M{"_id": params["id"]}, bson.M{"$pull" : bson.M{"videolist" : nil}});
}

func client_video_change(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database) {
	result := bson.M{}
	db.C("video_list").Find(bson.M{"_id": bson.ObjectIdHex( params["videoid"])}).One(&result);
	fmt.Println(params["id"]);
	db.C("video_client").Update(bson.M{"_id": params["id"]}, bson.M{"$set":bson.M{"videolist." + params["idx"]:bson.M{"_id":params["videoid"], "src":result["src"]}}})
}

