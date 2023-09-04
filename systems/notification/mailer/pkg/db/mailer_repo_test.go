package db_test

import (
	"log"
	"testing"
	"time"

	"github.com/tj/assert"
	int_db "github.com/ukama/ukama/systems/notification/mailer/pkg/db"

	uuid "github.com/ukama/ukama/systems/common/uuid"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UkamaDbMock struct {
	GormDb *gorm.DB
}

func (u UkamaDbMock) Init(model ...interface{}) error {
	panic("implement me: Init()")
}

func (u UkamaDbMock) Connect() error {
	panic("implement me: Connect()")
}

func (u UkamaDbMock) GetGormDb() *gorm.DB {
	return u.GormDb
}

func (u UkamaDbMock) InitDB() error {
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func() error) error {
	log.Fatal("implement me: ExecuteInTransaction()")
	return nil
}

func (u UkamaDbMock) ExecuteInTransaction2(dbOperation func(tx *gorm.DB) *gorm.DB,
	nestedFuncs ...func(tx *gorm.DB) error) error {
	log.Fatal("implement me: ExecuteInTransaction2()")
	return nil
}

func Test_SendEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	gdb, _ := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	repo := int_db.NewMailerRepo(&UkamaDbMock{
		GormDb: gdb,
	})

	email := int_db.Mailing{
		MailId:       uuid.NewV4(),
		Email:        "brackley@ukama.com",
		TemplateName: "test_template",
		SentAt:       nil,
		Status:       "pending",
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		DeletedAt:             nil,
	}

	mock.ExpectBegin()
			mock.ExpectExec("INSERT INTO \"mailings\"").WithArgs(
				email.MailId,
				email.Email,
				email.TemplateName,
				email.SentAt,
				email.Status,
				email.CreatedAt,
				email.UpdatedAt,
				email.DeletedAt).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	err = repo.SendEmail(&email)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

}

// func Test_GetEmailById(t *testing.T) {
//     db, mock, err := sqlmock.New()
//     if err != nil {
//         t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
//     }
//     defer db.Close()
//     gdb, _ := gorm.Open(postgres.New(postgres.Config{
//         DSN:                  "sqlmock_db_0",
//         DriverName:           "postgres",
//         Conn:                 db,
//         PreferSimpleProtocol: true,
//     }), &gorm.Config{})
//     repo := int_db.NewMailerRepo(&UkamaDbMock{
//         GormDb: gdb,
//     })

//     email := int_db.Mailing{
//         MailId:       uuid.NewV4(),
//         Email:        " brackley@ukama.com",
//         TemplateName: "test_template",
//         SentAt:       nil,
//         Status:       "pending",
//         CreatedAt:    time.Now(),
//         UpdatedAt:    time.Now(),
//         DeletedAt:    nil,
//     }

//     rows := sqlmock.NewRows([]string{"mail_id", "email", "template_name", "sent_at", "status", "created_at", "updated_at", "deleted_at"}).
//         AddRow(email.MailId, email.Email, email.TemplateName, email.SentAt, email.Status, email.CreatedAt, email.UpdatedAt, email.DeletedAt)

//     mock.ExpectQuery("SELECT * FROM \"mailings\"").WithArgs(email.MailId).WillReturnRows(rows)
//     result, err := repo.GetEmailById(email.MailId)
//     assert.NotNil(t, result)
//     assert.NoError(t, mock.ExpectationsWereMet())
// }
func Test_GetEmailById(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    gdb, _ := gorm.Open(postgres.New(postgres.Config{
        DSN:                  "sqlmock_db_0",
        DriverName:           "postgres",
        Conn:                 db,
        PreferSimpleProtocol: true,
    }), &gorm.Config{})

    repo := int_db.NewMailerRepo(&UkamaDbMock{
        GormDb: gdb,
    })

    email := int_db.Mailing{
        MailId:       uuid.NewV4(),
        Email:        "brackley@ukama.com",
        TemplateName: "test_template",
        SentAt:       nil,
        Status:       "pending",
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
        DeletedAt:    nil,
    }

    rows := sqlmock.NewRows([]string{"mail_id", "email", "template_name", "sent_at", "status", "created_at", "updated_at", "deleted_at"}).
        AddRow(email.MailId, email.Email, email.TemplateName, email.SentAt, email.Status, email.CreatedAt, email.UpdatedAt, email.DeletedAt)

    // Adjust the expectation to match the actual query with the WHERE condition.
    // Use `\` to escape special characters in the regular expression.
    mock.ExpectQuery(`SELECT \* FROM "mailings" WHERE mail_id = \$1`).WithArgs(email.MailId).WillReturnRows(rows)

    // Call the GetEmailById function with the UUID argument.
    result, err := repo.GetEmailById(email.MailId)

    // Perform your assertions on the result and error, if any.
    assert.NotNil(t, result)
    assert.NoError(t, err)

    // Ensure that all expectations were met.
    assert.NoError(t, mock.ExpectationsWereMet())
}