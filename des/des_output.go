package main

// @Des: 将解析结果解析成特定格式并输出 (StdOut/TreeOut/TrieOut)

import (
	"flag"
	"fmt"
	"github.com/kiqi007/notes-des/trie"
	"io"
	"log"
	"os"
	"strings"

	"github.com/xlab/treeprint"
)

var out = flag.String("out", "normal", "输出模式，可选值：normal, tree, trie")
var outPath = flag.String("out_path", "", "输出目录，若未设置则输出到控制台")
var nf = flag.String("empty", "", "如果文件无匹配注释，则结果填入此字符串")

type IDesOut interface {
	Output(parseDir string, parseRs map[string][]string)
}

type StdOut struct{}

func (s StdOut) Output(parseDir string, parseRs []ParseRs) {
	dir := strings.TrimRight(parseDir, "./")
	dir = dir[:strings.LastIndexByte(dir, '/')+1]

	wr := GetOutPath()
	for _, rs := range parseRs {
		path := strings.TrimPrefix(rs.Key, dir)
		_, _ = fmt.Fprintf(wr, "%s\n", formatValue(path, rs.Value))
	}
}

type TreeOut struct {
	tree treeprint.Tree
}

func (s TreeOut) Output(parseDir string, parseRs []ParseRs) {
	if s.tree == nil {
		s.tree = treeprint.New()
	}

	dir := strings.TrimRight(parseDir, "./")
	dir = dir[:strings.LastIndexByte(dir, '/')+1]

	for _, rs := range parseRs {
		path := strings.TrimPrefix(rs.Key, dir)
		s.addNode(s.tree, strings.Split(path, "/"), rs.Value)
	}

	_, _ = fmt.Fprintf(GetOutPath(), "%s", s.tree.String())
}

func (s TreeOut) addNode(tree treeprint.Tree, node []string, value []string) treeprint.Tree {
	if len(node) == 0 {
		return tree
	}

	if t := tree.FindByValue(node[0]); t != nil { // 节点已存在
		return s.addNode(t, node[1:], value)
	}

	if len(node) == 1 { // 文件处理
		return tree.AddNode(formatValue(node[0], value))
	}

	return s.addNode(tree.AddBranch(node[0]), node[1:], value) // 节点不存在
}

// 压缩树实现
type TrieOut struct {
	tree trie.Tree
}

func (s TrieOut) Output(parseDir string, parseRs []ParseRs) {
	if s.tree == nil {
		s.tree = trie.New()
	}

	dir := strings.TrimRight(parseDir, "./")
	dir = dir[:strings.LastIndexByte(dir, '/')+1]

	for _, rs := range parseRs {
		path := strings.TrimPrefix(rs.Key, dir)
		s.addNode(s.tree, strings.Split(path, "/"), rs.Value)
	}

	_, _ = fmt.Fprintf(GetOutPath(), "%s", s.tree.Tidy("/").String())
}

func (s TrieOut) addNode(tree trie.Tree, node []string, value []string) trie.Tree {
	if len(node) == 0 {
		return tree
	}

	if t := tree.FindByValue(node[0]); t != nil { // 节点已存在
		return s.addNode(t, node[1:], value)
	}

	if len(node) == 1 { // 文件处理
		return tree.AddNode(formatValue(node[0], value))
	}

	return s.addNode(tree.AddBranch(node[0]), node[1:], value) // 节点不存在
}

func GetOutPath() io.Writer {
	if len(*outPath) == 0 {
		return os.Stdout
	}

	create, err := os.Create(*outPath)
	if err != nil {
		log.Panic(fmt.Sprintf("[error] failed when create file, path=%s, err=%+v", *outPath, err))
	}
	return create
}

func formatValue(node string, value []string) string {
	v := ""
	switch value[0] {
	case "-":
		v = fmt.Sprintf("%s", node)
	case "d", "delete", "废弃":
		v = fmt.Sprintf("%s(del)", node)
	case "":
		v = fmt.Sprintf("%s - [%s]", node, *nf)
	default:
		v = fmt.Sprintf("%s - [%s]", node, strings.Join(value, " / "))
	}
	return v
}
