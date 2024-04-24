package util

import "sync"

type wait struct{}

type waitGroup struct {
	sync.Mutex
	wg      sync.WaitGroup
	chanNum chan wait
	err     error
}

// 0:表示不限制协程数
func NewWaitGroup(maxNum int) *waitGroup {
	if maxNum > 0 {
		return &waitGroup{
			chanNum: make(chan wait, maxNum),
		}
	}
	return &waitGroup{}
}

func (w *waitGroup) Go(f func() (err error)) {
	w.Add()
	go func() {
		defer w.Done()
		err := f()
		if err != nil {
			w.err = err
		}
	}()
}

func (w *waitGroup) Add() {
	w.wg.Add(1)
	if w.chanNum != nil {
		w.chanNum <- wait{}
	}
}

func (w *waitGroup) Done() {
	w.wg.Done()
	if w.chanNum != nil {
		<-w.chanNum
	}
}

func (w *waitGroup) Wait() error {
	w.wg.Wait()
	return w.err
}

func (w *waitGroup) SetError(err error) {
	w.err = err
	return
}
