package upload

import (
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/letter"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"
)

const (
	BufSize = 1024 * 1024 // 缓冲大小
)

// 文件上传
// @params: fileHeader文件元信息 dst目标文件
// @return error
func Upload(fileHeader *multipart.FileHeader, dst string) error {
	buf := make([]byte, BufSize)

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	split := strings.Split(dst, "/")
	dirSlice := split[1 : len(split)-1]
	dir := strings.Join(dirSlice, "/")
	os.MkdirAll(dir, os.ModePerm)

	openFile, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer openFile.Close()

	for {
		read, _ := file.Read(buf)

		if read == 0 {
			break
		}

		_, err = openFile.Write(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

// 分片合并
// @params: dir目标目录
// @return string文件路径 error
func MergeChunk(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	// 创建文件
	name := files[0].Name()
	fileType := strings.Split(name, ".")[1]
	mergeFile, err := os.Create(dir + letter.GenerateCode(20) + "." + fileType)

	defer mergeFile.Close()

	for _, f := range files {
		readbytes, err := ioutil.ReadFile(dir + f.Name())
		if err != nil {
			return "", err
		}

		_, err = mergeFile.Write(readbytes)
		if err != nil {
			return "", err
		}
	}

	// 去除 "./" 相对路径
	dst := strings.Replace(mergeFile.Name(), "./", "", 1)
	return dst, nil
}
