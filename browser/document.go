package browser

import (
	"github.com/ericchiang/css"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type DocumentEvent = string

type StaticNodeList []Node

const (
	DocumentEventDOMContentLoaded DocumentEvent = "DOMContentLoaded"
	DocumentEventLoad             DocumentEvent = "load"
)

type Document interface {
	Node
	Append(Element) Element
	Body() Element
	CreateElement(string) Element
	DocumentElement() Element
	GetElementById(string) Element
	QuerySelector(string) (Node, error)
	QuerySelectorAll(string) (StaticNodeList, error)
	// unexported
	createElement(*html.Node) Element
}
type elementConstructor func(doc *document) Element

type document struct {
	node
}

func NewDocument() Document {
	return &document{node: newNode()}
}

func (d *document) Body() Element {
	root := d.DocumentElement()
	if root != nil {
		for _, child := range root.ChildNodes() {
			if e, ok := child.(Element); ok {
				if e.TagName() == "BODY" {
					return e
				}
			}
		}
	}
	return nil
}

func (d *document) CreateElement(name string) Element {
	node := &html.Node{
		Type:      html.ElementNode,
		DataAtom:  atom.Lookup([]byte(name)),
		Data:      name,
		Namespace: "",
	}
	return d.createElement(node)
}

func (d *document) createElement(node *html.Node) Element {
	return NewHTMLElement(node)
}

func (d *document) Append(element Element) Element {
	NodeHelper{d}.AppendChild(element)
	return element
}

func (d *document) DocumentElement() Element {
	for _, c := range d.ChildNodes() {
		if e, ok := c.(Element); ok {
			return e
		}
	}
	return nil
}

func (d *document) NodeName() string { return "#document" }

func (d *document) Connected() bool {
	return true
}

func (d *document) GetElementById(id string) Element {
	var search func(node Node) Element
	search = func(node Node) Element {
		if elm, ok := node.(Element); ok && elm.GetAttribute("id") == id {
			return elm
		}
		for _, child := range node.ChildNodes() {
			if found := search(child); found != nil {
				return found
			}
		}
		return nil
	}
	return search(d)
}

func (d *document) QuerySelector(pattern string) (Node, error) {
	nodes, err := d.QuerySelectorAll(pattern)
	var result Node
	if len(nodes) > 0 {
		result = nodes[0]
	}
	return result, err
}

func (d *document) QuerySelectorAll(pattern string) (StaticNodeList, error) {
	sel, err := css.Parse(pattern)
	if err != nil {
		return nil, err
	}
	htmlNode, m := toHtmlNodeAndMap(d)
	nodes := sel.Select(htmlNode)
	result := make(StaticNodeList, len(nodes))
	for i, node := range nodes {
		resultNode := m[node]
		result[i] = resultNode
	}
	return result, nil
}

func (d *document) createHtmlNode() *html.Node {
	return &html.Node{
		Type: html.DocumentNode,
	}
}
