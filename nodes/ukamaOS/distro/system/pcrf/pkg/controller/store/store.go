package store

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"github.com/ukama/ukama/nodes/ukamaOS/distro/system/pcrf/pkg/api"
	"github.com/ukama/ukama/systems/common/uuid"
)

type Store struct {
	db *sql.DB
}

// Initialization of the SQLite database and tables (assumed to be done separately)
// var db *sql.DB

// Function to create tables if they don't exist

func NewStore(name string) (*Store, error) {
	repo := &Store{}
	sql.Register("sqlite3_with_extensions",
		&sqlite3.SQLiteDriver{
			Extensions: []string{
				"libsqlite3_uuid",
			},
		})
	// Open the SQLite database file
	database, err := sql.Open("sqlite3", "/tmp/pcrf.db")
	if err != nil {
		log.Errorf("Error opening database: %s. Error %s", name, err.Error())
		return nil, err
	}
	repo.db = database

	// Create tables if they don't exist
	err = repo.CreateTables()
	if err != nil {
		log.Errorf("Error creating tables %s", err.Error())
		return nil, err
	}
	return repo, nil
}

func (s *Store) createPolicyTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS policies ( 
			id BLOB PRIMARY KEY CHECK(length(id) = 16),
			data INTEGER,
			dlbr INTEGER,
			ulbr INTEGER,
			starttime INTEGER,
			endtime INTEGER,
			burst INTEGER
		);
	`)
	if err != nil {
		log.Errorf("Error creating Policies table.Error %s", err.Error())
		return err
	}
	return nil
}

func (s *Store) createRerouteTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS reroutes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			ipaddr TEXT UNIQUE
		);
	`)
	if err != nil {
		log.Errorf("Error creating Reroute table.Error %s", err.Error())
		return err
	}
	return nil
}

func (s *Store) createSubscriberTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS subscribers (
			id INTEGER PRIMARY KEY,
			imsi TEXT UNIQUE,
			policy_id BLOB CHECK(length(policy_id) = 16),
			reroute_id INTEGER,
			FOREIGN KEY(policy_id) REFERENCES policies(id),
			FOREIGN KEY(reroute_id) REFERENCES reroutes(id)
		);
	`)
	if err != nil {
		log.Errorf("Error creating Subscriber table.Error %s", err.Error())
		return err
	}
	return nil
}

func (s *Store) createUsageTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS usages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			subscriber_id INTEGER,
			data INTEGER,
			updatedat INTEGER, 
			FOREIGN KEY(subscriber_id) REFERENCES subscribers(id)
		);
	`)
	if err != nil {
		log.Errorf("Error creating Usage table.Error %s", err.Error())
		return err
	}
	return nil
}

func (s *Store) createMeterTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS meters (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			rate INTEGER,
			burst INTEGER,
			type INTEGER
		);
	`)
	if err != nil {
		log.Errorf("Error creating Meter table.Error %s", err.Error())
		return err
	}
	return nil
}

func (s *Store) createFlowTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS flows (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		tableid INTEGER,
		cookie INTEGER UNIQUE,
		priority INTEGER,
		ueipaddr TEXT,
		reroute_id INTEGER,
		meter_id INTEGER,
		FOREIGN KEY(reroute_id) REFERENCES reroutes(id),
		FOREIGN KEY(meter_id) REFERENCES meters(id),
		CHECK (cookie >= 0)
	);
`)
	if err != nil {
		log.Errorf("Error creating Flow table.Error %s", err.Error())
		return err
	}
	return nil
}

func (s *Store) createSessionTable() error {
	_, err := s.db.Exec(`
			CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			subscriber_id INTEGER,
			apnname TEXT,
			ueipaddr TEXT,
			starttime INTEGER,
			endtime INTEGER,
			txbytes INTEGER,
			rxbytes INTEGER,
			totalbytes INTEGER,
			txmeter_id INTEGER,
			rxmeter_id INTEGER,
			state INTEGER,
			sync INTEGER,
			FOREIGN KEY(subscriber_id) REFERENCES subscribers(id),
			FOREIGN KEY(txmeter_id) REFERENCES meters(id),
			FOREIGN KEY(rxmeter_id) REFERENCES meters(id)
		);
	`)
	if err != nil {
		log.Errorf("Error creating Session table.Error %s", err.Error())
		return err
	}
	return nil
}

