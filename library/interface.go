package em_library

import (
	"context"
	"time"
)

var (
	JwtToken = Interface_Jwt(NewJwtGo(Config.App.Key))
	I18n     = Interface_I18n(&Go_i18n{})
	Cache    = Interface_Cache(NewRedis())
	Log      = Interface_Log(NewLogrus())
	Captcha  = Interface_Captcha(NewRecaptcha())
	Validator  = Interface_Validator(NewValidator())
	ServiceDiscovery = Interface_ServiceDiscovery(NewConsul())
	CircuitBreaker = Interface_CircuitBreaker(NewHystrixGo())
)


// Jwt token interface
// JWT 令牌接口
type Interface_Jwt interface {
	CreateToken(interface{}) (string,  error)
	ParseToken(string) (interface{}, error)
	GetIdByToken(string) (uint, error)						// Get user ID
	GetIssuerByToken(string) (issuer string, err error)		// Get Username
}


// i18n interface
// i18n 接口
type Interface_I18n interface {
	TranslateString (str string, language string) string
	TranslateFromRequest (ctx context.Context, str string) string
}


// Cache interface
// 缓存接口
type Interface_Cache interface {
	GetString (key string) (string, error)
	SetString (key string, value string, time time.Duration)
	DeleteString (list ...string)
	GetHash (key string, field string) (string, error)
	SetHash (key string, value map[string]string)
	DeleteHash (key string, list ...string)
	ClearAllCache()
}


// Instance_Logrus interface
// 日志接口
type Interface_Log interface {
	Panic(args ...interface{})
	Fatal(args ...interface{})
	Error(args ...interface{})
	Warning(args ...interface{})
	Info(args ...interface{})
	Debug(args ...interface{})
	Trace(args ...interface{})
}

// Captcha interface
// 验证码接口
type Interface_Captcha interface {
	Verify(string, string) bool
}

// Validator interface
// 验证器接口
type Interface_Validator interface {
	Validate(request interface{}, my_struct interface{}) error
	ValidateStruct(interface{}) error
}

type Interface_ServiceDiscovery interface {
	GetServiceAddr(service_name string, options map[string]interface{}) (string, error)
	CancelService() error
}

type Interface_CircuitBreaker interface {
	Sync(name string, run func() error, fallBack func(error) error) error
	Async(name string, run func() error, fallBack func(error) error) chan error
}
