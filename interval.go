package seen

import "time"

type IntervalID int

var intervals = struct {
	giid    IntervalID
	runners map[IntervalID]TaskRunner
}{giid: 0, runners: make(map[IntervalID]TaskRunner)}

func SetInterval(f func(t, dt time.Duration) bool, due time.Duration) IntervalID {
	runner := Scheduler.ScheduleFutureRecursive(func(self func(time.Duration)) {
		if f(time.Duration(Scheduler.Now().UnixNano()), due) {
			self(due)
		}
	}, due)
	if runner != nil {
		intervals.giid++
		intervals.runners[intervals.giid] = runner
		return intervals.giid
	}
	return 0
}

func ClearInterval(id IntervalID) {
	if runner, present := intervals.runners[id]; present {
		runner.Cancel()
		delete(intervals.runners, id)
	}
}
