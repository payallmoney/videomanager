package admin

import (
	"net/http"
	"github.com/martini-contrib/render"
	"log"
)

func index(r render.Render, w http.ResponseWriter) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	r.HTML(200, "admin/index", nil)
}

func login(r render.Render, w http.ResponseWriter) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	log.Println("11111111111111111")
	r.HTML(200, "admin/index", nil)
}