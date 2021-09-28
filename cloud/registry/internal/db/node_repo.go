package db

import (
	sql2 "database/sql"
	uuid "github.com/satori/go.uuid"
	"github.com/ukama/ukamaX/common/sql"
	"github.com/ukama/ukamaX/common/ukama"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

// NodeID must be lowercase
type NodeRepo interface {
	Add(node *Node) error
	Get(id ukama.NodeID) (*Node, error)
	GetByOrg(orgName string, ownerId uuid.UUID) ([]Node, error)
	Delete(id ukama.NodeID) error
	Update(id ukama.NodeID, state NodeState) error
	GetByUser(ownerId uuid.UUID) ([]Node, error)
}

type nodeRepo struct {
	Db sql.Db
}

func NewNodeRepo(db sql.Db) NodeRepo {
	return &nodeRepo{
		Db: db,
	}
}

func (r *nodeRepo) Add(node *Node) error {
	node.NodeID = strings.ToLower(node.NodeID)
	d := r.Db.GetGormDb().Create(node)
	return d.Error
}

func (r *nodeRepo) Get(id ukama.NodeID) (*Node, error) {
	var node Node
	result := r.Db.GetGormDb().Preload(clause.Associations).First(&node, "node_id=?", id.StringLowercase())
	if result.Error != nil {
		return nil, result.Error
	}
	return &node, nil
}

func (r *nodeRepo) GetByOrg(orgName string, ownerId uuid.UUID) ([]Node, error) {

	db := r.Db.GetGormDb()
	rows, err := db.Raw(`select * from nodes 
									inner join orgs ON orgs.id = nodes.org_id
									where orgs.name=? and orgs.owner=? `, orgName, ownerId).Rows()
	if err != nil {
		return nil, err
	}
	nodes, err := r.mapNodesToOrgs(rows, db)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r *nodeRepo) GetByUser(ownerId uuid.UUID) ([]Node, error) {
	db := r.Db.GetGormDb()
	rows, err := db.Raw(`select * from nodes 
									inner join orgs ON orgs.id = nodes.org_id
									where orgs.owner=?`, ownerId).Rows()
	if err != nil {
		return nil, err
	}
	nodes, err := r.mapNodesToOrgs(rows, db)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

func (r *nodeRepo) mapNodesToOrgs(rows *sql2.Rows, db *gorm.DB) ([]Node, error) {
	var nodes []Node
	defer rows.Close()

	for rows.Next() {
		var node Node
		var org Org
		err := db.ScanRows(rows, &node)
		if err != nil {
			return nil, err
		}
		err = db.ScanRows(rows, &org)
		if err != nil {
			return nil, err
		}
		node.Org = &org
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (r *nodeRepo) Delete(id ukama.NodeID) error {
	res := r.Db.GetGormDb().Delete(&Node{}, id.StringLowercase())
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (r *nodeRepo) Update(id ukama.NodeID, state NodeState) error {
	result := r.Db.GetGormDb().Where("node_id=?", id.StringLowercase()).Updates(Node{State: state})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}
