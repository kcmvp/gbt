package profile

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatasource(t *testing.T) {
	assert.True(t, true, "True is true!")
	ds, err := GetDatasource()
	assert.Nil(t, err)
	assert.Equal(t, "mysql", ds.Driver, "Should read application.yml as default")
	assert.Equal(t, "abc", ds.Username, "Username should be abc")
	dsName, err := ds.DsName()
	assert.Equal(t, "abc:@tcp:/?parseTime=True", dsName)
	assert.Empty(t, ds.Host, "host should be emtpy")
	assert.Empty(t, ds.Url, "Url should be emtpy")
	With("test")
	ds, err = GetDatasource()
	dsName, err = ds.DsName()
	assert.Equal(t, "sqlite3", ds.Driver, "Driver should be sqlite3")
	assert.Equal(t, "file:ent?mode=memory&cache=shared&_fk=1", ds.Url, "Driver should be sqlite3")
	assert.NotEmpty(t, ds.Url, "Url should not be emtpy")
	assert.Equal(t, "def", ds.Username, "Username should be def")
	assert.Equal(t, "file:ent?mode=memory&cache=shared&_fk=1", dsName)

}
