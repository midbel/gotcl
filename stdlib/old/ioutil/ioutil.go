package ioutil

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"syscall"
)

func Copy(list []string, target string, force bool) error {
	return doCopy(list, target, force)
}

func doCopy(list []string, target string, force bool) error {
	i, err := os.Stat(target)
	if err == nil && i.IsDir() {
		return copyInto(list, target, force)
	}
	for _, i := range list {
		fi, err := os.Stat(i)
		if err != nil {
			return err
		}
		if fi.IsDir() {
			err = copyDir(i, target, force)
		} else {
			err = copyFile(i, target, force)
		}
		if err != nil {
			break
		}
	}
	return err
}

func copyInto(list []string, target string, force bool) error {
	err := os.MkdirAll(target, 0755)
	if err != nil {
		return err
	}
	for _, i := range list {
		err = doCopy([]string{i}, filepath.Join(target, filepath.Base(i)), force)
		if err != nil {
			break
		}
	}
	return err
}

func copyDir(dir, target string, force bool) error {
	if _, err := os.Stat(target); err == nil && !force {
		return pathError("copy", target)
	}
	return filepath.Walk(dir, func(file string, fi fs.FileInfo, err error) error {
		if err != nil || dir == file {
			return err
		}
		rel, err := filepath.Rel(file, dir)
		if err != nil {
			return err
		}
		to := filepath.Join(target, rel)
		if fi.IsDir() {
			return os.MkdirAll(to, 0755)
		}
		return copyFile(file, to, force)
	})
}

func copyFile(file, target string, force bool) error {
	if _, err := os.Stat(target); err == nil && !force {
		return pathError("copy", target)
	}
	r, err := os.Open(file)
	if err != nil {
		return err
	}
	defer r.Close()

	w, err := os.Create(target)
	if err != nil {
		return err
	}
	defer w.Close()

	if _, err = io.Copy(w, r); err != nil {
		return err
	}
	return changePerms(r, w)
}

func changePerms(r, w *os.File) error {
	s, err := r.Stat()
	if err != nil {
		return err
	}
	if err := w.Chmod(s.Mode()); err != nil {
		return err
	}
	sys, ok := s.Sys().(*syscall.Stat_t)
	if ok {
		err = w.Chown(int(sys.Uid), int(sys.Gid))
	}
	return err
}

func pathError(op, file string) error {
	return &fs.PathError{
		Op:   op,
		Path: file,
	}
}

func Move(list []string, target string, force bool) error {
	return nil
}

func Type(file string) (string, error) {
	fi, err := os.Stat(file)
	if err != nil {
		return "", err
	}
	var typ string
	switch mod := fi.Mode().Type(); {
	default:
		return "", fmt.Errorf("%s: unknown file type")
	case mod.IsRegular():
		typ = "file"
	case mod.IsDir():
		typ = "directory"
	case mod == fs.ModeCharDevice:
		typ = "character device"
	case mod == fs.ModeSymlink:
		typ = "symlink"
	case mod == fs.ModeDevice:
		typ = "device"
	case mod == fs.ModeNamedPipe:
		typ = "named pipe"
	case mod == fs.ModeSocket:
		typ = "socket"
	}
	return typ, nil
}

func Owned(file string, uid int) bool {
	fi, err := os.Stat(file)
	if err != nil {
		return false
	}
	sys, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return false
	}
	return sys.Uid == uint32(uid)
}

const (
	readable   = 0b100
	writable   = 0b010
	executable = 0b001
)

func Executable(file string, uid int) bool {
	perm := getPerm(file, uid)
	return perm&executable == executable
}

func Readable(file string, uid int) bool {
	perm := getPerm(file, uid)
	return perm&readable == readable
}

func Writable(file string, uid int) bool {
	perm := getPerm(file, uid)
	return perm&writable == writable
}

func getPerm(file string, uid int) os.FileMode {
	fi, err := os.Stat(file)
	if err != nil {
		return 0
	}
	sys, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return 0
	}
	perm := fi.Mode().Perm()
	if sys.Uid == uint32(uid) {
		perm = perm >> 6
	} else if isMember(uid, int(sys.Gid)) {
		perm = perm >> 3
	}
	return perm & 0b111
}

func isMember(uid int, gid int) bool {
	u, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		return false
	}
	gs, err := u.GroupIds()
	if err != nil {
		return false
	}
	sort.Strings(gs)
	i := sort.SearchStrings(gs, strconv.Itoa(gid))
	return i < len(gs) && gs[i] == strconv.Itoa(gid)
}
