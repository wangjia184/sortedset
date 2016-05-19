# Sorted Set in Go

Sorted Set is an ordered collection of nodes. Every node is associated with these properties.

| Property | Type | Description |
|---|---|---|
| `key` | `string` | The identifier of the node. It must be unique within the set. |
| `value` | `interface {}` | value associated with this node |
| `score` | `float64` | The score associated with this node, that in order to take the sorted set ordered. score may be repeated within the set. |

Each node in the set is associated with a `key`. While `key`s are unique, `score`s may be repeated.

With sorted sets you can add, remove, or update elements in a very fast way (in a time proportional to the logarithm of the number of elements). Since elements are __taken in order instead of ordered afterwards__, you can also get ranges by score or by rank (position) in a very fast way. Accessing the middle of a sorted set is also very fast, so you can use Sorted Sets as a smart list of non repeating elements where you can quickly access everything you need: elements in order, fast existence test, fast access to elements in the middle!