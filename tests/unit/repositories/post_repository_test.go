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

func TestCreatePost(t *testing.T) {
	db, mock, cleanup := test_helpers.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewPostRepository(db)

	post := &models.Post{Content: "Hello world", MoodTagID: 1}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "posts"`)).
		WillReturnRows(sqlmock.NewRows([]string{"internal_id"}).AddRow(1))
	mock.ExpectCommit()

	err := repo.Create(post)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), post.InternalID)
}

func TestGetAllPosts(t *testing.T) {
	db, mock, cleanup := test_helpers.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewPostRepository(db)

	rows := sqlmock.NewRows([]string{"internal_id", "content", "mood_tag_id"}).
		AddRow(1, "Hello world", 1).
		AddRow(2, "Another post", 2)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "posts"`)).
		WillReturnRows(rows)

	posts, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, posts, 2)
	assert.Equal(t, "Hello world", posts[0].Content)
}

func TestGetPostDetail(t *testing.T) {
	db, mock, cleanup := test_helpers.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewPostRepository(db)

	rows := sqlmock.NewRows([]string{"internal_id", "public_id", "content", "mood_tag_internal_id"}).
		AddRow(1, "123e4567-e89b-12d3-a456-426614174000", "Hello world", 1)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "posts" WHERE public_id = $1 ORDER BY "posts"."internal_id" LIMIT $2`)).
		WithArgs("123e4567-e89b-12d3-a456-426614174000", 1).
		WillReturnRows(rows)

	post, err := repo.GetPostDetail("123e4567-e89b-12d3-a456-426614174000")
	assert.NoError(t, err)
	assert.Equal(t, "Hello world", post.Content)
}

func TestGetPostsByUserID(t *testing.T) {
	db, mock, cleanup := test_helpers.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewPostRepository(db)

	rows := sqlmock.NewRows([]string{"internal_id", "content", "mood_tag_internala_id"}).
		AddRow(1, "Hello world", 1).
		AddRow(2, "Another post", 2)

	mock.ExpectQuery(
		regexp.QuoteMeta(`SELECT * FROM "posts" WHERE user_internal_id = $1 ORDER BY created_at desc`),
	).
		WithArgs(int64(42)).
		WillReturnRows(rows)

	posts, err := repo.GetPostsByUserID(42)
	assert.NoError(t, err)
	assert.Len(t, posts, 2)
}

func TestGetPostInternalIDByPublicID(t *testing.T) {
	db, mock, cleanup := test_helpers.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewPostRepository(db)

	rows := sqlmock.NewRows([]string{"internal_id"}).
		AddRow(1)

	mock.ExpectQuery(`(?i)SELECT\s+"?internal_id"?\s+FROM\s+"posts"\s+WHERE\s+public_id\s*=\s*\$1.*LIMIT\s*\$2`).
		WithArgs("123e4567-e89b-12d3-a456-426614174000", 1).
		WillReturnRows(rows)

	internalID, err := repo.GetPostInternalIDByPublicID("123e4567-e89b-12d3-a456-426614174000")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), internalID)
}

func TestUpdatePost(t *testing.T) {
	db, mock, cleanup := test_helpers.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewPostRepository(db)

	updatedPost := &models.Post{Content: "Updated content", MoodTagID: 2}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "posts" SET "content"=$1,"mood_tag_internal_id"=$2 WHERE public_id = $3`)).
		WithArgs(updatedPost.Content, updatedPost.MoodTagID, "123e4567-e89b-12d3-a456-426614174000").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Update("123e4567-e89b-12d3-a456-426614174000", updatedPost)
	assert.NoError(t, err)
}

func TestDeletePost(t *testing.T) {
	db, mock, cleanup := test_helpers.SetupMockDB(t)
	defer cleanup()

	repo := repositories.NewPostRepository(db)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "posts" WHERE public_id = $1`)).
		WithArgs("123e4567-e89b-12d3-a456-426614174000").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.Delete("123e4567-e89b-12d3-a456-426614174000")
	assert.NoError(t, err)
}
