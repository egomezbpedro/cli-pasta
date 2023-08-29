package main

import (
	"log"
    "os"
    "sort"
	c "github.com/egomezbpedro/cli-pasta/clipboard"
	db "github.com/egomezbpedro/cli-pasta/database"
	"github.com/koki-develop/go-fzf"
)

var (
    stdlog, errlog *log.Logger
)
func init() {
    stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
    errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

func initializeDatabase() db.Database{
    // Initialize the Database
        db := db.Database{};
        db.BucketName = "data";
        db.DatabaseName = "/usr/local/var/pasta.db";
        return db;
}

func fuzzySearch(channel chan map[int]string, clip c.Clipboard) {
    
    searchMap := <-channel

    items := make([]string, 0, len(searchMap))
    keys := make([]int, 0, len(searchMap))
    
    for k := range searchMap{
        keys = append(keys, k)
    }
    // Reverse sort
    sort.Sort(sort.Reverse(sort.IntSlice(keys)))
 
    for _, k := range keys {
        items = append(items, searchMap[k])
    }

    stdlog.Printf("Items: %v", items)

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

    result := make(chan map[int]string)
    clip := c.Clipboard{}

    go db.ReadAllValues(db.DatabaseName, db.BucketName, result);

    fuzzySearch(result, clip)
}
