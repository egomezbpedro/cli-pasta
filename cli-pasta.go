package main

import (
	"log"
	c "github.com/egomezbpedro/cli-pasta/clipboard"
	db "github.com/egomezbpedro/cli-pasta/database"
	"github.com/koki-develop/go-fzf"
)

func main() {

    // Initialize the Database
    db := db.Database{}
    db.BucketName = "data"
    db.DatabaseName = "pasta.db"

    result := make(chan []string)
    clip := c.Clipboard{}

    go db.ReadAllValues(db.DatabaseName, db.BucketName, result);
    
    items := <-result

    f, err := fzf.New()
    if err != nil {
        log.Fatal(err)
    }

    idxs, err := f.Find(items, func(i int) string { return items[i] })
    if err != nil {
        log.Fatal(err)
    }

    for _, i := range idxs {
        clip.WriteClipboard(items[i])
    }
    close(result)
}