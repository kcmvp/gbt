package datasource

import (
	"fmt"
	env "github.com/kcmvp/gbt/env"
)

const (
	MySQL    = "mysql"
	SQLite   = "sqlite3"
	Postgres = "postgres"
)

const datasourceKey = "datasource"

type DataSource struct {
	Driver   string
	Host     string
	Port     string
	Username string
	Password string
	Db       string
	Url      string
}

func (ds *DataSource) DsName() (string, error) {
	if len(ds.Url) > 0 {
		return ds.Url, nil
	} else {
		switch ds.Driver {
		case MySQL:
			//<user>:<pass>@tcp(<host>:<port>)/<database>?parseTime=True
			return fmt.Sprintf("%s:%s@tcp%s:%s/%s?parseTime=True", ds.Username, ds.Password, ds.Host, ds.Port, ds.Db), nil
		case SQLite:
			return fmt.Sprintf("%s", ds.Driver), nil
		case Postgres:
			//host=<host> port=<port> user=<user> dbname=<database> password=<pass>
			return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", ds.Host, ds.Port, ds.Username, ds.Db, ds.Password), nil
		default:
			return "", fmt.Errorf("unsupported driver: %q", ds.Driver)
		}
	}
}

func Datasource() (*DataSource, error) {
	ds := &DataSource{}
	err := env.ActiveProfile().UnmarshalKey(datasourceKey, ds)
	return ds, err
}
