package define

const (
	/* App */
	KvApp = "app/"
	KvAppCommunicationTimeout = "app/communication-timeout"
	KvAppKey = "app/key"
	KvAppTokenExpirationTime = "app/token-expiration-time"
	KvAppUseHttpCode = "app/use-http-code"
	KvAppMenu = "app/menu"
	KvAppMenuBackup = "app/menu_backup"
	KvAppLog = "app/log/"
	/* Log */
	KvLog = "log/"
	KvLogLevel = "log/level"
	KvLogTrace = "log/trace"
	KvLogDebug = "log/debug"
	KvLogInfo = "log/info"
	KvLogWarning = "log/warning"
	KvLogError = "log/error"
	KvLogFatal = "log/fatal"
	KvLogPanic = "log/panic"
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