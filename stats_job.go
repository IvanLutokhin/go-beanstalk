package beanstalk

type StatsJob struct {
	ID       int    `json:"id" yaml:"id"`              // is the job id
	Tube     string `json:"tube" yaml:"tube"`          // is the name of the tube that contains this job
	State    string `json:"state" yaml:"state"`        // is "ready" or "delayed" or "reserved" or "buried"
	Priority int    `json:"priority" yaml:"priority"`  // is the priority value set by the put, release, or bury commands
	Age      int    `json:"age" yaml:"age"`            // is the time in seconds since the put command that created this job
	Delay    int    `json:"delay" yaml:"delay"`        // is the integer number of seconds to wait before putting this job in the ready queue
	TTR      int    `json:"ttr" yaml:"ttr"`            // time to run -- is the integer number of seconds a worker is allowed to run this job
	TimeLeft int    `json:"timeLeft" yaml:"time-left"` // is the number of seconds left until the server puts this job into the ready queue
	File     int    `json:"file" yaml:"file"`          // is the number of the earliest binlog file containing this job
	Reserves int    `json:"reserves" yaml:"reserves"`  // is the number of times this job has been reserved
	Timeouts int    `json:"timeouts" yaml:"timeouts"`  // is the number of times this job has timed out during a reservation
	Releases int    `json:"releases" yaml:"releases"`  // is the number of times a client has released this job from a reservation
	Buries   int    `json:"buries" yaml:"buries"`      // is the number of times this job has been buried
	Kicks    int    `json:"kicks" yaml:"kicks"`        // is the number of times this job has been kicked
}
