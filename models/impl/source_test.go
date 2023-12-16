package impl

import (
	"context"
	"database/sql"
	"testing"
	"time"
	"xspends/models/interfaces"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
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
	//test for generic query error
	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), errors.New("executing insert for source")).
		Times(1)

	err = mockModelService.SourceModel.InsertSource(ctx, source)
	assert.Error(t, err)
	//test for invalid source type
	source = &interfaces.Source{
		UserID:    1,
		Name:      "Test Source",
		Type:      "Invalid type",
		Balance:   100.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = mockModelService.SourceModel.InsertSource(ctx, source)
	assert.Error(t, err)

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

	//test for generic query error
	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), errors.New("executing insert for source")).
		Times(1)

	err = mockModelService.SourceModel.UpdateSource(ctx, source)
	assert.Error(t, err)
	//test for invalid source type
	source = &interfaces.Source{
		UserID:    1,
		Name:      "Test Source",
		Type:      "Invalid type",
		Balance:   100.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = mockModelService.SourceModel.UpdateSource(ctx, source)
	assert.Error(t, err)
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

	//test for generic query error
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
	// Set up the expected query with sqlmock to return an error
	mock.ExpectExec(`DELETE FROM sources WHERE id = \? AND user_id = \?`).
		WithArgs(sourceID, userID).
		WillReturnError(errors.New("execution error"))

	ctx := context.Background()
	err = mockModelService.SourceModel.DeleteSource(ctx, sourceID, userID)

	assert.Error(t, err) // Expecting an error

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
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
		t.Errorf("source not found")
	}

	//test for query error with no rows found
	mock.ExpectQuery("^SELECT (.+) FROM sources WHERE").WithArgs(1, 1).
		WillReturnError(sql.ErrNoRows)

	exists, err1 := mockModelService.SourceModel.GetSourceByID(ctx, 1, 1)

	assert.Error(t, err1)
	assert.True(t, exists == nil) // Source does not exist

	//test for generic query error
	mock.ExpectQuery("^SELECT (.+) FROM sources WHERE").WithArgs(1, 1).
		WillReturnError(errors.New("query execution error"))

	_, err = mockModelService.SourceModel.GetSourceByID(ctx, 1, 1)

	assert.Error(t, err)

}

func TestGetSources(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	userID := int64(1)

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

	mockRows := sqlmock.NewRows([]string{"id", "user_id", "name", "type", "balance", "created_at", "updated_at"}).
		AddRow(1, userID, "Source 1", "CREDIT", 100.00, time.Now(), time.Now()).
		AddRow(2, userID, "Source 2", "SAVINGS", 200.00, time.Now(), time.Now())

	// Set up the expected query with sqlmock
	mock.ExpectQuery(`SELECT id, user_id, name, type, balance, created_at, updated_at FROM sources WHERE user_id = \?`).
		WithArgs(userID).
		WillReturnRows(mockRows)

	ctx := context.Background()
	sources, err := mockModelService.SourceModel.GetSources(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, sources)
	assert.Equal(t, 2, len(sources)) // Assuming 2 sources are returned

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	//test for generic query error
	mock.ExpectQuery(`SELECT id, user_id, name, type, balance, created_at, updated_at FROM sources WHERE user_id = \?`).
		WithArgs(userID).
		WillReturnError(errors.New("query execution error"))

	_, err = mockModelService.SourceModel.GetSources(ctx, userID)

	assert.Error(t, err)

	//test row processing error
	rows := sqlmock.NewRows([]string{"id", "user_id", "name", "type", "balance", "created_at", "updated_at"}).
		AddRow(1, userID, "Source 1", "CREDIT", 100.00, time.Now(), time.Now()).
		AddRow(2, userID, "Source 2", "SAVINGS", 200.00, time.Now(), time.Now()).
		RowError(1, errors.New("row processing error"))

	mock.ExpectQuery(`SELECT id, user_id, name, type, balance, created_at, updated_at FROM sources WHERE user_id = \?`).
		WithArgs(userID).
		WillReturnRows(rows)

	_, err = mockModelService.SourceModel.GetSources(ctx, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "during row processing for sources: row processing error")

}

func TestSourceIDExists(t *testing.T) {
	tearDown := setUp(t)
	defer tearDown()

	sourceID := int64(1)
	userID := int64(1)

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

	// Prepare a row with a single column that represents 'exists'
	mockRows := sqlmock.NewRows([]string{"1"}).AddRow(1) // Assuming the source exists

	// Set up the expected query with sqlmock
	mock.ExpectQuery(`SELECT 1 FROM sources WHERE id = \? AND user_id = \? LIMIT 1`).
		WithArgs(sourceID, userID).
		WillReturnRows(mockRows)

	ctx := context.Background()
	exists, err := mockModelService.SourceModel.SourceIDExists(ctx, sourceID, userID)

	assert.NoError(t, err)
	assert.True(t, exists) // Assuming the source exists

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
	//test for query error with no rows found
	mock.ExpectQuery(`SELECT 1 FROM sources WHERE id = \? AND user_id = \? LIMIT 1`).
		WithArgs(sourceID, userID).
		WillReturnError(sql.ErrNoRows)

	exists, err = mockModelService.SourceModel.SourceIDExists(ctx, sourceID, userID)

	assert.NoError(t, err)
	assert.False(t, exists) // Source does not exist

	//test for generic query error
	mock.ExpectQuery(`SELECT 1 FROM sources WHERE id = \? AND user_id = \? LIMIT 1`).
		WithArgs(sourceID, userID).
		WillReturnError(errors.New("query execution error"))

	_, err = mockModelService.SourceModel.SourceIDExists(ctx, sourceID, userID)

	assert.Error(t, err)
}
