package queue

import (
	"encoding/json"
)

type Packet[T any] struct {
	Key     string `json:"key"`
	Content T      `json:"content"`
}

func (p *Packet[T]) GetContent() T {
	return p.Content
}

func (p *Packet[T]) FromJson(data []byte) Packet[T] {
	_ = json.Unmarshal(data, p)
	return *p
}

func (p *Packet[T]) ToByte() []byte {
	rs, err := json.Marshal(p)
	if err != nil {
		return nil
	}
	return rs
}

func (p *Packet[T]) ToString() string {
	rs, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(rs)
}

func (p *Packet[T]) ToMap() map[string]interface{} {
	mp := map[string]interface{}{}
	_ = json.Unmarshal(p.ToByte(), &mp)
	return mp
}
