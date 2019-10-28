# Changelog

## v1.3.1
* Moved `conf.Options` out of internal package.

## v1.3.0

* Added support for functions with variadic arguments.
* Added fast call opcode for special case functions.

## v1.2.0
* Fixed comparison between nil and simple types (int, bool, string).
* Fixed equal operation on different nil types.
* Fixed bug with nil vars, now is possible to compare fields with `nil` literal.
* Fixed default type of map created in expr.
* Fixed type checker for interface's methods.
* Added bytecode virtual machine.
* Added optimizing compiler.
* Added builtin functions: all, none, any, one, filter, map.
* Added operator overloading.
* Improved error messages.

## v1.1.4
* Add support for method on maps.

## v1.1.3
* Improved speed of &&, || operators.
* Added support for go1.8.
* Added more benchmarks.

## v1.1.2
* Fixed work with field functions.

## v1.1.1
* Fixed work of methods with pointer receiver.

## v1.1.0
* Added support for struct's methods as functions.
* Fixed type checker for `~` operator.
* Fixed type checker for `?:` operator.

## v1.0.7
* Refactored embedded structs to types table conversion.
* Deprecated ~expr.With~, use _expr.Env_ instead.

## v1.0.6
* Added type checks for `in` operator.

## v1.0.5
* Fixed unquoting strings.

## v1.0.4
* Added support for embedded structs.
* Fixed type checks for functions and methods.

## v1.0.3
* Added possibility to use `in` operator with maps and structs.

## v1.0.2
* Fixed type checks for binary operators.
* Improved error messages for nil dereference.

## v1.0.1
* Added more type checks.

## v1.0.0
* Added advance type checks in strict mode.

## v0.0.5
* Improved error position reporting.

## v0.0.4
* Added regex compile while parse if possible.

## v0.0.3
* Added len builtin.

## v0.0.2
* Added pointer dereference for field extraction.

## v0.0.1
* First release.
