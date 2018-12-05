package db

import (
	"fmt"
	// blank import for MSSQL driver
	_ "github.com/denisenkom/go-mssqldb"
)

type msSQLDialect struct {
	BaseDialect
}

const (
	insertMigrationMSSQLDialectSql = "insert into %v.%v (name, source_dir, filename, type, db_schema) values (@p1, @p2, @p3, @p4, @p5)"
	insertTenantMSSQLDialectSql    = "insert into %v.%v (name) values (@p1)"
)

func (md *msSQLDialect) GetMigrationInsertSql() string {
	return fmt.Sprintf(insertMigrationMSSQLDialectSql, migratorSchema, migratorMigrationsTable)
}

func (md *msSQLDialect) GetTenantInsertSql() string {
	return fmt.Sprintf(insertTenantMSSQLDialectSql, migratorSchema, migratorTenantsTable)
}

func (md *msSQLDialect) GetCreateTenantsTableSql() string {
	createTenantsTableSql := `
IF NOT EXISTS (select * from information_schema.tables where table_schema = '%v' and table_name = '%v')
BEGIN
  create table [%v].%v (
    id int identity (1,1) primary key,
    name varchar(200) not null,
    created datetime default CURRENT_TIMESTAMP
  );
END
`
	return fmt.Sprintf(createTenantsTableSql, migratorSchema, migratorTenantsTable, migratorSchema, migratorTenantsTable)
}

func (md *msSQLDialect) GetCreateMigrationsTableSql() string {
	createMigrationsTableSql := `
IF NOT EXISTS (select * from information_schema.tables where table_schema = '%v' and table_name = '%v')
BEGIN
  create table [%v].%v (
    id int identity (1,1) primary key,
    name varchar(200) not null,
    source_dir varchar(200) not null,
    filename varchar(200) not null,
    type int not null,
    db_schema varchar(200) not null,
    created datetime default CURRENT_TIMESTAMP
  );
END
`
	return fmt.Sprintf(createMigrationsTableSql, migratorSchema, migratorMigrationsTable, migratorSchema, migratorMigrationsTable)
}

func (md *msSQLDialect) GetCreateSchemaSql(schema string) string {
	createSchemaSql := `
IF NOT EXISTS (select * from information_schema.schemata where schema_name = '%v')
BEGIN
	EXEC sp_executesql N'create schema %v';
END
`
	return fmt.Sprintf(createSchemaSql, schema, schema)
}
