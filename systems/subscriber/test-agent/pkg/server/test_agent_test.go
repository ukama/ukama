package server_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/subscriber/test-agent/mocks"
	pb "github.com/ukama/ukama/systems/subscriber/test-agent/pb/gen"
	"github.com/ukama/ukama/systems/subscriber/test-agent/pkg/server"
	"github.com/ukama/ukama/systems/subscriber/test-agent/pkg/storage"
)

func TestTestAgentServer_GetSim(t *testing.T) {
	t.Run("SimFound", func(t *testing.T) {
		t.Parallel()

		const (
			iccid  = "b8f04217beabf6a19e7eb5b3"
			imsi   = "eabf6a19e7eb5b3"
			status = "inactive"
		)

		store := &mocks.Storage{}

		store.On("Get", iccid).Return(
			&storage.SimInfo{
				Iccid:  iccid,
				Imsi:   imsi,
				Status: status,
			}, nil).Once()

		s := server.NewTestAgentServer(store)
		resp, err := s.GetSim(context.TODO(), &pb.GetSimRequest{
			Iccid: iccid})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, iccid, resp.GetSimInfo().Iccid)
		assert.Equal(t, imsi, resp.GetSimInfo().Imsi)
		assert.Equal(t, status, resp.GetSimInfo().Status)
		store.AssertExpectations(t)
	})

	t.Run("SimUnknownErrorOnGet", func(t *testing.T) {
		t.Parallel()

		const iccid = "b8f04217beabf6a19e7eb5b3"

		store := &mocks.Storage{}

		store.On("Get", iccid).Return(nil, storage.ErrInternal).Once()

		s := server.NewTestAgentServer(store)
		resp, err := s.GetSim(context.TODO(), &pb.GetSimRequest{
			Iccid: iccid})

		assert.Error(t, err)
		assert.Nil(t, resp)
		store.AssertExpectations(t)
	})

	t.Run("SimNotFoundAndNoErrorOnCreate", func(t *testing.T) {
		t.Parallel()

		const (
			iccid  = "b8f04217beabf6a19e7eb5b3"
			imsi   = "eabf6a19e7eb5b3"
			status = "inactive"
		)

		store := &mocks.Storage{}

		sim := &storage.SimInfo{
			Iccid:  iccid,
			Imsi:   imsi,
			Status: status,
		}

		store.On("Get", iccid).Return(nil, storage.ErrNotFound).Once()
		store.On("Put", iccid, sim).Return(nil).Once()

		s := server.NewTestAgentServer(store)
		resp, err := s.GetSim(context.TODO(), &pb.GetSimRequest{
			Iccid: iccid})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, iccid, resp.GetSimInfo().Iccid)
		assert.Equal(t, imsi, resp.GetSimInfo().Imsi)
		assert.Equal(t, status, resp.GetSimInfo().Status)
		store.AssertExpectations(t)
	})

	t.Run("SimNotFoundAndErrorOnCreate", func(t *testing.T) {
		t.Parallel()

		const (
			iccid  = "b8f04217beabf6a19e7eb5b3"
			imsi   = "eabf6a19e7eb5b3"
			status = "inactive"
		)

		store := &mocks.Storage{}

		sim := &storage.SimInfo{
			Iccid:  iccid,
			Imsi:   imsi,
			Status: status,
		}

		store.On("Get", iccid).Return(nil, storage.ErrNotFound).Once()
		store.On("Put", iccid, sim).Return(storage.ErrInternal).Once()

		s := server.NewTestAgentServer(store)
		resp, err := s.GetSim(context.TODO(), &pb.GetSimRequest{
			Iccid: iccid})

		assert.Error(t, err)
		assert.Nil(t, resp)
		store.AssertExpectations(t)
	})
}

