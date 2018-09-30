package lifeline

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	startSection = "Start"
)

type Story struct {
	raw      string
	rootNode Node
	sections map[string]Node
}

func newStory(data string) *Story {
	return &Story{
		raw:      data,
		rootNode: newBaseNode(""),
		sections: make(map[string]Node),
	}
}

func (st *Story) Init() error {
	_, err := st.rootNode.Parse(st.raw)
	if err != nil {
		return err
	}

	for i := 0; i < st.rootNode.ChildCount(); i++ {
		node := st.rootNode.Child(i)
		if _, ok := node.(*SectionNode); ok {
			st.sections[node.Content()] = node
		}
	}

	sectionCount := strings.Count(st.raw, sectionPrefix)
	if len(st.sections) != sectionCount {
		return fmt.Errorf("parsed section count mismatch, got %d, total: %d", len(st.sections), sectionCount)
	}

	return nil
}

func (st *Story) Play(ctx *Context) error {
	section, ok := st.sections[ctx.currentSection]
	if !ok {
		return fmt.Errorf("invalid section %s", ctx.currentSection)
	}

	start := ctx.currentSection
	section.Play(ctx)
	if ctx.currentSection != start {
		return st.Play(ctx)
	}
	return nil
}

func (st *Story) Reply(ctx *Context) error {
	section, ok := st.sections[ctx.currentSection]
	if !ok {
		return fmt.Errorf("invalid section %s", ctx.currentSection)
	}

	answerID, err := strconv.Atoi(ctx.data)
	if err != nil {
		return nil
	}

	lastNode := section.LastChild()
	if lastNode == nilNode {
		return fmt.Errorf("section %s has no child", section.Content())
	}

	node := lastNode.Eval(ctx)
	if node == nilNode {
		return fmt.Errorf("invalid answer id=%d in section %s", answerID, section.Content())
	}

	ctx.data = ""

	start := ctx.currentSection
	node.Play(ctx)
	if ctx.currentSection != start {
		return st.Play(ctx)
	}
	return nil
}
