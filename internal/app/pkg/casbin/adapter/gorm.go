package adapter

import (
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// casbinRuleTableName casbin rule's table name
// for consistency with the ent adapter
const casbinRuleTableName = "casbin_rules"

// NewGormAdapter build casin gorm adapter
func NewGormAdapter(db *gorm.DB) (adp *gormadapter.Adapter, err error) {
	// turn on automatic migration for consistency with the ent adapter
	// gormadapter.TurnOffAutoMigrate(db)

	adp, err = gormadapter.NewAdapterByDBUseTableName(db, "", casbinRuleTableName)
	if err != nil {
		return nil, err
	}

	return
}
