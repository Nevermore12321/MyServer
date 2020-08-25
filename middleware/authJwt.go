package middleware

import (
	"crypto/rsa"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

//  Claims 其实就是 将 payload 中的 哪些字段用于加密
type MapClaims map[string]interface{}

//  用于创建 jwt 中间件的 结构体
type GinJWTMiddleware struct {
	// 给用户显示的范围名称,JWT标识，必填
	Realm string

	// 使用的签名算法：HS256, HS384, HS512, RS256, RS384 or RS512， 可选，默认为HS256
	SigningAlgorithm string

	//  用于签名的服务端密钥。必填
	Key []byte

	//  token 的过期时间，可选，默认 1小时
	Timeout time.Duration

	//  token 更新时间, 在token过期后，MaxRefresh时间内可以更新token，  可选，默认为 0
	MaxRefresh time.Duration

	//  (认证者，也就是用户)回调函数， 在登录接口中使用的验证方法，也就是登录的controller，并返回验证成功后的用户对象。   必填
	Authenticator func(ctx *gin.Context) (interface{}, error)

	//  (授权人，也就是认证者有了token后的校验)回调函数，登录成功后后验证传入的 token 的controller方法，可在此处写权限验证逻辑。  可选
	//  可选，默认是 成功，也就是 true
	Authorizator func(data interface{}, ctx *gin.Context) bool

	// 回调函数， 在 login 期间会被调用，因为 jwt 由三部分组成，分别是 header、payload、signature
	// header 包含使用的签名算法和token类型
	// payload 通常包含一些用户信息，例如用户ID，用户名，用户邮箱等
	// signature 通过 header 和 payload 两部分 再加上 私钥 生成的签名
	// PayloadFunc 可以向 payload 中添加数据，在请求上下文可以通过 ctx.Get("JWT_PAYLOAD")获取，
	// 可选，默认 不会设置其他值
	PayloadFunc func(data interface{}) MapClaims

	//  自定义的 未验证通过的 处理函数
	//  code 是 返回的 status，message 为 错误信息
	Unauthorized func(ctx *gin.Context, code int, message string)

	//  自定义的 登录 login 的 response
	LoginResponse func(ctx *gin.Context, code int, token string, expire time.Time)

	//  自定义的 登出 logout 的 response
	LogoutResponse func(ctx *gin.Context, code int)

	//  自定义的 更新token的 refresh 的 response
	RefreshResponse func(ctx *gin.Context, code int, token string, expire time.Time)

	//  设置 identify key 的 处理函数
	IdentityHandler func(ctx *gin.Context) interface{}

	//  Set the identity key 设置密钥
	IdentityKey string

	//  设置 token 的获取位置，
	//  可选的值：
	//  "header:<name>" 头部的某字段 / "query:<name>" 参数中的某字段 / "cookie:<name>" cookie中的某字段
	//  默认值 为 header 中的 Authorization 可用的值， 即 header:Authorization
	TokenLookup string

	//  设置在 请求的携带token字段中的前缀
	//  默认值 ： Authorization: Bearer
	TokenHeadName string

	//  获取当前时间。可以重写它以使用另一个时间值。这对于测试如果服务器使用的时区与令牌不同。
	TimeFunc func() time.Time

	//  当在 gin jwt 中间 有错误时，可以根据不同的错误，返回不同的错误消息
	HTTPStatusMessageFunc func(e error, ctx *gin.Context) string

	//  用于非对称算法的私钥文件 （https）, 文件地址，读取内容后 保存在 privKey 中
	PrivKeyFile string

	//  用于非对称算法的公钥文件 （https），文件地址，读取内容后 保存在 pubKey 中
	PubKeyFile string

	//  用于对称加密算法的 私钥， 只能 当前 package 使用，小写
	privKey *rsa.PrivateKey

	//  用于对称加密算法的 公钥， 只能 当前 package 使用，小写
	pubKey *rsa.PublicKey

	//  可选，将 token 放入 cookie 中返回
	SendCookie bool

	//  Cookie 的有效时间，可选，默认值 等于 Timeout 的值
	CookieMaxAge time.Duration

	//  Cookie 是否 可以 不通过安全的加密通信 传输，也就是 https
	//  如果为 true ，表示 可以通过 http 传输，  如果为 false ，则表示只能通过 https 传输
	SecureCookie bool

	//  表示 指定 Cookie 是否可通过客户端脚本访问
	//  true， 表示不能通过脚本获取Cookie;  false。 表示 客户端可以修改
	CookieHTTPOnly bool

	//  允许 客户端 更改cookie domain
	CookieDomain string

	//  允许 为每一个 request 的 response 的header中添加 authorization 字段
	SendAuthorization bool

	//  禁止在上下文中 使用 abort()
	DisableAbort bool

	//  修改  Cookie 的名称
	CookieName string

	//  允许使用 http.SameSite 的cookie 参数
	//  Cookie 的SameSite属性用来限制第三方 Cookie，从而减少安全风险。
	//  可以有三个值：
	//  1. Strict ： 完全禁止第三方 Cookie，跨站点时，任何情况下都不会发送 Cookie。
	//  2. Lax ： 规则稍稍放宽， 导航到目标网址的 GET 请求，只包括三种情况：链接，预加载请求，GET 表单。
	//  3，None:  关闭 ，但是需要同事设置 Secure 属性为 true
	CookieSameSite http.SameSite
}

//  所有的 错误 定义
var (
	//  表示 Secret key 必填
	ErrMissingSecretKey = errors.New("secret key is required")

	//  表示 403 forbidden
	ErrForbidden = errors.New("you don't have permission to access the resource")

	//  表示 Authenticator 必填，也就是必须有 登录的 controller
	ErrMissingAuthenticatorFunc = errors.New("authenticator is undefined")

	//  表示 用户在登录验证时， 缺少用户名或者密码
	ErrMissingLoginValues = errors.New("missing Username or Password")

	//  表示 验证失败， 原因是 错误的 用户名和密码
	ErrFailedAuthentication = errors.New("incorrect Username or Password")

	//  表示 用户登录时，用户名 未注册
	ErrUsernameNotFound = errors.New("username is not registered")

	//  表示 创建jwt token 失败
	ErrFailedTokenCreation = errors.New("failed to create JWT token")

	//  表示 token 已过期，不能 更新
	ErrExpiredToken = errors.New("token is expired")

	//  表示 请求的头部中，没有 认证 字段
	ErrEmptyAuthHeader = errors.New("auth header is empty")

	//  表示 token 中 缺少 exp 字段
	ErrMissingExpField = errors.New("missing exp field")

	//  表示 exp字段必须是 float64 格式
	ErrWrongFormatOfExp = errors.New("exp must be float64 format")

	//  认证 的 header 无效， 例如 有错误的 Realm
	ErrInvalidAuthHeader = errors.New("auth header is invalid")

	//  表示 如果 认证token 在 URL 的参数中，该参数不存在
	ErrEmptyQueryToken = errors.New("query token is empty")

	//  表示 如果 认证token 在 Cookie 中，token不存在
	ErrEmptyCookieToken = errors.New("cookie token is empty")

	//  表示 如果 认证 token 中 path中， token不存在
	ErrEmptyParamToken = errors.New("parameter token is empty")

	//  表示 签名算法无效，需要为HS256、HS384、HS512、RS256、RS384或RS512
	ErrInvalidSigningAlgorithm = errors.New("invalid signing algorithm")

	//  表示 无法读取 给定的 (非对称)私钥
	ErrNoPrivKeyFile = errors.New("private key file unreadable")

	//  表示 无法读取 给定的 (非对称)公钥
	ErrNoPubKeyFile = errors.New("public key file unreadable")

	// 表示  给定的 (对称)私钥 无效
	ErrInvalidPrivKey = errors.New("private key invalid")

	// 表示  给定的 (对称)公钥 无效
	ErrInvalidPubKey = errors.New("public key invalid")
)

// IdentityKey 的默认值为 identity
var IdentityKey = "identity"

//  结构体  GinJWTMiddleware  的方法

//  ========================初始化 GinJWTMiddleware 结构体==============================

//  GinJWTMiddleware 结构体 的初始化方法， 也就是  指定 默认值
func (mw *GinJWTMiddleware) MiddlewareInit() error {

	//  设置默认 token 的获取位置 为 header 中的 Authorization
	if mw.TokenLookup == "" {
		mw.TokenLookup = "header:Authorization"
	}

	//  设置默认的 签名算法为 HS256
	if mw.SigningAlgorithm == "" {
		mw.SigningAlgorithm = "HS256"
	}

	//  设置 默认的 token 的过期时间
	if mw.Timeout == 0 {
		mw.Timeout = time.Hour
	}

	//  设置 默认的 获取当前时间 的函数
	if mw.TimeFunc == nil {
		mw.TimeFunc = time.Now
	}

	//  从 TokenHeadName (去掉头尾空格），设置默认为 Bearer
	mw.TokenHeadName = strings.TrimSpace(mw.TokenHeadName)
	if len(mw.TokenHeadName) == 0 {
		mw.TokenHeadName = "Bearer"
	}

	//  设置默认的 Authorizator 方法，也就是 登录成功后 验证token的方法，默认为true
	if mw.Authorizator == nil {
		mw.Authorizator = func(data interface{}, ctx *gin.Context) bool {
			return true
		}
	}

	//  设置 默认的 验证不通过 返回的 方法
	if mw.Unauthorized == nil {
		mw.Unauthorized = func(ctx *gin.Context, code int, message string) {
			ctx.JSON(code, gin.H{
				"status":  code,
				"message": message,
			})
		}
	}

	//  设置默认的 登录 login 方法的response 返回响应
	if mw.LoginResponse == nil {
		mw.LoginResponse = func(ctx *gin.Context, code int, token string, expire time.Time) {
			ctx.JSON(code, gin.H{
				"status": code,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		}
	}

	//  设置默认的 登出 logout 方法 的 响应
	if mw.LogoutResponse == nil {
		mw.LogoutResponse = func(ctx *gin.Context, code int) {
			ctx.JSON(http.StatusOK, gin.H{
				"status": http.StatusOK,
			})
		}
	}

	//  设置默认的 更新token 的 refresh 方法的响应
	if mw.RefreshResponse == nil {
		mw.RefreshResponse = func(ctx *gin.Context, code int, token string, expire time.Time) {
			ctx.JSON(code, gin.H{
				"status": code,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
			})
		}
	}

	//  设置 默认的 IdentityKey ,默认值为 identity
	if mw.IdentityKey == "" {
		mw.IdentityKey = IdentityKey
	}

	//  设置默认的 IdentityKey 处理函数, 也就是 获取 IdentityKey 的方法
	if mw.IdentityHandler == nil {
		mw.IdentityHandler = func(ctx *gin.Context) interface{} {
			//  从 ctx 中提取 jwt payload
			claims := ExtractClaims(ctx)
			//  从 payload 中 提取identityKey
			return claims[mw.IdentityKey]
		}
	}

	//  在 gin jwt 插件中有错误时，设置 默认的 错误处理，直接返回错误
	if mw.HTTPStatusMessageFunc == nil {
		mw.HTTPStatusMessageFunc = func(e error, ctx *gin.Context) string {
			return e.Error()
		}
	}

	//  设置 默认的 Realm 的值，默认为 "gin jwt"
	if mw.Realm == "" {
		mw.Realm = "gin jwt"
	}

	//  设置 默认的 cookie 有效时间, 默认值 等于 Timeout 的值
	if mw.CookieMaxAge == 0 {
		mw.CookieMaxAge = mw.Timeout
	}

	//  设置默认的 Cookie name 的值， 默认为 jwt
	if mw.CookieName == "" {
		mw.CookieName = "jwt"
	}

	//  判断 签名算法是否是 非对称算法
	//  如果是 非对称算法，那么需要解析公钥和私钥
	if mw.usingPublicKeyAlgorithm() {
		return mw.readKeys()
	}

	//  用于签名的 key （Secret) 必填，如果没有，报错
	if mw.Key == nil {
		return ErrMissingSecretKey
	}

	return nil
}

//  判断签名算法 是否是 非对称算法
//  对称加密算法： HS256/HS384/HS512
//  非对称加密算法： RS256/RS384/RS512
//  椭圆曲线数据签名算法： （ES256/ES384/ES512）
func (mw *GinJWTMiddleware) usingPublicKeyAlgorithm() bool {
	switch mw.SigningAlgorithm {
	case "RS256", "RS512", "RS384":
		return true
	}
	return false
}

//  分别将 公钥和私钥文件 中的内容 解析成 密钥类型 存放在 privkey 和 pubkey 中。
func (mw *GinJWTMiddleware) readKeys() error {
	//  读取 私钥 文件 并存放在 privkey 中
	err := mw.readPrivateKey()
	if err != nil {
		return err
	}
	//  读取 公钥 文件 并存放在 pubkey 中
	err = mw.readPublicKey()
	if err != nil {
		return err
	}
	return nil
}

//  将 PrivKeyFile 指定 的 公钥文件中的内容 解析成 rsa.PrivateKey 类型的 公钥 存放在 privkey 中
func (mw *GinJWTMiddleware) readPrivateKey() error {
	//  读取 mw.PrivKeyFile 私钥文件
	privKeyData, err := ioutil.ReadFile(mw.PrivKeyFile)
	if err != nil {
		return ErrNoPrivKeyFile
	}

	//   将 私钥 字符串 解析成 rsa.PrivateKey 类型
	key, errParse := jwt.ParseRSAPrivateKeyFromPEM(privKeyData)
	if errParse != nil {
		return ErrInvalidPrivKey
	}
	mw.privKey = key
	return nil
}

//  将 PubKeyFile 指定 的 公钥文件中的内容 解析成 rsa.PublicKey 类型的 公钥 存放在 pubkey 中
func (mw *GinJWTMiddleware) readPublicKey() error {
	//  读取 mw.PrivKeyFile 私钥文件
	pubKeyData, err := ioutil.ReadFile(mw.PubKeyFile)
	if err != nil {
		return ErrNoPubKeyFile
	}

	//   将 私钥 字符串 解析成 rsa.PublicKey 类型
	key, errParse := jwt.ParseRSAPublicKeyFromPEM(pubKeyData)
	if errParse != nil {
		return ErrInvalidPubKey
	}
	mw.pubKey = key
	return nil
}

//  从 ctx 中获取  jwt 的payload 的信息， 通过 ctx.Get("JWT_PAYLOAD")
func ExtractClaims(ctx *gin.Context) MapClaims {
	//  从 ctx 中提取 jwt payload 如果不存在就创建一个空的
	claims, ok := ctx.Get("JWT_PAYLOAD")
	if !ok {
		return make(MapClaims)
	}
	//  类型 强制转换
	return claims.(MapClaims)
}

//   创建 一个 新的 GinJWTMiddleware 结构体，并且没有赋值的 字段使用默认值
func New(mw *GinJWTMiddleware) (*GinJWTMiddleware, error) {
	if err := mw.MiddlewareInit(); err != nil {
		return nil, err
	}

	return mw, nil
}

// ================================完成 Gin JWT 中间件的功能======================================

//  创建gin 中间件 的函数，返回 gin.HandlerFunc
//  MiddlewareUseFunc 用来处理，携带了token，需要验证的请求
func (mw *GinJWTMiddleware) MiddlewareUseFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		mw.MiddlewareImplement(ctx)
	}
}

//    实现 gin jwt 中间的 具体步骤
//  1. 从 request 中 拿到 token string ，解析成 jwt.Token 类型，然后解析出 jwt.Token.Claims 到 实例
func (mw *GinJWTMiddleware) MiddlewareImplement(ctx *gin.Context) {
	//  步骤1, 如果 出错，则表示 token 不正确，返回 401
	claims, err := mw.GetClaimsFromJWT(ctx)
	if err != nil {
		mw.Unauthorized(ctx, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, ctx))
		return
	}

	//  如果 解析出来的 claims ，判断 每个字段 是否符合要求
	//  判断 claims 中有没有 exp 过期时间
	if claims["exp"] == nil {
		mw.Unauthorized(ctx, http.StatusBadRequest, mw.HTTPStatusMessageFunc(ErrMissingExpField, ctx))
		return
	}

	//  判断 claims 中的 exp 过期时间字段，是否是 float 类型
	if _, ok := claims["exp"].(float64); !ok {
		mw.Unauthorized(ctx, http.StatusBadRequest, mw.HTTPStatusMessageFunc(ErrWrongFormatOfExp, ctx))
		return
	}

	//  判断 claims 中的 exp 过期时间，如果超过时间， 则token无效
	if int64(claims["exp"].(float64)) < mw.TimeFunc().Unix() {
		mw.Unauthorized(ctx, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(ErrExpiredToken, ctx))
		return
	}

	//  将 解析出来 的 claim（也就是 payload） 存放在 gin 的 上下文 Context 中，以便使用，
	ctx.Set("JWT_PAYLOAD", claims)

	//  IdentityHandler  方法是   从 gin 的 Context 的 取出 JWT_PAYLOAD，也就是 claims ，从 claims 中取出 identityKey
	//  取出 identity 存在 ctx 中
	identity := mw.IdentityHandler(ctx)
	if identity != nil {
		ctx.Set(mw.IdentityKey, identity)
	}

	//  身份验证，如果身份验证失败，返回 403 forbidden
	//   用于验证 identityKey
	if !mw.Authorizator(identity, ctx) {
		mw.Unauthorized(ctx, http.StatusForbidden, mw.HTTPStatusMessageFunc(ErrForbidden, ctx))
		return
	}

	ctx.Next()
}

