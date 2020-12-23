package tools

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
)

//DownFile 下载文件
func DownFile(url, path string) error {
	return DownFileFun(url, path, nil)
}

//DownFileFun 下载文件
func DownFileFun(url, path string, progHand ProgressHand) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Get url fail,url:%s,err:%v", url, err)
	}

	var fsize int64
	if len(resp.Header.Values("Content-Length")) > 0 {
		fsize, _ = strconv.ParseInt(resp.Header.Values("Content-Length")[0], 10, 64)
	}

	defer resp.Body.Close()
	file, err := CreateFile(path)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = CopyFun(fsize, file, resp.Body, progHand)
	if err != nil {
		os.Remove(path)
		return fmt.Errorf("Copy body fail,url:%s,err:%v", url, err)
	}
	return nil
}

//IsValidURL  url是否有效
func IsValidURL(url string) bool {
	res, err := http.Get(url)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	return 200 == res.StatusCode
}

//URLHeaderAttr url返回的数据头的属性值
func URLHeaderAttr(url, name string) string {
	resp, err := http.Get(url)
	if err != nil {
		return ``
	}
	return resp.Header.Values(name)[0]
}
