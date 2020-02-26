package handler

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/PPIO/pi-cloud-monitor-backend/common"
	"github.com/PPIO/pi-cloud-monitor-backend/database"
	"github.com/labstack/echo/v4"
)

// 用户成功登录后下发的token的有效期，30天
const userTokenExpireDuration = 30 * 24 * time.Hour

func NewUser(dbUser database.User) User {
	return User{
		dbUser: dbUser,
	}
}

type User struct {
	dbUser database.User
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录
// @Tags 用户
// @Accept  json
// @Produce json
// @Param ReqUserLogin body handler.ReqUserLogin true "用户名和密码"
// @Success 200 {object} RspUserLogin
// @Failure 401 "用户名或密码错误"
// @Failure 500 "内部错误"
// @Router /user/login [post]
func (u *User) Login(c echo.Context) error {
	req := new(ReqUserLogin)
	if err := c.Bind(req); err != nil {
		body, err1 := ioutil.ReadAll(c.Request().Body)
		log.Errorf("can bind login data: %+v, err: %v, err1: %v", string(body), err, err1)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	ctx, cancel := context.WithTimeout(context.Background(), databaseOperationTimeout)
	defer cancel()
	account, err := u.dbUser.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err == database.ErrNoRecord {
			return c.JSON(http.StatusUnauthorized, nil)
		}
		return c.JSON(http.StatusInternalServerError, nil)
	}
	var token *string
	expireAt := time.Now().Add(userTokenExpireDuration).Unix()
	if token, err = common.GetToken(account.ID, expireAt); err != nil {
		log.Errorf("can not get token for uid(%d) err: %v", account.ID, err)
		return c.JSON(http.StatusInternalServerError, nil)
	}
	res := RspUserLogin{
		Token: *token,
	}
	return c.JSON(http.StatusOK, res)
}
