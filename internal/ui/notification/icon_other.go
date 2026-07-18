//go:build !darwin

package notification

import (
	_ "embed"
)

//go:embed nextcode-icon-solo.png
var Icon []byte
