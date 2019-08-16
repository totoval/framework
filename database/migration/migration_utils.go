package migration

import (
	"fmt"

	"github.com/totoval/framework/console"
	"github.com/totoval/framework/helpers/m"
	"github.com/totoval/framework/model"
)

type MigrationUtils struct {
	Migration
	model.BaseModel
}

// Project initialize
func (mu *MigrationUtils) SetUp() {

	console.Println(console.CODE_WARNING, "initializing:migration table")

	mu.Migration.up(mu.DB())

	console.Println(console.CODE_SUCCESS, "initialized:migration table")
}

// all migration list
func (mu *MigrationUtils) migrationList() (migrationList []Migration) {
	mu.DB().Find(&migrationList)
	return
}

// calc need migrate list
func (mu *MigrationUtils) needMigrateList() (_migratorList []Migrator) {
	for _, migrator := range migratorList {
		found := false
		for _, migration := range mu.migrationList() {
			if migrator.Name(&migrator) == migration.Migration {
				found = true
				break
			}
		}

		if !found {
			_migratorList = append(_migratorList, migrator)
		}
	}
	return
}

func (mu *MigrationUtils) currentBatch() uint {
	migration := &Migration{}
	if !mu.DB().Order("batch desc").First(&migration).RecordNotFound() {
		return migration.Batch
	}
	return 0
}
func (mu *MigrationUtils) addMigration(migratorName string, batch uint) bool {
	migration := &Migration{Migration: migratorName, Batch: batch}
	if nil != mu.DB().Create(&migration).Error {
		return false
	}
	return true
}
func (mu *MigrationUtils) needRollbackMigrationList(batch uint) (migrationList []Migration) {
	mu.DB().Where("batch = ?", batch).Find(&migrationList)
	return
}
func (mu *MigrationUtils) delMigration(migration *Migration) bool {
	if nil != mu.DB().Delete(&migration).Error {
		return false
	}
	return true
}
func (mu *MigrationUtils) errorRollback() {
	if err := recover(); err != nil {
		if _err, ok := err.(error); ok {
			console.Println(console.CODE_WARNING, "error:"+_err.Error())
		} else {
			console.Println(console.CODE_WARNING, "error:"+fmt.Sprint(err))
		}
	}
}

func (mu *MigrationUtils) Migrate() {

	defer mu.errorRollback()

	m.Transaction(func(h *m.Helper) {
		mu.SetTX(h.DB())
		batch := mu.currentBatch() + 1

		for _, migrator := range mu.needMigrateList() {
			migrationName := migrator.Name(&migrator)

			console.Println(console.CODE_WARNING, "migrating:"+migrationName)

			if err := migrator.Up(mu.DB()).Error; err != nil {
				panic(err)
			}

			// add migration
			if !mu.addMigration(migrationName, batch) {
				panic("migration added failed!")
			}
			console.Println(console.CODE_SUCCESS, "migrated:"+migrationName)
		}
	}, 1)
}
func (mu *MigrationUtils) Rollback() {

	defer mu.errorRollback()

	m.Transaction(func(h *m.Helper) {
		mu.SetTX(h.DB())
		for _, migration := range mu.needRollbackMigrationList(mu.currentBatch()) {
			console.Println(console.CODE_WARNING, "rollbacking:"+migration.Name())

			migrator := newMigrator(migration.Name())
			if migrator == nil {
				panic("migration has not been defined yet!")
			}

			if err := migrator.Down(mu.DB()).Error; err != nil {
				panic(err)
			}

			if !mu.delMigration(&migration) {
				panic("migration deleted failed!")
			}

			console.Println(console.CODE_SUCCESS, "rollbacked:"+migration.Name())
		}
	}, 1)
}
func (mu *MigrationUtils) Fresh() {

}
func (mu *MigrationUtils) Install() {

	//   --database[=DATABASE]  The database connection to use
	//@todo
	//  Create the migration repository
}
func (mu *MigrationUtils) Refresh() {

}
func (mu *MigrationUtils) Reset() {

}
func (mu *MigrationUtils) Status() {

	//+------+--------------------------------------------------------------+-------+
	//| Ran? | Migration                                                    | Batch |
	//+------+--------------------------------------------------------------+-------+
	//| Yes  | 2014_10_12_000000_create_users_table                         | 3     |
	//| Yes  | 2014_10_12_100000_create_password_resets_table               | 1     |
	//| Yes  | 2016_06_01_000001_create_oauth_auth_codes_table              | 1     |
	//| Yes  | 2016_06_01_000002_create_oauth_access_tokens_table           | 1     |
	//| Yes  | 2016_06_01_000003_create_oauth_refresh_tokens_table          | 1     |
	//| Yes  | 2016_06_01_000004_create_oauth_clients_table                 | 1     |
	//| Yes  | 2016_06_01_000005_create_oauth_personal_access_clients_table | 1     |
	//| Yes  | 2019_01_10_081308_create_user_verification_table             | 2     |
	//| Yes  | 2019_01_10_165704_create_data_area_table                     | 2     |
	//| No   | 2019_01_22_112905_create_customer_wechat_table               |       |
	//| No   | 2019_01_22_112909_create_customer_table                      |       |
	//+------+--------------------------------------------------------------+-------+
}
