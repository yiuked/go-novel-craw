package storege

import "gorm.io/gorm"

func migrate() {
	err := DBClient.Transaction(func(tx *gorm.DB) error {
		if !tx.Migrator().HasTable(&BookList{}) {
			if err := tx.Migrator().AutoMigrate(&BookList{}); err != nil {
				return err
			}
		}
		if !tx.Migrator().HasTable(&BookDetail{}) {
			if err := tx.Migrator().AutoMigrate(&BookDetail{}); err != nil {
				return err
			}
		}
		if !tx.Migrator().HasTable(&BookChapter{}) {
			if err := tx.Migrator().AutoMigrate(&BookChapter{}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return
	}
}
