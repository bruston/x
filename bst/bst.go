package bst

type Tree struct {
	Less func(a, b interface{}) bool
	Root *Node
}

type Node struct {
	Key         interface{}
	Value       interface{}
	Left, Right *Node
}

type LessFunc func(a, b interface{}) bool

func IntLess(a, b interface{}) bool {
	if a.(int) <= b.(int) {
		return true
	}
	return false
}

func New(fn LessFunc) *Tree {
	return &Tree{Less: fn}
}

func (t *Tree) Search(k interface{}) (interface{}, bool) {
	var cur = t.Root
	for cur != nil {
		if k == cur.Key {
			return cur.Value, true
		}
		if t.Less(k, cur.Key) {
			cur = cur.Left
		} else {
			cur = cur.Right
		}
	}
	return nil, false
}

func (t *Tree) Insert(k, v interface{}) {
	var cur = t.Root
	var prev *Node
	leftTurn := false
	for cur != nil {
		if k == cur.Key {
			cur.Value = v
			return
		}
		prev = cur
		if t.Less(k, cur.Key) {
			cur = cur.Left
			leftTurn = true
		} else {
			cur = cur.Right
			leftTurn = false
		}
	}
	node := &Node{Key: k, Value: v}
	if t.Root == nil {
		t.Root = node
		return
	}
	if leftTurn {
		prev.Left = node
	} else {
		prev.Right = node
	}
}
