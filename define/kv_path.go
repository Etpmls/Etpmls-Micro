package em_define

// Define KV path
// 定义KV路径
const (
	/* App */
	KvApp                     = "app/"
	KvAppCommunicationTimeout = "app/communication-timeout"
	KvAppKey                  = "app/key"
	KvAppTokenExpirationTime  = "app/token-expiration-time"
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
	KvTranslateEnable = "translate/enable"
	KvTranslateLanguage = "translate/language/"
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
	KvService        = "service/"
	KvServicePort    = "/port"
	KvServiceTag     = "/tag"
	KvServiceAddress = "/address"
	/* Captcha */
	KvCaptcha = "captcha/"
	KvCaptchaEnable = "captcha/enable"
	KvCaptchaSecret = "captcha/secret"
	KvCaptchaTimeout = "captcha/timeout"
	KvCaptchaHost = "captcha/host"
)

// Obtain the full path of the KV service according to the service name and field name
// 根据服务名和字段名获取KV服务完整路径
func GetPathByFieldName(serviceName, fieldName string) string {
	return KvService + serviceName + fieldName
}