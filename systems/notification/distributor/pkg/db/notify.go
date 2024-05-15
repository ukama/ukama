package db

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"

	log "github.com/sirupsen/logrus"

	uconf "github.com/ukama/ukama/systems/common/config"
	"github.com/ukama/ukama/systems/common/uuid"
	pb "github.com/ukama/ukama/systems/notification/distributor/pb/gen"
	"github.com/ukama/ukama/systems/notification/distributor/pkg/providers"
	enpb "github.com/ukama/ukama/systems/notification/event-notify/pb/gen"
)

const (
	DbDriverName   = "postgres"
	BufferCapacity = 10
)

type Sub struct {
	Id           uuid.UUID
	OrgId        string
	NetworkId    string
	UserId       string
	SubscriberId string
	Scopes       []pb.NotificationScope
	DataChan     chan *pb.Notification
	QuitChan     chan bool
}

type Subs map[string]Sub

type notifyHandler struct {
	Db                      *uconf.Database
	c                       enpb.EventToNotifyServiceClient
	minReconnectionInterval time.Duration `default:"10s"`
	maxReconnectionInterval time.Duration `default:"1m"`
	done                    chan bool
	subs                    Subs
}

type NotifyHandler interface {
	Register(orgId string, networkId string, subscriberId string, userId string, scopes []string) (string, *Sub)
	Deregister(id string) error
	Start()
	Stop()
}

func NewNotifyHandler(db *uconf.Database, c providers.EventNotifyClientProvider) *notifyHandler {
	svc, err := c.GetClient()
	if err != nil {
		log.Fatalf("Failed to get event notifu client: %v", err)
	}

	return &notifyHandler{
		Db:   db,
		c:    svc,
		done: make(chan bool),
		subs: make(Subs),
	}
}

func (h *notifyHandler) Register(orgId string, networkId string, subscriberId string, userId string, scopes []string) (string, *Sub) {
	id := uuid.NewV4().String()
	sub := Sub{
		OrgId:        orgId,
		NetworkId:    networkId,
		SubscriberId: subscriberId,
		UserId:       userId,
		Scopes:       make([]pb.NotificationScope, len(scopes)),
		DataChan:     make(chan *pb.Notification, BufferCapacity),
		QuitChan:     make(chan bool),
	}

	for _, s := range scopes {
		sid := pb.NotificationScope(pb.NotificationScope_value[s])
		if sid != pb.NotificationScope_UNKOWN_SCOPE {
			sub.Scopes = append(sub.Scopes, sid)
		}
	}

	h.subs[id] = sub

	log.Infof("Registerd %s sub with %+v to the notify handler", id, sub)

	return id, &sub
}

func (h *notifyHandler) Deregister(id string) error {
	s, ok := h.subs[id]
	if !ok {
		log.Errorf("Sub with id %s not found", id)
		return fmt.Errorf("sub with id %s not found", id)
	}

	log.Infof("Deleting sub %s with %+v from notify handler", id, s)

	delete(h.subs, id)

	return nil
}

func (h *notifyHandler) Start() {
	go h.notifyHandlerRoutine()
}

func (h *notifyHandler) Stop() {
	log.Infof("Stopping the notify handler routine")
	h.done <- true

	/* Cleaning all the sub */
	for k, s := range h.subs {
		log.Infof("Stopping sub %s with %+v", k, s)
		s.QuitChan <- true
	}
}

func (h *notifyHandler) processNotification(n *pb.Notification) {
	for k, s := range h.subs {
		log.Infof("Processing notification %+v for sub %s with %+v", n, k, s)

		for _, scope := range s.Scopes {
			if n.Scope == scope {
				/* Send over channel */
				s.DataChan <- n
			}
		}

	}
}

func (h *notifyHandler) notifyHandlerRoutine() {
	log.Info("DB notify handler routine")

	db, err := sql.Open(DbDriverName, "postgresql://"+h.Db.Username+":"+h.Db.Password+"@"+h.Db.Host+":"+strconv.Itoa(h.Db.Port)+"/"+h.Db.DbName+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbCS := fmt.Sprintf("dbname=%s user=%s password=%s sslmode=disable", h.Db.DbName, h.Db.Username, h.Db.Password)

	listener := pq.NewListener(dbCS, h.minReconnectionInterval, h.maxReconnectionInterval, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Println(err.Error())
		}
	})

	err = listener.Listen("user_notifications_channel")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	/*TODO: - Close the stream
	- May be check where the notifivcatins are getting filtered based on userid/kid or subscriberid
	- This will only report notifcation when websocket is connected if we have any old notification(stores 8Gb)  that had to be reterived by
	anyother API method
	- Looks like if this is session/user based we might not get trigeer properly beacuse all of the listner will be reading form the same notify queue.
	*/
	for {
		select {
		case notification := <-listener.Notify:
			log.Infof("DB notify received for %+v", notification)

			/* Parse DB trigger details */
			params := strings.Split(notification.Extra, ",")
			isRead, _ := strconv.ParseBool(params[2])

			/* Get notifcation detaild fron event-notify service */
			res, err := h.c.Get(context.Background(), &enpb.GetRequest{Id: params[1]})
			if err != nil {
				log.Errorf("Error getting notification: %v", err)
				continue
			}

			un := &pb.Notification{
				IsRead:       isRead,
				Id:           res.Notification.Id,
				OrgId:        res.Notification.OrgId,
				Title:        res.Notification.Title,
				UserId:       res.Notification.UserId,
				NetworkId:    res.Notification.NetworkId,
				Description:  res.Notification.Description,
				SubscriberId: res.Notification.SubscriberId,
				ForRole:      pb.RoleType(res.Notification.ForRole),
				Type:         pb.NotificationType(res.Notification.Type),
				Scope:        pb.NotificationScope(res.Notification.Scope),
			}
			log.Infof("Notification is %+v", un)

			h.processNotification(un)
			/* send notificaton to the registered clients */
			// err = srv.Send(&un)
			// if err != nil {
			// 	log.Errorf("Error sending notification: %v", err)
			// 	continue
			// }

		case <-h.done:
			log.Infof("Stopping Db notify handler routine.")
		}
	}
}
