package libs

import (
	"os"
	"path/filepath"
)

var (
	RootPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
)
