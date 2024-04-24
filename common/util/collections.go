package util

// list转map
func ListObjToMap[T, V any, K comparable](slice []T, f func(obj T) (K, V)) map[K]V {
	objMap := make(map[K]V, len(slice))
	for _, v := range slice {
		key, val := f(v)
		objMap[key] = val
	}
	return objMap
}

// list转map[][]obj
func ListObjToMapList[T, V any, K comparable](slice []T, f func(obj T) (K, V)) map[K][]V {
	objMap := make(map[K][]V, len(slice))
	for _, v := range slice {
		key, val := f(v)
		objMap[key] = append(objMap[key], val)
	}
	return objMap
}

// map转list
func MapToList[T, V any, K comparable](curMap map[K]V, f func(k K, v V) T) []T {
	list := make([]T, 0)
	for k, v := range curMap {
		list = append(list, f(k, v))
	}
	return list
}

// list转换其他类型的list
func ListObjToListObj[T any, U any](slice []T, f func(obj T) U) []U {
	list := make([]U, len(slice))
	for i, v := range slice {
		list[i] = f(v)
	}
	return list
}

// list元素去重
func RemoveRepeatFromList[T comparable](slice []T) []T {
	itemMap := make(map[T]bool, len(slice))
	for _, v := range slice {
		itemMap[v] = true
	}
	list := []T{}
	for k, _ := range itemMap {
		list = append(list, k)
	}
	return list
}

// listObj元素去重
func RemoveRepeatFromListObj[T any, U comparable](slice []T, f func(obj T) U) []T {
	itemMap := make(map[U]T, len(slice))
	for _, v := range slice {
		itemMap[f(v)] = v
	}
	list := []T{}
	for _, v := range itemMap {
		list = append(list, v)
	}
	return list
}

// 拷贝数组，
func CopyList[T any](src []T) []T {
	list := make([]T, len(src))
	for i, v := range src {
		list[i] = v
	}
	return list
}

// 拷贝map
func CopyMap[T comparable, U any](src map[T]U) map[T]U {
	dtsMap := make(map[T]U, len(src))
	for k, v := range src {
		dtsMap[k] = v
	}
	return dtsMap
}

// 拷贝Chan，暂时不考虑使用
func CopyChan[T any](src chan T) chan T {
	dtsChan := make(chan T, len(src))
	for v := range src {
		dtsChan <- v
	}
	return dtsChan
}

// 从map中获取keys
func GetKeysFromMap[T comparable](aMap map[T]interface{}) []T {
	list := make([]T, len(aMap))
	i := 0
	for k := range aMap {
		list[i] = k
		i++
	}
	return list
}

// 从map中获取keys
func GetValuesFromMap[K comparable, V any](aMap map[K]V) []V {
	list := make([]V, len(aMap))
	i := 0
	for _, v := range aMap {
		list[i] = v
		i++
	}
	return list
}

// 根据某个值去重后的list
func ListToDeduplicationList[T, V any, K comparable](list []T, f func(t T) (K, V)) []V {
	temMap := make(map[K]V, 0)
	for _, t := range list {
		k, v := f(t)
		temMap[k] = v
	}
	temList := make([]V, 0)
	for _, v := range temMap {
		temList = append(temList, v)
	}
	return temList
}

// 获取两个集合的交集、差集1，差集2
func GetThreeSetFromList[T, V any, K comparable](list1, list2 []T, f func(t T) (K, V)) (set1 []V, set2 []V, set3 []V) {
	temMap1 := make(map[K]V)
	for _, t := range list1 {
		k, v := f(t)
		temMap1[k] = v
	}
	temMap2 := make(map[K]V)
	for _, t := range list2 {
		k, v := f(t)
		temMap2[k] = v
	}
	for k1, v1 := range temMap1 {
		if _, ok := temMap2[k1]; ok {
			set2 = append(set2, v1)
		} else {
			set1 = append(set1, v1)
		}
	}
	for k2, v2 := range temMap2 {
		if _, ok := temMap1[k2]; !ok {
			set3 = append(set3, v2)
		}
	}
	return
}
