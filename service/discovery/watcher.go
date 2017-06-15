package discovery

import (
	"encoding/json"
	"strings"

	"fmt"

	etcd "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc/naming"
)

// EtcdWatcher is the implementaion of grpc.naming.Watcher
type EtcdWatcher struct {
	prefix  string
	srvName string
	client  *etcd.Client
	kapi    etcd.KeysAPI
	addrs   []string
}

// Close do nothing
func (ew *EtcdWatcher) Close() {
}

// Next to return the updates
func (ew *EtcdWatcher) Next() ([]*naming.Update, error) {
	dirKey := getDir(ew.prefix, ew.srvName)

	// ew.addrs is nil means it is intially called
	// first time query return instantlly
	if ew.addrs == nil {
		// query addresses from etcd
		resp, _ := ew.kapi.Get(context.Background(), dirKey, &etcd.GetOptions{Recursive: true})
		addrs := ew.extractAddrs(resp)
		ew.addrs = addrs

		if addrs != nil {
			return getUpdates([]string{}, addrs), nil
		}
	}

	// after first time, use Watcher to block until changes happen
	// generate etcd Watcher
	w := ew.kapi.Watcher(dirKey, &etcd.WatcherOptions{Recursive: true})
	for {
		_, err := w.Next(context.Background())
		if err == nil {
			// query addresses from etcd
			resp, err := ew.kapi.Get(context.Background(), dirKey, &etcd.GetOptions{Recursive: true})
			if err != nil {
				continue
			}
			addrs := ew.extractAddrs(resp)

			updates := getUpdates(ew.addrs, addrs)
			// update ew.addrs
			ew.addrs = addrs
			// if addrs updated, return it
			if len(updates) > 0 {
				return updates, nil
			}
		}
	}
}

// helper function to extract addrs rom etcd response
func (ew *EtcdWatcher) extractAddrs(resp *etcd.Response) (addrs []string) {
	if resp == nil || resp.Node == nil || resp.Node.Nodes == nil {
		return nil
	}

	if len(resp.Node.Nodes) == 0 {
		ew.dropEmptyDir(resp.Node.Key)
		return nil
	}

	for _, node := range resp.Node.Nodes {
		if !node.Dir && !stringEndsWith(node.Key, "ID") {
			service := new(Service)
			if err := json.Unmarshal([]byte(node.Value), service); err == nil {
				addrs = append(addrs, fmt.Sprintf("%s:%d", service.Host, service.Port))
			}
		} else {

		}
	}

	return
}

func (ew *EtcdWatcher) dropEmptyDir(emptyDir string) {
	if len(emptyDir) == 0 {
		return
	}

	_, err := ew.kapi.Delete(context.Background(), emptyDir, &etcd.DeleteOptions{Dir: false})
	if err != nil {
	}
}

func getUpdates(old, new []string) []*naming.Update {
	var updates []*naming.Update

	intersect := intersectString(old, new)

	deleted := diffString(old, intersect)
	for _, addr := range deleted {
		update := &naming.Update{Op: naming.Delete, Addr: addr}
		updates = append(updates, update)
	}

	added := diffString(new, intersect)
	for _, addr := range added {
		update := &naming.Update{Op: naming.Add, Addr: addr}
		updates = append(updates, update)
	}
	return updates
}

func intersectString(a, b []string) (i []string) {
	for _, va := range a {
		found := false
		for _, vb := range b {
			if va == vb {
				found = true
			}
		}
		if found {
			i = append(i, va)
		}
	}

	return
}

// a has, b not have
func diffString(a, b []string) (d []string) {
	for _, va := range a {
		found := false
		for _, vb := range b {
			if va == vb {
				found = true
			}
		}

		if !found {
			d = append(d, va)
		}
	}

	return
}

func stringEndsWith(haystack, needle string) bool {
	return strings.LastIndex(haystack, needle) == len(haystack)-len(needle)
}
