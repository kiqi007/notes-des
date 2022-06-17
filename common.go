package xast

// @Des: ast可复用的共有抽象部分

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	debug = flag.Bool("debug", false, "是否输出debug日志")
)

type Parser struct {
	ast.Visitor
	Fset *token.FileSet // 对应当前解析的文件目录

	Debug    bool
	Dir      string
	Packages []string
}

// BaseParser 初始化一个*Parser
func BaseParser() *Parser {
	p := &Parser{
		Debug: *debug,
	}
	return p
}

// WithDir : 允许通过传入参数从而覆盖命令行参数
func (pvv *Parser) WithDir(drs string, eKeys string) *Parser {
	if len(drs) == 0 {
		drs = "./..."
	}

	path, _ := filepath.Abs(drs)
	pvv.Dir = path
	pvv.Packages = excludeKeys(listPackages(pvv.Debug, pvv.Dir), eKeys)
	return pvv
}

func (pvv *Parser) WithVisitor(p ast.Visitor) *Parser {
	pvv.Visitor = p
	return pvv
}

func (pvv *Parser) Do() {
	log.Printf("[info] ast codes start, dir=%s\n", pvv.Dir)

	for _, pkg := range pvv.Packages {
		pvv.Fset = token.NewFileSet()
		fpkgs, err := parser.ParseDir(pvv.Fset, pkg, nil, parser.ParseComments)
		if err != nil {
			log.Panic("[panic] pkg parse panic, pkg=", pkg, ", err=", err)
		}

		for _, f := range fpkgs {
			ast.Walk(pvv, f)
		}
	}
}

func excludeKeys(dir []string, keys string) []string {
	keys = strings.TrimSpace(keys)
	// 如果没有过滤设置，返回原路径
	if len(keys) == 0 {
		return dir
	}

	// 分割
	key := strings.Split(keys, ",")
	for i := 0; i < len(key); i++ {
		key[i] = strings.TrimSpace(key[i])
	}

	var rs []string

Outer:
	for _, path := range dir {
		for _, pattern := range key {
			if strings.Contains(path, pattern) {
				continue Outer
			}
		}
		rs = append(rs, path)
	}
	return rs
}

// 对应指令：go list -f {{ .Dir }}
// - 如果dir为空，则：dir = "./..."
func listPackages(debug bool, dir string) []string {
	if len(dir) == 0 {
		dir = "./..."
	}

	if debug {
		log.Printf("[debug] path: %s\n", dir)
	}

	// 1. go list获取go文件目录列表
	output, err := exec.Command("go", "list", "-f", "{{ .Dir }}", dir).Output()
	if err != nil {
		log.Panic("[error] go list failed, output: ", string(output))
	}
	str := string(output)

	if debug {
		log.Printf("[debug] dir:\n%s", str)
	}

	// 2. 解析结果，同时过滤掉空目录
	tmp := strings.Split(str, "\n")
	var arrs []string
	for _, arr := range tmp {
		if len(arr) != 0 {
			arrs = append(arrs, arr)
		}
	}
	return arrs
}
