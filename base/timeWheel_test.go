package base

import (
	//"fmt"
	"testing"
	"time"
)

/*
时间轮算法定时器
*/

type TestData struct {
}

func (this *TestData) Timeout(id uint32, now int64) {

}

func TestTimeWheelAdd(t *testing.T) {
	manager := NewTimeWheel(2)
	e := &TestData{}
	now := time.Now().Unix()
	delay := []int64{1, 60, 11}
	end := now + delay[0]
	manager.Add(e, end, 0)
	wheel := manager.wheel[2]
	if wheel.hNode == nil {
		t.Fatal("add error not add")
	}
	if wheel.hNode.e != e {
		t.Errorf("add error add is wrong place")
	}
	e2 := &TestData{}
	end = now + delay[1]
	manager.Add(e2, end, 0)
	wheel = manager.wheel[0]
	if wheel.hNode == nil {
		t.Fatal("add e2 error not add")
	}
	if wheel.hNode.e != e2 {
		t.Errorf("add e2 error add is wrong place")
	}

}

func BenchmarkAddTimeWheel(b *testing.B) {
	manager := NewTimeWheel(2)
	now := time.Now().Unix()
	delay := []int64{1, 60, 61, 120, 2, 3, 53, 70, 10, 11, 9}
	intval := 0
	t := 0
	for i := 0; i < b.N; i++ {
		//b.StopTimer()
		e := &TestData{}
		t++
		t = t % len(delay)
		//b.StartTimer()
		manager.Add(e, now+delay[t], intval)

	}
}
