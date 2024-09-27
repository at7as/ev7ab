# ev7ab / examples

This dir contains several solutions for different typical tasks.

- [simple](https://github.com/at7as/ev7ab/tree/master/examples/simple): basic classification task
- [bezier](https://github.com/at7as/ev7ab/tree/master/examples/bezier): recover control points by curve
- [tictactoe](https://github.com/at7as/ev7ab/tree/master/examples/tictactoe): train game ai in duel mode
- [track](https://github.com/at7as/ev7ab/tree/master/examples/track): search fast path
- [digits](https://github.com/at7as/ev7ab/tree/master/examples/digits): an attempt to solve the MNIST task

To run example use (e.g. simple):
```bash
make simple-app
```
to check result:
```bash
make simple-try
```
if `make` is not available:
```bash
go run ./examples/simple/main.go -config=./examples/simple/app.config.json
```
and
```bash
go run ./examples/simple/main.go -try
```