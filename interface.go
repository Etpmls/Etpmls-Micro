package em

import (
	"context"
	"github.com/Etpmls/Etpmls-Micro/v2/library"
	"time"
)

var (
	/* Must run */
	Kv               = Interface_KV(em_library.NewConsul())
	Log              = Interface_Log(em_library.NewLogrus())
	JwtToken         = Interface_Jwt(em_library.NewJwtGo())
	/* Optional */
	I18n             = Interface_I18n(em_library.NewGoI18n())
	Validator        = Interface_Validator(em_library.NewValidator())
	CircuitBreaker   = Interface_CircuitBreaker(em_library.NewHystrixGo())
	Cache            = Interface_Cache(em_library.NewRedis())
	ServiceDiscovery = Interface_ServiceDiscovery(em_library.NewConsul())
	Captcha          = Interface_Captcha(em_library.NewRecaptcha())
)

// Jwt token interface
// JWT 令牌接口
type Interface_Jwt interface {
	CreateToken(c interface{}, secret string) (string, error)
	ParseToken(tokenString string, secret string) (interface{}, error)
	GetIdByToken(tokenString string, secret string) (int, error)                  // Get user ID
	GetIssuerByToken(tokenString string, secret string) (issuer string, err error) // Get Username
}

// i18n interface
// i18n 接口
type Interface_I18n interface {
	TranslateString(str string, language string) string
	TranslateFromRequest(ctx context.Context, str string) string
}

// Cache interface
// 缓存接口
type Interface_Cache interface {
	GetString(key string) (string, error)
	SetString(key string, value string, time time.Duration)
	DeleteString(list ...string)
	GetHash(key string, field string) (string, error)
	SetHash(key string, value map[string]string)
	DeleteHash(key string, list ...string)
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

type Interface_KV interface {
	ReadKey(key string) (string, error)
	CrateOrUpdateKey(key, value string) error
	DeleteKey(key string) error
	List(prefix string) (map[string]string, error)
}

type Interface_CircuitBreaker interface {
	Sync(name string, run func() error, fallBack func(error) error) error
	Async(name string, run func() error, fallBack func(error) error) chan error
}

