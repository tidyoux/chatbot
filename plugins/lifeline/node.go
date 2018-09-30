package lifeline

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	commentPrefix = "//"

	// section
	sectionPrefix = "::"

	// jump
	jumpPrefix      = "[["
	jumpDelayPrefix = "delay"
	jumpEndPrefix   = "]]"

	// cmd
	cmdPrefix    = "<<"
	setPrefix    = "<<set"
	choisePrefix = "<<choice"
	ifPrefix     = "<<if"
	elsePrefix   = "<<else"
	elseifPrefix = "<<elseif"
	endifPrefix  = "<<endif"
	cmdEndPrefix = ">>"
)

// Context def.
type Context struct {
	channel        string
	currentSection string
	msgch          chan<- string
	data           string
}

// Node interface.
type Node interface {
	Parse(string) (int, error)

	Content() string
	AddChild(child Node)
	Child(idx int) Node
	LastChild() Node
	ChildCount() int
	Eval(ctx *Context) Node

	Play(ctx *Context)
}

var (
	nilNode Node = newBaseNode("")
)

// BaseNode def.
type BaseNode struct {
	typ      string
	content  string
	children []Node
}

func newBaseNode(typ string) *BaseNode {
	return &BaseNode{
		typ: typ,
	}
}

func lineRemain(raw string) string {
	for i := range raw {
		if raw[i] == '\n' || strings.HasPrefix(raw[i:], jumpPrefix) || strings.HasPrefix(raw[i:], cmdPrefix) {
			return raw[:i]
		}
	}
	return raw
}

func (n *BaseNode) Parse(raw string) (int, error) {
	for i := 0; i < len(raw); {
		if raw[i] == ' ' || raw[i] == '\n' {
			i++
			continue
		}

		var node Node
		switch {
		case strings.HasPrefix(raw[i:], commentPrefix):
			node = newCommentNode()
		case strings.HasPrefix(raw[i:], sectionPrefix):
			if n.typ == "SectionNode" {
				return i, nil
			}
			node = newSectionNode()
		case strings.HasPrefix(raw[i:], jumpPrefix):
			node = newJumpNode()
		case strings.HasPrefix(raw[i:], cmdPrefix):
			switch {
			case strings.HasPrefix(raw[i:], setPrefix):
				node = newSetNode()
			case strings.HasPrefix(raw[i:], choisePrefix):
				node = newChoiseNode()
			case strings.HasPrefix(raw[i:], ifPrefix):
				node = newIfNode()
			case strings.HasPrefix(raw[i:], elseifPrefix), strings.HasPrefix(raw[i:], elsePrefix), strings.HasPrefix(raw[i:], endifPrefix):
				return i, nil
			default:
				node = newUnHandleCmdNode()
			}
		default:
			node = newTextNode()
		}

		k, err := node.Parse(raw[i:])
		if err != nil {
			return 0, err
		}

		n.AddChild(node)
		i += k
	}
	return len(raw), nil
}

func (n *BaseNode) Content() string {
	return n.content
}

func (n *BaseNode) AddChild(child Node) {
	n.children = append(n.children, child)
}

func (n *BaseNode) Child(idx int) Node {
	if idx < 0 || len(n.children) <= idx {
		return nilNode
	}
	return n.children[idx]
}

func (n *BaseNode) LastChild() Node {
	return n.Child(n.ChildCount() - 1)
}

func (n *BaseNode) ChildCount() int {
	return len(n.children)
}

func (n *BaseNode) Eval(ctx *Context) Node {
	child := n.LastChild()
	if child != nilNode {
		return child.Eval(ctx)
	}
	return n
}

func (n *BaseNode) Play(ctx *Context) {
	for _, child := range n.children {
		child.Play(ctx)
	}
}

// CommentNode def.
type CommentNode struct {
	*BaseNode
}

func newCommentNode() *CommentNode {
	return &CommentNode{
		newBaseNode("CommentNode"),
	}
}

func (n *CommentNode) Parse(raw string) (int, error) {
	i := len(commentPrefix)
	s := lineRemain(raw[i:])
	n.content = strings.TrimSpace(s)
	return i + len(s), nil
}

