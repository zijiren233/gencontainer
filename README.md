# <u>gencontainer</u>

- generic container
- more containers
- faster speed
- more methods

# <u>containers</u>

- `doubly linked list`
  > A doubly linked list is a data structure that consists of nodes, where each node contains a value and two pointers: one pointing to the previous node and another pointing to the next node. This structure allows traversal in both directions, forward and backward.
  > 
  > The main purpose of a doubly linked list is to provide efficient insertion and deletion operations at any position in the list. Unlike a singly linked list, which only allows traversal in one direction, a doubly linked list enables easy access to both the previous and next nodes, making it suitable for scenarios where frequent insertions and deletions are required.
  >
  > The principle behind a doubly linked list is that each node contains references to both the previous and next nodes, forming a chain-like structure. This allows for efficient insertion and deletion by updating the pointers of adjacent nodes.
  >
  > Doubly linked lists find applications in various scenarios, such as implementing data structures like queues, stacks, and hash tables. They are also useful in scenarios where bidirectional traversal is necessary, such as in certain graph algorithms or text editors.
  ```go
  l := dllist.New[int]()
  e1 := l.PushBack(1)
  e2 := l.PushBack(2)
  e := l.Front()
  l.Remove(e)
  l.Remove(e)
  ```

- `set`
  > Set is a fundamental data structure in programming that represents an unordered collection of unique elements. It is commonly used to store a group of distinct values without any specific order.
  >
  > The main purpose of using a set is to efficiently check for membership and eliminate duplicates. The underlying principle of a set is based on mathematical set theory, where elements are either present or absent. This allows for fast operations like adding, removing, and checking for existence of elements.
  >
  > Sets are particularly useful in scenarios where you need to perform operations like union, intersection, and difference between multiple sets. They provide a convenient way to handle tasks such as finding common elements, removing duplicates, or checking for uniqueness.
  >
  > In programming, sets find applications in various domains, such as data analysis, graph algorithms, and solving problems that involve unique values. They offer a powerful tool to manage collections of items with unique characteristics efficiently.
  ```go
  s1 := set.New[int]()
  s2 := set.New[int]()
  for i := 0; i < 10; i++ {
  s1.Insert(i)
  s2.Insert(i)
  }
  s1.Remove(0)
  s2.Remove(2)
  s3 := s1.Union(s2)
  s4 := s1.Difference(s2)
  ```

- `vector`
  > The structure of vector and slice is very similar, but vector provides more methods to manipulate data more conveniently
  ```go
  v := vec.New[int](WithValues(1, 2, 3), WithValues(4, 5, 6), WithCap[int](10))
  v.Insert(1, 10, 11, 12)
  if e, ok := v.Remove(1); !ok || e != 10 {
    t.Fatal("wrong remove")
  }
  ```

- `hashring`
  > Hashring, also known as consistent hashing algorithm, is a technique used in programming. It provides a way to distribute data across multiple nodes in a scalable and efficient manner.
  >
  > The key features of hashring include load balancing and fault tolerance. It achieves this by mapping data to a ring-like structure using a hash function.
  >
  > Hashring finds applications in distributed caching, distributed databases, and content delivery networks (CDNs), enabling efficient data storage and retrieval across a distributed system.
  ```go
  hr := hashring.New[string](50, WithNodes("node1", "node2", "node3", "node4"))
  n := hr.GetNode("somehash to get node")
  // n = "node3"
  ```

- `rwmap`
  > This rwmap is similar to sync.Map
  >
  > Support for generics, more utility methods
  ```go
  m := rwmap.New[int, int]()
	m.Store(1, 1)
	m.Store(2, 2)
	m.Store(3, 3)
	if v := m.LoadAndDeleteAll(); len(v) != 3 {
		t.Errorf("LoadAndDeleteAll() = %v, want 3", v)
	}
	if v := m.LoadAndDeleteAll(); len(v) != 0 {
		t.Errorf("LoadAndDeleteAll() = %v, want 0", v)
	}
  ```