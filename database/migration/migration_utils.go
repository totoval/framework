package migration

import (
	"github.com/jinzhu/gorm"
	"github.com/totoval/framework/cmd"
	"github.com/totoval/framework/model"
)

type MigrationUtils struct {
	db    *gorm.DB
	chLog chan interface{}
	Migration
}

func (m *MigrationUtils) Init(chLog chan interface{}) {
	m.setLog(chLog)
	m.setDB()
}

func (m *MigrationUtils) setDB() {
	m.db = model.DB()
}
func (m *MigrationUtils) setLog(ch chan interface{}) {
	m.chLog = ch
}

// 项目初始化
func (m *MigrationUtils) SetUp() {
	defer m.closeLog()
	m.log(cmd.CODE_WARNING, "initializing:migration table")

	m.Migration.up(m.db)

	m.log(cmd.CODE_SUCCESS, "initialized:migration table")
}

// 所有migrate过的任务列表
func (m *MigrationUtils) migrationList() (migrationList []Migration) {
	m.db.Find(&migrationList)
	return
}

// 计算需要migrate的任务
func (m *MigrationUtils) needMigrateList() (_migratorList []Migrator) {
	for _, migrator := range migratorList {
		found := false
		for _, migration := range m.migrationList() {
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

func (m *MigrationUtils) currentBatch() uint {
	migration := &Migration{}
	if !m.db.Order("batch desc").First(&migration).RecordNotFound() {
		return migration.Batch
	}
	return 0
}
func (m *MigrationUtils) addMigration(migratorName string, batch uint) bool {
	migration := &Migration{Migration: migratorName, Batch: batch}
	if nil != m.db.Create(&migration).Error {
		return false
	}
	return true
}
func (m *MigrationUtils) needRollbackMigrationList(batch uint) (migrationList []Migration) {
	m.db.Where("batch = ?", batch).Find(&migrationList)
	return
}
func (m *MigrationUtils) delMigration(migration *Migration) bool {
	if nil != m.db.Delete(&migration).Error {
		return false
	}
	return true
}
func (m *MigrationUtils) errorRollback(tx *gorm.DB) {
	if err := recover(); err != nil {
		tx.Rollback()
		if _err, ok := err.(error); ok {
			m.log(cmd.CODE_WARNING, "error:"+_err.Error())
		}else{
			m.log(cmd.CODE_WARNING, "error:"+err.(string))  //@todo err.(string) may be down when `panic(123)`
		}
	}
}

func (m *MigrationUtils) log(code interface{}, message string) {
	if _code, ok := code.(cmd.Attribute); ok {
		m.chLog <- cmd.TermLog{
			Code:    _code,
			Message: message,
		}
	}
}
func (m *MigrationUtils) closeLog() {
	m.chLog <- nil
}

func (m *MigrationUtils) Migrate() {
	defer m.closeLog()
	tx := m.db.Begin()
	{
		defer m.errorRollback(tx)

		batch := m.currentBatch() + 1
		for _, migrator := range m.needMigrateList() {
			migrationName := migrator.Name(&migrator)

			m.log(cmd.CODE_WARNING, "migrating:"+migrationName)

			if err := migrator.Up(m.db).Error; err != nil {
				panic(err)
			}

			// add migration
			if !m.addMigration(migrationName, batch) {
				panic("migration added failed!")
			}
			m.log(cmd.CODE_SUCCESS, "migrated:"+migrationName)
		}
	}
	tx.Commit()
}
func (m *MigrationUtils) Rollback() {
	defer m.closeLog()
	tx := m.db.Begin()
	{
		defer m.errorRollback(tx)

		for _, migration := range m.needRollbackMigrationList(m.currentBatch()) {
			m.log(cmd.CODE_WARNING, "rollbacking:"+migration.Name())

			migrator := newMigrator(migration.Name())
			if migrator == nil {
				panic("migration has not been defined yet!")
			}

			if err := migrator.Down(m.db).Error; err != nil {
				panic(err)
			}

			if !m.delMigration(&migration) {
				panic("migration deleted failed!")
			}

			m.log(cmd.CODE_SUCCESS, "rollbacked:"+migration.Name())
		}
	}
	tx.Commit()
}
func (m *MigrationUtils) Fresh() {
	defer m.closeLog()

}
func (m *MigrationUtils) Install() {
	defer m.closeLog()
	//   --database[=DATABASE]  The database connection to use
	//@todo
	//  Create the migration repository
}
func (m *MigrationUtils) Refresh() {
	defer m.closeLog()

}
func (m *MigrationUtils) Reset() {
	defer m.closeLog()

}
func (m *MigrationUtils) Status() {
	defer m.closeLog()
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
