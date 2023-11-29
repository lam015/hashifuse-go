package main

import (
    "github.com/hairyhenderson/go-fsimpl/vaultfs"
    "fmt"
    "log"
    "net/url"
    "io/fs"
    "encoding/json"
    "bytes"
)

var prettyJSON bytes.Buffer

func create_vaultfs(vault_base_url string, secret_engine string) {
    
    url, err := url.Parse(vault_base_url)
    if err != nil {
        log.Fatal(err)
    }

    vfs, err := vaultfs.New(url)
    if err != nil {
        log.Fatal(err)
    }

    fs.WalkDir(vfs, secret_engine, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            log.Fatal(err)
        }

        if !d.IsDir() {
            fmt.Println("[secret] " + path)
            
            secret_content, err := fs.ReadFile(vfs, path)
            if err != nil {
                log.Fatal(err)
            }

            err = json.Indent(&prettyJSON, secret_content, "", "\t")
            if err != nil {
                log.Fatal(err)
            }

            vault_secrets[path] = prettyJSON.String()
            prettyJSON.Reset()

        } else {
            fmt.Println("[   dir] " + path)
        }

        return nil
    })
    
    /*
    fmt.Println("---\nvault_secrets map:")
    fmt.Println(vault_secrets)
    */
      
}
