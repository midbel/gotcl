package interp

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	stdout = "stdout"
	stderr = "stderr"
)

const fdprefix = "file"

const (
	modeReadOnly   = "r"
	modeReadBoth   = "r+"
	modeWriteOnly  = "w"
	modeWriteBoth  = "w+"
	modeAppendOnly = "a"
	modeAppendBoth = "a+"
)

type FileSet struct {
	files map[string]*os.File
	next  int
}

func Stdio() *FileSet {
	fs := FileSet{
		files: make(map[string]*os.File),
	}
	fs.register("0", os.Stdin)
	fs.register("1", os.Stdout)
	fs.register("2", os.Stderr)

	return &fs
}

func (fs *FileSet) Print(fd, str string) error {
	w, err := fs.lookup(fd)
	if err != nil {
		return err
	}
	fmt.Fprint(w, str)
	return nil
}

func (fs *FileSet) Println(fd, str string) error {
	if err := fs.Print(fd, str); err != nil {
		return err
	}
	w, err := fs.lookup(fd)
	if err != nil {
		return err
	}
	fmt.Fprintln(w)
	return nil
}

func (fs *FileSet) Open(file, mode string) (string, error) {
	var (
		f   *os.File
		err error
	)
	switch mode {
	default:
		return "", fmt.Errorf("%s: unknown mode given", mode)
	case modeReadOnly, "":
		f, err = os.Open(file)
	case modeReadBoth:
	case modeWriteOnly:
		f, err = os.Create(file)
	case modeWriteBoth:
	case modeAppendOnly:
	case modeAppendBoth:
	}
	if err != nil {
		return "", err
	}
	fd := fdprefix + strconv.Itoa(fs.next)
	fs.register(fd, f)
	return fd, nil
}

func (fs *FileSet) Close(fd string) error {
	w, err := fs.lookup(fd)
	if err != nil {
		return err
	}
	delete(fs.files, fd)
	return w.Close()
}

func (fs *FileSet) Seek(fd string, offset, whence int) (int64, error) {
	w, err := fs.lookup(fd)
	if err != nil {
		return 0, err
	}
	return w.Seek(int64(offset), whence)
}

func (fs *FileSet) Tell(fd string) (int64, error) {
	return fs.Seek(fd, 0, io.SeekCurrent)
}

func (fs *FileSet) Gets(fd string) (string, error) {
	w, err := fs.lookup(fd)
	if err != nil {
		return "", err
	}
	var (
		buf = make([]byte, 4096)
		ret []byte
	)
	off, err := w.Seek(0, io.SeekCurrent)
	if err != nil {
		return "", err
	}
	for {
		n, err := w.Read(buf)
		if err != nil && !errors.Is(err, io.EOF) {
			return "", err
		}
		if n == 0 {
			break
		}
		if x := bytes.IndexByte(buf[:n], '\n'); x >= 0 {
			ret = append(ret, buf[:x]...)
			if _, err = w.Seek(off+int64(x+1), io.SeekStart); err != nil {
				return "", err
			}
			break
		} else {
			ret = append(ret, buf[:n]...)
		}
	}
	return strings.TrimSpace(string(ret)), nil
}

func (fs *FileSet) Read(fd string, length int) (string, error) {
	w, err := fs.lookup(fd)
	if err != nil {
		return "", err
	}
	var b []byte
	if length <= 0 {
		b, err = io.ReadAll(w)
		return string(b), err
	}
	b = make([]byte, length)
	n, err := io.ReadFull(w, b)
	if err == nil || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) {
		return string(b[:n]), nil
	}
	return "", err
}

func (fs *FileSet) Eof(fd string) (bool, error) {
	w, err := fs.lookup(fd)
	if err != nil {
		return false, err
	}
	tell, err := w.Seek(0, io.SeekCurrent)
	if err != nil {
		return false, err
	}
	s, err := w.Stat()
	if err != nil {
		return false, err
	}
	return tell == s.Size(), nil
}

func (fs *FileSet) register(fd string, f *os.File) {
	fs.files[fd] = f
	fs.next++
}

func (fs *FileSet) lookup(fd string) (*os.File, error) {
	switch fd {
	case stdout, "":
		fd = "1"
	case stderr:
		fd = "2"
	default:
	}
	w, ok := fs.files[fd]
	if !ok {
		return nil, fmt.Errorf("%s: undefined channel", fd)
	}
	return w, nil
}