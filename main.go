package main

import (
    "fmt"
    "flag"
    "os"

    "git.fuwafuwa.moe/x3/ngfshare/config"
    "git.fuwafuwa.moe/x3/ngfshare/net"
    "git.fuwafuwa.moe/x3/ngfshare/db"
    "git.fuwafuwa.moe/x3/ngfshare/view"
);

func main() {

    confFilePath := flag.String("config", "", "The path for the config.json file")
    createAuthKey := flag.Bool("genauth", false, "Generate a new auth key and exit")

    flag.Parse()

    if *confFilePath == "" {
        fmt.Println("Failure: --config argument is required")
        os.Exit(1)
    }

    conf, err := config.LoadConfig(*confFilePath);
    if err != nil {
        fmt.Printf("Failed to load: %+v", err)
        os.Exit(1)
    }

    err = view.LoadTemplates()
    if err != nil {
        return
    }

    dbctx, err := db.Open(conf.DBpath)
    if err != nil {
        return
    }
    defer dbctx.Close()

    if *createAuthKey {
        key, err := dbctx.CreateNewAuthKey()
        if err != nil {
            fmt.Println("Failed to create auth key:", err)
            os.Exit(1)
        }
        fmt.Println(key)
        return
    }

    err = net.Start(conf)
    if err != nil {
        fmt.Printf("Failed to start listen: ", err)
        return
    }
}
