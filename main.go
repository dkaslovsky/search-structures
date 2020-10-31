package main

import (
	"fmt"

	"github.com/dkaslovsky/search-structures/bst"
)

func main() {
	b := bst.NewBST(20, "val20",
		bst.NewBST(10, "val10",
			nil,
			bst.NewBST(15, "val15", nil, nil),
		),
		bst.NewBST(30, "val30",
			bst.NewBST(25, "val25", nil, nil),
			bst.NewBST(40, "val40", nil, nil),
		),
	)

	fmt.Println("validate =", b.Validate())
}
