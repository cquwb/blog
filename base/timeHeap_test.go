package base

import (
	//"fmt"
	"testing"
	"time"
)

/*
时间轮算法定时器
*/

type TestDataHeap struct {
}

func (this *TestDataHeap) Timeout(id uint32, now int64) {

}

func BenchmarkAddTimeHeap(b *testing.B) {
	manager := NewTimeManager(1024)
	now := time.Now().Unix()
	delay := []int64{1, 60, 61, 120, 2, 3, 53, 70, 10, 11, 9, 8, 30, 35, 43, 25, 40, 60, 70, 56, 45, 65, 74}
	intval := 0
	t := 0
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//b.StopTimer()
		e := &TestDataHeap{}
		t++
		t = t % len(delay)
		//b.StartTimer()
		manager.Add(e, now+delay[t], intval)

	}
}

func BenchmarkTimingWheel2(b *testing.B) {
	manager := NewTimeManager(1024)
	now := time.Now().Unix()

	cases := []struct {
		name string
		N    int64 // the data size (i.e. number of existing timers)
	}{
		{"N-0m", int64(0)},
		{"N-1m", int64(1000000)},
		{"N-10m", int64(10000000)},
	}
	for _, c := range cases {
		b.Run(c.name, func(b *testing.B) {
			//base := make([]int, c.N)
			for i := int64(0); i < c.N; i++ {
				manager.Add(&TestDataHeap{}, now+i, 0)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				manager.Add(&TestDataHeap{}, now+60, 0)
			}
		})
	}
}
