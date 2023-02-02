package storage

import (
	"testing"

	"github.com/tj/assert"
)

func TestMemStorage_Get(t *testing.T) {
	t.Run("IccidFound", func(t *testing.T) {
		const (
			iccid  = "b8f04217beabf6a19e7eb5b3"
			imsi   = "eabf6a19e7eb5b3"
			status = "inactive"
		)

		sim := &SimInfo{
			Iccid:  iccid,
			Imsi:   imsi,
			Status: status,
		}

		s := NewMemStorage()
		s.data = map[string]*SimInfo{
			iccid: sim,
		}

		res, err := s.Get(iccid)

		assert.NoError(t, err)
		assert.Equal(t, res, sim)
	})

	t.Run("IccidNotFound", func(t *testing.T) {
		const iccid = "b8f04217beabf6a19e7eb5b3"
		s := NewMemStorage()

		res, err := s.Get(iccid)

		assert.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestMemStorage_Put(t *testing.T) {
	t.Run("PutIccid", func(t *testing.T) {
		const (
			iccid  = "b8f04217beabf6a19e7eb5b3"
			imsi   = "eabf6a19e7eb5b3"
			status = "inactive"
		)

		sim := &SimInfo{
			Iccid:  iccid,
			Imsi:   imsi,
			Status: status,
		}

		s := NewMemStorage()

		err := s.Put(iccid, sim)

		assert.NoError(t, err)
	})
}

func TestMemStorage_Delete(t *testing.T) {
	t.Run("IccidFound", func(t *testing.T) {
		const (
			iccid  = "b8f04217beabf6a19e7eb5b3"
			imsi   = "eabf6a19e7eb5b3"
			status = "inactive"
		)

		sim := &SimInfo{
			Iccid:  iccid,
			Imsi:   imsi,
			Status: status,
		}

		s := NewMemStorage()
		s.data = map[string]*SimInfo{
			iccid: sim,
		}

		err := s.Delete(iccid)

		assert.NoError(t, err)
	})

	t.Run("IccidNotFound", func(t *testing.T) {
		const iccid = "b8f04217beabf6a19e7eb5b3"
		s := NewMemStorage()

		err := s.Delete(iccid)

		assert.NoError(t, err)
	})
}
