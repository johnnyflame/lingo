package driver

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/juju/errors"
	"github.com/lingo-reviews/lingo/tenet/driver/binary"
	"github.com/lingo-reviews/lingo/util"
	"github.com/lingo-reviews/tenets/go/dev/api"
	"github.com/lingo-reviews/tenets/go/dev/tenet/log"
)

// Binary is a tenet driver to execute binary tenets found in ~/.lingo/tenets/<repo>/<tenet>
type Binary struct {
	*Base
}

// Check that a file exists at the expected location and is executable. Setup
// the service, but don't start it.
func (b *Binary) Service() (Service, error) {
	tenetPath, err := b.binPath()
	if err != nil {
		return nil, errors.Trace(err)
	}

	file, err := os.Open(tenetPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.Errorf("The tenet %q could not be found. Has it been installed?", b.Name)
		}
		return nil, errors.Trace(err)
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, errors.Trace(err)
	}
	if fi.Mode().Perm()&0x49 == 0 {
		return nil, errors.Errorf("%s not exectuable", tenetPath)
	}

	// Note: the service needs to be started and stopped.
	return binary.NewService(tenetPath), nil
}

func (b *Binary) binPath() (string, error) {
	if dir := os.Getenv("LINGO_BIN"); dir != "" {
		return filepath.Join(dir, b.Name), nil
	}
	lHome, err := util.LingoHome()
	if err != nil {
		return "", errors.Trace(err)
	}
	return filepath.Join(lHome, "tenets", b.Name), nil
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func (*Binary) EditFilename(filename string) (editedFilename string) {
	var absPath string
	var err error
	if absPath, err = filepath.Abs(filename); err == nil {
		return absPath
	}
	log.Printf("could not get absolute path for %q:%v", filename, err)
	return filename
}

func (*Binary) EditIssue(issue *api.Issue) (editedIssue *api.Issue) {
	start := issue.Position.Start.Filename
	end := issue.Position.End.Filename

	pwd, err := os.Getwd()
	if err != nil {
		log.Printf("could not get relative path for %q:%v", start, err)
	}

	issue.Position.Start.Filename, err = filepath.Rel(pwd, start)
	if err != nil {
		log.Printf("could not get relative path for %q:%v", start, err)
	}

	issue.Position.End.Filename, err = filepath.Rel(pwd, end)
	if err != nil {
		log.Printf("could not get relative path for %q:%v", end, err)
	}

	return issue
}
