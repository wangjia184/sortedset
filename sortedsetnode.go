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
	"golang.org/x/exp/constraints"
)

// SortedSetLevel ...
type SortedSetLevel[K Indexable, V any, ScoreType constraints.Ordered] struct {
	forward *SortedSetNode[K, V, ScoreType]
	span    int64
}

// SortedSetNode Node in skip list
type SortedSetNode[K Indexable, V any, ScoreType constraints.Ordered] struct {
	key      K         // unique key of this node
	Value    V         // associated data
	score    ScoreType // score to determine the order of this node in the set
	backward *SortedSetNode[K, V, ScoreType]
	level    []SortedSetLevel[K, V, ScoreType]
}

func NewSortedSetNode[K Indexable, V any, ScoreType constraints.Ordered](
	level int, key K, value V, score ScoreType,
) *SortedSetNode[K, V, ScoreType] {
	return &SortedSetNode[K, V, ScoreType]{
		key:   key,
		Value: value,
		score: score,
		level: make([]SortedSetLevel[K, V, ScoreType], level),
	}
}

// Key Get the key of the node
func (n *SortedSetNode[K, V, ScoreType]) Key() K {
	return n.key
}

// Score Get the node of the node
func (n *SortedSetNode[K, V, ScoreType]) Score() ScoreType {
	return n.score
}