//  从 JWT 中 获取 cliams, 步骤：
//  1. 先从 request 中拿到 token string -->
//  2. 将 token string 解析成 jwt.token 类型， 里面包含 Raw，Method，Header，Claims，Signature，Valid 字段
//  3. 将 解析好的 jwt.Token 中的 claims 字段 解析出来， claims 也就是 payload
func (mw *GinJWTMiddleware) GetClaimsFromJWT(ctx *gin.Context) (MapClaims, error) {
	//  1。 先从 request 中拿到 token string, 并且解析成 jwt.Token 类型
	//  并且这个  token 是通过 keyfunc 中 返回的 密钥 验证过的
	jwtToken, err := mw.parseToken(ctx)
	if err != nil {
		return nil, err
	}

	//  是否为每一个 reqeust 的响应 response 的header 中添加 认证token
	//  如果是，则 在 header 中添加 Authorization 字段
	if mw.SendAuthorization {
		if v, ok := ctx.Get("JWT_TOKEN"); ok {
			ctx.Header("Authorization", mw.TokenHeadName+" "+v.(string))
		}
	}

	//  要解析到的 claims 结构体 实例
	claims := MapClaims{}

	for key, value := range jwtToken.Claims.(jwt.MapClaims) {
		claims[key] = value
	}

	return claims, nil
}

