package migrations

import  ."github.com/grafana/grafana/pkg/services/sqlstore/migrator"

func addMachineMigrations(mg *Migrator)  {

  machineV1:= Table{
    Name: "machine",
    Columns: []*Column{
      {Name: "machine_id", Type: DB_Int, IsPrimaryKey: true, IsAutoIncrement: true},
      {Name: "machine_name", Type: DB_NVarchar,Length: 45, Nullable: false},
      {Name: "process_id", Type: DB_Int, Nullable: false},
      {Name: "org_id", Type: DB_Int, Nullable: false},
      {Name: "description", Type: DB_NVarchar, Length: 80, Nullable: false},
      {Name: "created", Type: DB_DateTime},
      {Name: "updated", Type: DB_DateTime},
      {Name: "updated_by", Type: DB_NVarchar, Length: 45, Nullable: true},
      {Name: "vendor", Type: DB_NVarchar, Length: 45, Nullable: true},
    },
    Indices: []*Index{
      {Cols: []string{"updated_by"}, Type: IndexType},
      {Cols: []string{"vendor"}, Type: IndexType},
    },

  }
  mg.AddMigration("create process  table v1-7", NewAddTableMigration(machineV1))
  addTableIndicesMigrations(mg, "v1-7", machineV1)



}
