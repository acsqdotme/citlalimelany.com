package album

import (
	"log"

	bolt "go.etcd.io/bbolt"
)

const (
  pathToDB = "portfolio.db"
)

var (
  portBucket = []byte("Portfolio")
)

type Collection struct {
  Name string `json:"name"`
  Date string `json:"date"`
  Description string `json:"description"`
  Thumbnail Img `json:"thumbnail"`
  Photos []Img `json:"photos"`
}

type Img struct {
  Src string `json:"src"`
  Alt string `json:"alt"`
  Title string `json:"title"`
}

func MakeDB() error {
  log.Println("Initializing DB")
  db, err := bolt.Open(pathToDB, 0600, nil)
  if err != nil {
    return err
  }
  defer db.Close()

  err = db.Update(func(tx *bolt.Tx) error {
    _, err := tx.CreateBucketIfNotExists(portBucket)
    return err
  })

  return err
}
