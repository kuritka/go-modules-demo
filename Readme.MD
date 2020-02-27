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
`//indirect` says that dependency is not used yet. After using mux and rebuild , `//indirect` dismiss 



 - `go list` - shows module we are working on

 - `go list all` - shows all dependencies which we rely

 - `go list -m all` - shows all modules we use and their versions
     ```bash
    github.com/kuritka/gomodules
    github.com/gorilla/mux v1.7.4
     ```
 
 - `go list -m -versions github.com/gorilla/mux` touch repo and shows all available versions, i.e.
 `github.com/gorilla/mux v1.2.0 v1.3.0 v1.4.0 v1.5.0 v1.6.0 v1.6.1 v1.6.2 v1.7.0 v1.7.1 v1.7.2 v1.7.3 v1.7.4`
 
 
 ## go.sum
 
 go.sum contains info like :
 ```bash
github.com/gorilla/mux v1.7.4 h1:VuZ8uybHlWmqV03+zRzdwKL4tUnIp1MAQtp1mIFE1bc=
github.com/gorilla/mux v1.7.4/go.mod h1:DVbg23sWSpFRCP0SfiEN6jmj59UnW/n46BH5rLB71So=
``` 

hash is retreived after go downloaded dependency. It is uniqueue based on source code.
We can verify by `go mod verify` 


`%GOPATH%/pkg/mod/` - after `go get ` package is downloaded here. This is module chash . It is loaded here and build system looks here for downloaded packages.

If I modify file manually within `%GOPATH%/pkg/mod/github.com/gorilla.MUX/..` than `go verify` will generate another hash and saying that's wrong.

BUT: `go build .` will still work!! , thre reason why verifying is splitted from building is that **verifying is expensive operation** because all dependencies must be recalculated.  
hash is also computed from file date, so even if you have two identical files with different date, hash will differ! 

### it means that in the build pipeline `go mod verify` must be called before `go build` in PROD 
     
Only one way how to fix corrupted module is to delete from  `%GOPATH%/pkg/mod/` and than tun `go get \ go mod vendor` again.
    
     
## managing modules

Let's imagine that after refactoring app you don't need some library packages anymore (i.e. remove gorilla/mux).
 After investigate `go build .` you find that `go.mod` is not modified and and references are still there. And reference is not `//indirect`

go.mod looks like this:  
```bash
module github.com/kuritka/gomodules

//gomodules will be the name of executable

go 1.12

require github.com/gorilla/mux v1.7.4
```

To remove unused dependencies we call 

```bash
go mod tidy
```
        
        
## Versions

