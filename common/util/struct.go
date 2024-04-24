package util

import (
	bytes2 "bytes"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
)

func structToMap(values reflect.Value, keys reflect.Type, moreMap ...bool) map[string]interface{} {
	numField := keys.NumField()
	objMap := make(map[string]interface{}, numField)
	for i := 0; i < numField; i++ {
		k := keys.Field(i).Name
		v := values.Field(i)
		if len(moreMap) > 0 && moreMap[0] && v.Type().Kind() == reflect.Struct {
			objMap[k] = structToMap(v, v.Type())
		} else {
			objMap[k] = v
		}
	}
	return objMap
}

// 结构体转map，禁止循环引用
func StructToMap(obj any, moreMap ...bool) (map[string]interface{}, error) {
	keys := reflect.TypeOf(obj)
	values := reflect.ValueOf(obj)
	if keys.Kind() != reflect.Struct {
		return nil, errors.New("不支持非结构类型")
	}
	return structToMap(values, keys, moreMap...), nil
}

// 引用类型 - slice、map、channel、interface、function

// 使用反射进行对象之间同属性同名称复制，反射仅支持结构体 (flag true-深拷贝，false-浅拷贝)，默认浅拷贝
func ObjToObjByReflect[T any, U any](a *T, b *U, flag ...bool) {
	aElem := reflect.ValueOf(a).Elem()
	bElem := reflect.ValueOf(b).Elem()

	aType := aElem.Type()
	bType := bElem.Type()

	if aType.Kind() != reflect.Struct || bType.Kind() != reflect.Struct {
		return
	}

	aFieldNum := aElem.NumField()
	bFieldNum := bElem.NumField()

	aMap := make(map[string]reflect.Value, aFieldNum)
	for i := 0; i < aFieldNum; i++ {
		aFieldName := aType.Field(i).Name
		aMap[aFieldName] = aElem.Field(i)
	}

	for i := 0; i < bFieldNum; i++ {
		bField := bElem.Field(i)
		bFieldName := bType.Field(i).Name
		if v, ok := aMap[bFieldName]; ok {
			if v.Type() == bField.Type() && bField.CanSet() {
				if len(flag) > 0 && flag[0] {
					bField.Set(copyDfs(v, true))
				} else {
					bField.Set(v)
				}
			}
		}
	}
}

func copyDfs(srcVal reflect.Value, flag bool) reflect.Value {
	srcType := srcVal.Type()

	if srcType.Kind() == reflect.Ptr {
		srcVal = srcVal.Elem()
	}

	switch srcType.Kind() {
	case reflect.Slice:
		slice := reflect.MakeSlice(srcType, 0, 0)
		for i := 0; i < srcVal.Len(); i++ {
			if flag {
				slice = reflect.Append(slice, copyDfs(srcVal.Index(i), true))
			} else {
				slice = reflect.Append(slice, srcVal.Index(i))
			}
		}
		return slice
	case reflect.Map:
		dstMap := reflect.MakeMap(srcType)
		for _, key := range srcVal.MapKeys() {
			if flag {
				dstMap.SetMapIndex(key, copyDfs(srcVal.MapIndex(key), true))
			} else {
				dstMap.SetMapIndex(key, srcVal.MapIndex(key))
			}
		}
		return dstMap
	case reflect.Struct:
		dst := reflect.New(srcType).Elem()
		allFieldPrivate := true
		for i := 0; i < srcVal.NumField(); i++ {
			field := dst.Field(i)
			if field.CanSet() {
				if flag {
					field.Set(copyDfs(srcVal.Field(i), true))
				} else {
					field.Set(srcVal.Field(i))
				}
				allFieldPrivate = false
			}
		}
		// 所有字段都是私有的，浅拷贝，特殊处理
		if allFieldPrivate {
			dst.Set(srcVal)
		}
		return dst
	default:
		dst := reflect.New(srcType).Elem()
		dst.Set(srcVal)
		return dst
	}
}

// 使用反射进行拷贝(flag true-深拷贝，false-浅拷贝)，默认浅拷贝
func CopyObj[T any](a *T, flag ...bool) T {
	elem := reflect.ValueOf(a).Elem()
	var dts reflect.Value
	if len(flag) > 0 && flag[0] {
		dts = copyDfs(elem, true)
	} else {
		dts = copyDfs(elem, false)
	}
	return dts.Interface().(T)
}

// 使用json序列化进行结构体对象的复制操作，深拷贝，支持所有类型
func ObjToObjByJson(a interface{}, b interface{}) error {
	buf := new(bytes2.Buffer)
	encoder := json.NewEncoder(buf)
	decoder := json.NewDecoder(buf)

	if err := encoder.Encode(a); err != nil {
		return err
	}
	if err := decoder.Decode(b); err != nil {
		return err
	}
	return nil
}

// json数据过大时此方法，简单且不大的数据用json.Unmarshal，否则可能造成内存溢出
func JsonToStruct(data string, obj interface{}) error {
	reader := strings.NewReader(data)
	dec := json.NewDecoder(reader)
	dec.UseNumber()
	err := dec.Decode(obj)
	return err
}
