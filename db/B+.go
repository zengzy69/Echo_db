package db

import (
	"fmt"
	"strings"

	"sort"
)

// 定义 B+树节点结构
type BPlusNode struct {
	isLeaf   bool
	keys     []string
	children []*BPlusNode
	next     *BPlusNode // 叶子节点通过next指针连接
}

// BPlusTree 结构
type BPlusTree struct {
	root   *BPlusNode
	degree int
}

// 创建新的B+树
func NewBPlusTree(degree int) *BPlusTree {
	return &BPlusTree{
		root:   &BPlusNode{isLeaf: true},
		degree: degree,
	}
}

// 插入一个键值对到B+树
func (tree *BPlusTree) Insert(key string) {
	root := tree.root
	if len(root.keys) == tree.degree {
		// 如果根节点已满，进行分裂
		newRoot := &BPlusNode{isLeaf: false}
		newRoot.children = append(newRoot.children, root)
		tree.split(newRoot, 0)
		tree.root = newRoot
	}
	tree.insertNonFull(tree.root, key)
}

// 分裂节点
func (tree *BPlusTree) split(node *BPlusNode, index int) {
	child := node.children[index]
	midIndex := len(child.keys) / 2
	midKey := child.keys[midIndex]

	// 创建新节点并分裂
	newNode := &BPlusNode{isLeaf: child.isLeaf}
	newNode.keys = child.keys[midIndex+1:]
	child.keys = child.keys[:midIndex]

	if child.isLeaf {
		newNode.next = child.next
		child.next = newNode
	} else {
		newNode.children = child.children[midIndex+1:]
		child.children = child.children[:midIndex+1]
	}

	node.keys = append(node.keys[:index], append([]string{midKey}, node.keys[index+1:]...)...)
	node.children = append(node.children[:index+1], append([]*BPlusNode{newNode}, node.children[index+1:]...)...)
}

// 在非满节点插入
func (tree *BPlusTree) insertNonFull(node *BPlusNode, key string) {
	if node.isLeaf {
		node.keys = append(node.keys, key)
		sort.Strings(node.keys)
		return
	}

	i := 0
	for i < len(node.keys) && key > node.keys[i] {
		i++
	}
	if len(node.children[i].keys) == tree.degree {
		tree.split(node, i)
		if key > node.keys[i] {
			i++
		}
	}
	tree.insertNonFull(node.children[i], key)
}

// 查询B+树
func (tree *BPlusTree) Search(key string) bool {
	node := tree.root
	for node != nil {
		i := 0
		for i < len(node.keys) && key > node.keys[i] {
			i++
		}

		if i < len(node.keys) && key == node.keys[i] {
			return true
		}

		if node.isLeaf {
			return false
		}

		node = node.children[i]
	}
	return false
}

// 打印树的结构
func (tree *BPlusTree) PrintTree(node *BPlusNode, level int) {
	if node == nil {
		return
	}
	fmt.Printf("%sLevel %d: ", strings.Repeat(" ", level*2), level)

	for _, key := range node.keys {
		fmt.Printf("%s ", key)
	}
	fmt.Println()
	for _, child := range node.children {
		tree.PrintTree(child, level+1)
	}
}

// 删除节点
func (tree *BPlusTree) Delete(key string) {
	if tree.root == nil {
		return
	}

	// 调用递归删除函数
	tree.delete(tree.root, key)

	// 如果根节点没有键了，并且不是叶子节点，更新根节点
	if len(tree.root.keys) == 0 && !tree.root.isLeaf {
		tree.root = tree.root.children[0]
	}
}

