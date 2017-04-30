errors
----

`errors` adds a richness to the error type that allows errors to be easily
composed and extended. `errors` was created to help with sitations where a
complex process may return an error such as: "EOF". Using `errors`, you can
extend errors with context each time you receive them using `Extend`. The
`Contains` function will inform you if a rich error contains the provided error.
`IsOSNotExist` will return true if any of the underlying errors of a rich error
return true for `os.IsNotExist`.

Example:

```go
var errOne = errors.New("one")
var errTwo = errors.New("two")
_, errDNE = os.Open("file.txt")

extended := errors.Extend(errOne, errDNE)    // "[one; open file.txt: no such file or directory]"
extended2 := errors.Extend(errTwo, extended) // "[two; one; open file.txt: no such file or directory]"
errors.IsOSNotExist(extended)                // true
errors.Contains(extended, errDNE)            // true
errors.Contains(extended, errOne)            // true
errors.Contains(extended, errTwo)            // false
errors.Contains(extended2, errTwo)           // true
```

The `Compose` function works similarly to extend, however instead of simply
combining the two errors into a joined set, a new RichErr is created that
contains both underlying errors individually. The `Contains` and `IsOSNotExist`
functions will still work as expected, returning true if they apply to any of
the nested underlying errors.

```go
var errOne = errors.New("one")
var errTwo = errors.New("two")
composed := errors.Compose(errOne, errTwo)                       // "[[one]; [two]]"
composed2 := errors.Compose(errOne, errors.Extend(errTwo, errTwo)) // "[[one]; [two; two]]"
errors.Contains(composed2, errTwo) // true
```
