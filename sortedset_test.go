package sortedset

import (
	"fmt"
	"testing"
)

func checkOrder(t *testing.T, nodes []*SortedSetNode, expectedOrder []int32) {
	if len(expectedOrder) != len(nodes) {
		t.Errorf("nodes does not contain %d elements", len(expectedOrder))
	}
	for i := 0; i < len(expectedOrder); i++ {
		if nodes[i].Key() != expectedOrder[i] {
			fmt.Println(nodes[i].Key(), expectedOrder[i])
			t.Errorf("nodes[%d] is %q, but the expected key is %q", i, nodes[i].Key(), expectedOrder[i])
		}
	}
}

func checkIterByRankRange(t *testing.T, sortedset *SortedSet, start int, end int, expectedOrder []int32) {
	var keys []int32

	// check nil callback should do nothing
	sortedset.IterFuncByRankRange(start, end, nil)

	sortedset.IterFuncByRankRange(start, end, func(key int32, _ interface{}) bool {
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
	keys = []int32{}
	var i int
	sortedset.IterFuncByRankRange(start, end, func(key int32, _ interface{}) bool {
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

func checkRankRangeIterAndOrder(t *testing.T, sortedset *SortedSet, start int, end int, remove bool, expectedOrder []int32) {
	checkIterByRankRange(t, sortedset, start, end, expectedOrder)
	nodes := sortedset.GetByRankRange(start, end, remove)
	checkOrder(t, nodes, expectedOrder)
}

func TestCase1(t *testing.T) {
	sortedset := New()

	sortedset.AddOrUpdate(12, 89, "Kelly")
	sortedset.AddOrUpdate(22, 100, "Staley")
	sortedset.AddOrUpdate(33, 100, "Jordon")
	sortedset.AddOrUpdate(1000, -321, "Park")
	sortedset.AddOrUpdate(1000111, 101, "Albert")
	sortedset.AddOrUpdate(10001112, 99, "Lyman")
	sortedset.AddOrUpdate(10001113, 99, "Singleton")
	sortedset.AddOrUpdate(10001114, 70, "Audrey")

	sortedset.AddOrUpdate(1000111, 99, "ntrnrt")

	sortedset.Remove(22)

	node := sortedset.GetByRank(3, false)
	if node == nil || node.Key() != 12 {
		t.Error("GetByRank() does not return expected value `12`")
	}

	node = sortedset.GetByRank(-3, false)
	if node == nil || node.Key() != 10001112 {
		t.Error("GetByRank() does not return expected value `f`")
	}

	// get all nodes since the first one to last one
	checkRankRangeIterAndOrder(t, sortedset, 1, -1, false, []int32{1000, 10001114, 12, 1000111, 10001112, 10001113, 33})

	// get & remove the 2nd/3rd nodes in reserve order
	checkRankRangeIterAndOrder(t, sortedset, -2, -3, true, []int32{10001113, 10001112})

	// get all nodes since the last one to first one
	checkRankRangeIterAndOrder(t, sortedset, -1, 1, false, []int32{33, 1000111, 12, 10001114, 1000})

}

func TestCase2(t *testing.T) {

	// create a new set
	sortedset := New()

	// fill in new node
	sortedset.AddOrUpdate(12, 89, "Kelly")
	sortedset.AddOrUpdate(22, 100, "Staley")
	sortedset.AddOrUpdate(33, 100, "Jordon")
	sortedset.AddOrUpdate(1000, -321, "Park")
	sortedset.AddOrUpdate(1000111, 101, "Albert")
	sortedset.AddOrUpdate(10001112, 99, "Lyman")
	sortedset.AddOrUpdate(10001113, 99, "Singleton")
	sortedset.AddOrUpdate(10001114, 70, "Audrey")

	// update an existing node
	sortedset.AddOrUpdate(1000111, 99, "ntrnrt")

	// remove node
	sortedset.Remove(22)

	nodes := sortedset.GetByScoreRange(-500, 500, nil)
	checkOrder(t, nodes, []int32{1000, 10001114, 12, 1000111, 10001112, 10001113, 33})

	nodes = sortedset.GetByScoreRange(500, -500, nil)
	//t.Logf("%v", nodes)
	checkOrder(t, nodes, []int32{33, 10001113, 10001112, 1000111, 12, 10001114, 1000})

	nodes = sortedset.GetByScoreRange(600, 500, nil)
	checkOrder(t, nodes, []int32{})

	nodes = sortedset.GetByScoreRange(500, 600, nil)
	checkOrder(t, nodes, []int32{})

	rank := sortedset.FindRank(10001112)
	if rank != 5 {
		t.Error("FindRank() does not return expected value `5`")
	}

	rank = sortedset.FindRank(1000)
	if rank != 1 {
		t.Error("FindRank() does not return expected value `1`")
	}

	nodes = sortedset.GetByScoreRange(99, 100, nil)
	checkOrder(t, nodes, []int32{1000111, 10001112, 10001113, 33})

	nodes = sortedset.GetByScoreRange(90, 50, nil)
	checkOrder(t, nodes, []int32{12, 10001114})

	nodes = sortedset.GetByScoreRange(99, 100, &GetByScoreRangeOptions{
		ExcludeStart: true,
	})
	checkOrder(t, nodes, []int32{33})

	nodes = sortedset.GetByScoreRange(100, 99, &GetByScoreRangeOptions{
		ExcludeStart: true,
	})
	checkOrder(t, nodes, []int32{10001113, 10001112, 1000111})

	nodes = sortedset.GetByScoreRange(99, 100, &GetByScoreRangeOptions{
		ExcludeEnd: true,
	})
	checkOrder(t, nodes, []int32{1000111, 10001112, 10001113})

	nodes = sortedset.GetByScoreRange(100, 99, &GetByScoreRangeOptions{
		ExcludeEnd: true,
	})
	checkOrder(t, nodes, []int32{33})

	nodes = sortedset.GetByScoreRange(50, 100, &GetByScoreRangeOptions{
		Limit: 2,
	})
	checkOrder(t, nodes, []int32{10001114, 12})

	nodes = sortedset.GetByScoreRange(100, 50, &GetByScoreRangeOptions{
		Limit: 2,
	})
	checkOrder(t, nodes, []int32{33, 10001113})

	minNode := sortedset.PeekMin()
	if minNode == nil || minNode.Key() != 1000 {
		t.Error("PeekMin() does not return expected value `d`")
	}

	minNode = sortedset.PopMin()
	if minNode == nil || minNode.Key() != 1000 {
		t.Error("PopMin() does not return expected value `d`")
	}

	nodes = sortedset.GetByScoreRange(-500, 500, nil)
	checkOrder(t, nodes, []int32{10001114, 12, 1000111, 10001112, 10001113, 33})

	maxNode := sortedset.PeekMax()
	if maxNode == nil || maxNode.Key() != 33 {
		t.Error("PeekMax() does not return expected value `c`")
	}

	maxNode = sortedset.PopMax()
	if maxNode == nil || maxNode.Key() != 33 {
		t.Error("PopMax() does not return expected value `c`")
	}

	nodes = sortedset.GetByScoreRange(500, -500, nil)
	checkOrder(t, nodes, []int32{10001113, 10001112, 1000111, 12, 10001114})
}
