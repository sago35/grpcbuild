package grpcbuild

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// MkSendRequest ...
func MkSendRequest(files ...string) (*SendRequest, error) {
	f, err := StoreFiles(".", files...)
	if err != nil {
		return nil, err
	}

	return &SendRequest{Files: f}, nil
}

// StoreFiles ...
func StoreFiles(dir string, files ...string) ([]*File, error) {
	ret := []*File{}

	for _, file := range files {
		fi, err := os.Stat(filepath.Join(dir, file))
		if err != nil {
			return nil, err
		}

		if fi.IsDir() {

			fis, err := ioutil.ReadDir(filepath.Join(dir, file))
			if err != nil {
				return nil, err
			}

			for _, fi := range fis {
				sf, err := StoreFiles(dir, filepath.Join(file, fi.Name()))
				if err != nil {
					return nil, err
				}
				ret = append(ret, sf...)
			}

		} else {
			content, err := ioutil.ReadFile(filepath.Join(dir, file))
			if err != nil {
				return nil, err
			}

			ret = append(ret, &File{
				Filename: filepath.Base(file),
				Dir:      filepath.Dir(file),
				Data:     content,
			})
		}
	}

	return ret, nil
}

// StoreFilesChan ...
func StoreFilesChan(dir string, files ...string) (chan []*File, chan error) {
	ch := make(chan []*File)
	errCh := make(chan error)

	fis, err := _storeFiles(dir, files...)
	if err != nil {
		errCh <- err
		close(ch)
		close(errCh)
		return ch, errCh
	}

	go func() {
		defer close(ch)
		defer close(errCh)

		sf := []*File{}
		length := 0

		for _, f := range fis {
			f := f

			content, err := ioutil.ReadFile(filepath.Join(f.dir, f.filename))
			if err != nil {
				errCh <- err
				return
			}

			sf = append(sf, &File{
				Filename: f.filename,
				Dir:      f.dir,
				Data:     content,
			})
			length += len(content)

			if 10*1024*1024 < length {
				length = 0
				ch <- sf
				sf = []*File{}
			}
		}

		if 0 < length {
			ch <- sf
			sf = sf[:0]
		}

	}()

	return ch, errCh
}

type fileInfo struct {
	filename string
	dir      string
}

func _storeFiles(dir string, files ...string) ([]fileInfo, error) {
	ret := []fileInfo{}
	for _, file := range files {
		fi, err := os.Stat(filepath.Join(dir, file))
		if err != nil {
			return nil, err
		}

		if fi.IsDir() {
			fis, err := ioutil.ReadDir(filepath.Join(dir, file))
			if err != nil {
				return nil, err
			}

			for _, fi := range fis {
				sf, err := _storeFiles(dir, filepath.Join(file, fi.Name()))
				if err != nil {
					return nil, err
				}
				ret = append(ret, sf...)
			}
		} else {
			ret = append(ret, fileInfo{
				filename: filepath.Base(file),
				dir:      filepath.Join(dir, filepath.Dir(file)),
			})
		}
	}
	return ret, nil
}

// RetrieveFiles ...
func RetrieveFiles(dir string, files []*File) error {
	for _, f := range files {
		filename := filepath.Join(dir, f.GetDir(), f.GetFilename())
		err := os.MkdirAll(filepath.Dir(filename), 0777)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(filename, f.GetData(), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteFiles ...
func (res *ExecResponse) WriteFiles() error {
	return RetrieveFiles(".", res.GetFiles())
}
