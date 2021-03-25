package worker

import (
	"context"
	"testing"
	"time"

	"github.com/SQLek/wp-interview/model"
	"github.com/stretchr/testify/assert"
)

func TestScheluding(t *testing.T) {
	assert := assert.New(t)

	mock := new(MockModel)
	sche := MakeSheduler(Config{}, mock)

	// we want simple task to check out cheluder
	task := model.Task{
		URL:      testUrlPrefix + "/base64/TWVvdw==",
		Interval: 10,
	}

	entry := model.Entry{
		Duration: 1,
		Content:  "Meow",
		Code:     200,
		TaskID:   1,
	}

	// checking if spawning worker will add task to database
	mock.On("PutTask", task).Return(1, nil)
	mock.On("PutEntry", entry).Return(1, nil)
	mock.On("PutEntry", entry).Return(2, nil)
	mock.On("PutEntry", entry).Return(3, nil)
	mock.On("DeleteTask", uint(1)).Return(nil)
	id, err := sche.SpawnWorker(context.Background(), task)
	assert.Nil(err)

	// now we sleep something more than two iterations
	time.Sleep(12 * time.Second)
	err = sche.KillWorker(context.Background(), uint(id))
	assert.Nil(err)
	mock.AssertExpectations(t)

}
