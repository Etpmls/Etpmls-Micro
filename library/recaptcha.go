package em_library

import (
	"encoding/json"
	"github.com/Etpmls/Etpmls-Micro/v3/define"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type recaptcha struct {

}

func NewRecaptcha() *recaptcha {
	return &recaptcha{}
}

// Verify whether the id and content of the verification code are associated
// 验证验证码的id及内容是否关联
func (this *recaptcha) Verify(response string) (bool, []byte) {
	var tmout time.Duration
	// TODO: 方法需要优化，注释需要优化，不配置timeout本地测试会超时
	pair, _, err := kv.Get(em_define.KvCaptchaTimeout, nil)
	if err != nil || pair == nil {
		Instance_Logrus.Error(err)
		tmout = em_define.DefaultCaptchaTimout
	} else {
		tmout, err = time.ParseDuration(string(pair.Value))
		if err != nil {
			Instance_Logrus.Error(err)
			tmout = em_define.DefaultCaptchaTimout
		}
	}

	c := http.Client{
		Timeout: tmout,
	}

	var host string
	pairHost, _, err := kv.Get(em_define.KvCaptchaHost, nil)
	if err != nil || pairHost == nil || len(pairHost.Value) == 0 {
		Instance_Logrus.Error(err)
		host = "www.google.com"
	} else {
		host = string(pairHost.Value)
	}

	var secret string
	pairSecret, _, err := kv.Get(em_define.KvCaptchaSecret, nil)
	if err != nil || pairSecret == nil || len(pairSecret.Value) == 0 {
		Instance_Logrus.Error(err)
	} else {
		secret = string(pairSecret.Value)
	}

	resp, err := c.PostForm("https://" + host + "/recaptcha/api/siteverify", url.Values{"secret": []string{secret}, "response": []string{response}})
	if err != nil {
		Instance_Logrus.Error(err.Error())
		return false, []byte(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var m = make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		Instance_Logrus.Error(err.Error())
		return false, []byte(err.Error())
	}
	v, ok := m["success"]
	if ok && v == true {
		return true, nil
	}

	Instance_Logrus.Warning("Recaptcha verification failed!")
	return false, body
}