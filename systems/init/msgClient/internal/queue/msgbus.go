package queue

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/init/msgClient/internal/db"
	"google.golang.org/protobuf/types/known/anypb"
)

type MsgBusHandlerInterface interface {
	CreateServiceMsgBusHandler() error
	StopServiceQueueHandler(service string) (err error)
	UpdateServiceQueueHandler(s *db.Service) (err error)
	Publish(service string, key string, msg *anypb.Any) error
}

type MsgBusHandler struct {
	ql  map[string]*QueueListener
	qp  map[string]*QueuePublisher
	s   db.ServiceRepo
	r   db.RouteRepo
	mia uint32
	pHC time.Duration
	mR  chan bool
}

func NewMessageBusHandler(s db.ServiceRepo, r db.RouteRepo, miss uint32, period time.Duration) *MsgBusHandler {

	h := &MsgBusHandler{
		s:   s,
		r:   r,
		mia: miss,
		pHC: period,
	}
	h.mR = make(chan bool, 1)
	h.ql = make(map[string]*QueueListener)
	h.qp = make(map[string]*QueuePublisher)
	return h

}

func (m *MsgBusHandler) CreateServiceMsgBusHandler() error {

	services, err := m.s.List()
	if err != nil {
		log.Errorf("Error reading services. Error %s", err.Error())
		return err
	}

	/* Start routine to monitor */
	m.monitor(m.mR)

	if len(services) <= 0 {
		log.Errorf("No services available.")
	}

	for _, s := range services {

		log.Infof("Creating message bus handler  for %s service.", s.ServiceUuid)

		/* Create publisher */
		publisher, err := NewQueuePublisher(s)
		if err != nil {
			log.Errorf("Failed to create Publisher for %s. Error %s", s.Name, err.Error())
		} else {
			m.qp[s.ServiceUuid] = publisher
		}

		/*  Create a queue listner for each service */
		listener, err := NewQueueListener(s)
		if err != nil {
			log.Errorf("Failed to create listneer for %s. Error %s", s.Name, err.Error())
		} else {
			m.ql[s.ServiceUuid] = listener
		}

		log.Debugf("Service: %s \n Listner: %+v  \n Publisher: %+v", s.Name, listener, publisher)

	}

	m.StartQueueListeners()

	return nil
}

func (m *MsgBusHandler) StartQueueListeners() {

	for _, q := range m.ql {
		/*  Create a queue listner for each service */
		log.Infof("Starting new queue listener routine for service %s on %v routes", q.serviceName, q.routes)
		go q.startQueueListening()

	}

}

func (m *MsgBusHandler) StopQueueListener() {

	for _, q := range m.ql {
		q.stopQueueListening()
	}
}

func (m *MsgBusHandler) RestartServiceQueueListening(service string) (err error) {
	q, ok := m.ql[service]
	if ok {
		q.stopQueueListening()
		time.Sleep(500 * time.Millisecond)
		if !q.state {
			q.startQueueListening()
		}
	}
	return nil
}

func (m *MsgBusHandler) StopServiceQueueListening(service string) (err error) {
	q, ok := m.ql[service]
	if ok {
		q.stopQueueListening()
		time.Sleep(500 * time.Millisecond)
		if q.state {
			return fmt.Errorf("failed to stop queue listening service for %s", q.serviceName)
		}
	} else {
		return fmt.Errorf("no service with id %s registered", service)
	}

	return nil
}

func (m *MsgBusHandler) RemoveServiceQueueListening(service string) error {
	_, ok := m.ql[service]
	if ok {
		err := m.StopServiceQueueListening(service)
		if err != nil {
			return err
		}
		delete(m.ql, service)
		log.Infof("Removed queue listener for %s service", service)
	}

	return nil
}

func (m *MsgBusHandler) RemoveServiceQueuePublisher(service string) error {
	p, ok := m.qp[service]
	if ok {
		err := p.Close()
		if err != nil {
			return err
		}
		delete(m.qp, service)
		log.Infof("Removed queue listener for %s service", service)
	}

	return nil
}

func (m *MsgBusHandler) StopServiceQueueHandler(service string) (err error) {

	/* Stop publisher */
	p, ok := m.qp[service]
	if ok {
		err := p.Close()
		if err != nil {
			return err
		}
	} else {
		log.Errorf("No service with id %s registered in publisher", service)
	}

	/* Stop listener */
	_, ok = m.ql[service]
	if ok {
		err := m.StopServiceQueueListening(service)
		if err != nil {
			return err
		}
	} else {
		log.Errorf("No service with id %s registered in listener", service)
	}

	return err
}

/* start/Update Message queue parameters */
func (m *MsgBusHandler) UpdateServiceQueueHandler(s *db.Service) (err error) {

	log.Debugf("Removing old listener and publisher if any for service %s.", s.Name)
	/* Listner */
	m.RemoveServiceQueueListening(s.ServiceUuid)

	/* Publisher */
	m.RemoveServiceQueuePublisher(s.ServiceUuid)

	log.Debugf("Removing old listener and publisher if any for service %s completed.", s.Name)

	listener, err := NewQueueListener(*s)
	if err != nil {
		log.Errorf("Failed to create listener for %s. Error %s", s.Name, err.Error())
		return err
	} else {
		m.ql[s.ServiceUuid] = listener
	}

	go listener.startQueueListening()

	/* Check listner state before returning */
	time.Sleep(500 * time.Millisecond)

	if !listener.state {
		return fmt.Errorf("failed to start listener for service %s", listener.serviceName)
	}

	publisher, err := NewQueuePublisher(*s)
	if err != nil {
		log.Errorf("Failed to create publisher for %s. Error %s", s.Name, err.Error())
		return err
	} else {
		m.qp[s.ServiceUuid] = publisher
	}

	log.Debugf("Started listener and publisher if any for service %s.", s.Name)
	return nil
}

func (m *MsgBusHandler) Publish(service string, key string, msg *anypb.Any) error {
	p, ok := m.qp[service]
	if ok {

		err := p.Publish(key, msg)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("no publisher for service %s found", service)
	}

	return nil
}

func (m *MsgBusHandler) doHealthCheck() error {
	log.Debugf("[Health Check Monitor] Starting HealthCheck at %s", time.Now().Format(time.RFC1123))
	for id, q := range m.ql {
		if q.state {
			q.healthCheck()
			if q.continuousMiss > m.mia {
				if err := m.RemoveServiceQueueListening(id); err != nil {
					log.Errorf("[Health Check Monitor] Failed to remove listener for %s with id %s . Error %s", q.serviceName, id, err.Error())
				}

				if err := m.RemoveServiceQueuePublisher(id); err != nil {
					log.Errorf("[Health Check Monitor] Failed to remove publisher for %s with id %s. Error %s", q.serviceName, id, err.Error())
				}
			}
		}
	}
	log.Debugf("[Health Check Monitor] Completed HealthCheck at %s.", time.Now().Format(time.RFC1123))
	return nil
}

func (m *MsgBusHandler) monitor(q chan bool) {
	log.Infof("Starting health check routine with period %s.", m.pHC)
	t := time.NewTicker(m.pHC)

	go func() {
		for {
			select {
			case <-t.C:
				m.doHealthCheck()
			case <-q:
				t.Stop()
				return
			}
		}
	}()

}
