package bst

import (
	"errors"
	"math"

	"github.com/dkaslovsky/search-structures/queue"
)

// Errors returned from bst package
var (
	ErrKeyNotFound error = errors.New("key not found in BST")
	ErrDeleteRoot  error = errors.New("cannot delete root node of BST")
)

// BST is a node of a binary search tree indexed by Key containing value Val
type BST struct {
	Key   int64
	Val   string
	Left  *BST
	Right *BST
}

// NewBST constructs a BST
func NewBST(key int64, val string, left *BST, right *BST) *BST {
	return &BST{
		Key:   key,
		Val:   val,
		Left:  left,
		Right: right,
	}
}

// NewBSTRoot constructs a root node of a BST with both children set to nil
func NewBSTRoot(key int64, val string) *BST {
	return NewBST(key, val, nil, nil)
}

// Insert inserts a key/value pair
func (b *BST) Insert(key int64, val string) {
	curTree := b
	for {
		if key == curTree.Key {
			// allow an existing value to be overwritten
			curTree.Val = val
			return
		}
		if key < curTree.Key {
			if curTree.Left == nil {
				curTree.Left = NewBST(key, val, nil, nil)
				return
			}
			curTree = curTree.Left
			continue
		}
		if curTree.Right == nil {
			curTree.Right = NewBST(key, val, nil, nil)
			return
		}
		curTree = curTree.Right
	}
}

// Delete deletes a key/value pair
func (b *BST) Delete(key int64) error {
	target, parent, found := b.search(key)
	if !found {
		return ErrKeyNotFound
	}

	// target is a leaf node
	if target.Left == nil && target.Right == nil {
		// cannot delete if target is also root node
		if parent == nil {
			return ErrDeleteRoot
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

	// target has both left and right children, delete by overwriting with leftmost (min) value from
	// right branch (or rightmost (max) value from left branch)
	left, leftParent := target.Right.findLeftMost()
	// parent will be nil when leftmost is the starting node (target.Right) so set parent to target
	if leftParent == nil {
		leftParent = target
	}
	// overwrite target's key/value with left's key/value
	target.Key = left.Key
	target.Val = left.Val
	// if left has a child it must be a right-child, so set leftParent.Left to be the child
	if left.Right != nil {
		leftParent.Left = left.Right
	}
	left = nil
	return nil
}

// Search searches a BST for a key
func (b *BST) Search(key int64) (val string, found bool) {
	target, _, found := b.search(key)
	if !found {
		return "", false
	}
	return target.Val, true
}

// Validate determines if a BST satisfies the BST property
func (b *BST) Validate() bool {
	type node struct {
		*BST
		minKey int64
		maxKey int64
	}

	q := queue.NewQueue()
	q.Push(&node{
		BST:    b,
		minKey: math.MinInt64,
		maxKey: math.MaxInt64,
	})

	for {
		item, err := q.Pop()
		if err == queue.ErrEmptyQueue {
			return true
		}
		curNode, _ := item.(*node)
		if curNode.BST == nil {
			continue
		}
		if (curNode.Key < curNode.minKey) || (curNode.Key > curNode.maxKey) {
			return false
		}

		left := &node{
			BST:    curNode.Left,
			minKey: curNode.minKey,
			maxKey: curNode.Key - 1,
		}
		right := &node{
			BST:    curNode.Right,
			minKey: curNode.Key + 1,
			maxKey: curNode.maxKey,
		}
		q.Push(left)
		q.Push(right)
	}
}

// search searches for a key and returns the node, the parent node, and success bool
func (b *BST) search(key int64) (target *BST, parent *BST, found bool) {
	curTree := b
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

func (b *BST) replaceChild(child *BST, newChild *BST) {
	if child.Key < b.Key {
		b.Left = newChild
		return
	}
	b.Right = newChild
}

func (b *BST) findLeftMost() (left *BST, parent *BST) {
	left, parent = b, nil
	for {
		if left.Left == nil {
			return left, parent
		}
		parent = left
		left = left.Left
	}
}