//  从 gin 的上下文 Context 中解析 jwt token
func (mw *GinJWTMiddleware) parseToken(ctx *gin.Context) (*jwt.Token, error) {

	var token string
	var err error

	//   如果 TokenLookup 中定义了多种方式，是以 , 隔开
	//  例如 如果 header 中没有 Authorization 字段，那么就从 cookie中找 token，直到找到为止
	methods := strings.Split(mw.TokenLookup, ",")
	for _, method := range methods {
		//  如果 有 token ，则直接跳出循环
		if len(token) > 0 {
			break
		}

		//  否则，就从每一种 方法中，找 token
		//  将 TokenLookup 格式 以 : 分割， 例如 header: Authorization
		//  前半部分 是 key 也就是 位置， 后半部分 是value 也就是 字段
		parts := strings.Split(strings.TrimSpace(method), ":")
		//  将取出的 key value 前后空格 去掉
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])

		//  根据 token 的位置，也就是前半部分 key 的值，来寻找token
		switch k {
		case "header":
			//  如果 token 存在 header中，根据 header 中的 具体字段，来查找 token
			token, err = mw.jwtFromHeader(ctx, v)
		case "query":
			//  如果 token 存在 URL 的 参数中，根据 参数的 具体名称，来查找 token
			token, err = mw.jwtFromQuery(ctx, v)
		case "param":
			//  如果 token 存在 Path 中，根据 path  中的 具体字段，来查找 token
			token, err = mw.jwtFromParam(ctx, v)
		case "cookie":
			//  如果 token 存在 cookie 中，根据 cookie 名称，来查找 token
			token, err = mw.jwtFromCookie(ctx, v)
		}
	}

	//  如果在 找 token 过程中 有错误， 直接返回
	if err != nil {
		return nil, err
	}

	//  如果没有错误，将 token 解析后返回
	//  jwt.Parse 解析 token ，验证并返回令牌. keyFunc将接收解析的令牌，并应返回需要验证的密钥,
	//  jwt.Token 的 结构体为：
	//    Raw       string                 // 原始的 token 字符串
	//    Method    SigningMethod          // 签名算法
	//    Header    map[string]interface{} // jwt 的 header
	//    Claims    Claims                 // jwt 的 payload
	//    Signature string                 // jwt 的 signature
	//    Valid     bool				   // token 是否有效
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		//  如果 解析后的 token 中的签名算法 与 之前 定义好的 不同 ，则返回错误
		if jwt.GetSigningMethod(mw.SigningAlgorithm) != t.Method {
			return nil, ErrInvalidSigningAlgorithm
		}

		//  如果 使用的 非对称 加密算法，  则返回 public key
		if mw.usingPublicKeyAlgorithm() {
			return mw.pubKey, nil
		}

		//   判断有效后，将 token 存入 gin 的context 上下文，以便后面可以直接拿到
		ctx.Set("JWT_TOKEN", token)

		// 否则，返回 对称加密的 密钥，secret ，也就是 mw.key
		return mw.Key, nil
	})
}

