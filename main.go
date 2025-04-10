package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/benbjohnson/immutable"

	"github.com/Snawoot/deprank/hasher"
)

type Node struct {
	Name     string
	Children []*Node
}

type NodeSet = immutable.Set[*Node]

var nodeHasher = hasher.NewHasher[*Node]()

func NewNodeSet(values ...*Node) NodeSet {
	return immutable.NewSet(nodeHasher, values...)
}

func (n *Node) prettyWrite(prefix string, visited NodeSet, writer io.Writer) {
	if visited.Has(n) {
		fmt.Fprintf(writer, "%s... recursion goes to Node<Name=%q> ...\n", prefix, n.Name)
		return
	}
	suffix := ""
	if len(n.Children) > 0 {
		suffix = ":"
	}
	fmt.Fprintf(writer, "%sNode<Name=%q>%s\n", prefix, n.Name, suffix)
	visited = visited.Add(n)
	for _, child := range n.Children {
		child.prettyWrite(prefix+"\t", visited, writer)
	}
}

func (n *Node) String() string {
	var b strings.Builder
	n.prettyWrite("", NewNodeSet(), &b)
	return b.String()
}

type stringPair struct {
	a, b string
}

func ReadTree(from io.Reader, rootName *string) (*Node, error) {
	nameIdx := make(map[string]*Node)
	seenEdges := make(map[stringPair]struct{})
	scanner := bufio.NewScanner(from)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		a, b, found := strings.Cut(line, " ")
		if !found {
			return nil, fmt.Errorf("got line without space delimiter: %q", line)
		}

		if _, found := seenEdges[stringPair{a, b}]; found {
			continue
		} else {
			seenEdges[stringPair{a, b}] = struct{}{}
		}

		src, found := nameIdx[a]
		if !found {
			src = &Node{
				Name: a,
			}
			nameIdx[a] = src
		}
		dst, found := nameIdx[b]
		if !found {
			dst = &Node{
				Name: b,
			}
			nameIdx[b] = dst
		}

		if rootName == nil {
			rootName = &a
		}

		src.Children = append(src.Children, dst)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("unable to read graph edges: %w", err)
	}

	if len(nameIdx) == 0 {
		return nil, errors.New("no edges were read")
	}
	rootNode, found := nameIdx[*rootName]
	if !found {
		return nil, errors.New("specified root node was not found")
	}
	return rootNode, nil
}

var (
	rootName *string
)

func init() {
	flag.Func(
		"root",
		"specifies name of root node. Default: first (source) node of the first edge.",
		func(s string) error {
			rootName = &s
			return nil
		},
	)
}

func main() {
	flag.Parse()
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	tree, err := ReadTree(os.Stdin, rootName)
	if err != nil {
		return fmt.Errorf("tree read failed: %w", err)
	}
	fmt.Println(tree)
	return nil
}
