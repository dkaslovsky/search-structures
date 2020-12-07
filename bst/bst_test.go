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

		b1, ok := item1.(*Node)
		if !ok {
			return false, fmt.Sprintf("error casting first argument tree node: %v", b1)
		}
		b2, ok := item2.(*Node)
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
			tree1:      NewBst(NewNode(10, "val10", nil, nil)),
			tree2:      NewBst(NewNode(10, "val10", nil, nil)),
			expectedEq: true,
		},
		"equal multiple node trees": {
			tree1: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
						NewNode(40, "val40", nil, nil),
					),
				),
			),
			tree2: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
						NewNode(40, "val40", nil, nil),
					),
				),
			),
			expectedEq: true,
		},
		"unequal single node trees": {
			tree1:      NewBst(NewNode(1, "val1", nil, nil)),
			tree2:      NewBst(NewNode(2, "val2", nil, nil)),
			expectedEq: false,
		},
		"unequal single node trees with key mismatch": {
			tree1:      NewBst(NewNode(1, "val1", nil, nil)),
			tree2:      NewBst(NewNode(2, "val1", nil, nil)),
			expectedEq: false,
		},
		"unequal single node trees with value mismatch": {
			tree1:      NewBst(NewNode(1, "val1", nil, nil)),
			tree2:      NewBst(NewNode(1, "val2", nil, nil)),
			expectedEq: false,
		},
		"unequal trees with one having no children": {
			tree1: NewBst(NewNode(1, "val1", nil, nil)),
			tree2: NewBst(
				NewNode(1, "val1",
					nil,
					NewNode(2, "val2", nil, nil),
				),
			),
			expectedEq: false,
		},
		"unequal trees with children key mismatch": {
			tree1: NewBst(
				NewNode(1, "val1",
					nil,
					NewNode(2, "val2", nil, nil),
				),
			),
			tree2: NewBst(
				NewNode(1, "val1",
					nil,
					NewNode(3, "val1", nil, nil),
				),
			),
			expectedEq: false,
		},
		"unequal trees with children value mismatch": {
			tree1: NewBst(
				NewNode(1, "val1",
					nil,
					NewNode(2, "val2", nil, nil),
				),
			),
			tree2: NewBst(
				NewNode(1, "val1",
					nil,
					NewNode(2, "val3", nil, nil),
				),
			),
			expectedEq: false,
		},
		"unequal trees with different children structure": {
			tree1: NewBst(
				NewNode(1, "val1",
					nil,
					NewNode(2, "val2", nil, nil),
				),
			),
			tree2: NewBst(
				NewNode(1, "val1",
					NewNode(0, "val0", nil, nil),
					nil,
				),
			),
			expectedEq: false,
		},
		"unequal trees with deep different children structure": {
			tree1: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
						NewNode(40, "val40", nil, nil),
					),
				),
			),
			tree2: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
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
			tree:      NewBst(NewNode(10, "val10", nil, nil)),
			insertKey: 9,
			insertVal: "val9",
			expectedTree: NewBst(
				NewNode(10, "val10",
					NewNode(9, "val9", nil, nil),
					nil,
				),
			),
		},
		"insert should be on right": {
			tree:      NewBst(NewNode(10, "val10", nil, nil)),
			insertKey: 11,
			insertVal: "val11",
			expectedTree: NewBst(
				NewNode(10, "val10",
					nil,
					NewNode(11, "val11", nil, nil),
				),
			),
		},
		"deep insert to left child of left leaf": {
			tree: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
						NewNode(40, "val40", nil, nil),
					),
				),
			),
			insertKey: 24,
			insertVal: "val24",
			expectedTree: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25",
							NewNode(24, "val24", nil, nil),
							nil,
						),
						NewNode(40, "val40", nil, nil),
					),
				),
			),
		},
		"deep insert to right child of left leaf": {
			tree: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
						NewNode(40, "val40", nil, nil),
					),
				),
			),
			insertKey: 26,
			insertVal: "val26",
			expectedTree: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25",
							nil,
							NewNode(26, "val26", nil, nil),
						),
						NewNode(40, "val40", nil, nil),
					),
				),
			),
		},
		"deep insert to left child of right leaf": {
			tree: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
						NewNode(40, "val40", nil, nil),
					),
				),
			),
			insertKey: 39,
			insertVal: "val39",
			expectedTree: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
						NewNode(40, "val40",
							NewNode(39, "val39", nil, nil),
							nil,
						),
					),
				),
			),
		},
		"deep insert to right child of right leaf": {
			tree: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
						NewNode(40, "val40", nil, nil),
					),
				),
			),
			insertKey: 41,
			insertVal: "val41",
			expectedTree: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
						NewNode(40, "val40",
							nil,
							NewNode(41, "val41", nil, nil),
						),
					),
				),
			),
		},
		"insert of existing key overwrites value": {
			tree:         NewBst(NewNode(10, "val10", nil, nil)),
			insertKey:    10,
			insertVal:    "newVal10",
			expectedTree: NewBst(NewNode(10, "newVal10", nil, nil)),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			a := assert.New(t)
			test.tree.Insert(test.insertKey, test.insertVal)
			assertBstEqual(t, test.tree, test.expectedTree)

			valid, err := test.tree.Validate()
			a.NoError(err)
			a.True(valid)
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := map[string]struct {
		tree          *Bst
		expectedEmpty bool
	}{
		"empty tree": {
			tree:          NewBst(nil),
			expectedEmpty: true,
		},
		"nonempty tree": {
			tree:          NewBst(NewNode(10, "val10", nil, nil)),
			expectedEmpty: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expectedEmpty, test.tree.IsEmpty())
		})
	}
}

