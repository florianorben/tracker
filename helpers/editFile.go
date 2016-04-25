package helpers

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"os/exec"
)

const defaultEditor = "vim"

func OpenInEditor(b []byte) ([]byte, error) {
	return createAndModify(b)
}
func OpenStringInEditor(s string) (string, error) {
	tmp, err := createAndModify([]byte(s))
	if err != nil {
		return "", err
	}

	return string(tmp), nil

}

func createAndModify(b []byte) ([]byte, error) {
	tmpFile, err := createTempFile(b)
	if err != nil {
		return make([]byte, 0), err
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	contents, err := modifyTempFile(tmpFile)
	if err != nil {
		return make([]byte, 0), err
	}

	return contents, nil
}

func createTempFile(b []byte) (*os.File, error) {
	tmpDir := os.TempDir()
	tmpFile, tmpFileErr := ioutil.TempFile(tmpDir, "timer")

	if tmpFileErr != nil {
		fmt.Printf("Error: Creating temp file failed: %s\n", PrintRed(tmpFileErr.Error()))
		return nil, tmpFileErr
	}

	tmpFile.Write(b)
	return tmpFile, nil
}

func modifyTempFile(f *os.File) ([]byte, error) {
	editor := viper.GetString("core.editor")
	if editor == "" {
		editor = defaultEditor
	}
	path, err := exec.LookPath(editor)
	if err != nil {
		fmt.Printf("Error: %s\n", PrintRed(fmt.Sprintf("Looking up %s in %s failed", editor, path)))
		return nil, err
	}

	cmd := exec.Command(path, f.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		fmt.Printf("Error: Start failed: %s", PrintRed(err.Error()))
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Printf("Error: %s\n", PrintRed(err.Error()))
		return nil, err
	}

	f.Seek(0, 0)
	tmpFileContents, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Error: %s\n", PrintRed(err.Error()))
		return nil, err
	}

	return tmpFileContents, nil
}
