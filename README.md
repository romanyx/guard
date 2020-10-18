# guard

`guard` allows to generate guard function for constructors which are checking possible nil inputs in parameters.

```go
type Fooer interface {
	Foo() error
}

type Bar struct {
	foo Fooer
	r io.Reader
}

func NewBar(foo Fooer, r io.Reader) Bar {
	bar := Bar {
		foo: foo,
		r: r,
	}

	return &bar
}
```

```sh
$ go-guard func NewStr
// guardNewBar allows to guard NewStr constructor.
func guardNewBar(foo Fooer, r io.Reader) {
	gcheck.MustNotNil(1, "foo", foo)
	gcheck.MustNotNil(2, "r", r)
}
```

```sh
$ go-guard call NewStr
guardNewBar(foo, r) {
```

`gcheck` import path is `github.com/romanyx/guard/gcheck`