func TestSearch(t *testing.T) {
	tests := map[string]struct {
		tree           *Bst
		searchKey      int64
		expectedValue  string
		expectedExists bool
	}{
		"empty tree": {
			tree:           NewBst(nil),
			searchKey:      1,
			expectedValue:  "",
			expectedExists: false,
		},
		"single node tree without searchKey": {
			tree:           NewBst(NewNode(1, "val1", nil, nil)),
			searchKey:      2,
			expectedValue:  "",
			expectedExists: false,
		},
		"multi node tree without searchKey": {
			tree: NewBst(
				NewNode(10, "val1",
					NewNode(8, "val8", nil, nil),
					NewNode(12, "val12", nil, nil),
				),
			),
			searchKey:      2,
			expectedValue:  "",
			expectedExists: false,
		},
		"single node tree with searchKey": {
			tree:           NewBst(NewNode(1, "val1", nil, nil)),
			searchKey:      1,
			expectedValue:  "val1",
			expectedExists: true,
		},
		"multi node tree with searchKey": {
			tree: NewBst(
				NewNode(10, "val1",
					NewNode(8, "val8", nil, nil),
					NewNode(12, "val12", nil, nil),
				),
			),
			searchKey:      12,
			expectedValue:  "val12",
			expectedExists: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			a := assert.New(t)
			val, exists := test.tree.Search(test.searchKey)
			a.Equal(test.expectedExists, exists)
			if !test.expectedExists {
				return
			}
			a.Equal(test.expectedValue, val)
		})
	}
}

func TestDelete(t *testing.T) {
	// Delete() calls deleteBySide() in all cases except for an empty tree
	t.Run("empty tree", func(t *testing.T) {
		tree := NewBst(nil)
		err := tree.Delete(1)
		assert.Equal(t, ErrEmpty, err)
	})
}

