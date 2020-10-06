package mycompress

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func ReadFileBufferio(filePath string) []byte {
	fp, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("文件打开失败.", err)
		return nil
	}

	return fp
}

func WriteFileBufferio(filePath string, content []byte) {
	err := ioutil.WriteFile(filePath, content, 0666)
	if err != nil {
		fmt.Println("文件写入错误.", err)
		return
	}
}

func GetFileSize(filePath string) int64 {
	fp, err := os.Open(filePath)
	if err != nil {
		fmt.Println("文件打开失败.", err)
		return 0
	}
	file, _ := fp.Stat()
	return file.Size()
}

func CompareFile(file1, file2 string) bool {
	fp1, err := os.Open(file1)
	if err != nil {
		fmt.Println("file1打开失败.", err)
		return false
	}
	r1, err := getMD5SumString(fp1)
	if err != nil {
		fmt.Println("获取file1 md5失败.", err)
		return false
	}

	fp2, err := os.Open(file2)
	if err != nil {
		fmt.Println("file2打开失败.", err)
		return false
	}
	r2, err := getMD5SumString(fp2)
	if err != nil {
		fmt.Println("获取file2 md5失败.", err)
		return false
	}
	return compareCheckSum(r1, r2)
}

func getMD5SumString(f *os.File) (string, error) {
	file1Sum := md5.New()
	_, err := io.Copy(file1Sum, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", file1Sum.Sum(nil)), nil
}

func compareCheckSum(sum1, sum2 string) bool {
	match := true
	if sum1 != sum2 {
		match = false
	}
	return match
}
