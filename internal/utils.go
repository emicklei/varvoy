package internal

import (
	"log/slog"
	"os"
)

func osCopy(src, dst string) error {
	slog.Debug("copy", "src", src, "dst", dst)
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(dst, input, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
func osEnsureDir(dir string) error {
	// if err := os.Mkdir(dir, os.ModePerm); err != nil && errors.Unwrap(err) != fs.ErrExist {
	// 	slog.Error("failed to create temp dir", "err", err, "err.type", fmt.Sprintf("%T", err))
	// 	panic(err)
	// }
	os.Mkdir(dir, os.ModePerm)
	return nil
}
