/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package store

import (
	"database/sql"
	"errors"
)

const (
	PDNLive    = "live"
	PDNClosing = "closing"
	PDNClosed  = "closed"
)

var (
	ErrPDNNotFound      = errors.New("pdn_not_found")
	ErrInstanceRetired  = errors.New("epc_instance_retired")
	ErrInstanceConflict = errors.New("epc_instance_conflict")
)

// PDN is the EPC-owned binding, independent of policy and CDR sessions.
// Closed rows are tombstones: a delayed POST cannot recreate a deleted PDN.
type PDN struct {
	InstanceID string
	ID         string
	Imsi       string
	ApnName    string
	IP         string
	State      string
}

func (s *Store) createPDNTables() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS epc_instances (
			id TEXT PRIMARY KEY,
			active INTEGER NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS epc_active_instance
		ON epc_instances(active) WHERE active = 1;
		CREATE TABLE IF NOT EXISTS pdns (
			instance_id TEXT NOT NULL,
			pdn_id TEXT NOT NULL,
			imsi TEXT NOT NULL,
			apn TEXT NOT NULL,
			ip TEXT NOT NULL,
			state TEXT NOT NULL,
			PRIMARY KEY(instance_id, pdn_id)
		);
		CREATE UNIQUE INDEX IF NOT EXISTS pdn_live_imsi
		ON pdns(imsi) WHERE state != 'closed';
		CREATE UNIQUE INDEX IF NOT EXISTS pdn_live_ip
		ON pdns(ip) WHERE state != 'closed';
	`)
	return err
}

// BeginEPC retires the previous process and reserves its PDNs for cleanup.
// Retired process identifiers cannot be reused, including after a restart.
func (s *Store) BeginEPC(id, expected string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var active bool
	err = tx.QueryRow("SELECT active FROM epc_instances WHERE id = ?",
		id).Scan(&active)
	if err == nil {
		if !active {
			return ErrInstanceRetired
		}
		return tx.Commit()
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	var current string
	err = tx.QueryRow("SELECT id FROM epc_instances WHERE active = 1").
		Scan(&current)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if current != expected {
		return ErrInstanceConflict
	}
	if _, err = tx.Exec("UPDATE epc_instances SET active = 0"); err != nil {
		return err
	}
	if _, err = tx.Exec(`UPDATE pdns SET state = ? WHERE state = ?`,
		PDNClosing, PDNLive); err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO epc_instances VALUES (?, 1)", id)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) EPCActive(id string) (bool, error) {
	var active bool
	err := s.db.QueryRow("SELECT active FROM epc_instances WHERE id = ?",
		id).Scan(&active)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	return active, err
}

func (s *Store) RetireEPCs() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err = tx.Exec("UPDATE epc_instances SET active = 0"); err != nil {
		return err
	}
	if _, err = tx.Exec("UPDATE pdns SET state = ? WHERE state = ?",
		PDNClosing, PDNLive); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) InsertPDN(p *PDN) error {
	_, err := s.db.Exec(`INSERT INTO pdns
		(instance_id, pdn_id, imsi, apn, ip, state)
		VALUES (?, ?, ?, ?, ?, ?)`,
		p.InstanceID, p.ID, p.Imsi, p.ApnName, p.IP, p.State)
	return err
}

func (s *Store) SetPDNState(p *PDN, state string) error {
	_, err := s.db.Exec(`UPDATE pdns SET state = ?
		WHERE instance_id = ? AND pdn_id = ?`,
		state, p.InstanceID, p.ID)
	if err == nil {
		p.State = state
	}
	return err
}

func scanPDN(row *sql.Row) (*PDN, error) {
	p := &PDN{}
	err := row.Scan(&p.InstanceID, &p.ID, &p.Imsi,
		&p.ApnName, &p.IP, &p.State)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrPDNNotFound
	}
	return p, err
}

func (s *Store) GetPDN(instance, id string) (*PDN, error) {
	return scanPDN(s.db.QueryRow(`SELECT instance_id, pdn_id,
		imsi, apn, ip, state FROM pdns
		WHERE instance_id = ? AND pdn_id = ?`, instance, id))
}

func (s *Store) GetLivePDN(imsi string) (*PDN, error) {
	return scanPDN(s.db.QueryRow(`SELECT instance_id, pdn_id,
		imsi, apn, ip, state FROM pdns
		WHERE imsi = ? AND state != ?`, imsi, PDNClosed))
}

func (s *Store) GetPDNByIP(ip string) (*PDN, error) {
	return scanPDN(s.db.QueryRow(`SELECT instance_id, pdn_id,
		imsi, apn, ip, state FROM pdns
		WHERE ip = ? AND state != ?`, ip, PDNClosed))
}

func (s *Store) GetLivePDNs() ([]PDN, error) {
	rows, err := s.db.Query(`SELECT instance_id, pdn_id, imsi,
		apn, ip, state FROM pdns WHERE state != ?`, PDNClosed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pdns []PDN
	for rows.Next() {
		var p PDN
		if err := rows.Scan(&p.InstanceID, &p.ID, &p.Imsi,
			&p.ApnName, &p.IP, &p.State); err != nil {
			return nil, err
		}
		pdns = append(pdns, p)
	}
	return pdns, rows.Err()
}

func (s *Store) CurrentEPC() (string, error) {
	var id string
	err := s.db.QueryRow("SELECT id FROM epc_instances WHERE active = 1").
		Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return id, err
}
