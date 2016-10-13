package migrator

import (
  "strconv"

    "github.com/grafana/grafana/pkg/log"
 // "strings"


  "fmt"
  "github.com/go-xorm/core"
  _ "github.com/denisenkom/go-mssqldb"
  //"github.com/jinzhu/gorm"
  //"github.com/jinzhu/gorm"
 // "database/sql"
)

type Mssql struct {
  BaseDialect
}

func NewMssqlDialect() *Mssql {
  logger := log.New("main")
  logger.Info("ms sql dialect111 loading")
  d := Mssql{}
  d.BaseDialect.dialect = &d
  d.BaseDialect.driverName = MSSQL
  return &d
}

func (db *Mssql) GetName() string {
  return "mssql"
}

func (db *Mssql) BindVar(i int) string {
  return "$$" // ?
}
func (db *Mssql) SupportEngine() bool {
  return false
}
func (db *Mssql) ForUpdateSql(query string) string {
  return query
}
func (db *Mssql) IndexOnTable() bool {
  return true
}
func (db *Mssql) Quote(name string) string {
  return "\"" + name + "\""
}
func (db *Mssql) SupportCharset() bool {
  return false
}
func (db *Mssql) QuoteStr() string {
  return "\""
}

func (db *Mssql) AutoIncrStr() string {
  return "IDENTITY"
}
func (db *Mssql) SupportInsertMany() bool {
  return true
}

func (db *Mssql) SqlType(c *Column) string {
  var res string
  switch c.Type {
  case DB_Bool:
    res = DB_TinyInt
    if c.Default == "true" {
      c.Default = "1"
    } else if c.Default == "false" {
      c.Default = "0"
    }
  case DB_Serial:
    c.IsAutoIncrement = true
    c.IsPrimaryKey = true
    c.Nullable = false
    res = DB_Int
  case DB_BigSerial:
    c.IsAutoIncrement = true
    c.IsPrimaryKey = true
    c.Nullable = false
    res = DB_BigInt
  case DB_Bytea, DB_Blob, DB_Binary, DB_TinyBlob, DB_MediumBlob, DB_LongBlob:
    res = DB_VarBinary
    if c.Length == 0 {
      c.Length = 50
    }
  case DB_TimeStamp:
    res = DB_DateTime
  case DB_TimeStampz:
    res = "DATETIMEOFFSET"
    c.Length = 7
  case DB_MediumInt:
    res = DB_Int
  case DB_MediumText,DB_TinyText, DB_LongText:
    res = DB_Text
  case DB_Double:
    res = DB_Real
  default:
    res = c.Type
  }

  var hasLen1 bool = (c.Length > 0)
  var hasLen2 bool = (c.Length2 > 0)

  if res == DB_BigInt && !hasLen1 && !hasLen2 {
    c.Length = 20
    hasLen1 = true
  }

  if hasLen2 {
    res += "(" + strconv.Itoa(c.Length) + "," + strconv.Itoa(c.Length2) + ")"
  } else if hasLen1 {
    res += "(" + strconv.Itoa(c.Length) + ")"
  }
  return res
}

func (db *Mssql) IndexCheckSql(tableName, idxName string) (string, []interface{}) {
  args := []interface{}{idxName}
  sql := "select name from sysindexes where id=object_id('" + tableName + "') and name=?"
  return sql, args
}
func (db *Mssql) TableCheckSql(tableName string) (string, []interface{}) {
  args := []interface{}{}
  sql := "SELECT * FROM sysobjects WHERE id = object_id(N'" + tableName + "') AND OBJECTPROPERTY(id, N'IsUserTable') = 1"
  return sql, args
}

func (db *Mssql) DropTableSql(tableName string) string {
  return fmt.Sprintf("IF EXISTS (SELECT * FROM sysobjects WHERE id = "+
  "object_id(N'%s') and OBJECTPROPERTY(id, N'IsUserTable') = 1) "+
  "DROP TABLE \"%s\"", tableName, tableName)
}
func (db *Mssql) Filters() []core.Filter {
  return []core.Filter{&core.IdFilter{}, &core.QuoteFilter{}}
}
