package main

import (
  "os"
  "log"
  "./service"
)

func main() {
  var device service.DeviceService
  var err error
  
  switch os.Getenv("ServiceType") {
  case "google":
    device = service.NewGoogleService()
  case "dropbox": fallthrough
  default:
    device = service.NewDropboxService()
  }
  
  if err = device.Authorize(); err != nil {
    log.Fatal(err)
  }
  
  if err = service.InitDBEnvironment(); err != nil {
    log.Fatal(err)
  }
  
  if err = service.BackupDatabase(device); err != nil {
    log.Fatal(err)
  }
}
