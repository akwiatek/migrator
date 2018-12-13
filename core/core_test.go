package core

import (
	"testing"

	"github.com/lukaszbudnik/migrator/config"
	"github.com/stretchr/testify/assert"
)

const (
	unknownAction = "unknown"
	configFile    = "../test/migrator.yaml"
)

func TestPrintConfig(t *testing.T) {
	config, err := config.FromFile(configFile)
	assert.Nil(t, err)
	executeFlags := ExecuteFlags{}
	executeFlags.Action = PrintConfigAction
	doExecuteMigrator(config, executeFlags, createMockedConnector, createMockedDiskLoader)
}

func TestGetDiskMigrations(t *testing.T) {
	config, err := config.FromFile(configFile)
	assert.Nil(t, err)
	executeFlags := ExecuteFlags{}
	executeFlags.Action = GetDiskMigrationsAction
	doExecuteMigrator(config, executeFlags, createMockedConnector, createMockedDiskLoader)
}

func TestGetDBTenants(t *testing.T) {
	config, err := config.FromFile(configFile)
	assert.Nil(t, err)
	executeFlags := ExecuteFlags{}
	executeFlags.Action = GetDBTenantsAction
	doExecuteMigrator(config, executeFlags, createMockedConnector, createMockedDiskLoader)
}

func TestGetDBMigrations(t *testing.T) {
	config, err := config.FromFile(configFile)
	assert.Nil(t, err)
	executeFlags := ExecuteFlags{}
	executeFlags.Action = GetDBMigrationsAction
	doExecuteMigrator(config, executeFlags, createMockedConnector, createMockedDiskLoader)
}

func TestApplyMigrations(t *testing.T) {
	config, err := config.FromFile(configFile)
	assert.Nil(t, err)
	executeFlags := ExecuteFlags{}
	executeFlags.Action = ApplyAction
	doExecuteMigrator(config, executeFlags, createMockedConnector, createMockedDiskLoader)
}

func TestApplyMigrationsVerificationFailed(t *testing.T) {
	config, err := config.FromFile(configFile)
	assert.Nil(t, err)
	executeFlags := ExecuteFlags{}
	executeFlags.Action = ApplyAction
	doExecuteMigrator(config, executeFlags, createMockedConnector, createBrokenCheckSumMockedDiskLoader)
}

func TestAddTenant(t *testing.T) {
	config, err := config.FromFile(configFile)
	assert.Nil(t, err)
	executeFlags := ExecuteFlags{}
	executeFlags.Action = AddTenantAction
	doExecuteMigrator(config, executeFlags, createMockedConnector, createMockedDiskLoader)
}

func TestAddTenantVerificationFailed(t *testing.T) {
	config, err := config.FromFile(configFile)
	assert.Nil(t, err)
	executeFlags := ExecuteFlags{}
	executeFlags.Action = AddTenantAction
	doExecuteMigrator(config, executeFlags, createMockedConnector, createBrokenCheckSumMockedDiskLoader)
}

// todo 81%

func TestGetDBTenantsError(t *testing.T) {
	config, err := config.FromFile(configFile)
	assert.Nil(t, err)
	executeFlags := ExecuteFlags{}
	executeFlags.Action = GetDBTenantsAction
	err = doExecuteMigrator(config, executeFlags, createMockedErrorConnector, createMockedDiskLoader)
	assert.Equal(t, "trouble maker", err.Error())
}

func TestGetDBMigrationsError(t *testing.T) {
	config, err := config.FromFile(configFile)
	assert.Nil(t, err)
	executeFlags := ExecuteFlags{}
	executeFlags.Action = GetDBMigrationsAction
	err = doExecuteMigrator(config, executeFlags, createMockedErrorConnector, createMockedDiskLoader)
	assert.Equal(t, "trouble maker", err.Error())
}

func TestApplyMigrationsError(t *testing.T) {
	config, err := config.FromFile(configFile)
	assert.Nil(t, err)
	executeFlags := ExecuteFlags{}
	executeFlags.Action = ApplyAction
	err = doExecuteMigrator(config, executeFlags, createMockedErrorConnector, createMockedDiskLoader)
	assert.Equal(t, "trouble maker", err.Error())
}

func TestApplyMigrationsPassingVerificationError(t *testing.T) {
	config, err := config.FromFile(configFile)
	assert.Nil(t, err)
	executeFlags := ExecuteFlags{}
	executeFlags.Action = ApplyAction
	err = doExecuteMigrator(config, executeFlags, createMockedPassingVerificationErrorConnector, createMockedDiskLoader)
	// assert.Equal(t, "trouble maker", err.Error())
}

func TestAddTenantError(t *testing.T) {
	config, err := config.FromFile(configFile)
	assert.Nil(t, err)
	executeFlags := ExecuteFlags{}
	executeFlags.Action = AddTenantAction
	err = doExecuteMigrator(config, executeFlags, createMockedErrorConnector, createMockedDiskLoader)
	assert.Equal(t, "trouble maker", err.Error())
}

func TestAddTenantPassingVerificationError(t *testing.T) {
	config, err := config.FromFile(configFile)
	assert.Nil(t, err)
	executeFlags := ExecuteFlags{}
	executeFlags.Action = AddTenantAction
	err = doExecuteMigrator(config, executeFlags, createMockedPassingVerificationErrorConnector, createMockedDiskLoader)
	assert.Equal(t, "trouble maker", err.Error())
}
