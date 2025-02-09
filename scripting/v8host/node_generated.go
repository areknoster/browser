// This file is generated. Do not edit.

package v8host

import (
	"errors"
	log "github.com/gost-dom/browser/internal/log"
	v8 "github.com/tommie/v8go"
)

func init() {
	registerJSClass("Node", "EventTarget", createNodePrototype)
}

func createNodePrototype(scriptHost *V8ScriptHost) *v8.FunctionTemplate {
	iso := scriptHost.iso
	wrapper := newNodeV8Wrapper(scriptHost)
	constructor := v8.NewFunctionTemplateWithError(iso, wrapper.Constructor)

	instanceTmpl := constructor.InstanceTemplate()
	instanceTmpl.SetInternalFieldCount(1)

	prototypeTmpl := constructor.PrototypeTemplate()
	prototypeTmpl.Set("getRootNode", v8.NewFunctionTemplateWithError(iso, wrapper.getRootNode))
	prototypeTmpl.Set("cloneNode", v8.NewFunctionTemplateWithError(iso, wrapper.cloneNode))
	prototypeTmpl.Set("isSameNode", v8.NewFunctionTemplateWithError(iso, wrapper.isSameNode))
	prototypeTmpl.Set("contains", v8.NewFunctionTemplateWithError(iso, wrapper.contains))
	prototypeTmpl.Set("insertBefore", v8.NewFunctionTemplateWithError(iso, wrapper.insertBefore))
	prototypeTmpl.Set("appendChild", v8.NewFunctionTemplateWithError(iso, wrapper.appendChild))
	prototypeTmpl.Set("removeChild", v8.NewFunctionTemplateWithError(iso, wrapper.removeChild))

	prototypeTmpl.SetAccessorProperty("nodeType",
		v8.NewFunctionTemplateWithError(iso, wrapper.nodeType),
		nil,
		v8.None)
	prototypeTmpl.SetAccessorProperty("nodeName",
		v8.NewFunctionTemplateWithError(iso, wrapper.nodeName),
		nil,
		v8.None)
	prototypeTmpl.SetAccessorProperty("isConnected",
		v8.NewFunctionTemplateWithError(iso, wrapper.isConnected),
		nil,
		v8.None)
	prototypeTmpl.SetAccessorProperty("ownerDocument",
		v8.NewFunctionTemplateWithError(iso, wrapper.ownerDocument),
		nil,
		v8.None)
	prototypeTmpl.SetAccessorProperty("parentElement",
		v8.NewFunctionTemplateWithError(iso, wrapper.parentElement),
		nil,
		v8.None)
	prototypeTmpl.SetAccessorProperty("childNodes",
		v8.NewFunctionTemplateWithError(iso, wrapper.childNodes),
		nil,
		v8.None)
	prototypeTmpl.SetAccessorProperty("firstChild",
		v8.NewFunctionTemplateWithError(iso, wrapper.firstChild),
		nil,
		v8.None)
	prototypeTmpl.SetAccessorProperty("previousSibling",
		v8.NewFunctionTemplateWithError(iso, wrapper.previousSibling),
		nil,
		v8.None)
	prototypeTmpl.SetAccessorProperty("nextSibling",
		v8.NewFunctionTemplateWithError(iso, wrapper.nextSibling),
		nil,
		v8.None)
	prototypeTmpl.SetAccessorProperty("textContent",
		v8.NewFunctionTemplateWithError(iso, wrapper.textContent),
		v8.NewFunctionTemplateWithError(iso, wrapper.setTextContent),
		v8.None)

	return constructor
}

func (n nodeV8Wrapper) Constructor(info *v8.FunctionCallbackInfo) (*v8.Value, error) {
	return nil, v8.NewTypeError(n.scriptHost.iso, "Illegal Constructor")
}

