package migrations

import . "github.com/grafana/grafana/pkg/services/sqlstore/migrator"

func addPlaylistMigrations(mg *Migrator) {
	mg.AddMigration("Drop old table playlist table", NewDropTableMigration("playlist"))
	mg.AddMigration("Drop old table playlist_item table", NewDropTableMigration("playlist_item"))

	playlistV2 := Table{
		Name: "playlist",
		Columns: []*Column{
			{Name: "id", Type: DB_Int, IsPrimaryKey: true, IsAutoIncrement: true},
			{Name: "name", Type: DB_NVarchar, Length: 255, Nullable: false},
			{Name: "interval", Type: DB_NVarchar, Length: 255, Nullable: false},
			{Name: "org_id", Type: DB_Int, Nullable: false},
		},
	}

	// create table
	mg.AddMigration("create playlist table v2", NewAddTableMigration(playlistV2))

	playlistItemV2 := Table{
		Name: "playlist_item",
		Columns: []*Column{
			{Name: "id", Type: DB_Int, IsPrimaryKey: true, IsAutoIncrement: true},
			{Name: "playlist_id", Type: DB_Int, Nullable: false},
			{Name: "type", Type: DB_NVarchar, Length: 255, Nullable: false},
			{Name: "value", Type: DB_Text, Nullable: false},
			{Name: "title", Type: DB_Text, Nullable: false},
			{Name: "order", Type: DB_Int, Nullable: false},
		},
	}

	mg.AddMigration("create playlist item table v2", NewAddTableMigration(playlistItemV2))
}
