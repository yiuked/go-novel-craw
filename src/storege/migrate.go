package storege

func migrate() {
	err := DB.AutoMigrate(
		&BookList{},
		&BookDetail{},
		&BookChapter{})
	if err != nil {
		return
	}
}
