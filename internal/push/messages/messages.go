/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package messages

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/oleiade/reflections"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
)

var RoutingMap = map[string]func(string) (protoreflect.ProtoMessage, error){
	"subscriber.registry.subscriber.create":   NewSubscriberCreate,
	"subscriber.registry.subscriber.update":   NewSubscriberUpdate,
	"subscriber.registry.subscriber.delete":   NewSubscriberDelete,
	"dataplan.package.package.create":         NewPackageCreate,
	"subscriber.simmanager.sim.allocate":      NewSimAllocate,
	"subscriber.simmanager.sim.activepackage": NewSetActivePackageForSim,
	"subscriber.simmanager.sim.expirepackage": NewSimPackageExpire,
	"subscriber.simmanager.sim.usage":         NewSimUsage,
}

func WrapProto(f func(string) (protoreflect.ProtoMessage, error), data string) (*anypb.Any, error) {
	p, err := f(data)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated event message as proto: %w", err)
	}

	anyE, err := anypb.New(p)
	if err != nil {
		return nil, fmt.Errorf("failed to marshall event message as proto: %w", err)
	}

	return anyE, nil
}

func getData(msg string) (map[string]any, error) {
	m := make(map[string]any)

	err := json.Unmarshal([]byte(msg), &m)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal provided payload %q. Error: %w", msg, err)
	}

	return m, nil
}

func updateNumericField(obj any, fieldname string, val float64) error {
	rv := reflect.ValueOf(obj).Elem()

	fieldValue := rv.FieldByName(fieldname)
	if !fieldValue.IsValid() || !fieldValue.CanSet() {
		return fmt.Errorf("field info [%s] not found from event message", fieldname)
	}

	switch fieldValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fieldValue.OverflowInt(int64(val)) {
			return fmt.Errorf("can't assign value due to %s-overflow", fieldValue.Kind())
		}
		fieldValue.SetInt(int64(val))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if fieldValue.OverflowUint(uint64(val)) {
			return fmt.Errorf("can't assign value due to %s-overflow", fieldValue.Kind())
		}
		fieldValue.SetUint(uint64(val))
	case reflect.Float32, reflect.Float64:
		if fieldValue.OverflowFloat(val) {
			return fmt.Errorf("can't assign value due to %s-overflow", fieldValue.Kind())
		}
		fieldValue.SetFloat(val)
	default:
		return fmt.Errorf("can't assign value to a non-number type")
	}

	return nil
}

func updateProto(p any, d string) error {
	data, err := getData(d)
	if err != nil {
		return fmt.Errorf("failed to get data from event message: %w", err)
	}

	for k, v := range data {
		fk, err := reflections.GetFieldKind(p, k)
		if err != nil {
			return fmt.Errorf("failed to get field info [%s] from event message: %w", k, err)
		}

		if val, ok := v.(float64); ok && fk != reflect.Float64 {
			updateNumericField(p, k, val)
		} else {
			err = reflections.SetField(p, k, v)
			if err != nil {
				return fmt.Errorf("failed to update field info [%s: %v] from event message: %w", k, v, err)
			}
		}
	}

	return nil
}
