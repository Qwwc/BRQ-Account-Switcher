package platform

import (
	"os"
	"path/filepath"
	"strings"
)

// Epic ships two launchers that keep their sign-in in different places, so
// Platforms.json carries one entry for each. Both answer to the same Start
// Menu shortcut, which made first-launch detection turn both on and left users
// with two "Epic Games" tiles, one of which can never work on their PC.
// BRQ looks at what is actually installed and, when only one launcher is,
// keeps the other entry out of the home screen and the platform list.
const (
	epicNewPlatformName    = "Epic Games"
	epicLegacyPlatformName = "Epic Games (Legacy)"
)

// epicEnv is the environment the detection reads; tests replace it.
var epicEnv = os.Getenv

func epicPathExists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	return err == nil
}

func epicNewInstalled() bool {
	local := epicEnv("LOCALAPPDATA")
	if local == "" {
		return false
	}
	return epicPathExists(filepath.Join(local, "com.epicgames.eglv2")) ||
		epicPathExists(filepath.Join(local, "Epic Games", "Launcher", "EpicGamesLauncher.exe"))
}

func epicLegacyInstalled() bool {
	if local := epicEnv("LOCALAPPDATA"); local != "" &&
		epicPathExists(filepath.Join(local, "EpicGamesLauncher", "Saved", "Config")) {
		return true
	}
	for _, env := range []string{"ProgramFiles", "ProgramFiles(x86)"} {
		root := epicEnv(env)
		if root == "" {
			continue
		}
		for _, arch := range []string{"Win64", "Win32"} {
			if epicPathExists(filepath.Join(root, "Epic Games", "Launcher", "Portal", "Binaries", arch, "EpicGamesLauncher.exe")) {
				return true
			}
		}
	}
	return false
}

// epicHiddenVariant returns the Epic entry that has no launcher on this PC
// together with the one that does, or two empty strings when that can't be
// told (neither or both installed), in which case nothing is hidden.
func epicHiddenVariant() (hidden, installed string) {
	newOK, legacyOK := epicNewInstalled(), epicLegacyInstalled()
	switch {
	case newOK && !legacyOK:
		return epicLegacyPlatformName, epicNewPlatformName
	case legacyOK && !newOK:
		return epicNewPlatformName, epicLegacyPlatformName
	default:
		return "", ""
	}
}

func isHiddenEpic(name, hidden string) bool {
	return hidden != "" && strings.EqualFold(strings.TrimSpace(name), hidden)
}

// filterEpicNames drops the uninstalled Epic entry from a platform list.
func filterEpicNames(names []string, hidden string) []string {
	if hidden == "" {
		return names
	}
	out := make([]string, 0, len(names))
	for _, n := range names {
		if !isHiddenEpic(n, hidden) {
			out = append(out, n)
		}
	}
	return out
}

// applyEpicVariant corrects saved settings when the uninstalled Epic is turned
// on: it is turned off, and if the installed one was off, the installed one
// takes its place in the home order so the user keeps a working tile in the
// same spot. Reports whether settings changed and need saving.
func applyEpicVariant(s *AppSettings, hidden, installed string) bool {
	if s == nil || hidden == "" {
		return false
	}
	disabled := sliceToSet(s.DisabledPlatforms)
	if _, off := disabled[hidden]; off {
		return false
	}
	disabled[hidden] = struct{}{}
	if _, installedOff := disabled[installed]; installedOff {
		delete(disabled, installed)
		order := make([]string, 0, len(s.PlatformOrder)+1)
		placed := false
		for _, n := range s.PlatformOrder {
			switch {
			case strings.EqualFold(n, installed):
				continue
			case strings.EqualFold(n, hidden):
				order = append(order, installed, n)
				placed = true
			default:
				order = append(order, n)
			}
		}
		if !placed {
			order = append(order, installed)
		}
		s.PlatformOrder = order
	}
	s.DisabledPlatforms = setToSortedSlice(disabled)
	return true
}
