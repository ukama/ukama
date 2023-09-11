package msgbus

import (
	"encoding/json"
	"fmt"
	"net/url"

	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/systems/common/rest"
)

const shovelEndpoint = "/api/shovels/"
const putShovelEndpoint = "/api/parameters/shovel/"

type MsgBusShovelProvider interface {
	AddShovel(name string, s *Shovel) error
	GetShovel(name string) (s *Shovel, err error)
	RemoveShovel(name string) error
	CreateShovel(name string, s *Shovel) error
	RestartShovel(name string) error
}

type Shovel struct {
	SrcProtocol     string `json:"src-protocol" default:"amqp091"`
	DestProtocol    string `default:"amqp091" json:"dest-protocol"`
	SrcExchange     string `default:"amq.topic" json:"src-exchange"`
	SrcExchangeKey  string `json:"src-exchange-key,omitempty"`
	DestExchange    string `default:"amq.topic" json:"dest-exchange,omitempty"`
	DestExchangeKey string `json:"dest-exchange-key,omitempty"`
	DestQueue       string `json:"dest-queue,omitempty"`
	SrcQueue        string `json:"src-queue,omitempty"`
	SrcUri          string `json:"src-uri"`
	DestUri         string `json:"dest-uri"`
	Name            string `json:"name,omitempty"`
	Status          string `json:"status,omitempty"`
	Vhost           string `json:"vhost,omitempty"`
	Type            string `json:"type,omitempty"`
}

type msgBusShovelClient struct {
	R        *rest.RestClient
	user     string
	password string
	name     string
	s        *Shovel
}

func NewShovelProvider(url string, debug bool, name, user, password, srcUri, destUri, srcExchange, destExchange, srcExchangeKey string) MsgBusShovelProvider {

	f, err := rest.NewRestClient(url, debug)
	if err != nil {
		log.Fatalf("Can't connect to %s url. Error %s", url, err.Error())
	}

	s := &Shovel{
		SrcExchange:    srcExchange,
		DestExchange:   destExchange,
		SrcUri:         srcUri,
		DestUri:        destUri,
		SrcExchangeKey: srcExchangeKey,
		SrcProtocol:    "amqp091",
		DestProtocol:   "amqp091",
	}

	if srcExchange == "" && destExchange == "" && srcExchangeKey == "" && srcUri == "" && destUri == "" {
		log.Fatalf("Required shovel parameter missing for shovel provider: %+v.", s)
	}

	p := &msgBusShovelClient{
		R:        f,
		name:     name,
		user:     user,
		password: password,
		s:        s,
	}

	return p
}

func (c *msgBusShovelClient) AddShovel(name string, s *Shovel) error {
	errStatus := &rest.ErrorMessage{}

	resp, err := c.R.C.R().
		SetBasicAuth(c.user, c.password).
		SetError(errStatus).
		SetBody(map[string]interface{}{"value": *s}).
		Put(c.R.URL.String() + putShovelEndpoint + url.PathEscape("/") + "/" + name)

	if err != nil {
		log.Errorf("Failed to send api request to msgbus. Error %s", err.Error())
		return fmt.Errorf("api request to msgbus system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to add shovel %s to msgbus. HTTP resp code %d and Error message is %s", name, resp.StatusCode(), errStatus.Message)
		return fmt.Errorf("failed adding shovel %s to msgbus. Error %s", name, errStatus.Message)
	}

	log.Infof("Shovel %s added to msgbus.", name)

	return nil
}

func (c *msgBusShovelClient) GetShovel(name string) (*Shovel, error) {
	errStatus := &rest.ErrorMessage{}

	s := &[]Shovel{}
	resp, err := c.R.C.R().
		SetBasicAuth(c.user, c.password).
		SetError(errStatus).
		Get(c.R.URL.String() + shovelEndpoint + "vhost/" + url.PathEscape("/") + "/" + name)

	if err != nil {
		log.Errorf("Failed to send api request to msgbus. Error %s", err.Error())
		return nil, fmt.Errorf("api request to msgbus system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to reading shovel %s to msgbus. HTTP resp code %d and Error message is %s", name, resp.StatusCode(), errStatus.Message)
		return nil, fmt.Errorf("failed reading shovel %s to msgbus. Error %s", name, errStatus.Message)
	}

	err = json.Unmarshal(resp.Body(), s)
	if err != nil {
		log.Tracef("Failed to deserialize shovel info. Error message is %s", err.Error())

		return nil, fmt.Errorf("shovel info deserailization failure: %w", err)
	}

	log.Infof("Shovel %s info %+v.", name, s)

	return &(*s)[0], nil
}

func (c *msgBusShovelClient) RemoveShovel(name string) error {
	errStatus := &rest.ErrorMessage{}

	resp, err := c.R.C.R().
		SetBasicAuth(c.user, c.password).
		SetError(errStatus).
		Delete(c.R.URL.String() + shovelEndpoint + "vhost/" + url.PathEscape("/") + "/" + name)

	if err != nil {
		log.Errorf("Failed to send api request to msgbus. Error %s", err.Error())
		return fmt.Errorf("api request to msgbus system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to remove shovel %s to msgbus. HTTP resp code %d and Error message is %s", name, resp.StatusCode(), errStatus.Message)
		return fmt.Errorf("failed deleting shovel %s to msgbus. Error %s", name, errStatus.Message)
	}

	log.Infof("Shovel %s remove from msgbus.", name)

	return nil
}

func (c *msgBusShovelClient) RestartShovel(name string) error {
	errStatus := &rest.ErrorMessage{}

	resp, err := c.R.C.R().
		SetBasicAuth(c.user, c.password).
		SetError(errStatus).
		Delete(c.R.URL.String() + shovelEndpoint + "vhost/" + url.PathEscape("/") + "/" + name + "/restart")

	if err != nil {
		log.Errorf("Failed to send api request to msgbus. Error %s", err.Error())
		return fmt.Errorf("api request to msgbus system failure: %w", err)
	}

	if !resp.IsSuccess() {
		log.Errorf("Failed to restart shovel %s to msgbus. HTTP resp code %d and Error message is %s", name, resp.StatusCode(), errStatus.Message)
		return fmt.Errorf("failed restart shovel %s to msgbus. Error %s", name, errStatus.Message)
	}

	log.Infof("Shovel %s restarted in msgbus.", name)

	return nil
}

func (c *msgBusShovelClient) CreateShovel(name string, ns *Shovel) error {

	s, err := c.GetShovel(name)
	if err == nil && s != nil {
		log.Infof("Shovel %s already exists with %+v ", name, s)
		return nil
	} else {
		log.Infof("Creating shovel %s with %+v", name, s)
	}

	if s != nil {
		err = c.AddShovel(name, ns)
		if err == nil {
			log.Infof("Created shovel %s with %+v", name, ns)
		}
	} else {
		err = c.AddShovel(name, c.s)
		if err == nil {
			log.Infof("Created shovel %s with %+v", name, c.s)
		}
	}

	return err
}
