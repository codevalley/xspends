package impl

import (
	"context"
	"database/sql"
	"testing"
	"time"
	"xspends/models/interfaces"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInsertSource(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	source := &interfaces.Source{
		UserID:    1,
		Name:      "Test Source",
		Type:      SourceTypeCredit,
		Balance:   100.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := mockModelService.SourceModel.InsertSource(ctx, source)
	assert.NoError(t, err)
}

func TestUpdateSource(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	source := &interfaces.Source{
		ID:        1,
		UserID:    1,
		Name:      "Updated Source",
		Type:      SourceTypeSavings,
		Balance:   200.0,
		UpdatedAt: time.Now(),
	}

	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := mockModelService.SourceModel.UpdateSource(ctx, source)
	assert.NoError(t, err)
}

func TestDeleteSource(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	sourceID := int64(1)
	userID := int64(1)

	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := mockModelService.SourceModel.DeleteSource(ctx, sourceID, userID)
	assert.NoError(t, err)
}

func TestGetSourceByID(t *testing.T) {
	// Create a new sqlmock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mockDBService := &DBService{Executor: db}
	mockModelService = &ModelsServiceContainer{
		DBService:   mockDBService,
		SourceModel: &SourceModel{},
	}
	ModelsService = mockModelService
	// Set up expectations
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "type", "balance", "created_at", "updated_at"}).
		AddRow(1, 1, "Test Source", SourceTypeCredit, 100.0, time.Now(), time.Now())
	mock.ExpectQuery("^SELECT (.+) FROM sources WHERE").WithArgs(1, 1).WillReturnRows(rows)

	// Call the function under test
	ctx := context.Background()
	source, err := mockModelService.SourceModel.GetSourceByID(ctx, 1, 1)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, source)
	assert.Equal(t, int64(1), source.ID)
	assert.Equal(t, int64(1), source.UserID)
	assert.Equal(t, "Test Source", source.Name)
	assert.Equal(t, SourceTypeCredit, source.Type)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Add similar tests for GetSources and SourceIDExists
