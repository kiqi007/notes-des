package main

// @Des: des 命令行入口

import (
	"flag"
	"fmt"
	xast "github.com/kiqi007/notes-des"
	"os"
	"strings"
)

var descript = ` - 此应用程序将扫描目录下的所有文件，查找到特定注释("@Des: {describe text}" as default)，并构建成易读的结构并输出
   标准输出格式:                        	' filename.go - {describe text}'
   无需备注的文件( @Des: -):             	' filename.go'
   已经废弃的文件( @Des: [d|delete|废弃]):    	' filename.go(del)'
`

func main() {
	flag.Usage = Usage
	flag.Parse()

	dir := "./..."
	if l := len(os.Args); l != 1 {
		if !strings.HasPrefix(os.Args[l-1], "-") {
			dir = os.Args[l-1]
		}
	}

	r := NewDesVisitor(xast.BaseParser().WithDir(dir, *excludeKey))
	r.Do()
	parseRs := r.Output()

	switch *out {
	case "tree":
		TreeOut{}.Output(r.Dir, parseRs)
	case "trie":
		TrieOut{}.Output(r.Dir, parseRs)
	default:
		StdOut{}.Output(r.Dir, parseRs)
	}
}

func Usage() {
	fmt.Fprintf(os.Stderr, "usage: des [options] dir\n")

	fmt.Fprintf(os.Stderr, "%s\n", descript)

	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr)
	os.Exit(0)
}
