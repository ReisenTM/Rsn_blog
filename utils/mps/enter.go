package mps

import (
	"encoding/json"
	"reflect"
)

// StructToMap 结构体转map
func StructToMap(src any, typ string) (res map[string]interface{}) {
	res = make(map[string]interface{})
	content := reflect.ValueOf(src)
	for i := 0; i < content.NumField(); i++ {
		val := content.Field(i)
		tag := content.Type().Field(i).Tag.Get(typ)
		if tag == "" || tag == "-" {
			continue
		}
		if val.IsNil() {
			continue
		}
		if val.Kind() == reflect.Ptr {
			v1 := val.Elem().Interface()
			if val.Elem().Kind() == reflect.Slice {
				byteData, _ := json.Marshal(v1)
				res[tag] = string(byteData)
			} else {
				res[tag] = v1
			}
			continue
		}
		res[tag] = val.Interface()

	}
	return res
}
