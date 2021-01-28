package define

const (
	/* App */
	KvApp                     = "app/"
	KvAppCommunicationTimeout = "app/communication-timeout"
	KvAppKey                  = "app/key"
	KvAppTokenExpirationTime  = "app/token-expiration-time"
	KvAppUseHttpCode          = "app/use-http-code"
	/* Log */
	KvLog                    = "log/"
	KvLogLevel               = "log/level"
	KvLogLog                 = "log/log/"
	KvLogOutputMethod        = "log/output-method/"
	KvLogOutputMethodTrace   = "log/output-method/trace"
	KvLogOutputMethodDebug   = "log/output-method/debug"
	KvLogOutputMethodInfo    = "log/output-method/info"
	KvLogOutputMethodWarning = "log/output-method/warning"
	KvLogOutputMethodError   = "log/output-method/error"
	KvLogOutputMethodFatal   = "log/output-method/fatal"
	KvLogOutputMethodPanic   = "log/output-method/panic"
	/* I18n */
	KvI18nEnable = "i18n/enable"
	KvI18nLanguage = "i18n/language/"
	/* Validator */
	KvValidatorEnable = "validator/enable"
	/* Circuit Breaker */
	KvCircuitBreakerEnable = "circuit-breaker/enable"
	/* Cache */
	KvCache = "cache/"
	KvCacheEnable = "cache/enable"
	KvCacheDb = "cache/db"
	KvCacheAddress = "cache/address"
	KvCachePassword = "cache/password"
	/* Service discovery */
	KvServiceDiscovery        = "service-discovery/"
	KvServiceDiscoveryEnable  = "service-discovery/enable"
	KvServiceDiscoveryAddress = "service-discovery/address"
	/* Service */
	KvService = "service/"
	KvServiceRpcPort = "/rpc-port"
	KvServiceRpcTag = "/rpc-tag"
	KvServiceHttpId = "/http-id"
	KvServiceHttpName = "/http-name"
	KvServiceHttpPort = "/http-port"
	KvServiceHttpTag = "/http-tag"
	KvServiceAddress = "/address"
	KvServiceCheckInterval = "/check-interval"
	KvServiceCheckUrl = "/check-url"
	// Service Database
	KvServiceDatabase         = "/database/"		// /service/rpcName/database/
	KvServiceDatabaseEnable   = "/database/enable"
	KvServiceDatabaseHost     = "/database/host"
	KvServiceDatabaseUser     = "/database/user"
	KvServiceDatabasePassword = "/database/password"
	KvServiceDatabaseDbName   = "/database/dbname"
	KvServiceDatabasePort     = "/database/port"
	KvServiceDatabaseTimezone = "/database/timezone"
	KvServiceDatabasePrefix   = "/database/prefix"
	/* Captcha */
	KvCaptcha = "captcha/"
	KvCaptchaEnable = "captcha/enable"
	KvCaptchaSecret = "captcha/secret"
	KvCaptchaTimeout = "captcha/timeout"
	KvCaptchaHost = "captcha/host"
)

func MakeServiceConfField(serviceName, fieldName string) string {
	return KvService + serviceName + fieldName
}