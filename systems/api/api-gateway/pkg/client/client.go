package client

import (
	"errors"
	"net/http"

	"github.com/ukama/ukama/systems/common/rest"

	log "github.com/sirupsen/logrus"
)

type Client interface {
	GetNetwork(string) (*NetworkInfo, error)
	CreateNetwork(string, string, []string, []string, float64, float64, uint32, bool) (*NetworkInfo, error)

	GetPackage(string) (*PackageInfo, error)
	AddPackage(string, string, string, string, string, string, bool, bool, int64, int64, int64, string,
		string, string, string, string, uint64, float64, float64, float64, uint32, []string) (*PackageInfo, error)

	GetSim(string) (*SimInfo, error)
	ConfigureSim(string, string, string, string, string, string, string, string, string, string,
		string, string, string, string, uint32) (*SimInfo, error)
}

type clients struct {
	network    NetworkClient
	pkg        PackageClient
	subscriber SubscriberClient
	sim        SimClient
}

func NewClientsSet(network NetworkClient, pkg PackageClient, subscriber SubscriberClient,
	sim SimClient) Client {
	c := &clients{
		network:    network,
		pkg:        pkg,
		subscriber: subscriber,
		sim:        sim,
	}

	return c
}

func (c *clients) GetNetwork(id string) (*NetworkInfo, error) {
	net, err := c.network.Get(id)
	if err != nil {
		return nil, handleRestErrorStatus(err)
	}

	if !net.IsSynced {
		log.Warn("partial content. request is still ongoing")

		return net, rest.HttpError{
			HttpCode: http.StatusPartialContent,
			Message:  "partial content. request is still ongoing",
		}
	}

	return net, nil
}

func (c *clients) CreateNetwork(orgName, NetworkName string, allowedCountries,
	allowedNetworks []string, budget, overdraft float64, trafficPolicy uint32,
	paymentLinks bool) (*NetworkInfo, error) {
	net, err := c.network.Add(AddNetworkRequest{
		OrgName:          orgName,
		NetName:          NetworkName,
		AllowedCountries: allowedCountries,
		AllowedNetworks:  allowedNetworks,
		Budget:           budget,
		Overdraft:        overdraft,
		TrafficPolicy:    trafficPolicy,
		PaymentLinks:     paymentLinks,
	})
	if err != nil {
		return nil, handleRestErrorStatus(err)
	}

	return net, nil
}

func (c *clients) GetPackage(id string) (*PackageInfo, error) {
	pkg, err := c.pkg.Get(id)
	if err != nil {
		return nil, handleRestErrorStatus(err)
	}

	if !pkg.IsSynced {
		log.Warn("partial content. request is still ongoing")

		return pkg, rest.HttpError{
			HttpCode: http.StatusPartialContent,
			Message:  "partial content. request is still ongoing",
		}
	}

	return pkg, nil
}

func (c *clients) AddPackage(name, orgId, ownerId, from, to, baserateId string,
	isActive, flatRate bool, smsVolume, voiceVolume, dataVolume int64, voiceUnit, dataUnit,
	simType, apn, pType string, duration uint64, markup, amount, overdraft float64, trafficPolicy uint32,
	networks []string) (*PackageInfo, error) {

	pkg, err := c.pkg.Add(AddPackageRequest{
		Name:          name,
		OrgId:         orgId,
		OwnerId:       ownerId,
		From:          from,
		To:            to,
		BaserateId:    baserateId,
		Active:        isActive,
		SmsVolume:     smsVolume,
		VoiceVolume:   voiceVolume,
		DataVolume:    dataVolume,
		VoiceUnit:     voiceUnit,
		DataUnit:      dataUnit,
		SimType:       simType,
		Apn:           apn,
		Markup:        markup,
		Type:          pType,
		Flatrate:      flatRate,
		Amount:        amount,
		Overdraft:     overdraft,
		TrafficPolicy: trafficPolicy,
		Networks:      networks,
	})
	if err != nil {
		return nil, handleRestErrorStatus(err)
	}

	return pkg, nil
}

func (c *clients) GetSim(id string) (*SimInfo, error) {
	sim, err := c.sim.Get(id)
	if err != nil {
		return nil, handleRestErrorStatus(err)
	}

	if !sim.IsSynced {
		log.Warn("partial content. request is still ongoing")

		return sim, rest.HttpError{
			HttpCode: http.StatusPartialContent,
			Message:  "partial content. request is still ongoing",
		}
	}

	return sim, nil
}

func (c *clients) ConfigureSim(subscriberId, orgId, networkId, firstName, lastName,
	email, phoneNumber, address, dob, proofOfID, idSerial, packageId, simType,
	simToken string, trafficPolicy uint32) (*SimInfo, error) {
	if subscriberId == "" {
		subscriber, err := c.subscriber.Add(
			AddSubscriberRequest{
				OrgId:                 orgId,
				NetworkId:             networkId,
				FirstName:             firstName,
				LastName:              lastName,
				Email:                 email,
				PhoneNumber:           phoneNumber,
				Address:               address,
				Dob:                   dob,
				ProofOfIdentification: proofOfID,
				IdSerial:              idSerial,
			})
		if err != nil {
			log.Error("Failed to create new subscriber while configuring sim")

			return nil, err
		}

		subscriberId = subscriber.SubscriberId.String()
	}

	sim, err := c.sim.Add(AddSimRequest{
		SubscriberId:  subscriberId,
		NetworkId:     networkId,
		PackageId:     packageId,
		SimType:       simType,
		SimToken:      simToken,
		TrafficPolicy: trafficPolicy,
	})
	if err != nil {
		return nil, handleRestErrorStatus(err)
	}

	return sim, nil
}

func handleRestErrorStatus(err error) error {
	e := ErrorStatus{}

	if errors.As(err, &e) {
		return rest.HttpError{
			HttpCode: e.StatusCode,
			Message:  err.Error(),
		}
	}

	return err
}
