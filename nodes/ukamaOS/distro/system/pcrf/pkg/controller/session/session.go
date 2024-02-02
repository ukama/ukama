package session

import "time"

type session struct {
	period time.Duration `default:"1s"`
	imsi string
	store *store.Store
	 
}
