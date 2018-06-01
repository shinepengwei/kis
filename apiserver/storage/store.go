package storage

import (
	"github.com/coreos/etcd/clientv3"
	"path"
	"fmt"
	"context"
	"kis/apiserver/runtime"
	"encoding/json"
	"reflect"
	"kis/apis/meta"
)

func NewStore(c *clientv3.Client, prefix string) Interface{
	result := &store{
		client:        c,
		pathPrefix: path.Join("/", prefix),
	}
	return result
}

//这个store持有一个etcd的client，是registry/generic/store中Store的BackendStore
type store struct {
	client *clientv3.Client
	// getOpts contains additional options that should be passed
	// to all Get() calls.
	getOps        []clientv3.OpOption
	pathPrefix    string
}


// Get implements storage.Interface.Get.
func (s *store) Get(ctx context.Context, key string, resourceVersion string, out runtime.Object, ignoreNotFound bool) error {
	key = path.Join(s.pathPrefix, key)
	getResp, err := s.client.KV.Get(ctx, key, s.getOps...)
	if err != nil {
		return err
	}
	kv := getResp.Kvs[0]
	//TODO 序列化反序列化
	fmt.Println(kv.Value)

	return nil
}

// Create implements storage.Interface.Create.
func (s *store) Create(ctx context.Context, key string, obj, out runtime.Object) error {
	key = path.Join(s.pathPrefix, key)
	fmt.Println(key)
	newData, _ := json.Marshal(obj)
	fmt.Println("create "+key+", store into etcd:"+string(newData))
	txnResp, _ := s.client.KV.Txn(ctx).If(
		notFound(key),
	).Then(
		clientv3.OpPut(key, string(newData)),
	).Commit()
	fmt.Print("stored success:")
	fmt.Println(txnResp)
	return nil
}

// Create implements storage.Interface.List.
func (s *store) List(ctx context.Context, key string, resourceVersion string, listObj runtime.Object) error{
	key = path.Join(s.pathPrefix, key)
	fmt.Println("get from etcd:"+key)
	options := make([]clientv3.OpOption, 0, 4)
	options = append(options, clientv3.WithPrefix())
	getResp, err := s.client.KV.Get(ctx, key,options...)
	fmt.Println("List "+key+":")
	fmt.Println(getResp)
	listPtr, _:= meta.GetItemPtr(listObj)
	v := reflect.ValueOf(listPtr).Elem()
	for _, kv := range getResp.Kvs {
		fmt.Println(string(kv.Key)+","+string(kv.Value))
		appendListItem(v,kv.Value)
	}
	fmt.Println(listObj)
	return err
}

func notFound(key string) clientv3.Cmp {
	return clientv3.Compare(clientv3.ModRevision(key), "=", 0)
}


//v是PodList中的Item的反射Value值
func appendListItem(v reflect.Value, data []byte){
	obj, _:= reflect.New(v.Type().Elem()).Interface().(runtime.Object)
	json.Unmarshal(data,obj)
	v.Set(reflect.Append(v, reflect.ValueOf(obj).Elem()))
}