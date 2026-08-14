package sortedset

import (
	"sync"
	"testing"
)

func checkOrder(t *testing.T, nodes []*SortedSetNode[string, int64, string], expectedOrder []string) {
	if len(expectedOrder) != len(nodes) {
		t.Errorf("nodes does not contain %d elements", len(expectedOrder))
	}
	for i := 0; i < len(expectedOrder); i++ {
		if nodes[i].Key() != expectedOrder[i] {
			t.Errorf("nodes[%d] is %q, but the expected key is %q", i, nodes[i].Key(), expectedOrder[i])
		}
	}
}

func checkIterByRankRange(t *testing.T, sortedset *SortedSet[string, int64, string], start int, end int, expectedOrder []string) {
	var keys []string

	// check nil callback should do nothing
	sortedset.IterFuncRangeByRank(start, end, nil)

	sortedset.IterFuncRangeByRank(start, end, func(key string, _ string) bool {
		keys = append(keys, key)
		return true
	})
	if len(expectedOrder) != len(keys) {
		t.Errorf("keys does not contain %d elements", len(expectedOrder))
	}
	for i := 0; i < len(expectedOrder); i++ {
		if keys[i] != expectedOrder[i] {
			t.Errorf("keys[%d] is %q, but the expected key is %q", i, keys[i], expectedOrder[i])
		}
	}

	// check return early
	if len(expectedOrder) < 1 {
		return
	}
	keys = []string{}
	var i int
	sortedset.IterFuncRangeByRank(start, end, func(key string, _ string) bool {
		keys = append(keys, key)
		i++
		return i < len(expectedOrder)-1
	})
	if len(expectedOrder)-1 != len(keys) {
		t.Errorf("keys does not contain %d elements", len(expectedOrder)-1)
	}
	for i := 0; i < len(expectedOrder)-1; i++ {
		if keys[i] != expectedOrder[i] {
			t.Errorf("keys[%d] is %q, but the expected key is %q", i, keys[i], expectedOrder[i])
		}
	}
}

func checkRankRangeIterAndOrder(t *testing.T, sortedset *SortedSet[string, int64, string], start int, end int, remove bool, expectedOrder []string) {
	checkIterByRankRange(t, sortedset, start, end, expectedOrder)
	nodes := sortedset.GetRangeByRank(start, end, remove)
	checkOrder(t, nodes, expectedOrder)
}

func TestCase1(t *testing.T) {
	sortedset := New[string, int64, string]()

	sortedset.AddOrUpdate("a", 89, "Kelly")
	sortedset.AddOrUpdate("b", 100, "Staley")
	sortedset.AddOrUpdate("c", 100, "Jordon")
	sortedset.AddOrUpdate("d", -321, "Park")
	sortedset.AddOrUpdate("e", 101, "Albert")
	sortedset.AddOrUpdate("f", 99, "Lyman")
	sortedset.AddOrUpdate("g", 99, "Singleton")
	sortedset.AddOrUpdate("h", 70, "Audrey")

	sortedset.AddOrUpdate("e", 99, "ntrnrt")

	sortedset.Remove("b")

	node := sortedset.GetByRank(3, false)
	if node == nil || node.Key() != "a" {
		t.Error("GetByRank() does not return expected value `a`")
	}

	node = sortedset.GetByRank(-3, false)
	if node == nil || node.Key() != "f" {
		t.Error("GetByRank() does not return expected value `f`")
	}

	// get all nodes since the first one to last one
	checkRankRangeIterAndOrder(t, sortedset, 1, -1, false, []string{"d", "h", "a", "e", "f", "g", "c"})

	// get & remove the 2nd/3rd nodes in reserve order
	checkRankRangeIterAndOrder(t, sortedset, -2, -3, true, []string{"g", "f"})

	// get all nodes since the last one to first one
	checkRankRangeIterAndOrder(t, sortedset, -1, 1, false, []string{"c", "e", "a", "h", "d"})
}

