package main

// @Des: visitor实现，ast获取文件注释并匹配合规的注释

import (
	"flag"
	xast "github.com/kiqi007/notes-des"
	"go/ast"
	"log"
	"sort"
	"strings"
)

var needTest = flag.Bool("t", false, "解析测试文件")
var sign = flag.String("sign", "@Des: ", "匹配标志位，程序将检索所有存在此字符串的注释作为结果")
var excludeKey = flag.String("ekey", "", "如果文件路径存在关键词，则排除此目录，多个关键词之间用','分割")

type DesVisitor struct {
	*xast.Parser
	ParseRs []ParseRs // []{key, []value}

	needTest  bool
	parseSign string
}
type ParseRs struct {
	Key   string
	Value []string
}

func NewDesVisitor(parser *xast.Parser) *DesVisitor {
	vs := &DesVisitor{
		Parser:    parser,
		ParseRs:   make([]ParseRs, 0),
		needTest:  *needTest,
		parseSign: *sign,
	}
	parser.WithVisitor(vs)
	return vs
}

// Visit 因为如果直接监听file的话，存在map的不定序问题
func (d *DesVisitor) Visit(node ast.Node) (w ast.Visitor) {
	switch nd := node.(type) {
	case *ast.Package:
		var rs []*ast.File
		for _, f := range nd.Files {
			rs = append(rs, f)
		}
		sort.Slice(rs, func(i, j int) bool {
			in := d.Fset.Position(rs[i].Pos()).Filename
			jn := d.Fset.Position(rs[j].Pos()).Filename
			return in < jn
		})
		for _, r := range rs {
			d.visitFile(r)
		}
	}
	return w
}

func (d *DesVisitor) visitFile(node ast.Node) (w ast.Visitor) {
	switch nd := node.(type) {

	// 是文件级别的注释
	case *ast.File:
		fname := d.Fset.Position(nd.Pos()).Filename

		// 如果不检索测试文件，则跳过
		if !d.needTest && strings.HasSuffix(fname, "_test.go") {
			break
		}

		var hasDes bool
		rs := ParseRs{Key: fname}
		for _, r := range nd.Comments {
			for _, c := range r.List {
				if idx := strings.Index(c.Text, d.parseSign); idx != -1 {
					rs.Value = append(rs.Value, strings.TrimSpace(c.Text[idx+len(d.parseSign):]))
					hasDes = true
				}
			}
		}
		if !hasDes {
			rs.Value = []string{""}
		}

		d.ParseRs = append(d.ParseRs, rs)
		if d.Debug {
			log.Printf("[debug] %s - %+v\n", rs.Key, rs.Value)
		}
	}
	return d
}

func (d *DesVisitor) Output() []ParseRs {
	return d.ParseRs
}
