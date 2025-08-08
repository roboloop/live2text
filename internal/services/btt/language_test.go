package btt_test

import (
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/roboloop/live2text/internal/services/btt"
)

func TestSelectLanguage(t *testing.T) {
	t.Parallel()

	lc := setupLanguageComponent(t, func(mc *minimock.Controller, s *btt.SettingsComponentMock) {
		s.SelectSettingsMock.Expect(minimock.AnyContext, "Selected Language", "LIVE2TEXT_SELECTED_LANGUAGE", "foo").
			Return(nil)
	})
	err := lc.SelectLanguage(t.Context(), "foo")
	require.NoError(t, err)
}

func TestSelectedLanguage(t *testing.T) {
	t.Parallel()

	lc := setupLanguageComponent(t, func(mc *minimock.Controller, s *btt.SettingsComponentMock) {
		s.SelectedSettingMock.Expect(minimock.AnyContext, "LIVE2TEXT_SELECTED_LANGUAGE").Return("foo", nil)
	})
	language, err := lc.SelectedLanguage(t.Context())
	require.NoError(t, err)
	require.Equal(t, "foo", language)
}

func setupLanguageComponent(
	t *testing.T,
	setupMocks func(mc *minimock.Controller, s *btt.SettingsComponentMock),
) btt.LanguageComponent {
	mc := minimock.NewController(t)
	s := btt.NewSettingsComponentMock(mc)

	if setupMocks != nil {
		setupMocks(mc, s)
	}

	return btt.NewLanguageComponent(s)
}
