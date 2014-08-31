package service

import (
  "os"
  "fmt"
  "path"
  "time"
  "strconv"
)

const (
  defaultOAuth2FolderPath = "config"
  defaultBackupFolderPath = "gobackup"
  
  // Database
  defaultDatabaseOptionsFolderPath = "config"
  defaultPGFileName = "pg.json"
  
  // HTTP
  defaultHTTPServerPort = 3300
  
  // Goole Drive Relative Config
  defaultGoogleOAuth2Options = "google_oauth2_options.json"
  defaultGoogleOAuth2Token = "google_oauth2_token.json"
  
  // Dropbox Config
  defaultDropboxOAuth2Options = "dropbox_oauth2_options.json"
  defaultDropboxOAuth2Token = "dropbox_oauth2_token.json"
)

//
// OAuth2 Config
//

func oauth2FolderPath() (string) {
  return defaultOAuth2FolderPath;
}

func (sType ServiceType) optionsFilePath() (string) {
  var filename string
  switch sType {
  case GoogleService: filename = defaultGoogleOAuth2Options
  case DropboxService: fallthrough
  default: filename = defaultDropboxOAuth2Options
  }
  
  return path.Join(oauth2FolderPath(), filename)
}

func (sType ServiceType) tokenFilePath() (string) {
  var filename string
  switch sType {
  case GoogleService: filename = defaultGoogleOAuth2Token
  case DropboxService: fallthrough
  default: filename = defaultDropboxOAuth2Token
  }
  
  return path.Join(oauth2FolderPath(), filename)
}

//
// Backup Config
//

func backupFolderPath() (folder string) {
  return defaultBackupFolderPath
}

func newBackupFileName() (filename string) {
  return fmt.Sprintf("backup_%s.txt", time.Now().UTC().String())
}

func newBackupFilePath() (filepath string) {
  return path.Join(backupFolderPath(), newBackupFileName())
}

//
// HTTP Server
//

func isGetCodeFromServer() (res bool) {
  switch os.Getenv("FetchCode") {
  case "server": res = true
  case "cosole": res = false
  default:
    if isFg, err := isForeground(); err == nil && !isFg {
      res = true
    }
  }
  return
}

func httpServerPort() (port string) {
  return strconv.Itoa(defaultHTTPServerPort)
}

//
// Database Config
//

func databaseOptionsFolderPath() (string) {
  return defaultDatabaseOptionsFolderPath;
}

func databaseOptionsFileName() (string) {
  return defaultPGFileName;
}

func databaseOptionsFilePath() (string) {
  return path.Join(databaseOptionsFolderPath(), databaseOptionsFileName())
}
