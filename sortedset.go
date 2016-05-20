// Copyright (c) 2016, Jerry.Wang
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//  list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//  this list of conditions and the following disclaimer in the documentation
//  and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package sortedset

import (
	"math/rand"
)

type SCORE int64 // the type of score

const SKIPLIST_MAXLEVEL = 32 /* Should be enough for 2^32 elements */
const SKIPLIST_P = 0.25      /* Skiplist P = 1/4 */

type SortedSet struct {
	header *SortedSetNode
	tail   *SortedSetNode
	length int64
	level  int
	dict   map[string]*SortedSetNode
}

func createNode(level int, score SCORE, key string, value interface{}) *SortedSetNode {
	node := SortedSetNode{
		score: score,
		key:   key,
		Value: value,
		level: make([]SortedSetLevel, level),
	}
	return &node
}

// Returns a random level for the new skiplist node we are going to create.
// The return value of this function is between 1 and SKIPLIST_MAXLEVEL
// (both inclusive), with a powerlaw-alike distribution where higher
// levels are less likely to be returned.
func randomLevel() int {
	level := 1
	for float64(rand.Int31()&0xFFFF) < float64(SKIPLIST_P*0xFFFF) {
		level += 1
	}
	if level < SKIPLIST_MAXLEVEL {
		return level
	}

	return SKIPLIST_MAXLEVEL
}

func (this *SortedSet) insertNode(score SCORE, key string, value interface{}) *SortedSetNode {
	var update [SKIPLIST_MAXLEVEL]*SortedSetNode
	var rank [SKIPLIST_MAXLEVEL]int64

	x := this.header
	for i := this.level - 1; i >= 0; i-- {
		/* store rank that is crossed to reach the insert position */
		if this.level-1 == i {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		for x.level[i].forward != nil &&
			(x.level[i].forward.score < score ||
				(x.level[i].forward.score == score && // score is the same but the key is different
					x.level[i].forward.key < key)) {
			rank[i] += x.level[i].span
			x = x.level[i].forward
		}
		update[i] = x
	}

	/* we assume the key is not already inside, since we allow duplicated
	 * scores, and the re-insertion of score and redis object should never
	 * happen since the caller of Insert() should test in the hash table
	 * if the element is already inside or not. */
	level := randomLevel()

	if level > this.level { // add a new level
		for i := this.level; i < level; i++ {
			rank[i] = 0
			update[i] = this.header
			update[i].level[i].span = this.length
		}
		this.level = level
	}

	x = createNode(level, score, key, value)
	for i := 0; i < level; i++ {
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x

		/* update span covered by update[i] as x is inserted here */
		x.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	/* increment span for untouched levels */
	for i := level; i < this.level; i++ {
		update[i].level[i].span++
	}

	if update[0] == this.header {
		x.backward = nil
	} else {
		x.backward = update[0]
	}
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x
	} else {
		this.tail = x
	}
	this.length++
	return x
}

/* Internal function used by delete, DeleteByScore and DeleteByRank */
func (this *SortedSet) deleteNode(x *SortedSetNode, update [SKIPLIST_MAXLEVEL]*SortedSetNode) {
	for i := 0; i < this.level; i++ {
		if update[i].level[i].forward == x {
			update[i].level[i].span += x.level[i].span - 1
			update[i].level[i].forward = x.level[i].forward
		} else {
			update[i].level[i].span -= 1
		}
	}
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x.backward
	} else {
		this.tail = x.backward
	}
	for this.level > 1 && this.header.level[this.level-1].forward == nil {
		this.level--
	}
	this.length--
	delete(this.dict, x.key)
}

/* Delete an element with matching score/key from the skiplist. */
func (this *SortedSet) delete(score SCORE, key string) bool {
	var update [SKIPLIST_MAXLEVEL]*SortedSetNode

	x := this.header
	for i := this.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.score < score ||
				(x.level[i].forward.score == score &&
					x.level[i].forward.key < key)) {
			x = x.level[i].forward
		}
		update[i] = x
	}
	/* We may have multiple elements with the same score, what we need
	 * is to find the element with both the right score and object. */
	x = x.level[0].forward
	if x != nil && score == x.score && x.key == key {
		this.deleteNode(x, update)
		// free x
		return true
	}
	return false /* not found */
}

