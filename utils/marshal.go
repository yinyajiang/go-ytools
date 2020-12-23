package tools

import "encoding/json"

//MarshalToJSONFile 序列化到文件
func MarshalToJSONFile(s interface{}, path string) error {
	file, err := CreateFile(path)
	if err != nil {
		return err
	}
	jsonMarsh1, err := json.Marshal(s)
	if err != nil {
		return err
	}
	file.Write(jsonMarsh1)
	return nil
}

//UnMarshalJSONFile 反序列化JSON
func UnMarshalJSONFile(path string, out interface{}) error {
	data, err := ReadFileAll(path)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, out); err != nil {
		return err
	}
	return nil
}
