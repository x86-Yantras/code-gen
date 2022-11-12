package filesys

import (
	"fmt"
	"os"
)

type Fs struct{}

func (f *Fs) CreateDir(name string) error {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		err := os.MkdirAll(name, 0700)

		if err != nil {
			return err
		}
		return nil
	}
	fmt.Printf("Skipping: Dir %s already exists\n", name)
	return nil
}
