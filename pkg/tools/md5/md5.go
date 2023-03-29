package md5

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

// 返回一个32位md5加密后的字符串
func Md532(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 返回一个16位md5加密后字符串
func Md5To16(str string) string {
	return Md532(str)[8:24]
}

// GetFileMD5 计算文件的 MD5 值
// @param: path string文件路径 /data/file1.png
func GetFileMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := md5.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	md5sum := hex.EncodeToString(hasher.Sum(nil))

	return md5sum, nil
}
