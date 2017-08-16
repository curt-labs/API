package mocks

import (
)

// use the default pool for dockertest.
// @see [dockertest.NewPool()](https://github.com/ory/dockertest/blob/v3/dockertest.go#L63)
//const _DEFAULT_DOCKER_POOL = ""

//var mysqlContainer = DockertestContainer{Image: "mysql",
//	Tag: "5.6",
//	Host: "localhost",
//	Port:PublishedPort{Number: 3306, Protocol:"tcp"}}
//
// Connection information for the running MongoDB container
//type DockertestMySql struct {
//	Pool *dockertest.Pool
//	Resource *dockertest.Resource
//	Session *sql.DB
//}

//func NewSessionFromResource(pool *dockertest.Pool, mgoResource *dockertest.Resource) (session *mgo.Session, err error) {
	//err = mgoPool.Retry(attemptDialWith(session, mgoResource))
	//err = pool.Retry(func () (err error) {
	//	session, err = mgo.Dial(fmt.Sprintf("%s:%s", HOST, mgoResource.GetPort(PORT)))
	//	if err != nil {
	//		return err
	//	}
	//
	//	return session.Ping()
	//})
	//return session, err
//}
