/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package ukama

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeId_GetNodeType(t *testing.T) {
	t.Run("NodeTypeValid", func(tt *testing.T) {
		ntypes := map[string]string{
			"HNode":     NODE_ID_TYPE_HOMENODE,
			"hnode":     NODE_ID_TYPE_HOMENODE,
			"HOMENODE":  NODE_ID_TYPE_HOMENODE,
			"Towernode": NODE_ID_TYPE_TOWERNODE,
			"tnode":     NODE_ID_TYPE_TOWERNODE,
			"TNODE":     NODE_ID_TYPE_TOWERNODE,
			"ANode":     NODE_ID_TYPE_AMPNODE,
			"anode":     NODE_ID_TYPE_AMPNODE,
			"ANODE":     NODE_ID_TYPE_AMPNODE}

		for k, v := range ntypes {
			nodeId := NewVirtualNodeId(k)
			assert.Equal(tt, v, nodeId.GetNodeType())
		}
	})

	t.Run("NodeTypeUndefined", func(tt *testing.T) {
		nodeId := NewVirtualNodeId("some_weird_type")
		assert.Equal(tt, NODE_ID_TYPE_UNDEFINED, nodeId.GetNodeType())
	})
}

func TestModuleId_StringLowercase(t *testing.T) {
	t.Run("ModuleTypeValid", func(tt *testing.T) {
		ntypes := map[string]string{
			"comv1": MODULE_ID_TYPE_COMP,
			"trx":   MODULE_ID_TYPE_TRX,
			"ctrl":  MODULE_ID_TYPE_CTRL,
			"fe":    MODULE_ID_TYPE_FE,
			"undef": MODULE_ID_TYPE_UNDEFINED}

		for k, v := range ntypes {
			moduleId := NewVirtualModuleId(k)
			assert.Contains(tt, moduleId.StringLowercase(), v)
		}
	})

	t.Run("ModuleTypeUndefined", func(tt *testing.T) {
		moduleId := NewVirtualModuleId("some_weird_type")
		assert.Contains(tt, moduleId.StringLowercase(), NODE_ID_TYPE_UNDEFINED)
	})
}

func TestGetNodeType(t *testing.T) {
	t.Run("NodeTypeValid", func(tt *testing.T) {
		ntypes := map[string]string{
			"HNode":     NODE_ID_TYPE_HOMENODE,
			"hnode":     NODE_ID_TYPE_HOMENODE,
			"HOMENODE":  NODE_ID_TYPE_HOMENODE,
			"Towernode": NODE_ID_TYPE_TOWERNODE,
			"tnode":     NODE_ID_TYPE_TOWERNODE,
			"TNODE":     NODE_ID_TYPE_TOWERNODE,
			"ANode":     NODE_ID_TYPE_AMPNODE,
			"anode":     NODE_ID_TYPE_AMPNODE,
			"ANODE":     NODE_ID_TYPE_AMPNODE}

		for k, v := range ntypes {
			nodeId := NewVirtualNodeId(k)
			nodeType := GetNodeType(nodeId.StringLowercase())
			assert.Equal(tt, v, *nodeType)
		}
	})
}

func TestValidateNodeId(t *testing.T) {
	t.Run("NodeIdIsValid", func(tt *testing.T) {
		nodeId := "UK-SA2156-HNODE-A1-XXXX"

		uid, err := ValidateNodeId(string(nodeId))

		assert.NoError(tt, err)
		assert.Equal(tt, strings.ToLower(nodeId), string(uid))
	})

	t.Run("ValidateNodeIdCase1", func(tt *testing.T) {
		nodeId := "UK-SA2156"

		_, err := ValidateNodeId(string(nodeId))
		assert.Error(tt, err)
	})

	t.Run("ValidateNodeIdCase2", func(tt *testing.T) {
		nodeId := "UK-SA2156-CNODE-A1-XXXX"

		_, err := ValidateNodeId(string(nodeId))
		assert.Error(tt, err)
	})
}
