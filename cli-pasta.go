package main

import (
	"log"
	c "github.com/egomezbpedro/cli-pasta/clipboard"
	db "github.com/egomezbpedro/cli-pasta/database"
	"github.com/koki-develop/go-fzf"
)

func initializeDatabase() db.Database{
    // Initialize the Database
        db := db.Database{};
        db.BucketName = "data";
        db.DatabaseName = "/usr/local/var/pasta.db";
        return db;
}

func fuzzySearch(channel chan []string, clip c.Clipboard) {
    items := <-channel

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
    close(channel)
}

func main() {

    db := initializeDatabase();

    result := make(chan []string)
    clip := c.Clipboard{}

    go db.ReadAllValues(db.DatabaseName, db.BucketName, result);

    fuzzySearch(result, clip)
}