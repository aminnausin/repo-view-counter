package db

func NewDatabase(driver string) Database {
	switch driver {
	case "sqlite":
		return NewSqliteDB()
	case "postgres":
		return NewPostgresDB()
	default:
		// fmt.Errorf("unsupported driver: %s", driver)
		return nil
	}
}
