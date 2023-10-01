package album

import (
	"encoding/json"
	"errors"
	"log"

	bolt "go.etcd.io/bbolt"
)

const (
	pathToDB = "portfolio.db"
)

var (
	portBucket = []byte("Portfolio")
)

type Album struct {
	Name        string `json:"name"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Thumbnail   Img    `json:"thumbnail"`
	Photos      []Img  `json:"photos"`
}

type Img struct {
	Src   string `json:"src"`
	Alt   string `json:"alt"`
	Title string `json:"title"`
}

func openDB() (db *bolt.DB) {
	db, err := bolt.Open(pathToDB, 0600, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	return db
}

func closeDB(db *bolt.DB) {
	if db != nil {
		db.Close()
	}
}

func MakeDB() error {
	log.Println("Initializing DB")
	db := openDB()
	defer closeDB(db)

	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(portBucket)
		return err
	})

	return err
}

func AddToDB(key string, album Album) error {
	jason, err := json.Marshal(album)
	if err != nil {
		return err
	}

	db := openDB()
	defer closeDB(db)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(portBucket)
		v := b.Get([]byte(key))
		if v != nil {
			return errors.New("key " + key + " already exists!")
		}
		err := b.Put([]byte(key), jason)
		return err
	})

	return err
}

func FetchAlbum(key string) (album Album, err error) {
	db := openDB()
	defer closeDB(db)

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(portBucket)
		v := b.Get([]byte(key))
		err = json.Unmarshal(v, &album)
		return err
	})
	if err != nil {
		return Album{}, err
	}

	return album, nil
}

func DoesAlbumExist(key string) (exists bool, err error) {
	db := openDB()
	defer closeDB(db)

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(portBucket)
		exists = b.Get([]byte(key)) != nil
		return nil
	})

	return exists, err
}
