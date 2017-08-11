# merr

multiple error

## Installation

```sh
$ go get github.com/gnue/goutils/merr
```

## Usage

```go
import "github.com/gnue/goutils/merr"
```

## Examples

### Errors

```go
if err != nil {
    if errs, ok := err.(merr.Errors); ok {
        for _, err := range errs.Errors() {
            fmt.Println(err)
        }
    }
}
```

get error list

### New

```go
var errs []error

for _, s := range files {
    if _, err := os.Stat(s); err != nil {
        errs = append(errs, err)
    }
}

err = merr.New(errs...)
```

create multiple error

## Author

[gnue](https://github.com/gnue)

## License

[MIT](LICENSE.txt)

