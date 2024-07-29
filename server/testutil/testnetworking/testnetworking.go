package testnetworking

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/buildbuddy-io/buildbuddy/server/testutil/testfs"
	"github.com/stretchr/testify/require"
)

// Setup sets up the test to be able to call networking functions.
// It skips the test if the required net tools aren't available.
func Setup(t *testing.T) {
	// Ensure ip tools are in PATH
	os.Setenv("PATH", os.Getenv("PATH")+":/usr/sbin:/sbin")

	// Make sure the 'ip' tool is available and that we have the necessary
	// permissions to use it.
	cmd := []string{"ip", "link"}
	if os.Getuid() != 0 {
		cmd = append([]string{"sudo", "--non-interactive"}, cmd...)
	}
	if b, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput(); err != nil {
		t.Logf("%s failed: %s: %s", cmd, err, strings.TrimSpace(string(b)))
		t.Skipf("test requires passwordless sudo for 'ip' command - run ./tools/enable_local_firecracker.sh")
	}

	// Set up a symlink in PATH so that 'iptables' points to 'iptables-legacy'.
	// Our Firecracker setup does not yet have nftables enabled and can't use
	// the newer iptables.
	iptablesLegacyPath, err := exec.LookPath("iptables-legacy")
	require.NoError(t, err)
	overrideBinDir := testfs.MakeTempDir(t)
	err = os.Symlink(iptablesLegacyPath, filepath.Join(overrideBinDir, "iptables"))
	require.NoError(t, err)
	err = os.Setenv("PATH", overrideBinDir+":"+os.Getenv("PATH"))
	require.NoError(t, err)
}
