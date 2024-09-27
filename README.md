# ev7ab

Evolution laboratory of genetic algorithm.

Lab generates and variates parameters of neural network to reach the target. To do this need to prepare a Producer which defines execution conditions. Producer rates result value of the solution by set of values or algorithm, like game, or whatever.

To try lab in action build one of the examples, that utilize a ready-to-use application.


## Documentation

[![go.dev reference](https://pkg.go.dev/badge/github.com/at7as/ev7ab/)](https://pkg.go.dev/github.com/at7as/ev7ab)


## Installation

```bash
go get -u github.com/at7as/ev7ab/...
```


## Usage

First of all you need to implement Producer interface.
Than create lab instance, add project layout and run examine, like this:
```go
l := lab.New(&producer{})
l.ProjectAdd(project)
l.Run()
```
But this example does not work as expected because lab.Run() execute goroutine.
For really works complex examples see [next section](#examples).


## Examples

[App](https://github.com/at7as/ev7ab/tree/master/app) package is an example for using [lab](https://github.com/at7as/ev7ab/tree/master/lab) package.
[Examples](https://github.com/at7as/ev7ab/tree/master/examples) dir contains several solutions for different typical tasks.


## Roadmap

Some features may be released later:
- Convolution Node type
- Linear Node type
- Backpropogation
- Gradient descent


## License

This project is licensed under the terms of the BSD-style
license that can be found in the [LICENSE](https://github.com/at7as/ev7ab/blob/master/LICENSE) file.
