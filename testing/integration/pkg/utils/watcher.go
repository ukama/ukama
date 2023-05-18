package utils

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/testing/integration/pkg/messaging"
)

type EventValidator struct {
	key       string
	Validator func(string, []byte) bool
}

type Watcher struct {
	v []EventValidator
	l messaging.Listener
}

func NewWatcher(v []EventValidator) *Watcher {
	c := messaging.NewListenerConfig()
	return &Watcher{
		v: v,
		l: messaging.NewListener(c),
	}
}

func DummyValidator(e string, b []byte) bool {
	log.Debugf("Event: %s \n Body: %s", e, b)
	return true
}

func (w *Watcher) Start() {
	go w.l.StartListener()
}

func (w *Watcher) Stop() {
	w.l.StopListener()
}

func (w *Watcher) Expections() bool {
	time.Sleep(1 * time.Second)
	for _, e := range w.v {
		/* For now jsut checking event name  */
		i, ok := w.l.GetEvent(e.key)
		if !ok {
			log.Errorf("Event for %s is missing", e.key)
			return false
		}

		if e.Validator != nil {
			b, ok := i.([]byte)
			if ok {
				if !e.Validator(e.key, b) {
					log.Debugf("Got event for %s with validation failure", e.key)
					return false
				}
			} else {
				log.Debugf("Got event for %s with unexpected type", e.key)
				return false
			}
		}

	}

	return true
}

func SetupWatcher(events []string) *Watcher {
	Validator := []EventValidator{}
	for _, e := range events {
		v := EventValidator{
			key:       e,
			Validator: DummyValidator,
		}

		Validator = append(Validator, v)
	}

	w := NewWatcher(Validator)

	w.Start()
	time.Sleep(1 * time.Second)
	log.Debugf("Watcher created for %+v", w)
	return w
}
