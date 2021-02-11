package define

import "time"

const (
	DefaultAppCommunicationTimeout = "60s"
	DefaultLogLevel = "info"
	DefaultAppTokenExpirationTime = time.Hour * 12
	DefaultAppTokenExpirationTimeText = "12h"
	DefaultAppUseHttpCode = "false"
	DefaultCaptchaTimout = time.Second * 10
)