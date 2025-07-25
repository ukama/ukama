// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	amqp "github.com/streadway/amqp"
	mock "github.com/stretchr/testify/mock"

	msgbus "github.com/ukama/ukama/systems/common/msgbus"
)

// Consumer is an autogenerated mock type for the Consumer type
type Consumer struct {
	mock.Mock
}

// Close provides a mock function with no fields
func (_m *Consumer) Close() {
	_m.Called()
}

// IsClosed provides a mock function with no fields
func (_m *Consumer) IsClosed() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for IsClosed")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Subscribe provides a mock function with given fields: queueName, exchangeName, exchangeType, routingKeys, consumerName, handlerFunc
func (_m *Consumer) Subscribe(queueName string, exchangeName string, exchangeType string, routingKeys []msgbus.RoutingKey, consumerName string, handlerFunc func(amqp.Delivery, chan<- bool)) error {
	ret := _m.Called(queueName, exchangeName, exchangeType, routingKeys, consumerName, handlerFunc)

	if len(ret) == 0 {
		panic("no return value specified for Subscribe")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, []msgbus.RoutingKey, string, func(amqp.Delivery, chan<- bool)) error); ok {
		r0 = rf(queueName, exchangeName, exchangeType, routingKeys, consumerName, handlerFunc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubscribeToQueue provides a mock function with given fields: queueName, consumerName, handlerFunc
func (_m *Consumer) SubscribeToQueue(queueName string, consumerName string, handlerFunc func(amqp.Delivery, chan<- bool)) error {
	ret := _m.Called(queueName, consumerName, handlerFunc)

	if len(ret) == 0 {
		panic("no return value specified for SubscribeToQueue")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, func(amqp.Delivery, chan<- bool)) error); ok {
		r0 = rf(queueName, consumerName, handlerFunc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubscribeToServiceQueue provides a mock function with given fields: serviceName, exchangeName, routingKeys, consumerId, handlerFunc
func (_m *Consumer) SubscribeToServiceQueue(serviceName string, exchangeName string, routingKeys []msgbus.RoutingKey, consumerId string, handlerFunc func(amqp.Delivery, chan<- bool)) error {
	ret := _m.Called(serviceName, exchangeName, routingKeys, consumerId, handlerFunc)

	if len(ret) == 0 {
		panic("no return value specified for SubscribeToServiceQueue")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, []msgbus.RoutingKey, string, func(amqp.Delivery, chan<- bool)) error); ok {
		r0 = rf(serviceName, exchangeName, routingKeys, consumerId, handlerFunc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SubscribeWithArgs provides a mock function with given fields: queueName, exchangeName, exchangeType, routingKeys, consumerName, queueArgs, handlerFunc
func (_m *Consumer) SubscribeWithArgs(queueName string, exchangeName string, exchangeType string, routingKeys []msgbus.RoutingKey, consumerName string, queueArgs map[string]interface{}, handlerFunc func(amqp.Delivery, chan<- bool)) error {
	ret := _m.Called(queueName, exchangeName, exchangeType, routingKeys, consumerName, queueArgs, handlerFunc)

	if len(ret) == 0 {
		panic("no return value specified for SubscribeWithArgs")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, []msgbus.RoutingKey, string, map[string]interface{}, func(amqp.Delivery, chan<- bool)) error); ok {
		r0 = rf(queueName, exchangeName, exchangeType, routingKeys, consumerName, queueArgs, handlerFunc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewConsumer creates a new instance of Consumer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConsumer(t interface {
	mock.TestingT
	Cleanup(func())
}) *Consumer {
	mock := &Consumer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
