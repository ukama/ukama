package pkg

import (
	"context"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

const ETCD_TIMEOUT = 3

type EtcdStorage struct {
	etcd *clientv3.Client
}

func NewEtcdStorage(etcdHost string) *EtcdStorage {
	client, err := clientv3.New(clientv3.Config{
		DialTimeout: 5 * time.Second,
		Endpoints:   []string{etcdHost},
	})
	if err != nil {
		logrus.Fatalf("Cannot connect to etcd: %v", err)
	}

	return &EtcdStorage{
		etcd: client,
	}
}

func (e EtcdStorage) Get(key string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), ETCD_TIMEOUT*time.Second)
	defer cancel()
	val, err := e.etcd.Get(ctx, key, clientv3.WithLimit(1))
	if err != nil {
		logrus.Errorf("Cannot get sim info from etcd: %v", err)
		return nil, err
	}
	if len(val.Kvs) == 0 {
		return nil, nil
	}

	return val.Kvs[0].Value, nil
}

func (e EtcdStorage) Put(key string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ETCD_TIMEOUT*time.Second)
	defer cancel()
	_, err := e.etcd.Put(ctx, key, value)
	return err
}

func (e EtcdStorage) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ETCD_TIMEOUT*time.Second)
	defer cancel()
	_, err := e.etcd.Delete(ctx, key)
	return err
}