//   从 header 中 取出 jwt token, key 表示 具体存放在什么字段中，例如 header: Authorization, key 应该就是 Authorization
func (mw *GinJWTMiddleware) jwtFromHeader(ctx *gin.Context, fieldName string) (string, error) {
	//  获取到 请求 头部 中 的 具体的 存放 token 的字段
	authHeaderValue := ctx.Request.Header.Get(fieldName)

	//  如果 字段为空，那么就返回 头部没有 token 的错误
	if authHeaderValue == "" {
		return "", ErrEmptyAuthHeader
	}

	//  将 找到的 Header中的 存放 token 的字段 键值对中的值以空格分成两部分，例如 Authorization: Bearer asdfasdf
	//  判断 分开的前半部分，也就是 Bearer 和 定义的 TokenHeadName 是否一致
	parts := strings.SplitN(authHeaderValue, " ", 2)
	if !(len(parts) == 2 && parts[0] == mw.TokenHeadName) {
		return "", ErrInvalidAuthHeader
	}

	//  否则，返回取出的 token
	return parts[1], nil
}

//  从 URL的参数中获取 jwt token
func (mw *GinJWTMiddleware) jwtFromQuery(ctx *gin.Context, fieldName string) (string, error) {
	//   从 URL 的参数 中找到 fieldName 对应的值,也就是 token
	token := ctx.Query(fieldName)

	if token == "" {
		return "", ErrEmptyQueryToken
	}

	return token, nil
}

