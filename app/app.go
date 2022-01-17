package app

import (
	"github.com/vugu/vgrouter"
	"github.com/vugu/vugu"
	"github.com/vugu-examples/simple/app/components"
	"github.com/vugu-examples/simple/app/pages"
)

type VuguSetupOptions struct {
	AutoReload bool
}

// VuguSetup performs UI setup and wiring.
func VuguSetup(buildEnv *vugu.BuildEnv, eventEnv vugu.EventEnv, opts *VuguSetupOptions) (*App, vugu.Builder) {

	if opts == nil {
		opts = &VuguSetupOptions{}
	}

	app := &App{
		Router:   vgrouter.New(eventEnv),
	}

	//router := vgrouter.New(eventEnv)

	pageMap := pages.MakeRoutes().WithRecursive(true).WithClean(true).Map()
	// pageSeq := &state.PageSeq{
	// 	PageMap:  pageMap,
	// 	PathList: SiteNavPathList,
	// }
	// app.PageSeq = pageSeq

	buildEnv.SetWireFunc(func(b vugu.Builder) {

		if c, ok := b.(vgrouter.NavigatorSetter); ok {
			c.NavigatorSet(app.Router)
		}

		// if c, ok := b.(state.PageInfoSetter); ok {
		// 	c.PageInfoSet(app.PageInfo)
		// }

		// if c, ok := b.(state.PageSeqSetter); ok {
		// 	c.PageSeqSet(app.PageSeq)
		// }

	})

	root := &components.Root{
		AutoReload: opts.AutoReload,
	}
	buildEnv.WireComponent(root)

	// pages - add automatically from generated routes
	for path, inst := range pageMap {
		instBuilder := inst.(vugu.Builder)
		app.Router.MustAddRouteExact(path, vgrouter.RouteHandlerFunc(func(rm *vgrouter.RouteMatch) {
			root.Body = instBuilder
		}))
	}

	app.Router.MustAddRouteExact("/", vgrouter.RouteHandlerFunc(func(rm *vgrouter.RouteMatch) {
		root.Body = &pages.Index{}
	}))

	app.Router.MustAddRoute("/detail", vgrouter.RouteHandlerFunc(func(rm *vgrouter.RouteMatch) {
		root.Body = &pages.Detail{}
	}))

	app.Router.MustAddRoute("/login", vgrouter.RouteHandlerFunc(func(rm *vgrouter.RouteMatch) {
		root.Body = &pages.Login{}
	}))

	app.Router.MustAddRoute("/register", vgrouter.RouteHandlerFunc(func(rm *vgrouter.RouteMatch) {
		root.Body = &pages.Register{}
	}))

	app.Router.MustAddRoute("/movies", vgrouter.RouteHandlerFunc(func(rm *vgrouter.RouteMatch) {
		root.Body = &pages.Movies{}
	}))

	if app.Router.BrowserAvail() {
		err := app.Router.ListenForPopState()
		if err != nil {
			panic(err)
		}

		err = app.Router.Pull()
		if err != nil {
			panic(err)
		}
	}

	return app, root
}

// App holds overall application state.
type App struct {
	*vgrouter.Router
}