func (s *Store) CreateTables() error { // Enable the UUID extension
	// _, err := s.db.Exec("SELECT load_extension('/usr/lib/libsqlite3_mod_uuid.so')")
	// if err != nil {
	// 	log.Errorf("Failed to load uuid extension. Error: %s", err.Error())
	// 	log.Fatal(err)
	// }

	err := s.createPolicyTable()
	if err != nil {
		return err
	}

	err = s.createUsageTable()
	if err != nil {
		return err
	}

	err = s.createRerouteTable()
	if err != nil {
		return err
	}

	err = s.createSubscriberTable()
	if err != nil {
		return err
	}

	err = s.createMeterTable()
	if err != nil {
		return err
	}

	err = s.createFlowTable()
	if err != nil {
		return err
	}

	err = s.createMeterTable()
	if err != nil {
		return err
	}

	err = s.createSessionTable()
	if err != nil {
		return err
	}
	return nil
}

/* Create a policy */
func (s *Store) CreatePolicy(p *api.Policy) (*Policy, error) {
	policy := Policy{
		ID:        p.Uuid,
		Burst:     p.Burst,
		Data:      p.Data,
		Dlbr:      p.Dlbr,
		Ulbr:      p.Ulbr,
		StartTime: p.StartTime,
		EndTime:   p.EndTime,
	}

	err := s.InsertPolicy(&policy)
	if err != nil {
		log.Errorf("Error inserting policy %v.Error: %v", policy.ID.Bytes(), err)
		return nil, err
	}

	log.Infof("Created policy %v", policy)
	return &policy, nil
}

/* Create  a new route */
func (s *Store) CreateReroute(r *api.ReRoute) (*ReRoute, error) {
	reroute := ReRoute{
		IpAddr: r.Ip,
	}

	rr, err := s.GetReRouteByIP(r.Ip)
	if err != nil {
		if err == sql.ErrNoRows {
			err = s.InsertReRoute(&reroute)
			if err != nil {
				return nil, err
			}
			return s.GetReRouteByIP(r.Ip)
		} else {
			return nil, err
		}
	}

	log.Infof("Created route %+v", rr)
	return rr, nil
}

/* Create a new meter */
func (s *Store) CreateMeter(sub *Subscriber, p *Policy, typeM int) (*Meter, error) {
	var r sql.Result
	var err error
	if typeM == RX_PATH {
		r, err = s.db.Exec(`
			INSERT INTO meters (rate, type, burst)
			VALUES (?, ?, ?);
		`, p.Ulbr, typeM, p.Burst)
	} else if typeM == TX_PATH {
		r, err = s.db.Exec(`
			INSERT INTO meters (rate, type, burst)
			VALUES (?, ?, ?);
		`, p.Dlbr, typeM, p.Burst)
	}
	if err != nil {
		log.Errorf("Failed to insert meter for subscriber %s.Error: %s", sub.Imsi, err.Error())
		return nil, err
	}

	// Get the ID of the last inserted Subscriber
	mId, err := r.LastInsertId()
	if err != nil {
		log.Errorf("Failed to get last inserted meter. Error: %v", err)
		return nil, err
	}

	meter, err := s.GetMeter(uint32(mId))
	if err != nil {
		return nil, err
	}

	log.Infof("Created meter %+v", meter)
	return meter, nil
}