func (n *CommentNode) Eval(ctx *Context) Node {
	return n
}

// SectionNode def.
type SectionNode struct {
	*BaseNode
}

func newSectionNode() *SectionNode {
	return &SectionNode{
		newBaseNode("SectionNode"),
	}
}

func (n *SectionNode) Parse(raw string) (int, error) {
	i := len(sectionPrefix)
	section := lineRemain(raw[i:])
	n.content = strings.TrimSpace(section)

	i += len(section)
	k, err := n.BaseNode.Parse(raw[i:])
	if err != nil {
		return 0, err
	}

	return i + k, nil
}

func (n *SectionNode) Eval(ctx *Context) Node {
	return n
}

// TextNode def.
type TextNode struct {
	*BaseNode
}

func newTextNode() *TextNode {
	return &TextNode{
		newBaseNode("TextNode"),
	}
}

func (n *TextNode) Parse(raw string) (int, error) {
	s := lineRemain(raw)
	n.content = s
	return len(s), nil
}

func (n *TextNode) Eval(ctx *Context) Node {
	return n
}

func (n *TextNode) Play(ctx *Context) {
	ctx.msgch <- n.Content()
	time.Sleep(time.Second * 3)
}

// JumpNode def.
type JumpNode struct {
	*BaseNode
	target string
	delay  time.Duration
}

func newJumpNode() *JumpNode {
	return &JumpNode{
		BaseNode: newBaseNode("JumpNode"),
	}
}

func (n *JumpNode) Parse(raw string) (int, error) {
	k := strings.Index(raw, jumpEndPrefix)
	if k < 0 {
		return 0, fmt.Errorf("can't find jump end tag")
	}

	i := len(jumpPrefix)
	exps := strings.Split(raw[i:k], "|")
	if len(exps) == 0 {
		return 0, fmt.Errorf("invalid jump format: %s", raw[:k])
	}

	if len(exps) == 1 {
		n.target = exps[0]
	} else {
		n.target = exps[1]

		if strings.HasPrefix(exps[0], jumpDelayPrefix) {
			n.content = ""

			exps = strings.Split(exps[0], " ")
			if len(exps) != 2 {
				return 0, fmt.Errorf("invalid jump delay format: %s", strings.Join(exps, " "))
			}

			delay := exps[1]
			delayValue, err := strconv.Atoi(delay[:len(delay)-1])
			if err != nil {
				return 0, fmt.Errorf("invalid jump delay format: %s", strings.Join(exps, " "))
			}

			switch delay[len(delay)-1] {
			case 's':
				n.delay = time.Second * time.Duration(delayValue)
			case 'm':
				n.delay = time.Minute * time.Duration(delayValue)
			case 'h':
				n.delay = time.Hour * time.Duration(delayValue)
			default:
				return 0, fmt.Errorf("invalid jump delay format: %s", strings.Join(exps, " "))
			}
		} else {
			n.delay = 0
			n.content = exps[0]
		}
	}
	return k + len(jumpEndPrefix), nil
}

func (n *JumpNode) Eval(ctx *Context) Node {
	return n
}

func (n *JumpNode) Play(ctx *Context) {
	if n.delay > 0 {
		time.Sleep(time.Second * 10)
	}
	ctx.currentSection = n.target
}

// SetNode def.
type SetNode struct {
	*BaseNode
	key   string
	value []string
}

func newSetNode() *SetNode {
	return &SetNode{
		BaseNode: newBaseNode("SetNode"),
	}
}

func (n *SetNode) Parse(raw string) (int, error) {
	k := strings.Index(raw, cmdEndPrefix)
	if k < 0 {
		return 0, fmt.Errorf("can't find cmd-set end tag")
	}

	exp := strings.TrimSpace(raw[len(setPrefix):k])
	exps := strings.Split(exp, "=")
	if len(exps) != 2 {
		return 0, fmt.Errorf("invalid cmd-set format: %s", exp)
	}

	key := strings.TrimSpace(exps[0])
	if !strings.HasPrefix(key, "$") {
		return 0, fmt.Errorf("invalid cmd-set format: %s", exp)
	}

	key = key[1:]

	value := strings.TrimSpace(exps[1])
	valueExps := strings.Split(value, " ")
	if len(valueExps)%2 != 1 {
		return 0, fmt.Errorf("invalid cmd-set format: %s", exp)
	}

	n.key = key
	n.value = valueExps
	return k + len(cmdEndPrefix), nil
}

