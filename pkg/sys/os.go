// Package sys provides os level helper functions for interacting with the system
package sys

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/phR0ze/n/pkg/opt"
	"github.com/pkg/errors"
)

// Copy copies src to dst recursively, creating destination directories as needed.
// Handles globbing e.g. Copy("./*", "../")
// The dst will be copied to if it is an existing directory.
// The dst will be a clone of the src if it doesn't exist.
// Doesn't follow links by default but can be turned on with &Opt{"follow", true}
func Copy(src, dst string, opts ...*opt.Opt) (err error) {
	clone := true
	var sources []string

	// Set following links to off by default
	defaultFollowOpt(&opts, false)

	// Get Abs src and dst roots
	var dstAbs, srcAbs string
	if dstAbs, err = Abs(dst); err != nil {
		return
	}
	if srcAbs, err = Abs(src); err != nil {
		return
	}

	// Handle globbing
	if sources, err = filepath.Glob(srcAbs); err != nil {
		err = errors.Wrapf(err, "failed to get glob for %s", srcAbs)
		return
	}

	// Fail no sources were found
	if len(sources) == 0 {
		err = errors.Errorf("failed to get any sources for %s", srcAbs)
		return
	}

	// Clone given src as dst vs copy into dst
	if IsDir(dstAbs) || len(sources) > 1 {
		clone = false
	}

	// Copy all sources to dst
	for _, srcRoot := range sources {

		// Walk over file structure
		err = Walk(srcRoot, func(srcPath string, srcInfo *FileInfo, e error) error {
			if e != nil {
				return e
			}

			// Set proper dst path
			var dstPath string
			if clone {
				dstPath = path.Join(dstAbs, strings.TrimPrefix(srcPath, srcRoot))
			} else {
				dstPath = path.Join(dstAbs, strings.TrimPrefix(srcPath, path.Dir(srcRoot)))
			}

			// Handle individual copies
			switch {

			// Create destination directories as needed
			case srcInfo.IsDir():
				if e = os.MkdirAll(dstPath, srcInfo.Mode()); e != nil {
					return e
				}

			// Copy dir links
			case srcInfo.IsSymlinkDir():
				var target string
				if target, e = srcInfo.SymlinkTarget(); e != nil {
					return e
				}
				if e = os.Symlink(target, dstPath); e != nil {
					return e
				}

			// Copy file
			default:
				CopyFile(srcPath, dstPath, newInfoOpt(srcInfo))
			}
			return nil
		}, newFollowOpt(false))
	}
	return
}

// CopyFile copies a single file from src to dsty, creating destination directories as needed.
// The dst will be copied to if it is an existing directory.
// The dst will be a clone of the src if it doesn't exist.
// Doesn't follow links by default but can be turned on with &Opt{"follow", true}
// Returns the destination path for copied file
func CopyFile(src, dst string, opts ...*opt.Opt) (result string, err error) {
	var srcPath, dstPath string
	var srcInfo, srcDirInfo *FileInfo

	// Set following links to off by default
	defaultFollowOpt(&opts, false)

	// Check the source for issues
	if srcInfo = infoOpt(opts); srcInfo != nil {
		if srcPath, err = Abs(srcInfo.path); err != nil {
			return
		}
	} else {
		if srcPath, err = Abs(src); err != nil {
			return
		}
		if srcInfo, err = Lstat(srcPath); err != nil {
			return
		}
	}

	// Source dir permissions to use for destination directories
	if srcDirInfo, err = Lstat(path.Dir(srcPath)); err != nil {
		return
	}

	// Error out if not a regular file or symlink
	if srcInfo.IsDir() || srcInfo.IsSymlinkDir() {
		err = errors.Errorf("src target is not a regular file or a symlink to a file")
		return
	}

	// Get correct destination path
	if dstPath, err = Abs(dst); err != nil {
		return
	}
	dstInfo, e := os.Stat(dstPath)
	switch {

	// Doesn't exist so this is the new destination name, ensure all paths exist
	case os.IsNotExist(e):
		if err = os.MkdirAll(path.Dir(dstPath), srcDirInfo.Mode()); err != nil {
			return
		}

	// Destination exists and is either a file to overwrite or a dir to copy into
	case e == nil:
		if dstInfo.IsDir() {
			dstPath = path.Join(dstPath, path.Base(srcPath))
		}

	// unknown error case
	default:
		err = errors.Wrapf(e, "failed to Stat destination %s", dst)
		return
	}

	// Handle links a bit differently
	if srcInfo.IsSymlink() {
		var target string
		if target, err = srcInfo.SymlinkTarget(); err != nil {
			return
		}
		if err = os.Symlink(target, dstPath); err != nil {
			return
		}
	} else {
		// Open srcPath for reading
		var fin *os.File
		if fin, err = os.Open(srcPath); err != nil {
			err = errors.Wrapf(err, "failed to open file %s for reading", srcPath)
			return
		}
		defer fin.Close()

		// Create dstPath for writing
		var fout *os.File
		if fout, err = os.Create(dstPath); err != nil {
			err = errors.Wrapf(err, "failed to create file %s", dstPath)
			return
		}

		// Copy srcPath to dstPath
		if _, err = io.Copy(fout, fin); err != nil {
			err = errors.Wrapf(err, "failed to copy data to file %s", dstPath)
			if e := fout.Close(); e != nil {
				err = errors.Wrapf(err, "failed to close file %s", dstPath)
			}
			return
		}

		// Sync to disk
		if err = fout.Sync(); err != nil {
			err = errors.Wrapf(err, "failed to sync data to file %s", dstPath)
			if e := fout.Close(); e != nil {
				err = errors.Wrapf(err, "failed to close file %s", dstPath)
			}
			return
		}

		// Close file for writing
		if err = fout.Close(); err != nil {
			err = errors.Wrapf(err, "failed to close file %s", dstPath)
			return
		}

		// Set permissions of dstPath same as srcPath
		if err = os.Chmod(dstPath, srcInfo.Mode()); err != nil {
			err = errors.Wrapf(err, "failed to chmod file %s", dstPath)
			return
		}
	}

	result = dstPath
	return
}