func TestTestAgentServer_ActivateSim(t *testing.T) {
	t.Run("SimFoundAndSimStatusInactive", func(t *testing.T) {
		t.Parallel()

		const (
			iccid  = "922a4f72922a775acd978a75"
			status = "inactive"
		)

		store := &mocks.Storage{}

		store.On("Get", iccid).Return(
			&storage.SimInfo{
				Iccid:  iccid,
				Status: status,
			}, nil).Once()

		store.On("Put", iccid, mock.Anything).Return(nil).Once()

		s := server.NewTestAgentServer(store)
		resp, err := s.ActivateSim(context.TODO(), &pb.ActivateSimRequest{
			Iccid: iccid})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		store.AssertExpectations(t)
	})

	t.Run("SimStatusInactiveAndFailToUpdateStatus", func(t *testing.T) {
		t.Parallel()

		const (
			iccid  = "922a4f72922a775acd978a75"
			status = "inactive"
		)

		store := &mocks.Storage{}

		store.On("Get", iccid).Return(
			&storage.SimInfo{
				Iccid:  iccid,
				Status: status,
			}, nil).Once()

		store.On("Put", iccid, mock.Anything).Return(storage.ErrInternal).Once()

		s := server.NewTestAgentServer(store)
		resp, err := s.ActivateSim(context.TODO(), &pb.ActivateSimRequest{
			Iccid: iccid})

		assert.Error(t, err)
		assert.Nil(t, resp)
		store.AssertExpectations(t)
	})

	t.Run("SimFoundAndSimStatusActive", func(t *testing.T) {
		t.Parallel()

		const (
			iccid  = "922a4f72922a775acd978a75"
			status = "active"
		)

		store := &mocks.Storage{}

		store.On("Get", iccid).Return(
			&storage.SimInfo{
				Iccid:  iccid,
				Status: status,
			}, nil).Once()

		s := server.NewTestAgentServer(store)
		resp, err := s.ActivateSim(context.TODO(), &pb.ActivateSimRequest{
			Iccid: iccid})

		assert.Error(t, err)
		assert.Nil(t, resp)
		store.AssertExpectations(t)
	})

	t.Run("SimNotFound", func(t *testing.T) {
		t.Parallel()

		const (
			iccid = "922a4f72922a775acd978a75"
		)

		store := &mocks.Storage{}

		store.On("Get", iccid).Return(nil, storage.ErrNotFound).Once()

		s := server.NewTestAgentServer(store)
		resp, err := s.ActivateSim(context.TODO(), &pb.ActivateSimRequest{
			Iccid: iccid})

		assert.Error(t, err)
		assert.Nil(t, resp)
		store.AssertExpectations(t)
	})
}

func TestTestAgentServer_DeactivateSim(t *testing.T) {
	t.Run("SimFoundAndSimStatusActive", func(t *testing.T) {
		t.Parallel()

		const (
			iccid  = "50f54f9082fbe245b91924d4"
			status = "active"
		)

		store := &mocks.Storage{}

		store.On("Get", iccid).Return(
			&storage.SimInfo{
				Iccid:  iccid,
				Status: status,
			}, nil).Once()

		store.On("Put", iccid, mock.Anything).Return(nil).Once()

		s := server.NewTestAgentServer(store)
		resp, err := s.DeactivateSim(context.TODO(), &pb.DeactivateSimRequest{
			Iccid: iccid})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		store.AssertExpectations(t)
	})

	t.Run("SimStatusActiveAndFailToUpdateStatus", func(t *testing.T) {
		t.Parallel()

		const (
			iccid  = "50f54f9082fbe245b91924d4"
			status = "active"
		)

		store := &mocks.Storage{}

		store.On("Get", iccid).Return(
			&storage.SimInfo{
				Iccid:  iccid,
				Status: status,
			}, nil).Once()

		store.On("Put", iccid, mock.Anything).Return(storage.ErrInternal).Once()

		s := server.NewTestAgentServer(store)
		resp, err := s.DeactivateSim(context.TODO(), &pb.DeactivateSimRequest{
			Iccid: iccid})

		assert.Error(t, err)
		assert.Nil(t, resp)
		store.AssertExpectations(t)
	})

	t.Run("SimFoundAndSimStatusInactive", func(t *testing.T) {
		t.Parallel()

		const (
			iccid  = "50f54f9082fbe245b91924d4"
			status = "inactive"
		)

		store := &mocks.Storage{}

		store.On("Get", iccid).Return(
			&storage.SimInfo{
				Iccid:  iccid,
				Status: status,
			}, nil).Once()

		s := server.NewTestAgentServer(store)
		resp, err := s.DeactivateSim(context.TODO(), &pb.DeactivateSimRequest{
			Iccid: iccid})

		assert.Error(t, err)
		assert.Nil(t, resp)
		store.AssertExpectations(t)
	})

	t.Run("SimNotFound", func(t *testing.T) {
		t.Parallel()

		const (
			iccid = "50f54f9082fbe245b91924d4"
		)

		store := &mocks.Storage{}

		store.On("Get", iccid).Return(nil, storage.ErrNotFound).Once()

		s := server.NewTestAgentServer(store)
		resp, err := s.DeactivateSim(context.TODO(), &pb.DeactivateSimRequest{
			Iccid: iccid})

		assert.Error(t, err)
		assert.Nil(t, resp)
		store.AssertExpectations(t)
	})
}

