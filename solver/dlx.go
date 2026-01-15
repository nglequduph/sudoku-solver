package solver

type Node struct {
	left, right, up, down *Node
	col                   *Node
	rowID, colID          int
	size                  int
}

type DLX struct {
	header *Node
	result []int
	steps  int
}

func NewDLX(numCols int) *DLX {
	root := &Node{rowID: -1, colID: -1}
	root.left = root
	root.right = root
	root.up = root
	root.down = root
	last := root
	for i := 0; i < numCols; i++ {
		col := &Node{rowID: -1, colID: i}
		col.up = col
		col.down = col
		col.left = last
		col.right = root
		last.right = col
		root.left = col
		last = col
	}
	return &DLX{header: root}
}

func (dlx *DLX) AddRow(rowID int, cols []int) {
	var first *Node
	for _, colIdx := range cols {
		c := dlx.header.right
		for c.colID != colIdx {
			c = c.right
		}
		newNode := &Node{rowID: rowID, colID: colIdx, col: c, up: c.up, down: c}
		c.up.down = newNode
		c.up = newNode
		c.size++
		if first == nil {
			first = newNode
			newNode.left = newNode
			newNode.right = newNode
		} else {
			newNode.left = first.left
			newNode.right = first
			first.left.right = newNode
			first.left = newNode
		}
	}
}

func (dlx *DLX) cover(c *Node) {
	c.right.left = c.left
	c.left.right = c.right
	for i := c.down; i != c; i = i.down {
		for j := i.right; j != i; j = j.right {
			j.down.up = j.up
			j.up.down = j.down
			j.col.size--
		}
	}
}

func (dlx *DLX) uncover(c *Node) {
	for i := c.up; i != c; i = i.up {
		for j := i.left; j != i; j = j.left {
			j.col.size++
			j.up.down = j
			j.down.up = j
		}
	}
	c.left.right = c
	c.right.left = c
}

func (dlx *DLX) Solve() bool {
	dlx.steps++
	if dlx.header.right == dlx.header {
		return true
	}
	c := dlx.header.right
	for temp := c.right; temp != dlx.header; temp = temp.right {
		if temp.size < c.size {
			c = temp
		}
	}
	if c.size == 0 {
		return false
	}
	dlx.cover(c)
	for r := c.down; r != c; r = r.down {
		dlx.result = append(dlx.result, r.rowID)
		for j := r.right; j != r; j = j.right {
			dlx.cover(j.col)
		}
		if dlx.Solve() {
			return true
		}
		for j := r.left; j != r; j = j.left {
			dlx.uncover(j.col)
		}
		dlx.result = dlx.result[:len(dlx.result)-1]
	}
	dlx.uncover(c)
	return false
}

func SolveSudoku(grid [9][9]int) ([9][9]int, bool, int) {
	dlx := NewDLX(324)
	for r := 0; r < 9; r++ {
		for c := 0; c < 9; c++ {
			for v := 1; v <= 9; v++ {
				if grid[r][c] == 0 || grid[r][c] == v {
					rowID := r*81 + c*9 + (v - 1)
					b := (r/3)*3 + c/3
					cols := []int{r*9 + c, 81 + r*9 + (v - 1), 162 + c*9 + (v - 1), 243 + b*9 + (v - 1)}
					dlx.AddRow(rowID, cols)
				}
			}
		}
	}
	if dlx.Solve() {
		resGrid := [9][9]int{}
		for _, rowID := range dlx.result {
			r := rowID / 81
			c := (rowID % 81) / 9
			v := (rowID % 9) + 1
			resGrid[r][c] = v
		}
		return resGrid, true, dlx.steps
	}
	return grid, false, dlx.steps
}
