package workflow

import (
	"context"
	"fmt"
)

func StartWorkFlow(f Event, opts ...Options) (*WorkData, error) {
	return StartWorkFlowContext(nil, f, opts...)
}

func StartWorkFlowContext(ctx context.Context, f Event, opts ...Options) (data *WorkData, err error) {
	if ctx == nil {
		ctx = context.Background()
	}
	data = NewWorkData()
	for _, opt := range opts {
		opt.Apply(data)
	}

	err = data.Begin(ctx)
	if err != nil || data.IsAborted() {
		return
	}

	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}

		if err != nil || data.IsAborted() {
			re := data.Rollback(ctx)
			if re != nil {
				err = re
				return
			}
			return
		}

		if data.IsAborted() {
			return
		}
		data.Finish(ctx)
	}()

	err = f(ctx, data)
	if err != nil || data.IsAborted() {
		return
	}

	err = data.BeforeCommit(ctx)
	if err != nil || data.IsAborted() {
		return
	}

	err = data.Commit(ctx)
	if err != nil || data.IsAborted() {
		return
	}

	return
}
