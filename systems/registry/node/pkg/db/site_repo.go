package db

import (
	"fmt"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"gorm.io/gorm"
)

type SiteRepo interface {
	GetNodes(uuid.UUID) ([]Node, error)
	AddNode(*Site, func(*Site, *gorm.DB) error) error
	RemoveNodes([]string) error
	// RemoveNodeFromNetwork(nodeId ukama.NodeID) error
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
		Where("sites.site_id=?", siteID.String()).Find(&nodes)

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

func (s *siteRepo) RemoveNodes(detachedNodes []string) error {
	var nodes []Site

	err := s.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Where("node_id IN ?", detachedNodes).Find(&nodes)
		if result.Error != nil {
			return result.Error
		}

		if len(nodes) != len(detachedNodes) {
			return fmt.Errorf("invalid arguments: got %d items to match %d rows", len(detachedNodes), len(nodes))
		}

		result = tx.Delete(&nodes)
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

// func (r *nodeRepo) RemoveNodeFromNetwork(nodeId ukama.NodeID) error {
// node, err := r.Get(nodeId)
// if err != nil {
// return err
// }

// if !node.Allocation {
// return status.Errorf(codes.FailedPrecondition, "node is not yet assigned to network")
// }

// res := r.Db.GetGormDb().Exec("select * from attached_nodes where attached_id=(select id from nodes where node_id=?) OR node_id=(select id from nodes where node_id=?)",
// node.Id, node.Id)

// if res.Error != nil {
// return status.Errorf(codes.Internal, "failed to get node grouping result. error %s", res.Error.Error())
// }

// if res.RowsAffected > 0 {
// return status.Errorf(codes.FailedPrecondition, "node is grouped with other nodes.")
// }

// nd := Node{
// Network:    uuid.NullUUID{Valid: false},
// Allocation: false,
// }

// result := r.Db.GetGormDb().Where("node_id=?", node.Id).Select("network", "allocation").Updates(nd)
// if result.Error != nil {
// return fmt.Errorf("failed to remove  node from network id for %s. error %s", nodeId, result.Error)
// }

// return nil
// }
