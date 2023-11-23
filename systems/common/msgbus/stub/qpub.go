/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package stub

type QPubStub struct {
}

func (q QPubStub) Publish(payload any, routingKey string) error {
	return nil
}

func (q QPubStub) PublishToQueue(queueName string, payload any) error {
	return nil
}
func (q QPubStub) Close() error {
	return nil
}
