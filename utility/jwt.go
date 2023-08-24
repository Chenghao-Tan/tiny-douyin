package utility

import (
	"douyin/service/type/response"

	"crypto/rand"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 自定义错误类型
var ErrorTokenInvalid = errors.New("token无效")

const signKey = "tiny-douyin"
const expiry = 24 // token过期时间(单位为小时)

type customClaims struct {
	jwt.RegisteredClaims
	User_ID  uint   `json:"user_id"`
	Username string `json:"username"`
}

func randStr(length int) (str string) {
	temp := make([]byte, length)
	_, _ = rand.Read(temp)
	return string(temp)
}

func GenerateToken(userID uint, username string) (token string, err error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, &customClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Tiny-DouYin",
			Subject:   username,
			Audience:  []string{"Tiny-DouYin"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * expiry)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        randStr(10),
		},
		User_ID:  userID,
		Username: username,
	}).SignedString([]byte(signKey))
}

func ParseToken(tokenString string) (claims *customClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrorTokenInvalid
		}
		return []byte(signKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrorTokenInvalid
	}

	claims, ok := token.Claims.(*customClaims)
	if !ok {
		return nil, ErrorTokenInvalid
	}

	return claims, nil
}

// gin中间件
// jwt鉴权 验证token并提取user_id与username
func MiddlewareAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 尝试从GET中提取token
		tokenStr := ctx.Query("token")
		// 若失败则尝试从POST中提取token
		if tokenStr == "" {
			tokenStr = ctx.PostForm("token")
		}
		// 若无法提取token
		if tokenStr == "" {
			Logger().Warnf("MiddlewareAuth warn: 未授权请求")
			ctx.JSON(http.StatusUnauthorized, response.Status{Status_Code: -1, Status_Msg: "需要token"})
			ctx.Abort()
			return
		}

		// 解析/校验token (自动验证有效期等)
		claims, err := ParseToken(tokenStr)
		if err != nil {
			if err == ErrorTokenInvalid {
				Logger().Warnf("MiddlewareAuth warn: 未授权请求")
				ctx.JSON(http.StatusUnauthorized, response.Status{
					Status_Code: -1,
					Status_Msg:  "token无效",
				})
				ctx.Abort()
				return
			} else {
				Logger().Errorf("MiddlewareAuth err: %v", err)
				ctx.JSON(http.StatusInternalServerError, response.Status{
					Status_Code: -1,
					Status_Msg:  "token解析失败",
				})
				ctx.Abort()
				return
			}
		}

		// 提取user_id和username
		ctx.Set("req_id", claims.User_ID)
		ctx.Set("username", claims.Username)

		ctx.Next()
	}
}
