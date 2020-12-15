package em_library

import (
	"encoding/json"
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
func (this *recaptcha) Verify(secret string, response string) bool {
	c := http.Client{
		Timeout: time.Second * Config.Captcha.Timeout,
	}

	resp, err := c.PostForm("https://" + Config.Captcha.Host + "/recaptcha/api/siteverify", url.Values{"secret": []string{secret}, "response": []string{response}})
	if err != nil {
		Log.Error(err.Error())
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var m = make(map[string]interface{})
	err = json.Unmarshal(body, &m)
	if err != nil {
		Log.Error(err.Error())
		return false
	}
	v, ok := m["success"]
	if ok && v == true {
		return true
	}
	return false
}