//  从 path 中获取 jwt token
func (mw *GinJWTMiddleware) jwtFromParam(ctx *gin.Context, fieldName string) (string, error) {
	//  从 path 中找到对应的 fieldName 对应的值，也就是 token
	token := ctx.Param(fieldName)

	if token == "" {
		return "", ErrEmptyParamToken
	}

	return token, nil
}

//  从 cookie 中获取 jwt token
func (mw *GinJWTMiddleware) jwtFromCookie(ctx *gin.Context, fieldName string) (string, error) {
	//  从 Cookie 中找到对应的 fieldName 对应的值，也就是 token
	token, _ := ctx.Cookie(fieldName)

	if token == "" {
		return "", ErrEmptyCookieToken
	}

	return token, nil
}

//  ==========================中间件中使用的 Controller 方法 ==============================

//  LoginHandler  表示 登录请求的处理函数，登录请求流程
//  payload 是 json 格式的 {"username": "xxx", "password": "xxx"}
//  响应结果是  {"token": "xxx"}
func (mw *GinJWTMiddleware) LoginHandler(ctx *gin.Context) {
	//  如果 对用户名密码没有验证，则返回 500 错误
	if mw.Authenticator == nil {
		mw.Unauthorized(ctx, http.StatusInternalServerError, mw.HTTPStatusMessageFunc(ErrMissingAuthenticatorFunc, ctx))
		return
	}

	//  如果 定义 了 对用户名密码的验证操作，则 返回验证后的数据
	//  需要将 验证后的 数据 写入 到token 中
	data, err := mw.Authenticator(ctx)
	if err != nil {
		mw.Unauthorized(ctx, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, ctx))
		return
	}

	//  当 登录请求 收到后，并且用户名密码验证成功后，需要生成一个 jwt token，下次的请求 携带 token
	//  使用 jwt-go 生成 新的 token，并且将 token中的 各个字段赋值，例如 ：
	//  Raw 是原始 token string; Method 是 签名算法; Header 是 jwt header; Claims 是 payload 结构体; Valid 表示是否有效

	//  根据 指定的 签名算法  生成 token
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	//  从 token 结构体中，解析出  payload 结构体 claims
	//  修改 claims 也就是 修改 token.Claims
	claims := token.Claims.(jwt.MapClaims)

	//  判断 指定的 PayloadFunc 中有没有 定义新的 payload， 如果有，需要添加到 claims 中
	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(data) {
			claims[key] = value
		}
	}

	//  添加 过期时间
	expire := mw.TimeFunc().Add(mw.Timeout)
	claims["exp"] = expire.Unix()
	//  添加 当前时间
	claims["orig_iat"] = mw.TimeFunc().Unix()

	//  生成 token 字符串, 如果失败，表示 创建 token string 错误
	tokenStr, toStrErr := mw.SignedToString(token)
	if toStrErr != nil {
		mw.Unauthorized(ctx, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, ctx))
		return
	}

	//  如果 需要将 token 字符串 通过cookie 返回给 浏览器
	if mw.SendCookie {
		//  cookie 的 过期时间
		expireCookie := mw.TimeFunc().Add(mw.CookieMaxAge)
		maxAge := int(expireCookie.Unix() - mw.TimeFunc().Unix())

		//  设置 cookie 的 限制第三方权限
		if mw.CookieSameSite != 0 {
			ctx.SetSameSite(mw.CookieSameSite)
		}

		//  设置 cookie 中保存的具体的内容, 也是 key-value
		ctx.SetCookie(
			mw.CookieName,
			tokenStr,
			maxAge,
			"/",
			mw.CookieDomain,
			mw.SecureCookie,
			mw.CookieHTTPOnly,
		)
	}

	//  如果 不设置 cookie ，则默认 也会 将 token  作为 response data 返回
	mw.LoginResponse(ctx, http.StatusOK, tokenStr, expire)
}

