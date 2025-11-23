package domain

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Validator struct {
	Path string
}

func (v *Validator) ValidateXML(x string) error {
	if v.Path == "" {
		v.Path = "/usr/bin/virt-xml-validate"
	}

	tmpDir := os.TempDir()

	f, err := os.CreateTemp(tmpDir, "ev-")
	if err != nil {
		return fmt.Errorf("create temp file, %w", err)
	}

	filePath := filepath.Join(tmpDir, f.Name())

	_, err = io.Copy(f, strings.NewReader(x))
	_ = f.Close()

	if err != nil {
		return fmt.Errorf("write temp file, %w", err)
	}
	defer func() { _ = os.Remove(filePath) }()

	c := exec.Command(v.Path, filePath)
	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr

	err = c.Run()
	if err != nil {
		return fmt.Errorf("run command, %w", err)
	}

	return nil
}
