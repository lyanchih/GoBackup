package service

import (
  "io"
  "os"
  "errors"
  "strings"
  "encoding/json"
  "github.com/golang/oauth2"
  "github.com/stacktic/dropbox"
)

const (
  uriDropboxAuth  = "https://www.dropbox.com/1/oauth2/authorize"
  uriDropboxToken = "https://api.dropbox.com/1/oauth2/token"
)

type dropboxService struct {
  db *dropbox.Dropbox
}

func NewDropboxService() (device DeviceService) {
  svc := &dropboxService{
    db: dropbox.NewDropbox(),
  }
  svc.db.RootDirectory = "sandbox"
  return svc
}

func (*dropboxService) Type() (ServiceType) {
  return DropboxService
}

func (svc *dropboxService) Authorize() (err error) {
  var opts *oauth2.Options
  if file, err := os.Open(svc.Type().optionsFilePath()); err != nil {
    return err
  } else {
    defer file.Close()
    
    if err = json.NewDecoder(file).Decode(&opts); err != nil {
      return err
    }
  }
  svc.db.SetAppInfo(opts.ClientID, opts.ClientSecret)
  
  tokenPath := svc.Type().tokenFilePath()
  var token *oauth2.Token
  if token = readAccessToken(tokenPath); token == nil {
    var conf *oauth2.Config
    if conf, err = oauth2.NewConfig(opts, uriDropboxAuth, uriDropboxToken); err != nil {
      return
    }
    
    if token, err = fetchAccessToken(conf); err != nil {
      return
    }
  }
  
  if token == nil {
    return errors.New("Can't receive oauth2 access token")
  }
  svc.db.SetAccessToken(token.AccessToken)
  
  saveAccessToken(tokenPath, token)
  return
}

func (svc *dropboxService) GetAccount() (err error) {
  return
}

func (svc *dropboxService) GetQuota() (usaged int64, total int64) {
  usaged = -1
  total = -1
  
  if account, err := svc.db.GetAccountInfo(); err != nil {
    return
  } else {
    usaged = account.QuotaInfo.Shared + account.QuotaInfo.Normal
    total = account.QuotaInfo.Quota
  }
  return
}

func (svc *dropboxService) createFolder(path string) (err error) {
  _, err = svc.db.CreateFolder(path)
  if err != nil && strings.Contains(err.Error(), "a file or folder already exists at path") {
    err = nil
  }
  return
}

func (svc *dropboxService) createBackupFolder() (err error) {
  folderPath := backupFolderPath()
  err = svc.createFolder(folderPath)
  return
}

func (svc *dropboxService) BackupFromReader(reader io.ReadCloser) (err error) {
  defer reader.Close()
  
  if err = svc.createBackupFolder(); err != nil {
    return
  }

  _, err = svc.db.UploadByChunk(reader, -1, newBackupFilePath(), true, "")
  return
}
