package migrations

import  ."github.com/grafana/grafana/pkg/services/sqlstore/migrator"

func addProcessMigrations(mg *Migrator) {
/* for all databases expect mssql and azure sql
  processV1 := Table{
    Name: "process",
    Columns: []*Column{
      {Name: "process_id", Type: DB_Int, Length: 20, IsPrimaryKey: true, IsAutoIncrement: true, Nullable: true},
      {Name: "process_name", Type: DB_NVarchar, Length: 255, Nullable: true},
      {Name: "org_id", Type: DB_Int, Length: 20, Nullable: true},
      {Name: "created", Type: DB_DateTime, Nullable: false},
      {Name: "updated", Type: DB_DateTime, Nullable: false},
      {Name: "updated_by", Type: DB_NVarchar, Length: 45, Nullable: true},
    },
    Indices: []*Index{
      {Cols: []string{"process_id"}, Type: IndexType},
      {Cols: []string{"org_id"}, Type: IndexType},
      {Cols: []string{"updated_by"}, Type: IndexType},
    },

  }
  */
  processV1 := Table{
    Name: "process",
    Columns: []*Column{
      {Name: "process_id", Type: DB_Int,  IsPrimaryKey: true, IsAutoIncrement: true},
      {Name: "process_name", Type: DB_NVarchar, Length: 255, Nullable: true},
      {Name: "org_id", Type: DB_Int,  Nullable: true},
      {Name: "created", Type: DB_DateTime, Nullable: false},
      {Name: "updated", Type: DB_DateTime, Nullable: false},
      {Name: "updated_by", Type: DB_NVarchar, Length: 45, Nullable: true},
    },
    Indices: []*Index{
      {Cols: []string{"process_id"}, Type: IndexType},
      {Cols: []string{"org_id"}, Type: IndexType},
      {Cols: []string{"updated_by"}, Type: IndexType},
    },

  }
  mg.AddMigration("create process  table v1-7", NewAddTableMigration(processV1))
  addTableIndicesMigrations(mg, "v1-7", processV1)

}
