package util

import (
	"encoding/json"
	"path/filepath"
	"os"
	"log"
	"io"
	"runtime/debug"
	"reflect"
"strings"
	"strconv"
	"math"
)


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

func Jsonp(obj interface{}, callback string) string {
	b, _ := json.Marshal(obj)
	return callback + "(" + string(b) + ")"
}

func IsZero(val interface{}) bool {
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && IsZero(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		for i := 0; i < v.NumField(); i++ {
			z = z && IsZero(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}

func Version(version string) float64{
	var ret float64
	ret = 0
	versions := strings.Split(version,".")
	vints := make([]int,len(versions))
	vfloats := make([]float64,len(versions))
	for i:=0 ;i<len(versions);i++{
		vints[i],_ = strconv.Atoi(versions[i])
		vfloats[i] = float64(vints[i])
		var f1 float64
		var f2 float64
		f1 = 100
		f2 = float64(2-i)
		ret += vfloats[i]*math.Pow(f1,f2)

	}
	return ret
}