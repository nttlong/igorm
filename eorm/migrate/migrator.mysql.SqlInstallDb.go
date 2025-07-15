package migrate

func (m *migratorMySql) GetSqlInstallDb() ([]string, error) {
	return []string{}, nil // no thing to do for MySQL
}
