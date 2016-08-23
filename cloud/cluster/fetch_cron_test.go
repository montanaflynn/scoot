package cluster_test

import (
	"github.com/scootdev/scoot/cloud/cluster"
	"reflect"
	"sort"
	"sync"
	"testing"
	"time"
)

func TestFetchCron(t *testing.T) {
	h := makeCronHelper(t)
	h.assertFetch(t, h.setupTest())
	h.assertFetch(t, h.setupTest("host1:1234"))
	h.assertFetch(t, h.setupTest("host1:1234", "host2:8888"))
	h.assertFetch(t, h.setupTest("host1:1234"))
	h.assertFetch(t, h.setupTest())
}

func (h *cronHelper) setupTest(nodeNames ...string) []cluster.Node {
	nodes := nodes(nodeNames)
	sort.Sort(cluster.NodeSorter(nodes))
	h.f.setResult(nodes)
	time.Sleep(100000 * time.Nanosecond)
	return nodes
}

type cronHelper struct {
	t    *testing.T
	cl   *cluster.Cluster
	time time.Duration
	f    *fakeFetcher
	ch   chan []cluster.Node
	c    *cluster.FetchCron
}

func makeCronHelper(t *testing.T) *cronHelper {
	h := &cronHelper{t: t, time: time.Nanosecond}
	h.cl = cluster.NewCluster([]cluster.Node{}, make(chan interface{}))
	h.f = &fakeFetcher{}
	h.ch = make(chan []cluster.Node)
	h.c = cluster.NewFetchCron(h.f, h.time, h.ch)
	return h
}

func (h *cronHelper) assertFetch(t *testing.T, expected []cluster.Node) {
	actual := <-h.ch
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Unequal, expected %v, received %v", expected, actual)
	}
}

func nodes(ids []string) []cluster.Node {
	n := []cluster.Node{}
	for _, name := range ids {
		n = append(n, cluster.NewIdNode(name))
	}
	return n
}

// fakeFetcher for testing fetch cron
type fakeFetcher struct {
	mutex sync.Mutex
	nodes []cluster.Node
}

func (f *fakeFetcher) Fetch() ([]cluster.Node, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	return f.nodes, nil
}

func (f *fakeFetcher) setResult(nodes []cluster.Node) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.nodes = nodes
}