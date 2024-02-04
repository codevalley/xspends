package impl

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"testing"
	"time"
	"xspends/models/interfaces"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestInsertSource(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.SourceModel = NewSourceModel()
	})
	defer tearDown()

	source := &interfaces.Source{
		UserID:    1,
		ScopeID:   1,
		Name:      "Test Source",
		Type:      "CREDIT",
		Balance:   100.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := ModelsService.SourceModel.InsertSource(ctx, source)
	assert.NoError(t, err)
	//test for generic query error
	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), errors.New("executing insert for source")).
		Times(1)

	err = ModelsService.SourceModel.InsertSource(ctx, source)
	assert.Error(t, err)
	//test for invalid source type
	source = &interfaces.Source{
		UserID:    1,
		ScopeID:   1,
		Name:      "Test Source",
		Type:      "Invalid type",
		Balance:   100.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = ModelsService.SourceModel.InsertSource(ctx, source)
	assert.Error(t, err)

	//test for missing source name
	source = &interfaces.Source{
		UserID:    1,
		ScopeID:   1,
		Name:      "",
		Type:      "Invalid type",
		Balance:   100.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = ModelsService.SourceModel.InsertSource(ctx, source)
	assert.Error(t, err)

	//test for missing user ID
	source = &interfaces.Source{
		Name:      "Source name",
		Type:      "Invalid type",
		Balance:   100.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = ModelsService.SourceModel.InsertSource(ctx, source)
	assert.Error(t, err)
	//TODO: Missing scopeID test case

}

func TestUpdateSource(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.SourceModel = NewSourceModel()
	})
	defer tearDown()

	source := &interfaces.Source{
		ID:        1,
		UserID:    1,
		ScopeID:   1,
		Name:      "Updated Source",
		Type:      "SAVINGS",
		Balance:   200.0,
		UpdatedAt: time.Now(),
	}

	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := ModelsService.SourceModel.UpdateSource(ctx, source)
	assert.NoError(t, err)

	//test for generic query error
	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), errors.New("executing insert for source")).
		Times(1)

	err = ModelsService.SourceModel.UpdateSource(ctx, source)
	assert.Error(t, err)
	//test for invalid source type
	source = &interfaces.Source{
		UserID:    1,
		ScopeID:   1,
		Name:      "Test Source",
		Type:      "Invalid type",
		Balance:   100.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = ModelsService.SourceModel.UpdateSource(ctx, source)
	assert.Error(t, err)

	//test for missing source name
	source = &interfaces.Source{
		UserID:    1,
		ScopeID:   1,
		Name:      "",
		Type:      "Invalid type",
		Balance:   100.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = ModelsService.SourceModel.UpdateSource(ctx, source)
	assert.Error(t, err)

	//test for missing user ID
	source = &interfaces.Source{
		Name:      "Source name",
		Type:      "Invalid type",
		Balance:   100.0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = ModelsService.SourceModel.UpdateSource(ctx, source)
	assert.Error(t, err)
	//TODO: Missing scopeID test case
}

func TestDeleteSource(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.SourceModel = NewSourceModel()
	})
	defer tearDown()

	sourceID := int64(1)
	scopeID := []int64{1}

	mockExecutor.EXPECT().
		ExecContext(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(sql.Result(nil), nil).
		Times(1)

	err := ModelsService.SourceModel.DeleteSourceNew(ctx, sourceID, scopeID)
	assert.NoError(t, err)

	//test for generic query error
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	mockDBService := &DBService{Executor: db}
	mockModelService := &ModelsServiceContainer{
		DBService:   mockDBService,
		SourceModel: NewSourceModel(),
	}
	ModelsService = mockModelService
	// Set up the expected query with sqlmock to return an error
	mock.ExpectExec(`DELETE FROM sources WHERE id = \? AND scope_id IN (?)`).
		WithArgs(sourceID, scopeID).
		WillReturnError(errors.New("execution error"))

	ctx := context.Background()
	err = ModelsService.SourceModel.DeleteSourceNew(ctx, sourceID, scopeID)

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
	mockModelService := &ModelsServiceContainer{
		DBService:   mockDBService,
		SourceModel: NewSourceModel(),
	}
	ModelsService = mockModelService
	// Set up expectations
	rows := sqlmock.NewRows([]string{"source_id", "user_id", "name", "type", "balance", "scope_id", "created_at", "updated_at"}).
		AddRow(1, 1, "Test Source", "CREDIT", 100.0, 1, time.Now(), time.Now())
	mock.ExpectQuery("^SELECT (.+) FROM sources WHERE").WithArgs(1, 1).WillReturnRows(rows)

	// Call the function under test
	ctx := context.Background()
	source, err := ModelsService.SourceModel.GetSourceByIDNew(ctx, 1, []int64{1})

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, source)
	assert.Equal(t, int64(1), source.ID)
	assert.Equal(t, int64(1), source.UserID)
	assert.Equal(t, "Test Source", source.Name)
	assert.Equal(t, "CREDIT", source.Type)
	assert.Equal(t, 100.0, source.Balance)
	assert.Equal(t, int64(1), source.ScopeID)

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("source not found")
	}

	//test for query error with no rows found
	mock.ExpectQuery("^SELECT (.+) FROM sources WHERE").WithArgs(1, 1).
		WillReturnError(sql.ErrNoRows)

	exists, err1 := ModelsService.SourceModel.GetSourceByIDNew(ctx, 1, []int64{1})

	assert.Error(t, err1)
	assert.True(t, exists == nil) // Source does not exist

	//test for generic query error
	mock.ExpectQuery("^SELECT (.+) FROM sources WHERE").WithArgs(1, 1).
		WillReturnError(errors.New("query execution error"))

	_, err = ModelsService.SourceModel.GetSourceByIDNew(ctx, 1, []int64{1})

	assert.Error(t, err)

}

