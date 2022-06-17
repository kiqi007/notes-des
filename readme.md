# go，文件注释解析工具
- 此应用程序将扫描目录下的所有文件，查找到特定注释("@Des: {describe text}" as default)，并构建成易读的结构并输出

## 使用
```
 go install github.com/kiqi007/notes-des/des
  
 des -out=trie -out_path=./des.tree -ekey="/pkg/,/pb/" ./...
```

## 结果展示
```
./notes-des
├── common.go - [ast可复用的共有抽象部分]
├── des
│   ├── des_output.go - [将解析结果解析成特定格式并输出 (StdOut/TreeOut/TrieOut)]
│   ├── des_visitor.go - [visitor实现，ast获取文件注释并匹配合规的注释]
│   └── main.go - [des 命令行入口]
└── trie/treeprint.go - [树状结构构建与输出，为原实现扩展了压缩单叉节点的功能 / copy from github.com/xlab/treeprint]
```

## flag介绍
```
// des -help

usage: des [options] dir
 - 此应用程序将扫描目录下的所有文件，查找到特定注释("@Des: {describe text}" as default)，并构建成易读的结构并输出
   标准输出格式:                           	' filename.go - {describe text}'
   无需备注的文件 	( @Des: -):        			' filename.go'
   已废弃的文件  	( @Des: [d|delete|废弃]):   	' filename.go(del)'

Flags:
  -debug
        是否输出debug日志
  -ekey string
        如果文件路径存在关键词，则排除此目录，多个关键词之间用','分割
  -empty string
        如果文件无匹配注释，则结果填入此字符串
  -out string
        输出模式，可选值：normal, tree, trie (default "normal")
  -out_path string
        输出目录，若未设置则输出到控制台
  -sign string
        匹配标志位，程序将检索所有存在此字符串的注释作为结果 (default "@Des: ")
  -t    解析测试文件
```