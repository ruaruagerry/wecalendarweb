package main

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// 数据
type stAssets struct {
	Tag  string //标识数据唯一
	Data []byte //数据
}

const (
	goData = `package %s

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
)

// 数据
type stAssets struct {
	Tag  string //标识数据唯一
	Data []byte //数据
}

var (
	zipReader = &zip.Reader{}
	zipFile = map[string]*zip.File{}
)

%s
%s`
	goFunc = `
func init() {
	buf := bytes.NewBuffer(data)
	assets := &stAssets{}
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(assets); err != nil {
		msg := fmt.Sprintf("Decode Assets data err:%v", err)
		panic(msg)
	}
	tag := fmt.Sprintf("%x", md5.Sum(assets.Data))
	if tag != assets.Tag {
		msg := fmt.Sprintf("check assets data failed! tag:%s data:%s", tag, assets.Tag)
		panic(msg)
	}

	reader := bytes.NewReader(assets.Data)
	zipReader, err := zip.NewReader(reader, int64(reader.Len()))
	if err != nil {
		msg := fmt.Sprintf("NewReader err:%v", err)
		panic(msg)
	}

	for _, file := range zipReader.File {
		zipFile[file.Name] = file
	}
}

// OpenFile 打开文件
func OpenFile(fileName string) (io.ReadCloser, error) {
	file, has := zipFile[fileName]
	if !has {
		return nil, fmt.Errorf("file %s not find!", fileName)
	}
	return file.Open()
}

// ReadFile 读取文件
func ReadFile(fileName string) ([]byte, error) {
	reader, err := OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// 导出配置
func Export() []byte {
	buf := bytes.NewBuffer(data)
	assets := &stAssets{}
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(assets); err != nil {
		msg := fmt.Sprintf("Decode Assets data err:%v", err)
		panic(msg)
	}
	tag := fmt.Sprintf("%x", md5.Sum(assets.Data))
	if tag != assets.Tag {
		msg := fmt.Sprintf("check assets data failed! tag:%s data:%s", tag, assets.Tag)
		panic(msg)
	}
	return assets.Data
}

// 获得所有文件
func GetFiles() map[string]*zip.File {
	return zipFile
}
`
)

func genData(data []byte, buf *bytes.Buffer) {
	buf.WriteString("\n//配置数据\n")
	buf.WriteString("var data = []byte {\n    ")
	for i, v := range data {
		buf.WriteString(fmt.Sprintf("0x%02x, ", v))
		if (i+1)%20 == 0 {
			buf.WriteString("\n    ")
		}
	}
	buf.WriteString("\n}\n")
}

// ConvertZipToGo 转zip为go文件
func ConvertZipToGo(fileName string, outFileName string, packageName string) {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		msg := fmt.Sprintf("ReadFile %s err:%v", fileName, err)
		panic(msg)
	}
	assets := &stAssets{
		Tag:  fmt.Sprintf("%x", md5.Sum(buf)),
		Data: buf,
	}
	contentBuf := &bytes.Buffer{}
	enc := gob.NewEncoder(contentBuf)
	if err := enc.Encode(assets); err != nil {
		msg := fmt.Sprintf("Encode Assets err:%v", err)
		panic(msg)
	}

	writeBuf := &bytes.Buffer{}
	genData(contentBuf.Bytes(), writeBuf)

	fileBuf := &bytes.Buffer{}
	fileBuf.WriteString(fmt.Sprintf(goData, packageName, writeBuf.String(), goFunc))

	if err := ioutil.WriteFile(outFileName, fileBuf.Bytes(), os.ModePerm); err != nil {
		msg := fmt.Sprintf("Write Faile %s err:%v", outFileName, err)
		panic(msg)
	}
}

var (
	genConfigInputDir = "" //生成配置文件输入的xlsx表目录
	genConfigOutFile  = "" //生成配置文件名
	genPackage        = "" //生成包名
)

func init() {
	flag.StringVar(&genConfigInputDir, "input", "", "输入的xlsx表目录")
	flag.StringVar(&genConfigOutFile, "out", "", "输出配置文件名")
	flag.StringVar(&genPackage, "package", "", "生成包名")
}

func main() {
	flag.Parse()

	log.Println("ConvertZip")
	ConvertZipToGo(genConfigInputDir, genConfigOutFile, genPackage)
	log.Println("ConvertZip Success!")
}
