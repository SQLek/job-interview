package worker

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/SQLek/wp-interview/model"
)

// Perform a GET request on url with context and pack result in model.Entry.
func request(ctx context.Context, url string) (e model.Entry, err error) {
	start := time.Now()
	defer func() {
		e.Duration = float32(time.Since(start).Seconds())
	}()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	e.Code = resp.StatusCode
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	e.Content = string(body)

	return
}

// Fetch url and put it into model
func fetchAndPut(ctx context.Context, model model.Model, task model.Task) {
	entry, err := request(ctx, task.URL)
	if err != nil {
		log.Printf("Request failed! '%s' %s", task.URL, err)
	}
	entry.TaskID = task.ID

	_, err = model.PutEntry(entry)
	if err != nil {
		log.Printf("Database upload failed! %s", err)
		log.Println("Task: ", task)
		log.Println("Entry: ", entry)
	}

}

// Entrypoint of goroutine managing one cyclic fetch task
//
// This function will perform fetch imedietly,
// and then every task.Interval seconds
//
// Each fetch task will have timeout as specified in argument
//
// This function will not exit until cancelation of context.
// Encountered errors in fetching will not stop its execution.
func WorkOnRequests(ctx context.Context, model model.Model,
	task model.Task, timeout time.Duration) {

	// Opinion: Making struct to pass that much arguments,
	// would be good idea in most cases, but this function is
	// called only from one place (look SpawnWorker)

	ticker := time.NewTicker(time.Second * time.Duration(task.Interval))
	defer ticker.Stop()

	for {
		// cancelation propagates so no need to have cancel() here
		fetchCtx, cancelFunc := context.WithTimeout(ctx, timeout)
		go fetchAndPut(fetchCtx, model, task)

		select {

		case <-ctx.Done():
			cancelFunc()
			return
		case <-ticker.C:
			// nothing, just iterate
		}

	}

}
