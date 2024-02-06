package store

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/api"
)

type Store struct {
	db *sql.DB
}

// Initialization of the SQLite database and tables (assumed to be done separately)
// var db *sql.DB

// Function to create tables if they don't exist


func NewStore(name string) (*Store, error) {
	repo := &Store{}
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

func (s *Store) createPolicyTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS Policy (
			ID UUID PRIMARY KEY,
			Data INTEGER,
			Dlbr INTEGER,
			Ulbr INTEGER,
			StartTime INTEGER,
			EndTime INTEGER
		);
	`)
	if err != nil {
		log.Errorf("Error creating Policies table:", err)
		return err
	}
	return nil
}

func (s *Store) createRerouteTable() error {
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS reroutes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ipaddr TEXT UNIQUE
		);
	`)
	if err != nil {
		log.Errorf("Error creating Reroute table:", err)
		return err
	}
	return nil
}

func (s *Store) createSubscriberTable() error {
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS subscribers (
			id UUID PRIMARY KEY,
			imsi TEXT UNIQUE
		);
	`)
	if err != nil {
		log.Errorf("Error creating Subscriber table:", err)
		return err
	}
	return nil
}

func (s *Store) createUsageTable() error {
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS usages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			subscriber_id UUID,
			data INTEGER,
			updates_at INTEGER, 
			FOREIGN KEY(subscriber_id) REFERENCES subscribers(id)
		);
	`)
	if err != nil {
		log.Errorf("Error creating Usage table:", err)
		return err
	}
	return nil
}

func (s *Store) createMeterTable() error {
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS meters (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			rate INTEGER,
			type INTEGER
		);
	`)
	if err != nil {
		log.Errorf("Error creating Meter table:", err)
		return err
	}
	return nil
}

func (s *Store) createFlowTable() error {
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS flows (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		table INTEGER,
		cookie INTEGER CHECK(cookie >= 0),
		priority INTEGER,
		ueipaddr TEXT,
		reroute_id INTEGER,
		meter_id INTEGER,
		FOREIGN KEY(reroute_id) REFERENCES reroutes(id),
		FOREIGN KEY(meter_id) REFERENCES meters(id)
	);
`)
	if err != nil {
		log.Errorf("Error creating Flow table:", err)
		return err
	}
	return nil
}

func (s *Store) createSessionTable() error {
	_, err = s.db.Exec(`
	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY,
		subscriber_id UUID,
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
	`)
	if err != nil {
		log.Errorf("Error creating Session table:", err)
		return err
	}
	return nil
}

func (s *Store) CreateTables() {// Enable the UUID extension
	_, err := db.Exec("SELECT load_extension('libsqlite3_mod_uuid.so')")
	if err != nil {
		log.Fatal(err)
	}
	
	err = s.createPolicyTable();
	if err != nil {
		return err
	}

	err = s.createUsageTable();
	if err != nil {
		return err
	}
	
	err = s.createRerouteTable();
	if err != nil {
		return err
	}
	
	err = s.createSubscriberTable();
	if err != nil {
		return err
	}
	
	err = s.createMeterTable();
	if err != nil {
		return err
	}

	err = s.createFlowTable();
	if err != nil {err = s.createMeterTable();
		if err != nil {
			return err
		}
		return err
	}

}

