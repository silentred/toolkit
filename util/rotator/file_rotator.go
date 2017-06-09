package rotator

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type FileRotator struct {
	path     string
	prefix   string
	ext      string
	currFile string

	fd      io.WriteCloser
	spliter Spliter
	// clean older log file
	Clean bool
}

func NewFileRotator(path, prefix, ext string, splt Spliter) *FileRotator {
	if prefix == "" {
		prefix = "app"
	}

	if ext == "" {
		ext = "log"
	}

	r := &FileRotator{
		path:    path,
		prefix:  prefix,
		ext:     ext,
		spliter: splt,
		Clean:   false,
	}
	_, err := r.getNextWriter()
	if err != nil {
		log.Fatal(err)
	}

	return r
}

func (r *FileRotator) getNextName() string {
	return filepath.Join(r.path, r.spliter.getNextName(r.prefix, r.ext))
}

func (r *FileRotator) removeOlderFile() error {
	pattern := fmt.Sprintf("%s_*.%s", r.prefix, r.ext)
	files, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file == r.currFile {
			continue
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		info, err := f.Stat()
		if err != nil {
			return err
		}

		t := atime(info)
		if time.Now().Sub(t) > 24*time.Hour {
			err = os.Remove(file)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *FileRotator) getNextWriter() (io.Writer, error) {
	if r.Clean {
		err := r.removeOlderFile()
		if err != nil {
			fmt.Println(err)
		}
	}

	file := r.getNextName()

	perm, err := strconv.ParseInt("0755", 8, 64)
	if err != nil {
		return nil, err
	}
	fd, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(perm))
	if err == nil {
		// Make sure file perm is user set perm cause of `os.OpenFile` will obey umask
		os.Chmod(file, os.FileMode(perm))

		// close old fd
		if r.fd != nil {
			r.fd.Close()
		}
		r.fd = fd

		// reset currSize
		r.spliter.reset()
		// set currFileName
		r.currFile = file
	} else {
		return nil, err
	}

	return fd, nil
}

func (r *FileRotator) Write(p []byte) (n int, err error) {
	n, err = r.fd.Write(p)
	if err != nil {
		return n, err
	}

	if err == nil && r.spliter.reachLimit(n) {
		_, err := r.getNextWriter()
		if err != nil {
			return n, err
		}
	}

	return n, nil
}
