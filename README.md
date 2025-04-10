# deprank

Topological ranks from Go dependency graph

## Example

```
user@ws:~/src/deprank> go mod graph | deprank 
Ranking<len=4>:
	Rank 0:
		golang.org/x/tools@v0.1.10
		gopkg.in/yaml.v3@v3.0.1
		toolchain@go1.24.2
		github.com/pmezard/go-difflib@v1.0.0
		golang.org/x/xerrors@v0.0.0-20200804184101-5ec99f83aff1
		golang.org/x/mod@v0.6.0-dev.0.20220106191415-9b9b3d81d5e3
		golang.org/x/sys@v0.0.0-20211019181941-9d821ace8654
		github.com/davecgh/go-spew@v1.1.1
		github.com/stretchr/testify@v1.8.1
	Rank 1:
		golang.org/x/exp@v0.0.0-20220518171630-0b5c67f07fdf
		github.com/dolthub/maphash@v0.1.0
		go@1.24.2
	Rank 2:
		github.com/benbjohnson/immutable@v0.4.3
	Rank 3:
		github.com/Snawoot/deprank

```
