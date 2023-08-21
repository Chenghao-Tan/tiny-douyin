package utility

import (
	"douyin/service/type/response"

	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 自定义错误类型
var ErrorTokenInvalid = errors.New("token无效")

const sign_key = "tiny-douyin"

func GenerateToken(user_id uint, username string) (token string, err error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user_id,
		"username": username,
	}).SignedString([]byte(sign_key))
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrorTokenInvalid
		}
		return []byte(sign_key), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, ErrorTokenInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
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
		ctx.Set("user_id", uint(claims["user_id"].(float64))) // token中解析数字默认float64
		ctx.Set("username", claims["username"])

		ctx.Next()
	}
}
