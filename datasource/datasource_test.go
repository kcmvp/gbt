package datasource

import (
	env "github.com/kcmvp/gbt/env"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatasource(t *testing.T) {
	assert.True(t, true, "True is true!")
	ds, _ := Datasource()
	dsName, _ := ds.DsName()
	assert.Equal(t, "sqlite3", ds.Driver, "Driver should be sqlite3")
	assert.Equal(t, "file:ent?mode=memory&cache=shared&_fk=1", ds.Url, "Driver should be sqlite3")
	assert.NotEmpty(t, ds.Url, "Url should not be emtpy")
	assert.Equal(t, "def", ds.Username, "Username should be def")
	assert.Equal(t, "file:ent?mode=memory&cache=shared&_fk=1", dsName)
	assert.Equal(t, "test", env.ActiveProfile().Name())
}