// Exists return true if the given path exists
func Exists(src string) bool {
	if target, err := Abs(src); err == nil {
		if _, err := os.Stat(target); err == nil {
			return true
		}
	}
	return false
}

// Darwin returns true if the OS is OSX
func Darwin() (result bool) {
	if runtime.GOOS == "darwin" {
		result = true
	}
	return
}

// Linux returns true if the OS is Linux
func Linux() (result bool) {
	if runtime.GOOS == "linux" {
		result = true
	}
	return
}

// Windows returns true if the OS is Windows
func Windows() (result bool) {
	if runtime.GOOS == "windows" {
		result = true
	}
	return
}

// MD5 returns the md5 of the given file
func MD5(target string) (result string, err error) {
	if target, err = Abs(target); err != nil {
		return
	}
	if !Exists(target) {
		return "", os.ErrNotExist
	}

	// Open target file for reading
	var f *os.File
	if f, err = os.Open(target); err != nil {
		err = errors.Wrapf(err, "failed opening target file %s", target)
		return
	}
	defer f.Close()

	// Create a new md5 hash and copy in file bits
	hash := md5.New()
	if _, err = io.Copy(hash, f); err != nil {
		err = errors.Wrapf(err, "failed copying file data into hash from %s", target)
		return
	}

	// Compute 32 byte hash
	result = hex.EncodeToString(hash.Sum(nil))

	return
}

// MkdirP creates the target directory and any parent directories needed
// and returns the ABS path of the created directory
func MkdirP(target string, perms ...uint32) (dir string, err error) {
	if dir, err = Abs(target); err != nil {
		return
	}

	// Get/set default permission
	perm := os.FileMode(0755)
	if len(perms) > 0 {
		perm = os.FileMode(perms[0])
	}

	// Create directory
	if err = os.MkdirAll(dir, perm); err != nil {
		err = errors.Wrapf(err, "failed creating directories for %s", dir)
		return
	}

	return
}

// Move the src path to the dst path. If the dst already exists and is not a directory
// src will replace it. If there is an error it will be of type *LinkError. Wraps
// os.Rename but fixes the issue where dst name is required. Returns the new location
func Move(src, dst string) (result string, err error) {

	// Add src base name to dst directory to fix golang oversight
	if IsDir(dst) {
		dst = path.Join(dst, path.Base(src))
	}
	if err = os.Rename(src, dst); err != nil {
		err = errors.Wrapf(err, "failed renaming file %s", src)
		return
	}
	result = dst
	return
}

// Pwd returns the current working directory
func Pwd() (pwd string) {
	pwd, _ = os.Getwd()
	return
}

// ReadBytes returns the entire file as []byte
func ReadBytes(filepath string) (result []byte, err error) {
	if filepath, err = Abs(filepath); err != nil {
		return
	}

	if result, err = ioutil.ReadFile(filepath); err != nil {
		err = errors.Wrapf(err, "failed reading the file %s", filepath)
		return
	}
	return
}

// ReadLines returns a new slice of string representing lines
func ReadLines(filepath string) (result []string, err error) {
	if filepath, err = Abs(filepath); err != nil {
		return
	}

	var data []byte
	if data, err = ioutil.ReadFile(filepath); err != nil {
		err = errors.Wrapf(err, "failed reading the file %s", filepath)
		return
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}
	return
}

// ReadString returns the entire file as a string
func ReadString(filepath string) (result string, err error) {
	if filepath, err = Abs(filepath); err != nil {
		return
	}

	var data []byte
	if data, err = ioutil.ReadFile(filepath); err != nil {
		err = errors.Wrapf(err, "failed reading the file %s", filepath)
		return
	}
	result = string(data)
	return
}