func TestTestAgentServer_TerminateSim(t *testing.T) {
	t.Run("SimFoundAndSimStatusInactive", func(t *testing.T) {
		t.Parallel()

		const (
			iccid  = "e828484bac9ca46995e5617e"
			status = "inactive"
		)

		store := &mocks.Storage{}

		store.On("Get", iccid).Return(
			&storage.SimInfo{
				Iccid:  iccid,
				Status: status,
			}, nil).Once()

		store.On("Delete", iccid, mock.Anything).Return(nil).Once()

		s := server.NewTestAgentServer(store)
		resp, err := s.TerminateSim(context.TODO(), &pb.TerminateSimRequest{
			Iccid: iccid})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		store.AssertExpectations(t)
	})

	t.Run("SimStatusInactiveAndFailToUpdateStatus", func(t *testing.T) {
		t.Parallel()

		const (
			iccid  = "e828484bac9ca46995e5617e"
			status = "inactive"
		)

		store := &mocks.Storage{}

		store.On("Get", iccid).Return(
			&storage.SimInfo{
				Iccid:  iccid,
				Status: status,
			}, nil).Once()

		store.On("Delete", iccid, mock.Anything).Return(storage.ErrInternal).Once()

		s := server.NewTestAgentServer(store)
		resp, err := s.TerminateSim(context.TODO(), &pb.TerminateSimRequest{
			Iccid: iccid})

		assert.Error(t, err)
		assert.Nil(t, resp)
		store.AssertExpectations(t)
	})

	t.Run("SimFoundAndSimStatusActive", func(t *testing.T) {
		t.Parallel()

		const (
			iccid  = "e828484bac9ca46995e5617e"
			status = "active"
		)

		store := &mocks.Storage{}

		store.On("Get", iccid).Return(
			&storage.SimInfo{
				Iccid:  iccid,
				Status: status,
			}, nil).Once()

		s := server.NewTestAgentServer(store)
		resp, err := s.TerminateSim(context.TODO(), &pb.TerminateSimRequest{
			Iccid: iccid})

		assert.Error(t, err)
		assert.Nil(t, resp)
		store.AssertExpectations(t)
	})

	t.Run("SimNotFound", func(t *testing.T) {
		t.Parallel()

		const (
			iccid = "e828484bac9ca46995e5617e"
		)

		store := &mocks.Storage{}

		store.On("Get", iccid).Return(nil, storage.ErrNotFound).Once()

		s := server.NewTestAgentServer(store)
		resp, err := s.TerminateSim(context.TODO(), &pb.TerminateSimRequest{
			Iccid: iccid})

		assert.Error(t, err)
		assert.Nil(t, resp)
		store.AssertExpectations(t)
	})
}
