package test_helpers

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupMockDB membuat instance mock database menggunakan sqlmock
// dan mengembalikan *gorm.DB, sqlmock, dan fungsi cleanup
func SetupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	// 1️⃣ Buat koneksi mock
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	// 2️⃣ Bungkus mock ke dalam GORM dialector (pakai postgres dialector)
	dialector := postgres.New(postgres.Config{
		Conn: db,
	})

	// 3️⃣ Inisialisasi GORM tanpa log (biar output test rapi)
	gdb, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	require.NoError(t, err)

	// 4️⃣ Cleanup function (pastikan semua ekspektasi terpenuhi)
	cleanup := func() {
		err := mock.ExpectationsWereMet()
		require.NoError(t, err)
	}

	return gdb, mock, cleanup
}
