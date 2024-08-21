package main

import (
	"context"
	"fmt"
	"time"

	"github.com/reugn/go-quartz/quartz"
)


// the job interface is not implemented in the quartz library so we'll implement it here
type Job struct{}

func (job *Job) Execute(ctx context.Context) error {
	fmt.Println("Hello, World!")
	return nil
}

func (job *Job) Description() string {
	return "This is a print job that prints 'Hello, World!'"
}

func main() {
	// Create a new job
	job := &Job{}
	jobKey := quartz.NewJobKey("print_job")

	// Create a new cron trigger 	
	cronTrigger, err := quartz.NewCronTrigger("*/3 * * * * *") // Every 3 seconds
	if err != nil {
		fmt.Println("Failed to create CronTrigger:", err)
		return
	}

	// Create a scheduler
	scheduler := quartz.NewStdScheduler()
	
	// Schedule the job using trigger
	jobDetail := quartz.NewJobDetail(job, jobKey)
	err = scheduler.ScheduleJob(jobDetail, cronTrigger)
	if err != nil {
		fmt.Println("Failed to schedule job:", err)
		return
	}

	// Start the scheduler
	ctx := context.Background()
	scheduler.Start(ctx)

	// Keep the program running for demonstration
	time.Sleep(30 * time.Second)

	// Stop the scheduler gracefully
	scheduler.Stop()
}