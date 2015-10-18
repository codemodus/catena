# catena

    go get "github.com/codemodus/catena"

Package catena aids the composition of http.Handler wrapper catenas.

Nesting functions is a simple concept.  If your handler wrapper order does not 
need to be composable, do not use a package and avoid adding a dependency 
to your project.  However, nesting functions quickly becomes burdensome as the 
need for flexibility increases. Enter Catena.

## Usage

```go
type Catena
    func New(hws ...func(http.Handler) http.Handler) Catena
    func (c Catena) Append(hws ...func(http.Handler) http.Handler) Catena
    func (c Catena) End(h http.Handler) http.Handler
    func (c Catena) EndFn(h http.HandlerFunc) http.Handler
    func (c Catena) Merge(cs ...Catena) Catena
```

### Setup

```go
import (
    // ...

    "github.com/codemodus/catena"
)

func main() {
    // ...

    cat0 := catena.New(firstWrapper, secondWrapper)
    cat1 := cat0.Append(handlerWrapper, fourthWrapper)

    cat2 := catena.New(beforeFirstWrapper)
    cat2 = cat2.Merge(cat1)

    m := http.NewServeMux()
    m.Handle("/1w2w_End1", cat0.EndFn(ctxHandler))
    m.Handle("/1w2w_End2", cat0.EndFn(anotherCtxHandler))
    m.Handle("/1w2wHw4w_End1", cat1.EndFn(ctxHandler))
    m.Handle("/0w1w2wHw4w_End1", cat2.EndFn(ctxHandler))

    // ...
}
```

### http.Handler Wrapper

```go
func firstWrapper(n http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // ...

        n.ServeHTTP(w, r)

        // ...
    })
}
```
This function signature will make wrappers compatible with catena.

End-point functions will need to be adapted using http.HandlerFunc.  As a 
convenience, EndFn will adapt functions with compatible signatures.

## More Info

### chain

If a project ends up needing a request context, consider using 
[chain](https://github.com/codemodus/chain). The API is nearly identical to 
catena, so adding a request context is easy.

## Documentation

View the [GoDoc](http://godoc.org/github.com/codemodus/catena)

## Benchmarks

These results are for comparison of normally nested functions, and catenated 
functions.  Each benchmark includes 10 functions prior to the final handler.

    benchmark            iter      time/iter   bytes alloc         allocs
    ---------            ----      ---------   -----------         ------
    BenchmarkCatena10   20000    56.89 μs/op     3534 B/op   54 allocs/op
    BenchmarkNest10     20000    56.56 μs/op     3511 B/op   54 allocs/op
