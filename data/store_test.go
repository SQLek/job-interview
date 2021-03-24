package data

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFetchTask(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(testModel)

	// This test will perform typical use case

	// Creating task
	task := Task{
		Interval: 60,
		URL:      "https://httpbin.org/range/10",
	}
	tid, err := testModel.PutTask(task)
	assert.Nil(err)
	assert.Greater(tid, uint(0), "ID should not be zero")

	// TODO: Maybe add Update task here

	// Put two example entries
	entry1 := Entry{
		TaskID:   tid,
		Duration: 0.7,
		Content:  "qwertyuiop",
		Code:     200,
	}
	eid1, err := testModel.PutEntry(entry1)
	assert.Nil(err)
	assert.Greater(eid1, uint(0), "ID should not be zero")

	// inserting two entries at once meses one test bellow
	time.Sleep(5 * time.Second)

	entry2 := Entry{
		TaskID:   tid,
		Duration: 0.7,
		Content:  "poiuytrewq",
		Code:     200,
	}
	eid2, err := testModel.PutEntry(entry2)
	assert.Nil(err)
	assert.Greater(eid2, eid1, "ID should greater than previeus")

	// Lets try list our entries
	entries, err := testModel.ListEntries(tid)
	assert.Nil(err)
	assert.Equal(2, len(entries), "We expected to got two entries")
	//assert.Greater(entries[1].CreatedAt, entries[0].CreatedAt, "Entries are not sorted!")
	// !HACK! problems with gorm.CreatedAt. Look model.go at Model definition
	assert.Greater(entries[1].Created, entries[0].Created, "Entries are not sorted!")

	// Now lets list entries and make sure they are no preloaded
	tasks, err := testModel.ListTasks()
	assert.Nil(err)
	assert.Greater(len(tasks), 0, "Task list should not be empty")

	// TODO add iteration to enable concurrent testing
	// do not extend api just for testing
	ourTask := tasks[0]
	assert.Equal(tid, ourTask.ID, "Did this test run concurently?")
	assert.Empty(ourTask.Entries, "Entries should not be populated.")

	err = testModel.DeleteTask(tid)
	assert.Nil(err)
	// TODO maybe double chceck if it is deleted?

	// Lets chceck if entries are deleted
	entries2, err := testModel.ListEntries(tid)
	assert.Nil(err)
	assert.Empty(entries2)

}
