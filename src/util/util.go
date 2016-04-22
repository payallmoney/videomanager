package util

import (
	"encoding/json"
	"path/filepath"
	"os"
	"log"
	"io"
	"runtime/debug"
	"strings"
	"strconv"
	"math"
	"io/ioutil"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"crypto/md5"
	"fmt"
)

type JsonRet struct {
	Success bool
	Msg    string
	Data   interface{}
}

func CheckErr(err error) {
	var rootpath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		var logfile, logfileerr = os.OpenFile(rootpath + "/client.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
		if logfileerr != nil {
			log.Fatalf("error opening file: %v", logfileerr)
		}
		mWriter := io.MultiWriter(os.Stdout, logfile)
		log.SetOutput(mWriter)

		log.Println(err)
		log.Println(string(debug.Stack()))
		logfile.Close();
		log.SetOutput(os.Stdout)
	}
}

func Jsonp(obj interface{}, req *http.Request) string {
	b, _ := json.Marshal(obj)
	callback := req.FormValue("callback")
	return callback + "(" + string(b) + ")"
}

func IsZero(val interface{}) bool {
	if val == nil {
		return true
	}
	return false
}

func Version(version string) float64 {
	var ret float64
	ret = 0
	versions := strings.Split(version, ".")
	vints := make([]int, len(versions))
	vfloats := make([]float64, len(versions))
	for i := 0; i < len(versions); i++ {
		vints[i], _ = strconv.Atoi(versions[i])
		vfloats[i] = float64(vints[i])
		var f1 float64
		var f2 float64
		f1 = 100
		f2 = float64(2 - i)
		ret += vfloats[i] * math.Pow(f1, f2)
	}
	return ret
}

func JsonBody(req *http.Request) bson.M {
	body, _ := ioutil.ReadAll(req.Body)
	var params bson.M
	json.Unmarshal(body, &params)
	return params
}

func Js(obje interface{}) string {
	bytes ,_ := json.Marshal(obje)
	return string(bytes)
}

func IsJson(req *http.Request) bool {
	return req.Header.Get("accept")[:16] == "application/json"
}

func ComputeMd5(filePath string) (string, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x",hash.Sum(result)), nil
}