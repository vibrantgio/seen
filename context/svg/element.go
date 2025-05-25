package svg

import (
	"io"
	"strings"
)

// Element
type Element struct {
	*Dom

	Namespace string
	Tag       string

	attributes  map[string]string
	ChildNodes  []*Element
	TextContent string
}

func (e *Element) AppendChild(child *Element) {
	e.ChildNodes = append(e.ChildNodes, child)
}

func (e *Element) ReplaceChild(newChild, oldChild *Element) {
	for i, child := range e.ChildNodes {
		if child == oldChild {
			e.ChildNodes[i] = newChild
			return
		}
	}
}

func (e *Element) Attribute(name string) (string, bool) {
	v, present := e.attributes[name]
	return v, present
}

func (e *Element) SetAttribute(name, value string) {
	if e.attributes == nil {
		e.attributes = make(map[string]string)
	}
	e.attributes[name] = value
	if strings.ToLower(name) == "id" {
		e.Dom.SetElementById(value, e)
	}
}

func (e *Element) WriteTo(w io.Writer) (n int64, err error) {
	sn, err := io.WriteString(w, "<"+e.Tag)
	if err != nil {
		return
	}
	n += int64(sn)
	for key, val := range e.attributes {
		sn, err = io.WriteString(w, " "+key+"=\""+val+"\"")
		if err != nil {
			return
		}
		n += int64(sn)
	}
	if e.TextContent != "" {
		sn, err = io.WriteString(w, ">"+e.TextContent)
		if err != nil {
			return
		}
	} else {
		sn, err = io.WriteString(w, ">")
		if err != nil {
			return
		}
	}
	n += int64(sn)
	for _, child := range e.ChildNodes {
		var en int64
		en, err = child.WriteTo(w)
		if err != nil {
			return
		}
		n += en
	}
	sn, err = io.WriteString(w, "</"+e.Tag+">\n")
	if err != nil {
		return
	}
	n += int64(sn)
	return
}
