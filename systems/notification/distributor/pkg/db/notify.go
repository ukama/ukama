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
	"github.com/ukama/ukama/systems/common/notification"
	upb "github.com/ukama/ukama/systems/common/pb/gen/ukama"
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
	NodeId       string
	SubscriberId string
	Scopes       []notification.NotificationScope
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
	Register(orgId string, networkId string, subscriberId string, userId string, nodeId string, scopes []notification.NotificationScope) (string, *Sub)
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

func (h *notifyHandler) Register(orgId string, networkId string, subscriberId string, userId string, nodeId string, scopes []notification.NotificationScope) (string, *Sub) {

	sub := Sub{
		Id:           uuid.NewV4(),
		OrgId:        orgId,
		NetworkId:    networkId,
		SubscriberId: subscriberId,
		UserId:       userId,
		NodeId:       nodeId,
		Scopes:       scopes,
		DataChan:     make(chan *pb.Notification, BufferCapacity),
		QuitChan:     make(chan bool),
	}

	id := sub.Id.String()

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
			if n.Scope == upb.NotificationScope(scope) {
				/* Send over channel */
				s.DataChan <- n
			}
		}

	}
}

func (h *notifyHandler) notifyHandlerRoutine() {
	log.Infof("DB notify handler routine for %+v", h.Db)

	log.Infof("DB notify handler routine for %+v", h.Db)

	db, err := sql.Open(DbDriverName, "postgresql://"+h.Db.Username+":"+h.Db.Password+"@"+h.Db.Host+":"+strconv.Itoa(h.Db.Port)+"/"+h.Db.DbName+"?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	dbCS := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable", h.Db.Host, h.Db.Port, h.Db.DbName, h.Db.Username, h.Db.Password)

	log.Infof("Listening to user_notifications_channel from %s", dbCS)

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

	for {
		select {
		case notification := <-listener.Notify:
			log.Infof("DB notify received: %+v", notification)

			if notification == nil {
				log.Warn("Received nil notification, skipping")
				continue
			}

			params := strings.Split(notification.Extra, ",")
			log.Infof("Parsed notification params: %v", params)

			if len(params) < 5 {
				log.Errorf("Invalid notification format: %v", notification.Extra)
				continue
			}

			isRead, err := strconv.ParseBool(params[2])
			if err != nil {
				log.Errorf("Error parsing isRead: %v", err)
				continue
			}

			notificationId := params[1]
			log.Infof("Fetching notification details for ID: %s", notificationId)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			res, err := h.c.Get(ctx, &enpb.GetRequest{Id: notificationId})
			cancel()

			if err != nil {
				log.Errorf("Error getting notification: %v", err)
				continue
			}
			fmt.Println("VANESSA DATA :",res)

			un := &pb.Notification{
				IsRead:       isRead,
				Id:           res.Notification.Id,
				OrgId:        res.Notification.OrgId,
				Title:        res.Notification.Title,
				UserId:       res.Notification.UserId,
				NetworkId:    res.Notification.NetworkId,
				NodeId:       res.Notification.NodeId,
				Description:  res.Notification.Description,
				SubscriberId: res.Notification.SubscriberId,
				NodeStateId:  res.Notification.NodeStateId,
				Type:         upb.NotificationType(res.Notification.Type),
				Scope:        upb.NotificationScope(res.Notification.Scope),
				CreatedAt:    res.Notification.CreatedAt.AsTime().Format(time.RFC3339),
			}

			if res.Notification.NodeState != nil {
				un.NodeState = &pb.NodeState{
					Id:           res.Notification.NodeState.Id,
					Name:         res.Notification.NodeState.Name,
					NodeId:       res.Notification.NodeState.NodeId,
					CurrentState: res.Notification.NodeState.CurrentState,
					Latitude:     res.Notification.NodeState.Latitude,
					Longitude:    res.Notification.NodeState.Longitude,
				}
			}
			fmt.Println("VANESSA DATA :",un)

			h.processNotification(un)

		case <-h.done:
			log.Info("Stopping DB notify handler routine.")
			return
		}
	}
}
