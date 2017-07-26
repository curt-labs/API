package database

import "fmt"

type MySqlConnectionInfo struct {
	Database string
	Host     string
	Password string
	Protocol string
	Username string
	TimeZone string
}

// Convert MySqlConnectionInfo to a data source name that can be used by the 'mysql' driver to open
// @see `sql.Open`
// @see FIXME driver documentation on their dataSourceName
// TODO document minimum fields required to create a valid data source name
func (info MySqlConnectionInfo) String() string {
	return fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=true&loc=%s",
		info.Username,
		info.Password,
		info.Protocol,
		info.Host,
		info.Database,
		info.TimeZone)
}

// TODO test String() output correct
// TODO test String() output can be used to open a mysql connection pool (integration/docker test) @see https://stackoverflow.com/a/41407042/5782298, http://peter.bourgon.org/go-in-production/#testing-and-validation
