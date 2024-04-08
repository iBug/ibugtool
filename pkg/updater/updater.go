package updater

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/iBug/ibugtool/internal/version"
	"golang.org/x/sys/unix"
)

const (
	Repository = "https://github.com/iBug/ibugtool"
	Filename   = "ibugtool"
)

func InstallDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}

func CanUpdate() error {
	dir, err := InstallDir()
	if err != nil {
		return err
	}
	if err := unix.Access(dir, unix.W_OK); err != nil {
		return err
	}
	return nil
}

func GetLatestTag(repo string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	url := fmt.Sprintf("%s/releases/latest", repo)
	resp, err := client.Head(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusFound {
		return "", fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	loc := resp.Header.Get("Location")
	tag, ok := strings.CutPrefix(loc, repo+"/releases/tag/")
	if !ok {
		return "", fmt.Errorf("unexpected redirect location %s", loc)
	}
	return tag, nil
}

func DownloadFile(w *os.File, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}
	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return err
	}
	defer gz.Close()
	_, err = io.Copy(w, gz)
	return err
}

func UpdateBinary(w io.Writer, force bool) error {
	err := CanUpdate()
	if err != nil {
		if !force {
			return err
		}
		fmt.Fprintf(w, "Error: %v, continuing anyway\n", err)
	}
	dir, _ := InstallDir()
	out, err := os.CreateTemp(dir, Filename+".*")
	defer os.Remove(out.Name())
	if err != nil {
		return err
	}
	tag, err := GetLatestTag(Repository)
	if err != nil {
		return fmt.Errorf("GetLatestTag: %w", err)
	}
	if tag == version.Version {
		if !force {
			fmt.Fprintf(w, "Already up to date: %s\n", tag)
			return nil
		}
		fmt.Fprintf(w, "Already up to date: %s, continuing anyway\n", tag)
	}
	url := fmt.Sprintf("%s/releases/download/%s/%s.gz", Repository, tag, Filename)
	fmt.Fprintf(w, "Downloading from %s\n", url)
	err = DownloadFile(out, url)
	if err != nil {
		return err
	}
	out.Close()
	err = os.Chmod(out.Name(), 0755)
	if err != nil {
		return err
	}
	err = os.Rename(out.Name(), filepath.Join(dir, Filename))
	if err != nil {
		return err
	}
	return nil
}
