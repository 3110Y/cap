package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/3110Y/cap/internal/application/service"
	"github.com/3110Y/cap/internal/domain/entity"
	"github.com/3110Y/cap/internal/integration/mocks"
)

func TestRepositoryService_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepositoryRepository(ctrl)
	mockRepo.EXPECT().
		Add(gomock.Any()).
		DoAndReturn(func(r *entity.Repository) error {
			assert.Equal(t, "https://github.com/example/repo", r.URL)
			assert.Len(t, r.ID, 8)
			return nil
		})

	svc := service.NewRepositoryService(mockRepo)
	got, err := svc.Add("https://github.com/example/repo")

	require.NoError(t, err)
	assert.Equal(t, "https://github.com/example/repo", got.URL)
	assert.Len(t, got.ID, 8)
}

func TestRepositoryService_Add_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepositoryRepository(ctrl)
	mockRepo.EXPECT().Add(gomock.Any()).Return(errors.New("disk full"))

	svc := service.NewRepositoryService(mockRepo)
	_, err := svc.Add("https://example.com/repo")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to add repository")
}

func TestRepositoryService_Add_GeneratesUniqueIDs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Two calls should produce different IDs (due to timestamp salt).
	var ids []string
	mockRepo := mocks.NewMockRepositoryRepository(ctrl)
	mockRepo.EXPECT().Add(gomock.Any()).Times(2).
		DoAndReturn(func(r *entity.Repository) error {
			ids = append(ids, r.ID)
			return nil
		})

	svc := service.NewRepositoryService(mockRepo)
	_, err := svc.Add("https://example.com/repo")
	require.NoError(t, err)
	_, err = svc.Add("https://example.com/repo")
	require.NoError(t, err)

	// IDs may collide within the same nanosecond in theory, but the format must be 8 chars.
	assert.Len(t, ids[0], 8)
	assert.Len(t, ids[1], 8)
}

func TestRepositoryService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	want := []*entity.Repository{
		{ID: "aaaa1111", URL: "https://a.com"},
		{ID: "bbbb2222", URL: "https://b.com"},
	}
	mockRepo := mocks.NewMockRepositoryRepository(ctrl)
	mockRepo.EXPECT().List().Return(want, nil)

	svc := service.NewRepositoryService(mockRepo)
	got, err := svc.List()

	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestRepositoryService_List_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepositoryRepository(ctrl)
	mockRepo.EXPECT().List().Return(nil, errors.New("file not found"))

	svc := service.NewRepositoryService(mockRepo)
	_, err := svc.List()

	require.Error(t, err)
}

func TestRepositoryService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepositoryRepository(ctrl)
	mockRepo.EXPECT().Delete("aaaa1111").Return(nil)

	svc := service.NewRepositoryService(mockRepo)
	err := svc.Delete("aaaa1111")

	require.NoError(t, err)
}

func TestRepositoryService_Delete_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepositoryRepository(ctrl)
	mockRepo.EXPECT().Delete("missing").Return(errors.New("not found"))

	svc := service.NewRepositoryService(mockRepo)
	err := svc.Delete("missing")

	require.Error(t, err)
}
