package lwm2m

import (
	"bytes"
	"encoding/binary"
	cfg "lwm2m-gateway/pkg/config"
	"math/rand"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	CONNECTION_TIMEOUT time.Duration = 5 * time.Second
)

type RequestMsg struct {
	Token   uint32
	Length  uint32
	Message string
}

type ResponseMsg struct {
	Token   uint32
	Status  uint32
	Format  uint32
	Length  uint32
	Message string
}

type EventMsg struct {
	Token   uint32
	Uuid    string
	Uri     string
	Count   uint32
	Format  uint32
	Length  uint32
	Message string
}

type EventMsgResponse struct {
	Token  uint32
	Status uint32
}

// Prepare message for LwM2M Server
func PrepareMsgForLwm2mServer(op Lwm2mop, uuid *string, urlext *string, extraParam *string) *RequestMsg {

	// Handles URI like read UUID /3328/0/5821 or write UUID //3328/0/5821 1900
	msg := string(op) + " " + *uuid + " " + *urlext + " "
	if extraParam != nil {
		msg = msg + *extraParam
	}

	//Length of the message
	length := uint32(len(msg))

	//Request message
	reqMsg := &RequestMsg{
		Token:   rand.Uint32(),
		Length:  length,
		Message: msg,
	}

	log.Infof("TRX::Message prepared for Lwm2m Server %d Device %s [%d] and Data is %s Length is %d", reqMsg.Token, *uuid, len(*uuid), reqMsg.Message, length)
	return reqMsg
}

func EncodeReqToBytes(p RequestMsg) []byte {

	buf := make([]byte, 8)
	// Encode token
	binary.LittleEndian.PutUint32(buf[0:], p.Token)

	// Encode length
	binary.LittleEndian.PutUint32(buf[4:], p.Length)

	// Convert string to byte array
	log.Debugf("TRX::Length of message is  %d", len(p.Message))
	str := []byte(p.Message)

	//Append to final slice
	txstr := append(buf, str...)

	log.Debugf("TRX::Encoded message is %+v with length %d", txstr, len(txstr))

	return txstr
}

//Decode response message
func DecodeRespFromBytes(buf []byte) ResponseMsg {

	var resp ResponseMsg

	// Decode token
	resp.Token = binary.LittleEndian.Uint32(buf[0:])

	// Decode status
	resp.Status = binary.LittleEndian.Uint32(buf[4:])

	//Decode Format
	resp.Format = binary.LittleEndian.Uint32(buf[8:])

	// Decode length
	resp.Length = binary.LittleEndian.Uint32(buf[12:])

	//Decode message
	resp.Message = string(buf[16 : 16+resp.Length])

	log.Debugf("TRX::Decoded response message is %+v", resp)

	return resp
}

/* Decode Event Message */
func DecodeEventMsg(buf []byte) *EventMsg {
	var evt EventMsg

	// Decode token
	evt.Token = binary.LittleEndian.Uint32(buf[0:])

	// Decode UUID c style string
	n := bytes.Index(buf[4:35], []byte{0})
	evt.Uuid = string(buf[4 : n+4])

	//Decode URL c stye string
	n = bytes.Index(buf[36:68], []byte{0})
	evt.Uri = string(buf[36 : n+36])

	// Decode Count
	evt.Count = binary.LittleEndian.Uint32(buf[68:])

	// Decode Format
	evt.Format = binary.LittleEndian.Uint32(buf[72:])

	// Decode Length
	evt.Length = binary.LittleEndian.Uint32(buf[76:])

	//Decode message
	evt.Message = string(buf[80 : 80+evt.Length])

	log.Debugf("TRX::Decoded Response message is %+v", evt)

	return &evt
}

/* Encode Event response */
func EncodeEventMsg(r EventMsgResponse) []byte {

	res := make([]byte, 8)
	// Encode token
	binary.LittleEndian.PutUint32(res[0:], r.Token)

	// Encode length
	binary.LittleEndian.PutUint32(res[4:], r.Status)

	return res
}

