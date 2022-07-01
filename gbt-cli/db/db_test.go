package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQuery(t *testing.T) {
	type args struct {
		artifact string
	}
	tests := []struct {
		name string
		args args
		want *Template
	}{
		{
			name: "datasource(mysql)",
			args: args{"github.com/go-sql-driver/mysql"},
			want: &Template{
				"datasource",
				"github.com/go-sql-driver/mysql;github.com/lib/pq;github.com/mattn/go-sqlite3",
				"",
			},
		},
		{
			name: "ent",
			args: args{"entgo.io/ent"},
			want: &Template{
				"ent",
				"entgo.io/ent",
				"datasource",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Query(tt.args.artifact)
			assert.NoErrorf(t, err, tt.name, tt.args.artifact)
			assert.Equalf(t, tt.want, got, "Query(%v)", tt.args.artifact)
		})
	}
}

func TestName(t *testing.T) {
	assert.Equal(t, 1, 2)
}
