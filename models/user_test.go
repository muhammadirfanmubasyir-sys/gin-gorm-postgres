package models

import "testing"

func TestUserTableName(t *testing.T) {
	user := User{}
	if got := user.TableName(); got != "table_of_user" {
		t.Errorf("TableName() = %q, want %q", got, "table_of_user")
	}
}
