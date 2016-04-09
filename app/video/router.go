package video

import (
	"github.com/go-martini/martini"
//	"strings"
//	"encoding/json"
)
func Router(m martini.Router) {
	m.Any("/index", videoindex)
	m.Any("/program/version", programVersion)
	m.Any("/program/upload", programUpload)
	m.Any("/program/delete/:version", programDel)
	m.Any("/program/list", programList)
	m.Any("/program/reset", reset)
	m.Any("/status/:id", client_status)
}


