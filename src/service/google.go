package service

import (
  "io"
  "os"
  "fmt"
  "time"
  "errors"
  "net/http"
  "encoding/json"
  "github.com/golang/oauth2"
  "code.google.com/p/google-api-go-client/drive/v2"
)

const (
  uriGoogleAuth = "https://accounts.google.com/o/oauth2/auth"
  uriGoogleToken = "https://accounts.google.com/o/oauth2/token"
  uriGoogleDriveScope = "https://www.googleapis.com/auth/drive"
)

type googleService struct {
  service *drive.Service
}

type googleOAuth2Options struct {
  Installed *struct{
    AuthUri string `json:"AuthUri"`
    ClientSecret string `json:"client_Secret"`
    TokenUri string `json:"token_uri"`
    ClientEmail string `json:"client_email"`
    RedirectUris []string `json:"redirect_uris"`
    ClientX509CertUrl string `json:"clientX509CertUrl"`
    ClientID string `json:"client_id"`
    AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
  } `json:"installed,omitempty"`
}

func (googleOpts *googleOAuth2Options) toOAuth2Options() (*oauth2.Options, error) {
  if googleOpts.Installed == nil {
    return nil, errors.New("File doesn't have enought information")
  }
  
  if len(googleOpts.Installed.RedirectUris) == 0 {
    return nil, errors.New("redirect_uris is empty")
  }
  
  return &oauth2.Options{
    ClientID: googleOpts.Installed.ClientID,
    ClientSecret: googleOpts.Installed.ClientSecret,
    RedirectURL: googleOpts.Installed.RedirectUris[0],
    Scopes: []string {uriGoogleDriveScope, "email", "profile"},
  }, nil
}

func NewGoogleService() (device DeviceService) {
  return &googleService{}
}

func (*googleService) Type() (ServiceType) {
  return GoogleService
}

func (svc *googleService) Authorize() (err error) {
  var opts *oauth2.Options
  var conf *oauth2.Config
  if file, err := os.Open(svc.Type().optionsFilePath()); err != nil {
    return err
  } else {
    defer file.Close()
    
    var googleOpts *googleOAuth2Options
    if err = json.NewDecoder(file).Decode(&googleOpts); err != nil {
      return err
    }
    if opts, err = googleOpts.toOAuth2Options(); err != nil {
      return err
    } else {
      if conf, err = oauth2.NewConfig(opts, uriGoogleAuth, uriGoogleToken); err != nil {
        return err
      }
    }
  }
  
  tokenPath := svc.Type().tokenFilePath()
  var token *oauth2.Token
  if token = readAccessToken(tokenPath); token == nil {
    
    if token, err = fetchAccessToken(conf); err != nil {
      return
    }
  }
  
  t := conf.NewTransport()
  t.SetToken(token)
  
  var driveService *drive.Service
  if driveService, err = drive.New(&http.Client{
    Transport: t,
  }); err != nil {
    return
  }
  svc.service = driveService
  
  saveAccessToken(tokenPath, token)
  return
}

func (svc *googleService) GetQuota() (usaged int64, total int64) {
  usaged = -1
  total = -1
  
  if about, err := svc.service.About.Get().Do(); err != nil {
    return
  } else {
    usaged = about.QuotaBytesUsedAggregate
    total = about.QuotaBytesTotal
  }
  return
}

func (svc *googleService) createFolder(path string) (file *drive.File, err error) {
  var folder *drive.File = &drive.File{
    Title: path,
    Parents: []*drive.ParentReference{
      &drive.ParentReference{
        Id: "root",
        IsRoot: true,
        Kind: "drive#parentReference",
      },
    },
    MimeType: "application/vnd.google-apps.folder",
  }
  
  file, err = svc.service.Files.Insert(folder).Do()
  return
}

func (svc *googleService) createBackupFolder() (parent *drive.ParentReference, err error) {
  folder := backupFolderPath()
  
  if folder == "" {
    parent = &drive.ParentReference{
      Id: "root",
      IsRoot: true,
      Kind: "drive#fileLink",
    }
    return parent, nil
  }
  
  if list, err := svc.service.Children.List("root").Q(fmt.Sprintf("title contains '%s'", folder)).Do(); err != nil {
    return nil, err
  } else if len(list.Items) == 0 {
    if file, err := svc.createFolder(folder); err != nil {
      return nil, err
    } else {
      parent = &drive.ParentReference{
        Id: file.Id,
        Kind: "drive#fileLink",
      }
    }
  } else {
    item := list.Items[0]
    parent = &drive.ParentReference{
      Id: item.Id,
      Kind: "drive#fileLink",
    }
  }
  
  return parent, nil
}

func (svc *googleService) BackupFromReader(reader io.ReadCloser) (err error) {
  defer reader.Close()
  
  var parent *drive.ParentReference
  if parent, err = svc.createBackupFolder(); err != nil {
    return
  }
  
  var file *drive.File = &drive.File{
    Title: newBackupFileName(),
    Description: fmt.Sprintf("Database backup at %s", time.Now().UTC().String()),
    Parents: []*drive.ParentReference{ parent },
  }
  
  _, err = svc.service.Files.Insert(file).Media(reader).Do()
  return
}
