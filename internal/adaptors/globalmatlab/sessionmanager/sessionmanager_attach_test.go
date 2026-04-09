// Copyright 2026 The MathWorks, Inc.

package sessionmanager_test

import (
	"context"
	"path/filepath"
	"testing"
	"testing/synctest"
	"time"

	"github.com/matlab/matlab-mcp-core-server/internal/adaptors/globalmatlab/sessionmanager"
	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/testutils"
	configmocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/application/config"
	mocks "github.com/matlab/matlab-mcp-core-server/mocks/adaptors/globalmatlab/sessionmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSessionManager_StartSession_AttachMode_HappyPath(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	type contextKeyType string
	const contextKey contextKeyType = "testKey"
	const contextKeyValue = "testValue"

	ctx := context.WithValue(t.Context(), contextKey, contextKeyValue)
	expectedSessionID := entities.SessionID(456)
	expectedMATLABRoot := filepath.Join("some", "matlab", "root")
	expectedMATLABStartingDir := filepath.Join("some", "starting", "dir")
	discoveryTimeout := 5 * time.Second

	expectedAttachSessionDetails := entities.AttachToExistingSession{}

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return(expectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMATLABStartingDir(mockLogger.AsMockArg()).
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		MATLABSessionMode().
		Return(entities.MATLABSessionModeExisting).
		Once()

	mockConfig.EXPECT().
		MATLABSessionDiscoveryTimeout().
		Return(discoveryTimeout).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.MatchedBy(func(ctx context.Context) bool {
			return ctx.Value(contextKey) == contextKeyValue
		}), mockLogger.AsMockArg(), expectedAttachSessionDetails).
		Return(expectedSessionID, nil).
		Once()

	starter := sessionmanager.New(
		mockMATLABManager,
		mockConfigFactory,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	sessionID, err := starter.StartSession(ctx, mockLogger)

	// Assert
	require.NoError(t, err)
	require.Equal(t, expectedSessionID, sessionID)
}

func TestSessionManager_StartSession_AttachMode_StartMATLABSessionError(t *testing.T) {
	// Arrange
	mockLogger := testutils.NewInspectableLogger()

	mockMATLABManager := &mocks.MockMATLABManager{}
	defer mockMATLABManager.AssertExpectations(t)

	mockConfigFactory := &mocks.MockConfigFactory{}
	defer mockConfigFactory.AssertExpectations(t)

	mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
	defer mockMATLABRootSelector.AssertExpectations(t)

	mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
	defer mockMATLABStartingDirSelector.AssertExpectations(t)

	mockConfig := &configmocks.MockConfig{}
	defer mockConfig.AssertExpectations(t)

	type contextKeyType string
	const contextKey contextKeyType = "testKey"
	const contextKeyValue = "testValue"

	ctx := context.WithValue(t.Context(), contextKey, contextKeyValue)
	expectedMATLABRoot := filepath.Join("some", "matlab", "root")
	expectedMATLABStartingDir := filepath.Join("some", "starting", "dir")
	discoveryTimeout := 100 * time.Millisecond
	expectedError := assert.AnError

	expectedAttachSessionDetails := entities.AttachToExistingSession{}

	mockMATLABRootSelector.EXPECT().
		SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
		Return(expectedMATLABRoot, nil).
		Once()

	mockMATLABStartingDirSelector.EXPECT().
		SelectMATLABStartingDir(mockLogger.AsMockArg()).
		Return(expectedMATLABStartingDir, nil).
		Once()

	mockConfigFactory.EXPECT().
		Config().
		Return(mockConfig, nil).
		Once()

	mockConfig.EXPECT().
		MATLABSessionMode().
		Return(entities.MATLABSessionModeExisting).
		Once()

	mockConfig.EXPECT().
		MATLABSessionDiscoveryTimeout().
		Return(discoveryTimeout).
		Once()

	mockMATLABManager.EXPECT().
		StartMATLABSession(mock.MatchedBy(func(ctx context.Context) bool {
			return ctx.Value(contextKey) == contextKeyValue
		}), mockLogger.AsMockArg(), expectedAttachSessionDetails).
		Return(entities.SessionID(0), expectedError).
		Maybe()

	starter := sessionmanager.New(
		mockMATLABManager,
		mockConfigFactory,
		mockMATLABRootSelector,
		mockMATLABStartingDirSelector,
	)

	// Act
	sessionID, err := starter.StartSession(ctx, mockLogger)

	// Assert
	require.ErrorIs(t, err, expectedError)
	require.Equal(t, entities.SessionID(0), sessionID)
}

func TestSessionManager_StartSession_AttachMode_ImmediateSuccess(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Arrange
		mockLogger := testutils.NewInspectableLogger()

		mockMATLABManager := &mocks.MockMATLABManager{}
		defer mockMATLABManager.AssertExpectations(t)

		mockConfigFactory := &mocks.MockConfigFactory{}
		defer mockConfigFactory.AssertExpectations(t)

		mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
		defer mockMATLABRootSelector.AssertExpectations(t)

		mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
		defer mockMATLABStartingDirSelector.AssertExpectations(t)

		mockConfig := &configmocks.MockConfig{}
		defer mockConfig.AssertExpectations(t)

		type contextKeyType string
		const contextKey contextKeyType = "testKey"
		const contextKeyValue = "testValue"

		ctx := context.WithValue(t.Context(), contextKey, contextKeyValue)
		expectedSessionID := entities.SessionID(123)
		expectedMATLABRoot := filepath.Join("some", "matlab", "root")
		expectedMATLABStartingDir := filepath.Join("some", "starting", "dir")
		discoveryTimeout := 10 * time.Second

		expectedAttachSessionDetails := entities.AttachToExistingSession{}

		mockMATLABRootSelector.EXPECT().
			SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
			Return(expectedMATLABRoot, nil).
			Once()

		mockMATLABStartingDirSelector.EXPECT().
			SelectMATLABStartingDir(mockLogger.AsMockArg()).
			Return(expectedMATLABStartingDir, nil).
			Once()

		mockConfigFactory.EXPECT().
			Config().
			Return(mockConfig, nil).
			Once()

		mockConfig.EXPECT().
			MATLABSessionMode().
			Return(entities.MATLABSessionModeExisting).
			Once()

		mockConfig.EXPECT().
			MATLABSessionDiscoveryTimeout().
			Return(discoveryTimeout).
			Once()

		mockMATLABManager.EXPECT().
			StartMATLABSession(mock.MatchedBy(func(ctx context.Context) bool {
				return ctx.Value(contextKey) == contextKeyValue
			}), mockLogger.AsMockArg(), expectedAttachSessionDetails).
			Return(expectedSessionID, nil).
			Once()

		starter := sessionmanager.New(
			mockMATLABManager,
			mockConfigFactory,
			mockMATLABRootSelector,
			mockMATLABStartingDirSelector,
		)

		// Act
		sessionID, err := starter.StartSession(ctx, mockLogger)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedSessionID, sessionID)
	})
}

