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
Package sortedset provides the data-struct allows fast access the element in set by key or by score(order). It is inspired by Sorted Set from Redis.

Introduction

Every node in the set is associated with these properties.

    type SortedSetNode struct {
        key      string      // unique key of this node
        Value    interface{} // associated data
        score    SCORE       // int64 score to determine the order of this node in the set
    }

Each node in the set is associated with a key. While keys are unique, scores may be repeated. Nodes are taken in order (from low score to high score) instead of ordered afterwards. If scores are the same, the node is ordered by its key in lexicographic order. Each node in the set also can be accessed by rank, which represents the position in the sorted set.

Sorted Set is implemented basing on skip list and hash map internally. With sorted sets you can add, remove, or update nodes in a very fast way (in a time proportional to the logarithm of the number of nodes). You can also get ranges by score or by rank (position) in a very fast way. Accessing the middle of a sorted set is also very fast, so you can use Sorted Sets as a smart list of non repeating nodes where you can quickly access everything you need: nodes in order, fast existence test, fast access to nodes in the middle!

Use Case

A typical use case of sorted set is a leader board in a massive online game, where every time a new score is submitted you update it using AddOrUpdate() method. You can easily take the top users using GetByRankRange() method, you can also, given an user id, return its rank in the listing using FindRank() method. Using FindRank() and GetByRankRange() together you can show users with a score similar to a given user. All very quickly.


Examples

    // create a new set
    set := sortedset.New()

    // fill in new node
    set.AddOrUpdate("a", 89, "Kelly")
    set.AddOrUpdate("b", 100, "Staley")
    set.AddOrUpdate("c", 100, "Jordon")
    set.AddOrUpdate("d", -321, "Park")
    set.AddOrUpdate("e", 101, "Albert")
    set.AddOrUpdate("f", 99, "Lyman")
    set.AddOrUpdate("g", 99, "Singleton")
    set.AddOrUpdate("h", 70, "Audrey")

    // update an existing node by key
    set.AddOrUpdate("e", 99, "ntrnrt")

    // get the node by key
    set.GetByKey("b")

    // remove node by key
    set.Remove("b")

    // get the number of nodes in this set
    set.GetCount()

    // find the rank(postion) in the set.
    set.FindRank("d") // return 1 here

    // get and remove the node with minimum score
    set.PopMin()

    // get the node with maximum score
    set.PeekMax()

    // get the node at rank 1 (the node with minimum score)
    set.GetByRank(1, false)

    // get & remove the node at rank -1 (the node with maximum score)
    set.GetByRank(-1, true)

    // get the node with the 2nd highest maximum score
    set.GetByRank(-2, false)

    // get nodes with in rank range [1, -1],  that is all nodes actually
    set.GetByRankRange(1, -1, false)

    // get & remove the 2nd/3rd nodes in reserve order
    set.GetByRankRange(-2, -3, true)

    // get the nodes whose score are within the interval [60,100]
    set.GetByScoreRange(60, 100, nil)

    // get the nodes whose score are within the interval (60,100]
    set.GetByScoreRange(60, 100, &GetByScoreRangeOptions{
        ExcludeStart: true,
    })

    // get the nodes whose score are within the interval [60,100)
    set.GetByScoreRange(60, 100, &GetByScoreRangeOptions{
        ExcludeEnd: true,
    })

    // get the nodes whose score are within the interval [60,100] in reverse order
    set.GetByScoreRange(100, 60, nil)

    // get the top 2 nodes with lowest scores within the interval [60,100]
    set.GetByScoreRange(60, 100, &GetByScoreRangeOptions{
        Limit: 2,
    })

    // get the top 2 nodes with highest scores within the interval [60,100]
    set.GetByScoreRange(100, 60, &GetByScoreRangeOptions{
        Limit: 2,
    })

    // get the top 2 nodes with highest scores within the interval (60,100)
    set.GetByScoreRange(100, 60, &GetByScoreRangeOptions{
        Limit: 2,
        ExcludeStart: true,
        ExcludeEnd: true,
    })
*/
package sortedset