/* Delete meter */
func (s *Store) DeleteMeter(id uint32) error {

	_, err := s.db.Exec(`
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
func (s *Store) GetMeter(id uint32) (*Meter, error) {
	var meter Meter
	err := s.db.QueryRow(`
		SELECT id,rate,type,burst FROM meters 
		WHERE id = ?;
	`, id).Scan(&meter.ID, &meter.Rate, &meter.Type, &meter.Burst)
	if err != nil {
		log.Errorf("Failed to get meter %d. Error: %s", id, err.Error())
		return nil, err
	}

	return &meter, nil
}

/* Create a new flow */
func (s *Store) CreateFlow(m *Meter, r *ReRoute, ip string, table, priority uint32) (*Flow, error) {
	var ck uint64

	for {
		ck = uint64(rand.Uint32())
		b, err := s.CheckUniqueCookie(ck)
		if err != nil {
			log.Errorf("Failed to check unique cookie.Error: %s", err.Error())
			return nil, err
		}
		if b {
			break
		}
	}

	res, err := s.db.Exec(`
		INSERT INTO flows (cookie,tableid, priority, ueipaddr, reroute_id, meter_id)
		VALUES (?, ?, ?, ?, ?, ?);
	`, ck, table, priority, ip, r.ID, m.ID)
	if err != nil {
		log.Errorf("Failed to create flow for UE %s meter %d. Error: %v", ip, m.ID, err)
		return nil, err
	}

	// Get the ID of the last inserted Subscriber
	fId, err := res.LastInsertId()
	if err != nil {
		log.Errorf("Failed to get last inserted flow. Error: %v", err)
		return nil, err
	}

	flow, err := s.GetFlow(int(fId))
	if err != nil {
		return nil, err
	}

	log.Infof("Created flow %+v", flow)
	return flow, nil
}

/* Delete flow */
func (s *Store) DeleteFlow(id uint32) error {

	_, err := s.db.Exec(`
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
func (s *Store) GetFlow(id int) (*Flow, error) {
	var f Flow
	err := s.db.QueryRow(`
		SELECT id,cookie,tableid,priority,ueipaddr,meter_id,reroute_id FROM flows 
		WHERE id = ?;
	`, id).Scan(&f.ID, &f.Cookie, &f.Tableid, &f.Priority, &f.UeIpAddr, &f.MeterID.ID, &f.ReRouting.ID)
	if err != nil {
		log.Errorf("Failed to get flow %d. Error: %v", id, err)
		return nil, err
	}

	m, err := s.GetMeter(uint32(f.MeterID.ID))
	if err != nil {
		log.Errorf("Failed to get meter %d for flow %d. Error: %v", f.MeterID.ID, f.ID, err)
		return nil, err
	}
	f.MeterID = *m

	r, err := s.GetReRouteByID(f.ReRouting.ID)
	if err != nil {
		log.Errorf("Failed to get reeoute %d for flow %d. Error: %v", f.ReRouting.ID, f.ID, err)
		return nil, err
	}
	f.ReRouting = *r

	return &f, nil
}

/* Get flow */
func (s *Store) CheckUniqueCookie(cookie uint64) (bool, error) {
	var count int
	err := s.db.QueryRow(`
		SELECT COUNT(*) FROM flows 
		WHERE cookie = ?;
	`, cookie).Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (s *Store) CreateUsage(sub *Subscriber) error {
	_, err := s.db.Exec(`
		INSERT INTO usages (subscriber_id, updatedat, data)
		VALUES (?, ?, ?);
	`, sub.ID, time.Now().Unix(), 0)
	if err != nil {
		log.Errorf("Failed to create usage for subscriber %s. Error: %v", sub.Imsi, err)
		return err
	}
	return nil
}

/* Create a subscriber */
func (s *Store) CreateSubscriber(imsi string, p *api.Policy, ip *string) (*Subscriber, error) {

	reroute := &ReRoute{}

	/* create a policy */
	sp, err := s.CreatePolicy(p)
	if err != nil {
		return nil, err
	}

	/* Create a reroute if doen't exist */
	if ip != nil {
		reroute, err = s.CreateReroute(&api.ReRoute{
			Ip: *ip,
		})
		if err != nil {
			return nil, err
		}
	}

	err = s.CreateOrUpdateSubscriber(&api.CreateSubscriber{
		Imsi: imsi,
	}, &(sp.ID), &reroute.ID)
	if err != nil {
		log.Errorf("Failed to create subscriber with imsi %s. Error: %s", imsi, err.Error())
		return nil, err
	}
	// /* check if subscriber exist
	// /* Create subscriber */
	// subscriber := Subscriber{
	// 	Imsi:      imsi,
	// 	PolicyID:  *sp,
	// 	ReRouteID: *reroute,
	// }

	// err = s.InsertSubscriber(&subscriber)
	// if err != nil {
	// 	return nil, err
	// }

	sub, err := s.GetSubscriber(imsi)
	if err != nil {
		log.Errorf("Failed to get subscriber with imsi %s. Error: %s", imsi, err.Error())
		return nil, err
	}

	// Create initial Usage for the subscriber
	err = s.CreateUsage(sub)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

// CRUD operations for Session entity
func (s *Store) InsertSession(se *Session) (*Session, error) {
	res, err := s.db.Exec(`
		INSERT INTO sessions (subscriber_id, apnname, ueipaddr, starttime, endtime , txbytes , rxbytes , totalbytes , txmeter_id, rxmeter_id, state, sync)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?);
		`, se.SubscriberID.ID, se.ApnName, se.UeIpAddr, se.StartTime, se.EndTime, se.TxBytes, se.RxBytes, se.TotalBytes, se.TxMeterID.ID, se.RxMeterID.ID, se.State, se.Sync)
	if err != nil {
		log.Errorf("Failed to insert session.Error %v", err)
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Errorf("Failed to get last inserted session. Error %v", err)
		return nil, err
	}

	ns, err := s.GetSessionByID(int(id))
	if err != nil {
		log.Errorf("Failed to get session. Error %v", err)
		return nil, err
	}
	log.Infof("Session created in store %+v", ns)

	return ns, err
}

func (s *Store) DeleteSession(sub *Subscriber) error {
	_, err := s.db.Exec(`
	DELETE FROM sessions WHERE subscriber_id = ? ;
	`, sub.ID)

	if err != nil {
		log.Errorf("Failed to delete seession for subscriber %d: Error: %v", sub.ID, err.Error())
		return err
	}
	return nil
}

func (s *Store) UpdateSessionUsage(session *Session) error {
	_, err := s.db.Exec(`
		UPDATE sessions
		SET txbytes = ?, rxbytes = ?, totalbytes = ?
		WHERE id = ?;
	`, session.TxBytes, session.RxBytes, (session.TxBytes + session.RxBytes), session.ID)
	return err
}

func (s *Store) UpdateSessionEndUsage(session *Session) error {
	_, err := s.db.Exec(`
		UPDATE sessions
		SET endtime = ?, txbytes = ?, rxbytes = ?, totalbytes = ?, state = ?, sync= ?
		WHERE id = ?;
	`, session.EndTime, session.TxBytes, session.RxBytes, session.TotalBytes, session.State, session.Sync, session.ID)
	return err
}

func (s *Store) CreateSession(subscriber *Subscriber, ueIpAddr string) (*Session, *Flow, *Flow, error) {

	/* TODO: Check if required here vaildate if user has enough data */
	usage, err := s.GetUsageByImsi(subscriber.Imsi)
	if err != nil {
		log.Errorf("Error getting usage for subscriber: %s", err.Error())
		return nil, nil, nil, err
	}

	// Check if Data in Usage is less than Policy for rerouting
	if usage.Data >= subscriber.PolicyID.Data {
		log.Errorf("can't create flows. UE %s has reached max data cap.", subscriber.Imsi)
		return nil, nil, nil, fmt.Errorf("max data cap exceeded")
	}

	/* TODO: start create session
	tx, err := s.db.Begin()
	*/

	rxM, err := s.CreateMeter(subscriber, &subscriber.PolicyID, RX_PATH)
	if err != nil {
		return nil, nil, nil, err
	}

	txM, err := s.CreateMeter(subscriber, &subscriber.PolicyID, TX_PATH)
	if err != nil {
		return nil, nil, nil, err
	}

	session := Session{
		SubscriberID: *subscriber,
		UeIpAddr:     ueIpAddr,
		StartTime:    uint64(time.Now().Unix()), // Current epoch time
		TxMeterID:    *txM,
		RxMeterID:    *rxM,
		State:        SessionActive,
		Sync:         SessionSyncPending,
	}

	// Create Flow for RX
	flowRx := Flow{
		Tableid:  0,
		Priority: 100,
		UeIpAddr: ueIpAddr,
		MeterID:  session.RxMeterID,
	}

	// Create Flow for TX
	flowTx := Flow{
		Tableid:  0,
		Priority: 100,
		UeIpAddr: ueIpAddr,
		MeterID:  session.TxMeterID,
	}

	// Insert Flows
	rxF, err := s.CreateFlow(&flowRx.MeterID, &flowRx.ReRouting, flowRx.UeIpAddr, uint32(flowRx.Tableid), uint32(flowRx.Priority))
	if err != nil {
		return nil, nil, nil, err
	}

	txF, err := s.CreateFlow(&flowTx.MeterID, &flowTx.ReRouting, flowTx.UeIpAddr, uint32(flowTx.Tableid), uint32(flowTx.Priority))
	if err != nil {
		return nil, nil, nil, err
	}

	ns, err := s.InsertSession(&session)
	if err != nil {
		return nil, nil, nil, err
	}

	return ns, rxF, txF, nil
}

func (s *Store) EndSession(session *Session) error {
	// Update session with TX, RX, and Total bytes
	session.EndTime = uint64(time.Now().Unix())
	session.TotalBytes = session.TxBytes + session.RxBytes
	session.State = SessionCompleted
	session.Sync = SessionSyncReady

	// Update Usage for the subscriber
	subscriber, err := s.GetSubscriberByID(session.SubscriberID.ID)
	if err != nil {
		return err
	}

	err = s.UpdateSessionEndUsage(session)
	if err != nil {
		log.Errorf("Error updating session usage for subscriber %s.Error %s", subscriber.Imsi, err.Error())
		return err
	}

	usage, err := s.GetUsageByImsi(subscriber.Imsi)
	if err != nil {
		log.Errorf("Error getting usage for subscriber %s.Error %s", subscriber.Imsi, err.Error())
		return err
	}

	usage.Data += session.TotalBytes

	// Update subscriber and session
	err = s.UpdateUsage(usage)
	if err != nil {
		log.Errorf("Error updating usage for subscriber %s.Error %s", subscriber.Imsi, err.Error())
		return err
	}

	return nil
}

/* Update Usage */
func (s *Store) UpdateUsage(usage *Usage) error {
	_, err := s.db.Exec(`
		UPDATE usages
		SET data = ?, updatedat = ?,
		WHERE id = ?;
	`, usage.Data, usage.Updatedat, usage.ID)
	return err
}

/* Get usage by imsi */
func (s *Store) GetUsageByImsi(imsi string) (*Usage, error) {
	var usage Usage
	err := s.db.QueryRow("SELECT id, subscriber_id, updatedat, data FROM usages WHERE subscriber_id = (SELECT id FROM subscribers WHERE imsi = ?)", imsi).
		Scan(&usage.ID, &usage.SubscriberID.ID, &usage.Updatedat, &usage.Data)
	if err != nil {
		return nil, err
	}

	return &usage, nil
}

func (s *Store) GetPolicyByID(policyID uuid.UUID) (*Policy, error) {
	var policy Policy
	var id []byte
	err := s.db.QueryRow("SELECT id,data,dlbr,ulbr,burst,starttime,endtime FROM policies WHERE id = ?", policyID.Bytes()).
		Scan(&id, &policy.Data, &policy.Dlbr, &policy.Ulbr, &policy.Burst, &policy.StartTime, &policy.EndTime)
	if err != nil {
		return nil, err
	}
	policy.ID, err = uuid.FromBytes(id)
	return &policy, err
}

func (s *Store) GetApplicablePolicyByImsi(imsi string) (*Policy, error) {
	var policy Policy
	var id []byte
	err := s.db.QueryRow(`
		SELECT * FROM policies
		WHERE id = (SELECT policy_id FROM subscribers WHERE imsi = ?)
	`, imsi).
		Scan(&id, &policy.Data, &policy.Dlbr, &policy.Ulbr, &policy.Burst, &policy.StartTime, &policy.EndTime)
	if err != nil {
		return nil, err
	}

	policy.ID, err = uuid.FromBytes(id)
	return &policy, err
}

func (s *Store) GetSessionByID(sessionID int) (*Session, error) {
	var session Session

	err := s.db.QueryRow("SELECT subscriber_id, apnname, ueipaddr, starttime, endtime , txbytes , rxbytes , totalbytes , txmeter_id, rxmeter_id, state, sync FROM sessions WHERE id = ?", sessionID).
		Scan(&session.SubscriberID.ID, &session.ApnName, &session.UeIpAddr, &session.StartTime, &session.EndTime, &session.TxBytes, &session.RxBytes, &session.TotalBytes, &session.TxMeterID.ID, &session.RxMeterID.ID, &session.State, &session.Sync)
	if err != nil {
		return nil, err
	}

	// Fetch associated Subscriber
	sub, err := s.GetSubscriberByID(session.SubscriberID.ID)
	if err != nil {
		return nil, err
	}
	session.SubscriberID = *sub

	// Fetch associated Meters
	txM, err := s.GetMeter(uint32(session.TxMeterID.ID))
	if err != nil {
		return nil, err
	}
	session.TxMeterID = *txM

	rxM, err := s.GetMeter(uint32(session.RxMeterID.ID))
	if err != nil {
		return nil, err
	}
	session.RxMeterID = *rxM

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
		err := rows.Scan(&session.ID, &session.SubscriberID.ID, &session.UeIpAddr, &session.StartTime, &session.EndTime, &session.TxBytes, &session.RxBytes, &session.TotalBytes, &session.TxMeterID.ID, &session.RxMeterID.ID, &session.State)
		if err != nil {
			return nil, err
		}

		// Fetch associated Subscriber
		sub, err := s.GetSubscriberByID(session.SubscriberID.ID)
		if err != nil {
			return nil, err
		}
		session.SubscriberID = *sub

		// Fetch associated Meters
		txM, err := s.GetMeter(uint32(session.TxMeterID.ID))
		if err != nil {
			return nil, err
		}
		session.TxMeterID = *txM

		rxM, err := s.GetMeter(uint32(session.RxMeterID.ID))
		if err != nil {
			return nil, err
		}
		session.RxMeterID = *rxM

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
		Scan(&session.ID, &session.SubscriberID.ID, &session.UeIpAddr, &session.StartTime, &session.EndTime, &session.TxBytes, &session.RxBytes, &session.TotalBytes, &session.TxMeterID.ID, &session.RxMeterID.ID, &session.State)
	if err != nil {
		return nil, err
	}

	// Fetch associated Subscriber
	sub, err := s.GetSubscriberByID(session.SubscriberID.ID)
	if err != nil {
		return nil, err
	}
	session.SubscriberID = *sub

	// Fetch associated Meters
	txM, err := s.GetMeter(uint32(session.TxMeterID.ID))
	if err != nil {
		return nil, err
	}
	session.TxMeterID = *txM

	rxM, err := s.GetMeter(uint32(session.RxMeterID.ID))
	if err != nil {
		return nil, err
	}
	session.RxMeterID = *rxM

	return &session, nil
}

func (s *Store) GetFlowForMeter(id int) (*Flow, error) {
	var f *Flow

	err := s.db.QueryRow(`
		SELECT * FROM flows
		WHERE meter_id = (SELECT id FROM meters WHERE id = ?)
	`, id).
		Scan(&f.ID, &f.Cookie, &f.Tableid, &f.Priority, &f.UeIpAddr, &f.ReRouting.ID, &f.MeterID.ID)
	if err != nil {
		return nil, err
	}

	// Fetch associated Reroute
	route, err := s.GetReRouteByID(f.ReRouting.ID)
	if err != nil {
		return nil, err
	}
	f.ReRouting = *route
	// Fetch associated Meters
	m, err := s.GetMeter(uint32(f.MeterID.ID))
	if err != nil {
		return nil, err
	}
	f.MeterID = *m

	return f, nil
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
		err := rows.Scan(&session.ID, &session.SubscriberID.ID, &session.UeIpAddr, &session.StartTime, &session.EndTime, &session.TxBytes, &session.RxBytes, &session.TotalBytes, &session.TxMeterID.ID, &session.RxMeterID.ID, &session.State)
		if err != nil {
			return nil, err
		}

		// Fetch associated Subscriber
		sub, err := s.GetSubscriberByID(session.SubscriberID.ID)
		if err != nil {
			return nil, err
		}
		session.SubscriberID = *sub

		// Fetch associated Meters
		txM, err := s.GetMeter(uint32(session.TxMeterID.ID))
		if err != nil {
			return nil, err
		}
		session.TxMeterID = *txM

		rxM, err := s.GetMeter(uint32(session.RxMeterID.ID))
		if err != nil {
			return nil, err
		}
		session.RxMeterID = *rxM

		sessions = append(sessions, session)
	}

	return sessions, nil
}

func (s *Store) UpdatePolicy(policy *Policy) error {
	_, err := s.db.Exec(`
		UPDATE policies
		SET data = ?, dlbr = ?, ulbr = ?
		WHERE id = ?; 
		`, policy.Data, policy.Dlbr, policy.Ulbr, policy.ID.Bytes())
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
		SET tableid = ?, priority = ?, ueipaddr = ?, reroute_id = ?, meter_id = ?
		WHERE id = ?;
	`, flow.Tableid, flow.Priority, flow.UeIpAddr, flow.ReRouting.ID, flow.MeterID.ID, flow.ID)
	return err
}

/* CRUD operations for Policy entity */
func (s *Store) InsertPolicy(policy *Policy) error {
	_, err := s.db.Exec(`
		INSERT OR IGNORE INTO policies (id, data, dlbr, ulbr, starttime, endtime, burst)
		VALUES (?, ?, ?, ?, ?, ?, ?);
	`, policy.ID.Bytes(), policy.Data, policy.Dlbr, policy.Ulbr, policy.StartTime, policy.EndTime, policy.Burst)
	return err
}

/* CRUD operations for Reroute entity */
func (s *Store) InsertReRoute(reRoute *ReRoute) error {
	_, err := s.db.Exec(`
		INSERT OR IGNORE INTO reroutes (ipaddr)
		VALUES (?); `, reRoute.IpAddr)
	if err != nil {
		return err
	}
	return err
}

func (s *Store) GetReRouteByID(reRouteID int) (*ReRoute, error) {
	var reRoute ReRoute

	err := s.db.QueryRow("SELECT * FROM reroutes WHERE id = ?", reRouteID).
		Scan(&reRoute.ID, &reRoute.IpAddr)
	if err != nil {
		return nil, err
	}

	return &reRoute, nil
}

func (s *Store) GetReRouteByIP(ip string) (*ReRoute, error) {
	var reRoute ReRoute

	err := s.db.QueryRow("SELECT * FROM reroutes WHERE ipaddr = ?", ip).
		Scan(&reRoute.ID, &reRoute.IpAddr)
	if err != nil {
		return nil, err
	}

	return &reRoute, nil
}

/* Delete a route */
func (s *Store) DeleteReRoute(reRoute *ReRoute) error {
	_, err := s.db.Exec(`
		DELETE FROM reroutes 
		WHERE ipaddr = ?; `, reRoute.IpAddr)
	return err
}

func (s *Store) UpdateReroute(reRoute *ReRoute) error {
	_, err := s.db.Exec(`
		UPDATE reroutes
		SET ipaddr = ?
		WHERE id = ?;
	`, reRoute.IpAddr, reRoute.ID)
	return err
}

/* CRUD operations for Subscriber entity */
func (s *Store) InsertSubscriber(sub *Subscriber) error {
	_, err := s.db.Exec(`
		INSERT OR IGNORE INTO subscribers (imsi, policy_id, reroute_id)
		VALUES (?,?, ?);
	`, sub.Imsi, sub.PolicyID.ID.Bytes(), sub.ReRouteID.ID)
	return err
}

/* Update policy for Subscriber entity */
func (s *Store) UpdateSubscriberPolicy(subscriber *Subscriber, p uuid.UUID) error {
	_, err := s.db.Exec(`
		UPDATE subscribers
		SET policy_id = ?
		WHERE id = ?;
	`, p.Bytes(), subscriber.ID)
	return err
}

/* Update policy for Subscriber entity */
func (s *Store) UpdateSubscriberReRoute(subscriber *Subscriber, id int) error {
	_, err := s.db.Exec(`
		UPDATE subscribers
		SET reroute_id = ?
		WHERE id = ?;
	`, subscriber.ReRouteID.ID, subscriber.ID)
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

	query := "SELECT id, imsi, policy_id, reroute_id FROM subscribers WHERE Imsi = ?"
	row := s.db.QueryRow(query, imsi)

	var subscriber Subscriber
	var id []byte
	err := row.Scan(&subscriber.ID, &subscriber.Imsi, &id, &subscriber.ReRouteID.ID)
	if err != nil {
		// Subscriber not found
		return nil, fmt.Errorf("Subscriber not found: %v", err)
	}
	uuid, err := uuid.FromBytes(id)
	if err != nil {
		return nil, fmt.Errorf("policy id is not a valid uuid: %v", err)
	}

	p, err := s.GetPolicyByID(uuid)
	if err != nil {
		log.Errorf("failed to get policy for subscriber %s.Error: %v", subscriber.Imsi, err)
		return nil, err
	}
	subscriber.PolicyID = *p

	r, err := s.GetReRouteByID(subscriber.ReRouteID.ID)
	if err != nil {
		log.Errorf("failed to get reroute for subscriber %s.Error: %v", subscriber.Imsi, err)
		return nil, err
	}
	subscriber.ReRouteID = *r

	return &subscriber, err
}

func (s *Store) GetSubscriberByID(id int) (*Subscriber, error) {

	query := "SELECT imsi, policy_id, reroute_id FROM subscribers WHERE id = ?"
	row := s.db.QueryRow(query, id)

	var subscriber Subscriber
	var idb []byte
	err := row.Scan(&subscriber.Imsi, &idb, &subscriber.ReRouteID.ID)
	if err != nil {
		// Subscriber not found
		return nil, fmt.Errorf("Subscriber not found: %v", err)
	}

	uuid, err := uuid.FromBytes(idb)
	if err != nil {
		return nil, fmt.Errorf("policy id is not a valid uuid: %v", err)
	}

	p, err := s.GetPolicyByID(uuid)
	if err != nil {
		log.Errorf("failed to get policy for subscriber %s.Error: %v", subscriber.Imsi, err)
		return nil, err
	}
	subscriber.PolicyID = *p

	r, err := s.GetReRouteByID(subscriber.ReRouteID.ID)
	if err != nil {
		log.Errorf("failed to get reroute for subscriber %s.Error: %v", subscriber.Imsi, err)
		return nil, err
	}
	subscriber.ReRouteID = *r

	return &subscriber, err
}

func (s *Store) CreateOrUpdateSubscriber(ns *api.CreateSubscriber, p *uuid.UUID, id *int) error {

	// Check if the subscriber already exists
	subscriber := &Subscriber{}
	err := s.db.QueryRow("SELECT ID FROM Subscriber WHERE Imsi = ?", ns.Imsi).Scan(&subscriber.ID, &subscriber.Imsi)

	if err == nil && subscriber.ID != 0 {
		log.Infof("Subscriber already exists %s. Performing update", subscriber.Imsi)
	} else {
		err := s.InsertSubscriber(&Subscriber{
			Imsi: ns.Imsi,
		})
		if err != nil {
			log.Errorf("Error inserting subscriber %s.Error: %v", subscriber.Imsi, err.Error())
			return err
		}
		log.Infof("New subscriber with Imsi %s created.", ns.Imsi)

		// Get the subscriber
		subscriber, err = s.GetSubscriber(ns.Imsi)
		if err != nil {
			log.Errorf("Erorr while getting sunscriber with imsi %s. Error %s", subscriber.Imsi, err.Error())
			return err
		}
		/* Usage table */
		err = s.CreateUsage(subscriber)
		if err != nil {
			return err
		}
	}

	// Subscriber already exists, update the policy
	if p != nil {
		err := s.UpdateSubscriberPolicy(subscriber, *p)
		if err != nil {
			log.Errorf("Failed to update policy %s for the subscriber %s.Error %s", p.String(), ns.Imsi, err.Error())
			return err
		}
	}

	if id != nil {
		err = s.UpdateSubscriberReRoute(subscriber, *id)
		if err != nil {
			log.Errorf("Failed to update Reroute %d for the subscriber %s.Error %s", *id, ns.Imsi, err.Error())
			return err
		}
	}

	log.Infof("Policy %s assigned to the new subscriber %s.", p, subscriber.Imsi)

	return nil
}

func PrepareCDR(s *Session) *api.CDR {
	cdr := &api.CDR{
		Session:    s.ID,
		Imsi:       s.SubscriberID.Imsi,
		ApnName:    s.ApnName,
		Ip:         s.UeIpAddr,
		StartTime:  s.StartTime,
		EndTime:    s.EndTime,
		TxBytes:    s.TxBytes,
		RxBytes:    s.RxBytes,
		TotalBytes: s.TotalBytes,
	}
	return cdr
}
