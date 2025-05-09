package collection

type ibstNode[T IBSTItem] struct {
	item        T
	max         int
	left, right *ibstNode[T]
	height      int
}

func (n *ibstNode[T]) getLow() int {
	return n.item.GetLow()
}

func (n *ibstNode[T]) getHigh() int {
	return n.item.GetHigh()
}

// updateHeight recalculates the height of a node based on its children
func (n *ibstNode[T]) updateHeight() {
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
func (n *ibstNode[T]) updateMax() {
	n.max = n.item.GetHigh()

	if n.left != nil && n.left.max > n.max {
		n.max = n.left.max
	}

	if n.right != nil && n.right.max > n.max {
		n.max = n.right.max
	}
}

// balanceFactor returns the balance factor of a node
func (n *ibstNode[T]) balanceFactor() int {
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
func (n *ibstNode[T]) rotateRight() *ibstNode[T] {
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
func (n *ibstNode[T]) rotateLeft() *ibstNode[T] {
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
func (n *ibstNode[T]) findMin() *ibstNode[T] {
	for n.left != nil {
		return n.left
	}
	return n
}
