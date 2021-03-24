package worker

import (
	"github.com/SQLek/wp-interview/model"
	"github.com/stretchr/testify/mock"
)

// this is a mocked model.Model

type MockModel struct {
	mock.Mock
}

func (m *MockModel) PutTask(t model.Task) (ID uint, err error) {
	args := m.Called(t)
	return uint(args.Int(0)), args.Error(1)
}

func (m *MockModel) ListTasks() ([]model.Task, error) {
	args := m.Called()
	return args.Get(0).([]model.Task), args.Error(1)
}

func (m *MockModel) UpdateTask(t model.Task) error {
	args := m.Called(t)
	return args.Error(0)
}

func (m *MockModel) DeleteTask(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockModel) PutEntry(e model.Entry) (uint, error) {
	// we have to normalize duration to one second
	e.Duration = 1
	args := m.Called(e)
	return uint(args.Int(0)), args.Error(1)
}

func (m *MockModel) ListEntries(id uint) ([]model.Entry, error) {
	args := m.Called(id)
	return args.Get(0).([]model.Entry), args.Error(1)
}

func (m *MockModel) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockModel) Ping() error {
	args := m.Called()
	return args.Error(0)
}
