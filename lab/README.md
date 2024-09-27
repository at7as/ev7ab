# ev7ab / lab

Package lab is engine of genetic algorithm.


## Documentation

[![go.dev reference](https://pkg.go.dev/badge/github.com/at7as/ev7ab/)](https://pkg.go.dev/github.com/at7as/ev7ab/lab)


## Usage

First of all you need to implement Producer interface.
Than create lab instance, add project layout and run examine, like this:
```go
l := lab.New(&producer{})
l.ProjectAdd(project)
l.Run()
```
But this example does not work as expected because lab.Run() execute goroutine.

[App](https://github.com/at7as/ev7ab/tree/master/app) package is an example for using [lab](https://github.com/at7as/ev7ab/tree/master/lab) package.
