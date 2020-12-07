package bst

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/dkaslovsky/search-structures/queue"
)

// Errors returned from a Bst
var (
	ErrEmpty          error = errors.New("Bst is empty")
	ErrKeyNotFound    error = errors.New("key not found in Bst")
	ErrDeleteRootLeaf error = errors.New("cannot delete node that is both a leaf and root of Bst")
	ErrIteratorStop   error = errors.New("iterator stopped after iterating all nodes")
)

// Bst is a binary search tree
type Bst struct {
	Tree *Node
	r    *rand.Rand
}

// NewBst constructs a Bst
func NewBst(tree *Node) *Bst {
	return &Bst{
		Tree: tree,
		r:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Node is a node of a binary search tree indexed by Key containing value Val
type Node struct {
	Key   int64
	Val   string
	Left  *Node
	Right *Node
}

// NewNode constructs a node
func NewNode(key int64, val string, left *Node, right *Node) *Node {
	return &Node{
		Key:   key,
		Val:   val,
		Left:  left,
		Right: right,
	}
}

// IsEmpty evaluates is a Bst is empty
func (b *Bst) IsEmpty() bool {
	return b.Tree == nil
}

// Insert inserts a key/value pair
func (b *Bst) Insert(key int64, val string) {
	if b.IsEmpty() {
		b.Tree = NewNode(key, val, nil, nil)
		return
	}

	curTree := b.Tree
	for {
		if key == curTree.Key {
			// allow an existing value to be overwritten
			curTree.Val = val
			return
		}
		if key < curTree.Key {
			if curTree.Left == nil {
				curTree.Left = NewNode(key, val, nil, nil)
				return
			}
			curTree = curTree.Left
			continue
		}
		if curTree.Right == nil {
			curTree.Right = NewNode(key, val, nil, nil)
			return
		}
		curTree = curTree.Right
	}
}

// Delete deletes a key/value pair
func (b *Bst) Delete(key int64) error {
	if b.IsEmpty() {
		return ErrEmpty
	}

	// when the node-to-be-deleted has both left and right children, choose side to use for deleting
	// at random to avoid creating an unbalanced tree
	deleteSide := leftSide
	if b.r.Float64() > 0.5 {
		deleteSide = rightSide
	}

	return b.deleteBySide(key, deleteSide)
}

// Search searches a Bst for a key
func (b *Bst) Search(key int64) (val string, found bool) {
	if b.IsEmpty() {
		return val, false
	}

	target, _, found := b.search(key)
	if !found {
		return "", false
	}
	return target.Val, true
}

// Validate determines if a Bst satisfies the Bst property
func (b *Bst) Validate() (bool, error) {
	if b.IsEmpty() {
		return false, ErrEmpty
	}

	type validationNode struct {
		*Node
		minKey int64
		maxKey int64
	}

	q := queue.NewQueue()
	q.Push(&validationNode{
		Node:   b.Tree,
		minKey: math.MinInt64,
		maxKey: math.MaxInt64,
	})

	for {
		item, err := q.Pop()
		if err == queue.ErrEmptyQueue {
			return true, nil
		}
		curNode, ok := item.(*validationNode)
		if !ok {
			return false, errors.New("reached node of unknown type while traversing tree")
		}
		if curNode.Node == nil {
			continue
		}
		if (curNode.Key < curNode.minKey) || (curNode.Key > curNode.maxKey) {
			return false, nil
		}

		left := &validationNode{
			Node:   curNode.Left,
			minKey: curNode.minKey,
			maxKey: curNode.Key - 1,
		}
		right := &validationNode{
			Node:   curNode.Right,
			minKey: curNode.Key + 1,
			maxKey: curNode.maxKey,
		}
		q.Push(left)
		q.Push(right)
	}
}

// Iterator creates a function to iterate the nodes of the Bst by returning the next (breadth-first) node on each call
func (b *Bst) Iterator() func() (*Node, error) {
	if b.IsEmpty() {
		return func() (*Node, error) {
			return nil, ErrIteratorStop
		}
	}

	q := queue.NewQueue()
	q.Push(b.Tree)
	return func() (*Node, error) {
		item, err := q.Pop()
		if err == queue.ErrEmptyQueue {
			return nil, ErrIteratorStop
		}
		curB, ok := item.(*Node)
		if !ok {
			return nil, errors.New("reached node of unknown type while traversing tree")
		}
		if curB == nil {
			return curB, nil
		}
		if curB.Left != nil {
			q.Push(curB.Left)
		}
		if curB.Right != nil {
			q.Push(curB.Right)
		}
		return curB, nil
	}
}

// search searches for a key and returns the node, the parent node, and success bool
func (b *Bst) search(key int64) (target *Node, parent *Node, found bool) {
	curTree := b.Tree
	for {
		if key == curTree.Key {
			return curTree, parent, true
		}
		if key < curTree.Key {
			if curTree.Left == nil {
				return target, parent, false
			}
			parent = curTree
			curTree = curTree.Left
			continue
		}
		if curTree.Right == nil {
			return target, parent, false
		}
		parent = curTree
		curTree = curTree.Right
	}
}

func (b *Bst) deleteBySide(key int64, deleteSide side) error {
	target, parent, found := b.search(key)
	if !found {
		return ErrKeyNotFound
	}

	// target is a leaf node
	if target.Left == nil && target.Right == nil {
		// cannot delete if target is both a root and leaf node
		if parent == nil {
			return ErrDeleteRootLeaf
		}
		parent.replaceChild(target, nil)
		target = nil
		return nil
	}

	// target has only a right child
	if target.Left == nil {
		parent.replaceChild(target, target.Right)
		target = nil
		return nil
	}

	// target has only a left child
	if target.Right == nil {
		parent.replaceChild(target, target.Left)
		target = nil
		return nil
	}

	// target has both left and right children:
	// delete by overwriting with leftmost (min) value from right branch or rightmost (max) value
	// from left branch
	switch deleteSide {
	case leftSide:
		deleteOnLeft(target)
	case rightSide:
		deleteOnRight(target)
	default:
		return fmt.Errorf("deleteSide must be one of [%v, %v], received [%v]", leftSide, rightSide, deleteSide)
	}

	return nil
}

func (n *Node) replaceChild(child *Node, newChild *Node) {
	if child.Key < n.Key {
		n.Left = newChild
		return
	}
	n.Right = newChild
}

func deleteOnLeft(target *Node) {
	right, parent := target.Left.findRightMost()

	// overwrite target's key/value with left's key/value
	target.Key = right.Key
	target.Val = right.Val

	if parent != nil {
		// if right has a child it must be on the left and greater than parent, so
		// it becomes parent's child on the right
		parent.Right = right.Left
		right = nil
		return
	}

	// right and target.Left are the same, set both to nil
	target.Left = nil
	right = nil
}

func deleteOnRight(target *Node) {
	left, parent := target.Right.findLeftMost()

	// overwrite target's key/value with left's key/value
	target.Key = left.Key
	target.Val = left.Val

	if parent != nil {
		// if left has a child it must be on the right and less than parent, so
		// it becomes parent's child on the left
		parent.Left = left.Right
		left = nil
		return
	}

	// left and target.Right are the same, set both to nil
	target.Right = nil
	left = nil
}

func (n *Node) findLeftMost() (left *Node, parent *Node) {
	left, parent = n, nil
	for {
		if left.Left == nil {
			return left, parent
		}
		parent = left
		left = left.Left
	}
}

func (n *Node) findRightMost() (right *Node, parent *Node) {
	right, parent = n, nil
	for {
		if right.Right == nil {
			return right, parent
		}
		parent = right
		right = right.Right
	}
}

type side uint

const (
	leftSide side = iota
	rightSide
)
