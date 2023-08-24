package database

import (
	"fmt"
	"log"
	"strconv"
	"time"
	"github.com/boltdb/bolt"
)

type Database struct {
	BucketName   string
	DatabaseName string
}

/*
*
Open the connection to the key/value Database. By default if no Database exists it creates a new instance.

Return: Database connection - @type: *bolt.DB
*/
func (d *Database) OpenDatabaseRW(DatabaseName string) *bolt.DB {
	db, err := bolt.Open(DatabaseName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	return db
}
func (d *Database) OpenDatabaseRO(DatabaseName string) *bolt.DB {
	db, err := bolt.Open(DatabaseName, 0600, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

/*
*
Close the connection.

Args: Database connectionn - @type: *bolt.DB
*/
func (d *Database) CloseDatabase(db *bolt.DB) {
	db.Close()
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
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	d.CloseDatabase(db)
}

/*
*
Writes a key/valye pair to the Database.
Args:

	bucket name - @type string
	value       - @type string
*/
func (d *Database) WriteToBucket(database, bucket, value string) {

	db := d.OpenDatabaseRW(database)

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		id, _ := b.NextSequence()
		key := []byte(strconv.FormatUint(id, 10))
		err := b.Put([]byte(key), []byte(value))
		return err
	})
	d.CloseDatabase(db)
}

/*
*
Reads all the key/values pairs from the Database
Args:

	bucket name - @type string
*/
func (d *Database) ReadAllValues(database, bucket string, result chan []string) {
	db := d.OpenDatabaseRO(database)

	items := []string{}

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucket))

		c := b.Cursor()

		for _, v := c.First(); v != nil; _, v = c.Next() {
			if len(string(v)) > 0 {
				items = append(items, string(v))
			}
		}
		return nil
	})
	if len(items) > 0 {
		result <- items
	}
	d.CloseDatabase(db)
	return
}
