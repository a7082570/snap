package scheduler

import (
	"testing"

	"github.com/intelsdi-x/pulse/core"

	log "github.com/Sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	sum int
)

type dummyCatcher struct {
	count int
}

func (d *dummyCatcher) CatchCollection(m []core.Metric) {
	d.count++
	sum++
}

func TestTaskWatching(t *testing.T) {
	log.SetLevel(log.FatalLevel)
	Convey("", t, func() {
		twc := newTaskWatcherCollection()
		So(twc, ShouldHaveSameTypeAs, new(taskWatcherCollection))
		d1 := &dummyCatcher{}
		d2 := &dummyCatcher{}
		d3 := &dummyCatcher{}

		twc.add(1, d1)
		x, _ := twc.add(1, d2)
		y, _ := twc.add(2, d2)
		twc.add(3, d3)

		So(len(twc.coll[1]), ShouldEqual, 2)
		So(len(twc.coll[2]), ShouldEqual, 1)
		So(len(twc.coll), ShouldEqual, 3)

		twc.handleMetricCollected(1, nil)
		twc.handleMetricCollected(1, nil)
		twc.handleMetricCollected(2, nil)
		twc.handleMetricCollected(1, nil)
		twc.handleMetricCollected(2, nil)

		So(d1.count, ShouldEqual, 3)
		So(d2.count, ShouldEqual, 5)
		So(d3.count, ShouldEqual, 0)
		So(sum, ShouldEqual, 8)

		x.Close()
		y.Close()

		So(len(twc.coll[1]), ShouldEqual, 1)
		So(len(twc.coll[2]), ShouldEqual, 0)
		So(len(twc.coll), ShouldEqual, 2)

		twc.handleMetricCollected(1, nil)
		twc.handleMetricCollected(1, nil)
		twc.handleMetricCollected(2, nil)
		twc.handleMetricCollected(1, nil)
		twc.handleMetricCollected(2, nil)

		So(d1.count, ShouldEqual, 6)
		So(d2.count, ShouldEqual, 5)
		So(d3.count, ShouldEqual, 0)
		So(sum, ShouldEqual, 11)
	})
}
