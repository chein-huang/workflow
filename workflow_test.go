package workflow_test

import (
	"context"
	"fmt"

	"github.com/chein-huang/workflow"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	index = 0
)

var _ = Describe("Workflow", func() {
	index := 1

	BeforeEach(func() {
		index = 1
	})

	begin := successFunc("begin", &index)
	beforeCommit := successFunc("beforeCommit", &index)
	commit := successFunc("commit", &index)
	rollback := successFunc("rollback", &index)
	finish := successFunc("finish", &index)

	s := TestStruct{
		workBegin:        begin,
		workBeforeCommit: beforeCommit,
		workCommit:       commit,
		workRollback:     rollback,
		workFinish:       finish,
	}

	Context("success cases", func() {
		It("with interface", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithInterface(&s),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(3))
			Expect(data.MustGet("commit")).To(Equal(4))
			Expect(data.MustGet("finish")).To(Equal(5))
			Expect(data.Get("rollback")).Should(BeNil())
		})

		It("before begin success", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithBegin(begin),
				workflow.WithInterface(&s),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(2))
			Expect(data.MustGet("work")).To(Equal(3))
			Expect(data.MustGet("beforeCommit")).To(Equal(4))
			Expect(data.MustGet("commit")).To(Equal(5))
			Expect(data.MustGet("finish")).To(Equal(6))
			Expect(data.Get("rollback")).Should(BeNil())
		})

		It("after begin success", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithInterface(&s),
				workflow.WithBegin(begin),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(2))
			Expect(data.MustGet("work")).To(Equal(3))
			Expect(data.MustGet("beforeCommit")).To(Equal(4))
			Expect(data.MustGet("commit")).To(Equal(5))
			Expect(data.MustGet("finish")).To(Equal(6))
			Expect(data.Get("rollback")).Should(BeNil())
		})

		It("before beforeCommit success", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithBeforeCommit(beforeCommit),
				workflow.WithInterface(&s),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(4))
			Expect(data.MustGet("commit")).To(Equal(5))
			Expect(data.MustGet("finish")).To(Equal(6))
			Expect(data.Get("rollback")).Should(BeNil())
		})

		It("after beforeCommit success", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithInterface(&s),
				workflow.WithBeforeCommit(beforeCommit),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(4))
			Expect(data.MustGet("commit")).To(Equal(5))
			Expect(data.MustGet("finish")).To(Equal(6))
			Expect(data.Get("rollback")).Should(BeNil())
		})

		It("before commit success", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithCommit(commit),
				workflow.WithInterface(&s),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(3))
			Expect(data.MustGet("commit")).To(Equal(5))
			Expect(data.MustGet("finish")).To(Equal(6))
			Expect(data.Get("rollback")).Should(BeNil())
		})

		It("after commit success", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithInterface(&s),
				workflow.WithCommit(commit),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(3))
			Expect(data.MustGet("commit")).To(Equal(5))
			Expect(data.MustGet("finish")).To(Equal(6))
			Expect(data.Get("rollback")).Should(BeNil())
		})

		It("before finish success", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithFinish(finish),
				workflow.WithInterface(&s),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(3))
			Expect(data.MustGet("commit")).To(Equal(4))
			Expect(data.MustGet("finish")).To(Equal(6))
			Expect(data.Get("rollback")).Should(BeNil())
		})

		It("after finish success", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithInterface(&s),
				workflow.WithFinish(finish),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(3))
			Expect(data.MustGet("commit")).To(Equal(4))
			Expect(data.MustGet("finish")).To(Equal(6))
			Expect(data.Get("rollback")).Should(BeNil())
		})

		It("before rollback success", func() {
			data, err := workflow.StartWorkFlow(
				failedFunc("work", &index),
				workflow.WithRollback(rollback),
				workflow.WithInterface(&s),
			)
			Expect(err).ShouldNot(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("rollback")).To(Equal(4))
			Expect(data.Get("beforeCommit")).Should(BeNil())
			Expect(data.Get("commit")).Should(BeNil())
			Expect(data.Get("finish")).Should(BeNil())
		})

		It("after rollback success", func() {
			data, err := workflow.StartWorkFlow(
				failedFunc("work", &index),
				workflow.WithInterface(&s),
				workflow.WithRollback(rollback),
			)
			Expect(err).ShouldNot(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("rollback")).To(Equal(4))
			Expect(data.Get("beforeCommit")).Should(BeNil())
			Expect(data.Get("commit")).Should(BeNil())
			Expect(data.Get("finish")).Should(BeNil())
		})
	})

	beginFailed := failedFunc("begin", &index)
	beforeCommitFailed := failedFunc("beforeCommit", &index)
	commitFailed := failedFunc("commit", &index)
	rollbackFailed := failedFunc("rollback", &index)
	finishFailed := failedFunc("finish", &index)

	Context("failed cases", func() {
		It("before begin failed", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithBegin(beginFailed),
				workflow.WithInterface(&s),
			)
			Expect(err).ShouldNot(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.Get("work")).Should(BeNil())
			Expect(data.Get("beforeCommit")).Should(BeNil())
			Expect(data.Get("commit")).Should(BeNil())
			Expect(data.Get("rollback")).Should(BeNil())
			Expect(data.Get("finish")).Should(BeNil())
		})

		It("after begin failed", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithInterface(&s),
				workflow.WithBegin(beginFailed),
			)
			Expect(err).ShouldNot(BeNil())
			Expect(data.MustGet("begin")).To(Equal(2))
			Expect(data.Get("work")).Should(BeNil())
			Expect(data.Get("beforeCommit")).Should(BeNil())
			Expect(data.Get("commit")).Should(BeNil())
			Expect(data.Get("rollback")).Should(BeNil())
			Expect(data.Get("finish")).Should(BeNil())
		})

		It("work failed", func() {
			data, err := workflow.StartWorkFlow(
				failedFunc("work", &index),
				workflow.WithInterface(&s),
			)
			Expect(err).ShouldNot(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("rollback")).To(Equal(3))
			Expect(data.Get("beforeCommit")).Should(BeNil())
			Expect(data.Get("commit")).Should(BeNil())
			Expect(data.Get("finish")).Should(BeNil())
		})

		It("before beforeCommit failed", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithBeforeCommit(beforeCommitFailed),
				workflow.WithInterface(&s),
			)
			Expect(err).ShouldNot(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(3))
			Expect(data.MustGet("rollback")).To(Equal(4))
			Expect(data.Get("commit")).Should(BeNil())
			Expect(data.Get("finish")).Should(BeNil())
		})

		It("after beforeCommit failed", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithInterface(&s),
				workflow.WithBeforeCommit(beforeCommitFailed),
			)
			Expect(err).ShouldNot(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(4))
			Expect(data.MustGet("rollback")).To(Equal(5))
			Expect(data.Get("commit")).Should(BeNil())
			Expect(data.Get("finish")).Should(BeNil())
		})

		It("before commit failed", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithCommit(commitFailed),
				workflow.WithInterface(&s),
			)
			Expect(err).ShouldNot(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(3))
			Expect(data.MustGet("commit")).To(Equal(4))
			Expect(data.MustGet("rollback")).To(Equal(5))
			Expect(data.Get("finish")).Should(BeNil())
		})

		It("after commit failed", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithInterface(&s),
				workflow.WithCommit(commitFailed),
			)
			Expect(err).ShouldNot(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(3))
			Expect(data.MustGet("commit")).To(Equal(5))
			Expect(data.MustGet("rollback")).To(Equal(6))
			Expect(data.Get("finish")).Should(BeNil())
		})

		It("before rollback failed", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithCommit(commitFailed),
				workflow.WithRollback(rollbackFailed),
				workflow.WithInterface(&s),
			)
			Expect(err).ShouldNot(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(3))
			Expect(data.MustGet("commit")).To(Equal(4))
			Expect(data.MustGet("rollback")).To(Equal(5))
			Expect(data.Get("finish")).Should(BeNil())
		})

		It("after rollback failed", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithCommit(commitFailed),
				workflow.WithInterface(&s),
				workflow.WithRollback(rollbackFailed),
			)
			Expect(err).ShouldNot(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(3))
			Expect(data.MustGet("commit")).To(Equal(4))
			Expect(data.MustGet("rollback")).To(Equal(6))
			Expect(data.Get("finish")).Should(BeNil())
		})

		It("before finish failed", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithFinish(finishFailed),
				workflow.WithInterface(&s),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(3))
			Expect(data.MustGet("commit")).To(Equal(4))
			Expect(data.MustGet("finish")).To(Equal(5))
			Expect(data.Get("rollback")).Should(BeNil())
		})

		It("after finish failed", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithInterface(&s),
				workflow.WithFinish(finishFailed),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(3))
			Expect(data.MustGet("commit")).To(Equal(4))
			Expect(data.MustGet("finish")).To(Equal(6))
			Expect(data.Get("rollback")).Should(BeNil())
		})
	})

	beginAborted := abortedFunc("begin", &index)
	beforeCommitAborted := abortedFunc("beforeCommit", &index)
	commitAborted := abortedFunc("commit", &index)
	rollbackAborted := abortedFunc("rollback", &index)
	finishAborted := abortedFunc("finish", &index)

	Context("aborted cases", func() {
		It("before begin aborted", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithBegin(beginAborted),
				workflow.WithInterface(&s),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.Get("work")).Should(BeNil())
			Expect(data.Get("beforeCommit")).Should(BeNil())
			Expect(data.Get("commit")).Should(BeNil())
			Expect(data.Get("rollback")).Should(BeNil())
			Expect(data.Get("finish")).Should(BeNil())
		})

		It("after begin aborted", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithInterface(&s),
				workflow.WithBegin(beginAborted),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(2))
			Expect(data.Get("work")).Should(BeNil())
			Expect(data.Get("beforeCommit")).Should(BeNil())
			Expect(data.Get("commit")).Should(BeNil())
			Expect(data.Get("rollback")).Should(BeNil())
			Expect(data.Get("finish")).Should(BeNil())
		})

		It("work aborted", func() {
			data, err := workflow.StartWorkFlow(
				abortedFunc("work", &index),
				workflow.WithInterface(&s),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("rollback")).To(Equal(3))
			Expect(data.Get("beforeCommit")).Should(BeNil())
			Expect(data.Get("commit")).Should(BeNil())
			Expect(data.Get("finish")).Should(BeNil())
		})

		It("before beforeCommit aborted", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithBeforeCommit(beforeCommitAborted),
				workflow.WithInterface(&s),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(3))
			Expect(data.MustGet("rollback")).To(Equal(4))
			Expect(data.Get("commit")).Should(BeNil())
			Expect(data.Get("finish")).Should(BeNil())
		})

		It("after beforeCommit aborted", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithInterface(&s),
				workflow.WithBeforeCommit(beforeCommitAborted),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(4))
			Expect(data.MustGet("rollback")).To(Equal(5))
			Expect(data.Get("commit")).Should(BeNil())
			Expect(data.Get("finish")).Should(BeNil())
		})

		It("before commit aborted", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithCommit(commitAborted),
				workflow.WithInterface(&s),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(3))
			Expect(data.MustGet("commit")).To(Equal(4))
			Expect(data.MustGet("rollback")).To(Equal(5))
			Expect(data.Get("finish")).Should(BeNil())
		})

		It("after commit aborted", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithInterface(&s),
				workflow.WithCommit(commitAborted),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(3))
			Expect(data.MustGet("commit")).To(Equal(5))
			Expect(data.MustGet("rollback")).To(Equal(6))
			Expect(data.Get("finish")).Should(BeNil())
		})

		It("before rollback aborted", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithCommit(commitAborted),
				workflow.WithRollback(rollbackAborted),
				workflow.WithInterface(&s),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(3))
			Expect(data.MustGet("commit")).To(Equal(4))
			Expect(data.MustGet("rollback")).To(Equal(5))
			Expect(data.Get("finish")).Should(BeNil())
		})

		It("after rollback aborted", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithCommit(commitAborted),
				workflow.WithInterface(&s),
				workflow.WithRollback(rollbackAborted),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(3))
			Expect(data.MustGet("commit")).To(Equal(4))
			Expect(data.MustGet("rollback")).To(Equal(6))
			Expect(data.Get("finish")).Should(BeNil())
		})

		It("before finish aborted", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithFinish(finishAborted),
				workflow.WithInterface(&s),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(3))
			Expect(data.MustGet("commit")).To(Equal(4))
			Expect(data.MustGet("finish")).To(Equal(5))
			Expect(data.Get("rollback")).Should(BeNil())
		})

		It("after finish aborted", func() {
			data, err := workflow.StartWorkFlow(
				successFunc("work", &index),
				workflow.WithInterface(&s),
				workflow.WithFinish(finishAborted),
			)
			Expect(err).Should(BeNil())
			Expect(data.MustGet("begin")).To(Equal(1))
			Expect(data.MustGet("work")).To(Equal(2))
			Expect(data.MustGet("beforeCommit")).To(Equal(3))
			Expect(data.MustGet("commit")).To(Equal(4))
			Expect(data.MustGet("finish")).To(Equal(6))
			Expect(data.Get("rollback")).Should(BeNil())
		})
	})
})

func successFunc(key string, index *int) workflow.Event {
	return func(ctx context.Context, data *workflow.WorkData) error {
		data.Set(key, *index)
		(*index) = *index + 1
		return nil
	}
}

func failedFunc(key string, index *int) workflow.Event {
	return func(ctx context.Context, data *workflow.WorkData) error {
		data.Set(key, *index)
		(*index) = *index + 1
		return fmt.Errorf(key)
	}
}

func abortedFunc(key string, index *int) workflow.Event {
	return func(ctx context.Context, data *workflow.WorkData) error {
		data.Set(key, *index)
		(*index) = *index + 1
		data.Abort()
		return nil
	}
}
