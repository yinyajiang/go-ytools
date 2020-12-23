package tools

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
)

//GenMd5 生成MD5
func GenMd5(data []byte) (string, error) {
	buf := bytes.NewBuffer(data)
	md5g := md5.New()
	if _, err := io.Copy(md5g, buf); err != nil {
		return "", fmt.Errorf("Gen md5 fail,err:%v", err)
	}
	return fmt.Sprintf("%x", md5g.Sum(nil)), nil
}

//GenFileMd5 生成文件md5
func GenFileMd5(path string) (string, error) {
	file, err := OpenReadFile(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	md5g := md5.New()
	if _, err := io.Copy(md5g, file); err != nil {
		return "", fmt.Errorf("Gen md5 fail file:%s,err:%v", path, err)
	}
	return fmt.Sprintf("%x", md5g.Sum(nil)), err
}

//EncodeBase64 生成Base64
func EncodeBase64(data []byte) string {
	encoding := base64.StdEncoding.EncodeToString(data)
	return encoding
}

//DecodeBase64 解码Base64
func DecodeBase64(base string) []byte {
	dncoding, _ := base64.StdEncoding.DecodeString(base)
	return dncoding
}
