# Sorted Set in Go

A Redis-inspired sorted set backed by a skip list. Nodes are taken in order
from low score to high score; ties are ordered by key. Access by key is via a
hash index, so membership tests are O(1) and score/rank operations are O(log N).



## API

`SortedSet[K constraints.Ordered, SCORE constraints.Ordered, V any]`

| Method | Description |
| --- | --- |
| `AddOrUpdate(key K, score SCORE, value V) bool` | Insert or update; `true` when the key was new |
| `Remove(key K) *SortedSetNode[K, SCORE, V]` | Delete by key |
| `GetByKey(key K) *SortedSetNode[...]` | Look up a node by key |
| `GetCount() int` | Number of nodes |
| `PeekMin() / PopMin() / PeekMax() / PopMax()` | Extremes, with or without removal |
| `FindRank(key K) int` | 1-based rank of a key (0 when absent) |
| `GetRangeByScore(start, end SCORE, options *GetRangeByScoreOptions)` | Nodes whose score is in range |
| `GetRangeByRank(start, end int, remove bool)` | Nodes by 1-based rank range |
| `GetByRank(rank int, remove bool)` | Single node by rank |
| `IterFuncRangeByRank(start, end int, fn func(K, V) bool)` | Iterate a rank range |
| `Has(key K) bool` | Concurrent-safe membership test |

A node exposes `Key() K`, `Score() SCORE`, and the public `Value V` field.

## Concurrency contract

The set is **not** safe for concurrent use as a whole. Only `Has` may be
called from a goroutine other than the one owning the set: it reads the
`sync.Map`-backed index only, so a handler goroutine can poll membership while
the owner mutates. Every other method takes no lock and must be called from a
single goroutine.

## Requirements

Go 1.26+ (uses `golang.org/x/exp/constraints`).

## Usage

```go
import sortedset "github.com/wangjia184/sortedset"

set := sortedset.New[string, int64, string]()
set.AddOrUpdate("a", 89, "Kelly")
set.AddOrUpdate("b", 100, "Staley")

min := set.PeekMin() // Key() == "a"
rank := set.FindRank("b") // 2
all := set.GetRangeByRank(1, -1, false) // ascending
if set.Has("a") { // safe from any goroutine
    // ...
}
```

## Timing

All operations are O(log N), except `GetByKey`/`Has`/`Remove` lookups which
are O(1). GetRangeByScore is O(log N) to locate the start, then O(m) for m
returned nodes.


