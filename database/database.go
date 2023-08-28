package database

import (
    "os"
	"log"
	"strconv"
	"time"
    "encoding/base64"
    //"encoding/binary"
	"github.com/boltdb/bolt"
)

type Database struct {
	BucketName   string
	DatabaseName string
}
var (
    stdlog, errlog *log.Logger
)
func init() {
    stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
    errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}


/*
*
Open the connection to the key/value Database. By default if no Database exists it creates a new instance.

Return: Database connection - @type: *bolt.DB
*/
func (d *Database) OpenDatabaseRW(DatabaseName string) *bolt.DB {
	db, err := bolt.Open(DatabaseName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
        errlog.Fatal("Failed to open the databse in RW mode", err)
	}
	return db
}
func (d *Database) OpenDatabaseRO(DatabaseName string) *bolt.DB {
	db, err := bolt.Open(DatabaseName, 0600, &bolt.Options{ReadOnly: true})
	if err != nil {
        errlog.Fatal("Failed to open the databse in RO mode", err)
	}
	return db
}

/*
*
Close the connection.

Args: Database connectionn - @type: *bolt.DB
*/
func (d *Database) CloseDatabase(db *bolt.DB) {
    err := db.Close()
    if err != nil {
        errlog.Fatal("Failed to close  the databse", err)
    }
}

/*
*
Creates a new Bucket in the Database. Retruns an error if the bucket has no name,
if it fails to create it or if the bucket name is to long
*/
func (d *Database) CreateBucket(database, bucket string) {

	db := d.OpenDatabaseRW(database)

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(bucket))
		if err != nil {
            errlog.Println("Fail to create bucket", err)
			return err
		}
		return nil
	})
	d.CloseDatabase(db)
}

/*
Search for a key in the database

TODO: This function is similar to ReadAllValues
*/
func (d *Database) searchForValue(database, bucket, value string) bool{

	result := false

	db := d.OpenDatabaseRO(database)
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucket))
		c := b.Cursor()
		for _, v := c.First(); v != nil; _, v = c.Next() {
            if len(string(v)) > 0 {
                data, err := base64.StdEncoding.DecodeString(string(v))
                if err != nil {
                    errlog.Println("Fail to decode database:", err)
                }
                if string(data) == value {
                    result = true 
                    break
                }
            }
		}
		return nil
	})
	d.CloseDatabase(db)
	return result
}


/*
*
Writes a key/valye pair to the Database.
Args:

	bucket name - @type string
	value       - @type string
*/
func (d *Database) WriteToBucket(database, bucket, value string) {

    str := base64.StdEncoding.EncodeToString([]byte(value))
    exist := d.searchForValue(database, bucket, value)

    if !exist {
        db := d.OpenDatabaseRW(database)
        db.Update(func(tx *bolt.Tx) error {
            b := tx.Bucket([]byte(bucket))
            id, _ := b.NextSequence()
            key := []byte(strconv.FormatUint(id, 10))
            err := b.Put([]byte(key), []byte(str))
            return err
        })
        d.CloseDatabase(db)
    }
}

/*
*
Reads all the key/values pairs from the Database
Args:

	bucket name - @type string
*/
func (d *Database) ReadAllValues(database, bucket string, result chan []string) {

    items := []string{""}

    db := d.OpenDatabaseRO(database)
    db.View(func(tx *bolt.Tx) error {
        // Assume bucket exists and has keys
        b := tx.Bucket([]byte(bucket))

        b.ForEach(func(k, v []byte) error {
            if len(string(v)) > 0 {
                data, err := base64.StdEncoding.DecodeString(string(v))
                if err != nil {
                    errlog.Println("Fail to read/decode values from database", err)
                }
                intK, _ := strconv.Atoi(string(k))
                stdlog.Printf("%d/%s ", intK, string(data))
                // Insert elements into the items array using intK as the index
                items = append(items, string(data))
            }
            return nil
        })
        if len(items) > 0 {
            result <- items
        }
        d.CloseDatabase(db)
        return nil
    })
}
