package main

import (
	"dockerEdit/dockerEdit" //独自モジュール dockerEdit
	"fmt"
	"os"
)

func main() {
	//変数にDocker イメージのIDを指定
	did := os.Args[1]

	//dockerEdit.SaveImage(did) で，
	//ローカルに存在するイメージのID didをもとにイメージをtar形式で取得
	println("[1/3] Please wait for getting Docker Image ...")
	a, err := dockerEdit.SaveImage(did)
	if err != nil {
		println("[ERROR] FAILED to find Docker Image.")
		return
	}

	//dockerEdit.OpenTar(a, did) で，
	//tar形式のDocker イメージ a から イメージのIDから始まるjsonファイルを取得して j に格納
	//(ここにdockerfileで使用されたコマンドなどの情報が格納されている)
	println("[2/3] Please wait for extracting Docker Image ...")
	j := dockerEdit.OpenTar(a, did)

	//空のファイル fp をos.Createで作成しておく。
	//こちらをdockerfile として，取得したコマンドなどを書き込む。
	println("[3/3] Please wait for generating Dockerfile ...")
	fp, err2 := os.Create("dockerfile")
	if err2 != nil {
		fmt.Println("[ERROR] FAILED to make dockerfile.")
		return
	}
	defer fp.Close()

	//dockerEdit.Dunmarshal(j)で，
	//取得した内容からdockerfileで使用されたコマンドだけ選んで取得し，output に格納
	//続いてdockerEdit.Dstring(output, fp)で，
	//outputの内容をdockerfileのコマンドに即した形で整形し，dockerfileに出力して終了。
	var output []string
	output = dockerEdit.Dunmarshal(j)
	dockerEdit.Dstring(output, fp)
	println("DONE!!")
}
