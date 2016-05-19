package sortedset

import (
	"bytes"
	"fmt"
	"math/rand"
	"runtime/debug"
	"strconv"
	"testing"
	"time"
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
	nodes := sortedset.GetByRankRange(1, -1, false)
	checkOrder(t, nodes, []string{"d", "h", "a", "e", "f", "g", "c"})

	// get & remove the 2nd/3rd nodes in reserve order
	nodes = sortedset.GetByRankRange(-2, -3, true)
	checkOrder(t, nodes, []string{"g", "f"})

	// get all nodes since the last one to first one
	nodes = sortedset.GetByRankRange(-1, 1, false)
	checkOrder(t, nodes, []string{"c", "e", "a", "h", "d"})

}

func TestCase2(t *testing.T) {
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

	nodes := sortedset.GetByScoreRange(-500, 500, nil)
	checkOrder(t, nodes, []string{"d", "h", "a", "e", "f", "g", "c"})

	nodes = sortedset.GetByScoreRange(500, -500, nil)
	//t.Logf("%v", nodes)
	checkOrder(t, nodes, []string{"c", "g", "f", "e", "a", "h", "d"})

	nodes = sortedset.GetByScoreRange(600, 500, nil)
	checkOrder(t, nodes, []string{})

	nodes = sortedset.GetByScoreRange(500, 600, nil)
	checkOrder(t, nodes, []string{})

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
}

var sets [100]*SortedSet

func TestTimeComplexity(t *testing.T) {
	debug.SetGCPercent(-1)
	var buffer bytes.Buffer

	buffer.WriteString("\nCount\t\tAdd\t\tGetByScore\tGetByRank\tRemove")

	sets[0] = test(&buffer, 10000)
	sets[1] = test(&buffer, 100000)
	sets[2] = test(&buffer, 1000000)
	//sets[3] = test(&buffer, 10000000)
	//sets[4] = test(&buffer, 100000000)

	t.Log(buffer.String())
}

type testData struct {
	key   string
	score SCORE
	value string
}

func test(buffer *bytes.Buffer, rounds int) *SortedSet {

	sortedset := New()

	data := make([]testData, rounds)
	for i := 0; i < rounds; i++ {
		data[i] = testData{
			score: SCORE(rand.NormFloat64()),
			value: "d",
			key:   strconv.Itoa(i),
		}
	}

	buffer.WriteString(fmt.Sprintf("\n%d\t\t", rounds))

	start := time.Now()
	for _, d := range data {
		sortedset.AddOrUpdate(d.key, d.score, d.value)
	}
	seconds := time.Now().Sub(start).Seconds() * 1000000 / float64(rounds)
	buffer.WriteString(fmt.Sprintf("%.2f us/op\t", seconds))

	for i := 0; i < 1000; i++ {
		ff := sortedset.GetByScoreRange(data[i].score, data[i].score, nil)
		if len(ff) == 0 {
			panic("Unable to find the item")
		}
	}

	start = time.Now()
	for i := 0; i < 1000; i++ {
		ff := sortedset.GetByScoreRange(data[i].score, data[i].score, nil)
		if len(ff) == 0 {
			panic("Unable to find the item")
		}
	}
	seconds = time.Now().Sub(start).Seconds() * 1000
	buffer.WriteString(fmt.Sprintf("%.2f us/op\t", seconds))

	var foundItems [1000]*SortedSetNode

	start = time.Now()
	for i := 0; i < cap(foundItems); i++ {
		rank := rounds/2 + i
		node := sortedset.GetByRank(rank, false)
		if node == nil {
			panic("Unable to find the item")
		}
		foundItems[i] = node
	}
	seconds = time.Now().Sub(start).Seconds() * 1000
	buffer.WriteString(fmt.Sprintf("%.2f us/op\t", seconds))

	start = time.Now()
	for _, item := range foundItems {
		sortedset.Remove(item.Key())
	}
	seconds = time.Now().Sub(start).Seconds() * 1000
	buffer.WriteString(fmt.Sprintf("%.2f us/op\t\t", seconds))
	return sortedset
}