// Create a new SortedSet
func New() *SortedSet {
	sortedSet := SortedSet{
		level: 1,
		dict:  make(map[string]*SortedSetNode),
	}
	sortedSet.header = createNode(SKIPLIST_MAXLEVEL, 0, "", nil)
	return &sortedSet
}

// Get the number of elements
func (this *SortedSet) GetCount() int {
	return int(this.length)
}

// get the element with minimum score, nil if the set is empty
//
// Time complexity of this method is : O(log(N))
func (this *SortedSet) PeekMin() *SortedSetNode {
	return this.header.level[0].forward
}

// get and remove the element with minimal score, nil if the set is empty
//
// // Time complexity of this method is : O(log(N))
func (this *SortedSet) PopMin() *SortedSetNode {
	x := this.header.level[0].forward
	if x != nil {
		this.Remove(x.key)
	}
	return x
}

// get the element with maximum score, nil if the set is empty
// Time Complexity : O(1)
func (this *SortedSet) PeekMax() *SortedSetNode {
	return this.tail
}

// get and remove the element with maximum score, nil if the set is empty
//
// Time complexity of this method is : O(log(N))
func (this *SortedSet) PopMax() *SortedSetNode {
	x := this.tail
	if x != nil {
		this.Remove(x.key)
	}
	return x
}

// Add an element into the sorted set with specific key / value / score.
// if the element is added, this method returns true; otherwise false means updated
//
// Time complexity of this method is : O(log(N))
func (this *SortedSet) AddOrUpdate(key string, score SCORE, value interface{}) bool {
	var newNode *SortedSetNode = nil

	found := this.dict[key]
	if found != nil {
		// score does not change, only update value
		if found.score == score {
			found.Value = value
		} else { // score changes, delete and re-insert
			this.delete(found.score, found.key)
			newNode = this.insertNode(score, key, value)
		}
	} else {
		newNode = this.insertNode(score, key, value)
	}

	if newNode != nil {
		this.dict[key] = newNode
	}
	return found == nil
}

// Delete element specified by key
//
// Time complexity of this method is : O(log(N))
func (this *SortedSet) Remove(key string) *SortedSetNode {
	found := this.dict[key]
	if found != nil {
		this.delete(found.score, found.key)
		return found
	}
	return nil
}

type GetByScoreRangeOptions struct {
	Limit        int  // limit the max nodes to return
	ExcludeStart bool // exclude start value, so it search in interval (start, end] or (start, end)
	ExcludeEnd   bool // exclude end value, so it search in interval [start, end) or (start, end)
}

// Get the nodes whose score within the specific range
//
// If options is nil, it searchs in interval [start, end] without any limit by default
//
// Time complexity of this method is : O(log(N))
func (this *SortedSet) GetByScoreRange(start SCORE, end SCORE, options *GetByScoreRangeOptions) []*SortedSetNode {

	// prepare parameters
	var limit int = 2147483648
	if options != nil && options.Limit > 0 {
		limit = options.Limit
	}

	excludeStart := options != nil && options.ExcludeStart
	excludeEnd := options != nil && options.ExcludeEnd
	reverse := start > end
	if reverse {
		start, end = end, start
		excludeStart, excludeEnd = excludeEnd, excludeStart
	}

	//////////////////////////
	var nodes []*SortedSetNode

	//determine if out of range
	if this.length == 0 {
		return nodes
	}
	//////////////////////////

	if reverse { // search from end to start
		x := this.header

		if excludeEnd {
			for i := this.level - 1; i >= 0; i-- {
				for x.level[i].forward != nil &&
					x.level[i].forward.score < end {
					x = x.level[i].forward
				}
			}
		} else {
			for i := this.level - 1; i >= 0; i-- {
				for x.level[i].forward != nil &&
					x.level[i].forward.score <= end {
					x = x.level[i].forward
				}
			}
		}

		for x != nil && limit > 0 {
			if excludeStart {
				if x.score <= start {
					break
				}
			} else {
				if x.score < start {
					break
				}
			}

			next := x.backward

			nodes = append(nodes, x)
			limit--

			x = next
		}
	} else {
		// search from start to end
		x := this.header
		if excludeStart {
			for i := this.level - 1; i >= 0; i-- {
				for x.level[i].forward != nil &&
					x.level[i].forward.score <= start {
					x = x.level[i].forward
				}
			}
		} else {
			for i := this.level - 1; i >= 0; i-- {
				for x.level[i].forward != nil &&
					x.level[i].forward.score < start {
					x = x.level[i].forward
				}
			}
		}

		/* Current node is the last with score < or <= start. */
		x = x.level[0].forward

		for x != nil && limit > 0 {
			if excludeEnd {
				if x.score >= end {
					break
				}
			} else {
				if x.score > end {
					break
				}
			}

			next := x.level[0].forward

			nodes = append(nodes, x)
			limit--

			x = next
		}
	}

	return nodes
}

