package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/reugn/go-quartz/quartz"
)

// Job implements the quartz.Job interface
type Job struct{}

func (job *Job) Execute(ctx context.Context) error {
	fmt.Println("Hello, World!")
	return nil
}

func (job *Job) Description() string {
	return "This is a print job that prints 'Hello, World!'"
}

// RetryScheduler restarts the scheduler with retry logic
func RetryScheduler(scheduler quartz.Scheduler, maxRetries int, jobDetail *quartz.JobDetail, trigger quartz.Trigger, failedJobs *[]*quartz.JobDetail) {
	ctx := context.Background()
	for retries := 0; retries < maxRetries; retries++ {
		scheduler.Start(ctx)
		err := scheduler.ScheduleJob(jobDetail, trigger)
		if err != nil {
			fmt.Printf("Failed to schedule job on retry %d: %v\n", retries+1, err)
		} else {
			fmt.Println("Scheduler started successfully")
			return
		}
	}
	// max retries reached
	fmt.Println("Failed to start scheduler after max retries. Adding job to failed bucket.")
	Failed(jobDetail, failedJobs)
}

// Failed adds the job to the failed bucket
func Failed(jobDetail *quartz.JobDetail, failedJobs *[]*quartz.JobDetail) {
	jobKey := jobDetail.JobKey()
	fmt.Printf("Job %v added to failed bucket\n", jobKey)
	*failedJobs = append(*failedJobs, jobDetail)
}

func main() {
	var failedJobs []*quartz.JobDetail

	// Create a new job
	job := &Job{}
	jobKey := quartz.NewJobKey("print_job")
	jobDetail := quartz.NewJobDetail(job, jobKey)

	// Create a new cron trigger for every 3 seconds
	cronTrigger, err := quartz.NewCronTrigger("*/3 * * * * *")
	if err != nil {
		fmt.Println("Failed to create CronTrigger:", err)
		return
	}

	// Create a scheduler
	scheduler := quartz.NewStdScheduler()
	RetryScheduler(scheduler, 3, jobDetail, cronTrigger, &failedJobs)

	// Channel to listen for interrupt that'll be blocked until ctrl+c is recieved
	// buffered channel with size 1 
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	fmt.Printf("Received signal %s, stopping scheduler...\n", sig)
}
