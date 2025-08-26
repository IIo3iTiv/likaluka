package main

import "itc/proto/v1"

type Datsets struct {
	x map[string][]*proto.DataMessage
}

func (d *Datsets) Add(dm *proto.DataMessage) {
	d.x[dm.Guid] = append(d.x[dm.Guid], dm)
}

func (d *Datsets) Check() {
	for guid, v := range d.x {
		if len(v) == 25600 {
			_ = guid
		}
	}
}