// Send the request to LwM2M server.
func transmit(uuid string, message RequestMsg) (*ResponseMsg, error) {

	// Connect
	lwm2mserver := cfg.Config.Lwm2mServer.Address + ":" + cfg.Config.Lwm2mServer.Port
	conn, err := net.Dial("tcp", lwm2mserver)
	if err != nil {
		log.Errorf("TRX::Failed connecting to server %s . Error:: %s", lwm2mserver, err.Error())
		return nil, err
	}
	log.Debugf("TRX::Message for lwM2M server %s is %+v", lwm2mserver, message)

	// Encode
	msg := EncodeReqToBytes(message)

	//Set read time out
	err = conn.SetReadDeadline(time.Now().Add(CONNECTION_TIMEOUT))
	if err != nil {
		log.Errorf("TRX:: Failed to set response time out for the Lwm2m request.")
	}

	// Write to server
	_, err = conn.Write(msg)
	if err != nil {
		log.Errorf("TRX::Send to server failed:: Error %+v", err.Error())
		return nil, err
	}

	reply := make([]byte, 4096)

	lenr, rerr := conn.Read(reply)
	if err != nil {
		log.Debugf("TRX::Send to LwM2M server failed. Error:: %+v", rerr)
		return nil, err
	}
	log.Debugf("TRX::Reply message received from LwM2M server is :: %+v", reply[:lenr])

	// Decode Response message
	resp := DecodeRespFromBytes(reply)

	log.Infof("TRX::Response message received from LwM2M server for req id %d with status %d.", resp.Token, resp.Status)
	conn.Close()

	return &resp, nil
}

// Handle the event msg received
func receive_handler(data []byte) []byte {
	// Decode Event
	var resp EventMsgResponse
	if len(data) > 0 {
		evt := DecodeEventMsg(data)
		log.Infof("TRX::Event received from LwM2M server for unit %s URI: %s Event Id is %d", evt.Uuid, evt.Uri, evt.Token)

		//Add further handling.
		go handleEvent(evt)

		//Create response
		resp = EventMsgResponse{
			Token:  evt.Token,
			Status: COAP_201_CREATED,
		}
	} else {
		log.Debugf("TRX::Empty Event Received.")
		//Create response
		resp = EventMsgResponse{
			Token:  0,
			Status: COAP_402_BAD_OPTION,
		}
	}

	//Encode Response
	rep := EncodeEventMsg(resp)

	return rep
}

// Handling a notification connection
func handleNewNotification(conn net.Conn) {

	log.Debugf("TRX::Serving %s\n", conn.RemoteAddr().String())

	event := make([]byte, 1024)

	//Set read time out
	err := conn.SetDeadline(time.Now().Add(CONNECTION_TIMEOUT))
	if err != nil {
		log.Errorf("TRX::TRX:: Failed to set response time out for the Lwm2m request.")
	}

	// Get message, output
	lenr, err := conn.Read(event)
	if err != nil {
		log.Debugf("TRX::Error while waiting for message from LwM2M server:: %+v", err)
	}
	log.Debugf("TRX::Notification received from LwM2M server %+v received bytes %d", event[:lenr], lenr)

	// Handling Notification
	reply := receive_handler(event)

	// responding
	_, err = conn.Write(reply)
	if err != nil {
		log.Debugf("TRX::Error writingreply:: %+v", err)
	}

	log.Debugf("TRX::Message Responded to LwM2M server:: %+v", reply)

	conn.Close()
}

// Receiver for events
func Receiver() {
	// listen on port 3100
	gatewayserver := cfg.Config.Lwm2mGateway.Address + ":" + cfg.Config.Lwm2mGateway.Port
	log.Infof("TRX::Starting tcp server at %s to receive data from LwM2M server.", gatewayserver)
	ln, err := net.Listen("tcp", gatewayserver)
	if err != nil {
		log.Errorf("TRX::Failed to listen on %s notification server. Error :: %+v", gatewayserver, err)
		return
	}
	defer ln.Close()

	// accept connection
	for {
		log.Debugf("TRX::Waiting for new connection...")

		conn, err := ln.Accept()
		if err != nil {
			log.Errorf("TRX::Error while accepting new connection for LwM2M Gateway Notification server. Error: %+v", err)
			return
		}

		// Handling new notification request.
		go handleNewNotification(conn)
	}

}
