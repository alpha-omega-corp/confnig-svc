package pkg

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

type ConfigLock interface {
	Lock(ctx context.Context, ttl int64) (err error)
}

type configLock struct {
	ConfigLock

	Key        string
	Value      string
	LeaseID    clientv3.LeaseID
	etcdClient *clientv3.Client
}

func NewConfigLock(etcdClient *clientv3.Client, key string, value string) ConfigLock {
	return &configLock{
		Key:        key,
		Value:      value,
		etcdClient: etcdClient,
	}
}

func (h *configLock) Lock(ctx context.Context, ttl int64) (err error) {
	lease, err := h.etcdClient.Grant(ctx, ttl)
	if err != nil {
		return
	}

	_, err = h.etcdClient.Put(ctx, h.Key, h.Value, clientv3.WithLease(lease.ID))
	if err != nil {
		return
	}

	h.LeaseID = lease.ID

	return
}

func (h *configLock) Unlock(ctx context.Context) (err error) {
	_, err = h.etcdClient.Delete(ctx, h.Key)
	if err != nil {
		return
	}

	_, err = h.etcdClient.Revoke(ctx, h.LeaseID)
	if err != nil {
		return
	}

	log.Printf("Lock released: %s", h.Key)
	return
}
