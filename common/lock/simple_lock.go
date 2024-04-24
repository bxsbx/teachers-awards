package lock

import "sync"

// 非阻塞锁，由于是单服务，所以不需要加版本号（锁的添加删除是串行执行，且锁是保存在内存中的，服务宕机了也不怕）
type LockSet struct {
	sync.Mutex
	keySet map[string]struct{}
}

type FuncLockMap struct {
	sync.Mutex
	funcMap map[string]*LockSet
}

func NewFuncLockMap() *FuncLockMap {
	return &FuncLockMap{
		funcMap: make(map[string]*LockSet),
	}
}

func (f *FuncLockMap) GetFunMap(key string) *LockSet {
	f.Lock()
	defer f.Unlock()
	if temMap, ok := f.funcMap[key]; ok {
		return temMap
	} else {
		newFunMap := &LockSet{
			keySet: make(map[string]struct{}),
		}
		f.funcMap[key] = newFunMap
		return newFunMap
	}
}

func (l *LockSet) SetKey(key string) bool {
	l.Lock()
	defer l.Unlock()
	if _, ok := l.keySet[key]; ok {
		return false
	} else {
		l.keySet[key] = struct{}{}
		return true
	}
}

func (l *LockSet) DelKey(key string) {
	l.Lock()
	defer l.Unlock()
	delete(l.keySet, key)
}
