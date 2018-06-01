package factory

import (
	"kis/apiserver/storage"
	"kis/apiserver/storage/storagebackend"
	"github.com/coreos/etcd/clientv3"
	"time"
)


func CreateStorageBackend(c *storagebackend.Config) (storage.Interface,storagebackend.DestroyFunc,error){
	cfg := clientv3.Config{
		DialTimeout:          10 * time.Second,
		DialKeepAliveTime:    30 * time.Second,
		DialKeepAliveTimeout: 10 * time.Second,
		Endpoints:            c.ServerList,
	}
	client, err := clientv3.New(cfg)
	if err != nil {
		return nil, nil, err
	}
	//ctx, cancel := context.WithCancel(context.Background())
	//etcd3.StartCompactor(ctx, client, c.CompactionInterval)
	destroyFunc := func() {
		//cancel()
		client.Close()
	}

	return storage.NewStore(client, c.Prefix), destroyFunc, nil
}

