package adminActions

import (
	"MyServer/database/adminDB"
	midJWT "MyServer/middleware"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*

reqeust body 结构体为：
username; password
*/
type UserRequsetBody struct {
	Username string `json: "username"`
	Password string `json: "password"`
}

/*

UserClaims 将用户的某些字段放入到 payload 中的 claims 中
username，userid, Role
*/
type UserClaims struct {
	Username string
	UserID   string
	Role     string
}

/*

将 request body 解析到 结构体实例中，并且验证 request 的 username 和 password 是否正确
*/
func userLoginAuthenticator(ctx *gin.Context) (interface{}, error) {
	//  初始化 UserRequsetBody 实例
	var reqBodyVal UserRequsetBody

	//  将 request body 解析到 该实例中
	if err := ctx.ShouldBind(&reqBodyVal); err != nil {
		return nil, midJWT.ErrMissingLoginValues
	}

	username := reqBodyVal.Username
	password := reqBodyVal.Password
	fmt.Println(reqBodyVal)
	//  在数据库中 ，查询，判断 username 和 password 是否正确
	flag, err := adminDB.CheckAuth(username, password)
	if !flag {
		if err != nil && err.Type == "UsernameNotFound" {
			return nil, midJWT.ErrUsernameNotFound
		} else if err != nil && err.Type == "IncorrectPassword" {
			return nil, midJWT.ErrFailedAuthentication
		}
	}

	//  验证 密码后，返回 需要 加入到 payload 中的 字段
	user, idErr := adminDB.GetUserInfo(username)
	if idErr != nil && idErr.Type == "UsernameNotFound" {
		return nil, midJWT.ErrUsernameNotFound
	}
	return &UserClaims{
		Username: user.Username,
		UserID:   user.UserID,
		Role:     user.Role,
	}, nil
}

/*

claim 设定 为
userid; username;
向 claims 中 添加 新的 字段，也就是
role 表示 角色， value 表示 特定角色有特定的 值

*/
//  重新定义 登录时调用，可将载荷添加到token中, 返回要添加的 map 对象
func userPayloadFunc(data interface{}) midJWT.MapClaims {
	if v, ok := data.(*UserClaims); ok {
		var testStr string
		switch v.Role {
		case "admin":
			testStr = "1234"
		case "user":
			testStr = "5678"
		case "Guest":
			testStr = "9999"
		default:
			testStr = ""
		}

		return midJWT.MapClaims{
			identityKey: testStr,
			"Username":  v.Username,
			"UserID":    v.UserID,
			"Role":      v.Role,
		}

	}
	return midJWT.MapClaims{}
}

//   重新定义 ， 验证登录状态, 也就是验证 identityKey 中保存的内容是否 是 设置的
func newIdentityHandler(ctx *gin.Context) interface{} {
	//  从 ctx 中提取 jwt payload 如果不存在就创建一个空的
	claims := midJWT.ExtractClaims(ctx)

	//  从 payload 中 提取identityKey 对应的值，最终传给 Authorizator 做校验
	return claims[identityKey]
}

//  统一的授权接口，所有用户的 授权结构体都需要 实现 该接口方法
//  data 是 通过  Authenticator 验证完 用户名密码后 返回的 需要 授权的 字段
//  接收的data 是 newIdentityHandler 返回的字段
//  该函数用于 gin-jwt 中间件的传值
func AllUserAuthorizator(data interface{}, ctx *gin.Context) bool {
	//  判断 传进来的 data 是不是 gin-jwt 中的 MapClaims ，也就是 map[string]interface{}
	if v, ok := data.(string); ok {
		userClaims, exist := ctx.Get("JWT_PAYLOAD")
		if !exist {
			return false
		}
		role := userClaims.(midJWT.MapClaims)["Role"]

		switch role {
		case "admin":
			if v == "1234" {
				return true
			}
		case "user":
			if v == "5678" {
				return true
			}
		case "guest":
			if v == "9999" {
				return true
			}
		default:
			if v == "" {
				return true
			}
		}
	}
	return false
}

// @Summary 登录
// @Description 验证用户登录是否成功
// @Tags 校验
// @Security csrf-token
// @param csrf-token header string true "X-CSRF-TOKEN"
// @param body query  UserRequsetBody true "request body"
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"status":200,"data":{},"message":"ok"}"
// @Failure 400 {string} json "{"status":400,"data":{},"message":"bad request"}"
// @Failure 401 {string} json "{"status":401,"data":{},"message":"forbidden"}"
// @Router /admin/login [get]
func LoginAction(ctx *gin.Context) {

	//  如果 对用户名密码没有验证，则返回 500 错误
	if JwtAuth.Authenticator == nil {
		JwtAuth.Unauthorized(ctx, http.StatusInternalServerError, JwtAuth.HTTPStatusMessageFunc(midJWT.ErrMissingAuthenticatorFunc, ctx))
		return
	}

	//  如果 定义 了 对用户名密码的验证操作，则 返回验证后的数据
	//  需要将 验证后的 数据 写入 到token 中
	data, err := JwtAuth.Authenticator(ctx)
	if err != nil {
		fmt.Println(err)
		JwtAuth.Unauthorized(ctx, http.StatusUnauthorized, JwtAuth.HTTPStatusMessageFunc(err, ctx))
		return
	}

	//  当 登录请求 收到后，并且用户名密码验证成功后，需要生成一个 jwt token，下次的请求 携带 token
	//  使用 jwt-go 生成 新的 token，并且将 token中的 各个字段赋值，例如 ：
	//  Raw 是原始 token string; Method 是 签名算法; Header 是 jwt header; Claims 是 payload 结构体; Valid 表示是否有效

	//  根据 指定的 签名算法  生成 token
	token := jwt.New(jwt.GetSigningMethod(JwtAuth.SigningAlgorithm))
	//  从 token 结构体中，解析出  payload 结构体 claims
	//  修改 claims 也就是 修改 token.Claims
	claims := token.Claims.(jwt.MapClaims)

	//  判断 指定的 PayloadFunc 中有没有 定义新的 payload， 如果有，需要添加到 claims 中
	if JwtAuth.PayloadFunc != nil {
		for key, value := range JwtAuth.PayloadFunc(data) {
			claims[key] = value
		}
	}

	//  添加 过期时间
	expire := JwtAuth.TimeFunc().Add(JwtAuth.Timeout)
	claims["exp"] = expire.Unix()
	//  添加 当前时间
	claims["orig_iat"] = JwtAuth.TimeFunc().Unix()

	//  生成 token 字符串, 如果失败，表示 创建 token string 错误
	tokenStr, toStrErr := JwtAuth.SignedToString(token)
	if toStrErr != nil {
		JwtAuth.Unauthorized(ctx, http.StatusUnauthorized, JwtAuth.HTTPStatusMessageFunc(err, ctx))
		return
	}

	//  如果 需要将 token 字符串 通过cookie 返回给 浏览器
	if JwtAuth.SendCookie {
		//  cookie 的 过期时间
		expireCookie := JwtAuth.TimeFunc().Add(JwtAuth.CookieMaxAge)
		maxAge := int(expireCookie.Unix() - JwtAuth.TimeFunc().Unix())

		//  设置 cookie 的 限制第三方权限
		if JwtAuth.CookieSameSite != 0 {
			ctx.SetSameSite(JwtAuth.CookieSameSite)
		}

		//  设置 cookie 中保存的具体的内容, 也是 key-value
		ctx.SetCookie(
			JwtAuth.CookieName,
			tokenStr,
			maxAge,
			"/",
			JwtAuth.CookieDomain,
			JwtAuth.SecureCookie,
			JwtAuth.CookieHTTPOnly,
		)
	}

	//  如果 不设置 cookie ，则默认 也会 将 token  作为 response data 返回
	JwtAuth.LoginResponse(ctx, http.StatusOK, tokenStr, expire)
}
