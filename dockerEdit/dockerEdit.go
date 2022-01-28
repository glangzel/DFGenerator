package dockerEdit

import (
	"bytes"
	"encoding/json"
	"strings"
	//"io/ioutil"
	//"os"
	//"log"
	//"github.com/tidwall/gjson"
)

func Djson(beforeJ string) string {
	//j, err := ioutil.ReadFile("docker_debian/5852ca.json")
	var buf bytes.Buffer
	err := json.Indent(&buf, []byte(beforeJ), "", "  ")
	if err != nil {
		panic(err)
	}
	indentJ := buf.String()
	return indentJ
}

func Dstring(str string) string {
	str_replaced := strings.Replace(str, "CMD", "Replaced", -1)
	//同様にしてDockerfileと違う部分をどんどん置き換えていく
	return str_replaced
}
