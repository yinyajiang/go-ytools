package tools

import (
	"bytes"
	"fmt"
	"io"
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

//GetNetContent ...
func GetNetContent(url string) (ret []byte, e error) {
	resp, err := http.Get(url)
	if err != nil {
		e = err
		return
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		e = err
		return
	}
	ret = buf.Bytes()
	return
}

//SendNetRequest ...
func SendNetRequest(method, url string, head map[string]string, body io.Reader, recv io.Writer) (http.Header, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for key, val := range head {
		req.Header.Add(key, val)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if recv != nil {
		_, err = io.Copy(recv, resp.Body)
		if err != nil && err != io.EOF {
			return nil, err
		}
	}
	return resp.Header.Clone(), nil
}
