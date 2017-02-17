# catena

    go get -u github.com/codemodus/catena

Package catena aides gRPC interceptor catenation.

## Usage

```go
type UnaryServerCatena
    func NewUnaryServerCatena(is ...grpc.UnaryServerInterceptor) *UnaryServerCatena
    func (c *UnaryServerCatena) Append(is ...grpc.UnaryServerInterceptor) *UnaryServerCatena
    func (c *UnaryServerCatena) Copy(catena *UnaryServerCatena)
    func (c *UnaryServerCatena) Interceptor() grpc.UnaryServerInterceptor
    func (c *UnaryServerCatena) Merge(cs ...*UnaryServerCatena) *UnaryServerCatena
    func (c *UnaryServerCatena) ServerOption() grpc.ServerOption
```

### Setup

```go
soon...
```

## More Info

## Documentation

View the [GoDoc](http://godoc.org/github.com/codemodus/catena)
