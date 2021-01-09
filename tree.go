package validapi

import (
	"net/http"
)

type (

	//node represents a single position on the tree. its created as an interface
	// to allow both route and router to serve as stops.
	node interface {
		//find should accept a string array and locate the desired route handler,
		// wrapped in an array for the purpose of creating a middleware chain.
		find([]string) ([]*Middleware, *leaf)
		//add should implement the logic for creating new nodes and children.
		add([]string)
	}

	//leaf represents a combination of handler and propertygroup. It serves as the end of the tree.
	leaf struct{}

	//Route represents a position on the tree that can
	Route struct {
		Path           string
		handlers       map[string]*leaf
		staticChildren map[string]*node
		variableChild  *leaf
		defaultChild   *leaf
	}

	//Router an extension of Route that contains a list of middleware functions.
	Router struct {
		Route
		middleware []func(http.HandlerFunc) http.HandlerFunc
	}

	//Middleware wrapper type for the func(http.Handlerfunc) http.HandlerFunc
	Middleware func(http.HandlerFunc) http.HandlerFunc
)

func (r *Route) add(uri []string) {
	panic("not implemented")
}

func (r *Route) find(uri []string) ([]*Middleware, *leaf) {
	panic("not implemented")
	// return []*http.HandlerFunc{}, &leaf{}
}

func (r *Router) add(uri []string) {
	panic("not implemented")
}

func (r *Router) find(uri []string) ([]*Middleware, *leaf) {
	panic("not implemented")
	// return []*http.HandlerFunc{}, &leaf{}
}
