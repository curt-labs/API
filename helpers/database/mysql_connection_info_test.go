package database

import (
	"testing"
)

func TestMySqlConnectionInfo_String(t *testing.T) {
	testInfo := MySqlConnectionInfo{Database: "testDatabase",
		Host: "testHost",
		Password: "testPassword",
		Protocol: "testProtocol",
		Username: "testUsername",
		TimeZone: "testTimezone"}
	expected := "testUsername:testPassword@testProtocol(testHost)/testDatabase?parseTime=true&loc=testTimezone"
	actual := testInfo.String()

	if actual != expected {
		t.Errorf("MySqlConnectionInfo.String is not formatted correctly\nactual:   %s\nexpected: %s", actual, expected)
	}
}

func TestIntegrationMySqlConnectionInfo_String_CanConnectToServer(t *testing.T) {
	if testing.Short() {
		t.Skip("TestIntegrationMySqlConnectionInfo_String_CanConnectToServer")
	}

	// TODO setup 'mysql' container using dockertest
	// TODO see if the generated connection string can connect to the container and run a Ping()
}
