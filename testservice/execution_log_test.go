package testservice

import (
	"context"
	"testing"

	testservicev1 "github.com/annexhq/annex-proto/gen/go/rpc/testservice/v1"
	testv1 "github.com/annexhq/annex-proto/gen/go/type/test/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/annexhq/annex/internal/fake"
	"github.com/annexhq/annex/internal/ptr"
	"github.com/annexhq/annex/test"
)

func TestService_PublishTestExecutionLog(t *testing.T) {
	tests := []struct {
		name      string
		isCaseLog bool
	}{
		{
			name:      "publish test execution log",
			isCaseLog: false,
		},
		{
			name:      "publish case execution log",
			isCaseLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			s, fakes := newService()

			created, err := fakes.repo.CreateTest(ctx, fake.GenTestDefinition())
			require.NoError(t, err)

			te, err := fakes.repo.CreateScheduledTestExecution(ctx, fake.GenScheduledTestExec(created.ID))
			require.NoError(t, err)

			var reqCaseExecID *int32
			var wantCaseExecID *test.CaseExecutionID
			if tt.isCaseLog {
				ce, err := fakes.repo.CreateScheduledCaseExecution(ctx, fake.GenScheduledCaseExec(te.ID))
				require.NoError(t, err)
				wantCaseExecID = &ce.ID
				reqCaseExecID = ptr.Get(wantCaseExecID.Int32())
			}

			req := &testservicev1.PublishTestExecutionLogRequest{
				TestExecId: te.ID.String(),
				CaseExecId: reqCaseExecID,
				Level:      "INFO",
				Message:    "lorem ipsum",
				CreatedAt:  timestamppb.Now(),
			}
			res, err := s.PublishTestExecutionLog(ctx, req)
			require.NoError(t, err)
			assert.NotNil(t, res)

			execLogID, err := uuid.Parse(res.Id)
			require.NoError(t, err)

			got, err := fakes.repo.GetExecutionLog(ctx, execLogID)
			require.NoError(t, err)

			want := &test.ExecutionLog{
				ID:         execLogID,
				TestExecID: te.ID,
				CaseExecID: wantCaseExecID,
				Level:      req.Level,
				Message:    req.Message,
				CreateTime: req.CreatedAt.AsTime(),
			}

			assert.Equal(t, got, want)
		})
	}
}

func TestService_ListTestExecutionLogs(t *testing.T) {
	wantNumTestLogs := 15
	wantNumCaseLogs := 15

	ctx := context.Background()
	s, fakes := newService()

	tt, err := fakes.repo.CreateTest(ctx, fake.GenTestDefinition())
	require.NoError(t, err)

	te, err := fakes.repo.CreateScheduledTestExecution(ctx, fake.GenScheduledTestExec(tt.ID))
	require.NoError(t, err)

	ce, err := fakes.repo.CreateScheduledCaseExecution(ctx, fake.GenScheduledCaseExec(te.ID))
	require.NoError(t, err)

	var want []*testv1.ExecutionLog

	for range wantNumTestLogs {
		l := fake.GenTestExecLog(te.ID)
		err = fakes.repo.CreateExecutionLog(ctx, l)
		require.NoError(t, err)
		want = append(want, l.Proto())
	}

	for range wantNumCaseLogs {
		l := fake.GenCaseExecLog(te.ID, ce.ID)
		err = fakes.repo.CreateExecutionLog(ctx, l)
		require.NoError(t, err)
		want = append(want, l.Proto())
	}

	res, err := s.ListTestExecutionLogs(ctx, &testservicev1.ListTestExecutionLogsRequest{
		TestExecId: te.ID.String(),
	})
	require.NoError(t, err)

	got := res.Logs
	assert.Len(t, got, wantNumTestLogs+wantNumCaseLogs)
	assert.Equal(t, want, got)
}