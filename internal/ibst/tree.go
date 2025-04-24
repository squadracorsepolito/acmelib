// Package ibst contains the implementation of the interval binary search tree.
package ibst

import (
	"iter"

	"github.com/squadracorsepolito/acmelib/internal/stringer"
)

// Intervalable is an interface for intervalable items
// to be stored in the [Tree].
type Intervalable interface {
	Name() string
	GetLow() int
	SetLow(int)
	GetHigh() int
	SetHigh(int)
}

// Tree is a binary search tree that stores intervals.
type Tree[T Intervalable] struct {
	root *node[T]
	size int
}

// NewTree returns a new [Tree].
func NewTree[T Intervalable]() *Tree[T] {
	return &Tree[T]{
		root: nil,
		size: 0,
	}
}

// insertNode recursively inserts a new interval into the tree and balances it
func (t *Tree[T]) insertNode(root *node[T], item T) *node[T] {
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
	if low < root.getLow() {
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
		if item.GetLow() > root.left.getLow() {
			root.left = root.left.rotateLeft()
			return root.rotateRight()
		}
		// Left-Left case
		return root.rotateRight()
	}

	// Right heavy
	if balance < -1 {
		// Right-Left case
		if item.GetLow() < root.right.getLow() {
			root.right = root.right.rotateRight()
			return root.rotateLeft()
		}
		// Right-Right case
		return root.rotateLeft()
	}

	return root
}

// Insert adds a new intervalable item to the tree.
func (t *Tree[T]) Insert(item T) {
	if item.GetLow() > item.GetHigh() {
		// Invalid interval, silently ignore
		return
	}
	t.root = t.insertNode(t.root, item)
}

