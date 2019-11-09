package sugar

import (
	"encoding/json"
	"fmt"
)


func PrettifyStruct(i interface{}) string {
	je, _ := json.MarshalIndent(i, "", "\t")
	return string(je)
}

func PrintStruct(i interface{}) {
	fmt.Println(PrettifyStruct(i))
}