package service

import (
  "io"
  "fmt"
  "log"
  "net"
  "net/http"
)

func handleMethod(w http.ResponseWriter, req *http.Request, method string) (bool){
  if req.Method == method {
    return true
  }
  
  w.WriteHeader(404)
  w.Write([]byte("404 not found"))
  return false
}

func handleIndexPage(w http.ResponseWriter, req *http.Request) {
  if allow := handleMethod(w, req, "GET"); !allow {
    return
  }
  
  io.WriteString(w, `
<html><body>
<form method="POST" action="code">
  <label>OAuth2 Code</label>
  <input type="text" name="code"/>
  <input type="submit" value="submit"/>
</form>
</body></html>`)
}

func newFuncHandleCodePage(l net.Listener, codeChan chan<- string) (func (http.ResponseWriter, *http.Request)) {
  return func(w http.ResponseWriter, req *http.Request) {
    if allow := handleMethod(w, req, "POST"); !allow {
      return
    }
    
    if err := req.ParseForm(); err != nil || len(req.Form["code"]) == 0 || req.Form["code"][0] == "" {
      return
    }
    
    codeChan<- req.Form["code"][0]
    l.Close()
  }
}

func openHTTPServer(port, authURL string, codeChan chan<- string) (err error) {
  var l net.Listener
  if l, err = net.Listen("tcp", ":"+port); err != nil {
    return
  }
  
  http.Handle("/authorize", http.RedirectHandler(authURL, 301))
  http.HandleFunc("/", handleIndexPage)
  http.HandleFunc("/code", newFuncHandleCodePage(l, codeChan))
  
  server := &http.Server{}
  
  fmt.Println("HTTP Server open, Link to one of following URLs to get oauth2 code")
  for _, addr := range serverAddressList() {
    fmt.Println("http://" + addr + ":" + port + "/authorize")
  }
  fmt.Println("And link to one of following URLs to submit oauth2 code")
  for _, addr := range serverAddressList() {
    fmt.Println("http://" + addr + ":" + port + "/")
  }
  
  if err = server.Serve(l); err != nil {
    return
  }
  return
}

func serverAddressList() (addrList []string) {
  if addrs, err := net.InterfaceAddrs(); err != nil {
    return make([]string, 0)
  } else {
    index := 0
    list := make([]string, len(addrs))
    for _, addr := range addrs {
      if _, ipNet, err := net.ParseCIDR(addr.String()); err == nil && len(ipNet.IP) == 4 {
        list[index] = ipNet.IP.String()
        index += 1
      }
    }
    addrList = list[:index]
  }
  return addrList
}

func getCodeFromServer(authURL string) (code string, err error) {
  var codeChan chan string = make(chan string, 1)
  var port string = httpServerPort()
  go openHTTPServer(port, authURL, codeChan)
  fmt.Println(authURL)
  code = <-codeChan
  log.Println("Fetch oauth2 code", code)
  return code, nil
}
