JOB=>
type Job interface {
	Execute(context.Context) error
	Description() string
}

CronJobs=>
func NewCronTrigger(expression string) (*CronTrigger, error)
[resource](https://crontab.cronhub.io/)
<second> <minute> <hour> <day-of-month> <month> <day-of-week> <year>

JobKey=>
func NewJobKey(name string) *JobKey

Scheduler=>
type Scheduler interface {
	Start(context.Context)
	IsStarted() bool

	// ScheduleJob schedules a job using a specified trigger.
	ScheduleJob(jobDetail *JobDetail, trigger Trigger) error

	GetJobKeys(...Matcher[ScheduledJob]) ([]*JobKey, error)
	GetScheduledJob(jobKey *JobKey) (ScheduledJob, error)
	DeleteJob(jobKey *JobKey) error
	PauseJob(jobKey *JobKey) error
	ResumeJob(jobKey *JobKey) error

	// Clear removes all of the scheduled jobs.
	Clear() error
	Wait(context.Context)
	Stop()
}

ScheduledJob=>
type ScheduledJob interface {
	JobDetail() *JobDetail
	Trigger() Trigger
	NextRunTime() int64
}

NewStdScheduler=>
func NewStdScheduler() Scheduler
func (sched *StdScheduler) DeleteJob(jobKey *JobKey) error

ScheduleJob=>
func (sched *StdScheduler) ScheduleJob(
	jobDetail *JobDetail,
	trigger Trigger,
) error

Trigger=>
type Trigger interface {
	NextFireTime(prev int64) (int64, error)
	Description() string
}


SimpleTrigger=> 
type SimpleTrigger struct {
	Interval time.Duration
}

Crontrigger=>
has a cron expression inside it