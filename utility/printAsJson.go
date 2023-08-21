package utility

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// PrintAsJson 把嵌套结构体转为json并打印
func PrintAsJson(s any) {
	bs, _ := json.Marshal(s)
	var out bytes.Buffer
	err := json.Indent(&out, bs, "", "\t")
	if err != nil {
		return
	}
	fmt.Printf("\nConfig info: \n%v\n\n", out.String())
}
