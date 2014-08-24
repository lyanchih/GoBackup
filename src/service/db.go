package service

import (
  "os"
  "encoding/json"
)

func InitDBEnvironment() (err error) {
  var file *os.File
  if file, err = os.Open(databaseOptionsFilePath()); err != nil {
    return
  }
  var db_data *map[string]interface{}
  if err = json.NewDecoder(file).Decode(&db_data); err != nil {
    return
  }
  
  for k, v := range *db_data {
    if vv, ok := v.(string); !ok {
      continue
    } else if os.Getenv(k) == "" {
      os.Setenv(k, vv)
    }
  }
  
  return
}
