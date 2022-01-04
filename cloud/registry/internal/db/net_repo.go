package db

import (
	"github.com/jackc/pgtype"
	"github.com/pkg/errors"
	"github.com/ukama/ukamaX/common/sql"
	"github.com/ukama/ukamaX/common/ukama"
	"gorm.io/gorm/clause"
	"net"
)

type NetRepo interface {
	GetIP(nodeId ukama.NodeID) (*net.IP, error)
	SetIp(id ukama.NodeID, ip *net.IP) error
}

type netRepo struct {
	Db sql.Db
}

func NewNetRepo(db sql.Db) NetRepo {
	return &netRepo{
		Db: db,
	}
}

func (n netRepo) GetIP(nodeId ukama.NodeID) (*net.IP, error) {
	var node NodeIp
	res := n.Db.GetGormDb().Where("node_id = ?", nodeId).First(&node)
	if res.Error != nil {
		return nil, res.Error
	}

	return &node.IP.IPNet.IP, nil
}

func (n netRepo) SetIp(nodeId ukama.NodeID, ip *net.IP) error {
	node := NodeIp{
		NodeId: nodeId.StringLowercase(),
		IP:     pgtype.Inet{},
	}
	err := node.IP.Set(ip)
	if err != nil {
		return errors.Wrap(err, "error setting IP")
	}
	res := n.Db.GetGormDb().Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "lower(node_id::text)", Raw: true}},
		DoUpdates: clause.AssignmentColumns([]string{"ip"}),
	}).Create(&node)
	if res != nil {
		return res.Error
	}

	return nil
}
