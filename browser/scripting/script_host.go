package scripting

import (
	"fmt"
	"log/slog"
	"runtime"
	"strings"

	. "github.com/stroiman/go-dom/browser/dom"
	"github.com/stroiman/go-dom/browser/html"

	"github.com/tommie/v8go"
	v8 "github.com/tommie/v8go"
)

type globalInstall struct {
	name        string
	constructor *v8.FunctionTemplate
}

type globals struct {
	namedGlobals map[string]*v8.FunctionTemplate
}

type ScriptHost struct {
	iso             *v8.Isolate
	inspector       *v8.Inspector
	inspectorClient *v8.InspectorClient
	windowTemplate  *v8.ObjectTemplate
	globals         globals
	contexts        map[*v8.Context]*ScriptContext
}

func (h *ScriptHost) GetContext(v8ctx *v8.Context) (*ScriptContext, bool) {
	ctx, ok := h.contexts[v8ctx]
	return ctx, ok
}

func (h *ScriptHost) MustGetContext(v8ctx *v8.Context) *ScriptContext {
	if ctx, ok := h.GetContext(v8ctx); ok {
		return ctx
	}
	panic("Unknown v8 context")
}

type ScriptContext struct {
	host      *ScriptHost
	v8ctx     *v8.Context
	window    html.Window
	pinner    runtime.Pinner
	v8nodes   map[ObjectId]*v8.Value
	domNodes  map[ObjectId]Entity
	eventLoop *EventLoop
	disposers []Disposable
}

func (c *ScriptContext) CacheNode(obj *v8.Object, node Entity) (*v8.Value, error) {
	val := obj.Value
	objectId := node.ObjectId()
	c.v8nodes[objectId] = val
	c.domNodes[objectId] = node
	internal, err := v8.NewValue(c.host.iso, objectId)
	if err != nil {
		return nil, err
	}
	obj.SetInternalField(0, internal)
	return val, nil
}

func (c *ScriptContext) GetInstanceForNode(
	node Entity,
) (*v8.Value, error) {
	iso := c.host.iso
	if node == nil {
		return v8.Null(iso), nil
	}
	switch n := node.(type) {
	case CustomEvent:
		return c.GetInstanceForNodeByName("CustomEvent", n)
	case Event:
		return c.GetInstanceForNodeByName("Event", n)
	case Element:
		if constructor, ok := htmlElements[strings.ToLower(n.TagName())]; ok {
			return c.GetInstanceForNodeByName(constructor, n)
		}
		return c.GetInstanceForNodeByName("Element", n)
	case html.HTMLDocument:
		return c.GetInstanceForNodeByName("HTMLDocument", n)
	case Document:
		return c.GetInstanceForNodeByName("Document", n)
	case DocumentFragment:
		return c.GetInstanceForNodeByName("DocumentFragment", n)
	case Node:
		return c.GetInstanceForNodeByName("Node", n)
	case Attr:
		return c.GetInstanceForNodeByName("Attr", n)
	default:
		panic("Cannot lookup node")
	}
}

func (c *ScriptContext) GetInstanceForNodeByName(
	constructor string,
	node Entity,
) (*v8.Value, error) {
	iso := c.host.iso
	if node == nil {
		return v8.Null(iso), nil
	}
	prototype, ok := c.host.globals.namedGlobals[constructor]
	if !ok {
		panic("Bad constructor name")
	}
	value, err := prototype.InstanceTemplate().NewInstance(c.v8ctx)
	if err == nil {
		objectId := node.ObjectId()
		if cached, ok := c.v8nodes[objectId]; ok {
			return cached, nil
		}
		return c.CacheNode(value, node)
	}
	return nil, err
}

func (c *ScriptContext) GetCachedNode(this *v8.Object) (Entity, bool) {
	result, ok := c.domNodes[this.GetInternalField(0).Int32()]
	return result, ok
}

type JSConstructorFactory = func(*ScriptHost) *v8.FunctionTemplate

type class struct {
	globalIdentifier string
	constructor      JSConstructorFactory
	subClasses       []class
}

