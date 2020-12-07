package main

import (
	"fmt"
	"log"

	"github.com/dkaslovsky/search-structures/bst"
)

func main() {
	b := bst.NewBst(
		bst.NewNode(20, "val20",
			bst.NewNode(10, "val10",
				nil,
				bst.NewNode(15, "val15", nil, nil),
			),
			bst.NewNode(30, "val30",
				bst.NewNode(25, "val25", nil, nil),
				bst.NewNode(40, "val40",
					bst.NewNode(32, "val32", nil,
						bst.NewNode(34, "val34", nil, nil),
					),
					bst.NewNode(42, "val42", nil, nil),
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
