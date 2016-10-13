package migrations

import (
	. "github.com/grafana/grafana/pkg/services/sqlstore/migrator"
)

func addQuotaMigration(mg *Migrator) {

	var quotaV1 = Table{
		Name: "quota",
		Columns: []*Column{
			{Name: "id", Type: DB_Int, IsPrimaryKey: true, IsAutoIncrement: true},
			{Name: "org_id", Type: DB_Int, Nullable: true},
			{Name: "user_id", Type: DB_Int, Nullable: true},
			{Name: "target", Type: DB_NVarchar, Length: 255, Nullable: false},
			{Name: "limit", Type: DB_Int, Nullable: false},
			{Name: "created", Type: DB_DateTime, Nullable: false},
			{Name: "updated", Type: DB_DateTime, Nullable: false},
		},
		Indices: []*Index{
			{Cols: []string{"org_id", "user_id", "target"}, Type: UniqueIndex},
		},
	}
	mg.AddMigration("create quota table v1", NewAddTableMigration(quotaV1))

	//-------  indexes ------------------
	addTableIndicesMigrations(mg, "v1", quotaV1)
}
