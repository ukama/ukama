/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package messages

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"google.golang.org/protobuf/reflect/protoreflect"

	epb "github.com/ukama/ukama/systems/common/pb/gen/events"
	u "github.com/ukama/ukama/systems/common/ukama"
)

const (
	NodePort = 7070
	MeshPort = 7071
)

func NewNodeOnline(data string) (protoreflect.ProtoMessage, error) {
	nodeOnline := &epb.NodeOnlineEvent{
		NodeId:       string(u.NewVirtualNodeId(u.NODE_ID_TYPE_TOWERNODE)),
		NodeIp:       gofakeit.IPv4Address(),
		NodePort:     NodePort,
		MeshIp:       gofakeit.IPv4Address(),
		MeshPort:     MeshPort,
		MeshHostName: gofakeit.DomainName(),
	}

	if data != "" {
		err := updateProto(nodeOnline, data)
		if err != nil {
			return nil, fmt.Errorf("failed to update event proto: %w", err)
		}
	}

	return nodeOnline, nil
}

func NewNodeOffline(data string) (protoreflect.ProtoMessage, error) {
	nodeOffline := &epb.NodeOfflineEvent{
		NodeId: string(u.NewVirtualNodeId(u.NODE_ID_TYPE_TOWERNODE)),
	}

	if data != "" {
		err := updateProto(nodeOffline, data)
		if err != nil {
			return nil, fmt.Errorf("failed to update event proto: %w", err)
		}
	}

	return nodeOffline, nil
}

func NewNodeAssign(data string) (protoreflect.ProtoMessage, error) {
	nodeAssign := &epb.NodeAssignedEvent{
		NodeId:  string(u.NewVirtualNodeId(u.NODE_ID_TYPE_TOWERNODE)),
		Type:    u.NODE_ID_TYPE_TOWERNODE,
		Network: gofakeit.UUID(),
		Site:    gofakeit.UUID(),
	}

	if data != "" {
		err := updateProto(nodeAssign, data)
		if err != nil {
			return nil, fmt.Errorf("failed to update event proto: %w", err)
		}
	}

	return nodeAssign, nil
}

func NewNodeRelease(data string) (protoreflect.ProtoMessage, error) {
	nodeRelease := &epb.NodeReleasedEvent{
		NodeId:  string(u.NewVirtualNodeId(u.NODE_ID_TYPE_TOWERNODE)),
		Type:    u.NODE_ID_TYPE_TOWERNODE,
		Network: gofakeit.UUID(),
		Site:    gofakeit.UUID(),
	}

	if data != "" {
		err := updateProto(nodeRelease, data)
		if err != nil {
			return nil, fmt.Errorf("failed to update event proto: %w", err)
		}
	}

	return nodeRelease, nil
}
