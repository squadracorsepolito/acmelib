// Package internal contains internal types and functions used by the acmelib library.
package internal

import (
	"fmt"
	"strings"
)

// Intervalable is an interface for intervalable items
// to be stored in the [IntervalBST].
type Intervalable interface {
	GetLow() int
	GetHigh() int
}

type node[T Intervalable] struct {
	item        T
	max         int
	left, right *node[T]
	height      int
}

// updateHeight recalculates the height of a node based on its children
func (n *node[T]) updateHeight() {
	leftH := 0
	if n.left != nil {
		leftH = n.left.height
	}

	rightH := 0
	if n.right != nil {
		rightH = n.right.height
	}

	n.height = 1 + max(leftH, rightH)
}

// updateMax recalculates the max value of a node based on itself and its children
func (n *node[T]) updateMax() {
	n.max = n.item.GetHigh()

	if n.left != nil && n.left.max > n.max {
		n.max = n.left.max
	}

	if n.right != nil && n.right.max > n.max {
		n.max = n.right.max
	}
}

// balanceFactor returns the balance factor of a node
func (n *node[T]) balanceFactor() int {
	leftH := 0
	if n.left != nil {
		leftH = n.left.height
	}

	rightH := 0
	if n.right != nil {
		rightH = n.right.height
	}

	return leftH - rightH
}

// rotateRight performs a right rotation on the given node
func (n *node[T]) rotateRight() *node[T] {
	leftNode := n.left
	leftRightNode := leftNode.right

	// Perform rotation
	leftNode.right = n
	n.left = leftRightNode

	// Update heights and max values
	n.updateHeight()
	n.updateMax()
	leftNode.updateHeight()
	leftNode.updateMax()

	return leftNode
}

// rotateLeft performs a left rotation on the given node
func (n *node[T]) rotateLeft() *node[T] {
	rightNode := n.right
	rightLeftNode := rightNode.left

	// Perform rotation
	rightNode.left = n
	n.right = rightLeftNode

	// Update heights and max values
	n.updateHeight()
	n.updateMax()
	rightNode.updateHeight()
	rightNode.updateMax()

	return rightNode
}

// findMin returns the node with the minimum low value in the subtree
func (n *node[T]) findMin() *node[T] {
	for n.left != nil {
		return n.left
	}
	return n
}

// IntervalBST is a binary search tree that stores intervals.
type IntervalBST[T Intervalable] struct {
	root *node[T]
	size int
}

// NewIntervalBST returns a new [IntervalBST].
func NewIntervalBST[T Intervalable]() *IntervalBST[T] {
	return &IntervalBST[T]{
		root: nil,
		size: 0,
	}
}

// insertNode recursively inserts a new interval into the tree and balances it
func (t *IntervalBST[T]) insertNode(root *node[T], item T) *node[T] {
	// Standard BST insertion
	if root == nil {
		t.size++
		return &node[T]{
			item:   item,
			max:    item.GetHigh(),
			height: 1,
		}
	}

	low := item.GetLow()
	if low < root.item.GetLow() {
		root.left = t.insertNode(root.left, item)
	} else {
		root.right = t.insertNode(root.right, item)
	}

	// Update height and max value
	root.updateHeight()
	root.updateMax()

	// Get balance factor
	balance := root.balanceFactor()

	// Left heavy
	if balance > 1 {
		// Left-Right case
		if item.GetLow() > root.left.item.GetLow() {
			root.left = root.left.rotateLeft()
			return root.rotateRight()
		}
		// Left-Left case
		return root.rotateRight()
	}

	// Right heavy
	if balance < -1 {
		// Right-Left case
		if item.GetLow() < root.right.item.GetLow() {
			root.right = root.right.rotateRight()
			return root.rotateLeft()
		}
		// Right-Right case
		return root.rotateLeft()
	}

	return root
}

// Insert adds a new intervalable item to the tree.
func (t *IntervalBST[T]) Insert(item T) {
	if item.GetLow() > item.GetHigh() {
		// Invalid interval, silently ignore or could return an error
		return
	}
	t.root = t.insertNode(t.root, item)
}

// deleteNode recursively deletes a node with the given interval
func (t *IntervalBST[T]) deleteNode(root *node[T], item T) *node[T] {
	if root == nil {
		return nil
	}

	low := item.GetLow()
	high := item.GetHigh()

	// First locate the node to delete
	if low < root.item.GetLow() {
		root.left = t.deleteNode(root.left, item)
	} else if low > root.item.GetLow() {
		root.right = t.deleteNode(root.right, item)
	} else if high != root.item.GetHigh() {
		// Same low but different high, continue search
		root.right = t.deleteNode(root.right, item)
	} else {
		// Found the node to delete
		t.size--

		// Case with at most one child
		if root.left == nil {
			return root.right
		} else if root.right == nil {
			return root.left
		}

		// Node with two children: Get the inorder successor (smallest in the right subtree)
		successor := root.right.findMin()
		root.item = successor.item

		// Delete the inorder successor
		root.right = t.deleteNode(root.right, successor.item)
	}

	// Update height and max value
	root.updateHeight()
	root.updateMax()

	// Get balance factor
	balance := root.balanceFactor()

	// Left heavy
	if balance > 1 {
		// Left-Right case
		leftBalance := 0
		if root.left != nil {
			leftBalance = root.left.balanceFactor()
		}

		if leftBalance < 0 {
			root.left = root.left.rotateLeft()
			return root.rotateRight()
		}
		// Left-Left case
		return root.rotateRight()
	}

	// Right heavy
	if balance < -1 {
		// Right-Left case
		rightBalance := 0
		if root.right != nil {
			rightBalance = root.right.balanceFactor()
		}

		if rightBalance > 0 {
			root.right = root.right.rotateRight()
			return root.rotateLeft()
		}
		// Right-Right case
		return root.rotateLeft()
	}

	return root
}

