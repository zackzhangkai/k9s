package view

import (
	"context"
	"strings"

	"github.com/derailed/k9s/internal/resource"
	"github.com/derailed/k9s/internal/ui"
)

// Context presents a context viewer.
type Context struct {
	*Resource
}

// NewContext return a new context viewer.
func NewContext(title, gvr string, list resource.List) ResourceViewer {
	return &Context{
		Resource: NewResource(title, gvr, list),
	}
}

func (c *Context) Init(ctx context.Context) {
	c.extraActionsFn = c.extraActions
	c.enterFn = c.useCtx
	c.Resource.Init(ctx)

	c.masterPage().SetSelectedFn(c.cleanser)
}

func (c *Context) extraActions(aa ui.KeyActions) {
	c.masterPage().RmAction(ui.KeyShiftA)
}

func (c *Context) useCtx(app *App, _, res, sel string) {
	if err := c.useContext(sel); err != nil {
		app.Flash().Err(err)
		return
	}
	app.gotoResource("po", true)
}

func (*Context) cleanser(s string) string {
	name := strings.TrimSpace(s)
	if strings.HasSuffix(name, "*") {
		name = strings.TrimRight(name, "*")
	}
	if strings.HasSuffix(name, "(𝜟)") {
		name = strings.TrimRight(name, "(𝜟)")
	}
	return name
}

func (c *Context) useContext(name string) error {
	ctx := c.cleanser(name)
	if err := c.list.Resource().(*resource.Context).Switch(ctx); err != nil {
		return err
	}

	if err := c.app.switchCtx(name, false); err != nil {
		return err
	}
	c.refresh()
	c.masterPage().Select(1, 0)

	return nil
}
