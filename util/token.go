package util

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strings"
)

// Sign 生成签名
func Sign(params map[string]string) (result string) {
	result = ""
	token := params["app_token"]
	delete(params, "app_token")
	keys := make([]string, len(params))
	i := 0
	for k := range params {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	arr := make([]string, len(params))
	for i, k := range keys {
		arr[i] = (k + params[k])
	}
	result = GetMD5Hash(token + strings.Join(arr, "") + token)
	result = strings.ToUpper(result)
	return
}

// SignJSON 生成签名
func SignJSON(params map[string]interface{}) (result string) {
	result = ""
	token := params["app_token"].(string)
	delete(params, "app_token")
	strJSON, err := json.Marshal(params)
	if err != nil {
		return
	}
	result = GetMD5Hash(token + string(strJSON) + token)
	result = strings.ToUpper(result)
	return
}

// GetMD5Hash 获得md5加密
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
