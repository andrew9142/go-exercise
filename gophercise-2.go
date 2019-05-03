package main

import (
	"fmt"
	"time"

	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {
	//fmt.Println("Value: ", t.Value)
	ch <- t.Value
	if t.Left != nil {
		Walk(t.Left, ch)
	}

	if t.Right != nil {
		Walk(t.Right, ch)
	}

}
func Walk1(t *tree.Tree, ch chan int) {
	//fmt.Println("Value: ", t.Value)
	treeBuffer := make([]*tree.Tree, 1)
	treeBuffer[0] = t
	index := 0
	var current *tree.Tree
	for {
		if index < len(treeBuffer) {
			current = treeBuffer[index]
			index = index + 1
			//fmt.Println("Tree[", id, "]:", current.Value)
		} else {
			//fmt.Println("done")
			current = nil
		}

		if current != nil {
			if current.Value == 0 {
				close(ch)
				return
			}

			ch <- current.Value
			if current.Left != nil {
				treeBuffer = append(treeBuffer, current.Left)
			}
			if current.Right != nil {
				treeBuffer = append(treeBuffer, current.Right)
			}

		} else {
			close(ch)
			return
		}
	}

}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	buffer1 := make([]int, 0)
	buffer2 := make([]int, 0)

	channel1 := make(chan int)
	channel2 := make(chan int)
	ok1 := true
	ok2 := true
	var value1, value2 int

	go Walk1(t1, channel1)
	go Walk1(t2, channel2)
	for {
		select {
		case value1, ok1 = <-channel1:
			if ok1 {
				buffer1 = append(buffer1, value1)
			}
		case value2, ok2 = <-channel2:
			if ok2 {
				buffer2 = append(buffer2, value2)
			}
		}

		if ok1 == false && ok2 == false {
			break
		}
	}
	fmt.Println("Tree1: ", buffer1)
	fmt.Println("Tree2: ", buffer2)
	if len(buffer1) != len(buffer2) {
		return false
	}

	for i := range buffer1 {
		if buffer1[i] != buffer2[i] {
			return false
		}
		return true // done
	}

	return false // default
}

func Same1(t1, t2 *tree.Tree) bool {
	var buffer1, buffer2 [10]int
	index1 := 0
	index2 := 0
	channel1 := make(chan int)
	channel2 := make(chan int)
	go Walk(t1, channel1)
	go Walk(t2, channel2)
	tick := time.Tick(10 * time.Millisecond)
	stopFlag := false

	for {
		select {
		case value1 := <-channel1:
			stopFlag = false
			buffer1[index1] = value1
			index1 = index1 + 1
		case value2 := <-channel2:
			stopFlag = false
			buffer2[index2] = value2
			index2 = index2 + 1
		case <-tick:
			if stopFlag {
				fmt.Println("Tree1: ", buffer1)
				fmt.Println("Tree2: ", buffer2)

				if buffer1 != buffer2 {
					return false
				}
				/*
					for i := range buffer1 {
						if buffer1[i] != buffer2[i] {
							return false
						}
					}*/

				return true // done
			}

		default:
			stopFlag = true
		}
	}
	return false
}
func testSelect(output chan int) {
	x := 0
	for y := 0; y < 10; y = y + 1 {
		select {
		case output <- x:
			x = x + 1
		default:
			fmt.Println("Default")

		}
	}
}

func main() {

	tree1 := tree.New(10)
	tree2 := tree.New(10)
	fmt.Println("Resut: ", Same(tree1, tree1))
	fmt.Println("Resut: ", Same(tree1, tree2))

	/*
		testCh := make(chan int)
		testSelect(testCh)
		for {
			fmt.Println("select: ", <-testCh)
		}
	*/
}
