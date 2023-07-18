package data

type PostgresTestRepository struct {
}

func NewPostgresTestRepository() *PostgresTestRepository {
	return &PostgresTestRepository{}
}

func (u *PostgresTestRepository) GetAll() ([]*User, error) {
	return []*User{}, nil
}

func (u *PostgresTestRepository) GetByEmail(email string) (*User, error) {
	return &User{}, nil
}

func (u *PostgresTestRepository) GetOne(id int) (*User, error) {
	return &User{}, nil
}

func (u *PostgresTestRepository) Update(user User) error {
	return nil
}

func (u *PostgresTestRepository) DeleteByID(id int) error {
	return nil
}

func (u *PostgresTestRepository) Insert(user User) (int, error) {
	return 1, nil
}

func (u *PostgresTestRepository) ResetPassword(password string, user User) error {
	return nil
}

func (u *PostgresTestRepository) PasswordMatches(plainText string, user User) (bool, error) {
	return true, nil
}
