package fuzzy

type Node struct {
	word     string
	children map[int]*Node
}

func newNode(s string) *Node {
	return &Node{
		word:     s,
		children: make(map[int]*Node),
	}
}

type Tree struct {
	root *Node
}

func newTree(toInsert ...string) *Tree {
	tree := &Tree{}
	for _, w := range toInsert {
		tree.insert(w)
	}
	return tree
}

func (t *Tree) insert(word string) *Node {
	if t.root == nil {
		t.root = newNode(word)
		return t.root
	}

	u := t.root
	for u != nil {
		k := levenshtein(u.word, word)

		if k == 0 {
			return u
		}

		v, exists := u.children[k]
		if !exists {
			v = newNode(word)

			u.children[k] = v

			return v
		}
		u = v
	}

	return nil
}

func (tree *Tree) Lookup(word string) string {
	if tree.root == nil {
		return ""
	}

	stack := []*Node{tree.root}
	bestWord := ""
	bestDist := 100

	for len(stack) > 0 {
		u := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		dU := levenshtein(word, u.word)

		if dU < bestDist {
			bestWord = u.word
			bestDist = dU
		}

		for dist, child := range u.children {
			if dist <= bestDist {
				stack = append(stack, child)
			}
		}
	}

	return bestWord
}

func levenshtein(s, t string) int {
	m := len(s)
	n := len(t)

	v0 := make([]int, n+1)
	v1 := make([]int, n+1)

	for i := range v0 {
		v0[i] = i
	}

	for i := range m {
		v1[0] = i + 1

		for j := range n {
			deletionCost := v0[j+1] + 1
			insertionCost := v1[j] + 1
			substitutionCost := v0[j] + 1
			if s[i] == t[j] {
				substitutionCost = v0[j]
			}

			v1[j+1] = min(deletionCost, insertionCost, substitutionCost)
		}

		copy(v0, v1)
	}

	return v0[n]
}