func TestCase2(t *testing.T) {
	sortedset := New[string, int64, string]()

	sortedset.AddOrUpdate("a", 89, "Kelly")
	sortedset.AddOrUpdate("b", 100, "Staley")
	sortedset.AddOrUpdate("c", 100, "Jordon")
	sortedset.AddOrUpdate("d", -321, "Park")
	sortedset.AddOrUpdate("e", 101, "Albert")
	sortedset.AddOrUpdate("f", 99, "Lyman")
	sortedset.AddOrUpdate("g", 99, "Singleton")
	sortedset.AddOrUpdate("h", 70, "Audrey")

	sortedset.AddOrUpdate("e", 99, "ntrnrt")

	sortedset.Remove("b")

	nodes := sortedset.GetRangeByScore(-500, 500, nil)
	checkOrder(t, nodes, []string{"d", "h", "a", "e", "f", "g", "c"})

	nodes = sortedset.GetRangeByScore(500, -500, nil)
	checkOrder(t, nodes, []string{"c", "g", "f", "e", "a", "h", "d"})

	nodes = sortedset.GetRangeByScore(600, 500, nil)
	checkOrder(t, nodes, []string{})

	nodes = sortedset.GetRangeByScore(500, 600, nil)
	checkOrder(t, nodes, []string{})

	rank := sortedset.FindRank("f")
	if rank != 5 {
		t.Error("FindRank() does not return expected value `5`")
	}

	rank = sortedset.FindRank("d")
	if rank != 1 {
		t.Error("FindRank() does not return expected value `1`")
	}

	nodes = sortedset.GetRangeByScore(99, 100, nil)
	checkOrder(t, nodes, []string{"e", "f", "g", "c"})

	nodes = sortedset.GetRangeByScore(90, 50, nil)
	checkOrder(t, nodes, []string{"a", "h"})

	nodes = sortedset.GetRangeByScore(99, 100, &GetRangeByScoreOptions{
		ExcludeStart: true,
	})
	checkOrder(t, nodes, []string{"c"})

	nodes = sortedset.GetRangeByScore(100, 99, &GetRangeByScoreOptions{
		ExcludeStart: true,
	})
	checkOrder(t, nodes, []string{"g", "f", "e"})

	nodes = sortedset.GetRangeByScore(99, 100, &GetRangeByScoreOptions{
		ExcludeEnd: true,
	})
	checkOrder(t, nodes, []string{"e", "f", "g"})

	nodes = sortedset.GetRangeByScore(100, 99, &GetRangeByScoreOptions{
		ExcludeEnd: true,
	})
	checkOrder(t, nodes, []string{"c"})

	nodes = sortedset.GetRangeByScore(50, 100, &GetRangeByScoreOptions{
		Limit: 2,
	})
	checkOrder(t, nodes, []string{"h", "a"})

	nodes = sortedset.GetRangeByScore(100, 50, &GetRangeByScoreOptions{
		Limit: 2,
	})
	checkOrder(t, nodes, []string{"c", "g"})

	minNode := sortedset.PeekMin()
	if minNode == nil || minNode.Key() != "d" {
		t.Error("PeekMin() does not return expected value `d`")
	}

	minNode = sortedset.PopMin()
	if minNode == nil || minNode.Key() != "d" {
		t.Error("PopMin() does not return expected value `d`")
	}

	nodes = sortedset.GetRangeByScore(-500, 500, nil)
	checkOrder(t, nodes, []string{"h", "a", "e", "f", "g", "c"})

	maxNode := sortedset.PeekMax()
	if maxNode == nil || maxNode.Key() != "c" {
		t.Error("PeekMax() does not return expected value `c`")
	}

	maxNode = sortedset.PopMax()
	if maxNode == nil || maxNode.Key() != "c" {
		t.Error("PopMax() does not return expected value `c`")
	}

	nodes = sortedset.GetRangeByScore(500, -500, nil)
	checkOrder(t, nodes, []string{"g", "f", "e", "a", "h"})
}

func TestHas(t *testing.T) {
	sortedset := New[string, int64, string]()
	if sortedset.Has("a") {
		t.Fatal("Has() true for an empty set")
	}

	sortedset.AddOrUpdate("a", 10, "v-a")
	sortedset.AddOrUpdate("b", 20, "v-b")
	if !sortedset.Has("a") || !sortedset.Has("b") {
		t.Fatal("Has() false after AddOrUpdate")
	}
	if sortedset.Has("missing") {
		t.Fatal("Has() true for an unknown key")
	}

	sortedset.Remove("a")
	if sortedset.Has("a") {
		t.Fatal("Has() still true after Remove")
	}
	if !sortedset.Has("b") {
		t.Fatal("Has() false for the remaining key")
	}
}

func TestConcurrentReaders(t *testing.T) {
	sortedset := New[int64, int64, int64]()

	var wg sync.WaitGroup
	stop := make(chan struct{})

	// writer mutates the set on one goroutine.
	wg.Add(1)
	go func() {
		defer wg.Done()
		i := int64(0)
		for {
			select {
			case <-stop:
				return
			default:
			}
			key := i % 1000
			sortedset.AddOrUpdate(key, i, i)
			if i%7 == 0 {
				sortedset.Remove(key)
			}
			i++
		}
	}()

	// concurrent readers query membership via Has (the only concurrency-safe
	// method) without racing the writer.
	for r := 0; r < 4; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				_ = sortedset.Has(int64(r))
			}
		}()
	}

	for i := 0; i < 50; i++ {
		_ = sortedset.Has(int64(i))
	}
	close(stop)
	wg.Wait()
}
