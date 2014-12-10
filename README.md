
CURT Go API
=========
---------

> The new version of the CURT API used the [GoEngine Boilerplate](http://github.com/ninnemana/goengine-boilerplate)
for being Content-Type agnostic to XML and JSON. Some of the best features are listed below:

  - Concurrent MySQL access using [Goroutines](http://golang.org/doc/effective_go.html#concurrency)
  - JSON rendering powered by [encoding/json](http://golang.org/pkg/encoding/json/)
  - XML rendering powered by [encoding/xml](http://golang.org/pkg/encoding/xml/)
  - MySQL Persistence using [mymysql](https://github.com/ziutek/mymysql)
  - ACES Compliant vehicle lookup with product groups


--------
Endpoints
---------
---------

> Note: this application is still in heavy development and all endpoints/objects have the potential to change at any time.

> You can view example endpoints for all of the routes in the index_test.go file.

#### Vehicle

---

https://github.com/curt-labs/GoAPI/blob/master/docs/Vehicle.md

---
#### Parts

---

*Get Part by Part #

    GET - http://goapi.curtmfg.com/part/110003?key=[public api key]

*Reverse Lookup by Part #

    GET - http://goapi.curtmfg.com/part/110003/vehicles?key=[public api key]

----

#### Categories

---

https://github.com/curt-labs/GoAPI/blob/master/docs/Categories.md

----

#### Customer

---

*Authentication*

    POST - http://goapi.curtmfg.com/customer/auth

    Payload
    --------------------------
    email: user@example.com
    password: password

> The following GET route for the customer user authentication is only useful if in the last 6 hours this user has logged in through the POST directive of the /customer/auth endpoint.

    GET - http://goapi.curtmfg.com/customer/auth?key=c8bd5d89-8d16-11e2-801f-00155d47bb0a

*Customer Locations*

    POST - http://goapi.curtmfg.com/customer/locations

    Payload
    --------------------------
    key: CEB28F99-F03A-4568-B004-E4FFA87CBDF1

*Customer Users*

    POST - http://goapi.curtmfg.com/customer/users

    Payload
    --------------------------
    key: CEB28F99-F03A-4568-B004-E4FFA87CBDF1

> The customer users endpoint will only return data if the requesting user is marked as sudo user

Philoshopy
-

> This version if the API is meant to focus on data quantity while maintaining, if not improving performance, by leveraging concurrency. We would like the client to have the ability to make fewer requests to the API Server and be provided with a larger amount of data in the response.

Deployment
-

Deployment will be done using the master branch on Github. Once a commit is pushed to Github
it will route that commit to Drone.io, which will then running CI testing across
the project and then deploy new Docker containers to all CURT servers.

Contributors
-
* Alex Ninneman
    * [Github](http://github.com/ninnemana)
    * [Twitter](https://twitter.com/ninnemana)
* David Vaini
    * [Github](https://github.com/DavidVaini)
* John Shenk
    * [Github](https://github.com/stinkyfingers)
* Matt Mickelson
    * [Github](https://github.com/mickelsonm)

License
-

MIT

*Free Software, Fuck Yeah!*