// Get nodes within specific rank range [start, end]
// Note that the rank is 1-based integer. Rank 1 means the first node; Rank -1 means the last node;
//
// If start is greater than end, the returned array is in reserved order
// If remove is true, the returned nodes are removed
//
// Time complexity of this method is : O(log(N))
func (this *SortedSet) GetByRankRange(start int, end int, remove bool) []*SortedSetNode {

	/* Sanitize indexes. */
	if start < 0 {
		start = int(this.length) + start + 1
	}
	if end < 0 {
		end = int(this.length) + end + 1
	}
	if start <= 0 {
		start = 1
	}
	if end <= 0 {
		end = 1
	}

	reverse := start > end
	if reverse { // swap start and end
		start, end = end, start
	}

	var update [SKIPLIST_MAXLEVEL]*SortedSetNode
	var nodes []*SortedSetNode
	var traversed int = 0

	x := this.header
	for i := this.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			traversed+int(x.level[i].span) < start {
			traversed += int(x.level[i].span)
			x = x.level[i].forward
		}
		if remove {
			update[i] = x
		} else {
			if traversed+1 == start {
				break
			}
		}
	}

	traversed++
	x = x.level[0].forward
	for x != nil && traversed <= end {
		next := x.level[0].forward

		nodes = append(nodes, x)

		if remove {
			this.deleteNode(x, update)
		}

		traversed++
		x = next
	}

	if reverse {
		for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
			nodes[i], nodes[j] = nodes[j], nodes[i]
		}
	}
	return nodes
}

// Get node by rank.
// Note that the rank is 1-based integer. Rank 1 means the first node; Rank -1 means the last node;
//
// If remove is true, the returned nodes are removed
// If node is not found at specific rank, nil is returned
//
// Time complexity of this method is : O(log(N))
func (this *SortedSet) GetByRank(rank int, remove bool) *SortedSetNode {
	nodes := this.GetByRankRange(rank, rank, remove)
	if len(nodes) == 1 {
		return nodes[0]
	}
	return nil
}

// Get node by key
//
// If node is not found, nil is returned
// Time complexity : O(1)
func (this *SortedSet) GetByKey(key string) *SortedSetNode {
	return this.dict[key]
}

// Find the rank of the node specified by key
// Note that the rank is 1-based integer. Rank 1 means the first node
//
// If the node is not found, 0 is returned. Otherwise rank(> 0) is returned
//
// Time complexity of this method is : O(log(N))
func (this *SortedSet) FindRank(key string) int {
	var rank int = 0
	node := this.dict[key]
	if node != nil {
		x := this.header
		for i := this.level - 1; i >= 0; i-- {
			for x.level[i].forward != nil &&
				(x.level[i].forward.score < node.score ||
					(x.level[i].forward.score == node.score &&
						x.level[i].forward.key <= node.key)) {
				rank += int(x.level[i].span)
				x = x.level[i].forward
			}

			if x.key == key {
				return rank
			}
		}
	}
	return 0
}