// createGlobals returns an ordered list of constructors to be created in global
// scope. They must be installed in "order", as base classes must be installed
// before subclasses
func createGlobals(host *ScriptHost) []globalInstall {
	result := make([]globalInstall, 0)
	var iter func(class classSpec) *v8.FunctionTemplate
	uniqueNames := make(map[string]*v8.FunctionTemplate)
	iter = func(class classSpec) *v8.FunctionTemplate {
		if constructor, found := uniqueNames[class.name]; found {
			return constructor
		}
		var superClassConstructor *v8.FunctionTemplate
		if class.superClassName != "" {
			superClassSpec, found := classes[class.superClassName]
			if !found {
				panic(
					"Missing super class spec. Class: " + class.name + ". Super: " + class.superClassName,
				)
			}
			superClassConstructor = iter(superClassSpec)
		}
		constructor := class.factory(host)
		if superClassConstructor != nil {
			constructor.Inherit(superClassConstructor)
		}
		uniqueNames[class.name] = constructor
		result = append(result, globalInstall{class.name, constructor})
		return constructor
	}
	for _, class := range classes {
		iter(class)
	}

	// if htmlElement, ok := uniqueNames["HTMLElement"]; ok {
	// 	for _, cls := range htmlElements {
	// 		if _, ok := uniqueNames[cls]; !ok {
	// 			fn := NewIllegalConstructorBuilder[Element](host).constructor
	// 			fn.Inherit(htmlElement)
	// 			uniqueNames[cls] = fn
	// 			result = append(result, globalInstall{cls, fn})
	// 		}
	// 	}
	// }
	return result
}

func (host *ScriptHost) ConsoleAPIMessage(message v8.ConsoleAPIMessage) {
	fmt.Println("Message", message)
	switch message.ErrorLevel {
	case v8.ErrorLevelDebug:
		slog.Debug(message.Message)
	case v8.ErrorLevelInfo:
	case v8.ErrorLevelLog:
		slog.Info(message.Message)
	case v8.ErrorLevelWarning:
		slog.Warn(message.Message)
	case v8.ErrorLevelError:
		slog.Error(message.Message)
	}
}

type classSpec struct {
	name           string
	superClassName string
	factory        JSConstructorFactory
}

var classes map[string]classSpec = make(map[string]classSpec)

func RegisterJSClass(
	className string,
	superClassName string,
	constructorFactory JSConstructorFactory,
) {
	spec := classSpec{
		className, superClassName, constructorFactory,
	}
	if _, ok := classes[className]; ok {
		panic("Same class added twice: " + className)
	}
	if superClassName == "" {
		classes[className] = spec
		return
	}
	parent, parentFound := classes[superClassName]
	for parentFound {
		if parent.superClassName == className {
			panic("Recursive class parents" + className)
		}
		parent, parentFound = classes[parent.superClassName]
	}
	classes[className] = spec
}

func init() {
	RegisterJSClass("Event", "", CreateEvent)
	RegisterJSClass("CustomEvent", "Event", CreateCustomEvent)
	RegisterJSClass("NamedNodeMap", "", CreateNamedNodeMap)
	RegisterJSClass("Location", "", CreateLocationPrototype)
	RegisterJSClass("NodeList", "", CreateNodeList)
	RegisterJSClass("EventTarget", "", CreateEventTarget)

	RegisterJSClass("XMLHttpRequest", "EventTarget", CreateXmlHttpRequestPrototype)

	RegisterJSClass("Document", "Node", CreateDocumentPrototype)
	RegisterJSClass("HTMLDocument", "Document", CreateHTMLDocumentPrototype)
	RegisterJSClass("DocumentFragment", "Node", CreateDocumentFragmentPrototype)
	RegisterJSClass("ShadowRoot", "DocumentFragment", CreateShadowRootPrototype)
	RegisterJSClass("Element", "Node", CreateElement)
	RegisterJSClass("HTMLElement", "Element", CreateHtmlElement)
	RegisterJSClass("Attr", "Node", CreateAttr)

	RegisterJSClass("FormData", "", CreateFormData)
	RegisterJSClass("DOMParser", "", CreateDOMParserPrototype)

	for _, cls := range htmlElements {
		if _, found := classes[cls]; !found {
			RegisterJSClass(cls, "HTMLElement", CreateIllegalConstructor)
		}
	}
}