func (n nodeV8Wrapper) getRootNode(info *v8.FunctionCallbackInfo) (*v8.Value, error) {
	ctx := n.mustGetContext(info)
	log.Debug("V8 Function call: Node.getRootNode")
	args := newArgumentHelper(n.scriptHost, info)
	instance, err0 := n.getInstance(info)
	options, err1 := tryParseArgWithDefault(args, 0, n.defaultGetRootNodeOptions, n.decodeGetRootNodeOptions)
	if args.noOfReadArguments >= 1 {
		err := errors.Join(err0, err1)
		if err != nil {
			return nil, err
		}
		result := instance.GetRootNode(options)
		return ctx.getInstanceForNode(result)
	}
	return nil, errors.New("Node.getRootNode: Missing arguments")
}

func (n nodeV8Wrapper) cloneNode(info *v8.FunctionCallbackInfo) (*v8.Value, error) {
	ctx := n.mustGetContext(info)
	log.Debug("V8 Function call: Node.cloneNode")
	args := newArgumentHelper(n.scriptHost, info)
	instance, err0 := n.getInstance(info)
	subtree, err1 := tryParseArgWithDefault(args, 0, n.defaultboolean, n.decodeBoolean)
	if args.noOfReadArguments >= 1 {
		err := errors.Join(err0, err1)
		if err != nil {
			return nil, err
		}
		result := instance.CloneNode(subtree)
		return ctx.getInstanceForNode(result)
	}
	return nil, errors.New("Node.cloneNode: Missing arguments")
}

func (n nodeV8Wrapper) isSameNode(info *v8.FunctionCallbackInfo) (*v8.Value, error) {
	ctx := n.mustGetContext(info)
	log.Debug("V8 Function call: Node.isSameNode")
	args := newArgumentHelper(n.scriptHost, info)
	instance, err0 := n.getInstance(info)
	otherNode, err1 := tryParseArg(args, 0, n.decodeNode)
	if args.noOfReadArguments >= 1 {
		err := errors.Join(err0, err1)
		if err != nil {
			return nil, err
		}
		result := instance.IsSameNode(otherNode)
		return n.toBoolean(ctx, result)
	}
	return nil, errors.New("Node.isSameNode: Missing arguments")
}

func (n nodeV8Wrapper) contains(info *v8.FunctionCallbackInfo) (*v8.Value, error) {
	ctx := n.mustGetContext(info)
	log.Debug("V8 Function call: Node.contains")
	args := newArgumentHelper(n.scriptHost, info)
	instance, err0 := n.getInstance(info)
	other, err1 := tryParseArg(args, 0, n.decodeNode)
	if args.noOfReadArguments >= 1 {
		err := errors.Join(err0, err1)
		if err != nil {
			return nil, err
		}
		result := instance.Contains(other)
		return n.toBoolean(ctx, result)
	}
	return nil, errors.New("Node.contains: Missing arguments")
}

func (n nodeV8Wrapper) insertBefore(info *v8.FunctionCallbackInfo) (*v8.Value, error) {
	ctx := n.mustGetContext(info)
	log.Debug("V8 Function call: Node.insertBefore")
	args := newArgumentHelper(n.scriptHost, info)
	instance, err0 := n.getInstance(info)
	node, err1 := tryParseArg(args, 0, n.decodeNode)
	child, err2 := tryParseArg(args, 1, n.decodeNode)
	if args.noOfReadArguments >= 2 {
		err := errors.Join(err0, err1, err2)
		if err != nil {
			return nil, err
		}
		result, callErr := instance.InsertBefore(node, child)
		if callErr != nil {
			return nil, callErr
		} else {
			return ctx.getInstanceForNode(result)
		}
	}
	return nil, errors.New("Node.insertBefore: Missing arguments")
}

func (n nodeV8Wrapper) appendChild(info *v8.FunctionCallbackInfo) (*v8.Value, error) {
	ctx := n.mustGetContext(info)
	log.Debug("V8 Function call: Node.appendChild")
	args := newArgumentHelper(n.scriptHost, info)
	instance, err0 := n.getInstance(info)
	node, err1 := tryParseArg(args, 0, n.decodeNode)
	if args.noOfReadArguments >= 1 {
		err := errors.Join(err0, err1)
		if err != nil {
			return nil, err
		}
		result, callErr := instance.AppendChild(node)
		if callErr != nil {
			return nil, callErr
		} else {
			return ctx.getInstanceForNode(result)
		}
	}
	return nil, errors.New("Node.appendChild: Missing arguments")
}

