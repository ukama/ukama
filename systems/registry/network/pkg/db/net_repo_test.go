/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

package db_test

import (
	extsql "database/sql"
	"fmt"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/tj/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	net_db "github.com/ukama/ukama/systems/registry/network/pkg/db"
)

type UkamaDbMock struct {
	GormDb *gorm.DB
}

func (u UkamaDbMock) Init(model ...interface{}) error {
	panic("implement me: Init()")
}

func (u UkamaDbMock) Connect() error {
	panic("implement me: Connect()")
}

func (u UkamaDbMock) GetGormDb() *gorm.DB {
	return u.GormDb
}

func (u UkamaDbMock) InitDB() error {
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func() error) error {
	log.Fatal("implement me: ExecuteInTransaction()")
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func(tx *gorm.DB) error) error {
	log.Fatal("implement me: ExecuteInTransaction2()")
	return nil
}

// Test constants
const (
	testNetworkName1   = "network1"
	testNetworkName3   = "test"
	testCountryUSA     = "USA"
	testNetworkVerizon = "Verizon"
)

// Test data builders
type NetworkBuilder struct {
	network *net_db.Network
}

func NewNetworkBuilder() *NetworkBuilder {
	return &NetworkBuilder{
		network: &net_db.Network{
			Id:               uuid.NewV4(),
			Name:             testNetworkName1,
			Deactivated:      false,
			AllowedCountries: pq.StringArray{testCountryUSA},
			AllowedNetworks:  pq.StringArray{testNetworkVerizon},
			Budget:           1000.0,
			Overdraft:        100.0,
			TrafficPolicy:    1,
			PaymentLinks:     true,
			IsDefault:        false,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
			DeletedAt:        gorm.DeletedAt{},
			SyncStatus:       ukama.StatusTypePending,
		},
	}
}

func (b *NetworkBuilder) WithID(id uuid.UUID) *NetworkBuilder {
	b.network.Id = id
	return b
}

func (b *NetworkBuilder) WithName(name string) *NetworkBuilder {
	b.network.Name = name
	return b
}

func (b *NetworkBuilder) WithAllowedCountries(countries []string) *NetworkBuilder {
	b.network.AllowedCountries = pq.StringArray(countries)
	return b
}

func (b *NetworkBuilder) WithAllowedNetworks(networks []string) *NetworkBuilder {
	b.network.AllowedNetworks = pq.StringArray(networks)
	return b
}

func (b *NetworkBuilder) Build() *net_db.Network {
	return b.network
}

// Database setup helpers
func setupMockDB(t *testing.T) (sqlmock.Sqlmock, *gorm.DB, net_db.NetRepo) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gdb, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	repo := net_db.NewNetRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	return mock, gdb, repo
}

func setupMockDBWithExpectations(t *testing.T, expectations func(sqlmock.Sqlmock)) (sqlmock.Sqlmock, *gorm.DB, net_db.NetRepo) {
	mock, gdb, repo := setupMockDB(t)
	expectations(mock)
	return mock, gdb, repo
}

