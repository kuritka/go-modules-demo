# go modules

The demo is testing go modules 



### init

initialize  project 
```bash
go mod init github.com/kuritka/gomodules
```

which creates 
```bash

module github.com/kuritka/gomodules

go 1.12

```


after downloading  `go get github.com/gorilla/mux` file `go.sum` will appear and 
`require  github.com/gorilla/mux v1.7.3 //indirect` will appear.
`//indirect` says that dependency is not used yet 

