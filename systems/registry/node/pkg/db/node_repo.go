package db

import (
	"fmt"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const MaxAttachedNodes = 2

type NodeRepo interface {
	Add(*Node, func(*Node, *gorm.DB) error) error
	Get(ukama.NodeID) (*Node, error)
	GetForOrg(uuid.UUID) ([]Node, error)
	GetAll() ([]Node, error)
	Delete(ukama.NodeID, func(ukama.NodeID, *gorm.DB) error) error
	Update(*Node, func(*Node, *gorm.DB) error) error
	AttachNodes(nodeId ukama.NodeID, attachedNodeId []string) error
	DetachNode(detachNodeId ukama.NodeID) error
	GetNodeCount() (int64, int64, int64, error)
}

type nodeRepo struct {
	Db sql.Db
}

func NewNodeRepo(db sql.Db) NodeRepo {
	return &nodeRepo{
		Db: db,
	}
}

func (n *nodeRepo) Add(node *Node, nestedFunc func(node *Node, tx *gorm.DB) error) error {
	err := n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		if nestedFunc != nil {
			nestErr := nestedFunc(node, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		result := tx.Create(node)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	return err
}

func (n *nodeRepo) Get(id ukama.NodeID) (*Node, error) {
	var node Node

	result := n.Db.GetGormDb().Preload(clause.Associations).First(&node, "id=?", id.StringLowercase())

	if result.Error != nil {
		return nil, result.Error
	}

	return &node, nil
}

func (n *nodeRepo) GetForOrg(orgId uuid.UUID) ([]Node, error) {
	var nodes []Node

	result := n.Db.GetGormDb().Preload(clause.Associations).Where("org_id = ?", orgId.String()).Find(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	return nodes, nil
}

func (n *nodeRepo) GetAll() ([]Node, error) {
	var nodes []Node

	result := n.Db.GetGormDb().Preload(clause.Associations).Find(&nodes)

	if result.Error != nil {
		return nil, result.Error
	}

	return nodes, nil
}

// TODO: check for still allocated and attached nodes
func (n *nodeRepo) Delete(nodeId ukama.NodeID, nestedFunc func(ukama.NodeID, *gorm.DB) error) error {
	err := n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Select(clause.Associations, "NodeStatus").Delete(&Node{Id: nodeId.StringLowercase()})
		if result.Error != nil {
			return result.Error
		}

		if nestedFunc != nil {
			nestErr := nestedFunc(nodeId, tx)
			if nestErr != nil {
				return nestErr
			}
		}

		return nil
	})

	return err
}

// Update updated node with `id`. Only fields that are not nil are updated, eg name and state.
func (n *nodeRepo) Update(node *Node, nestedFunc func(*Node, *gorm.DB) error) error {
	err := n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Clauses(clause.Returning{}).Updates(node)

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

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

func (n *nodeRepo) AttachNodes(nodeId ukama.NodeID, attachedNodeId []string) error {
	batchGet := func(nodeIds []string) ([]Site, error) {
		var nodes []Site

		result := n.Db.GetGormDb().Where("id IN ?", nodeIds).Find(&nodes)

		if result.Error != nil {
			return nil, result.Error
		}

		return nodes, nil
	}

	attachedNodeSites, err := batchGet(attachedNodeId)
	if err != nil {
		return err
	}

	res, err := batchGet([]string{nodeId.StringLowercase()})
	if err != nil {
		return err
	}

	parentNodeSite := res[0]

	parentNode, err := n.Get(nodeId)
	if err != nil {
		return err
	}

	if parentNode.Type != ukama.NODE_ID_TYPE_TOWERNODE {
		return status.Errorf(codes.InvalidArgument, "node type must be a towernode")
	}

	if parentNode.Attached == nil {
		parentNode.Attached = make([]*Node, 0)
	}

	if len(attachedNodeId)+len(parentNode.Attached) > MaxAttachedNodes {
		return status.Errorf(codes.InvalidArgument, "max number of attached nodes is %d", MaxAttachedNodes)
	}

	err = n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		for _, aNds := range attachedNodeSites {
			aNd, err := ukama.ValidateNodeId(aNds.NodeId)
			if err != nil {
				return err
			}

			an, err := n.Get(aNd)
			if err != nil {
				return err
			}
			if an.Type != ukama.NODE_ID_TYPE_AMPNODE {
				return status.Errorf(codes.InvalidArgument, "cannot attach non amplifier node")
			}

			if parentNodeSite.SiteId != aNds.SiteId {
				return status.Errorf(codes.InvalidArgument, "cannot attach nodes from different sites")
			}

			parentNode.Attached = append(parentNode.Attached, an)
		}

		d := tx.Save(parentNode)
		if d.Error != nil {
			return d.Error
		}

		return nil
	})

	return err
}

func (n *nodeRepo) DetachNode(detachNodeId ukama.NodeID) error {
	err := n.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		result := tx.Exec("delete from attached_nodes where attached_id=(select id from nodes where node_id=?) OR node_id=(select id from nodes where node_id=?)",
			detachNodeId, detachNodeId)

		if result.Error != nil {
			return fmt.Errorf("failed to remove from group for %s node: error %s", detachNodeId.StringLowercase(), result.Error)
		}

		return nil
	})

	return err
}

func (r *nodeRepo) GetNodeCount() (nodeCount, onlineCount, offlineCount int64, err error) {
	db := r.Db.GetGormDb()

	if err := db.Model(&Node{}).Count(&nodeCount).Error; err != nil {
		return 0, 0, 0, err
	}

	res1 := db.Raw("select COUNT(*) from nodes LEFT JOIN node_statuses ON nodes.id = node_statuses.node_id WHERE node_statuses.conn = ?", Online).Scan(&onlineCount)
	if res1.Error != nil {
		return 0, 0, 0, err
	}

	res2 := db.Raw("select COUNT(*) from nodes LEFT JOIN node_statuses ON nodes.id = node_statuses.node_id WHERE node_statuses.conn = ?", Offline).Scan(&offlineCount)
	if res2.Error != nil {
		return 0, 0, 0, err
	}

	return nodeCount, onlineCount, offlineCount, nil
}