/* Create a policy */
func (s *Store) CreatePolicy(p *api.Policy) (*Policy, error) {
	policy := Policy{
		ID:   p.Uuid,
		Data: p.Data,
		Dlbr: p.Dlbr,
		Ulbr: p.Ulbr,
		StartTime: p.StartTime,
		EndTime: p.EndTime,
	}

	err := s.InsertPolicy(&policy)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

/* Create  a new route */
func (s *Store) CreateReroute(r *api.Reroute) (*ReRoute, error) {
	reroute := ReRoute{
		Ipaddr: r.Ip,
	}

	err := s.InsertReRoute(&reroute)
	if err != nil {
		return nil, err
	}

	return &reroute, nil
}



/* Create a new meter */
func (s *Store) CreateMeter(s *Subscriber, p *Policy, typeM int) (*Meter, error) {
	var r sql.Result

	if (typeM == RX_PATH) {
		r, err := s.db.Exec(`
			INSERT INTO meters (rate, type, burst)
			VALUES (?, ?, ?);
		`, p.Ulbr, typeM, p.burst)
	} else if (typeM == TX_PATH) {
		r, err := s.db.Exec(`
			INSERT INTO meters (rate, type, burst)
			VALUES (?, ?, ?);
		`, p.Dlbr, typeM, p.burst)
	}

	// Get the ID of the last inserted Subscriber
	mId, err := r.LastInsertId()
	if err != nil {
		log.Errorf("Failed to get last inserted meter. Error: %v", err)
		return nil, err
	}

	meter, err = s.GetMeter(mId)
	if err != nil {
		return nil, err
	}
	return &meter, nil
}

/* Delete meter */
func (s *Store) DeleteMeter(id uint32)  error {

	r, err := s.db.Exec(`
		DELETE FROM meters 
		WHERE id = ?;
	`, id)	
	if err != nil {
		log.Errorf("Failed to delete meter %d. Error: %v", id, err)
		return err
	}

	return nil
}

/* Get meter */
func (s *Store) GetMeter(id uint32)  (*Meter, error) {
	var meter Meter
	r, err := s.db.QueryRow(`
		SELECT FROM meters 
		WHERE id = ?;
	`, id).Scan(&meter.ID, &meter.Rate, &meter.Type, &meter.Burst)
	if err != nil {
		log.Errorf("Failed to get meter %d. Error: %v", id, err)
		return nil, err
	}

	return &meter, nil
}


/* Create a new flow */
func (s *Store) CreateFlow(m *Meter, r *Reroute, ip string, table, priority uint32 ) (*Flow, error) {
	r, err := s.db.Exec(`
		INSERT INTO flows (table, priority, ueipaddr, reroute_id, meter_id)
		VALUES (?, ?, ?, ?, ?);
	`, table, priority, ip, r, m)
	if err != nil {
		log.Errorf("Failed to create flow for UE %s meter %d. Error: %v", ip, m.ID, err)
		return err
	}
	
	// Get the ID of the last inserted Subscriber
	fId, err := r.LastInsertId()
	if err != nil {
		log.Errorf("Failed to get last inserted flow. Error: %v", err)
		return nil, err
	}

	flow, err = s.GetFlow(fId)
	if err != nil {
		return nil, err
	}
	return &flow, nil
}

/* Delete flow */
func (s *Store) DeleteFlow(id uint32)  error {

	r, err := s.db.Exec(`
		DELETE FROM flows 
		WHERE id = ?;
	`, id)
	if err != nil {
		log.Errorf("Failed to delete flow %d. Error: %v", id, err)
		return err
	}

	return nil
}

/* Get flow */
func (s *Store) GetFlow(id uint32)  (*Flow, error) {
	var f Flow
	r, err := s.db.QueryRow(`
		SELECT FROM flows 
		WHERE id = ?;
	`, id).Scan(&f.ID, &f.Cookie, &f.Table, &f.Priority, &f.UeIpaddr, &f.MeterID.ID, &f.ReRouting)
	if err != nil {
		log.Errorf("Failed to get flow %d. Error: %v", id, err)
		return nil, err
	}

	return &f, nil
}


/* Create a subscriber */
func (s *Store) CreateSubscriber(imsi string, p *api.Policy) (*Subscriber, error) {
	subscriber := Subscriber{
		Imsi: imsi,
	}

	err := s.InsertSubscriber(&subscriber)
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
	defaultPolicy, err := CreatePolicy()
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
func (s *Store) InsertSession(s *Sesssion) (*Session, error) {
	_, err := s.db.Exec(`
		INSERT INTO sessions (subscriber_id, apn_name, ueipaddr, starttime, endtime , txbytes , rxbytes , totalbytes , txmeter_id, rxmeter_id, state)
		VALUES (s.SubsciberID, s.ApnName, s.UeIpaddr, s.StartTime, s.EndTime, s.TxBytes, s.RXBytes, s.TotalBytes, s.TXMeterId.ID, s.RXMeterId.ID, s.State);
	`)
	return err
}

func (s *Store) DeleteSession(s *Subscriber) error {
	_, err := s.db.Exec(`
	DELETE FROM sessions WHERE subscriber_id = ? ;
	`, s.Subscriber)		
	
	if err != nil {
		log.Errorf("Failed to delete seession for subscriber %s: Error: %v", s.Subscriber, err.Error())
		return err
	}
	return nil
}


func (s *Store) UpdateSessionUsage(session *Session) error {
	_, err := s.db.Exec(`
		UPDATE sessions
		SET txbytes = ?, rxbytes = ?, totalbytes = ?
		WHERE id = ?;
	`, session.TxBytes, session.RXBytes, session.TotalBytes, session.ID)
	return err
}

func (s *Store) UpdateSession(session *Session) error {
	_, err := s.db.Exec(`
		UPDATE sessions
		SET endtime = ?, txbytes = ?, rxbytes = ?, totalbytes = ?, state = ?
		WHERE id = ?;
	`, session.EndTime, session.TxBytes, session.RXBytes, session.TotalBytes, session.State, session.ID)
	return err
}

func (s *Store) CreateSession(subscriber *Subscriber, ueIpAddr string) (*Session, error) {
	
	rxM , err := s.CreateMeter(subscriber, &subscriber.PolicyID, RX_PATH, ueIpAddr)
	if err != nil {
		return nil, err
	}

	txM , err := s.CreateMeter(subscriber, &subscriber.PolicyID, TX_PATH, ueIpAddr)
	if err != nil {
		return nil, err
	}

 	session := Session{
		Subscriber: subscriber,
		UeIpaddr:   ueIpAddr,
		StartTime:  uint64(time.Now().Unix()), // Current epoch time
		TXMeterId:  txM,
		RXMeterId:  rxM,
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
	if subscriber.UsageID.Data >= subscriber.Policy.Data {
		log.Errof("can't create flows. UE %s has reached max data cap.", subscriber.Imsi)
		return fmt.Errorf("max data cap exceeded")
	}

	// Insert Flows
	rxF, err := s.CreateFlow(&flowRX)
	if err != nil {
		return nil, err
	}

	txF, err := s.CreateFlow(&flowTX)
	if err != nil {
		return nil, err
	}

	err = s.InsertSession(&session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *Store) EndSession(session *Session) error {
	// Update session with TX, RX, and Total bytes
	// session.TxBytes = /* Set TX bytes */;
	// session.RXBytes = /* Set RX bytes */;
	session.EndTime = uint64(time.Now().Unix())
	session.TotalBytes = session.TxBytes + session.RXBytes
	session.State = Completed

	// Update Usage for the subscriber
	subscriber, err := s.GetSubscriberByID(session.Subscriber.ID)
	if err != nil {
		return err
	}
	subscriber.UsageID.Data += session.TotalBytes

	err = s.UpdateSession(session)
	if err != nil {
		return err
	}

	// Update subscriber and session
	err = s.UpdateUsage(subscriber.UsageID)
	if err != nil {
		return err
	}

	return nil
}

/* Update Usage */
func (s *Store) UpdateUsage(usage *Usage) error {
	_, err := s.db.Exec(`
		UPDATE usages
		SET data = ?, epoch = ?,
		WHERE id = ?;
	`, usage.Data, usage.Epoch, usage.ID)
	return err
}

/* Get usage by imsi */
func (s *Store) GetUsageByImsi(imsi string) (*Usage, error) {
	var usage Usage

	err := s.db.QueryRow("SELECT * FROM usages WHERE subscriber_id = (SELECT id FROM subscribers WHERE imsi = ?)", imsi).
		Scan(&usage.ID, &usage.Subscriber.ID, &usage.Data)
	if err != nil {
		return nil, err
	}

	return &usage, nil
}

func (s *Store) GetPolicyByID(policyID int) (*Policy, error) {
	var policy Policy

	err := s.db.QueryRow("SELECT * FROM policies WHERE id = ?", policyID).
		Scan(&policy.ID, &policy.Data, &policy.Dlbr, &policy.Ulbr)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

func (s *Store) GetApplicablePolicyByImsi(imsi string) (*Policy, error) {
	var policy Policy

	err := s.db.QueryRow(`
		SELECT * FROM policies
		WHERE id = 1 AND (SELECT data FROM usages WHERE subscriber_id = (SELECT id FROM subscribers WHERE imsi = ?)) >= 2000000000
	`, imsi).
		Scan(&policy.ID, &policy.Data, &policy.Dlbr, &policy.Ulbr)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}

func (s *Store) GetSessionByID(sessionID int) (*Session, error) {
	var session Session

	err := s.db.QueryRow("SELECT * FROM sessions WHERE id = ?", sessionID).
		Scan(&session.ID, &session.Subscriber.ID, &session.UeIpaddr, &session.StartTime, &session.EndTime, &session.TxBytes, &session.RXBytes, &session.TotalBytes, &session.TXMeterId.ID, &session.RXMeterId.ID, &session.State)
	if err != nil {
		return nil, err
	}

	// Fetch associated Subscriber
	session.Subscriber, err = s.GetSubscriberByID(session.Subscriber.ID)
	if err != nil {
		return nil, err
	}

	// Fetch associated Meters
	session.TXMeterId, err = s.GetMeter(session.TXMeterId.ID)
	if err != nil {
		return nil, err
	}

	session.RXMeterId, err = s.GetMeter(session.RXMeterId.ID)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *Store) GetSessionsByImsi(imsi string) ([]Session, error) {
	var sessions []Session

	rows, err := s.db.Query(`
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
		session.Subscriber, err = s.GetSubscriberByID(session.Subscriber.ID)
		if err != nil {
			return nil, err
		}

		// Fetch associated Meters
		session.TXMeterId, err = s.GetMeter(session.TXMeterId.ID)
		if err != nil {
			return nil, err
		}

		session.RXMeterId, err = s.GetMeter(session.RXMeterId.ID)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (s *Store) GetActiveSessionByImsi(imsi string) (*Session, error) {
	var session Session

	err := s.db.QueryRow(`
		SELECT * FROM sessions
		WHERE subscriber_id = (SELECT id FROM subscribers WHERE imsi = ?) AND state = 1
	`, imsi).
		Scan(&session.ID, &sesessionssion.Subscriber.ID, &session.UeIpaddr, &session.StartTime, &session.EndTime, &session.TxBytes, &session.RXBytes, &session.TotalBytes, &session.TXMeterId.ID, &session.RXMeterId.ID, &session.State)
	if err != nil {
		return nil, err
	}

	// Fetch associated Subscriber
	session.Subscriber, err = s.GetSubscriberByID(session.Subscriber.ID)
	if err != nil {
		return nil, err
	}

	// Fetch associated Meters
	session.TXMeterId, err = s.GetMeter(session.TXMeterId.ID)
	if err != nil {
		return nil, err
	}

	session.RXMeterId, err = s.GetMeter(session.RXMeterId.ID)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *Store) GetAllActiveSessions() ([]Session, error) {
	var sessions []Session

	rows, err := s.db.Query("SELECT * FROM sessions WHERE state = 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var session Session
		err := rows.Scan(&session.ID, &session.Subscriber.ID, &session.UeIpaddr, &session.StartTime, &session.EndTime, &session.TxBytes, &session.RXBytes, &session.TotalBytes, &session.TXMeterId.ID, &session.RXMeterId.ID, &session.State)
		if err != nil {
			return nil, errsession
		}

		// Fetch associated Subscriber
		session.Subscriber, err = s.GetSubscriberByID(session.Subscriber.ID)
		if err != nil {
			return nil, err
		}

		// Fetch associated Meters
		session.TXMeterId, err = Get
	}
}


func (s *Store) UpdatePolicy(policy *Policy) error {
	_, err := s.db.Exec(`
		UPDATE policies
		SET data = ?, dlbr = ?, ulbr = ?
		WHERE id = ?; 
		`, policy.Data, policy.Dlbr, policy.Ulbr, policy.ID)
		return err
}



func (s *Store) UpdateMeter(meter *Meter) error {
	_, err := s.db.Exec(`
		UPDATE meters
		SET rate = ?, type = ?
		WHERE id = ?;
	`, meter.Rate, meter.Type, meter.ID)
	return err
}

func (s *Store) UpdateFlow(flow *Flow) error {
	_, err := s.db.Exec(`
		UPDATE flows
		SET table = ?, priority = ?, ueipaddr = ?, reroute_id = ?, meter_id = ?
		WHERE id = ?;
	`, flow.Table, flow.Priority, flow.UeIpaddr, flow.ReRouting.ID, flow.MeterID.ID, flow.ID)
	return err
}



/* CRUD operations for Policy entity */
func (s *Store) InsertPolicy(policy *Policy) error {
	_, err := s.db.Exec(`
		INSERT OR IGNORE INTO policies (id, data, dlbr, ulbr)
		VALUES (?, ?, ?, ?);
	`, policy.ID, policy.Data, policy.Dlbr, policy.Ulbr)
	return err
}

func (s *Store) GetPolicyByID(policyID uuid.UUID) (*Policy, error) {
	var policy Policy

	err := s.db.QueryRow("SELECT * FROM policies WHERE id = ?", policyID).
		Scan(&policy.ID, &policy.Data, &policy.Dlbr, &policy.Ulbr)
	if err != nil {
		return nil, err
	}

	return &policy, nil
}


/* CRUD operations for Reroute entity */
func (s *Store) InsertReRoute(reRoute *ReRoute) error {
	_, err := s.db.Exec(`
		INSERT OR IGNORE INTO reroutes (ipaddr)
		VALUES (?); `, reRoute.Ipaddr)
	return err
}

func (s *Store) GetReRouteByID(reRouteID int) (*ReRoute, error) {
	var reRoute ReRoute

	err := s.db.QueryRow("SELECT * FROM reroutes WHERE id = ?", reRouteID).
		Scan(&reRoute.ID, &reRoute.Ipaddr)
	if err != nil {
		return nil, err
	}

	return &reRoute, nil
}

/* Delete a route */
func (s *Store) DeleteReRoute(reRoute *ReRoute) error
{
	_, err := s.db.Exec(`
		DELETE FROM reroutes 
		WHERE ipaddr = ?; `, reRoute.Ipaddr)
	return err
}

func (s *Store) UpdateReroute(reRoute *ReRoute) error {
	_, err := s.db.Exec(`
		UPDATE reroutes
		SET ipaddr = ?
		WHERE id = ?;
	`, reRoute.Ipaddr, reRoute.ID)
	return err
}

/* CRUD operations for Subscriber entity */
func (s *Store) InsertSubscriber(s *api.Subscriber ) error {
	_, err := s.db.Exec(`
		INSERT OR IGNORE INTO subscriber (s.ID, s.Imsi)
		VALUES (?, ?);
	`, s.ID, s.Imsi)
	return err
}

/* Update policy for Subscriber entity */
func (s *Store) UpdateSubscriber(subscriber *Subscriber, p uuid.UUID) error {
	_, err := s.db.Exec(`
		UPDATE subscribers
		SET policy_id = ?
		WHERE id = ?;
	`, p, subscriber.ID)
	return err
}


/* Update policy for Subscriber entity */
func (s *Store) DeleteSubscriber(subscriber *Subscriber) error {
	_, err := s.db.Exec(`
		DELETE FROM subscribers
		WHERE id = ?;
	`, subscriber.ID)
	return err
}

func (s *Store) GetSubscriber(imsi string) (*Subscriber, error) {

	query := "SELECT ID, Imsi FROM Subscriber WHERE Imsi = ?"
	row := s.db.QueryRow(query, imsi)

	var subscriber Subscriber
	err = row.Scan(&subscriber.ID, &subscriber.Imsi)
	if err != nil {
		// Subscriber not found
		return nil, fmt.Errorf("Subscriber not found: %v", err)
	}

	return &subscriber,nil
}

func (s *Store) CreateSubscriberOrUpdatePolicy(s *api.Subscriber, p uuid.UUID) error{
	
	// Check if the subscriber already exists
	var subscriberID uuid.UUID
	err = s.db.QueryRow("SELECT ID FROM Subscriber WHERE Imsi = ?", s.Imsi).Scan(&subscriberID)

	if err == nil && subscriberID != 0 {
		// Subscriber already exists, update the policy
		return s.UpdateSubscriber(subscriberID, p)
	} else {
		 err := s.InsertSubscriber(s)
		 if err != nil {
			return err
		 }
		 log.Infof("New subscriber with Imsi %s created. ID: %d\n", imsi, newSubscriberID)
		
		 // Get the ID of the last inserted Subscriber
		sub, err := s.GetSubscriber(s.Imsi)
		if err != nil {
			log.Errorf("Erorr while getting sunscriber with imsi %s. Error %s", s.Imsi, err.Error())
			return err
		}

		err = s.UpdateSubscriber(sub, p)
		if err != nil {
			return err
		}
		log.Infof("Policy %s assigned to the new subscriber %s.", p, sub.Imsi)
	}
}
