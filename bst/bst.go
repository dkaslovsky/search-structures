package bst

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/dkaslovsky/search-structures/queue"
)

// Errors returned from a BST
var (
	ErrKeyNotFound    error = errors.New("key not found in BST")
	ErrDeleteRootLeaf error = errors.New("cannot delete node that is both a leaf and root of BST")
	ErrIteratorStop   error = errors.New("iterator is empty")
)

// BST is a node of a binary search tree indexed by Key containing value Val
type BST struct {
	Key   int64
	Val   string
	Left  *BST
	Right *BST

	r *rand.Rand
}

// NewBST constructs a BST
func NewBST(key int64, val string, left *BST, right *BST) *BST {
	return &BST{
		Key:   key,
		Val:   val,
		Left:  left,
		Right: right,

		r: rand.New(rand.NewSource(time.Now().UTC().UnixNano())),
	}
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
	// when the node-to-be-deleted has both left and right children, choose side to use for deleting
	// at random to avoid creating an unbalanced tree
	deleteSide := leftSide
	if b.r.Float64() > 0.5 {
		deleteSide = rightSide
	}

	return b.delete(key, deleteSide)
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
func (b *BST) Validate() (bool, error) {
	type validationNode struct {
		*BST
		minKey int64
		maxKey int64
	}

	q := queue.NewQueue()
	q.Push(&validationNode{
		BST:    b,
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
		if curNode.BST == nil {
			continue
		}
		if (curNode.Key < curNode.minKey) || (curNode.Key > curNode.maxKey) {
			return false, nil
		}

		left := &validationNode{
			BST:    curNode.Left,
			minKey: curNode.minKey,
			maxKey: curNode.Key - 1,
		}
		right := &validationNode{
			BST:    curNode.Right,
			minKey: curNode.Key + 1,
			maxKey: curNode.maxKey,
		}
		q.Push(left)
		q.Push(right)
	}
}

// Iterator creates a function to iterate the nodes of the BST by returning the next (breadth-first) node on each call
func (b *BST) Iterator() func() (*BST, error) {
	q := queue.NewQueue()
	q.Push(b)
	return func() (*BST, error) {
		item, err := q.Pop()
		if err == queue.ErrEmptyQueue {
			return nil, ErrIteratorStop
		}
		curB, ok := item.(*BST)
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

// Equal evaluates if two BSTs are equal
func (b *BST) Equal(other *BST) bool {
	q1 := queue.NewQueue()
	q1.Push(b)

	q2 := queue.NewQueue()
	q2.Push(other)

	for {
		item1, errItem1 := q1.Pop()
		item2, errItem2 := q2.Pop()
		if errItem1 == queue.ErrEmptyQueue && errItem2 == queue.ErrEmptyQueue {
			return true
		}

		b1, ok := item1.(*BST)
		if !ok {
			return false
		}
		b2, ok := item2.(*BST)
		if !ok {
			return false
		}

		if b1.Key != b2.Key || b1.Val != b2.Val {
			return false
		}

		if b1 == nil {
			continue
		}

		if b1.Left != nil {
			q1.Push(b1.Left)
		}
		if b1.Right != nil {
			q1.Push(b1.Right)
		}
		if b2.Left != nil {
			q2.Push(b2.Left)
		}
		if b2.Right != nil {
			q2.Push(b2.Right)
		}
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

func (b *BST) delete(key int64, deleteSide side) error {
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

func (b *BST) replaceChild(child *BST, newChild *BST) {
	if child.Key < b.Key {
		b.Left = newChild
		return
	}
	b.Right = newChild
}

func deleteOnLeft(target *BST) {
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

func deleteOnRight(target *BST) {
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

func (b *BST) findRightMost() (right *BST, parent *BST) {
	right, parent = b, nil
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
