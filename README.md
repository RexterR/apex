# apex
A super lightweight and customizable minimalistic HTTP framework<br>
[![Go Report Card](https://goreportcard.com/badge/github.com/RexterR/apex)](https://goreportcard.com/report/github.com/RexterR/apex)

`Apex` is a BYO*, minimalistic web framework that builds on top of your router of choice and adds some additional functionality—namely, middleware chaining and route grouping. It's meant to be used on projects large and small that require flexibility, and varying degrees of custom code and architecture.
This package takes some inspiration from design decisions in
[`chi`](https://github.com/pressly/chi) and
[`gin`](https://github.com/gin-gonic/gin).

## Usage

### Hello World

You can initialize a new instance of the `Apex` type with whichever type that
implements `apex.Handler`. An adapter for
[`httprouter`](https://github.com/julienschmidt/httprouter) is included.

``` go
import (
    "fmt"
    "net/http"

    "github.com/julienschmidt/httprouter"

    "github.com/RexterR/apex"
    "github.com/RexterR/apex/adapter"
)

func main() {
    a := &adapter.Httprouter{Router: httprouter.New()}
    m := apex.New(a)

    m.Get("/", helloWorld)

    http.ListenAndServe(":8080", m)
}

func helloWorld(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "hello world!")
}
```

### Route Parameters

`apex` supports all the syntax variations for defining route parameters that
the underlying router does. For instance, in the case of `httprouter`:

```go
import (
    "fmt"
    "net/http"

    "github.com/julienschmidt/httprouter"

    "github.com/RexterR/apex"
    "github.com/RexterR/apex/adapter"
)

func main() {
    a := &adapter.Httprouter{Router: httprouter.New()}
    m := apex.New(a)

    m.Get("/:name", greet)

    http.ListenAndServe(":8080", m)
}

func greet(w http.ResponseWriter, r *http.Request) {
    name := httprouter.ParamsFromContext(r.Context()).ByName("name")
    fmt.Fprintf(w, "hello %s!", name)
}
```

### Route Grouping

``` go
import (
    "fmt"
    "net/http"

    "github.com/julienschmidt/httprouter"

    "github.com/RexterR/apex"
    "github.com/RexterR/apex/adapter"
)

func main() {
    a := &adapter.Httprouter{Router: httprouter.New()}
    m := apex.New(a)

    apiRouter := m.NewGroup("/api")
    {
        // GET /api
        apiRouter.Get("/", apiRoot)
        // GET /api/ignacio
        apiRouter.Get("/:name", greet)
    }

    http.ListenAndServe(":8080", m)
}

func apiRoot(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "api root")
}

func greet(w http.ResponseWriter, r *http.Request) {
    name := httprouter.ParamsFromContext(r.Context()).ByName("name")
    fmt.Fprintf(w, "hello %s!", name)
}
```

### Middleware

Middleware in `apex` are simply functions that take an `http.Handler` (the one
next in the chain) and return another one. They are resolved in the order that
they are chained. You can chain them together with the `Middleware.Then`
method.

`apex` users are meant to take advantage of `context` to make better use of
middleware.

``` go
import (
    "context"
    "fmt"
    "log"
    "net/http"

    "github.com/julienschmidt/httprouter"

    "github.com/RexterR/apex"
    "github.com/RexterR/apex/adapter"
)

func main() {
    a := &adapter.Httprouter{Router: httprouter.New()}
    m := apex.New(a)

    chain := apex.Middleware(logger).Then(printer)
    m.Use(chain)

    apiRouter := m.NewGroup("/api")
    {
        apiRouter.Get("/", apiRoot)
        nameRouter := apiRouter.NewGroup("/:name")
        {
            // Every request sent to routes defined on this sub-router will now
            // have a reference to a name in its context.
            // Useful for RESTful design.
            nameRouter.Use(nameExtractor)

            // GET /api/ignacio
            nameRouter.Get("/", greet)
            // GET /api/ignacio/goodbye
            nameRouter.Get("/goodbye", goodbye)
        }
    }

    http.ListenAndServe(":8080", m)
}

// -- Middleware --

// a simple logger
func logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("| %s %s", r.Method, r.URL)
        next.ServeHTTP(w, r)
    })
}

// a useless middleware that prints text
func printer(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println("this prints some text")
        next.ServeHTTP(w, r)
    })
}

// extracts a name from the URL and injects it into the request's context
func nameExtractor(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        name := httprouter.ParamsFromContext(r.Context()).ByName("name")
        ctx := context.WithValue(r.Context(), "name", name)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// -- Handlers --

func apiRoot(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "api root")
}

// greets the user with :name
func greet(w http.ResponseWriter, r *http.Request) {
    name := r.Context().Value("name").(string)
    fmt.Fprintf(w, "hello %s!", name)
}

// says "bye" to the user with :name
func goodbye(w http.ResponseWriter, r *http.Request) {
    name := r.Context().Value("name").(string)
    fmt.Fprintf(w, "bye %s!", name)
}
```
