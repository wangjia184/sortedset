package sortedset

import (
	"testing"
)

func checkOrder(t *testing.T, nodes []*SortedSetNode, expectedOrder []string) {
	if len(expectedOrder) != len(nodes) {
		t.Errorf("nodes does not contain %d elements", len(expectedOrder))
	}
	for i := 0; i < len(expectedOrder); i++ {
		if nodes[i].Key() != expectedOrder[i] {
			t.Errorf("nodes[%d] is %q, but the expected key is %q", i, nodes[i].Key(), expectedOrder[i])
		}

	}
}

func checkIterByRankRange(t *testing.T, sortedset *SortedSet, start int, end int, expectedOrder []string) {
	var keys []string

	// check nil callback should do nothing
	sortedset.IterFuncByRankRange(start, end, nil)

	sortedset.IterFuncByRankRange(start, end, func(key string, _ interface{}) bool {
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
	// reset data
	keys = []string{}
	var i int
	sortedset.IterFuncByRankRange(start, end, func(key string, _ interface{}) bool {
		keys = append(keys, key)
		i++
		// return early
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

func checkRankRangeIterAndOrder(t *testing.T, sortedset *SortedSet, start int, end int, remove bool, expectedOrder []string) {
	checkIterByRankRange(t, sortedset, start, end, expectedOrder)
	nodes := sortedset.GetByRankRange(start, end, remove)
	checkOrder(t, nodes, expectedOrder)
}

func TestCase1(t *testing.T) {
	sortedset := New()

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

	// create a new set
	sortedset := New()

	// fill in new node
	sortedset.AddOrUpdate("a", 89, "Kelly")
	sortedset.AddOrUpdate("b", 100, "Staley")
	sortedset.AddOrUpdate("c", 100, "Jordon")
	sortedset.AddOrUpdate("d", -321, "Park")
	sortedset.AddOrUpdate("e", 101, "Albert")
	sortedset.AddOrUpdate("f", 99, "Lyman")
	sortedset.AddOrUpdate("g", 99, "Singleton")
	sortedset.AddOrUpdate("h", 70, "Audrey")

	// update an existing node
	sortedset.AddOrUpdate("e", 99, "ntrnrt")

	// remove node
	sortedset.Remove("b")

	nodes := sortedset.GetByScoreRange(-500, 500, nil)
	checkOrder(t, nodes, []string{"d", "h", "a", "e", "f", "g", "c"})

	nodes = sortedset.GetByScoreRange(500, -500, nil)
	//t.Logf("%v", nodes)
	checkOrder(t, nodes, []string{"c", "g", "f", "e", "a", "h", "d"})

	nodes = sortedset.GetByScoreRange(600, 500, nil)
	checkOrder(t, nodes, []string{})

	nodes = sortedset.GetByScoreRange(500, 600, nil)
	checkOrder(t, nodes, []string{})

	rank := sortedset.FindRank("f")
	if rank != 5 {
		t.Error("FindRank() does not return expected value `5`")
	}

	rank = sortedset.FindRank("d")
	if rank != 1 {
		t.Error("FindRank() does not return expected value `1`")
	}

	nodes = sortedset.GetByScoreRange(99, 100, nil)
	checkOrder(t, nodes, []string{"e", "f", "g", "c"})

	nodes = sortedset.GetByScoreRange(90, 50, nil)
	checkOrder(t, nodes, []string{"a", "h"})

	nodes = sortedset.GetByScoreRange(99, 100, &GetByScoreRangeOptions{
		ExcludeStart: true,
	})
	checkOrder(t, nodes, []string{"c"})

	nodes = sortedset.GetByScoreRange(100, 99, &GetByScoreRangeOptions{
		ExcludeStart: true,
	})
	checkOrder(t, nodes, []string{"g", "f", "e"})

	nodes = sortedset.GetByScoreRange(99, 100, &GetByScoreRangeOptions{
		ExcludeEnd: true,
	})
	checkOrder(t, nodes, []string{"e", "f", "g"})

	nodes = sortedset.GetByScoreRange(100, 99, &GetByScoreRangeOptions{
		ExcludeEnd: true,
	})
	checkOrder(t, nodes, []string{"c"})

	nodes = sortedset.GetByScoreRange(50, 100, &GetByScoreRangeOptions{
		Limit: 2,
	})
	checkOrder(t, nodes, []string{"h", "a"})

	nodes = sortedset.GetByScoreRange(100, 50, &GetByScoreRangeOptions{
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

	nodes = sortedset.GetByScoreRange(-500, 500, nil)
	checkOrder(t, nodes, []string{"h", "a", "e", "f", "g", "c"})

	maxNode := sortedset.PeekMax()
	if maxNode == nil || maxNode.Key() != "c" {
		t.Error("PeekMax() does not return expected value `c`")
	}

	maxNode = sortedset.PopMax()
	if maxNode == nil || maxNode.Key() != "c" {
		t.Error("PopMax() does not return expected value `c`")
	}

	nodes = sortedset.GetByScoreRange(500, -500, nil)
	checkOrder(t, nodes, []string{"g", "f", "e", "a", "h"})

	nodes = sortedset.GetByScoreRange(69, 89, &GetByScoreRangeOptions{
		Remove: true,
	})
	checkOrder(t, nodes, []string{"h", "a"})

	nodes = sortedset.GetByScoreRange(500, -500, nil)
	checkOrder(t, nodes, []string{"g", "f", "e"})

}
