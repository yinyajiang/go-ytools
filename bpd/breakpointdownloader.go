package bpd

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	tools "github.com/yinyajiang/go-ytools/utils"
)

//BPDownloader 断点下载器
type BPDownloader struct {
	reader io.ReadCloser
}

//New ...
func New() *BPDownloader {
	return &BPDownloader{}
}

//Stop ...
func (b *BPDownloader) Stop() {
	if nil != b.reader {
		b.reader.Close()
	}
}

//Download ...
func (b *BPDownloader) Download(url, path string, progFun tools.ProgressHand) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	resp.Body.Close()

	err = tools.CreateDirs(tools.AbsParent(path))
	if err != nil {
		return
	}

	if !strings.EqualFold(resp.Header.Get("Accept-Ranges"), "bytes") {
		return fmt.Errorf("No support Accept-Ranges")
	}
	if resp.ContentLength <= 0 {
		return fmt.Errorf("Not content length")
	}

	size := resp.ContentLength
	editTag := resp.Header.Get("ETag") + resp.Header.Get("Last-modified")
	editTag, _ = tools.GenMd5([]byte(editTag))
	etagPath := path + ".yetag"

	var offset int64
	if tools.IsExist(path) &&
		tools.IsExist(etagPath) {
		if editTag != tools.ReadFileString(etagPath) {
			fmt.Println("Last edittag is not equal")
		} else {
			offset = tools.FileSize(path)
		}
	}

	if offset == 0 {
		tools.RemovePath(etagPath)
		tools.RemovePath(path)
	} else if offset == size {
		return nil
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", offset, size-1))
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	file, err := tools.OpenApptendFile(path)
	if err != nil {
		return
	}

	tools.WriteFileString(etagPath, editTag)
	b.reader = resp.Body
	_, err = tools.CopyFun(size-offset, file, resp.Body, func(total int64, prog float64) {
		if nil != progFun {
			if prog != 1.0 {
				prog = (float64(total)*prog + float64(offset)) / float64(size)
			}
			if prog > 1.0 {
				prog = 1.0
			}
			progFun(size, (float64(total)*prog+float64(offset))/float64(size))
		}
	})
	return
}
