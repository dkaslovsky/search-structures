package bst

import (
	"fmt"
	"testing"

	"github.com/dkaslovsky/search-structures/queue"
	"github.com/stretchr/testify/assert"
)

func equal(bst1 *BST, bst2 *BST) (eq bool, msg string) {
	q1 := queue.NewQueue()
	q1.Push(bst1)

	q2 := queue.NewQueue()
	q2.Push(bst2)

	for {
		item1, errItem1 := q1.Pop()
		item2, errItem2 := q2.Pop()
		if errItem1 == queue.ErrEmptyQueue && errItem2 == queue.ErrEmptyQueue {
			return true, ""
		}

		b1, ok := item1.(*BST)
		if !ok {
			return false, fmt.Sprintf("error casting first argument tree node: %v", b1)
		}
		b2, ok := item2.(*BST)
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
		tree1      *BST
		tree2      *BST
		expectedEq bool
	}{
		"equal single node trees": {
			tree1:      NewBST(10, "val10", nil, nil),
			tree2:      NewBST(10, "val10", nil, nil),
			expectedEq: true,
		},
		"equal multiple node trees": {
			tree1: NewBST(20, "val20",
				NewBST(10, "val10",
					nil,
					NewBST(15, "val15", nil, nil),
				),
				NewBST(30, "val30",
					NewBST(25, "val25", nil, nil),
					NewBST(40, "val40", nil, nil),
				),
			),
			tree2: NewBST(20, "val20",
				NewBST(10, "val10",
					nil,
					NewBST(15, "val15", nil, nil),
				),
				NewBST(30, "val30",
					NewBST(25, "val25", nil, nil),
					NewBST(40, "val40", nil, nil),
				),
			),
			expectedEq: true,
		},
		"unequal single node trees": {
			tree1:      NewBST(1, "val1", nil, nil),
			tree2:      NewBST(2, "val2", nil, nil),
			expectedEq: false,
		},
		"unequal single node trees with key mismatch": {
			tree1:      NewBST(1, "val1", nil, nil),
			tree2:      NewBST(2, "val1", nil, nil),
			expectedEq: false,
		},
		"unequal single node trees with value mismatch": {
			tree1:      NewBST(1, "val1", nil, nil),
			tree2:      NewBST(1, "val2", nil, nil),
			expectedEq: false,
		},
		"unequal trees with one having no children": {
			tree1: NewBST(1, "val1", nil, nil),
			tree2: NewBST(1, "val1",
				nil,
				NewBST(2, "val2", nil, nil),
			),
			expectedEq: false,
		},
		"unequal trees with children key mismatch": {
			tree1: NewBST(1, "val1",
				nil,
				NewBST(2, "val2", nil, nil),
			),
			tree2: NewBST(1, "val1",
				nil,
				NewBST(3, "val1", nil, nil),
			),
			expectedEq: false,
		},
		"unequal trees with children value mismatch": {
			tree1: NewBST(1, "val1",
				nil,
				NewBST(2, "val2", nil, nil),
			),
			tree2: NewBST(1, "val1",
				nil,
				NewBST(2, "val3", nil, nil),
			),
			expectedEq: false,
		},
		"unequal trees with different children structure": {
			tree1: NewBST(1, "val1",
				nil,
				NewBST(2, "val2", nil, nil),
			),
			tree2: NewBST(1, "val1",
				NewBST(0, "val0", nil, nil),
				nil,
			),
			expectedEq: false,
		},
		"unequal trees with deep different children structure": {
			tree1: NewBST(20, "val20",
				NewBST(10, "val10",
					nil,
					NewBST(15, "val15", nil, nil),
				),
				NewBST(30, "val30",
					NewBST(25, "val25", nil, nil),
					NewBST(40, "val40", nil, nil),
				),
			),
			tree2: NewBST(20, "val20",
				NewBST(10, "val10",
					nil,
					NewBST(15, "val15", nil, nil),
				),
				NewBST(30, "val30",
					NewBST(25, "val25", nil, nil),
					nil,
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
		tree         *BST
		insertKey    int64
		insertVal    string
		expectedTree *BST
	}{
		"insert should be on left": {
			tree:      NewBST(10, "val10", nil, nil),
			insertKey: 9,
			insertVal: "val9",
			expectedTree: NewBST(10, "val10",
				NewBST(9, "val9", nil, nil),
				nil,
			),
		},
		"insert should be on right": {
			tree:      NewBST(10, "val10", nil, nil),
			insertKey: 11,
			insertVal: "val11",
			expectedTree: NewBST(10, "val10",
				nil,
				NewBST(11, "val11", nil, nil),
			),
		},
		"deep insert to left child of left leaf": {
			tree: NewBST(20, "val20",
				NewBST(10, "val10",
					nil,
					NewBST(15, "val15", nil, nil),
				),
				NewBST(30, "val30",
					NewBST(25, "val25", nil, nil),
					NewBST(40, "val40", nil, nil),
				),
			),
			insertKey: 24,
			insertVal: "val24",
			expectedTree: NewBST(20, "val20",
				NewBST(10, "val10",
					nil,
					NewBST(15, "val15", nil, nil),
				),
				NewBST(30, "val30",
					NewBST(25, "val25",
						NewBST(24, "val24", nil, nil),
						nil,
					),
					NewBST(40, "val40", nil, nil),
				),
			),
		},
		"deep insert to right child of left leaf": {
			tree: NewBST(20, "val20",
				NewBST(10, "val10",
					nil,
					NewBST(15, "val15", nil, nil),
				),
				NewBST(30, "val30",
					NewBST(25, "val25", nil, nil),
					NewBST(40, "val40", nil, nil),
				),
			),
			insertKey: 26,
			insertVal: "val26",
			expectedTree: NewBST(20, "val20",
				NewBST(10, "val10",
					nil,
					NewBST(15, "val15", nil, nil),
				),
				NewBST(30, "val30",
					NewBST(25, "val25",
						nil,
						NewBST(26, "val26", nil, nil),
					),
					NewBST(40, "val40", nil, nil),
				),
			),
		},
		"deep insert to left child of right leaf": {
			tree: NewBST(20, "val20",
				NewBST(10, "val10",
					nil,
					NewBST(15, "val15", nil, nil),
				),
				NewBST(30, "val30",
					NewBST(25, "val25", nil, nil),
					NewBST(40, "val40", nil, nil),
				),
			),
			insertKey: 39,
			insertVal: "val39",
			expectedTree: NewBST(20, "val20",
				NewBST(10, "val10",
					nil,
					NewBST(15, "val15", nil, nil),
				),
				NewBST(30, "val30",
					NewBST(25, "val25", nil, nil),
					NewBST(40, "val40",
						NewBST(39, "val39", nil, nil),
						nil,
					),
				),
			),
		},
		"deep insert to right child of right leaf": {
			tree: NewBST(20, "val20",
				NewBST(10, "val10",
					nil,
					NewBST(15, "val15", nil, nil),
				),
				NewBST(30, "val30",
					NewBST(25, "val25", nil, nil),
					NewBST(40, "val40", nil, nil),
				),
			),
			insertKey: 41,
			insertVal: "val41",
			expectedTree: NewBST(20, "val20",
				NewBST(10, "val10",
					nil,
					NewBST(15, "val15", nil, nil),
				),
				NewBST(30, "val30",
					NewBST(25, "val25", nil, nil),
					NewBST(40, "val40",
						nil,
						NewBST(41, "val41", nil, nil),
					),
				),
			),
		},
		"insert of existing key overwrites value": {
			tree:         NewBST(10, "val10", nil, nil),
			insertKey:    10,
			insertVal:    "newVal10",
			expectedTree: NewBST(10, "newVal10", nil, nil),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.tree.Insert(test.insertKey, test.insertVal)
			assertBstEqual(t, test.tree, test.expectedTree)
		})
	}
}

func assertBstEqual(t *testing.T, bst1 *BST, bst2 *BST) {
	a := assert.New(t)
	eq, msg := equal(bst1, bst2)
	if !eq {
		a.Fail(msg)
	}
}
