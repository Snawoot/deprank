package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Snawoot/deprank/graph"
	"github.com/Snawoot/deprank/ranking"
)

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
	dag, err := graph.ReadDAG(os.Stdin, rootName)
	if err != nil {
		return fmt.Errorf("DAG read failed: %w", err)
	}

	r, err := ranking.RankGraph(dag)
	if err != nil {
		return fmt.Errorf("DAG ranking failed: %w", err)
	}
	fmt.Println(r)
	return nil
}
