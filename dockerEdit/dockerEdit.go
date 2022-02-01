package dockerEdit

import (
	"archive/tar"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

type Key struct {
	History []Key2 `json:"history"`
}

type Key2 struct {
	Created_by string `json:"created_by"`
}

func Dunmarshal(str string) []string {
	//fmt.Println(str)
	b := []byte(str)
	var k Key
	if err := json.Unmarshal(b, &k); err != nil {
		panic(err)
	}
	//fmt.Printf("%+v\n", k.History[0].Created_by)

	var hist []string

	for i := range k.History {
		hist = append(hist, k.History[i].Created_by)
	}
	return hist
}

func Dstring(output []string, fp *os.File) {
	for i := range output {
		output[i] = strings.Replace(output[i], "/bin/sh -c #(nop)", "", -1)
		output[i] = strings.Replace(output[i], "/bin/sh -c", "RUN", -1)
		fp.WriteString(output[i])
		fp.WriteString("\n")
	}
}

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

func OpenTar(dimage io.ReadCloser, did string) string {
	//buf := make([]byte, 1024)
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
