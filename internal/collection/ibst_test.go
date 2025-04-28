package collection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInterval is a simple implementation of the Intervalable interface for testing
type TestInterval struct {
	Low, High int
}

func (i *TestInterval) Name() string {
	return "name"
}

func (i *TestInterval) GetLow() int {
	return i.Low
}

func (i *TestInterval) SetLow(low int) {
	i.Low = low
}

func (i *TestInterval) GetHigh() int {
	return i.High
}

func (i *TestInterval) SetHigh(high int) {
	i.High = high
}

func newTestInterval(low, high int) *TestInterval {
	return &TestInterval{Low: low, High: high}
}

func Test_IBST_Insert(t *testing.T) {
	tree := NewIBST[*TestInterval]()

	// Insert single interval
	tree.Insert(newTestInterval(10, 20))
	if tree.Size() != 1 {
		t.Errorf("Tree should have size 1 after insertion, got %d", tree.Size())
	}
	if tree.IsEmpty() {
		t.Error("Tree should not be empty after insertion")
	}

	// Insert multiple intervals
	tree.Insert(newTestInterval(5, 15))
	tree.Insert(newTestInterval(25, 35))
	if tree.Size() != 3 {
		t.Errorf("Tree should have size 3 after multiple insertions, got %d", tree.Size())
	}

	// Test invalid interval (low > high)
	originalSize := tree.Size()
	tree.Insert(newTestInterval(40, 30)) // This shouldn't be inserted
	if tree.Size() != originalSize {
		t.Error("Invalid interval should not be inserted")
	}
}

func Test_IBST_Delete(t *testing.T) {
	tree := NewIBST[*TestInterval]()

	// Insert intervals
	intervals := []*TestInterval{
		newTestInterval(10, 20),
		newTestInterval(5, 15),
		newTestInterval(25, 35),
	}

	for _, interval := range intervals {
		tree.Insert(interval)
	}

	// Delete existing interval
	tree.Delete(intervals[1]) // Delete (5, 15)
	if tree.Size() != 2 {
		t.Errorf("Tree should have size 2 after deletion, got %d", tree.Size())
	}

	// Try to delete non-existent interval
	originalSize := tree.Size()
	tree.Delete(newTestInterval(40, 50))
	if tree.Size() != originalSize {
		t.Error("Deleting non-existent interval should not change tree size")
	}

	// Delete remaining intervals
	tree.Delete(intervals[0]) // Delete (10, 20)
	tree.Delete(intervals[2]) // Delete (25, 35)
	if !tree.IsEmpty() {
		t.Error("Tree should be empty after deleting all intervals")
	}
}

func Test_IBST_Intersects(t *testing.T) {
	tree := NewIBST[*TestInterval]()

	// Test empty tree
	_, ok := tree.Intersects(newTestInterval(10, 20))
	if ok {
		t.Error("Empty tree should not have intersections")
	}

	// Insert intervals
	tree.Insert(newTestInterval(10, 20))
	tree.Insert(newTestInterval(30, 40))
	tree.Insert(newTestInterval(50, 60))

	// Test intersecting interval
	_, ok = tree.Intersects(newTestInterval(15, 25))
	if !ok {
		t.Error("Should detect intersection with (10, 20)")
	}

	// Test non-intersecting interval
	_, ok = tree.Intersects(newTestInterval(21, 29))
	if ok {
		t.Error("Should not detect intersection with (21, 29)")
	}

	// Test interval that contains an existing interval
	_, ok = tree.Intersects(newTestInterval(5, 25))
	if !ok {
		t.Error("Should detect intersection when new interval contains existing one")
	}

	// Test interval that is contained by an existing interval
	_, ok = tree.Intersects(newTestInterval(12, 18))
	if !ok {
		t.Error("Should detect intersection when new interval is contained by existing one")
	}

	// Test boundary cases
	_, ok = tree.Intersects(newTestInterval(20, 30))
	if !ok {
		t.Error("Should detect intersection at boundary (20)")
	}

	_, ok = tree.Intersects(newTestInterval(40, 50))
	if !ok {
		t.Error("Should detect intersection at boundary (40)")
	}
}

