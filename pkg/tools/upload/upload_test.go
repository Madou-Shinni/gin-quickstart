package upload

import (
	"fmt"
	"github.com/Madou-Shinni/gin-quickstart/pkg/tools/md5"
	"testing"
)

func TestChunked(t *testing.T) {
	Chunked("C:\\Users\\sni\\Downloads\\寻美中国 ·浙里是宁波.mp4", 10*1024*1024)
}

func TestMergeChunk(t *testing.T) {
	chunk, err := MergeChunk("18667b1d3e1b7a366cabe5d4c9108072", "寻美中国 ·浙里是宁波.mp4")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(chunk)
}

func TestFileMd5(t *testing.T) {
	file1MD5, _ := md5.GetFileMD5("C:\\Users\\sni\\Downloads\\寻美中国 ·浙里是宁波.mp4")
	file2MD5, _ := md5.GetFileMD5("D:\\go-project\\gin-quickstart\\pkg\\tools\\upload\\file\\寻美中国 ·浙里是宁波.mp4")

	fmt.Printf("f1md5: %v\n", file1MD5)
	fmt.Printf("f2md5: %v\n", file2MD5)
	fmt.Printf("f1md5 == f2md5: %v", file1MD5 == file2MD5)
}

func TestRemoveChunk(t *testing.T) {
	RemoveChunk("18667b1d3e1b7a366cabe5d4c9108072")
}
