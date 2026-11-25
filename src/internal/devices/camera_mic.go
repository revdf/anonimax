package devices

import (
	"fmt"
	"os"
	"strings"

	"gitlab.com/parrotsec/packages/anonsurf/internal/command"
	"gitlab.com/parrotsec/packages/anonsurf/internal/fileops"
)

const StatePath = "/run/anonsurf/loaded-av-modules"

var runner = command.Runner{}

var cameraModules = []string{
	"uvcvideo",
	"videobuf2_v4l2",
	"videobuf2_vmalloc",
	"videobuf2_common",
}

var audioModules = []string{
	"snd_hda_intel",
	"snd_hda_codec",
	"snd_usb_audio",
	"snd_hda_core",
}

// isLoaded checks if a kernel module is currently loaded via /sys/module.
func isLoaded(name string) bool {
	_, err := os.Stat("/sys/module/" + name)
	return err == nil
}

// Disable unloads camera and audio kernel modules and saves the list of
// modules that were active so Enable can restore exactly those.
func Disable() error {
	all := append(cameraModules, audioModules...)
	var removed []string

	for _, mod := range all {
		if !isLoaded(mod) {
			continue
		}
		if err := runner.Run("rmmod", mod); err != nil {
			// Module may be in use — skip it, don't abort the whole operation.
			continue
		}
		removed = append(removed, mod)
	}

	if len(removed) == 0 {
		return nil
	}

	data := strings.Join(removed, "\n") + "\n"
	if err := fileops.WriteAtomicWithDir(StatePath, []byte(data), 0600, 0700); err != nil {
		return fmt.Errorf("save AV modules state: %w", err)
	}
	return nil
}

// Enable reloads the modules that were unloaded by Disable.
func Enable() error {
	data, err := os.ReadFile(StatePath)
	if err != nil {
		// No state file means Disable was not called or already cleaned up.
		return nil
	}

	modules := strings.Fields(strings.TrimSpace(string(data)))
	var lastErr error
	for _, mod := range modules {
		if err := runner.Run("modprobe", mod); err != nil {
			lastErr = err
		}
	}

	if err := fileops.RemoveIfExists(StatePath); err != nil {
		return fmt.Errorf("clear AV modules state: %w", err)
	}
	return lastErr
}
