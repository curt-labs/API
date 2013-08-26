package xml_helper

import (
	"bytes"
	"errors"
	"strings"
)

type (
	Node interface {
		Children() []Node
		Set(key, value string) error
		Add(child Node) error
		String() string
		buildString(buffer *bytes.Buffer, indent int)
	}
	Attribute struct {
		name, value string
	}
	Element struct {
		tag        string
		children   []Node
		attributes map[string]string
	}
	Text struct {
		value string
	}
)

func (this *Attribute) Add(child Node) error {
	return errors.New("Adding a child to an attribute is not supported.")
}
func (this *Attribute) Children() []Node {
	return make([]Node, 0)
}
func (this *Attribute) Set(key, value string) error {
	this.name = key
	this.value = value
	return nil
}
func (this *Attribute) String() string {
	return this.name + "=" + this.value
}
func (this *Attribute) buildString(buffer *bytes.Buffer, indent int) {
	buffer.WriteString(this.String())
}

func (this *Element) Add(child Node) error {
	this.children = append(this.children, child)
	return nil
}
func (this *Element) Children() []Node {
	return this.children
}
func (this *Element) Set(key, value string) error {
	this.attributes[key] = value
	return nil
}
func (this *Element) String() string {
	var buf bytes.Buffer
	this.buildString(&buf, 0)
	return buf.String()
}
func (this *Element) buildString(buffer *bytes.Buffer, indent int) {
	for i := 0; i < indent; i++ {
		buffer.WriteByte('\t')
	}
	buffer.WriteByte('<')
	buffer.WriteString(this.tag)
	if len(this.attributes) > 0 {
		for k, v := range this.attributes {
			buffer.WriteByte(' ')
			buffer.WriteString(k)
			buffer.WriteString("=\"")
			buffer.WriteString(strings.Replace(v, "\"", "&quot;", -1))
			buffer.WriteString("\"")
		}
	}
	buffer.WriteByte('>')

	if len(this.children) > 0 {

		hasNodes := false

		for _, n := range this.children {
			if _, ok := n.(*Element); ok {
				hasNodes = true
				break
			}
		}

		if hasNodes {
			buffer.WriteByte('\n')
			for _, n := range this.children {
				n.buildString(buffer, indent+1)
			}
			for i := 0; i < indent; i++ {
				buffer.WriteByte('\t')
			}
		} else {
			for _, n := range this.children {
				n.buildString(buffer, 0)
			}
		}
	}

	buffer.WriteString("</")
	buffer.WriteString(this.tag)
	buffer.WriteString(">\n")
}

func (this *Text) Add(child Node) error {
	return errors.New("Adding a child to a text node is not supported.")
}
func (this *Text) Children() []Node {
	return make([]Node, 0)
}
func (this *Text) Set(key, value string) error {
	return errors.New("Setting an attribute on a text node is not supported.")
}
func (this *Text) String() string {
	return this.value
}
func (this *Text) buildString(buffer *bytes.Buffer, indent int) {
	buffer.WriteString(strings.Replace(this.value, "<", "&lt;", -1))
}

func A(name, value string) *Attribute {
	return &Attribute{name, value}
}
func E(tag string, children ...Node) *Element {
	attributes := make(map[string]string)
	nodes := make([]Node, 0)
	for _, n := range children {
		if att, ok := n.(*Attribute); ok {
			attributes[att.name] = att.value
		} else {
			nodes = append(nodes, n)
		}
	}
	return &Element{tag, nodes, attributes}
}
func T(text string) *Text {
	return &Text{text}
}