func (n *SetNode) Eval(ctx *Context) Node {
	return n
}

func calculate(channel string, exps []string) (string, error) {
	var (
		value string
		op    string
	)

	for _, exp := range exps {
		switch exp {
		case "+", "-", "*", "/":
			op = exp
		default:
			v := exp
			if strings.HasPrefix(v, "$") {
				v, _ = getStatus(channel, v[1:])
			}

			if len(op) > 0 {
				v1, err := strconv.Atoi(value)
				if err != nil {
					return "", fmt.Errorf("invalid caculate expression: %s", strings.Join(exps, " "))
				}

				v2, err := strconv.Atoi(v)
				if err != nil {
					return "", fmt.Errorf("invalid caculate expression: %s", strings.Join(exps, " "))
				}

				switch op {
				case "+":
					value = strconv.Itoa(v1 + v2)
				case "-":
					value = strconv.Itoa(v1 - v2)
				case "*":
					value = strconv.Itoa(v1 * v2)
				case "/":
					value = strconv.Itoa(v1 / v2)
				}
				op = ""
			} else {
				value = v
			}
		}
	}
	return value, nil
}

func (n *SetNode) Play(ctx *Context) {
	value, err := calculate(ctx.channel, n.value)
	if err != nil {
		log.Printf("Error: play cmd-set failed, %v\n", err)
		return
	}

	setStatus(ctx.channel, n.key, value)
}

// ChoiseNode def.
type ChoiseNode struct {
	*BaseNode
}

func newChoiseNode() *ChoiseNode {
	return &ChoiseNode{
		newBaseNode("ChoiseNode"),
	}
}

func (n *ChoiseNode) Parse(raw string) (int, error) {
	for i := 0; i < len(raw); {
		switch {
		case raw[i] == ' ' || raw[i] == '|':
			i++
			continue
		case strings.HasPrefix(raw[i:], choisePrefix):
			i += len(choisePrefix)
			continue
		case strings.HasPrefix(raw[i:], cmdEndPrefix):
			i += len(cmdEndPrefix)
			continue
		case strings.HasPrefix(raw[i:], jumpPrefix):
			node := newJumpNode()
			k, err := node.Parse(raw[i:])
			if err != nil {
				return 0, err
			}

			n.AddChild(node)
			i += k
		default:
			return i, nil
		}
	}
	return len(raw), nil
}

func (n *ChoiseNode) Eval(ctx *Context) Node {
	if len(ctx.data) > 0 {
		answerID, err := strconv.Atoi(ctx.data)
		if err == nil {
			node := n.Child(answerID - 1)
			if node != nilNode {
				return node
			}
		}
	}

	return n
}

func (n *ChoiseNode) Play(ctx *Context) {
	ctx.msgch <- "--------------"
	for i, child := range n.children {
		ctx.msgch <- fmt.Sprintf("%d. %s", i+1, child.Content())
	}
}

// ConditionNode def.
type ConditionNode struct {
	*BaseNode
	left  string
	op    string
	right []string
}

func newConditionNode() *ConditionNode {
	return &ConditionNode{
		BaseNode: newBaseNode("ConditionNode"),
	}
}

func (n *ConditionNode) Parse(raw string) (int, error) {
	var i int
	switch {
	case strings.HasPrefix(raw, ifPrefix):
		i = len(ifPrefix)
	case strings.HasPrefix(raw, elseifPrefix):
		i = len(elseifPrefix)
	case strings.HasPrefix(raw, elsePrefix):
		i = len(elsePrefix)
	default:
		return 0, fmt.Errorf("can't find condition start tag")
	}

	k, err := n.parseCondition(raw[i:])
	if err != nil {
		return 0, err
	}

	i += k
	node := newBaseNode("")
	k, err = node.Parse(raw[i:])
	if err != nil {
		return 0, err
	}

	n.AddChild(node)
	return i + k, nil
}