func TestDeleteBySide(t *testing.T) {
	tests := map[string]struct {
		tree         *Bst
		side         side
		deleteKey    int64
		expectedTree *Bst
		expectedErr  error
	}{
		"single node tree": {
			tree:        NewBst(NewNode(10, "val10", nil, nil)),
			deleteKey:   10,
			expectedErr: ErrDeleteRootLeaf,
		},
		"multi node tree without deleteKey": {
			tree: NewBst(
				NewNode(10, "val10",
					NewNode(8, "val8", nil, nil),
					NewNode(12, "val12", nil, nil),
				),
			),
			deleteKey:   1,
			expectedErr: ErrKeyNotFound,
		},
		"multi node tree with deleteKey on left": {
			tree: NewBst(
				NewNode(10, "val10",
					NewNode(8, "val8", nil, nil),
					NewNode(12, "val12", nil, nil),
				),
			),
			deleteKey: 8,
			expectedTree: NewBst(
				NewNode(10, "val10",
					nil,
					NewNode(12, "val12", nil, nil),
				),
			),
			expectedErr: nil,
		},
		"multi node tree with deleteKey on right": {
			tree: NewBst(
				NewNode(10, "val10",
					NewNode(8, "val8", nil, nil),
					NewNode(12, "val12", nil, nil),
				),
			),
			deleteKey: 12,
			expectedTree: NewBst(
				NewNode(10, "val10",
					NewNode(8, "val8", nil, nil),
					nil,
				),
			),
			expectedErr: nil,
		},
		"multi node unbalanced tree with deleteKey having left child": {
			tree: NewBst(
				NewNode(10, "val10",
					NewNode(5, "val5",
						NewNode(2, "val2", nil, nil),
						nil,
					),
					nil,
				),
			),
			deleteKey: 5,
			expectedTree: NewBst(
				NewNode(10, "val10",
					NewNode(2, "val2", nil, nil),
					nil,
				),
			),
			expectedErr: nil,
		},
		"multi node unbalanced tree with deleteKey having right child": {
			tree: NewBst(
				NewNode(10, "val10",
					NewNode(5, "val5",
						nil,
						NewNode(7, "val7", nil, nil),
					),
					nil,
				),
			),
			deleteKey: 5,
			expectedTree: NewBst(
				NewNode(10, "val10",
					nil,
					NewNode(7, "val7", nil, nil),
				),
			),
			expectedErr: nil,
		},
		"leftSide delete on multi node tree with parent deleteKey": {
			tree: NewBst(
				NewNode(10, "val10",
					NewNode(8, "val8", nil, nil),
					NewNode(12, "val12", nil, nil),
				),
			),
			side:      leftSide,
			deleteKey: 10,
			expectedTree: NewBst(
				NewNode(8, "val8",
					nil,
					NewNode(12, "val12", nil, nil),
				),
			),
			expectedErr: nil,
		},
		"rightSide delete on multi node tree with parent deleteKey": {
			tree: NewBst(
				NewNode(10, "val10",
					NewNode(8, "val8", nil, nil),
					NewNode(12, "val12", nil, nil),
				),
			),
			side:      rightSide,
			deleteKey: 10,
			expectedTree: NewBst(
				NewNode(12, "val12",
					NewNode(8, "val8", nil, nil),
					nil,
				),
			),
			expectedErr: nil,
		},
		"leftSide delete on deep multi node tree": {
			tree: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
						NewNode(40, "val40",
							NewNode(32, "val32", nil,
								NewNode(34, "val34", nil, nil),
							),
							NewNode(42, "val42", nil, nil),
						),
					),
				),
			),
			side:      leftSide,
			deleteKey: 30,
			expectedTree: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(25, "val25",
						nil,
						NewNode(40, "val40",
							NewNode(32, "val32", nil,
								NewNode(34, "val34", nil, nil),
							),
							NewNode(42, "val42", nil, nil),
						),
					),
				),
			),
			expectedErr: nil,
		},
		"rightSide delete on deep multi node tree": {
			tree: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
						NewNode(40, "val40",
							NewNode(32, "val32", nil,
								NewNode(34, "val34", nil, nil),
							),
							NewNode(42, "val42", nil, nil),
						),
					),
				),
			),
			side:      rightSide,
			deleteKey: 30,
			expectedTree: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(32, "val32",
						NewNode(25, "val25", nil, nil),
						NewNode(40, "val40",
							NewNode(34, "val34", nil, nil),
							NewNode(42, "val42", nil, nil),
						),
					),
				),
			),
			expectedErr: nil,
		},
		"leftSide delete on deep multi node tree at parent": {
			tree: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
						NewNode(40, "val40",
							NewNode(32, "val32", nil,
								NewNode(34, "val34", nil, nil),
							),
							NewNode(42, "val42", nil, nil),
						),
					),
				),
			),
			side:      leftSide,
			deleteKey: 20,
			expectedTree: NewBst(
				NewNode(15, "val15",
					NewNode(10, "val10", nil, nil),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
						NewNode(40, "val40",
							NewNode(32, "val32", nil,
								NewNode(34, "val34", nil, nil),
							),
							NewNode(42, "val42", nil, nil),
						),
					),
				),
			),
			expectedErr: nil,
		},
		"rightSide delete on deep multi node tree at parent": {
			tree: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
						NewNode(40, "val40",
							NewNode(32, "val32", nil,
								NewNode(34, "val34", nil, nil),
							),
							NewNode(42, "val42", nil, nil),
						),
					),
				),
			),
			side:      rightSide,
			deleteKey: 20,
			expectedTree: NewBst(
				NewNode(25, "val25",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						nil,
						NewNode(40, "val40",
							NewNode(32, "val32", nil,
								NewNode(34, "val34", nil, nil),
							),
							NewNode(42, "val42", nil, nil),
						),
					),
				),
			),
			expectedErr: nil,
		},

		"deep multi node tree with leaf deleteKey": {
			tree: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
						NewNode(40, "val40",
							NewNode(32, "val32", nil,
								NewNode(34, "val34", nil, nil),
							),
							NewNode(42, "val42", nil, nil),
						),
					),
				),
			),
			deleteKey: 42,
			expectedTree: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
						NewNode(40, "val40",
							NewNode(32, "val32", nil,
								NewNode(34, "val34", nil, nil),
							),
							nil,
						),
					),
				),
			),
			expectedErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			a := assert.New(t)
			err := test.tree.deleteBySide(test.deleteKey, test.side)
			if err != nil {
				a.Equal(test.expectedErr, err)
				return
			}
			assertBstEqual(t, test.expectedTree, test.tree)

			valid, err := test.tree.Validate()
			a.NoError(err)
			a.True(valid)
		})
	}
}

