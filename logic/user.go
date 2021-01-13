package logic

/*
 处理用户模块的业务
*/
import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
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

func SignIn(p *models.ParamSignIn) (token string, err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err := mysql.Login(user); err != nil {
		return "", err
	}
	// user是指针在dao层里面已经被改变了 现在可以拿到userID了
	//生成JWT
	return jwt.GenToken(user.UserID, user.Username)
}
