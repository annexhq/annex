package event

import (
	testv1 "github.com/annexsh/annex-proto/gen/go/type/test/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var typeProto = map[Type]testv1.Event_Type{
	TypeUnspecified:            testv1.Event_TYPE_UNSPECIFIED,
	TypeTestExecutionScheduled: testv1.Event_TYPE_TEST_EXECUTION_SCHEDULED,
	TypeTestExecutionStarted:   testv1.Event_TYPE_TEST_EXECUTION_STARTED,
	TypeTestExecutionFinished:  testv1.Event_TYPE_TEST_EXECUTION_FINISHED,
	TypeCaseExecutionScheduled: testv1.Event_TYPE_CASE_EXECUTION_SCHEDULED,
	TypeCaseExecutionStarted:   testv1.Event_TYPE_CASE_EXECUTION_STARTED,
	TypeCaseExecutionFinished:  testv1.Event_TYPE_CASE_EXECUTION_FINISHED,
	TypeLogPublished:           testv1.Event_TYPE_LOG_PUBLISHED,
}

func (t Type) Proto() testv1.Event_Type {
	pb, ok := typeProto[t]
	if !ok {
		return testv1.Event_TYPE_UNSPECIFIED
	}
	return pb
}

var dataTypeProto = map[DataType]testv1.Event_Data_Type{
	DataTypeUnspecified:   testv1.Event_Data_TYPE_UNSPECIFIED,
	DataTypeNone:          testv1.Event_Data_TYPE_NONE,
	DataTypeTestExecution: testv1.Event_Data_TYPE_TEST_EXECUTION,
	DataTypeCaseExecution: testv1.Event_Data_TYPE_CASE_EXECUTION,
	DataTypeLog:           testv1.Event_Data_TYPE_LOG,
}

func (t DataType) Proto() testv1.Event_Data_Type {
	pb, ok := dataTypeProto[t]
	if !ok {
		return testv1.Event_Data_TYPE_UNSPECIFIED
	}
	return pb
}

func (e *ExecutionEvent) Proto() *testv1.Event {
	data := &testv1.Event_Data{
		Type: e.Data.Type.Proto(),
	}

	switch e.Data.Type {
	case DataTypeUnspecified, DataTypeNone:
	case DataTypeTestExecution:
		if e.Data.TestExecution != nil {
			data.Data = &testv1.Event_Data_TestExecution{
				TestExecution: e.Data.TestExecution.Proto(),
			}
		}
	case DataTypeCaseExecution:
		if e.Data.CaseExecution != nil {
			data.Data = &testv1.Event_Data_CaseExecution{
				CaseExecution: e.Data.CaseExecution.Proto(),
			}
		}
	case DataTypeLog:
		if e.Data.Log != nil {
			data.Data = &testv1.Event_Data_Log{
				Log: e.Data.Log.Proto(),
			}
		}
	}

	return &testv1.Event{
		EventId:         uuid.NewString(),
		TestExecutionId: e.TestExecID.String(),
		Type:            e.Type.Proto(),
		Data:            data,
		CreateTime:      timestamppb.New(e.CreateTime),
	}
}