func TestGetSources(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.SourceModel = NewSourceModel()
	})
	defer tearDown()

	userID := int64(1)
	scopes := []int64{1}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockDBService := &DBService{Executor: db}
	mockModelService := &ModelsServiceContainer{
		DBService:   mockDBService,
		SourceModel: NewSourceModel(),
	}
	ModelsService = mockModelService

	mockRows := sqlmock.NewRows([]string{"source_id", "user_id", "name", "type", "balance", "scope_id", "created_at", "updated_at"}).
		AddRow(1, userID, "Source 1", "CREDIT", 100.00, scopes[0], time.Now(), time.Now()).
		AddRow(2, userID, "Source 2", "SAVINGS", 200.00, scopes[0], time.Now(), time.Now())

	// Set up the expected query with sqlmock
	mock.ExpectQuery(`SELECT source_id, user_id, name, type, balance, scope_id, created_at, updated_at FROM sources WHERE scope_id IN (?)`).
		WithArgs(scopes[0]).
		WillReturnRows(mockRows)

	ctx := context.Background()
	sources, err := ModelsService.SourceModel.GetSourcesNew(ctx, scopes)

	assert.NoError(t, err)
	assert.NotNil(t, sources)
	assert.Equal(t, 2, len(sources)) // Assuming 2 sources are returned

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	//test for generic query error
	mock.ExpectQuery(`SELECT source_id, user_id, name, type, balance, scope_id, created_at, updated_at FROM sources WHERE scope_id IN (?)`).
		WithArgs(scopes[0]).
		WillReturnError(errors.New("query execution error"))

	_, err = ModelsService.SourceModel.GetSourcesNew(ctx, scopes)

	assert.Error(t, err)

	//test row processing error
	rows := sqlmock.NewRows([]string{"source_id", "user_id", "name", "type", "balance", "scope_id", "created_at", "updated_at"}).
		AddRow(1, userID, "Source 1", "CREDIT", 100.00, scopes[0], time.Now(), time.Now()).
		AddRow(2, userID, "Source 2", "SAVINGS", 200.00, scopes[0], time.Now(), time.Now()).
		RowError(1, errors.New("row processing error"))

	mock.ExpectQuery(`SELECT source_id, user_id, name, type, balance, scope_id, created_at, updated_at FROM sources WHERE scope_id IN (?)`).
		WithArgs(scopes[0]).
		WillReturnRows(rows)

	_, err = ModelsService.SourceModel.GetSourcesNew(ctx, scopes)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "during row processing for sources: row processing error")

}

func TestSourceIDExists(t *testing.T) {
	tearDown := setUp(t, func(config *ModelsConfig) {
		// Replace the mocked CategoryModel with a real one just for this test
		config.SourceModel = NewSourceModel()
	})
	defer tearDown()

	// Variables for the test
	sourceID := int64(1)
	scopes := []int64{1, 2, 3} // Example with multiple scope IDs

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock service setup
	mockDBService := &DBService{Executor: db}
	mockModelService := &ModelsServiceContainer{
		DBService:   mockDBService,
		SourceModel: NewSourceModel(),
	}
	ModelsService = mockModelService

	ctx := context.Background()

	// Test case 1: Source exists
	{
		mockRows := sqlmock.NewRows([]string{"1"}).AddRow(1) // Assuming the source exists
		args := prepareArgsForSQLMock(scopes, sourceID)
		mock.ExpectQuery(`SELECT 1 FROM sources WHERE .+ AND source_id = \? LIMIT 1`).
			WithArgs(args...).
			WillReturnRows(mockRows)

		exists, err := ModelsService.SourceModel.SourceIDExistsNew(ctx, sourceID, scopes)
		assert.NoError(t, err)
		assert.True(t, exists)
	}

	// Test case 2: No rows found
	{
		args := prepareArgsForSQLMock(scopes, sourceID)
		mock.ExpectQuery(`SELECT 1 FROM sources WHERE .+ AND source_id = \? LIMIT 1`).
			WithArgs(args...).
			WillReturnError(sql.ErrNoRows)

		exists, err := ModelsService.SourceModel.SourceIDExistsNew(ctx, sourceID, scopes)
		assert.NoError(t, err)
		assert.False(t, exists)
	}

	// Test case 3: Generic query execution error
	{
		args := prepareArgsForSQLMock(scopes, sourceID)
		mock.ExpectQuery(`SELECT 1 FROM sources WHERE .+ AND source_id = \? LIMIT 1`).
			WithArgs(args...).
			WillReturnError(errors.New("query execution error"))

		_, err := ModelsService.SourceModel.SourceIDExistsNew(ctx, sourceID, scopes)
		assert.Error(t, err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// Helper function to prepare args for SQLMock with dynamic scope IDs and sourceID
func prepareArgsForSQLMock(scopes []int64, sourceID int64) []driver.Value {
	args := make([]driver.Value, 0, len(scopes)+1)
	for _, scope := range scopes {
		args = append(args, scope)
	}
	args = append(args, sourceID)
	return args
}
