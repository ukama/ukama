/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package main

// import (
// "fmt"
// "syscall"

// "os"
// "os/signal"
// "time"

// "github.com/ukama/ukama/systems/notification/notify/internal"

// "github.com/ukama/ukama/systems/common/msgbus"

// "github.com/streadway/amqp"
// "github.com/ukama/ukama/systems/common/config"
// "google.golang.org/protobuf/proto"

// log "github.com/sirupsen/logrus"
// pb "github.com/ukama/ukama/systems/notification/notify/pb/gen"
// )

// var mockMsgClient msgbus.IMsgBus
// var mockHandlerName = "MockClient"
// var routingKey msgbus.RoutingKey

// [> Init read config <]
// func initConfig() {
// log.Infof("Initializing config")
// internal.ServiceConfig = internal.NewConfig()
// config.LoadConfig(internal.ServiceName, internal.ServiceConfig)
// }

// // Msgbus initialization
// func initMsgBus() {
// mockMsgClient, err := msgbus.NewConsumerClient(internal.ServiceConfig.Queue.Uri)
// if err != nil {
// failOnError(err, "Could not create a consumer. Error %s"+err.Error())
// }

// //
// key := msgbus.NewRoutingKeyBuilder().
// SetDeviceSource().
// SetContainer(internal.ServiceName).
// SetEventType().
// SetObject("notification").
// SetAction("*").
// MustBuild()

// routingKey = msgbus.RoutingKey(key)
// // Routing key
// routingKeys := []msgbus.RoutingKey{routingKey}

// log.Debugf("Mock:: msgClient: %+v", mockMsgClient)

// // Subscribe to exchange
// err = mockMsgClient.Subscribe(msgbus.DeviceQ.Queue, msgbus.DeviceQ.Exchange, msgbus.DeviceQ.ExchangeType, routingKeys, mockHandlerName, EvtMsgHandlerCB)
// failOnError(err, "Could not start subscribe to "+msgbus.DeviceQ.Exchange+msgbus.DeviceQ.ExchangeType)

// }

// //Handle Response message
// func EvtMsgHandlerCB(d amqp.Delivery, done chan<- bool) {

// log.Debugf("Mock::Event message handler for the mock %s  msg: %v", d.RoutingKey, d)

// // unmarshal
// switch msgbus.RoutingKey(d.RoutingKey) {

// // Add Device Response
// case routingKey:
// nMsg := &pb.NotificationMsg{}

// err := proto.Unmarshal(d.Body, nMsg)
// if err != nil {
// log.Errorf("Mock::Fail unmarshal: %s", d.Body)
// }

// [> Processing for the event can be done here <]
// log.Debugf("Mock::Received notification %s msg: %+v", routingKey, nMsg)

// }

// done <- true

// }

// // main
// func main() {
// // Log level
// log.SetLevel(log.TraceLevel)
// log.Debugf("Mock::Starting mock services..!!\n")

// // Read config
// initConfig()

// //Initialize msgbus
// initMsgBus()

// // Makes sure connection is closed when service exits.
// handleSigterm(func() {

// // Close connection
// if mockMsgClient != nil {
// mockMsgClient.Close()
// }

// })

// for {
// // Would wait for 5 seconds for reply
// time.Sleep(5 * time.Second)
// }

// }

// // Handles Ctrl+C or most other means of "controlled" shutdown gracefully.
// // Invokes the supplied func before exiting.
// func handleSigterm(handleExit func()) {
// c := make(chan os.Signal, 1)
// signal.Notify(c, os.Interrupt)
// signal.Notify(c, syscall.SIGTERM)

// go func() {
// <-c
// handleExit()
// os.Exit(1)
// }()

// }

// // Fatal error
// func failOnError(err error, msg string) {
// if err != nil {
// log.Errorf("Mock:: %s: %s", msg, err)
// panic(fmt.Sprintf("Mock:: %s: %s", msg, err))
// }
// }
