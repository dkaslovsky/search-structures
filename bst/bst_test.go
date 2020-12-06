package bst

import (
	"fmt"
	"testing"

	"github.com/dkaslovsky/search-structures/queue"
	"github.com/stretchr/testify/assert"
)

func equal(bst1 *Bst, bst2 *Bst) (eq bool, msg string) {
	if bst1.IsEmpty() || bst2.IsEmpty() {
		if bst1.IsEmpty() && bst2.IsEmpty() {
			return true, ""
		}
		return false, "found one empty and one nonempty tree"
	}

	if bst1.IsEmpty() {
		if bst2.IsEmpty() {
			return true, ""
		}
	}

	q1 := queue.NewQueue()
	q1.Push(bst1.Tree)

	q2 := queue.NewQueue()
	q2.Push(bst2.Tree)

	for {
		item1, errItem1 := q1.Pop()
		item2, errItem2 := q2.Pop()
		if errItem1 == queue.ErrEmptyQueue && errItem2 == queue.ErrEmptyQueue {
			return true, ""
		}

		b1, ok := item1.(*BstNode)
		if !ok {
			return false, fmt.Sprintf("error casting first argument tree node: %v", b1)
		}
		b2, ok := item2.(*BstNode)
		if !ok {
			return false, fmt.Sprintf("error casting second argument tree node: %v", b2)
		}

		if b1.Key != b2.Key || b1.Val != b2.Val {
			return false, fmt.Sprintf("unequal nodes:\n%v\n%v", b1, b2)
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

// test equal so it can be used in subsequent assertions
func TestEqual(t *testing.T) {

	tests := map[string]struct {
		tree1      *Bst
		tree2      *Bst
		expectedEq bool
	}{
		"empty trees": {
			tree1:      NewBst(nil),
			tree2:      NewBst(nil),
			expectedEq: true,
		},
		"equal single node trees": {
			tree1:      NewBst(NewBstNode(10, "val10", nil, nil)),
			tree2:      NewBst(NewBstNode(10, "val10", nil, nil)),
			expectedEq: true,
		},
		"equal multiple node trees": {
			tree1: NewBst(
				NewBstNode(20, "val20",
					NewBstNode(10, "val10",
						nil,
						NewBstNode(15, "val15", nil, nil),
					),
					NewBstNode(30, "val30",
						NewBstNode(25, "val25", nil, nil),
						NewBstNode(40, "val40", nil, nil),
					),
				),
			),
			tree2: NewBst(
				NewBstNode(20, "val20",
					NewBstNode(10, "val10",
						nil,
						NewBstNode(15, "val15", nil, nil),
					),
					NewBstNode(30, "val30",
						NewBstNode(25, "val25", nil, nil),
						NewBstNode(40, "val40", nil, nil),
					),
				),
			),
			expectedEq: true,
		},
		"unequal single node trees": {
			tree1:      NewBst(NewBstNode(1, "val1", nil, nil)),
			tree2:      NewBst(NewBstNode(2, "val2", nil, nil)),
			expectedEq: false,
		},
		"unequal single node trees with key mismatch": {
			tree1:      NewBst(NewBstNode(1, "val1", nil, nil)),
			tree2:      NewBst(NewBstNode(2, "val1", nil, nil)),
			expectedEq: false,
		},
		"unequal single node trees with value mismatch": {
			tree1:      NewBst(NewBstNode(1, "val1", nil, nil)),
			tree2:      NewBst(NewBstNode(1, "val2", nil, nil)),
			expectedEq: false,
		},
		"unequal trees with one having no children": {
			tree1: NewBst(NewBstNode(1, "val1", nil, nil)),
			tree2: NewBst(
				NewBstNode(1, "val1",
					nil,
					NewBstNode(2, "val2", nil, nil),
				),
			),
			expectedEq: false,
		},
		"unequal trees with children key mismatch": {
			tree1: NewBst(
				NewBstNode(1, "val1",
					nil,
					NewBstNode(2, "val2", nil, nil),
				),
			),
			tree2: NewBst(
				NewBstNode(1, "val1",
					nil,
					NewBstNode(3, "val1", nil, nil),
				),
			),
			expectedEq: false,
		},
		"unequal trees with children value mismatch": {
			tree1: NewBst(
				NewBstNode(1, "val1",
					nil,
					NewBstNode(2, "val2", nil, nil),
				),
			),
			tree2: NewBst(
				NewBstNode(1, "val1",
					nil,
					NewBstNode(2, "val3", nil, nil),
				),
			),
			expectedEq: false,
		},
		"unequal trees with different children structure": {
			tree1: NewBst(
				NewBstNode(1, "val1",
					nil,
					NewBstNode(2, "val2", nil, nil),
				),
			),
			tree2: NewBst(
				NewBstNode(1, "val1",
					NewBstNode(0, "val0", nil, nil),
					nil,
				),
			),
			expectedEq: false,
		},
		"unequal trees with deep different children structure": {
			tree1: NewBst(
				NewBstNode(20, "val20",
					NewBstNode(10, "val10",
						nil,
						NewBstNode(15, "val15", nil, nil),
					),
					NewBstNode(30, "val30",
						NewBstNode(25, "val25", nil, nil),
						NewBstNode(40, "val40", nil, nil),
					),
				),
			),
			tree2: NewBst(
				NewBstNode(20, "val20",
					NewBstNode(10, "val10",
						nil,
						NewBstNode(15, "val15", nil, nil),
					),
					NewBstNode(30, "val30",
						NewBstNode(25, "val25", nil, nil),
						nil,
					),
				),
			),
			expectedEq: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			eq, _ := equal(test.tree1, test.tree2)
			assert.Equal(t, test.expectedEq, eq)
		})
	}
}

func TestInsert(t *testing.T) {

	tests := map[string]struct {
		tree         *Bst
		insertKey    int64
		insertVal    string
		expectedTree *Bst
	}{
		"insert should be on left": {
			tree:      NewBst(NewBstNode(10, "val10", nil, nil)),
			insertKey: 9,
			insertVal: "val9",
			expectedTree: NewBst(
				NewBstNode(10, "val10",
					NewBstNode(9, "val9", nil, nil),
					nil,
				),
			),
		},
		"insert should be on right": {
			tree:      NewBst(NewBstNode(10, "val10", nil, nil)),
			insertKey: 11,
			insertVal: "val11",
			expectedTree: NewBst(
				NewBstNode(10, "val10",
					nil,
					NewBstNode(11, "val11", nil, nil),
				),
			),
		},
		"deep insert to left child of left leaf": {
			tree: NewBst(
				NewBstNode(20, "val20",
					NewBstNode(10, "val10",
						nil,
						NewBstNode(15, "val15", nil, nil),
					),
					NewBstNode(30, "val30",
						NewBstNode(25, "val25", nil, nil),
						NewBstNode(40, "val40", nil, nil),
					),
				),
			),
			insertKey: 24,
			insertVal: "val24",
			expectedTree: NewBst(
				NewBstNode(20, "val20",
					NewBstNode(10, "val10",
						nil,
						NewBstNode(15, "val15", nil, nil),
					),
					NewBstNode(30, "val30",
						NewBstNode(25, "val25",
							NewBstNode(24, "val24", nil, nil),
							nil,
						),
						NewBstNode(40, "val40", nil, nil),
					),
				),
			),
		},
		"deep insert to right child of left leaf": {
			tree: NewBst(
				NewBstNode(20, "val20",
					NewBstNode(10, "val10",
						nil,
						NewBstNode(15, "val15", nil, nil),
					),
					NewBstNode(30, "val30",
						NewBstNode(25, "val25", nil, nil),
						NewBstNode(40, "val40", nil, nil),
					),
				),
			),
			insertKey: 26,
			insertVal: "val26",
			expectedTree: NewBst(
				NewBstNode(20, "val20",
					NewBstNode(10, "val10",
						nil,
						NewBstNode(15, "val15", nil, nil),
					),
					NewBstNode(30, "val30",
						NewBstNode(25, "val25",
							nil,
							NewBstNode(26, "val26", nil, nil),
						),
						NewBstNode(40, "val40", nil, nil),
					),
				),
			),
		},
		"deep insert to left child of right leaf": {
			tree: NewBst(
				NewBstNode(20, "val20",
					NewBstNode(10, "val10",
						nil,
						NewBstNode(15, "val15", nil, nil),
					),
					NewBstNode(30, "val30",
						NewBstNode(25, "val25", nil, nil),
						NewBstNode(40, "val40", nil, nil),
					),
				),
			),
			insertKey: 39,
			insertVal: "val39",
			expectedTree: NewBst(
				NewBstNode(20, "val20",
					NewBstNode(10, "val10",
						nil,
						NewBstNode(15, "val15", nil, nil),
					),
					NewBstNode(30, "val30",
						NewBstNode(25, "val25", nil, nil),
						NewBstNode(40, "val40",
							NewBstNode(39, "val39", nil, nil),
							nil,
						),
					),
				),
			),
		},
		"deep insert to right child of right leaf": {
			tree: NewBst(
				NewBstNode(20, "val20",
					NewBstNode(10, "val10",
						nil,
						NewBstNode(15, "val15", nil, nil),
					),
					NewBstNode(30, "val30",
						NewBstNode(25, "val25", nil, nil),
						NewBstNode(40, "val40", nil, nil),
					),
				),
			),
			insertKey: 41,
			insertVal: "val41",
			expectedTree: NewBst(
				NewBstNode(20, "val20",
					NewBstNode(10, "val10",
						nil,
						NewBstNode(15, "val15", nil, nil),
					),
					NewBstNode(30, "val30",
						NewBstNode(25, "val25", nil, nil),
						NewBstNode(40, "val40",
							nil,
							NewBstNode(41, "val41", nil, nil),
						),
					),
				),
			),
		},
		"insert of existing key overwrites value": {
			tree:         NewBst(NewBstNode(10, "val10", nil, nil)),
			insertKey:    10,
			insertVal:    "newVal10",
			expectedTree: NewBst(NewBstNode(10, "newVal10", nil, nil)),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.tree.Insert(test.insertKey, test.insertVal)
			assertBstNodeEqual(t, test.tree, test.expectedTree)
		})
	}
}

func assertBstNodeEqual(t *testing.T, bst1 *Bst, bst2 *Bst) {
	a := assert.New(t)
	eq, msg := equal(bst1, bst2)
	if !eq {
		a.Fail(msg)
	}
}
