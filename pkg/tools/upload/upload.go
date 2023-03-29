package upload

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/md5"
	"io"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	BufSize      = 1024 * 1024 // 缓冲大小
	fileChunkDir = "./filechunk/"
	mergeFileDir = "./file/"
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
// @params: dir string当前目录 文件的MD5, fileName string文件名
// @return string文件路径 error
func MergeChunk(dir, fileName string) (string, error) {
	files, err := ioutil.ReadDir(fileChunkDir + dir)
	if err != nil {
		return "", err
	}

	// 创建目录
	os.MkdirAll(mergeFileDir, os.ModePerm)

	// 创建文件
	mergeFile, err := os.Create(mergeFileDir + fileName)
	if err != nil {
		return "", err
	}

	defer mergeFile.Close()

	for i := range files {
		content, err := ioutil.ReadFile(fileChunkDir + dir + "/" + fileName + "_" + strconv.Itoa(i))
		if err != nil {
			return "", err
		}

		_, err = mergeFile.Write(content)
		if err != nil {
			_ = os.Remove(mergeFileDir + fileName)
			return mergeFileDir + fileName, err
		}
	}

	// 去除 "./" 相对路径
	dst := strings.Replace(mergeFile.Name(), "./", "", 1)
	return dst, nil
}

// ChunkedUpload 文件分片
// @param: filePath string文件路径， chunkSize int 分片大小 1024*1024=1M
func ChunkedUpload(filePath string, chunkSize int64) {
	file, err := os.Open(filePath) // 打开要上传的大文件
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return
	}

	// 计算分片数量
	fileSize := fileInfo.Size()
	numChunks := int64(math.Ceil(float64(fileSize) / float64(chunkSize)))

	// 创建 HTTP 请求
	url := "http://localhost:8080/upload"
	contentType := "multipart/form-data"
	requestBody := &bytes.Buffer{}
	multipartWriter := multipart.NewWriter(requestBody)
	defer multipartWriter.Close()

	// 循环读取分片并上传
	for i := int64(0); i < numChunks; i++ {
		chunk := make([]byte, chunkSize)
		n, err := file.ReadAt(chunk, i*chunkSize)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading chunk:", err)
			return
		}

		// 添加分片到 HTTP 请求
		chunkBuffer := bytes.NewBuffer(chunk[:n])
		// 请求参数
		// multipartWriter.WriteField("key","value")
		chunkPart, err := multipartWriter.CreateFormFile("file", filepath.Base(file.Name())+"_"+strconv.Itoa(int(i)))
		if err != nil {
			fmt.Println("Error creating form file:", err)
			return
		}
		io.Copy(chunkPart, chunkBuffer)

		// 发送 HTTP 请求
		request, err := http.NewRequest("POST", url, requestBody)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		request.Header.Set("Content-Type", contentType)
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}

		// 处理响应
		if response.StatusCode != http.StatusOK {
			fmt.Println("Server returned non-200 status code:", response.StatusCode)
			return
		}
		fmt.Printf("Uploaded chunk %d/%d\n", i+1, numChunks)

		response.Body.Close()
	}

	fmt.Println("File uploaded successfully!")

	return
}

// Chunked 文件分片
// @param: filepath string文件路径， chunkSize int 分片大小 1024*1024=1M
func Chunked(filepath string, chunkSize int64) error {
	file, err := os.Open(filepath) // 打开要上传的大文件
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer file.Close()

	// 计算文件MD5
	fileMD5, err := md5.GetFileMD5(filepath)
	if err != nil {
		fmt.Println("Error getting file md5:", err)
		return err
	}

	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return err
	}

	// 计算分片数量
	fileSize := fileInfo.Size()
	numChunks := int64(math.Ceil(float64(fileSize) / float64(chunkSize)))

	// 循环读取分片并保存
	for i := int64(0); i < numChunks; i++ {
		chunk := make([]byte, chunkSize)
		n, err := file.ReadAt(chunk, i*chunkSize)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading chunk:", err)
			return err
		}

		outputFilename := fmt.Sprintf("%s_%d", fileChunkDir+fileMD5+"/"+fileInfo.Name(), i)
		err = os.MkdirAll(fileChunkDir+fileMD5, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return err
		}

		outputFile, err := os.Create(outputFilename)
		if err != nil {
			return err
		}

		writer := bufio.NewWriter(outputFile)
		_, err = writer.Write(chunk[:n])
		if err != nil {
			return err
		}
		err = writer.Flush()
		if err != nil {
			fmt.Println("Error writer flush:", err)
			return err
		}

		err = outputFile.Close()
		if err != nil {
			fmt.Println("Error outputFile close:", err)
			return err
		}
	}

	return nil
}

// RemoveChunk 删除分片
// @param: fileMD5 分片文件目录
func RemoveChunk(fileMD5 string) error {
	err := os.RemoveAll(fileChunkDir + fileMD5)
	if err != nil {
		return err
	}

	return nil
}
