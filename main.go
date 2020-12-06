package main

import (
	"fmt"
	"log"

	"github.com/dkaslovsky/search-structures/bst"
)

func main() {
	b := bst.NewBst(
		bst.NewBstNode(20, "val20",
			bst.NewBstNode(10, "val10",
				nil,
				bst.NewBstNode(15, "val15", nil, nil),
			),
			bst.NewBstNode(30, "val30",
				bst.NewBstNode(25, "val25", nil, nil),
				bst.NewBstNode(40, "val40",
					bst.NewBstNode(32, "val32", nil,
						bst.NewBstNode(34, "val34", nil, nil),
					),
					bst.NewBstNode(42, "val42", nil, nil),
				),
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