// Remove the given target file or empty directory. If there is an
// error it will be of type *PathError
func Remove(target string) error {
	return os.Remove(target)
}

// RemoveAll removes the target path and any children it contains.
// It removes everything it can but returns the first error it encounters.
// If the target path does not exist nil is returned
func RemoveAll(target string) error {
	return os.RemoveAll(target)
}

// Symlink creates newname as a symbolic link to link. If there is an error,
// it will be of type *LinkError.
func Symlink(src, link string) error {
	return os.Symlink(src, link)
}

// Touch creates an empty text file similar to the linux touch command
func Touch(filepath string) (path string, err error) {
	if path, err = Abs(filepath); err != nil {
		return
	}

	var f *os.File
	if f, err = os.Create(path); err != nil {
		err = errors.Wrapf(err, "failed creating/truncating file %s", filepath)
		return
	}

	// Ignoring close in the error case above is ok as the file pointer will be nil
	if err = f.Close(); err != nil {
		err = errors.Wrapf(err, "failed closing file %s", filepath)
		return
	}
	return
}

// WriteBytes is a pass through to ioutil.WriteBytes with default permissions
func WriteBytes(filepath string, data []byte, perms ...uint32) (err error) {
	if filepath, err = Abs(filepath); err != nil {
		return
	}

	perm := os.FileMode(0644)
	if len(perms) > 0 {
		perm = os.FileMode(perms[0])
	}
	if err = ioutil.WriteFile(filepath, data, perm); err != nil {
		err = errors.Wrapf(err, "failed writing bytes to file %s", filepath)
		return
	}
	return
}

// WriteLines is a pass through to ioutil.WriteFile with default permissions
func WriteLines(filepath string, lines []string, perms ...uint32) (err error) {
	if filepath, err = Abs(filepath); err != nil {
		return
	}

	perm := os.FileMode(0644)
	if len(perms) > 0 {
		perm = os.FileMode(perms[0])
	}
	if err = ioutil.WriteFile(filepath, []byte(strings.Join(lines, "\n")), perm); err != nil {
		err = errors.Wrapf(err, "failed writing lines to file %s", filepath)
		return
	}
	return
}

// WriteStream reads from the io.Reader and writes to the given file using io.Copy
// thus never filling memory i.e. streaming.  dest will be overwritten if it exists.
func WriteStream(reader io.Reader, filepath string, perms ...uint32) (err error) {
	if filepath, err = Abs(filepath); err != nil {
		return
	}

	perm := os.FileMode(0644)
	if len(perms) > 0 {
		perm = os.FileMode(perms[0])
	}

	var writer *os.File
	flags := os.O_CREATE | os.O_TRUNC | os.O_WRONLY
	if writer, err = os.OpenFile(filepath, flags, perm); err != nil {
		err = errors.Wrapf(err, "failed opening file %s for writing", filepath)
		return
	}

	if _, err = io.Copy(writer, reader); err != nil {
		err = errors.Wrap(err, "failed copying stream data")
		if e := writer.Close(); e != nil {
			err = errors.Wrapf(err, "failed to close file %s", filepath)
		}
		return
	}
	if err = writer.Sync(); err != nil {
		err = errors.Wrapf(err, "failed syncing stream to file %s", filepath)
		if e := writer.Close(); e != nil {
			err = errors.Wrapf(err, "failed to close file %s", filepath)
		}
		return
	}

	if err = writer.Close(); err != nil {
		err = errors.Wrapf(err, "failed to close file %s", filepath)
	}
	return
}

// WriteString is a pass through to ioutil.WriteFile with default permissions
func WriteString(filepath string, data string, perms ...uint32) (err error) {
	if filepath, err = Abs(filepath); err != nil {
		return
	}

	perm := os.FileMode(0644)
	if len(perms) > 0 {
		perm = os.FileMode(perms[0])
	}
	if err = ioutil.WriteFile(filepath, []byte(data), perm); err != nil {
		err = errors.Wrapf(err, "failed writing string to file %s", filepath)
		return
	}
	return
}

func slice(x []string, i, j int) (result []string) {

	// Convert to postive notation
	if i < 0 {
		i = len(x) + i
	}
	if j < 0 {
		j = len(x) + j
	}

	// Move start/end within bounds
	if i < 0 {
		i = 0
	}
	if j >= len(x) {
		j = len(x) - 1
	}

	// Specifically offsetting j to get an inclusive behavior out of Go
	j++

	// Only operate when indexes are within bounds
	// allow j to be len of s as that is how we include last item
	if i >= 0 && i < len(x) && j >= 0 && j <= len(x) {
		result = x[i:j]
	} else {
		result = []string{}
	}
	return
}
