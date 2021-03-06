package service

import (
  "io"
  "os"
  "fmt"
  "bytes"
  "errors"
  "encoding/json"
  "github.com/golang/oauth2"
)

func readAccessToken(filename string) (token *oauth2.Token) {
  if file, err := os.Open(filename); err == nil {
    defer file.Close()
    json.NewDecoder(file).Decode(&token)
  }
  
  return
}

func saveAccessToken(filepath string, token *oauth2.Token) (err error) {
  if token == nil {
    return errors.New("Can't save nil oauth2 token")
  }
  
  var file *os.File
  if file, err = os.OpenFile(filepath, os.O_WRONLY | os.O_TRUNC | os.O_CREATE, 0600); err != nil {
    return
  }
  defer file.Close()
  
  var buffer *bytes.Buffer
  if bs, err := json.Marshal(token); err != nil {
    return err
  } else {
    buffer = bytes.NewBuffer(bs)
  }
  
  io.Copy(file, buffer)
  
  return nil
}

func fetchAccessToken(conf *oauth2.Config) (token *oauth2.Token, err error) {
  if conf == nil {
    err = errors.New("oauth2 config shouldn't be nil")
    return
  }
  
  authURL := conf.AuthCodeURL("stat", "", "")
  
  var code string
  if isGetCodeFromServer() {
    if code, err = getCodeFromServer(authURL); err != nil {
      return nil, err
    }
  } else if code, err = getCodeFromConsole(authURL); err != nil {
    return nil, err
  }

  token, err = conf.Exchange(code)
  return
}

func getCodeFromConsole(authURL string) (code string, err error) {
  fmt.Printf("Go to the following link in your browser: \n%s\n", authURL)
  fmt.Println("Please enter your code")
  fmt.Scanln(&code)
  return
}
