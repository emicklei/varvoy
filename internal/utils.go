package internal

import (
	"math/rand"
	"net"
	"os"
)

func osEnsureDir(dir string) error {
	// if err := os.Mkdir(dir, os.ModePerm); err != nil && errors.Unwrap(err) != fs.ErrExist {
	// 	slog.Error("failed to create temp dir", "err", err, "err.type", fmt.Sprintf("%T", err))
	// 	panic(err)
	// }
	_ = os.Mkdir(dir, os.ModePerm)
	return nil
}

// from https://github.com/phayes/freeport/blob/master/freeport.go
func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
