# Exprgen

## Install
```
go install github.com/antonmedv/expr/exprgen
```
## Usage
Fetch methods generates for all struct/map/array/string named types(exception is map types with unnamed not basic key type like `map[struct{...}]int`).

To generate just call exprgen with pkg paths as arguments: 
```
exprgen pkg1 pkg2 ...
```

After call, file `*pkg_name*_exprgen.go` will be created in each packages from arguments.
