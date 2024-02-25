package config

import (
    _ "fmt"
    "encoding/json"
    "io/ioutil"
)

type Config struct {
    Port uint16
    Address string
    DBpath string
    StoreDir string
    UrlPrefix string
    IdLen int
    AuthKeyLen int
}

var Conf = Config{}

func LoadConfig(path string) (Config, error) {
    conf := Config{}

    cont, err := ioutil.ReadFile(path)
    if err != nil {
        return conf, err
    }
    err = json.Unmarshal(cont, &conf)
    Conf = conf
    return conf, err
}
