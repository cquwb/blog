package base

//"fmt"
//"time"

type Queue struct {
	i, j, l int
	q       []interface{}
}

func (this *Queue) Push(d interface{}) {
	n := len(this.q)
	if this.l >= n {
		this.Expand()
		n = len(this.q)
	}

	this.q[this.i] = d
	this.i++
	this.i = this.i % n
	this.l++
	//fmt.Printf("this push %d \n", i, j, l)
}

func (this *Queue) Pop() interface{} {
	if this.l <= 0 {
		return nil
	}
	n := len(this.q)
	d := this.q[this.j]
	this.j++
	this.j = this.j % n
	this.l--
	return d
}

func (this *Queue) Len() int {
	return len(this.q)
}

func (this *Queue) Expand() {
	n := len(this.q)
	var newq []interface{}
	if n <= 0 {
		newq = make([]interface{}, 8)
	} else if n <= 8 {
		newq = make([]interface{}, 16)
	} else {
		newq = make([]interface{}, n+n/4)
	}
	j := this.j
	if n > 0 {
		copy(newq[0:n-j], this.q[j:n])
		copy(newq[n-j:n], this.q[0:j])
		this.j = 0
		this.i = n
	}

	for k := range this.q {
		this.q[k] = nil
	}
	this.q = newq
}
