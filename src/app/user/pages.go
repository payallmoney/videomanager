package user

import (
	"net/http"
	"github.com/martini-contrib/render"
)

func index(r render.Render, w http.ResponseWriter) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	r.HTML(200, "user/index", nil)
}

