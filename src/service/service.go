package service

import (
  "io"
  "os"
  "os/exec"
  "errors"
)

type ServiceType int

type DeviceService interface {
  Authorize() (error)
  GetQuota() (int64, int64)
  BackupFromReader(io.ReadCloser) (error)
}

const (
  GoogleService ServiceType = 0
  DropboxService ServiceType = 1
)

func verifyQuota(device DeviceService, size int64) (err error) {
  usaged, total := device.GetQuota()
  if size > (total - usaged) {
    err = errors.New("Quota is not enought")
  }
  return
}

func BackupFile(device DeviceService, filepath string) (err error) {
  var file *os.File
  if file, err = os.Open(filepath); err != nil {
    return
  }

  var size int64
  if size, err = file.Seek(0, 2); err != nil {
    return
  } else if _, err = file.Seek(0, 0); err != nil {
    return
  }
  
  if err = verifyQuota(device, size); err != nil {
    return
  }
  
  err = device.BackupFromReader(file)
  return
}

func BackupDatabase(device DeviceService) (err error) {
  cmd := exec.Command("pg_dump")
  stdout, err := cmd.StdoutPipe()
  if err != nil {
    return err
  }
  defer stdout.Close()
  
  if err := cmd.Start(); err != nil {
    return err
  }
  
  return device.BackupFromReader(stdout)
}
