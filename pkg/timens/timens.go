package timens

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
)

// NsHandle is a handle to a network namespace. It can be cast directly
// to an int and used as a file descriptor.
type NsHandle int

// Close closes the NsHandle and resets its file descriptor to -1.
// It is not safe to use an NsHandle after Close() is called.
func (ns *NsHandle) Close() error {
	if err := unix.Close(int(*ns)); err != nil {
		return err
	}
	(*ns) = -1
	return nil
}

// Setns sets namespace using syscall. Note that this should be a method
// in syscall but it has not been added.
func Setns(ns NsHandle, nstype int) (err error) {
	return unix.Setns(int(ns), nstype)
}

// Set sets the current time namespace to the namespace represented
// by NsHandle.
func Set(ns NsHandle) (err error) {
	return Setns(ns, unix.CLONE_NEWTIME)
}

// New creates a new time namespace, sets it as current and returns
// a handle to it.
func New() (ns NsHandle, err error) {
	if err := unix.Unshare(unix.CLONE_NEWTIME); err != nil {
		return -1, err
	}
	return Get()
}
// Get gets a handle to the current threads network namespace.
func Get() (NsHandle, error) {
	return GetFromThread(os.Getpid(), unix.Gettid())
}

// GetFromPath gets a handle to a network namespace
// identified by the path
func GetFromPath(path string) (NsHandle, error) {
	fd, err := unix.Open(path, unix.O_RDONLY|unix.O_CLOEXEC, 0)
	if err != nil {
		return -1, err
	}
	return NsHandle(fd), nil
}

// GetFromThread gets a handle to the network namespace of a given pid and tid.
func GetFromThread(pid, tid int) (NsHandle, error) {
	return GetFromPath(fmt.Sprintf("/proc/%d/task/%d/ns/time_for_children", pid, tid))
}