// deleteNode recursively deletes a node with the given interval
func (t *Tree[T]) deleteNode(root *node[T], item T) *node[T] {
	if root == nil {
		return nil
	}

	low := item.GetLow()
	high := item.GetHigh()

	// First locate the node to delete
	if low < root.getLow() {
		root.left = t.deleteNode(root.left, item)
	} else if low > root.getLow() {
		root.right = t.deleteNode(root.right, item)
	} else if high != root.getHigh() {
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
func (t *Tree[T]) Delete(item T) {
	t.root = t.deleteNode(t.root, item)
}

// Size returns the number of intervals in the tree.
func (t *Tree[T]) Size() int {
	return t.size
}

// IsEmpty returns true if the tree is empty.
func (t *Tree[T]) IsEmpty() bool {
	return t.size == 0
}

// intersectsNode recursively checks for any intersection with the given interval
func (t *Tree[T]) intersectsNode(root *node[T], low, high int) bool {
	if root == nil {
		return false
	}

	// Early pruning: if root.max < low, then no interval in this subtree can overlap
	if root.max < low {
		return false
	}

	// If current node overlaps, return true immediately
	if root.getLow() <= high && low <= root.getHigh() {
		return true
	}

	// More efficient traversal based on BST properties
	// If low value is less than root's low, we need to check left subtree
	if low < root.getLow() && t.intersectsNode(root.left, low, high) {
		return true
	}

	// Always check right subtree (intervals with same low value might be on the right)
	return t.intersectsNode(root.right, low, high)
}

// Intersects states whether the given intervalable item intersects any already in the tree.
func (t *Tree[T]) Intersects(item T) bool {
	if t.IsEmpty() {
		return false
	}
	return t.intersectsNode(t.root, item.GetLow(), item.GetHigh())
}

// inOrderTraversal recursively traverses the tree in ascending order
func (t *Tree[T]) inOrderTraversal(root *node[T], yield func(T) bool) bool {
	if root == nil {
		return true
	}

	// Visit left subtree
	if !t.inOrderTraversal(root.left, yield) {
		return false
	}

	// Visit current node
	if !yield(root.item) {
		return false
	}

	// Visit right subtree
	return t.inOrderTraversal(root.right, yield)
}

// InOrder returns an iterator over all intervals in the tree in ascending order by low value.
func (t *Tree[T]) InOrder() iter.Seq[T] {
	return func(yield func(T) bool) {
		t.inOrderTraversal(t.root, yield)
	}
}

// GetInOrder returns all intervals in the tree in ascending order by low value.
func (t *Tree[T]) GetInOrder() []T {
	result := make([]T, 0, t.size)

	for item := range t.InOrder() {
		result = append(result, item)
	}

	return result
}

// reverseOrderTraversal recursively traverses the tree in descending order
func (t *Tree[T]) reverseOrderTraversal(root *node[T], yield func(T) bool) bool {
	if root == nil {
		return true
	}

	// Visit right subtree
	if !t.reverseOrderTraversal(root.right, yield) {
		return false
	}

	// Visit current node
	if !yield(root.item) {
		return false
	}

	// Visit left subtree
	return t.reverseOrderTraversal(root.left, yield)
}

// ReverseOrder returns an iterator over all intervals in the tree in descending order by low value.
func (t *Tree[T]) ReverseOrder() iter.Seq[T] {
	return func(yield func(T) bool) {
		t.reverseOrderTraversal(t.root, yield)
	}
}

// GetReverseOrder returns all intervals in the tree in descending order by low value.
func (t *Tree[T]) GetReverseOrder() []T {
	result := make([]T, 0, t.size)

	for item := range t.ReverseOrder() {
		result = append(result, item)
	}

	return result
}

// checkOtherIntervals a helper function to check if an interval intersects with any other interval
// except the one we are skipping.
func (t *Tree[T]) checkOtherIntervals(node *node[T], low, high int, skipItem T) bool {
	if node == nil {
		return false
	}

	// Early pruning: if node.max < low, no interval in this subtree can overlap
	if node.max < low {
		return false
	}

	// Check if current node overlaps with the new interval
	// Skip the node if it's the one we're updating
	isSameItem := node.getLow() == skipItem.GetLow() && node.getHigh() == skipItem.GetHigh()
	if !isSameItem && node.getLow() <= high && low <= node.getHigh() {
		return true // Found an intersection
	}

	// If low value is less than node's low, we need to check left subtree
	if low < node.getLow() && t.checkOtherIntervals(node.left, low, high, skipItem) {
		return true
	}

	// Check right subtree if needed
	return t.checkOtherIntervals(node.right, low, high, skipItem)
}

// CanUpdate checks if the given intervalable item can be updated
// without intersecting with any other interval.
func (t *Tree[T]) CanUpdate(item T, newLow, newHigh int) bool {
	// If tree is empty or has only one item (the one we're updating)
	if t.size <= 1 {
		return true
	}

	return !t.checkOtherIntervals(t.root, newLow, newHigh, item)
}

// Update updates an intervalable item in the tree with the new low and high values.
func (t *Tree[T]) Update(item T, newLow, newHigh int) {
	t.Delete(item)
	item.SetLow(newLow)
	item.SetHigh(newHigh)
	t.Insert(item)
}

// Clear removes all intervals from the tree
func (t *Tree[T]) Clear() {
	t.root = nil
	t.size = 0
}

func (t *Tree[T]) stringify(s *stringer.Stringer, node *node[T]) {
	if node == nil {
		return
	}

	// Traverse right subtree first (will appear at the top)
	s.Indent()
	t.stringify(s, node.right)
	s.Unindent()

	// Current node
	s.Write("[%d, %d] %s (max: %d)\n", node.getLow(), node.getHigh(), node.item.Name(), node.max)

	// Traverse left subtree
	s.Indent()
	t.stringify(s, node.left)
	s.Unindent()
}

// Stringify writes a string representation of the tree into
// a [stringer.Stringer].
func (t *Tree[T]) Stringify(s *stringer.Stringer) {
	s.Write("size: %d\n", t.size)
	t.stringify(s, t.root)
}

func (t *Tree[T]) String() string {
	s := stringer.New()
	s.Write("interval_bst\n")
	t.Stringify(s)
	return s.String()
}
