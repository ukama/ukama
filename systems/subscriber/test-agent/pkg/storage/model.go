package storage

type SimInfo struct {
	Iccid string `json:"iccid"`
	Imsi  string `json:"imsi"`
}

type Storage interface {
	Get(key string) ([]byte, error)
	Put(key string, value string) error
	Delete(key string) error
}
