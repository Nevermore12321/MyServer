package middleware

import (
	"MyServer/app"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dchest/uniuri"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io"
	"strings"
)

//var defaultIgnoreMethods = []string{"GET", "POST", "OPTIONS"}
var defaultIgnoreMethods = []string{"OPTIONS"}
var defaultIgnoreActions = []string{"GetCsrfTokenAction"}

const (
	//  csrfSecret 的值 和 csrfSalt 的值 生成 token 保存在 csrfToken 的键中
	csrfSecret = "csrfSecret" //  存在 session 中的 csrf Secret 字段 key 的名称
	csrfSalt   = "csrfSalt"   //  存在 session 中的 csrf Salt 字段 的 key 的名称
	csrfToken  = "csrfToken"  //  尊在 ctx 中的 生成 token 名称
)

//  对于 csrf middleware 中间件的 选项 结构体
type CsrfOptions struct {
	Secret        string                        // 生成 csrf token 的密钥
	IgnoreMethods []string                      // 对哪些跳过不加密
	IgnoreActions []string                      // 对哪些 controller 跳过
	ErrorFunc     gin.HandlerFunc               // csrf token 验证失败后的 错误处理函数
	TokenGetter   func(ctx *gin.Context) string // 获取 csrf token 的函数
}

//  表示  如果 csrf token 验证不正确，默认的处理函数
func defaultErrorFunc(ctx *gin.Context) {
	panic(errors.New("CSRF token mismatch"))
}

//  默认的 读取 csrf token 的函数， 也可以设置 自己的 TokenGetter 函数，例如下面，我们就从 session 中拿 token
func defaultTokenGetter(ctx *gin.Context) string {
	request := ctx.Request

	//  如果 提交的表单中，有_csrf的 key 值，就是 csrf token
	if token := request.FormValue("_csrf"); len(token) > 0 {
		return token
		//  如果 在 url 中具有 _csrf 的 参数，就是 csrf token
	} else if token := request.URL.Query().Get("_csrf"); len(token) > 0 {
		return token
		//  如果 在 请求的 头部 有携带 名叫 X-CSRF-TOKEN 的 header ，就是 csrf token
	} else if token := request.Header.Get("X-Csrf-Token"); len(token) > 0 {
		return token
		//  如果 在 请求的 头部 有携带 名叫 X-XSRF-TOKEN 的 header ，就是 csrf token
	} else if token := request.Header.Get("X-Xsrf-Token"); len(token) > 0 {
		return token
	}
	//  如果都没有 ，则 没有携带 csrf token
	return ""
}

//  生成 token 函数
func tokenize(secret, salt string) string {
	//  sha1 用于生成 散列值，例如 git 的版本控制，生成遗传字符串
	//  sha1.New() 生成一个 hash 值(各个属性都为空)
	newHash := sha1.New()

	//  将 salt-secret 写入到 生成的 hash 属性中
	hashData := salt + "-" + secret
	_, err := io.WriteString(newHash, hashData)
	if err != nil {
		panic(err)
	}

	//  进行 base64 编码
	hash := base64.URLEncoding.EncodeToString(newHash.Sum(nil))
	return hash
}

//  判断 key 是否在 slice 中
func inArray(arr []string, key string) bool {
	inArr := false

	for _, v := range arr {
		if v == key {
			inArr = true
			break
		}
	}

	return inArr
}

func AddErrorFunc(ctx *gin.Context) {
	ctx.JSON(400, gin.H{
		"detail": "CSRF token mismatch",
	})
	ctx.Abort()
}

func csrfTokenValidate(options CsrfOptions) gin.HandlerFunc {
	ignoreMethods := options.IgnoreMethods
	ignoreAction := options.IgnoreActions
	errorFunc := options.ErrorFunc
	tokenGetter := options.TokenGetter

	//  如果没传  options ，那么就使用默认的 配置
	if ignoreMethods == nil {
		ignoreMethods = defaultIgnoreMethods
	}

	if ignoreAction == nil {
		ignoreAction = defaultIgnoreActions
	}
	if errorFunc == nil {
		errorFunc = defaultErrorFunc
	}
	if tokenGetter == nil {
		tokenGetter = defaultTokenGetter
	}

	//  返回 gin 中间使用的函数类型
	return func(ctx *gin.Context) {
		//  获取 当前 url 处理的 controller name，例如 /v1/login 的 LoginAction，返回的是 MyServer/controller.LoginAction
		fn := ctx.HandlerName()
		//  main.LoginAction 找到 / 后的 所有 controller 名字
		fn = fn[strings.LastIndex(fn, "/"):]

		//  从 ctx 中拿到 session
		session := sessions.Default(ctx)
		//  将 csrfSecret 属性 添加到 ctx 中，以供后面使用
		ctx.Set(csrfSecret, options.Secret)

		//  遍历  ignoreActions 列表中的 controller name ，如果匹配，就直接进行 下一个 中间件，并且controller结束后直接返回。
		for _, action := range ignoreAction {
			if strings.Contains(fn, action) {
				msg := fmt.Sprintf("Controller %s ignore validate CSRF Token", action)
				app.Logger.Info(msg)
				// 进行 下一个 中间件
				ctx.Next()
				//  controller 执行完后， 直接返回
				return
			}
		}

		//  遍历 ignoreMethods 列表，如果 有当前的 请求 method 在其中，则直接进行 下一个 中间件，controller 结束后直接返回
		if inArray(ignoreMethods, ctx.Request.Method) {
			ctx.Next()
			return
		}

		//  从 session 中 获取 csrfSalt 的值，如果不存在，则表示验证失败，直接返回
		salt := session.Get(csrfSalt)
		fmt.Println("come in")
		fmt.Println(salt)
		if salt == nil {
			//  如果 请求第一次来，则拿不到 csrfSalt ，
			errorFunc(ctx)
			return
		}

		//  从 ctx 中 获取 csrfToken 的值，如果不存在，表示验证失败，直接返回
		token := tokenGetter(ctx)

		if tokenize(options.Secret, salt.(string)) != token {
			errorFunc(ctx)
			return
		}

		ctx.Next()

	}
}

//  从 session 中 获取 token
func GetToken(ctx *gin.Context) string {
	//  从 ctx 中获取 session
	session := sessions.Default(ctx)

	//  从 ctx 中 获取 csrfSecret 的值
	secret := ctx.MustGet(csrfSecret).(string)

	//  如果是该 地址 不是请求第一次到，则会将 token 存在 ctx 中，直接获取即可
	token, ok := ctx.Get(csrfToken)
	if ok {
		return token.(string)
	}

	//  如果没找到，那就 生成 新的 salt 并且将 csrfSalt 存入, 只有 salt 存在 session 中
	salt := session.Get(csrfSalt)
	if salt == nil {
		//  生成 salt 并且 存放在 session 中, 生成一个16位随机字符串
		salt = uniuri.New()
		session.Set(csrfSalt, salt)
		err := session.Save()
		if err != nil {
			panic(err)
		}
	}

	//  根据 secret 和 salt 生成   token, 并且 存入 ctx 中
	newToken := tokenize(secret, salt.(string))
	ctx.Set(csrfToken, newToken)

	return newToken
}
