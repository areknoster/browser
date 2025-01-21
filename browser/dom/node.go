package dom

import (
	"errors"
	"slices"
	"strings"

	. "github.com/stroiman/go-dom/browser/internal/dom"
	"github.com/stroiman/go-dom/browser/internal/entity"
	"github.com/stroiman/go-dom/browser/internal/log"
	"golang.org/x/net/html"
)

type NodeType int

const (
	NodeTypeElement               NodeType = 1
	NodeTypeAttribute             NodeType = 2
	NodeTypeText                  NodeType = 3
	NodeTypeCDataSection          NodeType = 4
	NodeTypeProcessingInstruction NodeType = 7
	NodeTypeComment               NodeType = 8
	NodeTypeDocument              NodeType = 9
	NodeTypeDocumentType          NodeType = 10
	NodeTypeDocumentFragment      NodeType = 11
)

type GetRootNodeOptions bool

type Node interface {
	entity.Entity
	EventTarget
	AppendChild(node Node) (Node, error)
	GetRootNode(options ...GetRootNodeOptions) Node
	ChildNodes() NodeList
	IsConnected() bool
	Contains(node Node) bool
	InsertBefore(newNode Node, referenceNode Node) (Node, error)
	NodeName() string
	NodeType() NodeType
	OwnerDocument() Document
	Parent() Node
	RemoveChild(node Node) (Node, error)
	NextSibling() Node
	PreviousSibling() Node
	FirstChild() Node
	TextContent() string
	SetTextContent(value string)
	// SetSelf must be called when creating instances of structs embedding a Node.
	//
	// If this is not called, the specialised type, which is itself a Node, will
	// not be returned from functions that should have returned it, e.g., through
	// ChildNodes. Only the embedded Node will be returned, and any specialised
	// behaviour, including HTML output, will not work.
	//
	// This function is a workaround to solve a fundamental problem. The DOM
	// specifies a model that is fundamentally object-oriented, with sub-classes
	// overriding behaviour in super-classes. This is not a behaviour that Go has.
	SetSelf(node Node)

	getSelf() Node
	createHtmlNode() *html.Node
	setParent(node Node)
	nodes() []Node
}

type node struct {
	eventTarget
	entity.Entity
	self       Node
	childNodes NodeList
	parent     Node
}

func newNode() node {
	return node{newEventTarget(), entity.New(), nil, newNodeList(), nil}
}

func (n *node) AppendChild(child Node) (Node, error) {
	_, err := n.self.InsertBefore(child, nil)
	return child, err
}

func (n *node) InsertBefore(newChild Node, referenceNode Node) (Node, error) {
	if fragment, ok := newChild.(DocumentFragment); ok {
		for fragment.ChildNodes().Length() > 0 {
			if _, err := n.InsertBefore(fragment.ChildNodes().Item(0), referenceNode); err != nil {
				return nil, err
			}
		}
		return fragment, nil
	}
	result, err := n.insertBefore(newChild, referenceNode)
	if err == nil {
		newChild.setParent(n.self)
	}
	return result, err
}

func (n *node) ChildNodes() NodeList { return n.childNodes }

func (n *node) GetRootNode(options ...GetRootNodeOptions) Node {
	if len(options) > 1 {
		log.Warn("Node.GetRootNode: composed not yet implemented")
	}
	if n.parent == nil {
		return n.self
	} else {
		return n.parent.GetRootNode(options...)
	}
}

func (n *node) Contains(node Node) bool {
	for _, c := range n.ChildNodes().All() {
		if c == node || c.Contains(node) {
			return true

		}
	}
	return false
}

func (n *node) Parent() Node { return n.parent }

func (n *node) setParent(parent Node) {
	n.parent = parent
	n.parentTarget = parent
}

func (n *node) IsConnected() (result bool) {
	if n.parent != nil {
		result = n.parent.IsConnected()
	}
	return
}

func (n *node) NodeName() string {
	return "#node"
}

// removeNodeFromParent removes the node from the current parent, _if_ it has
// one. Does nothing for disconnected nodes.
func removeNodeFromParent(node Node) {
	parent := node.Parent()
	if parent != nil {
		parent.RemoveChild(node)
	}
}

