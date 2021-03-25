package model

import "time"

// Representation of URL fetching task
type Task struct {
	ID       uint    `json:"id"`
	Interval uint    `json:"interval"`
	URL      string  `json:"url"`
	Entries  []Entry `json:"-"`
}

// Representation of single fetch attempt
type Entry struct {
	ID     uint `json:"-"`
	TaskID uint `json:"-"`
	//CreatedAt time.Time - for some reason it works in gorm.Model?
	// TODO investigate and remove following hack
	Created float64 `json:"created_at"`

	Duration float32 `json:"duration"`
	Content  string  `json:"response"`
	Code     int     `json:"-"`
}

// Puts task into database and returns its ID.
// It will not spawn cyclic worker!
// Look into worker module for more info.
func (d *data) PutTask(t Task) (ID uint, err error) {

	result := d.db.Create(&t)
	return t.ID, result.Error

}

// Obtain slice of tasks.
// This method is will not obtain coresponding entries.
// Use PopulateEntries to populate Task if needed
func (d *data) ListTasks() ([]Task, error) {

	tasks := make([]Task, 0)
	result := d.db.Find(&tasks)
	return tasks, result.Error

}

// Updates task in database.
// It will not affect cyclic worker!
// Look into worker module for more info.
func (d *data) UpdateTask(t Task) error {

	result := d.db.Save(&t)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNoRowsAfected
	}
	return nil

}

// Removes task with all its entries.
// It will not kill cyclic worker!
// Look into worker module for more info.
// This method is hard delete.
func (d *data) DeleteTask(id uint) error {

	result := d.db.Where("task_id = ?", id).Delete(&Entry{})
	if result.Error != nil {
		return result.Error
	}

	result2 := d.db.Delete(&Task{}, id)
	if result2.Error != nil {
		return result2.Error
	}

	if result2.RowsAffected == 0 {
		return ErrNoRowsAfected
	}
	return nil
}

// Put single entry into database.
func (d *data) PutEntry(e Entry) (uint, error) {

	// !HACK! createdAt don't work currencly
	// More info look Entry

	e.Created = float64(time.Now().UnixNano()) / 1000000

	result := d.db.Create(&e)
	return e.ID, result.Error

}

// Obtain slice of entries.
func (d *data) ListEntries(id uint) ([]Entry, error) {

	entries := make([]Entry, 0)
	result := d.db.Order("created").Find(&entries)
	return entries, result.Error

}
