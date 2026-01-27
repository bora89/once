package ui

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstallForm_FillAndSubmit(t *testing.T) {
	form := NewInstallForm()

	// Type into image ref field
	form = typeText(form, "nginx:latest")
	assert.Equal(t, "nginx:latest", form.ImageRef())

	// Press enter to move to hostname field
	form = pressEnter(form)
	assert.Equal(t, fieldHostname, form.focused)

	// Type into hostname field
	form = typeText(form, "myapp.example.com")
	assert.Equal(t, "myapp.example.com", form.Hostname())

	// Press enter to move to install button
	form = pressEnter(form)
	assert.Equal(t, fieldInstallButton, form.focused)

	// Press enter to submit
	form, cmd := form.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	require.NotNil(t, cmd, "expected a command from submit")

	msg := cmd()
	submitMsg, ok := msg.(InstallFormSubmitMsg)
	require.True(t, ok, "expected InstallFormSubmitMsg, got %T", msg)
	assert.Equal(t, "nginx:latest", submitMsg.ImageRef)
	assert.Equal(t, "myapp.example.com", submitMsg.Hostname)
}

func TestInstallForm_TabNavigation(t *testing.T) {
	form := NewInstallForm()
	assert.Equal(t, fieldImageRef, form.focused)

	// Tab through all fields
	form = pressTab(form)
	assert.Equal(t, fieldHostname, form.focused)

	form = pressTab(form)
	assert.Equal(t, fieldInstallButton, form.focused)

	form = pressTab(form)
	assert.Equal(t, fieldCancelButton, form.focused)

	// Tab wraps around
	form = pressTab(form)
	assert.Equal(t, fieldImageRef, form.focused)
}

func TestInstallForm_ShiftTabNavigation(t *testing.T) {
	form := NewInstallForm()

	// Shift+Tab goes backwards, wrapping to cancel button
	form = pressShiftTab(form)
	assert.Equal(t, fieldCancelButton, form.focused)

	form = pressShiftTab(form)
	assert.Equal(t, fieldInstallButton, form.focused)
}

func TestInstallForm_CancelButton(t *testing.T) {
	form := NewInstallForm()

	// Navigate to cancel button
	form.focused = fieldCancelButton

	// Press enter
	_, cmd := form.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	require.NotNil(t, cmd)

	msg := cmd()
	_, ok := msg.(InstallFormCancelMsg)
	assert.True(t, ok, "expected InstallFormCancelMsg, got %T", msg)
}

// Helpers

func typeText(form InstallForm, text string) InstallForm {
	for _, r := range text {
		form, _ = form.Update(tea.KeyPressMsg{Code: r, Text: string(r)})
	}
	return form
}

func pressEnter(form InstallForm) InstallForm {
	form, _ = form.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	return form
}

func pressTab(form InstallForm) InstallForm {
	form, _ = form.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	return form
}

func pressShiftTab(form InstallForm) InstallForm {
	form, _ = form.Update(tea.KeyPressMsg{Code: tea.KeyTab, Mod: tea.ModShift})
	return form
}
