package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/3110Y/cap/internal/application/service"
	"github.com/3110Y/cap/internal/integration/mocks"
)

func TestSelfUpdateService_AlreadyLatest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSelfUpdateRepository(ctrl)
	mockRepo.EXPECT().LatestVersion().Return("v1.2.3", nil)
	// Download must NOT be called.

	svc := service.NewSelfUpdateService(mockRepo, service.CurrentVersion("v1.2.3"))
	result, err := svc.Update()

	require.NoError(t, err)
	assert.True(t, result.AlreadyLatest)
	assert.Equal(t, "1.2.3", result.OldVersion)
	assert.Equal(t, "1.2.3", result.NewVersion)
}

func TestSelfUpdateService_AlreadyLatest_WithoutVPrefix(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSelfUpdateRepository(ctrl)
	mockRepo.EXPECT().LatestVersion().Return("v2.0.0", nil)

	// CurrentVersion without "v" prefix should also match.
	svc := service.NewSelfUpdateService(mockRepo, service.CurrentVersion("2.0.0"))
	result, err := svc.Update()

	require.NoError(t, err)
	assert.True(t, result.AlreadyLatest)
}

func TestSelfUpdateService_NewerVersionAvailable_CallsDownload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSelfUpdateRepository(ctrl)
	mockRepo.EXPECT().LatestVersion().Return("v1.3.0", nil)
	mockRepo.EXPECT().Download("v1.3.0").Return(nil)

	svc := service.NewSelfUpdateService(mockRepo, service.CurrentVersion("v1.2.3"))
	result, err := svc.Update()

	require.NoError(t, err)
	assert.False(t, result.AlreadyLatest)
	assert.Equal(t, "1.2.3", result.OldVersion)
	assert.Equal(t, "1.3.0", result.NewVersion)
}

func TestSelfUpdateService_LatestVersionError_Propagates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSelfUpdateRepository(ctrl)
	mockRepo.EXPECT().LatestVersion().Return("", errors.New("network unreachable"))

	svc := service.NewSelfUpdateService(mockRepo, service.CurrentVersion("v1.0.0"))
	_, err := svc.Update()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check for updates")
}

func TestSelfUpdateService_DownloadError_Propagates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockSelfUpdateRepository(ctrl)
	mockRepo.EXPECT().LatestVersion().Return("v2.0.0", nil)
	mockRepo.EXPECT().Download("v2.0.0").Return(errors.New("permission denied"))

	svc := service.NewSelfUpdateService(mockRepo, service.CurrentVersion("v1.0.0"))
	_, err := svc.Update()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to download update")
}
