package main

import (
	"context"
	"etcd-client/pkg"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
	"time"
)

func main() {
	endpoints := []string{"localhost:2379"}

	config := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}

	client, err := clientv3.New(config)
	if err != nil {
		panic(err)
	}

	store(client, "hosts")
	store(client, "auth")
	store(client, "github")

	ctx := context.Background()

	if err != nil {
		panic(err)
	}

	res, err := client.Get(ctx, "/config/github.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Print(res.Kvs)

	defer func(client *clientv3.Client) {
		err := client.Close()
		if err != nil {
			fmt.Printf("Error closing etcd client: %v", err)
			os.Exit(1)
		}
	}(client)

	l := pkg.NewConfigLock(client, "test", "test")

	if err := l.Lock(ctx, 10); err != nil {
		panic(err)
	}
}

func store(c *clientv3.Client, name string) {
	f, err := os.ReadFile("envs/" + name + ".yaml")

	if err != nil {
		panic(err)
	}

	_, err = c.Put(context.Background(), "/config/"+name+".yaml", string(f))
}
