package repositories_test

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/DeepSyyy/Spatium-Backend/models"
	"github.com/DeepSyyy/Spatium-Backend/repositories"
	test_helpers "github.com/DeepSyyy/Spatium-Backend/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func TestFindByUserAndPost(t *testing.T) {
	db, mock, cleanup := test_helpers.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewReactionRepository(db)

	// Pastikan gunakan sqlmock.NewRows, bukan mock.NewRows
	rows := sqlmock.NewRows([]string{"internal_id", "user_internal_id", "post_internal_id", "type"}).
		AddRow(1, 1, 1, int(0))

	// Gunakan regexp.QuoteMeta agar pattern SQL bisa diinterpretasi literal
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "reactions" WHERE user_internal_id = $1 AND post_internal_id = $2 ORDER BY "reactions"."internal_id" LIMIT $3`)).
		WithArgs(int64(1), int64(1), 1).
		WillReturnRows(rows)

	reaction, err := repo.FindByUserAndPost(1, 1)

	assert.NoError(t, err)
	assert.NotNil(t, reaction)
	assert.Equal(t, int64(1), reaction.InternalID)
	assert.Equal(t, int64(1), reaction.UserID)
	assert.Equal(t, int64(1), reaction.PostID)
	assert.Equal(t, int64(0), reaction.ReactionTypeInternalID)

	// Pastikan semua expectation sqlmock terpenuhi
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestCreateReaction(t *testing.T) {
	db, mock, cleanup := test_helpers.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewReactionRepository(db)

	reaction := &models.Reaction{
		UserID:                 1,
		PostID:                 1,
		ReactionTypeInternalID: 1,
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "reactions"`)).
		WillReturnRows(sqlmock.NewRows([]string{"internal_id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(reaction)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), reaction.InternalID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateReaction(t *testing.T) {
	db, mock, cleanup := test_helpers.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewReactionRepository(db)

	reaction := &models.Reaction{
		InternalID:             1,
		UserID:                 1,
		PostID:                 1,
		ReactionTypeInternalID: 2,
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "reactions" SET "public_id"=$1,"user_internal_id"=$2,"post_internal_id"=$3,"reaction_type_internal_id"=$4,"emoji"=$5,"created_at"=$6 WHERE "internal_id" = $7`)).
		WithArgs(
			sqlmock.AnyArg(), // public_id (UUID)
			int64(1),         // user_internal_id
			int64(1),         // post_internal_id
			int64(2),         // reaction_type_internal_id
			"",               // emoji
			sqlmock.AnyArg(), // created_at
			int64(1),         // internal_id (WHERE)
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update(reaction)

	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetReactionTypeIDByEmoji(t *testing.T) {
	db, mock, cleanup := test_helpers.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewReactionRepository(db)

	rows := sqlmock.NewRows([]string{"internal_id"}).
		AddRow(1)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "internal_id" FROM "reaction_types" WHERE emoji = $1 ORDER BY "reaction_types"."internal_id" LIMIT $2`)).
		WithArgs("👍", 1).
		WillReturnRows(rows)
	mock.ExpectCommit()

	reactionTypeID, err := repo.GetReactionTypeIDByEmoji("👍")

	assert.NoError(t, err)
	assert.Equal(t, int64(1), reactionTypeID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetReactionSummary(t *testing.T) {
	db, mock, cleanup := test_helpers.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewReactionRepository(db)

	rows := sqlmock.NewRows([]string{"emoji", "count"}).
		AddRow("👍", int64(5)).
		AddRow("❤️", int64(3))

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT rt.emoji, COUNT(r.internal_id) AS count FROM reaction_types rt JOIN reactions r ON rt.internal_id = r.reaction_type_internal_id WHERE r.post_internal_id = $1 GROUP BY "rt"."emoji"`,
	)).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	summary, err := repo.GetReactionSummary(1)
	assert.NoError(t, err)

	expected := map[string]int{
		"👍":  5,
		"❤️": 3,
	}
	assert.Equal(t, expected, summary)
	assert.NoError(t, mock.ExpectationsWereMet())
}
