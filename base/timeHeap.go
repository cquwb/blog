package base

//"fmt"
//"time"

type Timeouter interface {
	Timeout(uint32, int64)
}

type Timer struct {
	Timeouter
	id     uint32
	end    int64
	intval int
}

type TimeHeap []*Timer

func (this *TimeHeap) Add(t *Timer) {
	n := this.Len()
	stmp := *this
	stmp = stmp[0 : n+1]
	stmp[n] = t
	*this = stmp
	this.Up(n)
}

func (this *TimeHeap) Up(i int) {
	for i > 0 {
		parent := (i - 1) / 2
		if parent == i {
			break
		}
		if this.Less(i, parent) {
			this.Swap(i, parent)
			i = parent
		} else {
			break
		}
	}
}

func (this *TimeHeap) Down(i, n int) bool {
	t := i
	for {
		left := t*2 + 1
		if left >= n || left < 0 {
			break
		}
		s := left
		if right := left + 1; right < n && this.Less(right, left) {
			s = right
		}
		if this.Less(s, t) {
			this.Swap(t, s)
			t = s
		} else {
			break
		}
	}
	return t > i
}

func (this *TimeHeap) Remove(i int) *Timer {
	n := this.Len()
	stmp := *this
	if n-1 != i {
		this.Swap(i, n-1)
	}
	t := stmp[n-1]
	stmp[n-1] = nil
	stmp = stmp[0 : n-1]
	*this = stmp
	this.Fix(i, n-1)
	return t

}

func (this *TimeHeap) Fix(i, n int) {
	if !this.Down(i, n) {
		this.Up(i)
	}
}

func (this *TimeHeap) Pop() *Timer {
	return this.Remove(0)

}

func (this TimeHeap) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this TimeHeap) Less(i, j int) bool {
	if this[i].end < this[j].end {
		return true
	}
	return false
}

func (this TimeHeap) Len() int {
	return len(this)
}

func (this TimeHeap) Cap() int {
	return cap(this)
}

func (this *TimeHeap) Expand() {
	n := this.Len()
	c := this.Cap()
	n1 := make([]*Timer, n, c*2)
	copy(n1, *this)
	*this = n1
}

type TimeManager struct {
	id uint32
	th TimeHeap
}

func (this *TimeManager) Add(t Timeouter, end int64, intval int) uint32 {
	if this.th.Cap() <= this.th.Len() {
		//panic("time add panic")
		this.th.Expand()
		//return 0
	}
	this.id++
	timer := &Timer{t, 0, end, intval}
	timer.id = this.id
	this.th.Add(timer)
	//fmt.Printf("%d Add %d end:=%d \n", time.Now().Unix(), this.id, end)
	return this.id
}

func (this *TimeManager) Remove(id uint32) {
	for k, v := range this.th {
		if v.id == id {
			this.th.Remove(k)
			break
		}
	}
}

func (this *TimeManager) Run(now int64) {
	//fmt.Printf("begin run \n")
	//smtp := *this.th
	for this.th.Len() > 0 {
		d := this.th[0]
		if d.end <= now {
			tt := this.th.Pop()
			tt.Timeout(tt.id, now)
			if d.intval > 0 {
				d.end += int64(d.intval)
				this.th.Add(d)
			}
		} else {
			break
		}
	}
}

func NewTimeManager(size int) *TimeManager {
	return &TimeManager{
		id: 0,
		th: make([]*Timer, 0, size),
	}
}
