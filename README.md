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
$ go-guard func NewBar
// guardNewBar allows to guard NewBar constructor.
func guardNewBar(foo Fooer, r io.Reader) {
	guard.MustNotNil(1, "foo", foo)
	guard.MustNotNil(2, "r", r)
}
```

```sh
$ go-guard call NewStr
guardNewBar(foo, r)
```

`guard` import path is `github.com/romanyx/guard`

You can use go-guard from Vim with [vim-go-guard](http://github.com/romanyx/vim-go-guard) plugin
