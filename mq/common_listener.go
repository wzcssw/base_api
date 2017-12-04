package mq

import (
	"encoding/json"
	"fmt"
)

//CommonListener 公共管道监听者
func CommonListener(msg []byte) {
	var v map[string]interface{}
	if err := json.Unmarshal(msg, &v); err != nil {
		fmt.Printf("parse error is %s", err)
		return
	}
	ucid := v["ucid"].(string)
	fmt.Printf("ucid is %s\n", ucid)
	// 仅允许写方法调用，不允许编写逻辑
	switch ucid {
	case "say":
		go fmt.Printf("msg is %s\n", v["msg"].(string))
	}
}
