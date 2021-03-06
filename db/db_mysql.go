package db

import (
	"fmt"
	// blank import for MySQL driver
	_ "github.com/go-sql-driver/mysql"
)

type mySQLDialect struct {
	baseDialect
}

const (
	insertMigrationMySQLDialectSQL             = "insert into %v.%v (name, source_dir, filename, type, db_schema, contents, checksum, version_id) values (?, ?, ?, ?, ?, ?, ?, ?)"
	insertTenantMySQLDialectSQL                = "insert into %v.%v (name) values (?)"
	insertVersionMySQLDialectSQL               = "insert into %v.%v (name) values (?)"
	selectVersionsByFileMySQLDialectSQL        = "select mv.id as vid, mv.name as vname, mv.created as vcreated, mm.id as mid, mm.name, mm.source_dir, mm.filename, mm.type, mm.db_schema, mm.created, mm.contents, mm.checksum from %v.%v mv left join %v.%v mm on mv.id = mm.version_id where mv.id in (select version_id from %v.%v where filename = ?) order by vid desc, mid asc"
	selectVersionByIDMySQLDialectSQL           = "select mv.id as vid, mv.name as vname, mv.created as vcreated, mm.id as mid, mm.name, mm.source_dir, mm.filename, mm.type, mm.db_schema, mm.created, mm.contents, mm.checksum from %v.%v mv left join %v.%v mm on mv.id = mm.version_id where mv.id = ? order by mid asc"
	selectMigrationByIDMySQLDialectSQL         = "select id, name, source_dir, filename, type, db_schema, created, contents, checksum from %v.%v where id = ?"
	versionsTableSetupMySQLDropDialectSQL      = `drop procedure if exists migrator_create_versions`
	versionsTableSetupMySQLCallDialectSQL      = `call migrator_create_versions()`
	versionsTableSetupMySQLProcedureDialectSQL = `
create procedure migrator_create_versions()
begin
if not exists (select * from information_schema.tables where table_schema = '%v' and table_name = '%v') then
  create table %v.%v (
    id serial primary key,
    name varchar(200) not null,
    created timestamp default now()
  );
  alter table %v.%v add column version_id bigint unsigned;
  create index migrator_versions_version_id_idx on %v.%v (version_id);
  if exists (select * from %v.%v) then
    insert into %v.%v (name) values ('Initial version');
    -- initial version_id sequence is always 1
    update %v.%v set version_id = 1;
  end if;
  alter table %v.%v
    modify version_id bigint unsigned not null;
  alter table %v.%v
    add constraint migrator_versions_version_id_fk foreign key (version_id) references %v.%v (id) on delete cascade;
end if;
end;
`
)

// LastInsertIDSupported instructs migrator if Result.LastInsertId() is supported by the DB driver
func (md *mySQLDialect) LastInsertIDSupported() bool {
	return true
}

// GetMigrationInsertSQL returns MySQL-specific migration insert SQL statement
func (md *mySQLDialect) GetMigrationInsertSQL() string {
	return fmt.Sprintf(insertMigrationMySQLDialectSQL, migratorSchema, migratorMigrationsTable)
}

// GetTenantInsertSQL returns MySQL-specific migrator's default tenant insert SQL statement
func (md *mySQLDialect) GetTenantInsertSQL() string {
	return fmt.Sprintf(insertTenantMySQLDialectSQL, migratorSchema, migratorTenantsTable)
}

func (md *mySQLDialect) GetVersionInsertSQL() string {
	return fmt.Sprintf(insertVersionMySQLDialectSQL, migratorSchema, migratorVersionsTable)
}

// GetCreateVersionsTableSQL returns MySQL-specific SQLs which does:
// 1. drop procedure if exists
// 2. create procedure
// 3. calls procedure
// far from ideal MySQL in contrast to MS SQL and PostgreSQL does not support the execution of anonymous blocks of code
func (md *mySQLDialect) GetCreateVersionsTableSQL() []string {
	return []string{
		versionsTableSetupMySQLDropDialectSQL,
		fmt.Sprintf(versionsTableSetupMySQLProcedureDialectSQL, migratorSchema, migratorVersionsTable, migratorSchema, migratorVersionsTable, migratorSchema, migratorMigrationsTable, migratorSchema, migratorMigrationsTable, migratorSchema, migratorMigrationsTable, migratorSchema, migratorVersionsTable, migratorSchema, migratorMigrationsTable, migratorSchema, migratorMigrationsTable, migratorSchema, migratorMigrationsTable, migratorSchema, migratorVersionsTable),
		versionsTableSetupMySQLCallDialectSQL,
	}
}

func (md *mySQLDialect) GetVersionsByFileSQL() string {
	return fmt.Sprintf(selectVersionsByFileMySQLDialectSQL, migratorSchema, migratorVersionsTable, migratorSchema, migratorMigrationsTable, migratorSchema, migratorMigrationsTable)
}

func (md *mySQLDialect) GetVersionByIDSQL() string {
	return fmt.Sprintf(selectVersionByIDMySQLDialectSQL, migratorSchema, migratorVersionsTable, migratorSchema, migratorMigrationsTable)
}

func (md *mySQLDialect) GetMigrationByIDSQL() string {
	return fmt.Sprintf(selectMigrationByIDMySQLDialectSQL, migratorSchema, migratorMigrationsTable)
}
