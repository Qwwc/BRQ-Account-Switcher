package platform

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func withEpicEnv(t *testing.T, vars map[string]string) {
	t.Helper()
	prev := epicEnv
	epicEnv = func(k string) string { return vars[k] }
	t.Cleanup(func() { epicEnv = prev })
}

func mkdirAll(t *testing.T, parts ...string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(parts...), 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestEpicHiddenVariantClassicOnly(t *testing.T) {
	local := t.TempDir()
	mkdirAll(t, local, "EpicGamesLauncher", "Saved", "Config")
	withEpicEnv(t, map[string]string{"LOCALAPPDATA": local})

	hidden, installed := epicHiddenVariant()
	if hidden != epicNewPlatformName || installed != epicLegacyPlatformName {
		t.Fatalf("got hidden=%q installed=%q", hidden, installed)
	}
}

func TestEpicHiddenVariantNewOnly(t *testing.T) {
	local := t.TempDir()
	mkdirAll(t, local, "com.epicgames.eglv2")
	withEpicEnv(t, map[string]string{"LOCALAPPDATA": local})

	hidden, installed := epicHiddenVariant()
	if hidden != epicLegacyPlatformName || installed != epicNewPlatformName {
		t.Fatalf("got hidden=%q installed=%q", hidden, installed)
	}
}

func TestEpicHiddenVariantBothOrNeitherHidesNothing(t *testing.T) {
	local := t.TempDir()
	withEpicEnv(t, map[string]string{"LOCALAPPDATA": local})
	if h, _ := epicHiddenVariant(); h != "" {
		t.Fatalf("neither installed: hid %q", h)
	}
	mkdirAll(t, local, "com.epicgames.eglv2")
	mkdirAll(t, local, "EpicGamesLauncher", "Saved", "Config")
	if h, _ := epicHiddenVariant(); h != "" {
		t.Fatalf("both installed: hid %q", h)
	}
}

func TestApplyEpicVariantSwapsInstalledIntoPlace(t *testing.T) {
	s := AppSettings{
		PlatformOrder:     []string{"Steam", epicNewPlatformName, "Riot Games", epicLegacyPlatformName},
		DisabledPlatforms: []string{epicLegacyPlatformName, "Ubisoft"},
	}
	if !applyEpicVariant(&s, epicNewPlatformName, epicLegacyPlatformName) {
		t.Fatal("expected a change")
	}
	wantOrder := []string{"Steam", epicLegacyPlatformName, epicNewPlatformName, "Riot Games"}
	if !reflect.DeepEqual(s.PlatformOrder, wantOrder) {
		t.Fatalf("order = %v, want %v", s.PlatformOrder, wantOrder)
	}
	dis := sliceToSet(s.DisabledPlatforms)
	if _, ok := dis[epicNewPlatformName]; !ok {
		t.Fatal("uninstalled Epic should be disabled")
	}
	if _, ok := dis[epicLegacyPlatformName]; ok {
		t.Fatal("installed Epic should be enabled")
	}
	if applyEpicVariant(&s, epicNewPlatformName, epicLegacyPlatformName) {
		t.Fatal("second run should be a no-op")
	}
}

func TestApplyEpicVariantKeepsInstalledWhenBothOn(t *testing.T) {
	s := AppSettings{
		PlatformOrder:     []string{epicNewPlatformName, epicLegacyPlatformName, "Steam"},
		DisabledPlatforms: nil,
	}
	applyEpicVariant(&s, epicNewPlatformName, epicLegacyPlatformName)
	if !reflect.DeepEqual(s.PlatformOrder, []string{epicNewPlatformName, epicLegacyPlatformName, "Steam"}) {
		t.Fatalf("order changed unexpectedly: %v", s.PlatformOrder)
	}
	if !reflect.DeepEqual(s.DisabledPlatforms, []string{epicNewPlatformName}) {
		t.Fatalf("disabled = %v", s.DisabledPlatforms)
	}
}

func TestFilterEpicNames(t *testing.T) {
	got := filterEpicNames([]string{"Steam", epicNewPlatformName, epicLegacyPlatformName}, epicNewPlatformName)
	if !reflect.DeepEqual(got, []string{"Steam", epicLegacyPlatformName}) {
		t.Fatalf("got %v", got)
	}
}
