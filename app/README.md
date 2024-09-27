# ev7ab / app

Package app provides basic usage of lab package with CUI.


## Documentation

[![go.dev reference](https://pkg.go.dev/badge/github.com/at7as/ev7ab/)](https://pkg.go.dev/github.com/at7as/ev7ab/app)


## Usage
```go
cfgFile := flag.String("config", "./app.config.json", "path to app config file")
flag.Parse()
app.Run(&producer{}, *cfgFile)
```
[Examples](https://github.com/at7as/ev7ab/tree/master/examples) dir contains several solutions for different typical tasks.
