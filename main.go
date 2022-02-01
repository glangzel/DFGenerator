package main

import (
	"dockerEdit/dockerEdit"
	"fmt"
	"os"
)

func main() {
	did := os.Args[1]
	a, _ := dockerEdit.SaveImage(did)

	j := dockerEdit.OpenTar(a, did)

	fp, err2 := os.Create("dockerfile")
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	defer fp.Close()
	//println(j)

	var output []string
	output = dockerEdit.Dunmarshal(j)
	dockerEdit.Dstring(output, fp)

}
