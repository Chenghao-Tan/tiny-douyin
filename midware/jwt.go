package midware

import (
	"douyin/conf"
	"douyin/repo/redis"
	"douyin/service/type/response"
	"douyin/utility"

	"context"
	"crypto/rand"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 自定义错误类型
var ErrorTokenInvalid = errors.New("token无效")

var signKey = []byte("tiny-douyin")

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

// 生成token
func GenerateToken(userID uint, username string) (token string, err error) {
	expiration := time.Hour * time.Duration(conf.Cfg().System.AutoLogout).Abs()

	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, &customClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Tiny-DouYin",
			Subject:   username,
			Audience:  []string{"Tiny-DouYin"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        randStr(10),
		},
		User_ID:  userID,
		Username: username,
	}).SignedString(signKey)
	if err != nil {
		return "", err
	}

	// 设置为该用户当前唯一有效token
	err = redis.SetJWT(context.TODO(), userID, token, expiration)
	if err != nil {
		return "", err
	}

	return token, nil
}

// 解析/校验token
func ParseToken(tokenString string) (claims *customClaims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &customClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrorTokenInvalid
		}
		return signKey, nil
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

	// 检查token是否已被主动无效化
	if !redis.CheckJWT(context.TODO(), claims.User_ID, tokenString) {
		return nil, ErrorTokenInvalid
	}

	return claims, nil
}

// gin中间件
// jwt鉴权 验证token并提取user_id与username
func MiddlewareAuth(mandatorily bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 尝试从GET中提取token
		tokenString := ctx.Query("token")
		// 若失败则尝试从POST中提取token
		if tokenString == "" {
			tokenString = ctx.PostForm("token")
		}
		// 若无法提取token
		if tokenString == "" {
			if mandatorily {
				utility.Logger().Warnf("MiddlewareAuth warn: 未授权请求")
				ctx.JSON(http.StatusUnauthorized, response.Status{
					Status_Code: -1,
					Status_Msg:  "需要token",
				})
				ctx.Abort()
			} else {
				ctx.Next()
			}
			return
		}

		// 解析/校验token (自动验证有效期等)
		claims, err := ParseToken(tokenString)
		if err != nil {
			if mandatorily {
				if err == ErrorTokenInvalid {
					utility.Logger().Warnf("MiddlewareAuth warn: 未授权请求")
					ctx.JSON(http.StatusUnauthorized, response.Status{
						Status_Code: -1,
						Status_Msg:  "token无效",
					})
				} else {
					utility.Logger().Errorf("MiddlewareAuth err: %v", err)
					ctx.JSON(http.StatusInternalServerError, response.Status{
						Status_Code: -1,
						Status_Msg:  "token解析失败",
					})
				}
				ctx.Abort()
			} else {
				ctx.Next()
			}
			return
		}

		// 提取user_id和username
		ctx.Set("req_id", claims.User_ID)
		ctx.Set("username", claims.Username)

		ctx.Next()
	}
}
