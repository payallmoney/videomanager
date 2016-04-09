package auth

import ("database/sql"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessions"
	"github.com/martini-contrib/render"
	"net/http"
	"fmt"
	dbutil "github.com/hcsoft/helpsystem/db"
	erutil "github.com/hcsoft/helpsystem/error"
)

func AdminLogin(session sessions.Session, db *sql.DB, r render.Render, req *http.Request) {
	userid := req.FormValue("userid")
	if userid == "" {
		r.HTML(200, "login", "请登录")
		return
	}
	fmt.Println(userid)
	password := req.FormValue("password")
	rows, err := db.Query("select * from auth_user where userid= ? ", userid)
	erutil.CheckErr(err)
	if rows.Next() {
		values := dbutil.GetOneResult(rows)
		fmt.Println(values["password"]);
		if values["password"] == password {
			session.Set("admin_userid", values["userid"])
			session.Set("admin_username", values["username"])
			r.Redirect("/admin")
			//			r.HTML(200, "admin/index", nil)
		}else {
			r.HTML(200, "login", "用户名或密码错误")
		}
	}else {
		r.HTML(200, "login", "用户名或密码错误")
	}
}

func AdminLogout(session sessions.Session, r render.Render) {
	session.Delete("admin_userid")
	r.HTML(200, "login", "登出成功")
}

func AdminAuth(session sessions.Session, c martini.Context, r render.Render) {
	v := session.Get("admin_userid")
	if v == nil {
		r.Redirect("/login")
	}else {
		c.Next();
	}
}
