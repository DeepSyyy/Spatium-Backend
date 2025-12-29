package repositories_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/repositories"
	test_helpers "github.com/DeepSyyy/Spatium-Backend/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	// Implement test for Create method
	db, mock, cleanup := test_helpers.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(db)

	user := &models.User{}

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users"`).
		WillReturnRows(sqlmock.NewRows([]string{"internal_id"}).AddRow(1))
	mock.ExpectCommit()
	err := repo.Create(user)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.InternalID)
}

func TestFindByRecoveryCode(t *testing.T) {
	db, mock, cleanup := test_helpers.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(db)

	rows := sqlmock.NewRows([]string{"internal_id", "recovery_code"}).
		AddRow(1, "recovery123")

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE recovery_code = \$1 ORDER BY "users"."internal_id" LIMIT \$2`).
		WithArgs("recovery123", 1).
		WillReturnRows(rows)

	user, err := repo.FindByRecoveryCode("recovery123")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.InternalID)
}

func TestUpdateLastLogin(t *testing.T) {
	db, mock, cleanup := test_helpers.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewUserRepository(db)

	user := &models.User{InternalID: 1}

	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "users" SET "last_login"=\$1 WHERE internal_id = \$2`).
		WithArgs(sqlmock.AnyArg(), int64(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateLastLogin(user)
	assert.NoError(t, err)
}
