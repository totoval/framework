package dashboard

import (
	"sync/atomic"
)

type flow struct {
	data     int64
	chanData chan *FlowData
}
type FlowData struct {
	//Timestamp zone.Time `json:"timestamp"`
	Flows int64 `json:"flows"`
}

var Flow *flow

const flowChanBufSize = 1

func init() {
	Flow = newFlow()
}
func newFlow() *flow {
	return &flow{
		data:     0,
		chanData: make(chan *FlowData, flowChanBufSize),
	}
}
func (f *flow) Add() {
	atomic.AddInt64(&f.data, 1)
	f.notify()
}
func (f *flow) Dec() {
	if atomic.LoadInt64(&f.data) > 0 {
		atomic.AddInt64(&f.data, -1)
	}
	f.notify()
}
func (f *flow) notify() {
	if len(f.chanData) >= flowChanBufSize {
		<-f.chanData
	}
	f.chanData <- &FlowData{
		//Timestamp: zone.Now(),
		Flows: atomic.LoadInt64(&f.data),
	}
}
func (f *flow) Current() chan *FlowData {
	return f.chanData
}
func (f *flow) Close() {
	close(f.chanData)
}
