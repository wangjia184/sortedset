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

/*
Package sortedset provides an ordered collection implemented basing on skiplist for fast access


Every node in the set is associated with these properties.

| Property | Type | Description |
|---|---|---|
| `key` | `string` | The identifier of the node. It must be unique within the set. |
| `value` | `interface {}` | value associated with this node |
| `score` | `float64` | score is in order to take the sorted set ordered. It may be repeated. |

Each node in the set is associated with a `key`. While `key`s are unique, `score`s may be repeated.
Nodes are __taken in order instead of ordered afterwards__, from low score to high score. If scores are the same, the node is ordered by its key in lexicographic order. Each node in the set is associated with __rank__, which represents the position of the node in the sorted set. The __rank__ is 1-based, that is to say, rank 1 is the node with minimum score.

Sorted Set is implemented basing on skip list and hash map internally. With sorted sets you can add, remove, or update nodes in a very fast way (in a time proportional to the logarithm of the number of nodes). You can also get ranges by score or by rank (position) in a very fast way. Accessing the middle of a sorted set is also very fast, so you can use Sorted Sets as a smart list of non repeating nodes where you can quickly access everything you need: nodes in order, fast existence test, fast access to nodes in the middle!

A typical use case of sorted set is a leader board in a massive online game, where every time a new score is submitted you update it using `AddOrUpdate()`. You can easily take the top users using `GetByRankRange()`, you can also, given an user id, return its rank in the listing using `FindRank()`. Using `FindRank()` and `GetByRankRange()` together you can show users with a score similar to a given user. All very quickly.


*/
package sortedset
