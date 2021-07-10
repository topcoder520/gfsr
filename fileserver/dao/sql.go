package dao

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/topcoder520/gfsr/fileserver/logs"
	"github.com/topcoder520/gfsr/fileserver/model"
	"github.com/topcoder520/gossql"
)

var commonDao *gossql.GoSql

func init() {
	godb, err := sql.Open("sqlite3", "./fileserver.db")
	if err != nil {
		logs.Println(err)
		return
	}
	commonDao = gossql.New(godb)
}

func GetUserByUsername(username string) *model.LoginUser {
	loginUser := &model.LoginUser{}
	commonDao.Query("select * from LoginUser where username=?", username).Unique(loginUser)
	return loginUser
}

func ExistUser(username, password string) bool {
	count := 0
	err := commonDao.Query("select count(1) from LoginUser where username=? and password=?", username, password).Count(&count)
	if err != nil {
		logs.Error(err)
		return false
	}
	return count > 0
}

func AddLoginLog(loginLog *model.LoginLog) {
	commonDao.Insert("insert into LoginLog(LoginUserId,CreateTime) values(?,?)", loginLog.LoginUserId, loginLog.CreateTime)
}