func TestSessionManager_StartSession_AttachMode_RetryThenSuccess(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Arrange
		mockLogger := testutils.NewInspectableLogger()

		mockMATLABManager := &mocks.MockMATLABManager{}
		defer mockMATLABManager.AssertExpectations(t)

		mockConfigFactory := &mocks.MockConfigFactory{}
		defer mockConfigFactory.AssertExpectations(t)

		mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
		defer mockMATLABRootSelector.AssertExpectations(t)

		mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
		defer mockMATLABStartingDirSelector.AssertExpectations(t)

		mockConfig := &configmocks.MockConfig{}
		defer mockConfig.AssertExpectations(t)

		type contextKeyType string
		const contextKey contextKeyType = "testKey"
		const contextKeyValue = "testValue"

		ctx := context.WithValue(t.Context(), contextKey, contextKeyValue)
		expectedSessionID := entities.SessionID(123)
		expectedMATLABRoot := filepath.Join("some", "matlab", "root")
		expectedMATLABStartingDir := filepath.Join("some", "starting", "dir")
		retryInterval := 200 * time.Millisecond
		discoveryTimeout := 300 * time.Millisecond

		expectedAttachSessionDetails := entities.AttachToExistingSession{}

		mockMATLABRootSelector.EXPECT().
			SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
			Return(expectedMATLABRoot, nil).
			Once()

		mockMATLABStartingDirSelector.EXPECT().
			SelectMATLABStartingDir(mockLogger.AsMockArg()).
			Return(expectedMATLABStartingDir, nil).
			Once()

		mockConfigFactory.EXPECT().
			Config().
			Return(mockConfig, nil).
			Once()

		mockConfig.EXPECT().
			MATLABSessionMode().
			Return(entities.MATLABSessionModeExisting).
			Once()

		mockConfig.EXPECT().
			MATLABSessionDiscoveryTimeout().
			Return(discoveryTimeout).
			Once()

		mockMATLABManager.EXPECT().
			StartMATLABSession(mock.MatchedBy(func(ctx context.Context) bool {
				return ctx.Value(contextKey) == contextKeyValue
			}), mockLogger.AsMockArg(), expectedAttachSessionDetails).
			Return(entities.SessionID(0), assert.AnError).
			Once()

		mockMATLABManager.EXPECT().
			StartMATLABSession(mock.MatchedBy(func(ctx context.Context) bool {
				return ctx.Value(contextKey) == contextKeyValue
			}), mockLogger.AsMockArg(), expectedAttachSessionDetails).
			Return(expectedSessionID, nil).
			Once()

		starter := sessionmanager.New(
			mockMATLABManager,
			mockConfigFactory,
			mockMATLABRootSelector,
			mockMATLABStartingDirSelector,
		)
		starter.SetDiscoveryRetryInterval(retryInterval)

		// Act
		sessionID, err := starter.StartSession(ctx, mockLogger)

		// Assert
		require.NoError(t, err)
		require.Equal(t, expectedSessionID, sessionID)
	})
}

