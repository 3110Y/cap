package handler_test

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/3110Y/cap/internal/application/service"
	"github.com/3110Y/cap/internal/domain/entity"
	"github.com/3110Y/cap/internal/integration/mocks"
	"github.com/3110Y/cap/internal/presentation/handler"
)

func newRepositoryHandler(t *testing.T, mockRepo *mocks.MockRepositoryRepository) *handler.RepositoryHandler {
	t.Helper()
	svc := service.NewRepositoryService(mockRepo)
	return handler.NewRepositoryHandler(svc)
}

func newCmdWithOutput() (*cobra.Command, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	return cmd, buf
}

func TestRepositoryHandler_Add_PrintsIDAndURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepositoryRepository(ctrl)
	mockRepo.EXPECT().Add(gomock.Any()).DoAndReturn(func(r *entity.Repository) error {
		return nil
	})

	h := newRepositoryHandler(t, mockRepo)
	cmd, buf := newCmdWithOutput()

	err := h.Add(cmd, []string{"https://github.com/example/repo"})

	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "Added:")
	assert.Contains(t, out, "https://github.com/example/repo")
}

func TestRepositoryHandler_List_PrintsTable(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repos := []*entity.Repository{
		{ID: "abc12345", URL: "https://github.com/user/r1"},
		{ID: "def67890", URL: "https://github.com/user/r2"},
	}
	mockRepo := mocks.NewMockRepositoryRepository(ctrl)
	mockRepo.EXPECT().List().Return(repos, nil)

	h := newRepositoryHandler(t, mockRepo)
	cmd, buf := newCmdWithOutput()
	cmd.Flags().Int("page", 1, "")
	cmd.Flags().Int("size", 50, "")

	err := h.List(cmd, nil)

	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "abc12345")
	assert.Contains(t, out, "def67890")
	assert.Contains(t, out, "https://github.com/user/r1")
	assert.Contains(t, out, "https://github.com/user/r2")
}

func TestRepositoryHandler_List_Empty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepositoryRepository(ctrl)
	mockRepo.EXPECT().List().Return([]*entity.Repository{}, nil)

	h := newRepositoryHandler(t, mockRepo)
	cmd, buf := newCmdWithOutput()
	cmd.Flags().Int("page", 1, "")
	cmd.Flags().Int("size", 50, "")

	err := h.List(cmd, nil)

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "No repositories")
}

func TestRepositoryHandler_Delete_PrintsRemoved(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockRepositoryRepository(ctrl)
	mockRepo.EXPECT().Delete("abc12345").Return(nil)

	h := newRepositoryHandler(t, mockRepo)
	cmd, buf := newCmdWithOutput()

	err := h.Delete(cmd, []string{"abc12345"})

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Removed: abc12345")
}