func NewScriptHost() *ScriptHost {
	host := &ScriptHost{iso: v8.NewIsolate()}
	host.inspectorClient = v8.NewInspectorClient(host)
	host.inspector = v8.NewInspector(host.iso, host.inspectorClient)

	globalInstalls := createGlobals(host)
	host.globals = globals{make(map[string]*v8.FunctionTemplate)}
	for _, globalInstall := range globalInstalls {
		host.globals.namedGlobals[globalInstall.name] = globalInstall.constructor
	}
	constructors := host.globals.namedGlobals
	window := constructors["Window"]
	host.windowTemplate = window.InstanceTemplate()
	host.contexts = make(map[*v8.Context]*ScriptContext)
	installGlobals(window, host, globalInstalls)
	installEventLoopGlobals(host, host.windowTemplate)
	return host
}

func (host *ScriptHost) Dispose() {
	var undiposedContexts []*ScriptContext
	for _, ctx := range host.contexts {
		undiposedContexts = append(undiposedContexts, ctx)
	}
	undisposedCount := len(undiposedContexts)

	if undisposedCount > 0 {
		slog.Warn("Script host shutdown: Not all contexts disposed", "count", len(host.contexts))
		for _, ctx := range undiposedContexts {
			ctx.Dispose()
		}
	}
	host.inspectorClient.Dispose()
	host.inspector.Dispose()
	host.iso.Dispose()
}

var global *v8.Object

func (host *ScriptHost) NewContext(window html.Window) *ScriptContext {
	context := &ScriptContext{
		host:     host,
		v8ctx:    v8.NewContext(host.iso, host.windowTemplate),
		window:   window,
		v8nodes:  make(map[ObjectId]*v8.Value),
		domNodes: make(map[ObjectId]Entity),
	}
	host.inspector.ContextCreated(context.v8ctx)
	err := installPolyfills(context)
	if err != nil {
		// TODO: Handle
		panic(err)
	}
	global = context.v8ctx.Global()
	errorCallback := func(err error) {
		window.DispatchEvent(NewCustomEvent("error"))
	}
	context.eventLoop = NewEventLoop(global, errorCallback)
	host.contexts[context.v8ctx] = context
	context.CacheNode(global, window)
	context.disposers = append(context.disposers, context.eventLoop.Start())

	return context
}

type Wrapper ScriptHost

func (w *Wrapper) NewScriptEngine(window html.Window) html.ScriptEngine {
	host := (*ScriptHost)(w)
	return host.NewContext(window)
}

func (w *Wrapper) Dispose() {
	host := (*ScriptHost)(w)
	host.Dispose()
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func (ctx *ScriptContext) Dispose() {
	ctx.host.inspector.ContextDestroyed(ctx.v8ctx)
	slog.Debug("ScriptContext: Dispose")
	for _, dispose := range ctx.disposers {
		dispose.Dispose()
	}
	ctx.pinner.Unpin()
	// TODO: Synchronize
	delete(ctx.host.contexts, ctx.v8ctx)
	ctx.v8ctx.Close()
}

func (ctx *ScriptContext) AddDisposer(disposer Disposable) {
	ctx.disposers = append(ctx.disposers, disposer)
}

func (ctx *ScriptContext) RunScript(script string) (*v8.Value, error) {
	return ctx.v8ctx.RunScript(script, "")
}

func (ctx *ScriptContext) Run(script string) error {
	_, err := ctx.RunScript(script)
	return err
}

func (ctx *ScriptContext) Eval(script string) (interface{}, error) {
	result, err := ctx.RunScript(script)
	if err == nil {
		return v8ValueToGoValue(result)
	}
	return nil, err
}

func (ctx *ScriptContext) Window() html.Window {
	return ctx.window
}

func v8ValueToGoValue(result *v8go.Value) (interface{}, error) {
	if result == nil {
		return nil, nil
	}
	if result.IsBoolean() {
		return result.Boolean(), nil
	}
	if result.IsInt32() {
		return result.Int32(), nil
	}
	if result.IsString() {
		return result.String(), nil
	}
	if result.IsNull() {
		return nil, nil
	}
	if result.IsUndefined() {
		return nil, nil
	}
	return nil, fmt.Errorf("Value not yet supported: %v", *result)
}
