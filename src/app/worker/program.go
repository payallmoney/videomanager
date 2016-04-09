package video

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"net/http"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/payallmoney/videomanager/src/util"

	"fmt"
	"path/filepath"
	"os"
//	"strings"
	"io"
//	"encoding/json"
	"log"
	"strconv"
	"encoding/json"
)



func programVersion(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database) {
	var result  bson.M
	//db.collection.find().sort({age:-1}).limit(1)
	db.C("video_program").Find(bson.M{}).Sort("-version").One(&result);

	r.JSON(200, result)
}

func version(db *mgo.Database) int{
	var result  bson.M
	//db.collection.find().sort({age:-1}).limit(1)
	db.C("video_program_version").Find(bson.M{}).Sort("-version").One(&result);
	if result == nil{
		return 0
	}else{
		return result["version"].(int)
	}
}

func programUpload(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database,) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	util.CheckErr(err)
	err = req.ParseMultipartForm(100000)
	util.CheckErr(err)
	//get a ref to the parsed multipart form
	m := req.MultipartForm

	//get the *fileheaders
	//	fmt.Println(json.Marshal(m))
	files := m.File["file_data"]
	filenames := make([]string , len(files))
	log.Println(len(files))

	for i, _ := range files {
		version := version(db)+1
		//for each fileheader, get a handle to the actual file
		file, _ := files[i].Open()
		defer file.Close()
		//create destination file making sure the path is writeable.
		ext := filepath.Ext(files[i].Filename)
		newname := strconv.Itoa(version)
		newfilename := newname + ext

		fmt.Println(dir + "/static/program/" + newfilename)
		dst, _ := os.Create(dir + "/static/program/" + newfilename)
		defer dst.Close()
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Println("version==",version)
		program :=bson.M{"_id":version,"version": version, "src":"/program/" + newfilename}
		str ,_ := json.Marshal(program)
		filenames = append(filenames, string(str))
		//log.Println(string(str))
		//filenames = append(filenames, newfilename)
		db.C("video_program").Insert(program)
		db.C("video_program_version").Insert(bson.M{"_id":"version","version":version})
		db.C("video_program_version").Update(bson.M{"_id":"version"}, bson.M{"$set":bson.M{"version":version}})
	}

	log.Printf("%v\r\n",filenames)

	r.JSON(200, filenames)
}

func reset( db *mgo.Database){
	db.C("video_program_version").Insert(bson.M{"_id":"version","version":0})
	db.C("video_program_version").Update(bson.M{"_id":"version"}, bson.M{"$set":bson.M{"version":0}})
}

func programDel(r render.Render, params martini.Params, req *http.Request, w http.ResponseWriter, db *mgo.Database) {
	version ,_:=strconv.Atoi(params["version"])
	db.C("video_program").Remove(bson.M{"_id":version})
}

func programList(r render.Render, db *mgo.Database, params martini.Params, req *http.Request, w http.ResponseWriter) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	result := []bson.M{}
	db.C("video_program").Find(bson.M{}).Sort("version").All(&result)
	r.JSON(200, result)
}