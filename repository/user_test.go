package repository_test

import (
	"testing"

	"GoCleanArchitecture/entities"
	repo "GoCleanArchitecture/repository"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	gormDB, _ := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})

	mockUser := &entities.User{
		ID: 9,
	}

	rows := sqlmock.NewRows([]string{"id"}).
		AddRow(mockUser.ID)

	query := "^SELECT \\* FROM `users` WHERE id = \\?"
	mock.ExpectQuery(query).WithArgs("9").WillReturnRows(rows)

	d := repo.NewUserRepository(gormDB)
	user, err := d.GetUser("9")

	assert.NoError(t, err)
	assert.Equal(t, mockUser, &user)
}
