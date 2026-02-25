package ui

import (
	"testing"

	"github.com/basecamp/once/internal/docker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActionsMenu_ShowsStopWhenRunning(t *testing.T) {
	app := &docker.Application{Running: true}
	m := NewActionsMenu(app)

	assert.Contains(t, m.View(), "Stop")
}

func TestActionsMenu_ShowsStartWhenStopped(t *testing.T) {
	app := &docker.Application{Running: false}
	m := NewActionsMenu(app)

	assert.Contains(t, m.View(), "Start")
}

func TestActionsMenu_SelectStartStop(t *testing.T) {
	app := &docker.Application{Running: true}
	m := NewActionsMenu(app)

	// Shortcut key goes to menu, which returns MenuSelectMsg
	cmd := m.Update(runeKeyMsg('s'))
	require.NotNil(t, cmd)
	msg := cmd()

	// Feed MenuSelectMsg back to get ActionsMenuSelectMsg
	cmd = m.Update(msg)
	require.NotNil(t, cmd)
	msg = cmd()

	selectMsg, ok := msg.(ActionsMenuSelectMsg)
	require.True(t, ok, "expected ActionsMenuSelectMsg, got %T", msg)
	assert.Equal(t, ActionsMenuStartStop, selectMsg.action)
	assert.Equal(t, app, selectMsg.app)
}

func TestActionsMenu_SelectRemove(t *testing.T) {
	app := &docker.Application{}
	m := NewActionsMenu(app)

	cmd := m.Update(runeKeyMsg('r'))
	require.NotNil(t, cmd)
	msg := cmd()

	cmd = m.Update(msg)
	require.NotNil(t, cmd)
	msg = cmd()

	selectMsg, ok := msg.(ActionsMenuSelectMsg)
	require.True(t, ok, "expected ActionsMenuSelectMsg, got %T", msg)
	assert.Equal(t, ActionsMenuRemove, selectMsg.action)
}

func TestActionsMenu_EscCloses(t *testing.T) {
	app := &docker.Application{}
	m := NewActionsMenu(app)

	cmd := m.Update(keyPressMsg("esc"))
	require.NotNil(t, cmd)

	msg := cmd()
	_, ok := msg.(ActionsMenuCloseMsg)
	assert.True(t, ok, "expected ActionsMenuCloseMsg, got %T", msg)
}

func TestActionsMenu_KeyboardNavigation(t *testing.T) {
	app := &docker.Application{Running: true}
	m := NewActionsMenu(app)

	// Navigate down to Remove
	m.Update(runeKeyMsg('j'))
	assert.Equal(t, 1, m.menu.selected)

	// Navigate back up to Start/Stop
	m.Update(runeKeyMsg('k'))
	assert.Equal(t, 0, m.menu.selected)

	// Navigate down and select with enter
	m.Update(runeKeyMsg('j'))
	cmd := m.Update(keyPressMsg("enter"))
	require.NotNil(t, cmd)
	msg := cmd()

	cmd = m.Update(msg)
	require.NotNil(t, cmd)
	msg = cmd()

	selectMsg, ok := msg.(ActionsMenuSelectMsg)
	require.True(t, ok, "expected ActionsMenuSelectMsg, got %T", msg)
	assert.Equal(t, ActionsMenuRemove, selectMsg.action)
}
