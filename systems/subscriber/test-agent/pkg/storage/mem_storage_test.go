package storage_test

import (
	"testing"

	"github.com/tj/assert"
	"github.com/ukama/ukama/systems/subscriber/test-agent/pkg/storage"
)

func TestMemStorage_Get(t *testing.T) {
	t.Run("IccidFound", func(t *testing.T) {
		t.Parallel()

		const (
			iccid  = "b8f04217beabf6a19e7eb5b3"
			imsi   = "eabf6a19e7eb5b3"
			status = "inactive"
		)

		sim := &storage.SimInfo{
			Iccid:  iccid,
			Imsi:   imsi,
			Status: status,
		}

		data := map[string]*storage.SimInfo{
			iccid: sim,
		}

		s := storage.NewMemStorage(data)

		res, err := s.Get(iccid)

		assert.NoError(t, err)
		assert.Equal(t, res, sim)
	})

	t.Run("IccidNotFound", func(t *testing.T) {
		t.Parallel()

		const iccid = "b8f04217beabf6a19e7eb5b3"
		s := storage.NewMemStorage(make(map[string]*storage.SimInfo))

		res, err := s.Get(iccid)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestMemStorage_Put(t *testing.T) {
	t.Run("PutIccid", func(t *testing.T) {
		t.Parallel()

		const (
			iccid  = "b8f04217beabf6a19e7eb5b3"
			imsi   = "eabf6a19e7eb5b3"
			status = "inactive"
		)

		sim := &storage.SimInfo{
			Iccid:  iccid,
			Imsi:   imsi,
			Status: status,
		}

		s := storage.NewMemStorage(make(map[string]*storage.SimInfo))

		err := s.Put(iccid, sim)

		assert.NoError(t, err)
	})
}

func TestMemStorage_Delete(t *testing.T) {
	t.Run("IccidFound", func(t *testing.T) {
		t.Parallel()

		const (
			iccid  = "b8f04217beabf6a19e7eb5b3"
			imsi   = "eabf6a19e7eb5b3"
			status = "inactive"
		)

		sim := &storage.SimInfo{
			Iccid:  iccid,
			Imsi:   imsi,
			Status: status,
		}

		data := map[string]*storage.SimInfo{
			iccid: sim,
		}

		s := storage.NewMemStorage(data)

		err := s.Delete(iccid)

		assert.NoError(t, err)
	})

	t.Run("IccidNotFound", func(t *testing.T) {
		t.Parallel()

		const iccid = "b8f04217beabf6a19e7eb5b3"
		s := storage.NewMemStorage(make(map[string]*storage.SimInfo))

		err := s.Delete(iccid)

		assert.NoError(t, err)
	})
}