`v1.5.3-pre1` - true accross all module version

 - `v` - required
 
 - `1` - **major version**, within major version we guarantee backward compatibility! Anything with major version 1 should be compatible with major version 1
        
 - `5` - **minor revision** - if we add new features (no changes in current API, don't break existing contracts), Doesn't break backward compatibility
 
 - `3` - **patch** - bug fixes, no new features, doesn't break backward compatibility
  
 - `pre1` - **pre-release identifier** - text is arbitrary. ordered in alphabetical manner. `pre1` is less than `pre2` and `alfa` is less then `beta`. If we don't specify conrete version with identifier, lates version without identifier is downloaded by default.
 
    ```bash
    go list -m -versions rsc.io/quote
    rsc.io/quote v1.0.0 v1.1.0 v1.2.0 v1.2.1 v1.3.0 v1.4.0 v1.5.0 v1.5.1 v1.5.2 v1.5.3-pre1
    ```
   
   ```bash
   go get rsc.io/quote
   # downloads v1.5.2 as default 
   ```
   
   
     
 
 https://semver.org  
 
 i.e. `ithub.com/gorilla/mux v1.2.0 v1.3.0 v1.4.0 v1.5.0 v1.6.0 v1.6.1 v1.6.2 v1.7.0 v1.7.1 v1.7.2 v1.7.3 v1.7.4`
 
 ---
 
 ### rules
 There are still a lot of libraries out there which didn't follow go versioning rules.
 
####rules for version v2+####
 - backward compatibility should be preserved within major version
 
 - Each major vesrion has unique import path:
    `import github.com/gorrile/mux/v2`  
    
    btw: it is just illustration, gorrila doesnt have v2
    ```bash
        $ go list -m -versions github.com/gorilla/mux/v2
        go: finding github.com/gorilla/mux/v2 latest
        github.com/gorilla/mux/v2
    ```
   
   
```bash
go list -m -versions rsc.io/quote
rsc.io/quote v1.0.0 v1.1.0 v1.2.0 v1.2.1 v1.3.0 v1.4.0 v1.5.0 v1.5.1 v1.5.2 v1.5.3-pre1```

go list -m -versions rsc.io/quote/v2
rsc.io/quote/v2 v2.0.0 v2.0.1

go get rsc.io/quote
# within go.mod you find this: require rsc.io/quote v1.5.2 // indirect
```

## regardiong example - if you don't explicitly give version 2 or higher, it downloads always latest from v1 !!## 

you may have multiple versions simultaniously  
```bash
<go.mod>
    require (
	    rsc.io/quote v1.5.2
	    rsc.io/quote/v2 v2.0.1
    )

<go.mod>
    import (
    	"rsc.io/quote"
    	quotev2 "rsc.io/quote/v2"
    )

<main.go>
	fmt.Println(quote.Hello())
	fmt.Println(quotev2.Hello())
```

We can pull
 
 - specific version `@v1.7.2` 
 - version prefix `@v1`
 - Latest `@latest`;  `go get github.com/gorilla/mux@latest`
 - Specific commit `@c9b2812`
 - Specific commit `@master` ; `go get github.com/gorilla/mux@master` - shows v1.7.4-c9282.. because it presunmes that newer version is 1.7.4 (command was exexcuted during time when 1.7.3 was the highest)  
 - Comparison `@>=1.7.2`


I can download gorilla mux by `go get github.com/gorilla/mux@v1.7.2` and `go get github.com/gorilla/mux@v1.6.1`
What will be in go mod file ?
First `go.mod` contains `github.com/gorilla/mux@ v1.7.2` and than it is changed to `github.com/gorilla/mux@ v1.6.1` .
If the major version is same it means that there is backward compatibility 

  
## Advanced go mod techniquest

### get

`go mod get` downloads package from repo or cache and modify `go.mod`, `go.sum`

### why
analytical tool  explores where module comes from

This is shown for `// indirect` 
```bash
go mod why github.com/gorilla/mux
(main module does not need package github.com/gorilla/mux)
```


This is shown for `// indirect` 
```bash
go mod why github.com/gorilla/mux
(main module does not need package github.com/gorilla/mux)
```

and dependency graph is shown if gorrila mux is used :
```bash
go mod why github.com/gorilla/mux
# github.com/gorilla/mux
github.com/kuritka/gomodules
github.com/gorilla/mux
```

### graph

graph shows what is dependent on the left and dependencies on the right.
if package has dependencies, it is shown in the table. In the example below you see  rsc.io has dependencies. 
```bash
go mod graph
github.com/kuritka/gomodules github.com/gorilla/mux@v1.7.4
github.com/kuritka/gomodules rsc.io/quote@v1.5.2
rsc.io/quote@v1.5.2 rsc.io/sampler@v1.3.0
rsc.io/sampler@v1.3.0 golang.org/x/text@v0.0.0-20170915032832-14c0d48ead0c
```

### edit
 
`go mod edit` is commad used for scripting go mod file

- `go mod edit -require github.com/gorilla/context@v1.1.1` adds  line to require
- `go mod edit -go 1.10` change version in `go.mod` file to `go 1.10`
- `go mod edit -droprequire github.com/gorilla/context` drops line where `github.com/gorilla/context v1.1.1` appear
- `go mod edit -exclude  github.com/gorilla/mux@v1.7.2` excludes `gorilla/mux@v1.7.2`. If we downloadin `gorilla/mux@v>=1.7.1` it excludes `1.7.2` so `1.7.3` will be downloaded
- `go mod edit -print` prints mod file
- `go mod help` shows all posisible commands


### download

`go mod download` downloads all dependencies from `go.mod` file to `module cache`

### vendor
analyze project , and finds all references. Creates `vendor` directory and populate that. 
`go mod vendor` vendors all repos. If you don't want to download from internet in production pipeline, or you are out of `%GOPATH%` , you can vendor your dependencies.

If you use `vendor` directories, you import all references from `vendor` . Not some of them from `vendor` and some of them from `module cache`.

`vendor` is good thet you can deploy exact package versions with your code without downloading from the internet

because `vendor` analyze your references, it also updates `go.mod` file automatically by package version it downloads


### readonly 

`go build -mod=readonly .` says that nobody can modify mod file when build is under way. reccomendation is to use that in production build .