func Test_IBST_Traversal(t *testing.T) {
	assert := assert.New(t)

	tree := NewIBST[*TestInterval]()

	// Insert intervals (not in order)
	toInsert := []*TestInterval{
		newTestInterval(40, 50),
		newTestInterval(0, 10),
		newTestInterval(20, 30),
	}

	for _, interval := range toInsert {
		tree.Insert(interval)
	}

	expectedInOrder := [][]int{{0, 10}, {20, 30}, {40, 50}}

	// Visit all intervals in order
	idx := 0
	for interval := range tree.InOrder() {
		assert.Equal(expectedInOrder[idx][0], interval.GetLow())
		assert.Equal(expectedInOrder[idx][1], interval.GetHigh())
		idx++
	}

	// Check if we have visited all intervals
	assert.Equal(3, idx)

	// Visit all intervals in reverse order
	for interval := range tree.ReverseOrder() {
		idx--
		assert.Equal(expectedInOrder[idx][0], interval.GetLow())
		assert.Equal(expectedInOrder[idx][1], interval.GetHigh())
	}

	// Check if we have visited all intervals
	assert.Equal(0, idx)
}

func Test_IBST_CanUpdateInterval(t *testing.T) {
	tree := NewIBST[*TestInterval]()

	// Test empty tree
	if !tree.CanUpdate(newTestInterval(10, 20), 15, 25) {
		t.Error("Should be able to update interval in empty tree")
	}

	// Insert intervals
	interval1 := newTestInterval(10, 20)
	interval2 := newTestInterval(30, 40)
	interval3 := newTestInterval(50, 60)

	tree.Insert(interval1)
	tree.Insert(interval2)
	tree.Insert(interval3)

	// Test updating to non-intersecting interval
	if !tree.CanUpdate(interval1, 15, 25) {
		t.Error("Should be able to update to non-intersecting interval")
	}

	// Test updating to intersecting interval
	if tree.CanUpdate(interval1, 25, 35) {
		t.Error("Should not be able to update to interval that intersects with others")
	}

	// Test updating to same interval (should always be possible)
	if !tree.CanUpdate(interval2, interval2.GetLow(), interval2.GetHigh()) {
		t.Error("Should be able to update to the same interval")
	}

	// Test single element tree
	singleTree := NewIBST[*TestInterval]()
	singleInterval := newTestInterval(10, 20)
	singleTree.Insert(singleInterval)

	if !singleTree.CanUpdate(singleInterval, 15, 25) {
		t.Error("Should be able to update interval in single-element tree")
	}
}

func Test_IBST_TreeBalancing(t *testing.T) {
	tree := NewIBST[*TestInterval]()

	// Insert elements in ascending order which would create a skewed tree without balancing
	for i := 1; i <= 10; i++ {
		tree.Insert(newTestInterval(i*10, i*10+5))
	}

	// Function to check max height
	var checkHeight func(node *ibstNode[*TestInterval]) int
	checkHeight = func(node *ibstNode[*TestInterval]) int {
		if node == nil {
			return 0
		}
		leftHeight := checkHeight(node.left)
		rightHeight := checkHeight(node.right)
		if leftHeight > rightHeight {
			return leftHeight + 1
		}
		return rightHeight + 1
	}

	// Calculate height
	height := checkHeight(tree.root)

	// For 10 elements, a balanced BST should have height ≈ log₂(10) ≈ 3.32,
	// so height should be at most 5 for a reasonably balanced tree
	if height > 5 {
		t.Errorf("Tree is not well balanced. Height is %d for 10 elements", height)
	}
}

func Test_IBST_Update(t *testing.T) {
	assert := assert.New(t)

	tree := NewIBST[*TestInterval]()

	int0 := newTestInterval(10, 20)
	int1 := newTestInterval(30, 40)

	// Insert some intervals
	tree.Insert(int0)
	tree.Insert(int1)

	// Update intervals
	tree.Update(int0, 5, 15)
	tree.Update(int1, 35, 45)

	assert.Equal(2, tree.Size())

	assert.Equal(5, int0.GetLow())
	assert.Equal(15, int0.GetHigh())
	assert.Equal(35, int1.GetLow())
	assert.Equal(45, int1.GetHigh())

	// Check if the third interval doesn't intersect
	_, ok := tree.Intersects(newTestInterval(16, 34))
	assert.False(ok)
}