// 删除节点中的一个键
func (tree *BPlusTree) delete(node *BPlusNode, key string) {
	// 叶子节点删除
	if node.isLeaf {
		index := findKeyIndex(node.keys, key)
		if index != -1 {
			// 找到并删除该键
			node.keys = append(node.keys[:index], node.keys[index+1:]...)
		}
		return
	}

	// 内部节点删除
	// 找到该键在内部节点的索引
	index := findKeyIndex(node.keys, key)

	// 如果索引在范围内，且该键在当前节点
	if index < len(node.keys) && key == node.keys[index] {
		// 如果该节点的子节点是叶子节点，直接删除
		if len(node.children[index].keys) > 0 {
			node.keys[index] = tree.getPred(node, index)
			tree.delete(node.children[index], node.keys[index]) // 删除该键
		} else {
			node.keys = append(node.keys[:index], node.keys[index+1:]...)
		}
		return
	}

	// 如果该键在子节点中
	if len(node.children[index].keys) <= tree.degree/2 {
		tree.fix(node, index)
	}

	tree.delete(node.children[index], key)
}

// 获取前驱元素（最大值）
func (tree *BPlusTree) getPred(node *BPlusNode, index int) string {
	current := node.children[index]
	for !current.isLeaf {
		current = current.children[len(current.children)-1]
	}
	return current.keys[len(current.keys)-1]
}

// 获取后继元素（最小值）
func (tree *BPlusTree) getSucc(node *BPlusNode, index int) string {
	current := node.children[index+1]
	for !current.isLeaf {
		current = current.children[0]
	}
	return current.keys[0]
}

// 修复不平衡的节点
func (tree *BPlusTree) fix(node *BPlusNode, index int) {
	// 如果兄弟节点有足够的键，可以借一个
	if index > 0 && len(node.children[index-1].keys) > tree.degree/2 {
		tree.borrowFromPrev(node, index)
	} else if index < len(node.children)-1 && len(node.children[index+1].keys) > tree.degree/2 {
		tree.borrowFromNext(node, index)
	} else {
		// 否则，合并兄弟节点
		if index == len(node.children)-1 {
			tree.merge(node, index-1)
		} else {
			tree.merge(node, index)
		}
	}
}

// 从前一个兄弟节点借一个元素
func (tree *BPlusTree) borrowFromPrev(node *BPlusNode, index int) {
	child := node.children[index]
	sibling := node.children[index-1]

	// 将父节点的键移动到子节点
	child.keys = append([]string{node.keys[index-1]}, child.keys...)
	node.keys[index-1] = sibling.keys[len(sibling.keys)-1]

	// 如果兄弟节点不是叶子节点，父节点的子节点需要移动
	if !child.isLeaf {
		child.children = append([]*BPlusNode{sibling.children[len(sibling.children)-1]}, child.children...)
		sibling.children = sibling.children[:len(sibling.children)-1]
	}

	// 删除兄弟节点的键
	sibling.keys = sibling.keys[:len(sibling.keys)-1]
}

// 从后一个兄弟节点借一个元素
func (tree *BPlusTree) borrowFromNext(node *BPlusNode, index int) {
	child := node.children[index]
	sibling := node.children[index+1]

	// 将父节点的键移动到子节点
	child.keys = append(child.keys, node.keys[index])
	node.keys[index] = sibling.keys[0]

	// 如果兄弟节点不是叶子节点，父节点的子节点需要移动
	if !child.isLeaf {
		child.children = append(child.children, sibling.children[0])
		sibling.children = sibling.children[1:]
	}

	// 删除兄弟节点的键
	sibling.keys = sibling.keys[1:]
}

// 合并节点
func (tree *BPlusTree) merge(node *BPlusNode, index int) {
	left := node.children[index]
	right := node.children[index+1]

	// 将父节点的键移到左子节点
	left.keys = append(left.keys, node.keys[index])
	left.keys = append(left.keys, right.keys...)

	// 如果右子节点不是叶子节点，合并子节点
	if !left.isLeaf {
		left.children = append(left.children, right.children...)
	}

	// 删除父节点的键和右子节点
	node.keys = append(node.keys[:index], node.keys[index+1:]...)
	node.children = append(node.children[:index+1], node.children[index+2:]...)
}

// 查找键在节点中的位置
func findKeyIndex(keys []string, key string) int {
	for i, k := range keys {
		if key == k {
			return i
		}
		if key < k {
			break
		}
	}
	return -1
}