func TestValidate(t *testing.T) {
	tests := map[string]struct {
		tree          *Bst
		expectedValid bool
		expectedErr   error
	}{
		"empty tree": {
			tree:          NewBst(nil),
			expectedValid: false,
			expectedErr:   ErrEmpty,
		},
		"single node tree": {
			tree:          NewBst(NewNode(10, "val10", nil, nil)),
			expectedValid: true,
			expectedErr:   nil,
		},
		"multi node valid tree": {
			tree: NewBst(
				NewNode(10, "val10",
					NewNode(8, "val8", nil, nil),
					NewNode(12, "val12", nil, nil),
				),
			),
			expectedValid: true,
			expectedErr:   nil,
		},
		"multi node invalid tree": {
			tree: NewBst(
				NewNode(10, "val1",
					NewNode(11, "va11", nil, nil),
					NewNode(12, "val12", nil, nil),
				),
			),
			expectedValid: false,
			expectedErr:   nil,
		},
		"deep multi node valid tree": {
			tree: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
						NewNode(40, "val40", nil, nil),
					),
				),
			),
			expectedValid: true,
			expectedErr:   nil,
		},
		"deep multi node invalid tree": {
			tree: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
						NewNode(24, "val24", nil, nil),
					),
				),
			),
			expectedValid: false,
			expectedErr:   nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			a := assert.New(t)
			valid, err := test.tree.Validate()
			a.Equal(test.expectedValid, valid)
			a.Equal(test.expectedErr, err)
		})
	}
}

func TestIterator(t *testing.T) {
	type testNode struct {
		key int64
		val string
	}

	tests := map[string]struct {
		tree                  *Bst
		expectedIteratedNodes []testNode
	}{
		"empty tree": {
			tree:                  NewBst(nil),
			expectedIteratedNodes: []testNode{},
		},
		"single node tree": {
			tree:                  NewBst(NewNode(10, "val10", nil, nil)),
			expectedIteratedNodes: []testNode{testNode{10, "val10"}},
		},
		"multi node tree": {
			tree: NewBst(
				NewNode(10, "val10",
					NewNode(8, "val8", nil, nil),
					NewNode(12, "val12", nil, nil),
				),
			),
			expectedIteratedNodes: []testNode{
				testNode{10, "val10"},
				testNode{8, "val8"},
				testNode{12, "val12"},
			},
		},
		"deep multi node tree": {
			tree: NewBst(
				NewNode(20, "val20",
					NewNode(10, "val10",
						nil,
						NewNode(15, "val15", nil, nil),
					),
					NewNode(30, "val30",
						NewNode(25, "val25", nil, nil),
						NewNode(40, "val40", nil, nil),
					),
				),
			),
			expectedIteratedNodes: []testNode{
				testNode{20, "val20"},
				testNode{10, "val10"},
				testNode{30, "val30"},
				testNode{15, "val15"},
				testNode{25, "val25"},
				testNode{40, "val40"},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			a := assert.New(t)
			iter := test.tree.Iterator()
			nodes := []testNode{}
			for {
				node, err := iter()
				if err == ErrIteratorStop {
					break
				}
				a.NoError(err)
				nodes = append(nodes, testNode{node.Key, node.Val})
			}
			a.Equal(test.expectedIteratedNodes, nodes)
		})
	}
}

func assertBstEqual(t *testing.T, bst1 *Bst, bst2 *Bst) {
	a := assert.New(t)
	eq, msg := equal(bst1, bst2)
	if !eq {
		a.Fail(msg)
	}
}