//  将 设置好的 jwt.Token 格式的 结构体 token 转换成 字符串 token
func (mw *GinJWTMiddleware) SignedToString(token *jwt.Token) (string, error) {
	//  字符串 token 变量
	var tokenStr string
	var err error

	//  判断是否是 非对称加密，
	//  如果是 非对称加密算法，需要对 私钥进行 签名算法，生成 公钥 token
	//  如果是 对称加密算法，只需要对 Key 进行签名算法， 生成 token
	if mw.usingPublicKeyAlgorithm() {
		tokenStr, err = token.SignedString(mw.privKey)
	} else {
		tokenStr, err = token.SignedString(mw.Key)
	}

	return tokenStr, err
}

//  LogoutHandler 表示 登出 请求的处理函数，登出请求流程
//  删除 jwt token
func (mw *GinJWTMiddleware) LogoutHandler(ctx *gin.Context) {
	//   删除 cookie 中 jwt token
	if mw.SendCookie {
		if mw.CookieSameSite != 0 {
			ctx.SetSameSite(mw.CookieSameSite)
		}

		//  将 cookie 设置 为 空
		ctx.SetCookie(
			mw.CookieName,
			"",
			-1,
			"/",
			mw.CookieDomain,
			mw.SecureCookie,
			mw.CookieHTTPOnly,
		)
	}

	//  取消 cookie 后，返回 200 ok
	mw.LogoutResponse(ctx, http.StatusOK)
}