func Test_NetRepo_Get(t *testing.T) {
	t.Run("NetworkExist", func(t *testing.T) {
		// Arrange
		netID := uuid.NewV4()
		expectedNetwork := NewNetworkBuilder().
			WithID(netID).
			WithName(testNetworkName1).
			WithAllowedNetworks([]string{testNetworkVerizon}).
			WithAllowedCountries([]string{testCountryUSA}).
			Build()

		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"id", "name", "allowed_networks", "allowed_countries"}).
				AddRow(netID, testNetworkName1, expectedNetwork.AllowedNetworks, expectedNetwork.AllowedCountries)

			mock.ExpectQuery(`^SELECT.*networks.*`).
				WithArgs(netID, sqlmock.AnyArg()).
				WillReturnRows(rows)
		})

		// Act
		net, err := repo.Get(netID)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, net)
		assert.Equal(t, net.Id, netID)
		assert.Equal(t, net.AllowedNetworks, expectedNetwork.AllowedNetworks)
		assert.Equal(t, net.AllowedCountries, expectedNetwork.AllowedCountries)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NetworkNotFound", func(t *testing.T) {
		// Arrange
		netID := uuid.NewV4()

		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery(`^SELECT.*networks.*`).
				WithArgs(netID, sqlmock.AnyArg()).
				WillReturnError(extsql.ErrNoRows)
		})

		// Act
		net, err := repo.Get(netID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, net)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_NetRepo_GetByName(t *testing.T) {
	t.Run("NetworkExist", func(t *testing.T) {
		// Arrange
		netID := uuid.NewV4()

		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"id", "name"}).
				AddRow(netID, testNetworkName1)

			mock.ExpectQuery(`^SELECT.*networks.*`).
				WithArgs(testNetworkName1, sqlmock.AnyArg()).
				WillReturnRows(rows)
		})

		// Act
		network, err := repo.GetByName(testNetworkName1)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, network)
		assert.Equal(t, network.Id, netID)
		assert.Equal(t, network.Name, testNetworkName1)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NetworkNotFound", func(t *testing.T) {
		// Arrange
		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery(`^SELECT.*network.*`).
				WithArgs(testNetworkName1, sqlmock.AnyArg()).
				WillReturnError(extsql.ErrNoRows)
		})

		// Act
		network, err := repo.GetByName(testNetworkName1)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, network)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_NetRepo_GetAll(t *testing.T) {
	t.Run("networks exist", func(t *testing.T) {
		// Arrange
		netID := uuid.NewV4()

		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"id", "name"}).
				AddRow(netID, testNetworkName1)

			mock.ExpectQuery(`^SELECT.*networks.*`).
				WithArgs().
				WillReturnRows(rows)
		})

		// Act
		networks, err := repo.GetAll()

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, networks)
		assert.Len(t, networks, 1)
		assert.Equal(t, networks[0].Id, netID)
		assert.Equal(t, networks[0].Name, testNetworkName1)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NoNetworksFound", func(t *testing.T) {
		// Arrange
		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery(`^SELECT.*networks.*`).
				WithArgs().
				WillReturnError(extsql.ErrNoRows)
		})

		// Act
		networks, err := repo.GetAll()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, networks)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_NetRepo_GetDefault(t *testing.T) {
	t.Run("DefaultNetworkExists", func(t *testing.T) {
		// Arrange
		netID := uuid.NewV4()

		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"id", "name", "is_default"}).
				AddRow(netID, testNetworkName1, true)

			mock.ExpectQuery(`^SELECT.*networks.*`).
				WithArgs(true, sqlmock.AnyArg()).
				WillReturnRows(rows)
		})

		// Act
		network, err := repo.GetDefault()

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, network)
		assert.Equal(t, network.Id, netID)
		assert.Equal(t, network.Name, testNetworkName1)
		assert.True(t, network.IsDefault)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NoDefaultNetwork", func(t *testing.T) {
		// Arrange
		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery(`^SELECT.*networks.*`).
				WithArgs(true, sqlmock.AnyArg()).
				WillReturnError(extsql.ErrNoRows)
		})

		// Act
		network, err := repo.GetDefault()

		// Assert
		assert.Error(t, err)
		assert.Nil(t, network)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_NetRepo_GetNetworkCount(t *testing.T) {
	t.Run("GetNetworkCount", func(t *testing.T) {
		// Arrange
		expectedCount := int64(5)

		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			rows := sqlmock.NewRows([]string{"count"}).
				AddRow(expectedCount)

			mock.ExpectQuery(`^SELECT count\(\*\) FROM "networks"`).
				WithArgs().
				WillReturnRows(rows)
		})

		// Act
		count, err := repo.GetNetworkCount()

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCount, count)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("GetNetworkCountError", func(t *testing.T) {
		// Arrange
		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectQuery(`^SELECT count\(\*\) FROM "networks"`).
				WithArgs().
				WillReturnError(fmt.Errorf("database error"))
		})

		// Act
		count, err := repo.GetNetworkCount()

		// Assert
		assert.Error(t, err)
		assert.Equal(t, int64(0), count)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_NetRepo_Add(t *testing.T) {
	t.Run("AddNetwork", func(t *testing.T) {
		// Arrange
		network := NewNetworkBuilder().
			WithName(testNetworkName1).
			Build()

		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()

			mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
				WithArgs(network.Id, network.Name, sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectCommit()
		})

		// Act
		err := repo.Add(network, nil)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("AddNetworkWithNestedFunc", func(t *testing.T) {
		// Arrange
		network := NewNetworkBuilder().
			WithName(testNetworkName1).
			Build()

		nestedFunc := func(network *net_db.Network, tx *gorm.DB) error {
			// Simulate some nested operation
			return nil
		}

		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()

			mock.ExpectExec(regexp.QuoteMeta(`INSERT`)).
				WithArgs(network.Id, network.Name, sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
					sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectCommit()
		})

		// Act
		err := repo.Add(network, nestedFunc)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("AddNetworkWithNestedFuncError", func(t *testing.T) {
		// Arrange
		network := NewNetworkBuilder().
			WithName(testNetworkName1).
			Build()

		nestedFunc := func(network *net_db.Network, tx *gorm.DB) error {
			return fmt.Errorf("nested function error")
		}

		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()
			mock.ExpectRollback()
		})

		// Act
		err := repo.Add(network, nestedFunc)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nested function error")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_NetRepo_Delete(t *testing.T) {
	t.Run("NetworkExist", func(t *testing.T) {
		// Arrange
		network := NewNetworkBuilder().
			WithName(testNetworkName3).
			Build()

		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(`UPDATE "networks"`)).
				WithArgs(sqlmock.AnyArg(), network.Id).
				WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
		})

		// Act
		err := repo.Delete(network.Id)

		// Assert
		assert.NoError(t, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("NetworkNotFound", func(t *testing.T) {
		// Arrange
		networkID := uuid.NewV4()

		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(`UPDATE "networks"`)).
				WithArgs(sqlmock.AnyArg(), networkID).
				WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected
			mock.ExpectCommit()
		})

		// Act
		err := repo.Delete(networkID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, gorm.ErrRecordNotFound, err)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("DeleteDatabaseError", func(t *testing.T) {
		// Arrange
		networkID := uuid.NewV4()

		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()
			mock.ExpectExec(regexp.QuoteMeta(`UPDATE "networks"`)).
				WithArgs(sqlmock.AnyArg(), networkID).
				WillReturnError(fmt.Errorf("database connection lost"))
			mock.ExpectRollback()
		})

		// Act
		err := repo.Delete(networkID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database connection lost")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}

func Test_NetRepo_SetDefault(t *testing.T) {
	t.Run("SetDefaultSuccess", func(t *testing.T) {
		// Arrange
		netID := uuid.NewV4()

		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			mock.ExpectBegin()

			rows := sqlmock.NewRows([]string{"id", "name", "is_default"}).
				AddRow(netID, testNetworkName1, false)
			mock.ExpectQuery(`^SELECT.*networks.*`).
				WithArgs(netID, sqlmock.AnyArg()).
				WillReturnRows(rows)

			mock.ExpectExec(`UPDATE "networks" SET "is_default"=\$1,"updated_at"=\$2 WHERE is_default = \$3 AND "networks"\."deleted_at" IS NULL`).
				WithArgs(false, sqlmock.AnyArg(), true).
				WillReturnResult(sqlmock.NewResult(0, 2))

			mock.ExpectExec(`UPDATE "networks" SET "is_default"=\$1,"updated_at"=\$2 WHERE "networks"\."deleted_at" IS NULL AND "id" = \$3`).
				WithArgs(true, sqlmock.AnyArg(), netID).
				WillReturnResult(sqlmock.NewResult(1, 1))

			mock.ExpectCommit()
		})

		// Act
		network, err := repo.SetDefault(netID, true)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, network)
		assert.Equal(t, network.Id, netID)
		assert.True(t, network.IsDefault)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("SetDefaultNetworkNotFound", func(t *testing.T) {
		// Arrange
		netID := uuid.NewV4()

		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			// Expect transaction begin
			mock.ExpectBegin()

			// Expect find network by ID to fail (new implementation order)
			mock.ExpectQuery(`^SELECT.*networks.*`).
				WithArgs(netID, sqlmock.AnyArg()).
				WillReturnError(extsql.ErrNoRows)

			// Expect transaction rollback
			mock.ExpectRollback()
		})

		// Act
		network, err := repo.SetDefault(netID, true)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, network)
		assert.Contains(t, err.Error(), "failed to find network")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("SetDefaultAlreadyDefault", func(t *testing.T) {
		// Arrange
		netID := uuid.NewV4()

		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			// Expect transaction begin
			mock.ExpectBegin()

			// Expect find network by ID - already default
			rows := sqlmock.NewRows([]string{"id", "name", "is_default"}).
				AddRow(netID, testNetworkName1, true)
			mock.ExpectQuery(`^SELECT.*networks.*`).
				WithArgs(netID, sqlmock.AnyArg()).
				WillReturnRows(rows)

			// No additional updates expected due to early return optimization

			// Expect transaction commit
			mock.ExpectCommit()
		})

		// Act
		network, err := repo.SetDefault(netID, true)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, network)
		assert.Equal(t, network.Id, netID)
		assert.True(t, network.IsDefault)

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("SetDefaultClearExistingDefaultsError", func(t *testing.T) {
		// Arrange
		netID := uuid.NewV4()

		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			// Expect transaction begin
			mock.ExpectBegin()

			// Expect find network by ID - not default
			rows := sqlmock.NewRows([]string{"id", "name", "is_default"}).
				AddRow(netID, testNetworkName1, false)
			mock.ExpectQuery(`^SELECT.*networks.*`).
				WithArgs(netID, sqlmock.AnyArg()).
				WillReturnRows(rows)

			// Expect clear existing defaults to fail
			mock.ExpectExec(`UPDATE "networks" SET "is_default"=\$1,"updated_at"=\$2 WHERE is_default = \$3 AND "networks"\."deleted_at" IS NULL`).
				WithArgs(false, sqlmock.AnyArg(), true).
				WillReturnError(fmt.Errorf("failed to clear existing defaults"))

			// Expect transaction rollback
			mock.ExpectRollback()
		})

		// Act
		network, err := repo.SetDefault(netID, true)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, network)
		assert.Contains(t, err.Error(), "failed to clear existing default networks")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("SetDefaultUpdateNetworkError", func(t *testing.T) {
		// Arrange
		netID := uuid.NewV4()

		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			// Expect transaction begin
			mock.ExpectBegin()

			// Expect find network by ID - not default
			rows := sqlmock.NewRows([]string{"id", "name", "is_default"}).
				AddRow(netID, testNetworkName1, false)
			mock.ExpectQuery(`^SELECT.*networks.*`).
				WithArgs(netID, sqlmock.AnyArg()).
				WillReturnRows(rows)

			// Expect clear existing defaults to succeed
			mock.ExpectExec(`UPDATE "networks" SET "is_default"=\$1,"updated_at"=\$2 WHERE is_default = \$3 AND "networks"\."deleted_at" IS NULL`).
				WithArgs(false, sqlmock.AnyArg(), true).
				WillReturnResult(sqlmock.NewResult(0, 1))

			// Expect update network to fail
			mock.ExpectExec(`UPDATE "networks" SET "is_default"=\$1,"updated_at"=\$2 WHERE "networks"\."deleted_at" IS NULL AND "id" = \$3`).
				WithArgs(true, sqlmock.AnyArg(), netID).
				WillReturnError(fmt.Errorf("failed to update network"))

			// Expect transaction rollback
			mock.ExpectRollback()
		})

		// Act
		network, err := repo.SetDefault(netID, true)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, network)
		assert.Contains(t, err.Error(), "failed to update network")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})

	t.Run("SetDefaultTransactionBeginError", func(t *testing.T) {
		// Arrange
		netID := uuid.NewV4()

		mock, _, repo := setupMockDBWithExpectations(t, func(mock sqlmock.Sqlmock) {
			// Expect transaction begin to fail
			mock.ExpectBegin().WillReturnError(fmt.Errorf("transaction begin failed"))
		})

		// Act
		network, err := repo.SetDefault(netID, true)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, network)
		assert.Contains(t, err.Error(), "transaction begin failed")

		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	})
}
