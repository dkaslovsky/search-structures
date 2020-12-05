package main

import (
	"fmt"
	"log"

	"github.com/dkaslovsky/search-structures/bst"
)

func main() {
	// b := bst.NewBST(20, "val20",
	// 	bst.NewBST(10, "val10",
	// 		nil,
	// 		bst.NewBST(15, "val15", nil, nil),
	// 	),
	// 	bst.NewBST(30, "val30",
	// 		bst.NewBST(25, "val25", nil, nil),
	// 		bst.NewBST(40, "val40", nil, nil),
	// 	),
	// )

	b := bst.NewBST(20, "val20",
		bst.NewBST(10, "val10",
			nil,
			bst.NewBST(15, "val15", nil, nil),
		),
		bst.NewBST(30, "val30",
			bst.NewBST(25, "val25", nil, nil),
			bst.NewBST(40, "val40",
				bst.NewBST(32, "val32", nil,
					bst.NewBST(34, "val34", nil, nil),
				),
				bst.NewBST(42, "val42", nil, nil),
			),
		),
	)

	iter := b.Iterator()
	for {
		node, err := iter()
		if err == bst.ErrIteratorStop {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(node)
	}

	valid, err := b.Validate()
	if err != nil {
		log.Fatal(err)
	}
	if !valid {
		log.Fatal("tree not valid")
	}

	fmt.Println("=======")

	b.Delete(30)
	iter = b.Iterator()
	for {
		node, err := iter()
		if err == bst.ErrIteratorStop {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(node)
	}

	valid, err = b.Validate()
	if err != nil {
		log.Fatal(err)
	}
	if !valid {
		log.Fatal("tree not valid")
	}

}
