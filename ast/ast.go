package ast

import (
	"bytes"
	"curryLang/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

type AssignmentStatement struct {
	Name  *Identifier
	Value Expression
}

func (ls *AssignmentStatement) statementNode()       {}
func (ls *AssignmentStatement) TokenLiteral() string { return "<assignment>" }

func (ls *AssignmentStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

type WhileStatement struct {
	Token     token.Token // the token.WHILE token
	Condition Expression
	Body      []Statement
}

func (ls *WhileStatement) statementNode()       {}
func (ls *WhileStatement) TokenLiteral() string { return ls.Token.Literal }

func (ls *WhileStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " (")
	out.WriteString(ls.Condition.String())
	out.WriteString(")")
	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String() + ";"
	}

	return ""
}

type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type InfixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Left     Expression
	Operator string
	Right    Expression
}

func (pe *InfixExpression) expressionNode()      {}
func (pe *InfixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(" " + pe.Operator + " ")
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

func (i *Identifier) String() string { return i.Value }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type Boolean struct {
	Token token.Token
	Value bool
}

func (il *Boolean) expressionNode()      {}
func (il *Boolean) TokenLiteral() string { return il.Token.Literal }
func (il *Boolean) String() string       { return il.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (str *StringLiteral) expressionNode()      {}
func (str *StringLiteral) TokenLiteral() string { return str.Token.Literal }
func (str *StringLiteral) String() string       { return str.Value }

type ListExpression struct {
	Token token.Token
	Value []Expression
}

func (list *ListExpression) expressionNode()      {}
func (list *ListExpression) TokenLiteral() string { return list.Token.Literal }
func (list *ListExpression) String() string {
	var out bytes.Buffer

	out.WriteString("[")

	for i, s := range list.Value {
		out.WriteString(s.String())
		if i < len(list.Value)-1 {
			out.WriteString(", ")
		}
	}

	out.WriteString("]")

	return out.String()
}

type IndexAccessExpression struct {
	Token  token.Token
	Source Expression
	Value  Expression
}

func (index *IndexAccessExpression) expressionNode()      {}
func (index *IndexAccessExpression) TokenLiteral() string { return index.Token.Literal }
func (index *IndexAccessExpression) String() string {
	var out bytes.Buffer

	out.WriteString(index.Source.String())
	out.WriteString("[")
	out.WriteString(index.Value.String())
	out.WriteString("]")

	return out.String()
}

type IfElseExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence []Statement
	Alternative []Statement
}

func (il *IfElseExpression) expressionNode()      {}
func (il *IfElseExpression) TokenLiteral() string { return il.Token.Literal }
func (il *IfElseExpression) String() string       { return il.Token.Literal }
func (il *IfElseExpression) ConsequenceString() string {
	var out bytes.Buffer
	for _, s := range il.Consequence {
		out.WriteString(s.String())
	}
	return out.String()
}
func (il *IfElseExpression) AlternativeString() string {
	var out bytes.Buffer
	for _, s := range il.Alternative {
		out.WriteString(s.String())
	}
	return out.String()
}

type Parameter struct {
	Name string
}

func (p Parameter) String() string {
	return p.Name
}

type FunctionExpression struct {
	Token      token.Token
	Name       string
	Parameters []Parameter
	Body       []Statement
}

func (il *FunctionExpression) expressionNode()      {}
func (il *FunctionExpression) TokenLiteral() string { return il.Token.Literal }
func (il *FunctionExpression) String() string       { return il.Token.Literal }

func (il *FunctionExpression) ParametersString() string {
	var out bytes.Buffer
	totalParameters := len(il.Parameters) - 1
	out.WriteString("(")

	for i, s := range il.Parameters {
		out.WriteString(s.String())
		if i < totalParameters {
			out.WriteString(", ")
		}
	}

	out.WriteString(")")

	return out.String()
}

func (il *FunctionExpression) BodyString() string {
	var out bytes.Buffer
	for _, s := range il.Body {
		out.WriteString(s.String())
	}
	return out.String()
}

type FunctionCallExpression struct {
	Token        token.Token
	FunctionExpr Expression
	Parameters   []Expression
}

func (il *FunctionCallExpression) expressionNode()      {}
func (il *FunctionCallExpression) TokenLiteral() string { return il.Token.Literal }
func (il *FunctionCallExpression) String() string       { return il.Token.Literal }

func (il *FunctionCallExpression) ParametersString() string {
	var out bytes.Buffer
	totalParameters := len(il.Parameters) - 1
	out.WriteString("(")

	for i, s := range il.Parameters {
		out.WriteString(s.String())
		if i < totalParameters {
			out.WriteString(", ")
		}
	}

	out.WriteString(")")

	return out.String()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