func (n nodeV8Wrapper) removeChild(info *v8.FunctionCallbackInfo) (*v8.Value, error) {
	ctx := n.mustGetContext(info)
	log.Debug("V8 Function call: Node.removeChild")
	args := newArgumentHelper(n.scriptHost, info)
	instance, err0 := n.getInstance(info)
	child, err1 := tryParseArg(args, 0, n.decodeNode)
	if args.noOfReadArguments >= 1 {
		err := errors.Join(err0, err1)
		if err != nil {
			return nil, err
		}
		result, callErr := instance.RemoveChild(child)
		if callErr != nil {
			return nil, callErr
		} else {
			return ctx.getInstanceForNode(result)
		}
	}
	return nil, errors.New("Node.removeChild: Missing arguments")
}

func (n nodeV8Wrapper) nodeName(info *v8.FunctionCallbackInfo) (*v8.Value, error) {
	ctx := n.mustGetContext(info)
	log.Debug("V8 Function call: Node.nodeName")
	instance, err := n.getInstance(info)
	if err != nil {
		return nil, err
	}
	result := instance.NodeName()
	return n.toDOMString(ctx, result)
}

func (n nodeV8Wrapper) isConnected(info *v8.FunctionCallbackInfo) (*v8.Value, error) {
	ctx := n.mustGetContext(info)
	log.Debug("V8 Function call: Node.isConnected")
	instance, err := n.getInstance(info)
	if err != nil {
		return nil, err
	}
	result := instance.IsConnected()
	return n.toBoolean(ctx, result)
}

func (n nodeV8Wrapper) ownerDocument(info *v8.FunctionCallbackInfo) (*v8.Value, error) {
	ctx := n.mustGetContext(info)
	log.Debug("V8 Function call: Node.ownerDocument")
	instance, err := n.getInstance(info)
	if err != nil {
		return nil, err
	}
	result := instance.OwnerDocument()
	return ctx.getInstanceForNode(result)
}

func (n nodeV8Wrapper) parentElement(info *v8.FunctionCallbackInfo) (*v8.Value, error) {
	ctx := n.mustGetContext(info)
	log.Debug("V8 Function call: Node.parentElement")
	instance, err := n.getInstance(info)
	if err != nil {
		return nil, err
	}
	result := instance.ParentElement()
	return ctx.getInstanceForNode(result)
}

func (n nodeV8Wrapper) childNodes(info *v8.FunctionCallbackInfo) (*v8.Value, error) {
	ctx := n.mustGetContext(info)
	log.Debug("V8 Function call: Node.childNodes")
	instance, err := n.getInstance(info)
	if err != nil {
		return nil, err
	}
	result := instance.ChildNodes()
	return n.toNodeList(ctx, result)
}

func (n nodeV8Wrapper) firstChild(info *v8.FunctionCallbackInfo) (*v8.Value, error) {
	ctx := n.mustGetContext(info)
	log.Debug("V8 Function call: Node.firstChild")
	instance, err := n.getInstance(info)
	if err != nil {
		return nil, err
	}
	result := instance.FirstChild()
	return ctx.getInstanceForNode(result)
}

func (n nodeV8Wrapper) previousSibling(info *v8.FunctionCallbackInfo) (*v8.Value, error) {
	ctx := n.mustGetContext(info)
	log.Debug("V8 Function call: Node.previousSibling")
	instance, err := n.getInstance(info)
	if err != nil {
		return nil, err
	}
	result := instance.PreviousSibling()
	return ctx.getInstanceForNode(result)
}

func (n nodeV8Wrapper) nextSibling(info *v8.FunctionCallbackInfo) (*v8.Value, error) {
	ctx := n.mustGetContext(info)
	log.Debug("V8 Function call: Node.nextSibling")
	instance, err := n.getInstance(info)
	if err != nil {
		return nil, err
	}
	result := instance.NextSibling()
	return ctx.getInstanceForNode(result)
}
