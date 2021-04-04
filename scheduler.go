package seen

import (
	"sort"
	"time"
)

type Task func()

type RecursiveTask func(self func())

type FutureRecursiveTask func(self func(time.Duration))

type TaskRunner interface {
	Cancel()
}

type TaskScheduler interface {
	Now() time.Time
	Schedule(task Task) TaskRunner
	ScheduleFuture(task Task, due time.Duration) TaskRunner
	ScheduleRecursive(task RecursiveTask) TaskRunner
	ScheduleFutureRecursive(task FutureRecursiveTask, due time.Duration) TaskRunner
	Run() bool
	Len() int
}

var Scheduler TaskScheduler = &scheduler{}

type futuretask struct {
	at     time.Time
	run    func()
	cancel chan struct{}
}

func (t *futuretask) Cancel() {
	if t.cancel != nil {
		close(t.cancel)
	}
}

type scheduler struct {
	tasks   []futuretask
	current *futuretask
}

func (s *scheduler) Len() int {
	return len(s.tasks)
}

func (s *scheduler) Less(i, j int) bool {
	return s.tasks[i].at.Before(s.tasks[j].at)
}

func (s *scheduler) Swap(i, j int) {
	s.tasks[i], s.tasks[j] = s.tasks[j], s.tasks[i]
}

func (s *scheduler) Now() time.Time {
	return time.Now()
}

func (s *scheduler) Schedule(task Task) TaskRunner {
	t := futuretask{time.Now(), task, make(chan struct{})}
	s.tasks = append(s.tasks, t)
	sort.Stable(s)
	return &t
}

func (s *scheduler) ScheduleFuture(task Task, due time.Duration) TaskRunner {
	t := futuretask{time.Now().Add(due), task, make(chan struct{})}
	s.tasks = append(s.tasks, t)
	sort.Stable(s)
	return &t
}

func (s *scheduler) ScheduleRecursive(task RecursiveTask) TaskRunner {
	t := futuretask{cancel: make(chan struct{})}
	self := func() {
		t.at = time.Now()
		s.tasks = append(s.tasks, t)
		sort.Stable(s)
	}
	t.run = func() {
		task(self)
	}
	self()
	return &t
}

func (s *scheduler) ScheduleFutureRecursive(task FutureRecursiveTask, due time.Duration) TaskRunner {
	t := futuretask{cancel: make(chan struct{})}
	self := func(due time.Duration) {
		t.at = time.Now().Add(due)
		s.tasks = append(s.tasks, t)
		sort.Stable(s)
	}
	t.run = func() {
		task(self)
	}
	self(due)
	return &t
}

func (s *scheduler) Run() bool {
	for len(s.tasks) != 0 {
		s.current = &s.tasks[0]
		if time.Until(s.current.at) < time.Millisecond {
			s.tasks = s.tasks[1:]
			select {
			case <-s.current.cancel:
				// skip
			default:
				s.current.run()
			}
		} else {
			s.current = nil
			return true
		}
	}
	s.current = nil
	return false
}