//  RefreshHandler  用来 更新 token，注意： 在更新前，token 必须是有效的才可以
//  步骤：
//  1. 检查 token 是否正确，是否过期，需要解析出  token 中的 claims结构体
//  2. 如果过期，创建 新的 token，并且将 这个 claims 结构体放进去，并且重新 添加 过期时间
//  3， 如果没有过期，返回错误
func (mw *GinJWTMiddleware) RefreshTokenHandler(ctx *gin.Context) {
	//  更新 操作，返回 token string， 过期时间，错误
	tokenStr, expire, err := mw.RefreshToken(ctx)
	if err != nil {
		mw.Unauthorized(ctx, http.StatusUnauthorized, mw.HTTPStatusMessageFunc(err, ctx))
		return
	}

	mw.RefreshResponse(ctx, http.StatusOK, tokenStr, expire)
}

//  更新 token 的具体 操作
func (mw *GinJWTMiddleware) RefreshToken(ctx *gin.Context) (string, time.Time, error) {
	//  首先判断 是否可以 更新 token
	claims, err := mw.CheckIfTokenExpire(ctx)
	if err != nil {
		return "", time.Now(), err
	}

	//  如果可以更新，那么创建 新的 token，并且 更新 token 中的 claims
	newToken := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	newClaims := newToken.Claims.(jwt.MapClaims)

	//  将 以前的 claims 传入 新的 claims
	for key := range claims {
		newClaims[key] = claims[key]
	}

	//  设置新的过期时间 和 当前时间
	expire := mw.TimeFunc().Add(mw.Timeout)
	newClaims["exp"] = expire.Unix()
	newClaims["orig_iat"] = mw.TimeFunc().Unix()

	//  将 新的 token 转换成 token string
	tokenStr, strErr := mw.SignedToString(newToken)

	if strErr != nil {
		return "", time.Now(), err
	}

	//  重新设置 cookie
	if mw.SendCookie {
		expireCookie := mw.TimeFunc().Add(mw.CookieMaxAge)
		maxAge := int(expireCookie.Unix() - mw.TimeFunc().Unix())

		if mw.CookieSameSite != 0 {
			ctx.SetSameSite(mw.CookieSameSite)
		}

		ctx.SetCookie(
			mw.CookieName,
			tokenStr,
			maxAge,
			"/",
			mw.CookieDomain,
			mw.SecureCookie,
			mw.CookieHTTPOnly,
		)
	}
	return tokenStr, expire, nil
}

