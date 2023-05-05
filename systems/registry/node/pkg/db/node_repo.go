package db

import (
	"fmt"
	"strings"

	"github.com/ukama/ukama/systems/common/sql"
	"github.com/ukama/ukama/systems/common/ukama"
	"github.com/ukama/ukama/systems/common/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const MaxAttachedNodes = 2

// NodeID must be lowercase
type NodeRepo interface {
	Add(node *Node, nestedFunc ...func() error) error
	Get(id ukama.NodeID) (*Node, error)
	GetAll() (*[]Node, error)
	GetFreeNodes() (*[]Node, error)
	Delete(id ukama.NodeID, nestedFunc ...func() error) error
	Update(id ukama.NodeID, state *NodeState, nodeName *string, nestedFunc ...func() error) error
	AttachNodes(nodeId ukama.NodeID, attachedNodeId []ukama.NodeID, networkID uuid.UUID) error
	DetachNode(detachNodeId ukama.NodeID) error
	GetNodeCount() (int64, int64, int64, error)
	AddNodeToNetwork(nodeId ukama.NodeID, networkID uuid.UUID) error
	RemoveNodeFromNetwork(nodeId ukama.NodeID) error
}

type nodeRepo struct {
	Db sql.Db
}

func NewNodeRepo(db sql.Db) NodeRepo {
	return &nodeRepo{
		Db: db,
	}
}

func (r *nodeRepo) Add(node *Node, nestedFunc ...func() error) error {
	node.NodeID = strings.ToLower(node.NodeID)

	err := r.Db.ExecuteInTransaction(func(tx *gorm.DB) *gorm.DB {
		return tx.Create(node)
	}, nestedFunc...)

	return err
}

func (r *nodeRepo) Get(id ukama.NodeID) (*Node, error) {
	var node Node

	result := r.Db.GetGormDb().Preload(clause.Associations).First(&node, "node_id=?", id.StringLowercase())

	if result.Error != nil {
		return nil, result.Error
	}

	return &node, nil
}

func (r *nodeRepo) GetAll() (*[]Node, error) {
	var node []Node

	result := r.Db.GetGormDb().Preload(clause.Associations).Find(&node)

	if result.Error != nil {
		return nil, result.Error
	}

	return &node, nil
}

func (r *nodeRepo) GetFreeNodes() (*[]Node, error) {
	var node []Node

	result := r.Db.GetGormDb().Preload(clause.Associations).Where("allocation = ?", false).Find(&node)

	if result.Error != nil {
		return nil, result.Error
	}

	return &node, nil
}

func (r *nodeRepo) Delete(id ukama.NodeID, nestedFunc ...func() error) error {
	err := r.Db.ExecuteInTransaction(func(tx *gorm.DB) *gorm.DB {
		d := tx.Delete(&Node{}, "node_id = ?", id.StringLowercase())

		if d.Error != nil {
			return d
		}

		if d.RowsAffected == 0 {
			d.Error = gorm.ErrRecordNotFound

			return d
		}

		return d
	}, nestedFunc...)

	return err
}

// Update updated node with `id`. Only fields that are not nil are updated
func (r *nodeRepo) Update(id ukama.NodeID, state *NodeState, nodeName *string, nestedFunc ...func() error) error {
	var rowsAffected int64
	err := r.Db.ExecuteInTransaction(func(tx *gorm.DB) *gorm.DB {
		nd := Node{}

		if state != nil {
			nd.State = *state
		}

		if nodeName != nil {
			nd.Name = *nodeName
		}

		result := tx.Where("node_id=?", id.StringLowercase()).Updates(nd)
		rowsAffected = result.RowsAffected

		return result
	}, nestedFunc...)

	if rowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return err
}

func (r *nodeRepo) AddNodeToNetwork(nodeId ukama.NodeID, networkID uuid.UUID) error {
	node, err := r.Get(nodeId)
	if err != nil {
		return err
	}

	if node.Allocation {
		return status.Errorf(codes.InvalidArgument, "node is already assigned to network")
	}

	nd := Node{
		Allocation: true,
		Network: uuid.NullUUID{
			UUID:  networkID,
			Valid: true,
		},
	}

	result := r.Db.GetGormDb().Where("node_id=?", node.NodeID).Updates(nd)
	if result.Error != nil {
		return fmt.Errorf("failed to update network id for %s for node: error %s", nodeId, result.Error)
	}

	return nil

}

func (r *nodeRepo) RemoveNodeFromNetwork(nodeId ukama.NodeID) error {
	node, err := r.Get(nodeId)
	if err != nil {
		return err
	}

	if !node.Allocation {
		return status.Errorf(codes.FailedPrecondition, "node is not yet assigned to network")
	}

	res := r.Db.GetGormDb().Exec("select * from attached_nodes where attached_id=(select id from nodes where node_id=?) OR node_id=(select id from nodes where node_id=?)",
		node.NodeID, node.NodeID)

	if res.Error != nil {
		return status.Errorf(codes.Internal, "failed to get node grouping result. error %s", res.Error.Error())
	}

	if res.RowsAffected > 0 {
		return status.Errorf(codes.FailedPrecondition, "node is grouped with other nodes.")
	}

	nd := Node{
		Network:    uuid.NullUUID{Valid: false},
		Allocation: false,
	}

	result := r.Db.GetGormDb().Where("node_id=?", node.NodeID).Select("network", "allocation").Updates(nd)
	if result.Error != nil {
		return fmt.Errorf("failed to remove  node from network id for %s. error %s", nodeId, result.Error)
	}

	return nil
}

func (r *nodeRepo) AttachNodes(nodeId ukama.NodeID, attachedNodeId []ukama.NodeID, networkID uuid.UUID) error {
	parentNode, err := r.Get(nodeId)
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

	err = r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		for _, n := range attachedNodeId {
			an, err := r.Get(n)

			if err != nil {
				return err
			}

			if an.Type != ukama.NODE_ID_TYPE_AMPNODE {
				return status.Errorf(codes.InvalidArgument, "cannot attach non amplifier node")
			}

			parentNode.Attached = append(parentNode.Attached, an)

			nd := Node{
				Allocation: true,
				Network: uuid.NullUUID{
					UUID:  networkID,
					Valid: true,
				},
			}

			result := tx.Where("node_id=?", n).Updates(nd)
			if result.Error != nil {
				return fmt.Errorf("failed to update network id for %s node: error %s", n.StringLowercase(), result.Error)
			}
		}

		parentNode.Network = uuid.NullUUID{
			UUID:  networkID,
			Valid: true,
		}
		parentNode.Allocation = true
		d := tx.Save(parentNode)
		if d.Error != nil {
			return d.Error
		}

		return nil
	})

	return err

}

