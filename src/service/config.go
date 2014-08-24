package service

import (
  "fmt"
  "path"
  "time"
)

const (
  defaultOAuth2FolderPath = "config"
  defaultBackupFolderPath = "gobackup"
  
  // Database
  defaultDatabaseOptionsFolderPath = "config"
  defaultPGFileName = "pg.json"
  
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
