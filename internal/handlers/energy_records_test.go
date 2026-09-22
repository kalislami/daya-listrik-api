package handlers

import (
	"context"
	"daya-listrik-api/internal/models"
	"daya-listrik-api/internal/repository"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockRepo struct{ mock.Mock }

func (m *mockRepo) AddRecord(ctx context.Context, r *models.EnergyRecord) error {
	return m.Called(ctx, r).Error(0)
}
func (m *mockRepo) GetRecords(ctx context.Context) ([]models.EnergyRecord, error) {
	a := m.Called(ctx)
	return a.Get(0).([]models.EnergyRecord), a.Error(1)
}
func (m *mockRepo) GetByIdRecord(ctx context.Context, id int) (*models.EnergyRecord, error) {
	a := m.Called(ctx, id)
	if a.Get(0) == nil {
		return nil, a.Error(1)
	}
	return a.Get(0).(*models.EnergyRecord), a.Error(1)
}
func (m *mockRepo) UpdateRecord(ctx context.Context, r *models.EnergyRecord) error {
	return m.Called(ctx, r).Error(0)
}
func (m *mockRepo) DeleteRecord(ctx context.Context, id int) error {
	return m.Called(ctx, id).Error(0)
}

func testApp() (*fiber.App, *mockRepo) {
	app := fiber.New()
	repo := new(mockRepo)
	InitializeRoutes(app, &EnergyRecordHandler{Repo: repo})
	return app, repo
}

func call(t *testing.T, app *fiber.App, method, path, body string) *http.Response {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := app.Test(req)
	require.NoError(t, err)
	t.Cleanup(func() { resp.Body.Close() })
	return resp
}

func checkError(t *testing.T, resp *http.Response, status int, code string) {
	t.Helper()
	require.Equal(t, status, resp.StatusCode)
	var payload apiError
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&payload))
	require.Equal(t, code, payload.Error.Code)
	require.NotEmpty(t, payload.Error.Message)
	require.NotContains(t, payload.Error.Message, "database secret detail")
}

func TestCreateRecord(t *testing.T) {
	valid := `{"device":"  AC  ","usage":100,"duration":2}`
	t.Run("success", func(t *testing.T) {
		app, repo := testApp()
		repo.On("AddRecord", mock.Anything, mock.MatchedBy(func(r *models.EnergyRecord) bool {
			return r.Device == "AC" && r.Usage == 100 && r.Duration == 2 && r.ID == 0
		})).Run(func(a mock.Arguments) { a.Get(1).(*models.EnergyRecord).ID = 7 }).Return(nil).Once()
		resp := call(t, app, "POST", "/api/records", valid)
		require.Equal(t, 201, resp.StatusCode)
		var record models.EnergyRecord
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&record))
		require.Equal(t, 7, record.ID)
		require.Equal(t, "AC", record.Device)
		repo.AssertExpectations(t)
	})
	for name, body := range map[string]string{
		"invalid JSON":     `{`,
		"blank device":     `{"device":"   ","usage":100,"duration":2}`,
		"zero usage":       `{"device":"AC","usage":0,"duration":2}`,
		"negative usage":   `{"device":"AC","usage":-1,"duration":2}`,
		"invalid duration": `{"device":"AC","usage":100,"duration":0}`,
		"server field":     `{"id":1,"device":"AC","usage":100,"duration":2}`,
	} {
		t.Run(name, func(t *testing.T) {
			app, _ := testApp()
			checkError(t, call(t, app, "POST", "/api/records", body), 400, "VALIDATION_ERROR")
		})
	}
	t.Run("repository error", func(t *testing.T) {
		app, repo := testApp()
		repo.On("AddRecord", mock.Anything, mock.Anything).Return(errors.New("database secret detail")).Once()
		resp := call(t, app, "POST", "/api/records", valid)
		checkError(t, resp, 500, "INTERNAL_SERVER_ERROR")
	})
}

