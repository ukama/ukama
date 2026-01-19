/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package ukama

import "database/sql/driver"
 
 type NodeType string
 
 const (
	 NODE_TYPE_HOMENODE  = "hnode"
	 NODE_TYPE_TOWERNODE = "tnode"
	 NODE_TYPE_AMPNODE   = "anode"
	 NODE_TYPE_UNDEFINED = "undef"
 )
 
 func (s *NodeType) Scan(value interface{}) error {
	 *s = NodeType(value.(string))
 
	 return nil
 }
 
 func (s NodeType) Value() (driver.Value, error) {
	 return string(s), nil
 }
 
 func (s NodeType) String() string {
	 return string(s)
 }
 