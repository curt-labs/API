CURT API v3 - Setup Local Development
=========
---------
This document explains the steps to setup a developer environment for the Go API project.
> This document assumes you have already setup your local databases.
 
Required Access:
-
- GitHub Access
- Dealer Account API Key ([http://dealers.curtmfg.com/](http://dealers.curtmfg.com/))

Required Software:
-
- go

	`$ brew install go`
- mongodb
	
 	`$ brew install mongodb`
- redis

 	`$ brew install redis`
- mysql

 	`$ brew install mysql`

Setup Go Environment
-
Go requires a directory structure and setting up environment variables for a development environment. Go Documentation can be found here [https://golang.org/doc/code.html](https://golang.org/doc/code.html).
 
 Here is our short and sweet by example explanation
 
- Create "GOPATH" directory 
    
	`~ $ mkdir -p workspace/gocode/`

- Create source code path that will be used for the API project

 	`~ $ mkdir -p workspace/gocode/src/github.com/curt-labs/`
 	
- Get "GOPATH" path
 	
 	`~ $ cd workspace/gocode`
 	
 	```
 	~/workspace/gocode $ pwd
 	
 	<pwd path>
	```
- Add to `~/.bash_profile`

	```
	# Go
	export GOPATH=<pwd path>
	export PATH=$PATH:$GOPATH/bin
	```
- Add the environment variables to current session
	`~ $ source .bash_profile`

Checkout API Project Code
-
- Go to the curt-labs directory inside your Go source directory
`~/workspace/gocode/src/github.com/curt-labs/ $ git clone https://github.com/curt-labs/API.git`


Start Application
-
- Checkout Working Branch (goapi)
	
	`$ git checkout goapi`
- Install dependencies

	`$ go get`

- Start Application

	`$ go run index.go`