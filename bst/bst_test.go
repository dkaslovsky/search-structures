package bst

import "testing"

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
			eq := test.tree1.Equal(test.tree2)
			if eq != test.expectedEq {
				t.Fatalf("fail")
			}
		})
	}
}
