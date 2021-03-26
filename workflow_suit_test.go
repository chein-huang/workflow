package workflow_test

import (
	"context"
	"testing"

	"github.com/chein-huang/workflow"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWorkFlow(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WorkFlow Suite")
}

type TestStruct struct {
	workBegin        workflow.Event
	workBeforeCommit workflow.Event
	workCommit       workflow.Event
	workRollback     workflow.Event
	workFinish       workflow.Event
}

func (s *TestStruct) WorkBegin(ctx context.Context, data *workflow.WorkData) error {
	if s.workBegin != nil {
		return s.workBegin(ctx, data)
	}
	return nil
}

func (s *TestStruct) WorkBeforeCommit(ctx context.Context, data *workflow.WorkData) error {
	if s.workBeforeCommit != nil {
		return s.workBeforeCommit(ctx, data)
	}
	return nil
}

func (s *TestStruct) WorkCommit(ctx context.Context, data *workflow.WorkData) error {
	if s.workCommit != nil {
		return s.workCommit(ctx, data)
	}
	return nil
}

func (s *TestStruct) WorkRollback(ctx context.Context, data *workflow.WorkData) error {
	if s.workRollback != nil {
		return s.workRollback(ctx, data)
	}
	return nil
}

func (s *TestStruct) WorkFinish(ctx context.Context, data *workflow.WorkData) error {
	if s.workFinish != nil {
		return s.workFinish(ctx, data)
	}
	return nil
}
