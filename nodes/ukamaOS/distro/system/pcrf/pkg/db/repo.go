package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type Repo struct {
	db *sql.DB
}

type store interface {

} 
// Initialization of the SQLite database and tables (assumed to be done separately)
// var db *sql.DB

// Function to create tables if they don't exist

var db *sql.DB

func InitializeDataBase(name string) (*Repo, error) {
	repo := &repo{}
	// Open the SQLite database file
	database, err := sql.Open("sqlite3", name)
	if err != nil {
		log.Errorf("Error opening database: %s. Error %s",name, err.Error())
		return nil, err
	}
	repo.db = database

	// Create tables if they don't exist
	err := createTables()
	if err != nil {
		log.Errorf("Error creating tables %s", err.Error())
		return nil, err;
	}
	return repo, nil
}

// Function to create tables if they don't exist
func createTables() error {
	// Create Policies table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS policies (
			id INTEGER PRIMARY KEY,
			data INTEGER,
			dlbr INTEGER,
			ulbr INTEGER
		);
	`)
	if err != nil {
		log.Errorf("Error creating Policies table:", err)
		return err
	}

	// Create ReRoutes table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS reroutes (
			id INTEGER PRIMARY KEY,
			ipaddr TEXT
		);
	`)
	if err != nil {
		log.Errorf("Error creating ReRoutes table:", err)
		return err
	}

	// Create Subscribers table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS subscribers (
			id INTEGER PRIMARY KEY,
			imsi TEXT
		);
	`)
	if err != nil {
		log.Errorf("Error creating Subscribers table:", err)
		return err
	}

	// Create Usages table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS usages (
			id INTEGER PRIMARY KEY,
			subscriber_id INTEGER,
			data INTEGER,
			FOREIGN KEY(subscriber_id) REFERENCES subscribers(id)
		);
	`)
	if err != nil {
		log.Errorf("Error creating Usages table:", err)
		return err
	}

	// Create Meters table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS meters (
			id INTEGER PRIMARY KEY,
			rate INTEGER,
			type INTEGER
		);
	`)
	if err != nil {
		log.Errorf("Error creating Meters table:", err)
		return err
	}

	// Create Flows table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS flows (
			id INTEGER PRIMARY KEY,
			table INTEGER,
			priority INTEGER,
			ueipaddr TEXT,
			reroute_id INTEGER,
			meter_id INTEGER,
			FOREIGN KEY(reroute_id) REFERENCES reroutes(id),
			FOREIGN KEY(meter_id) REFERENCES meters(id)
		);
	`)
	if err != nil {
		log.Errorf("Error creating Flows table:", err)
		return err
	}

	// Create Sessions table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY,
			subscriber_id INTEGER,
			ueipaddr TEXT,
			starttime INTEGER,
			endtime INTEGER,
			txbytes INTEGER,
			rxbytes INTEGER,
			totalbytes INTEGER,
			txmeter_id INTEGER,
			rxmeter_id INTEGER,
			state INTEGER,
			FOREIGN KEY(subscriber_id) REFERENCES subscribers(id),
			FOREIGN KEY(txmeter_id) REFERENCES meters(id),
			FOREIGN KEY(rxmeter_id) REFERENCES meters(id)
		);
	`)
	if err != nil {
		log.Errorf("Error creating Sessions table:", err)
		return err
	}

	log.Infof("Tables created successfully.")
}

// CRUD operations for Policy entity

func CreateDefaultPolicy() (*Policy, error) {
	policy := Policy{
		ID:   1,
		Data: 0,
		Dlbr: 5000,
		Ulbr: 1000,
	}

	err := InsertPolicy(&policy)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

// CRUD operations for ReRoute entity

func CreateDefaultReRoute() (*ReRoute, error) {
	reroute := ReRoute{
		ID:     1,
		Ipaddr: "192.168.0.14",
	}

	err := InsertReRoute(&reroute)
	if err != nil {
		return nil, err
	}

	return &reroute, nil
}

// CRUD operations for Subscriber entity

func CreateSubscriber(imsi string) (*Subscriber, error) {
	subscriber := Subscriber{
		Imsi: imsi,
	}

	err := InsertSubscriber(&subscriber)
	if err != nil {
		return nil, err
	}

	// Create initial Usage for the subscriber
	initialUsage := Usage{
		Subscriber: subscriber,
		Data:       0,
	}

	err = InsertUsage(&initialUsage)
	if err != nil {
		return nil, err
	}

	// Allocate default policy to the subscriber
	defaultPolicy, err := CreateDefaultPolicy()
	if err != nil {
		return nil, err
	}

	subscriber.UsageID = initialUsage
	subscriber.Policy = append(subscriber.Policy, *defaultPolicy)

	err = UpdateSubscriber(&subscriber)
	if err != nil {
		return nil, err
	}

	return &subscriber, nil
}

// CRUD operations for Session entity

func CreateSession(subscriber *Subscriber, imsi string, ueIpAddr string) (*Session, error) {
	session := Session{
		Subscriber: subscriber,
		UeIpaddr:   ueIpAddr,
		StartTime:  uint64(time.Now().Unix()), // Current epoch time
		TXMeterId:  Meter{Type: TX_PATH},
		RXMeterId:  Meter{Type: RX_PATH},
		State:      Active,
	}

	// Create Flow for RX
	flowRX := Flow{
		Table:    0,
		Priority: 100,
		UeIpaddr: ueIpAddr,
		MeterID:  session.RXMeterId,
	}

	// Create Flow for TX
	flowTX := Flow{
		Table:    0,
		Priority: 100,
		UeIpaddr: ueIpAddr,
		MeterID:  session.TXMeterId,
	}

	// Check if Data in Usage is less than Policy for rerouting
	if subscriber.UsageID.Data < subscriber.Policy[0].Data {
		flowRX.ReRouting = ReRoute{} // Null rerouting
		flowTX.ReRouting = ReRoute{}
	} else {
		defaultReRoute, err := CreateDefaultReRoute()
		if err != nil {
			return nil, err
		}
		flowRX.ReRouting = *defaultReRoute
		flowTX.ReRouting = *defaultReRoute
	}

	// Insert Flows
	err := InsertFlow(&flowRX)
	if err != nil {
		return nil, err
	}

	err = InsertFlow(&flowTX)
	if err != nil {
		return nil, err
	}

	// Update session with Flow IDs
	session.RXMeterId.FlowID = flowRX.ID
	session.TXMeterId.FlowID = flowTX.ID

	err = InsertSession(&session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func EndSession(session *Session) error {
	// Update session with TX, RX, and Total bytes
	session.TxBytes = /* Set TX bytes */;
	session.RXBytes = /* Set RX bytes */;
	session.TotalBytes = session.TxBytes + session.RXBytes
	session.State = Completed

	// Update Usage for the subscriber
	subscriber, err := GetSubscriberByID(session.Subscriber.ID)
	if err != nil {
		return err
	}
	subscriber.UsageID.Data += session.TotalBytes

	// Update subscriber and session
	err = UpdateSubscriber(subscriber)
	if err != nil {
		return err
	}

	err = UpdateSession(session)
	if err != nil {
		return err
	}

	return nil
}

// Queries
// Queries

func GetUsageByImsi(imsi string) (*Usage, error) {
	var usage Usage

	err := db.QueryRow("SELECT * FROM usages WHERE subscriber_id = (SELECT id FROM subscribers WHERE imsi = ?)", imsi).
		Scan(&usage.ID, &usage.Subscriber.ID, &usage.Data)
	if err != nil {
		return nil, err
	}

	return &usage, nil
}

func GetPolicyByID(policyID int) (*Policy, error) {
	var policy Policy

	err := db.QueryRow("SELECT * FROM policies WHERE id = ?", policyID).
		Scan(&policy.ID, &policy.Data, &policy.Dlbr, &policy.Ulbr)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

func GetApplicablePolicyByImsi(imsi string) (*Policy, error) {
	var policy Policy

	err := db.QueryRow(`
		SELECT * FROM policies
		WHERE id = 1 AND (SELECT data FROM usages WHERE subscriber_id = (SELECT id FROM subscribers WHERE imsi = ?)) >= 2000000000
	`, imsi).
		Scan(&policy.ID, &policy.Data, &policy.Dlbr, &policy.Ulbr)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

func GetSessionByID(sessionID int) (*Session, error) {
	var session Session

	err := db.QueryRow("SELECT * FROM sessions WHERE id = ?", sessionID).
		Scan(&session.ID, &session.Subscriber.ID, &session.UeIpaddr, &session.StartTime, &session.EndTime, &session.TxBytes, &session.RXBytes, &session.TotalBytes, &session.TXMeterId.ID, &session.RXMeterId.ID, &session.State)
	if err != nil {
		return nil, err
	}

	// Fetch associated Subscriber
	session.Subscriber, err = GetSubscriberByID(session.Subscriber.ID)
	if err != nil {
		return nil, err
	}

	// Fetch associated Meters
	session.TXMeterId, err = GetMeterByID(session.TXMeterId.ID)
	if err != nil {
		return nil, err
	}

	session.RXMeterId, err = GetMeterByID(session.RXMeterId.ID)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func GetSessionsByImsi(imsi string) ([]Session, error) {
	var sessions []Session

	rows, err := db.Query(`
		SELECT * FROM sessions
		WHERE subscriber_id = (SELECT id FROM subscribers WHERE imsi = ?)
	`, imsi)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var session Session
		err := rows.Scan(&session.ID, &session.Subscriber.ID, &session.UeIpaddr, &session.StartTime, &session.EndTime, &session.TxBytes, &session.RXBytes, &session.TotalBytes, &session.TXMeterId.ID, &session.RXMeterId.ID, &session.State)
		if err != nil {
			return nil, err
		}

		// Fetch associated Subscriber
		session.Subscriber, err = GetSubscriberByID(session.Subscriber.ID)
		if err != nil {
			return nil, err
		}

		// Fetch associated Meters
		session.TXMeterId, err = GetMeterByID(session.TXMeterId.ID)
		if err != nil {
			return nil, err
		}

		session.RXMeterId, err = GetMeterByID(session.RXMeterId.ID)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)subscriber.Imsi
	}

	return sessions, nil
}

func GetActiveSessionByImsi(imsi string) (*Session, error) {
	var session Session

	err := db.QueryRow(`
		SELECT * FROM sessions
		WHERE subscriber_id = (SELECT id FROM subscribers WHERE imsi = ?) AND state = 1
	`, imsi).
		Scan(&session.ID, &session.Subscriber.ID, &session.UeIpaddr, &session.StartTime, &session.EndTime, &session.TxBytes, &session.RXBytes, &session.TotalBytes, &session.TXMeterId.ID, &session.RXMeterId.ID, &session.State)
	if err != nil {
		return nil, err
	}

	// Fetch associated Subscriber
	session.Subscriber, err = GetSubscriberByID(session.Subscriber.ID)
	if err != nil {
		return nil, err
	}

	// Fetch associated Meters
	session.TXMeterId, err = GetMeterByID(session.TXMeterId.ID)
	if err != nil {
		return nil, err
	}

	session.RXMeterId, err = GetMeterByID(session.RXMeterId.ID)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func GetAllActiveSessions() ([]Session, error) {
	var sessions []Session

	rows, err := db.Query("SELECT * FROM sessions WHERE state = 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var session Session
		err := rows.Scan(&session.ID, &session.Subscriber.ID, &session.UeIpaddr, &session.StartTime, &session.EndTime, &session.TxBytes, &session.RXBytes, &session.TotalBytes, &session.TXMeterId.ID, &session.RXMeterId.ID, &session.State)
		if err != nil {
			return nil, err
		}

		// Fetch associated Subscriber
		session.Subscriber, err = GetSubscriberByID(session.Subscriber.ID)
		if err != nil {
			return nil, err
		}

		// Fetch associated Meters
		session.TXMeterId, err = Get
	}
}


// Update operations

// Update operations

func UpdateReroute(reRoute *ReRoute) error {
	_, err := db.Exec(`
		UPDATE reroutes
		SET ipaddr = ?
		WHERE id = ?;
	`, reRoute.Ipaddr, reRoute.ID)
	return err
}

func UpdateSubscriber(subscriber *Subscriber) error {
	_, err := db.Exec(`
		UPDATE subscribers
		SET imsi = ?
		WHERE id = ?;
	`, subscriber.Imsi, subscriber.ID)
	return err
}

func UpdatePolicy(policy *Policy) error {
	_, err := db.Exec(`
		UPDATE policies
		SET data = ?, dlbr = ?, ulbr = ?
		WHERE id = ?; 
		`, policy.Data, policy.Dlbr, policy.Ulbr, policy.ID)
		return err
}

func UpdateUsage(usage *Usage) error {
	_, err := db.Exec(`
		UPDATE usages
		SET data = ?
		WHERE id = ?;
	`, usage.Data, usage.ID)
	return err
}

func UpdateMeter(meter *Meter) error {
	_, err := db.Exec(`
		UPDATE meters
		SET rate = ?, type = ?
		WHERE id = ?;
	`, meter.Rate, meter.Type, meter.ID)
	return err
}

func UpdateFlow(flow *Flow) error {
	_, err := db.Exec(`
		UPDATE flows
		SET table = ?, priority = ?, ueipaddr = ?, reroute_id = ?, meter_id = ?
		WHERE id = ?;
	`, flow.Table, flow.Priority, flow.UeIpaddr, flow.ReRouting.ID, flow.MeterID.ID, flow.ID)
	return err
}

func UpdateSession(session *Session) error {
	_, err := db.Exec(`
		UPDATE sessions
		SET ueipaddr = ?, starttime = ?, endtime = ?, txbytes = ?, rxbytes = ?, totalbytes = ?, txmeter_id = ?, rxmeter_id = ?, state = ?
		WHERE id = ?;
	`, session.UeIpaddr, session.StartTime, session.EndTime, session.TxBytes, session.RXBytes, session.TotalBytes, session.TXMeterId.ID, session.RXMeterId.ID, session.State, session.ID)
	return err
}


// ... (similar update operations for other entities)

// CRUD operations for Policy entity

func InsertPolicy(policy *Policy) error {
	_, err := db.Exec(`
		INSERT INTO policies (id, data, dlbr, ulbr)
		VALUES (?, ?, ?, ?);
	`, policy.ID, policy.Data, policy.Dlbr, policy.Ulbr)
	return err
}

func GetPolicyByID(policyID int) (*Policy, error) {
	var policy Policy

	err := db.QueryRow("SELECT * FROM policies WHERE id = ?", policyID).
		Scan(&policy.ID, &policy.Data, &policy.Dlbr, &policy.Ulbr)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

// CRUD operations for ReRoute entity

func InsertReRoute(reRoute *ReRoute) error {
	_, err := db.Exec(`
		INSERT INTO reroutes (id, ipaddr)
		VALUES (?, ?);
	`, reRoute.ID, reRoute.Ipaddr)
	return err
}

func GetReRouteByID(reRouteID int) (*ReRoute, error) {
	var reRoute ReRoute

	err := db.QueryRow("SELECT * FROM reroutes WHERE id = ?", reRouteID).
		Scan(&reRoute.ID, &reRoute.Ipaddr)
	if err != nil {
		return nil, err
	}

	return &reRoute, nil
}

// ... (similar CRUD operations for other entities)
