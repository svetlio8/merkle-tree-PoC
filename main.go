package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"strings"
)

func main() {

	data := []Hashable{
		Block("alice-test-20"),
		Block("bob-test-20"),
		Block("carol-test-20"),
		Block("dave-test-20"),
		Block("eric-test-20"),
		Block("fill-test-20"),
		Block("g-test-20"),
		Block("h-test-20"),
		Block("i-test-20"),
	}
	tree := buildTree(data)
	t := tree[0]
	test := t.(Node)
	rootHash := test.hash()
	log.Println("rootHash", rootHash)
	log.Println("Test Tree:")
	printTree(test)

	stringTestW := "WRONG-bob-test-20"
	testLeaf := crypto.Keccak256([]byte(stringTestW))
	stringHash := fmt.Sprintf("%x", testLeaf)
	log.Println("str: ", stringTestW)
	log.Println("stringHash: ", stringHash)
	isValidW := isValidLeaf(test, 0, testLeaf)
	log.Println("isValidLeaf wrong: ", isValidW)

	stringTestV := "bob-test-20"
	testLeafV := crypto.Keccak256([]byte(stringTestV))
	stringHashV := fmt.Sprintf("%x", testLeafV)
	log.Println("str: ", stringTestV)
	log.Println("stringHash: ", stringHashV)
	isValidV := isValidLeaf(test, 0, testLeafV)
	log.Println("isValidLeaf valid: ", isValidV)

	var leaf []Block
	getTreeLeaf(test, 0, &leaf)
	log.Println("leaf: ", leaf)

	leafHash := convertLeafToHash(leaf)
	log.Println("leafHash: ", leafHash)

}

// buildTree recursively builds the merkle tree, bottom-first. It makes a node-list for each level until it results in
// 1 node, meaning we reached the root node.
func buildTree(parts []Hashable) []Hashable {
	var nodes []Hashable
	var i int
	for i = 0; i < len(parts); i += 2 {
		if i + 1 < len(parts) {
			newNode := Node{
				left: parts[i],
				right: parts[i+1],
			}
			nodes = append(nodes, newNode)
		} else {
			newNode := Node{
				left: parts[i],
				//right: parts[i], //duplicate vs empty
				right: EmptyBlock{},
			}
			nodes = append(nodes, newNode)
		}
	}
	if len(nodes) == 1 {
		return nodes
	} else if len(nodes) > 1 {
		return buildTree(nodes)
	} else {
		panic("Not valid tree")
	}
}

type Hashable interface {
	hash() Hash
}

type Hash []byte

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

type Block string

func (b Block) hash() Hash {
	return hash([]byte(b)[:])
}

type EmptyBlock struct {
}

func (_ EmptyBlock) hash() Hash {
	return []byte{}
}

type Node struct {
	left  Hashable
	right Hashable
}

func (n Node) hash() Hash {
	var l, r []byte
	l = n.left.hash()
	r = n.right.hash()
	return hash(append(l[:], r[:]...))
}

func hash(data []byte) Hash {
	return crypto.Keccak256(data)
}

func printTree(node Node) {
	printNode(node, 0)
}

func printNode(node Node, level int) {
	fmt.Printf("(%d) %s %s\n", level, strings.Repeat(" ", level), node.hash())

	if l, ok := node.left.(Node); ok {
		printNode(l, level+1)
	} else if l, ok := node.left.(Block); ok {
		fmt.Printf("(%d) %s %s (data: %s)\n", level + 1, strings.Repeat(" ", level + 1), l.hash(), l)
	}
	if r, ok := node.right.(Node); ok {
		printNode(r, level+1)
	} else if r, ok := node.right.(Block); ok {
		fmt.Printf("(%d) %s %s (data: %s)\n", level + 1, strings.Repeat(" ", level + 1), r.hash(), r)
	}
}

//func isValidLeaf(node Node, level int, leaf [20]byte) bool {
func isValidLeaf(node Node, level int, leaf []byte) bool {

	if bytes.Equal(leaf, node.hash()) {
		return true
	}

	if l, ok := node.left.(Node); ok {
		if isValidLeaf(l, level+1, leaf) {
			return true
		}
	} else if l, ok := node.left.(Block); ok {
		if bytes.Equal(leaf, l.hash()) {
			return true
		}
	}

	if r, ok := node.right.(Node); ok {
		if isValidLeaf(r, level+1, leaf) {
			return true
		}
	} else if r, ok := node.right.(Block); ok {
		if bytes.Equal(leaf, r.hash()) {
			return true
		}
	}

	return false
}

func getTreeLeaf(node Node, level int, result *[]Block) {

	if l, ok := node.left.(Node); ok {
		getTreeLeaf(l, level+1, result)
	} else if l, ok := node.left.(Block); ok {
		*result = append(*result, l)
	}

	if r, ok := node.right.(Node); ok {
		getTreeLeaf(r, level+1, result)
	} else if r, ok := node.right.(Block); ok {
		*result = append(*result, r)
	}

	return
}

func convertLeafToHash(leaf []Block) []Hash {
	var res []Hash

	for _, v := range leaf {
		res = append(res, v.hash())
	}

	return res
}
