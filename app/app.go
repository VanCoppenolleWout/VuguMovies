package app

import (
	"github.com/vugu-examples/simple/app/pages"
	"github.com/vugu-examples/simple/app/components"
	"github.com/vugu/vgrouter"
	"github.com/vugu/vugu"
)

type VuguSetupOptions struct {
	AutoReload bool
}

// OVERALL APPLICATION WIRING IN vuguSetup
func VuguSetup(buildEnv *vugu.BuildEnv, eventEnv vugu.EventEnv, opts *VuguSetupOptions) (*App, vugu.Builder) {
	if opts == nil {
		
	}
	
	app := &App{
		Router: vgrouter.New(eventEnv),
	}

	pageMap := pages.MakeRoutes().WithRecursive(true).WithClean(true).Map()

	// tl := state.LoadTacoListAPI()
	// ca := state.LoadCartAPI()
	// CREATE A NEW ROUTER INSTANCE
	//router := vgrouter.New(eventEnv)

	// MAKE OUR WIRE FUNCTION POPULATE ANYTHING THAT WANTS A "NAVIGATOR".
	buildEnv.SetWireFunc(func(b vugu.Builder) {
		if c, ok := b.(vgrouter.NavigatorSetter); ok {
			c.NavigatorSet(app.Router)
		}
		// if s, ok := b.(state.TacoListAPISetter); ok {
		// 	s.TacoListAPISet(tl)
		// }
		// if s, ok := b.(state.CartAPISetter); ok {
		// 	s.CartAPISet(ca)
		// }
	})

	// CREATE THE ROOT COMPONENT
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
	// router.MustAddRouteExact("/cart",
	// 	vgrouter.RouteHandlerFunc(func(rm *vgrouter.RouteMatch) {
	// 		root.Body = &pages.Cart{}
	// 	}))
	// router.MustAddRouteExact("/checkout",
	// 	vgrouter.RouteHandlerFunc(func(rm *vgrouter.RouteMatch) {
	// 		root.Body = &pages.Checkout{}
	// 	}))
	// router.SetNotFound(vgrouter.RouteHandlerFunc(
	// 	func(rm *vgrouter.RouteMatch) {
	// 		root.Body = &pages.PageNotFound{} // A PAGE FOR THE NOT-FOUND CASE
	// 	}))

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
