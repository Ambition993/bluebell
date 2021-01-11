package logic

/*
 处理用户模块的业务
*/
import (
	"web_app_base/dao/mysql"
	"web_app_base/models"
	"web_app_base/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	//1判断用户是否已经存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		//数据库查询错误
		return err
	}
	//2生成UID
	userID := snowflake.GenID()
	//构造一个user实例
	u := models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//3 用户密码加密 在dao加密
	//4保存进数据库
	return mysql.InsertUser(&u)
}

func SignIn(p *models.ParamSignIn) (err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	return mysql.Login(user)
}