func TestGetRecord(t *testing.T) {
	for _, id := range []string{"abc", "0", "-1"} {
		t.Run("invalid "+id, func(t *testing.T) {
			app, _ := testApp()
			checkError(t, call(t, app, "GET", "/api/records/"+id, ""), 400, "INVALID_ID")
		})
	}
	for name, repoErr := range map[string]error{"not found": repository.ErrRecordNotFound, "database error": errors.New("database secret detail")} {
		t.Run(name, func(t *testing.T) {
			app, repo := testApp()
			repo.On("GetByIdRecord", mock.Anything, 1).Return(nil, repoErr).Once()
			status, code := 404, "RECORD_NOT_FOUND"
			if name == "database error" {
				status, code = 500, "INTERNAL_SERVER_ERROR"
			}
			checkError(t, call(t, app, "GET", "/api/records/1", ""), status, code)
		})
	}
	app, repo := testApp()
	repo.On("GetByIdRecord", mock.Anything, 1).Return(&models.EnergyRecord{ID: 1, Device: "AC"}, nil).Once()
	resp := call(t, app, "GET", "/api/records/1", "")
	require.Equal(t, 200, resp.StatusCode)
}

func TestUpdateRecord(t *testing.T) {
	valid := `{"device":"AC","usage":100,"duration":2}`
	app, repo := testApp()
	repo.On("UpdateRecord", mock.Anything, mock.MatchedBy(func(r *models.EnergyRecord) bool { return r.ID == 1 && r.Duration == 2 })).Return(nil).Once()
	require.Equal(t, 200, call(t, app, "PUT", "/api/records/1", valid).StatusCode)
	for _, tc := range []struct{ path, body, code string }{
		{"/api/records/0", valid, "INVALID_ID"},
		{"/api/records/1", `{`, "VALIDATION_ERROR"},
		{"/api/records/1", `{"device":"AC","usage":1,"duration":-2}`, "VALIDATION_ERROR"},
	} {
		app, _ := testApp()
		checkError(t, call(t, app, "PUT", tc.path, tc.body), 400, tc.code)
	}
	for _, tc := range []struct {
		err    error
		status int
		code   string
	}{
		{repository.ErrRecordNotFound, 404, "RECORD_NOT_FOUND"},
		{errors.New("db failed"), 500, "INTERNAL_SERVER_ERROR"},
	} {
		app, repo := testApp()
		repo.On("UpdateRecord", mock.Anything, mock.Anything).Return(tc.err).Once()
		checkError(t, call(t, app, "PUT", "/api/records/1", valid), tc.status, tc.code)
	}
}

func TestDeleteRecord(t *testing.T) {
	app, repo := testApp()
	repo.On("DeleteRecord", mock.Anything, 1).Return(nil).Once()
	require.Equal(t, 204, call(t, app, "DELETE", "/api/records/1", "").StatusCode)
	app, _ = testApp()
	checkError(t, call(t, app, "DELETE", "/api/records/-1", ""), 400, "INVALID_ID")
	for _, tc := range []struct {
		err    error
		status int
		code   string
	}{
		{repository.ErrRecordNotFound, 404, "RECORD_NOT_FOUND"},
		{errors.New("db failed"), 500, "INTERNAL_SERVER_ERROR"},
	} {
		app, repo := testApp()
		repo.On("DeleteRecord", mock.Anything, 1).Return(tc.err).Once()
		checkError(t, call(t, app, "DELETE", "/api/records/1", ""), tc.status, tc.code)
	}
}

func TestListRecords(t *testing.T) {
	for _, records := range [][]models.EnergyRecord{{{ID: 1}}, {}} {
		app, repo := testApp()
		repo.On("GetRecords", mock.Anything).Return(records, nil).Once()
		resp := call(t, app, "GET", "/api/records", "")
		require.Equal(t, 200, resp.StatusCode)
		var got []models.EnergyRecord
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&got))
		require.Len(t, got, len(records))
	}
	app, repo := testApp()
	repo.On("GetRecords", mock.Anything).Return([]models.EnergyRecord(nil), errors.New("db failed")).Once()
	checkError(t, call(t, app, "GET", "/api/records", ""), 500, "INTERNAL_SERVER_ERROR")
}
