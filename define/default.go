package em_define

import "time"

// Define the default value
// 定义默认值
const (
	DefaultAppCommunicationTimeout = "60s"
	DefaultLogLevel = "info"
	DefaultAppTokenExpirationTime = time.Hour * 12
	DefaultAppTokenExpirationTimeText = "12h"
	DefaultCaptchaTimout = time.Second * 10
)