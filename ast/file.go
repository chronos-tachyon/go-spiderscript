package ast

import (
	"strings"

	"github.com/chronos-tachyon/go-spiderscript/internal/util"
	"github.com/chronos-tachyon/go-spiderscript/token"
	"github.com/chronos-tachyon/go-spiderscript/tokenpredicate"
)

type File struct {
	Statements []Node
}

func (file *File) Init() {
	*file = File{
		Statements: make([]Node, 0, 32),
	}
}

func (file *File) Match(p *Parser) bool {
	for {
		if p.Consume(nil, tokenpredicate.Type(token.EOF)) {
			return true
		}

		var stmt Node
		if MatchStatement(&stmt, p) {
			file.Statements = append(file.Statements, stmt)
			continue
		}

		return true
	}
}

func (file *File) Children() []Node {
	list := make([]Node, len(file.Statements))
	for index, stmt := range file.Statements {
		list[index] = stmt
	}
	return list
}

func (file *File) String() string {
	return util.StringImpl(file)
}

func (file *File) GoString() string {
	return util.GoStringImpl(file)
}

func (file *File) EstimateStringLength() uint {
	var sum uint
	for _, stmt := range file.Statements {
		sum += stmt.EstimateStringLength()
	}
	return sum
}

func (file *File) EstimateGoStringLength() uint {
	var sum uint = 6
	for _, stmt := range file.Statements {
		sum += 1 + stmt.EstimateGoStringLength()
	}
	return sum
}

func (file *File) WriteStringTo(out *strings.Builder) {
	for _, stmt := range file.Statements {
		stmt.WriteStringTo(out)
	}
}

func (file *File) WriteGoStringTo(out *strings.Builder) {
	out.WriteString("File{")
	for index, stmt := range file.Statements {
		if index > 0 {
			out.WriteByte(',')
		}
		stmt.WriteGoStringTo(out)
	}
	out.WriteByte('}')
}

func (file *File) ComputeStringLengthEstimates() {
	for _, stmt := range file.Statements {
		stmt.ComputeStringLengthEstimates()
	}
}

var _ Node = (*File)(nil)
