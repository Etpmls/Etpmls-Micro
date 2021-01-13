package em

const (
	/* App */
	KvApp = "app/"
	KvAppCommunicationTimeout = "app/communication-timeout"
	KvAppKey = "app/key"
	KvAppTokenExpirationTime = "app/token-expiration-time"
	KvAppUseHttpCode = "app/use-http-code"
	/* Log */
	KvLogLevel = "log/level"
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
	KvServiceDiscovery = "service-discovery/"
	KvServiceDiscoveryEnable = "service-discovery/enable"
	KvServiceDiscoveryServerAddress = "service-discovery/server/address"
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
	/* Database */
	KvDatabase = "database/"
	KvDatabaseEnable = "database/enable"
	KvDatabaseHost = "database/host"
	KvDatabaseUser = "database/user"
	KvDatabasePassword = "database/password"
	KvDatabaseDbName = "database/db-name"
	KvDatabasePort = "database/port"
	KvDatabaseTimezone = "database/timezone"
	KvDatabasePrefix = "database/prefix"
)

func MakeServiceConfField(serviceName, fieldName string) string {
	return KvService + serviceName + fieldName
}