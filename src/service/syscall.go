package service

import (
  "errors"
  "unsafe"
  "syscall"
)

func isForeground() (bool, error) {
  pgrp := syscall.Getpgrp()
  for _, fd := range []uintptr{uintptr(syscall.Stdin), uintptr(syscall.Stdout), uintptr(syscall.Stderr)} {
    pid := 0
    if _, _, err := syscall.RawSyscall(syscall.SYS_IOCTL, fd, syscall.TIOCGPGRP, uintptr(unsafe.Pointer(&pid))); err != 0 {
      return false, errors.New(err.Error())
    } else if pid != pgrp {
      return false, nil
    }
  }
  return true, nil
}
