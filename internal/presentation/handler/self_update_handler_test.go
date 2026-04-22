package handler_test

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/3110Y/cap/internal/application/service"
	"github.com/3110Y/cap/internal/integration/mocks"
	"github.com/3110Y/cap/internal/presentation/handler"
)

func newSelfUpdateHandler(
	t *testing.T,
	mockSelfUpdate *mocks.MockSelfUpdateRepository,
	currentVersion string,
) *handler.SelfUpdateHandler {
	t.Helper()
	svc := service.NewSelfUpdateService(mockSelfUpdate, service.CurrentVersion(currentVersion))
	return handler.NewSelfUpdateHandler(svc)
}

func TestSelfUpdateHandler_AlreadyLatest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSU := mocks.NewMockSelfUpdateRepository(ctrl)
	mockSU.EXPECT().LatestVersion().Return("1.0.0", nil)

	h := newSelfUpdateHandler(t, mockSU, "1.0.0")
	cmd := &cobra.Command{}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	err := h.SelfUpdate(cmd, nil)

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "Already up to date")
}

func TestSelfUpdateHandler_Updated(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSU := mocks.NewMockSelfUpdateRepository(ctrl)
	mockSU.EXPECT().LatestVersion().Return("1.1.0", nil)
	mockSU.EXPECT().Download("1.1.0").Return(nil)

	h := newSelfUpdateHandler(t, mockSU, "1.0.0")
	cmd := &cobra.Command{}
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)

	err := h.SelfUpdate(cmd, nil)

	require.NoError(t, err)
	out := buf.String()
	assert.Contains(t, out, "Updated:")
	assert.Contains(t, out, "1.0.0")
	assert.Contains(t, out, "1.1.0")
}
