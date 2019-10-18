package base

import (
	//"fmt"
	"time"
)

/*
时间轮算法定时器
*/

const MaxWheel int = 120

type TimeWheel struct {
	id       uint32
	wheel    []*Wheel
	i        int
	lasttime int64
	si       int //每个轮子的间隔的倒数，以秒为单位，最小是一秒
	mapWheel map[uint32]*WheelNode
}

func NewTimeWheel(si int) *TimeWheel {
	r := &TimeWheel{
		id:       0,
		wheel:    make([]*Wheel, MaxWheel),
		lasttime: time.Now().Unix(),
		si:       si,
		mapWheel: make(map[uint32]*WheelNode),
	}
	for i := 0; i < MaxWheel; i++ {
		r.wheel[i] = NewWheel(i)
	}
	return r
}

//加入一个超时时间， t是回调函数
func (this *TimeWheel) Add(e TimeEntry, end int64, intval int) uint32 {
	if end < this.lasttime {
		return 0
	}
	this.id++
	index, round := this.GetIndexAndRound(end)
	wheel := this.wheel[index]
	node := &WheelNode{
		id:       this.id,
		e:        e,
		preNode:  nil,
		nextNode: nil,
		round:    round,
		index:    index,
		end:      end,
		intval:   intval,
		check:    0,
	}
	wheel.Add(node)
	//this.mapWheel[this.id] = node
	return this.id
}

func (this *TimeWheel) GetIndexAndRound(end int64) (int, int) {
	index := (int(end-this.lasttime)*this.si + this.i) % MaxWheel
	round := (int(end-this.lasttime-1) * this.si) / MaxWheel
	//fmt.Printf("GetIndexAndRound %d %d %d %d \n", end, this.lasttime, index, round)
	return index, round
}

//移除一个事件
func (this *TimeWheel) Remove(id uint32) TimeEntry {
	//fmt.Printf("begin remove id  %d \n", id)
	if n, ok := this.mapWheel[id]; ok {
		this.wheel[n.index].Remove(n)
		delete(this.mapWheel, n.id)
		return n.e
	}
	return nil
}

//定时驱动的
func (this *TimeWheel) Run(now int64) {
	this.i++
	this.i = this.i % MaxWheel
	this.lasttime = now
	//fmt.Printf("this.i is %d this.lasttime=%d \n", this.i, this.lasttime)
	wheel := this.wheel[this.i]
	node := wheel.hNode
	for node != nil {
		if node.check >= now {
			break
		}
		nextNode := node.nextNode
		if node.round <= 0 {
			node.e.Timeout(node.id, now)
			wheel.Remove(node)
			if node.intval > 0 {
				node.end += int64(node.intval)
				index, round := this.GetIndexAndRound(node.end)
				node.index = index
				node.round = round
				newWheel := this.wheel[index]
				node.check = now
				newWheel.Add(node)
			} else {
				delete(this.mapWheel, node.id)
			}
		} else {
			node.round--
			node.check = now
		}

		if nextNode == node {
			break //是头元素
		}
		node = nextNode
	}

}

//时间轮的每个轮子部分 双向列表
type Wheel struct {
	hNode *WheelNode
	index int
}

func NewWheel(index int) *Wheel {
	return &Wheel{
		hNode: nil,
		index: index,
	}
}

func (this *Wheel) Add(node *WheelNode) {
	if this.hNode == nil {
		this.hNode = node
		node.nextNode = node
		node.preNode = node
	} else {
		preNode := this.hNode.preNode
		node.nextNode = this.hNode
		node.preNode = preNode
		preNode.nextNode = node
		this.hNode.preNode = node
	}
}

func (this *Wheel) Remove(node *WheelNode) {
	preNode := node.preNode
	nextNode := node.nextNode
	if preNode == node && nextNode == node {
		this.hNode = nil
	} else {
		if node == this.hNode {
			this.hNode = nextNode //头元素替换
		}
		preNode.nextNode = nextNode
		nextNode.preNode = preNode
		node.preNode = nil
		node.nextNode = nil
	}
}

//每个节点
type WheelNode struct {
	id       uint32
	e        TimeEntry
	preNode  *WheelNode
	nextNode *WheelNode
	round    int //代表是第轮执行的
	index    int
	end      int64
	intval   int
	check    int64
}

type TimeEntry interface {
	Timeout(uint32, int64)
}