func (n *ConditionNode) parseCondition(raw string) (int, error) {
	k := strings.Index(raw, cmdEndPrefix)
	if k < 0 {
		return 0, fmt.Errorf("can't find condition end tag")
	}

	exp := strings.TrimSpace(raw[:k])
	if len(exp) > 0 {
		exps := strings.Split(exp, " ")
		if len(exps) < 3 || !(exps[1] == "is" || exps[1] == "eq" || exps[1] == "gte") {
			return 0, fmt.Errorf("invalid condition format: %s", exp)
		}

		n.left = exps[0]
		n.op = exps[1]
		n.right = exps[2:]
	}
	return k + len(cmdEndPrefix), nil
}

func (n *ConditionNode) Eval(ctx *Context) Node {
	if len(n.left) == 0 {
		return n.Child(0).Eval(ctx)
	}

	l, err := calculate(ctx.channel, []string{n.left})
	if err != nil {
		log.Printf("Error: invalid condition left format: %s\n", n.left)
		return nilNode
	}

	r, err := calculate(ctx.channel, n.right)
	if err != nil {
		log.Printf("Error: invalid condition right format: %s\n", strings.Join(n.right, " "))
		return nilNode
	}

	var is bool
	switch n.op {
	case "is", "eq":
		is = (l == r)
	case "gte":
		lInt, err := strconv.Atoi(l)
		if err != nil {
			log.Printf("Error: invalid condition left format for op=%s: %s\n", n.op, n.left)
			return nilNode
		}

		rInt, err := strconv.Atoi(r)
		if err != nil {
			log.Printf("Error: invalid condition right format for op=%s: %s\n", n.op, strings.Join(n.right, " "))
			return nilNode
		}
		is = (lInt >= rInt)
	}

	if is {
		return n.Child(0).Eval(ctx)
	}

	return nilNode
}

func (n *ConditionNode) Play(ctx *Context) {
	n.Eval(ctx).Play(ctx)
}

// IfNode def.
type IfNode struct {
	*BaseNode
}

func newIfNode() *IfNode {
	return &IfNode{
		newBaseNode("IfNode"),
	}
}

func (n *IfNode) Parse(raw string) (int, error) {
	for i := 0; i < len(raw); {
		if strings.HasPrefix(raw[i:], ifPrefix) || strings.HasPrefix(raw[i:], elseifPrefix) || strings.HasPrefix(raw[i:], elsePrefix) {
			node := newConditionNode()
			pos, err := node.Parse(raw[i:])
			if err != nil {
				return 0, err
			}

			n.AddChild(node)
			i += pos
		} else if strings.HasPrefix(raw[i:], endifPrefix) {
			return i + len(endifPrefix) + len(cmdEndPrefix), nil
		} else {
			i++
		}
	}

	return 0, fmt.Errorf("can't find cmd-endif tag")
}

func (n *IfNode) Eval(ctx *Context) Node {
	for _, child := range n.children {
		if node := child.Eval(ctx); node != nilNode {
			return node
		}
	}
	return nilNode
}

func (n *IfNode) Play(ctx *Context) {
	n.Eval(ctx).Play(ctx)
}

// UnHandleCmdNode def.
type UnHandleCmdNode struct {
	*BaseNode
}

func newUnHandleCmdNode() *UnHandleCmdNode {
	return &UnHandleCmdNode{
		newBaseNode("UnHandleCmdNode"),
	}
}

func (n *UnHandleCmdNode) Parse(raw string) (int, error) {
	k := strings.Index(raw, cmdEndPrefix)
	if k < 0 {
		return 0, fmt.Errorf("can't find cmdend tag")
	}

	n.content = raw[len(cmdPrefix):k]
	return k + len(cmdEndPrefix), nil
}

func (n *UnHandleCmdNode) Eval(ctx *Context) Node {
	return n
}

func (n *UnHandleCmdNode) Play(ctx *Context) {
	log.Printf("Warning: unhandle cmd: %s\n", n.Content())
}