func Test_IBST_Clear(t *testing.T) {
	tree := NewIBST[*TestInterval]()

	// Insert some intervals
	tree.Insert(newTestInterval(10, 20))
	tree.Insert(newTestInterval(30, 40))

	// Clear the tree
	tree.Clear()

	if !tree.IsEmpty() {
		t.Error("Tree should be empty after clear")
	}

	if tree.Size() != 0 {
		t.Errorf("Tree size should be 0 after clear, got %d", tree.Size())
	}

	if tree.root != nil {
		t.Error("Tree root should be nil after clear")
	}

	// Check that we can still insert after clearing
	tree.Insert(newTestInterval(50, 60))
	if tree.Size() != 1 {
		t.Errorf("Tree should have size 1 after inserting to cleared tree, got %d", tree.Size())
	}
}

// Test edge cases with overlapping intervals
func Test_IBST_EdgeCases(t *testing.T) {
	tree := NewIBST[*TestInterval]()

	// Insert intervals with same low value but different high values
	tree.Insert(newTestInterval(10, 20))
	tree.Insert(newTestInterval(10, 30))

	if tree.Size() != 2 {
		t.Errorf("Tree should have size 2 after inserting intervals with same low, got %d", tree.Size())
	}

	// Delete one of them
	tree.Delete(newTestInterval(10, 20))

	// Check if the correct one remains
	intervals := tree.GetInOrder()
	if len(intervals) != 1 || intervals[0].GetHigh() != 30 {
		t.Error("Wrong interval was deleted")
	}

	// Test with intervals that have exact same boundaries
	tree.Clear()
	tree.Insert(newTestInterval(10, 20))
	tree.Insert(newTestInterval(10, 20)) // Duplicate

	// BST implementation might or might not allow duplicates
	// This is more of a documentation test than a correctness test
	t.Logf("Tree size after inserting duplicate: %d (BST may or may not allow duplicates)", tree.Size())

	// Test deletion of non-existent intervals with same low but different high
	originalSize := tree.Size()
	tree.Delete(newTestInterval(10, 30)) // Non-existent
	if tree.Size() != originalSize {
		t.Error("Deleting non-existent interval should not change tree size")
	}
}

// Benchmark insertion performance
func Benchmark_IBST_Insert(b *testing.B) {
	tree := NewIBST[*TestInterval]()

	b.ResetTimer()
	for i := 0; b.Loop(); i++ {
		// Insert random intervals
		low := i * 2
		high := low + 10
		tree.Insert(newTestInterval(low, high))
	}
}

// Benchmark intersection checking performance
func Benchmark_IBST_Intersects(b *testing.B) {
	tree := NewIBST[*TestInterval]()

	// Insert some intervals first
	for i := 0; i < 1000; i += 20 {
		tree.Insert(newTestInterval(i, i+10))
	}

	testInterval := newTestInterval(500, 510) // Will hit middle of the tree

	b.ResetTimer()
	for b.Loop() {
		tree.Intersects(testInterval)
	}
}

// Benchmark CanUpdateInterval performance
func Benchmark_IBST_CanUpdateInterval(b *testing.B) {
	tree := NewIBST[*TestInterval]()

	// Insert some intervals first
	intervals := make([]*TestInterval, 0, 1000)
	for i := 0; i < 1000; i += 20 {
		interval := newTestInterval(i, i+10)
		intervals = append(intervals, interval)
		tree.Insert(interval)
	}

	// Use an interval from the middle
	testInterval := intervals[len(intervals)/2]
	newLow := testInterval.GetLow() + 5
	newHigh := testInterval.GetHigh() + 5

	b.ResetTimer()
	for b.Loop() {
		tree.CanUpdate(testInterval, newLow, newHigh)
	}
}
