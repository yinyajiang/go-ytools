package tools

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//ThePath 指定路径
func ThePath(root string, path ...string) string {
	root = AbsPath(root)
	tmp := append([]string{}, root)
	path = append(tmp, path...)
	return AbsJoinPath(path...)
}

//AbsParent 绝对父路径
func AbsParent(path string) string {
	return filepath.Dir(AbsPath(path))
}

//AbsPath 绝对路径
func AbsPath(path string) string {
	return AbsJoinPath(path)
}

//AbsJoinPath 拼接路径
func AbsJoinPath(paths ...string) string {
	if 0 == len(paths) {
		return ""
	}
	abs, err := filepath.Abs(paths[0])
	if err != nil {
		if strings.HasPrefix(paths[0], "./") || strings.HasPrefix(paths[0], ".\\") {
			return "./" + filepath.Join(paths...)
		}
	}
	paths[0] = abs
	return filepath.Join(paths...)
}

//PathName 从url或路径中获取文件名
func PathName(path string) string {
	index1 := strings.LastIndex(path, "/")
	index2 := strings.LastIndex(path, "\\")
	if -1 == index1 && -1 == index2 {
		return path
	}

	index := index1
	if index2 > index1 {
		index = index2
	}
	return path[index+1:]
}

//PathStem 从url或路径中获取不带dot的文件名
func PathStem(path string) string {
	dotname := PathName(path)
	index := strings.LastIndex(dotname, ".")
	if -1 != index {
		return dotname[:index]
	}
	return dotname
}

//PathAvailableSpace 可用空间
func PathAvailableSpace(path string) uint64 {
	usage, err := DiskUsage(path)
	if err != nil {
		return 0
	}
	return usage.Free
}

//ReplacePath 替换\为/
func ReplacePath(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

//GetHasFileRoot 找到第一个非空路径
func GetHasFileRoot(root string) (ret string) {
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.HasPrefix(info.Name(), ".") {
			return nil
		}
		if !info.IsDir() {
			ret = AbsParent(path)
			return fmt.Errorf("found")
		}
		return nil
	})
	return
}

//GetExtFilePath 从目录中获取指定后缀的文件路径
func GetExtFilePath(dirpath, ext string) string {
	ext = strings.ToLower(ext)
	dirpath = AbsPath(dirpath)
	rd, err := ioutil.ReadDir(dirpath)
	if nil != err {
		return ""
	}
	for _, fi := range rd {
		path := AbsJoinPath(dirpath, fi.Name())
		if fi.IsDir() {
			return GetExtFilePath(path, ext)
		} else if ext == strings.ToLower(filepath.Ext(path)) {
			return path
		}
	}
	return ""
}

//FilterFile 从目录中过滤文件,不会深度遍历
func FilterFile(dirpath string, filter []string) []string {
	dirpath = AbsPath(dirpath)
	rd, err := ioutil.ReadDir(dirpath)
	if nil != err {
		return []string{}
	}
	results := []string{}
	for _, fi := range rd {
		if fi.IsDir() {
			continue
		}
		if IsInFilter(fi.Name(), filter) {
			path := AbsJoinPath(dirpath, fi.Name())
			results = append(results, path)
		}
	}
	return results
}

//FilterDeepFile 从目录中过滤文件,会深度遍历
func FilterDeepFile(dirpath string, filter []string) []string {
	results := []string{}
	PathWalk(dirpath, func(path string, info os.FileInfo, postName string) error {
		if info.IsDir() {
			return nil
		}
		if IsInFilter(info.Name(), filter) {
			results = append(results, path)
		}
		return nil
	})
	return results
}

//IsInFilter 文件格式或名字是否在筛选器里面
func IsInFilter(file string, filter []string) bool {
	file = PathName(file)
	for _, pattern := range filter {
		b, _ := filepath.Match(pattern, file)
		if b {
			return b
		}
	}
	return false
}

//GetApplicationDir 获取当前程序目录
func GetApplicationDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return ""
	}
	return dir
}

//IsEqualPath 是否是相同路径
func IsEqualPath(p1, p2 string) bool {
	if strings.HasSuffix(p1, ":") {
		p1 += "/"
	}
	if strings.HasSuffix(p2, ":") {
		p2 += "/"
	}

	p1 = AbsJoinPath(p1, "t")
	p2 = AbsJoinPath(p2, "t")

	p1 = ReplacePath(p1)
	p2 = ReplacePath(p2)

	p1 = strings.ToLower(p1)
	p2 = strings.ToLower(p2)
	return p1 == p2
}

//PathWalk 遍历目录，回调带上去除根路径的name
func PathWalk(path string, f func(p string, info os.FileInfo, postName string) error) {
	pathlen := len(path)
	filepath.Walk(path, func(fpath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if len(fpath) <= pathlen {
			return nil
		}
		if info.Name() == "." || info.Name() == ".." {
			return nil
		}
		postName := PostPath(fpath, path)
		return f(fpath, info, postName)
	})
}

//SpecialDirType 特定目录
type SpecialDirType int32

var (
	//LocalAppdata 数据路径
	LocalAppdata SpecialDirType = 0x001c
)

//LocalPath 当前程序文件路径
func LocalPath(path string) string {
	return AbsJoinPath(AbsParent(os.Args[0]), path)
}

//PostPath full去除src的剩余路径路径
func PostPath(full, src string) string {
	i := strings.Index(full, src)
	if i == -1 {
		return ""
	}
	i = len(src)
	if '\\' == full[i] || '/' == full[i] {
		i++
	}
	return full[i:]
}

//TempPath 临时目录文件路径
func TempPath(path string) string {
	return AbsJoinPath(os.TempDir(), path)
}

//DataPath 数据目录文件路径
func DataPath(path string) string {
	return AbsJoinPath(GetSpecialDir(LocalAppdata), "Library/Preferences", ProductName, path)
}
