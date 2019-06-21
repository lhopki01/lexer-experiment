package ast

type JenkinsFile struct {
	Library  string
	Imports  []string
	Function string
	Values   map[string]interface{}
}

type ConcatenatedItem struct {
	Primary interface{}
	Append  interface{}
}

// base json Node
type Json interface {
	TokenLiteral()
}

type Array struct {
	Array []Json
}

type Object struct {
	Object map[string]Json
}

type String struct {
	Value string
}

type Integer struct {
	Value string
}

// construct interface
func (node String) TokenLiteral()  {}
func (node Integer) TokenLiteral() {}
func (node Array) TokenLiteral()   {}
func (node Object) TokenLiteral()  {}