// Delete removes an intervalable item from the tree.
func (t *IntervalBST[T]) Delete(item T) {
	t.root = t.deleteNode(t.root, item)
}

// Size returns the number of intervals in the tree.
func (t *IntervalBST[T]) Size() int {
	return t.size
}

// IsEmpty returns true if the tree is empty.
func (t *IntervalBST[T]) IsEmpty() bool {
	return t.size == 0
}

// intersectsNode recursively checks for any intersection with the given interval
func (t *IntervalBST[T]) intersectsNode(root *node[T], low, high int) bool {
	if root == nil {
		return false
	}

	// Early pruning: if root.max < low, then no interval in this subtree can overlap
	if root.max < low {
		return false
	}

	// If current node overlaps, return true immediately
	if root.item.GetLow() <= high && low <= root.item.GetHigh() {
		return true
	}

	// More efficient traversal based on BST properties
	// If low value is less than root's low, we need to check left subtree
	if low < root.item.GetLow() && t.intersectsNode(root.left, low, high) {
		return true
	}

	// Always check right subtree (intervals with same low value might be on the right)
	return t.intersectsNode(root.right, low, high)
}

// Intersects states whether the given intervalable item intersects any already in the tree.
func (t *IntervalBST[T]) Intersects(item T) bool {
	if t.IsEmpty() {
		return false
	}
	return t.intersectsNode(t.root, item.GetLow(), item.GetHigh())
}

// inOrderTraversal performs an in-order traversal and collects items
func (t *IntervalBST[T]) inOrderTraversal(root *node[T], result *[]T) {
	if root == nil {
		return
	}

	t.inOrderTraversal(root.left, result)
	*result = append(*result, root.item)
	t.inOrderTraversal(root.right, result)
}

// GetAllIntervals returns all intervals in the tree in ascending order by low value
func (t *IntervalBST[T]) GetAllIntervals() []T {
	result := make([]T, 0, t.size)
	t.inOrderTraversal(t.root, &result)
	return result
}

// checkOtherIntervals a helper function to check if an interval intersects with any other interval
// except the one we are skipping.
func (t *IntervalBST[T]) checkOtherIntervals(node *node[T], low, high int, skipItem T) bool {
	if node == nil {
		return false
	}

	// Early pruning: if node.max < low, no interval in this subtree can overlap
	if node.max < low {
		return false
	}

	// Check if current node overlaps with the new interval
	// Skip the node if it's the one we're updating
	isSameItem := node.item.GetLow() == skipItem.GetLow() && node.item.GetHigh() == skipItem.GetHigh()
	if !isSameItem && node.item.GetLow() <= high && low <= node.item.GetHigh() {
		return true // Found an intersection
	}

	// If low value is less than node's low, we need to check left subtree
	if low < node.item.GetLow() && t.checkOtherIntervals(node.left, low, high, skipItem) {
		return true
	}

	// Check right subtree if needed
	return t.checkOtherIntervals(node.right, low, high, skipItem)
}

// CanUpdateInterval checks if the given intervalable item can be updated.
func (t *IntervalBST[T]) CanUpdateInterval(item T, newLow, newHigh int) bool {
	// If tree is empty or has only one item (the one we're updating)
	if t.size <= 1 {
		return true
	}

	return !t.checkOtherIntervals(t.root, newLow, newHigh, item)
}

// Clear removes all intervals from the tree
func (t *IntervalBST[T]) Clear() {
	t.root = nil
	t.size = 0
}

func (t *IntervalBST[T]) stringify(sb *strings.Builder, node *node[T], level int) {
	if node == nil {
		return
	}

	// Traverse right subtree first (will appear at the top)
	t.stringify(sb, node.right, level+1)

	// Current node with proper indentation
	indent := strings.Repeat("\t", level)
	sb.WriteString(fmt.Sprintf("%s[%d, %d] (max: %d)\n",
		indent, node.item.GetLow(), node.item.GetHigh(), node.max))

	// Traverse left subtree
	t.stringify(sb, node.left, level+1)
}

func (t *IntervalBST[T]) String() string {
	if t.root == nil {
		return "Empty IntervalBST"
	}

	sb := new(strings.Builder)

	sb.WriteString(fmt.Sprintf("IntervalBST with %d intervals:\n", t.size))
	t.stringify(sb, t.root, 0)

	return sb.String()
}
