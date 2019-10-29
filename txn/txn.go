package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	etcd "github.com/coreos/etcd/clientv3"
)

func newClient() *clientv3.Client {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func main() {
	ctx := context.Background()

	client := newClient()
	defer client.Close()

	if _, err := client.Put(ctx, "foo", "bar"); err != nil {
		log.Fatal(err)
	}

	_, err := client.Txn(ctx).If(
		etcd.Compare(etcd.Value("foo"), "=", "bar"),
	).Then(
		etcd.OpPut("frob", "true"),
	).Else(
		etcd.OpPut("frob", "false"),
	).Commit()

	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Get(ctx, "frob")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(resp.Kvs[0].Value))
}
