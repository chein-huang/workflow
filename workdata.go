package workflow

import (
	"context"
	"fmt"
	"sync"
)

type WorkData struct {
	workBegin        *funcs
	workBeforeCommit *funcs
	workCommit       *funcs
	workFinish       *funcs
	workRollback     *funcs
	currentState     *funcs
	ctx              map[string]interface{}
	ctxLocker        sync.RWMutex
	isAborted        bool
}

func NewWorkData() *WorkData {
	d := &WorkData{
		workBegin:        NewFuncs(),
		workBeforeCommit: NewFuncs(),
		workCommit:       NewFuncs(),
		workFinish:       NewFuncs(),
		workRollback:     NewFuncs(),
		ctx:              make(map[string]interface{}),
	}
	d.ResetProgress()
	return d
}

func (d *WorkData) Get(key string) (interface{}, bool) {
	d.ctxLocker.RLock()
	defer d.ctxLocker.RUnlock()
	data, ok := d.ctx[key]
	return data, ok
}

func (d *WorkData) MustGet(key string) interface{} {
	if data, ok := d.Get(key); !ok {
		panic(fmt.Sprintf("key: %v is not exists", key))
	} else {
		return data
	}
}

func (d *WorkData) Set(key string, value interface{}) {
	d.ctxLocker.Lock()
	defer d.ctxLocker.Unlock()
	d.ctx[key] = value
}

func (d *WorkData) Next(ctx context.Context) error {
	if d.currentState == nil {
		return nil
	}

	defer func() {
		d.currentState = nil
	}()

	return d.currentState.next(ctx, d)
}

func (d *WorkData) Abort() {
	d.isAborted = true
	if d.currentState == nil {
		return
	}

	d.currentState.abort()
}

func (d *WorkData) IsAborted() bool {
	return (d.currentState != nil && d.currentState.isAborted()) || d.isAborted
}

func (d *WorkData) ResetProgress() {
	d.workBegin.reset()
	d.workBeforeCommit.reset()
	d.workCommit.reset()
	d.workRollback.reset()
	d.workFinish.reset()
}

func (d *WorkData) Begin(ctx context.Context) error {
	d.currentState = d.workBegin
	return d.Next(ctx)
}

func (d *WorkData) BeforeCommit(ctx context.Context) error {
	d.currentState = d.workBeforeCommit
	return d.Next(ctx)
}

func (d *WorkData) Commit(ctx context.Context) error {
	d.currentState = d.workCommit
	return d.Next(ctx)
}

func (d *WorkData) Rollback(ctx context.Context) error {
	d.currentState = d.workRollback
	return d.Next(ctx)
}

func (d *WorkData) Finish(ctx context.Context) error {
	d.currentState = d.workFinish
	return d.Next(ctx)
}
