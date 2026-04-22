package service_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/3110Y/cap/internal/application/service"
	"github.com/3110Y/cap/internal/domain/entity"
	"github.com/3110Y/cap/internal/integration/mocks"
)

func newInstallSvc(
	t *testing.T,
	mockCache *mocks.MockCacheRepository,
	mockProject *mocks.MockProjectRepository,
) *service.InstallService {
	t.Helper()
	mp := service.NewMarketplaceService(mockCache)
	return service.NewInstallService(mockCache, mockProject, mp)
}

func TestInstallService_Add_HappyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockProject := mocks.NewMockProjectRepository(ctrl)

	primitives := []*entity.Primitive{
		{
			Type:         "skills",
			Vendor:       "anthropic",
			Name:         "code-review",
			Version:      "1.2.0",
			Description:  "Code review skill",
			SourceRepoID: "abc12345",
		},
	}

	repoDir := "/cache/repos/abc12345"
	expectedSourceDir := filepath.Join(repoDir, "skills", "anthropic", "code-review", "1.2.0")

	mockCache.EXPECT().ReadIndexes().Return(primitives, nil)
	mockCache.EXPECT().RepoDir("abc12345").Return(repoDir)
	mockProject.EXPECT().Install(expectedSourceDir, "skills", "anthropic", "code-review", "1.2.0").Return(nil)

	svc := newInstallSvc(t, mockCache, mockProject)
	got, err := svc.Add("skills", "code-review", "1.2.0")

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "code-review", got.Name)
	assert.Equal(t, "anthropic", got.Vendor)
	assert.Equal(t, "1.2.0", got.Version)
}

func TestInstallService_Add_PrimitiveNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockProject := mocks.NewMockProjectRepository(ctrl)

	mockCache.EXPECT().ReadIndexes().Return([]*entity.Primitive{}, nil)

	svc := newInstallSvc(t, mockCache, mockProject)
	_, err := svc.Add("skills", "nonexistent", "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestInstallService_Add_EmptySourceRepoID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockProject := mocks.NewMockProjectRepository(ctrl)

	primitives := []*entity.Primitive{
		{
			Type:         "skills",
			Vendor:       "anthropic",
			Name:         "code-review",
			Version:      "1.2.0",
			SourceRepoID: "", // empty — signals cache not synced
		},
	}

	mockCache.EXPECT().ReadIndexes().Return(primitives, nil)

	svc := newInstallSvc(t, mockCache, mockProject)
	_, err := svc.Add("skills", "code-review", "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "cap update")
}

func TestInstallService_Add_InstallError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockProject := mocks.NewMockProjectRepository(ctrl)

	primitives := []*entity.Primitive{
		{
			Type:         "skills",
			Vendor:       "anthropic",
			Name:         "code-review",
			Version:      "1.2.0",
			SourceRepoID: "abc12345",
		},
	}

	repoDir := "/cache/repos/abc12345"
	expectedSourceDir := filepath.Join(repoDir, "skills", "anthropic", "code-review", "1.2.0")

	mockCache.EXPECT().ReadIndexes().Return(primitives, nil)
	mockCache.EXPECT().RepoDir("abc12345").Return(repoDir)
	mockProject.EXPECT().Install(expectedSourceDir, "skills", "anthropic", "code-review", "1.2.0").
		Return(errors.New("permission denied"))

	svc := newInstallSvc(t, mockCache, mockProject)
	_, err := svc.Add("skills", "code-review", "1.2.0")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "permission denied")
}
