package dockerEdit //独自モジュール dockerEditの実装。

import (
	"archive/tar"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

/*Docker イメージをtarファイル化し展開すると，
その中に名前がDocker IDから始まるjsonファイルが存在する。*/
/*このjsonファイルの"history"配列中の"created_by"オブジェクトに
Dockerfileで使用されたコマンドが記載されている。*/
/*そこで，以下の2つの配列を作成し，これらのコマンドを扱う。*/
type Key struct {
	History []Key2 `json:"history"`
}

type Key2 struct {
	Created_by string `json:"created_by"`
}

/*ローカルに存在するDocker イメージのID idを引数に受け取り，*/
/*対象のDocker イメージをtarファイルに変換し*/
/*その内容をio.ReadCloserで渡す関数。*/
func SaveImage(id string) (io.ReadCloser, error) {
	var err error
	var cli *client.Client

	ctx := context.Background()

	cli, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}

	readCloser, err := cli.ImageSave(ctx, []string{id})
	if err != nil {
		return nil, err
	}

	return readCloser, nil
}

/*tarファイル dimageと, dockerファイルのID didを引数に受け取り*/
/*archive/tarとstrings ライブラリを用いて*/
/*目的のjsonファイルの内容をstring形式で渡す関数。*/
func OpenTar(dimage io.ReadCloser, did string) string {
	var dfile []byte
	file := dimage

	defer file.Close()

	// tarの展開
	tarReader := tar.NewReader(file)
	for {
		tarHeader, err := tarReader.Next()
		if err == io.EOF {
			break
		}

		if tarHeader != nil {
			if strings.Contains(tarHeader.Name, did) == true {
				dfile, _ = ioutil.ReadAll(tarReader)
				break
			}
		}
	}
	return string(dfile)
}

/*string形式のjsonファイルの内容 strを受け取り，*/
/*encoding/json ライブラリを用いて Unmarshalで構造化した後に*/
/*appendで"Created_by"属性を持つ内容(つまりDockerfileのコマンド内容)だけ取り出し*/
/*string配列として受け渡す関数。*/
func Dunmarshal(str string) []string {
	b := []byte(str)
	var k Key
	if err := json.Unmarshal(b, &k); err != nil {
		panic(err)
	}

	var hist []string

	for i := range k.History {
		hist = append(hist, k.History[i].Created_by)
	}
	return hist
}

/*string形式のjsonファイルの内容 outputとファイル(*os.File) fpを受け取り，*/
/*dockerfileの"RUN"に相当する内容や不要な内容を整形し*/
/*Writestring でファイルに出力する関数。*/
func Dstring(output []string, fp *os.File) {
	output[0] = "FROM debian:latest as base"
	for i := range output {
		output[i] = strings.Replace(output[i], "/bin/sh -c #(nop) ", "", -1)
		output[i] = strings.Replace(output[i], "/bin/sh -c", "RUN", -1)
		fp.WriteString(output[i])
		fp.WriteString("\n")
	}
}
