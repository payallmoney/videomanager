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
)

func clients(r render.Render, db *mgo.Database, params martini.Params, req *http.Request, w http.ResponseWriter,session sessions.Session,) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	result := []bson.M{}
	db.C("video_client").Find(bson.M{"user": session.Get("user_userid")}).Sort("_id").All(&result)
	for _, value := range result {
		list := bson.M{}
		db.C("video_list").Find(bson.M{"_id":value["id"]}).One(&list)
		value["list"] = list["videolist"];
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
		dst, _ := os.Create(dir + "/static/uploadvideo/" + newfilename)
		filenames = append(filenames, newfilename)
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		db.C("video_list").Insert(bson.M{"user": session.Get("user_userid"), "name":newfilename, "src":"/uploadvideo/" + newfilename})
	}
	r.JSON(200, filenames)
}

func del(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database, ) {
	db.C("video_list").Remove(bson.M{"_id": bson.ObjectIdHex(params["id"])});
	r.JSON(200, nil)
}

func changename(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database) {
	db.C("video_list").Update(bson.M{"_id": bson.ObjectIdHex(params["id"])}, bson.M{"$set":bson.M{"name":params["name"]}})
}

func client_add(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database,session sessions.Session,) {
	db.C("video_client").Update(bson.M{"_id": params["id"]}, bson.M{"$set":bson.M{"user":session.Get("user_userid")}})
}
func client_del(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database) {
	db.C("video_client").Remove(bson.M{"_id": params["id"]})
}




func client_video_add(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database) {
	result := bson.M{}
	db.C("video_list").Find(bson.M{"_id": bson.ObjectIdHex(params["videoid"])}).One(&result);
	str, _ := result["src"].(string)
	db.C("video_client").Update(bson.M{"_id": params["id"]}, bson.M{"$push":bson.M{"videolist":bson.M{"_id":params["videoid"], "src":str}}})
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

