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

import "golang.org/x/exp/constraints"

type SortedSetLevel[K constraints.Ordered, SCORE constraints.Ordered, V any] struct {
	forward *SortedSetNode[K, SCORE, V]
	span    int64
}

// Node in skip list
type SortedSetNode[K constraints.Ordered, SCORE constraints.Ordered, V any] struct {
	key      K     // unique key of this node
	Value    V     // associated data
	score    SCORE // score to determine the order of this node in the set
	backward *SortedSetNode[K, SCORE, V]
	level    []SortedSetLevel[K, SCORE, V]
}

// Get the key of the node
func (this *SortedSetNode[K, SCORE, V]) Key() K {
	return this.key
}

// Get the score of the node
func (this *SortedSetNode[K, SCORE, V]) Score() SCORE {
	return this.score
}
