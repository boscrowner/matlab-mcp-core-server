// Copyright 2025-2026 The MathWorks, Inc.

package globalmatlab

import (
	"context"
	"errors"
	"sync"

	"github.com/matlab/matlab-mcp-core-server/internal/entities"
	"github.com/matlab/matlab-mcp-core-server/internal/messages"
)

var (
	ErrLostMATLABConnection = errors.New("lost connection to specified existing MATLAB session")
)

type MATLABManagerAdaptor interface {
	StartSession(ctx context.Context, logger entities.Logger) (entities.SessionID, error)
	ShouldRestart() (bool, messages.Error)
	StopMATLABSession(ctx context.Context, sessionLogger entities.Logger, sessionID entities.SessionID) error
	GetMATLABSessionClient(ctx context.Context, sessionLogger entities.Logger, sessionID entities.SessionID) (entities.MATLABSessionClient, error)
}

type GlobalMATLAB struct {
	matlabManagerAdaptor MATLABManagerAdaptor

	lock              *sync.Mutex
	startSessionError error

	sessionID entities.SessionID
}

func New(
	matlabManagerAdaptor MATLABManagerAdaptor,
) *GlobalMATLAB {
	return &GlobalMATLAB{
		matlabManagerAdaptor: matlabManagerAdaptor,

		lock: &sync.Mutex{},
	}
}

func (g *GlobalMATLAB) Client(ctx context.Context, logger entities.Logger) (entities.MATLABSessionClient, error) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if g.startSessionError != nil {
		return nil, g.startSessionError
	}

	return g.getOrCreateClient(ctx, logger)
}

func (g *GlobalMATLAB) getOrCreateClient(ctx context.Context, logger entities.Logger) (entities.MATLABSessionClient, error) {
	var sessionIDZeroValue entities.SessionID

	// Start MATLAB if we don't have a session
	if g.sessionID == sessionIDZeroValue {
		sessionID, err := g.matlabManagerAdaptor.StartSession(ctx, logger)
		if err != nil {
			g.startSessionError = err
			return nil, err
		}
		g.sessionID = sessionID
	}

	// Try to get the client
	client, err := g.matlabManagerAdaptor.GetMATLABSessionClient(ctx, logger, g.sessionID)
	if err != nil {
		// Retry: stop old session and start a new one
		if stopErr := g.matlabManagerAdaptor.StopMATLABSession(ctx, logger, g.sessionID); stopErr != nil {
			logger.WithError(stopErr).Warn("failed to stop MATLAB session")
		}

		sessionID, err := g.restartMATLABSession(ctx, logger)
		if err != nil {
			g.startSessionError = err
			return nil, err
		}
		g.sessionID = sessionID

		return g.matlabManagerAdaptor.GetMATLABSessionClient(ctx, logger, g.sessionID)
	}

	return client, nil
}
func (g *GlobalMATLAB) restartMATLABSession(ctx context.Context, logger entities.Logger) (entities.SessionID, error) {
	shouldRestart, messagesErr := g.matlabManagerAdaptor.ShouldRestart()
	if messagesErr != nil {
		return 0, messagesErr
	}

	if !shouldRestart {
		return 0, ErrLostMATLABConnection
	}

	sessionID, err := g.matlabManagerAdaptor.StartSession(ctx, logger)
	if err != nil {
		return 0, err
	}

	return sessionID, nil
}
