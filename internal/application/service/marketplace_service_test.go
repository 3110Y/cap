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

var testPrimitives = []*entity.Primitive{
	{Type: "skills", Vendor: "anthropic", Name: "code-review", Version: "1.2.0", Description: "Автоматический анализ кода и поиск проблем."},
	{Type: "skills", Vendor: "anthropic", Name: "code-review", Version: "1.1.0", Description: "Анализ кода (предыдущая версия)"},
	{Type: "agents", Vendor: "community", Name: "pdf-parser", Version: "0.5.1", Description: "Извлечение текста из PDF"},
	{Type: "rules", Vendor: "internal", Name: "lint-rules", Version: "2.0.0", Description: "Правила линтера для Go"},
}

func newMarketplaceSvc(t *testing.T, primitives []*entity.Primitive) (*service.MarketplaceService, *gomock.Controller) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockCache.EXPECT().ReadIndexes().Return(primitives, nil).AnyTimes()
	return service.NewMarketplaceService(mockCache), ctrl
}

// --- List ---

func TestMarketplaceService_List_FiltersByType(t *testing.T) {
	svc, ctrl := newMarketplaceSvc(t, testPrimitives)
	defer ctrl.Finish()

	got, err := svc.List("skills")

	require.NoError(t, err)
	assert.Len(t, got, 2)
	for _, p := range got {
		assert.Equal(t, "skills", p.Type)
	}
}

func TestMarketplaceService_List_ReturnsEmptyForUnknownType(t *testing.T) {
	svc, ctrl := newMarketplaceSvc(t, testPrimitives)
	defer ctrl.Finish()

	got, err := svc.List("unknown")

	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestMarketplaceService_List_CacheError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCache := mocks.NewMockCacheRepository(ctrl)
	mockCache.EXPECT().ReadIndexes().Return(nil, errors.New("io error"))

	svc := service.NewMarketplaceService(mockCache)
	_, err := svc.List("skills")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read indexes")
}

// --- Search ---

func TestMarketplaceService_Search_MatchesByDescription(t *testing.T) {
	svc, ctrl := newMarketplaceSvc(t, testPrimitives)
	defer ctrl.Finish()

	got, err := svc.Search("skills", "анализ")

	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestMarketplaceService_Search_CaseInsensitive(t *testing.T) {
	svc, ctrl := newMarketplaceSvc(t, testPrimitives)
	defer ctrl.Finish()

	got, err := svc.Search("skills", "АНАЛИЗ")

	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestMarketplaceService_Search_NoMatches(t *testing.T) {
	svc, ctrl := newMarketplaceSvc(t, testPrimitives)
	defer ctrl.Finish()

	got, err := svc.Search("skills", "ничего_нет")

	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestMarketplaceService_Search_InvalidRegex(t *testing.T) {
	svc, ctrl := newMarketplaceSvc(t, testPrimitives)
	defer ctrl.Finish()

	_, err := svc.Search("skills", "[invalid(")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid regex")
}

// --- Info ---

func TestMarketplaceService_Info_BySimpleName(t *testing.T) {
	svc, ctrl := newMarketplaceSvc(t, testPrimitives)
	defer ctrl.Finish()

	got, err := svc.Info("skills", "code-review", "")

	require.NoError(t, err)
	assert.Len(t, got, 2) // both versions
	assert.Equal(t, "code-review", got[0].Name)
}

func TestMarketplaceService_Info_ByVendorName(t *testing.T) {
	svc, ctrl := newMarketplaceSvc(t, testPrimitives)
	defer ctrl.Finish()

	got, err := svc.Info("skills", "anthropic/code-review", "")

	require.NoError(t, err)
	assert.Len(t, got, 2)
	for _, p := range got {
		assert.Equal(t, "anthropic", p.Vendor)
	}
}

func TestMarketplaceService_Info_WithVersion(t *testing.T) {
	svc, ctrl := newMarketplaceSvc(t, testPrimitives)
	defer ctrl.Finish()

	got, err := svc.Info("skills", "code-review", "1.1.0")

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, "1.1.0", got[0].Version)
}

func TestMarketplaceService_Info_NotFound(t *testing.T) {
	svc, ctrl := newMarketplaceSvc(t, testPrimitives)
	defer ctrl.Finish()

	_, err := svc.Info("skills", "nonexistent", "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found in marketplace")
}

func TestMarketplaceService_Info_WrongVendor(t *testing.T) {
	svc, ctrl := newMarketplaceSvc(t, testPrimitives)
	defer ctrl.Finish()

	_, err := svc.Info("skills", "wrongvendor/code-review", "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found in marketplace")
}