// DetachNode removes node from parent node
func (r *nodeRepo) DetachNode(detachNodeId ukama.NodeID) error {

	err := r.Db.GetGormDb().Transaction(func(tx *gorm.DB) error {
		nd := Node{
			Network:    uuid.NullUUID{Valid: false},
			Allocation: false,
		}

		result := tx.Where("node_id=?", detachNodeId).Select("network", "allocation").Updates(nd)
		if result.Error != nil {
			return fmt.Errorf("failed to update network id for %s node: error %s", detachNodeId.StringLowercase(), result.Error)
		}

		result = tx.Exec("delete from attached_nodes where attached_id=(select id from nodes where node_id=?)",
			detachNodeId.StringLowercase())

		if result.Error != nil {
			return fmt.Errorf("failed to update network id for %s node: error %s", detachNodeId.StringLowercase(), result.Error)
		}

		return nil
	})

	return err
}

func (r *nodeRepo) GetNodeCount() (nodeCount, activeNodeCount, inactiveNodeCount int64, err error) {
	db := r.Db.GetGormDb()

	if err := db.Model(&Node{}).Count(&nodeCount).Error; err != nil {
		return 0, 0, 0, err
	}

	if err := db.Model(&Node{}).Where("state != ?", Offline).Count(&activeNodeCount).Error; err != nil {
		return 0, 0, 0, err
	}

	if err := db.Model(&Node{}).Where("state = ?", Offline).Count(&inactiveNodeCount).Error; err != nil {
		return 0, 0, 0, err
	}

	return nodeCount, activeNodeCount, inactiveNodeCount, nil
}
