package fuzzy

import "strings"

type Stringer interface {
	String() string
}

type Node struct {
	value    Stringer
	children map[int]*Node
}

func newNode(v Stringer) *Node {
	return &Node{
		value:    v,
		children: make(map[int]*Node),
	}
}

type Tree struct {
	root *Node
}

func NewTree(toInsert []Stringer) *Tree {
	tree := &Tree{}
	for _, w := range toInsert {
		tree.insert(w)
	}
	return tree
}

func (t *Tree) insert(val Stringer) {
	if t.root == nil {
		t.root = newNode(val)
		return
	}

	u := t.root
	for u != nil {
		k := levenshtein(u.value.String(), val.String())
		if k == 0 {
			return
		}

		v, exists := u.children[k]
		if !exists {
			v = newNode(val)
			u.children[k] = v
			return
		}
		u = v
	}

	return
}

func (tree *Tree) Lookup(s string) []Stringer {
	const tol = 5
	if tree.root == nil {
		return nil
	}

	var ans [tol + 1][]Stringer

	stack := []*Node{tree.root}
	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		d := levenshtein(s, u.value.String())
		if d <= tol {
			ans[d] = append(ans[d], u.value)
		}

		for dist, child := range u.children {
			if dist > d-tol && dist < d+tol {
				stack = append(stack, child)
			}
		}
	}

	for i := 1; i <= tol; i++ {
		ans[0] = append(ans[0], ans[i]...)
	}

	return ans[0][:min(5, len(ans[0]))]
}

func levenshtein(s, t string) int {
	m, n := len(s), len(t)
	v0, v1 := make([]int, n+1), make([]int, n+1)

	for i := range v0 {
		v0[i] = i
	}

	for i := range m {
		v1[0] = i + 1

		for j := range n {
			deletionCost := v0[j+1] + 1
			insertionCost := v1[j] + 1
			substitutionCost := v0[j] + 1
			if strings.EqualFold(string(s[i]), string(t[j])) {
				substitutionCost = v0[j]
			}
			v1[j+1] = min(deletionCost, insertionCost, substitutionCost)
		}
		copy(v0, v1)
	}

	return v0[n]
}