func (n *node) RemoveChild(node Node) (Node, error) {
	idx := slices.Index(n.childNodes.All(), node)
	if idx == -1 {
		return nil, newDomError(
			"Node.removeChild: The node to be removed is not a child of this node",
		)
	}
	n.childNodes.setNodes(slices.Delete(n.childNodes.All(), idx, idx+1))
	return node, nil
}

func (n *node) insertBefore(newNode Node, referenceNode Node) (Node, error) {
	if _, isAttribute := newNode.(Attr); isAttribute {
		return nil, newDomError("Node.appendChild: May not add an Attribute as a child")
	}
	// TODO, Don't allow newNode to be inserted in it's own branch (circular tree)
	// TODO, Handle a fragment. Also returns nil
	if referenceNode == nil {
		n.childNodes.append(newNode)
	} else {
		i := slices.Index(n.childNodes.All(), referenceNode)
		if i == -1 {
			return nil, errors.New("Reference node not found")
		}
		n.childNodes.setNodes(slices.Insert(n.childNodes.All(), i, newNode))
	}
	removeNodeFromParent(newNode)
	return newNode, nil
}

type nodeIterator struct{ Node }

func toHtmlNode(node Node) *html.Node {
	return nodeIterator{node}.toHtmlNode(nil)
}
func toHtmlNodeAndMap(node Node) (*html.Node, map[*html.Node]Node) {
	m := make(map[*html.Node]Node)
	result := nodeIterator{node}.toHtmlNode(m)
	return result, m
}

func (n nodeIterator) toHtmlNode(m map[*html.Node]Node) *html.Node {
	htmlNode := n.Node.createHtmlNode()
	if m != nil {
		m[htmlNode] = n.Node
	}
	for _, child := range n.nodes() {
		htmlNode.AppendChild(nodeIterator{child}.toHtmlNode(m))
	}
	return htmlNode
}

func (n *node) OwnerDocument() Document {
	parent := n.Parent()
	if parent != nil {
		return parent.OwnerDocument()
	}
	return nil
}

func (n *node) FirstChild() Node {
	if n.childNodes.Length() == 0 {
		return nil
	}
	return n.childNodes.Item(0)
}

func (n *node) NextSibling() Node {
	children := n.Parent().nodes()
	idx := slices.IndexFunc(
		children,
		func(child Node) bool { return n.ObjectId() == child.ObjectId() },
	) + 1
	if idx == 0 {
		panic("We should exist in our parent's collection")
	}
	if idx >= len(children) {
		return nil
	}
	return children[idx]
}
func (n *node) PreviousSibling() Node {
	children := n.Parent().nodes()
	idx := slices.IndexFunc(
		children,
		func(child Node) bool { return n.ObjectId() == child.ObjectId() },
	) - 1
	if idx == -2 {
		panic("We should exist in our parent's collection")
	}
	if idx < 0 {
		return nil
	}
	return children[idx]
}

func (n *node) nodes() []Node {
	return n.childNodes.All()
}

func (n *node) SetSelf(node Node) { n.self = node }
func (n *node) getSelf() Node     { return n.self }

func (n *node) SetTextContent(val string) {
	for x := n.FirstChild(); x != nil; x = n.FirstChild() {
		x.RemoveChild(x)
	}
	n.AppendChild(NewText(val))
}

func (n *node) TextContent() string {
	b := &strings.Builder{}
	for _, node := range n.nodes() {
		b.WriteString(node.TextContent())
	}
	return b.String()
}

func (n *node) renderChildren(builder *strings.Builder) {
	if childRenderer, ok := n.self.(ChildrenRenderer); ok {
		childRenderer.RenderChildren(builder)
	}
}

func (n *node) RenderChildren(builder *strings.Builder) {
	for _, child := range n.childNodes.All() {
		if renderer, ok := child.(Renderer); ok {
			renderer.Render(builder)
		}
	}
}

func (n *node) String() string { return n.NodeName() }
