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
"runtime"
"os/exec"
"golang.org/x/text/transform"
"bytes"
"golang.org/x/text/encoding/simplifiedchinese"
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
func Cmd(cmdstr string) (string ,error){
	var commond string
	var param1 string
	if runtime.GOOS == "windows" {
		// "wmic cpu get ProcessorId /format:csv"
		commond = "cmd"
		param1 ="/c"
	}else {
		commond = "bash"
		param1 ="-c"
	}
	cmd := exec.Command(commond, param1 ,cmdstr)
	output, err := cmd.CombinedOutput()

	if runtime.GOOS == "windows" {
		log.Println("999999999999999999999")

		reader := transform.NewReader(bytes.NewReader(output), simplifiedchinese.GBK.NewDecoder())
		outputstr, encodeerr := ioutil.ReadAll(reader)
		CheckErr(encodeerr)
		//reader = transform.NewReader(bytes.NewReader(errs), simplifiedchinese.GBK.NewDecoder())
		//errstr, err := ioutil.ReadAll(reader)
		//CheckErr(err)
		log.Println(string(outputstr))
		//log.Println(string(errstr))

		return string(outputstr), err
	}else{
		log.Println("101010101010")

		log.Println(string(output))
		//log.Println(string(errs))
		return string(output), err
	}
}