func TestSessionManager_StartSession_AttachMode_RetryExhausted(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Arrange
		mockLogger := testutils.NewInspectableLogger()

		mockMATLABManager := &mocks.MockMATLABManager{}
		defer mockMATLABManager.AssertExpectations(t)

		mockConfigFactory := &mocks.MockConfigFactory{}
		defer mockConfigFactory.AssertExpectations(t)

		mockMATLABRootSelector := &mocks.MockMATLABRootSelector{}
		defer mockMATLABRootSelector.AssertExpectations(t)

		mockMATLABStartingDirSelector := &mocks.MockMATLABStartingDirSelector{}
		defer mockMATLABStartingDirSelector.AssertExpectations(t)

		mockConfig := &configmocks.MockConfig{}
		defer mockConfig.AssertExpectations(t)

		type contextKeyType string
		const contextKey contextKeyType = "testKey"
		const contextKeyValue = "testValue"

		ctx := context.WithValue(t.Context(), contextKey, contextKeyValue)
		expectedMATLABRoot := filepath.Join("some", "matlab", "root")
		expectedMATLABStartingDir := filepath.Join("some", "starting", "dir")
		retryInterval := 200 * time.Millisecond
		discoveryTimeout := 300 * time.Millisecond

		expectedAttachSessionDetails := entities.AttachToExistingSession{}
		expectedError := assert.AnError

		mockMATLABRootSelector.EXPECT().
			SelectMATLABRoot(ctx, mockLogger.AsMockArg()).
			Return(expectedMATLABRoot, nil).
			Once()

		mockMATLABStartingDirSelector.EXPECT().
			SelectMATLABStartingDir(mockLogger.AsMockArg()).
			Return(expectedMATLABStartingDir, nil).
			Once()

		mockConfigFactory.EXPECT().
			Config().
			Return(mockConfig, nil).
			Once()

		mockConfig.EXPECT().
			MATLABSessionMode().
			Return(entities.MATLABSessionModeExisting).
			Once()

		mockConfig.EXPECT().
			MATLABSessionDiscoveryTimeout().
			Return(discoveryTimeout).
			Once()

		mockMATLABManager.EXPECT().
			StartMATLABSession(mock.MatchedBy(func(ctx context.Context) bool {
				return ctx.Value(contextKey) == contextKeyValue
			}), mockLogger.AsMockArg(), expectedAttachSessionDetails).
			Return(entities.SessionID(0), expectedError).
			Twice()

		starter := sessionmanager.New(
			mockMATLABManager,
			mockConfigFactory,
			mockMATLABRootSelector,
			mockMATLABStartingDirSelector,
		)
		starter.SetDiscoveryRetryInterval(retryInterval)

		// Act
		_, err := starter.StartSession(ctx, mockLogger)

		// Assert
		require.ErrorIs(t, err, expectedError)
	})
}
