# execuser

enchant with RunUser

## Installation

```sh
$ go get github.com/gnue/goutils/execuser
```

## Usage

```go
import "github.com/gnue/goutils/execuser"
```

## Examples

### Lookup

```go
cmd := exec.Command("ps", "u")

if isRootUser() {
    if u, err := execuser.Lookup("username"); err == nil {
        u.RunUser(cmd)
    }
}

cmd.Run()
```

run user by username

### LookupPath

```go
cmd := exec.Command("ps", "u")

if isRootUser() {
    if u, err := execuser.LookupPath("."); err == nil {
        u.RunUser(cmd)
    }
}

cmd.Run()
```

run user by path

## Author

[gnue](https://github.com/gnue)

## License

[MIT](LICENSE.txt)

