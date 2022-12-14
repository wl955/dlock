package dlock

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/wl955/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

var _opts clientv3.Config

var client *clientv3.Client

func Init(opts ...Option) {
	custom := Options{}

	for _, o := range opts {
		o.apply(&custom)
	}

	var err error

	client, err = clientv3.New(
		_opts,
	)
	if err != nil {
		log.Fatalf("failed to create an etcd client: %s\n", err)
	}

	log.Info("creted etcd client\n")
}

func NewLocker(pfx string, ttl int) (locker sync.Locker, err error) {
	if nil == client {
		return nil, errors.New("init first")
	}
	// WithTTL configures the session's TTL in seconds.
	session, err := concurrency.NewSession(client, concurrency.WithTTL(ttl))
	if err != nil {
		return
	}

	locker = concurrency.NewLocker(session, pfx)

	return
}

const (
	// These const values might be need adjustment.
	nrGarbageObjects = 100 * 1000 * 1000
	sessionTTL       = 1
)

func test() {
	session, err := concurrency.NewSession(client, concurrency.WithTTL(sessionTTL))
	if err != nil {
		fmt.Printf("failed to create a session: %s\n", err)
		os.Exit(1)
	}

	locker := concurrency.NewLocker(session, "/lock")
	locker.Lock()
	defer locker.Unlock()

	version := session.Lease()
	fmt.Printf("acquired lock, version: %d\n", version)
}
