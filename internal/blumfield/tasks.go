package blumfield

import (
	"blumfield/internal/models"
	"context"
	"fmt"
	"math/rand/v2"
	"time"
)

func (b *Blumfield) GetTasks(ctx context.Context) (*[]models.TasksResponse, error) {
	tasks := []models.TasksResponse{}

	// Check if the context is already done before making the request
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Make the request and check for errors
	resp, err := b.client.R().
		SetHeaders(b.BaseHeaders).
		SetContext(ctx). // Set context to handle cancellation
		SetResult(&tasks).
		Get(earnURL)

	if err != nil {
		return nil, err
	}

	// Optional: Check the response status code for additional error handling
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode())
	}

	return &tasks, nil
}

func (b *Blumfield) CompleteTasks(ctx context.Context, tasks *[]models.TasksResponse) {
	if ctx.Err() != nil {
		return
	}
	for _, response := range *tasks {
		for _, section := range response.SubSections {
			if section.Title != "Frens" && section.Title != "Farming" && section.Title != "New" {
				b.log.Info("Section: ", section.Title)
				for _, task := range section.Tasks {
					select {
					case <-ctx.Done():
						b.log.Info("Shutting down...")
						return
					default:
					}
					if task.ValidationType != "KEYWORD" {
						switch task.Status {
						case "FINISHED":
							b.tools.LogTask(&task, "Finished")
							continue
						case "READY_FOR_CLAIM":
							resp, err := b.client.R().
								SetHeaders(b.BaseHeaders).
								SetResult(&models.Task{}).
								Post(ClaimTaskURL(task.ID))
							if err != nil {
								b.tools.LogTask(&task, "Failed")
								b.log.Error("Error claiming task: ", err)
							}
							b.tools.LogTask(resp.Result().(*models.Task), "Claimed!")
						default:
							respStart, err := b.client.R().
								SetHeaders(b.BaseHeaders).
								SetResult(&models.Task{}).
								Post(StartTaskURL(task.ID))
							if err != nil {
								b.tools.LogTask(&task, "Failed")
								b.log.Error("Error starting task: ", err)
							}
							b.tools.LogTask(respStart.Result().(*models.Task), "Started")
							time.Sleep(time.Duration(rand.IntN(6)+1) * time.Second)

							respClaim, err := b.client.R().
								SetHeaders(b.BaseHeaders).
								SetResult(&models.Task{}).
								Post(ClaimTaskURL(task.ID))
							if err != nil {
								b.tools.LogTask(&task, "Failed")
								b.log.Error("Error claiming task: ", err)
							}
							b.tools.LogTask(respClaim.Result().(*models.Task), "Claimed!")
						}
						time.Sleep(time.Duration(rand.IntN(4)+1) * time.Second)
					} else {
						continue
					}
				}
			} else {
				continue
			}
		}
	}
}
