package utils

import (
	"net/http"
	"regexp"
)

// Struct route - handling the patternmatching with regular expressions
type Route struct {
	Pattern *regexp.Regexp
	Word    string
	Handler http.Handler
}

// Struct regexHandler - slice of route struct objects
type RegexHandler struct {
	Routes []*Route
}

// Function Handler with a RegexHandler object - Uses the RegexHandler object to append a new instance
// into the slice of routes
func (h *RegexHandler) Handler(r string, v string, handler func(http.ResponseWriter, *http.Request)) {
	// MustCompile is like Compile but panics if the expression cannot be parsed
	res := regexp.MustCompile(r)
	// Appends the route into the slice of routes
	h.Routes = append(h.Routes, &Route{res, v, http.HandlerFunc(handler)})
}

func (h *RegexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Loops through the Route objects
	for _, route := range h.Routes {
		// Pattern matching the URL path, if true and the "word" is a Methof, use the function
		// and then exit the function
		if route.Pattern.MatchString(r.URL.Path) && route.Word == r.Method {
			route.Handler.ServeHTTP(w, r)
			return
		}
	}
	// Gives back an error if not found
	http.NotFound(w, r)
}
