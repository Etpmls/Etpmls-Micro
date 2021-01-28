package em

import (
	"encoding/json"
	em_protobuf "github.com/Etpmls/Etpmls-Micro/protobuf"
	"net/http"
)

type response struct {}

// Return a json format response according to ResponseWriter
// 根据ResponseWriter返回json格式响应
func (this *response) Http_Json(w http.ResponseWriter, httpCode int, code string, status string, message string, data interface{}, err error) {
	dataBytes, _ := json.Marshal(data)
	resBytes, err := json.Marshal(em_protobuf.Response{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    string(dataBytes),
	})
	if err != nil {
		this.Http_InternalError(w)
		return
	}
	w.WriteHeader(httpCode)
	_, _ = w.Write(resBytes)
}

// Return error response according to ResponseWriter
// 根据ResponseWriter返回error
func (this *response) Http_Error(w http.ResponseWriter, httpCode int, code string, message string, err error) {
	this.Http_Json(w, httpCode, code, "error", message, nil, err)
	return
}

// Return internal error response according to ResponseWriter
// 根据ResponseWriter返回内部error
func (this *response) Http_InternalError(w http.ResponseWriter) {
	b := []byte("{\"code\":\"500000\",\"status\":\"error\",\"message\":\"Internal Error\",\"data\":null}")
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write(b)
	return
}

// Return success response according to ResponseWriter
// 根据ResponseWriter返回success
func (this *response) Http_Success(w http.ResponseWriter, httpCode int, code string, message string, data interface{}) {
	this.Http_Json(w, httpCode, code, "success", message, data, nil)
	return
}

