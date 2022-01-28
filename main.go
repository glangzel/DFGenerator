package main

import (
	"dockerEdit/dockerEdit"
	"flag"
	"io/ioutil"
)

func main() {
	f := flag.String("flag1", "hoge", "flag 1")
	flag.Parse()
	j, err := ioutil.ReadFile(*f)
	//"/Users/hoge/fuga.json"
	k := string(j)
	dockerfile := dockerEdit.Djson(k)
	dockerfile = dockerEdit.Dstring(dockerfile)
	print(err)
	print(dockerfile)
}
