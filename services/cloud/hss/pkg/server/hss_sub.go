package server

import pb "github.com/ukama/ukama/services/cloud/hss/pb/gen"

type HssSubscriber interface {
	ImsiAdded(org string, imsi *pb.ImsiRecord)
	ImsiUpdated(org string, imsi *pb.ImsiRecord)
	ImsiDeleted(org string, imsi string)
	GutiAdded(org string, imsi string, guti *pb.Guti)
	TaiUpdated(org string, tai *pb.UpdateTaiRequest)
}

type HssEventsSubscribers struct {
	Subscriber []HssSubscriber
}

func NewHssEventsSubscribers(subscriber ...HssSubscriber) *HssEventsSubscribers {
	return &HssEventsSubscribers{Subscriber: subscriber}
}

func (s HssEventsSubscribers) ImsiAdded(org string, imsi *pb.ImsiRecord) {
	go func() {
		for _, subscriber := range s.Subscriber {
			subscriber.ImsiAdded(org, imsi)
		}
	}()
}

func (s HssEventsSubscribers) ImsiUpdated(org string, imsi *pb.ImsiRecord) {
	go func() {
		for _, subscriber := range s.Subscriber {
			subscriber.ImsiUpdated(org, imsi)
		}
	}()
}

func (s HssEventsSubscribers) ImsiDeleted(org string, imsi string) {
	go func() {
		for _, subscriber := range s.Subscriber {
			subscriber.ImsiDeleted(org, imsi)
		}
	}()
}

func (s HssEventsSubscribers) GutiAdded(org string, imsi string, guti *pb.Guti) {
	go func() {
		for _, subscriber := range s.Subscriber {
			subscriber.GutiAdded(org, imsi, guti)
		}
	}()
}

func (s HssEventsSubscribers) TaiUpdated(org string, tai *pb.UpdateTaiRequest) {
	go func() {
		for _, subscriber := range s.Subscriber {
			subscriber.TaiUpdated(org, tai)
		}
	}()
}
