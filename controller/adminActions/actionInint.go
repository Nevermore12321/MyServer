package adminActions

import (
	"MyServer/config"
	midJWT "MyServer/middleware"
	"time"
)

var identityKey = "testKey"

//  创建   gin-jwt 的中间件实例
var JwtAuth *midJWT.GinJWTMiddleware = &midJWT.GinJWTMiddleware{
	Realm:             config.GetStringFromConfig("jwt.realm"),
	SigningAlgorithm:  config.GetStringFromConfig("jwt.signingAlgorithm"),
	Key:               []byte(config.GetStringFromConfig("jwt.secret")),
	Timeout:           time.Hour,
	MaxRefresh:        time.Hour * 2,
	Authenticator:     userLoginAuthenticator,
	Authorizator:      AllUserAuthorizator,
	PayloadFunc:       userPayloadFunc,
	IdentityHandler:   newIdentityHandler,
	IdentityKey:       identityKey,
	TokenLookup:       "header: Authorization",
	TokenHeadName:     "Bearer",
	TimeFunc:          time.Now,
	SendCookie:        false,
	CookieMaxAge:      0,
	SecureCookie:      false,
	CookieHTTPOnly:    true,
	CookieDomain:      "",
	SendAuthorization: true,
	DisableAbort:      true,
	CookieName:        "jwt",
	CookieSameSite:    0,
}

func init() {
	// 初始化  JwtAuth
	initErr := JwtAuth.MiddlewareInit()
	if initErr != nil {
		panic(initErr)
	}

}
