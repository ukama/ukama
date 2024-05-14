/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
 
package policy_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	cmocks "github.com/ukama/ukama/systems/common/mocks"
	"github.com/ukama/ukama/systems/ukama-agent/asr/mocks"
	"github.com/ukama/ukama/systems/ukama-agent/asr/pkg/db"
	ip "github.com/ukama/ukama/systems/ukama-agent/asr/pkg/policy"
)

func TestPolicy_DataCapCheck(t *testing.T) {
	t.Run("DataCapChecks", func(t *testing.T) {
		valid := ip.DataCapCheck(sub)
		assert.Equal(t, true, valid)

		newSub := sub
		newSub.Policy.ConsumedData = 1024000000
		valid = ip.DataCapCheck(newSub)
		assert.Equal(t, false, valid)

		newSub.Policy.ConsumedData = 2024000000
		valid = ip.DataCapCheck(newSub)
		assert.Equal(t, false, valid)

	})
}

func TestPolicy_AllowedTimeOfService(t *testing.T) {
	t.Run("AllowedTimeOfService", func(t *testing.T) {
		valid := ip.AllowedTimeOfServiceCheck(sub)
		assert.Equal(t, true, valid)

		newSub := sub
		newSub.AllowedTimeOfService = 0
		valid = ip.AllowedTimeOfServiceCheck(newSub)
		assert.Equal(t, false, valid)

	})
}

func TestPolicy_ValidityCheck(t *testing.T) {
	t.Run("ValidityCheck", func(t *testing.T) {
		valid := ip.ValidityCheck(sub)
		assert.Equal(t, true, valid)

		newSub := sub
		newSub.Policy.StartTime = uint64(time.Now().Unix() - 10000000)
		newSub.Policy.EndTime = uint64(time.Now().Unix())
		valid = ip.ValidityCheck(newSub)
		assert.Equal(t, false, valid)

		newSub.Policy.StartTime = uint64(time.Now().Unix() + 1000)
		newSub.Policy.EndTime = uint64(time.Now().Unix() + 100000000)
		valid = ip.ValidityCheck(newSub)
		assert.Equal(t, false, valid)

	})
}

func TestPolicy_RemoveProfile(t *testing.T) {
	asrRepo := &mocks.AsrRecordRepo{}
	mbC := &cmocks.MsgBusServiceClient{}

	t.Run("RemoveProfile", func(t *testing.T) {

		asrRepo.On("Delete", sub.Imsi, db.POLICY_FAILURE).Return(nil).Once()

		mbC.On("PublishRequest", "request.cloud.local.ukama.ukamaagent.asr.nodefeeder.publish", mock.Anything).Return(nil).Once()
		mbC.On("PublishRequest", "event.cloud.local.ukama.ukamaagent.asr.activesubscriber.delete", mock.Anything).Return(nil).Once()

		pc := ip.NewPolicyController(asrRepo, mbC, dataplanHost, OrgName, OrgId, Reroute, MonitoringPeriod, false)
		assert.NotNil(t, pc)

		err, state := ip.RemoveProfile(pc, sub)
		assert.NoError(t, err)
		assert.Equal(t, true, state)

		asrRepo.AssertExpectations(t)
		assert.NoError(t, err)
	})
}