//  检查 token 是否过期，步骤：
// 1. 从 requset 中 解析出 token string，将 token string 转换成 jwt.Token 格式
// 2. 如果出现错误，并且不是 过期错误，直接返回
// 3。 如果出现了过期时间验证错误，或者 如果没有出现错误，那么就从 token 中 解析出 claims 结构体，orig_iat 当前时间，和过期时间， 判断是否 超过了 设置的最大 更新token时间，MaxRefresh
//  注意： MaxRefresh 时间 一般可以设置为 比 过期时间 大
// 4. 如果过期，则不能更新，如果不过期，则可以更新
func (mw *GinJWTMiddleware) CheckIfTokenExpire(ctx *gin.Context) (jwt.MapClaims, error) {
	//  1. 步骤1 ： 从request 中解析出 token string 并转成 jwt.Token
	token, err := mw.parseToken(ctx)

	if err != nil {
		//  这里出错的原因，如果是 token 过期的错误，需要捕捉，更新 token 过期时间
		//  如果 其他错误，直接返回
		//  将错误 转换成 jwt ValidationError
		validateErr, ok := err.(*jwt.ValidationError)
		//  如果不能转换，或者 错误不是 过期错误，那么就直接返回
		if !ok || validateErr.Errors != jwt.ValidationErrorExpired {
			return nil, err
		}
	}

	//  步骤2： 从token 中解析出 claims
	claims := token.Claims.(jwt.MapClaims)

	// 步骤3： 判断是否过期
	origIat := int64(claims["orig_iat"].(float64))
	if origIat < mw.TimeFunc().Add(-mw.MaxRefresh).Unix() {
		return nil, ErrExpiredToken
	}

	return claims, nil
}
