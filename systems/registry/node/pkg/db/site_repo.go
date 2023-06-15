package db

import (
	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SiteRepo interface {
	GetNodes(uuid.UUID) ([]Node, error)
	AddNode(*Site, func(*Site, *gorm.DB) error) error
	RemoveNode(ukama.NodeID) error
	GetFreeNodes() ([]Node, error)
	GetFreeNodesForOrg(uuid.UUID) ([]Node, error)
	IsAllocated(ukama.NodeID) bool
}

type siteRepo struct {
	Db sql.Db
}

func NewSiteRepo(db sql.Db) SiteRepo {
	return &siteRepo{
		Db: db,
	}
}

func (s *siteRepo) GetNodes(siteID uuid.UUID) ([]Node, error) {
	var nodes []Node

	result := s.Db.GetGormDb().Joins("JOIN sites on sites.node_id=nodes.id").
		Preload(clause.Associations).Where("sites.site_id=?", siteID.String()).Find(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return nodes, nil
}

func (s *siteRepo) AddNode(node *Site, nestedFunc func(node *Site, tx *gorm.DB) error) error {
	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Create(node)
		if result.Error != nil {
			return result.Error
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(node, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return err
}

func (s *siteRepo) RemoveNode(nodeId ukama.NodeID) error {
	if !s.IsAllocated(nodeId) {
		return status.Errorf(codes.FailedPrecondition, "node is not yet assigned to site/network")
	}

	res := s.Db.GetGormDb().Exec("select * from attached_nodes where attached_id=(select id from nodes where node_id=?) OR node_id=(select id from nodes where node_id=?)",
		nodeId.StringLowercase(), nodeId.StringLowercase())

	if res.Error != nil {
		return status.Errorf(codes.Internal,
			"failed to get node grouping result. error %s", res.Error.Error())
	}

	if res.RowsAffected > 0 {
		return status.Errorf(codes.FailedPrecondition,
			"node is grouped with other nodes.")
	}

	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&Site{NodeId: nodeId.StringLowercase()})
		if result.Error != nil {
			return result.Error
		}
		return nil
	})

	return err
}

func (s *siteRepo) GetFreeNodes() ([]Node, error) {
	var nodes []Node

	result := s.Db.GetGormDb().Raw("SELECT * from nodes WHERE id NOT IN ? AND deleted_at IS NULL",
		s.Db.GetGormDb().Raw("SELECT node_id from sites WHERE deleted_at IS NULL")).Scan(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	return nodes, nil
}

func (s *siteRepo) GetFreeNodesForOrg(orgId uuid.UUID) ([]Node, error) {
	var nodes []Node

	result := s.Db.GetGormDb().Raw("SELECT * from nodes WHERE id NOT IN ? AND org_Id= ? AND  deleted_at IS NULL",
		s.Db.GetGormDb().Raw("SELECT node_id from sites WHERE deleted_at IS NULL"), orgId).Scan(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	return nodes, nil
}

func (s *siteRepo) IsAllocated(nodeId ukama.NodeID) bool {
	var nd Site

	result := s.Db.GetGormDb().Where(&Site{NodeId: nodeId.StringLowercase()}).First(&nd)
	return result.Error == nil
}
