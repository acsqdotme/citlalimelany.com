package album

import (
	"encoding/json"
	"errors"
	"log"
	"sort"

	bolt "go.etcd.io/bbolt"
)

const (
	pathToDB = "portfolio.db"
)

var (
	portBucket = []byte("Portfolio")
)

type Album struct {
	FileName    string `json:"filename"`
	Title       string `json:"title"`
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

func AddToDB(album Album) error {
	jason, err := json.Marshal(album)
	if err != nil {
		return err
	}

	db := openDB()
	defer closeDB(db)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(portBucket)
		v := b.Get([]byte(album.FileName))
		if v != nil {
			return errors.New("key " + album.FileName + " already exists!")
		}
		err := b.Put([]byte(album.FileName), jason)
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

func AggregateAlbums() (albums []Album, err error) {
	db := openDB()
	defer closeDB(db)

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(portBucket)
		err = b.ForEach(func(k, v []byte) error {
			var a Album
			if err = json.Unmarshal(v, &a); err != nil {
				return err
			}
			albums = append(albums, a)
			return nil
		})
		return err
	})
	if err != nil {
		return []Album{}, err
	}

	sort.Slice(albums, func(i, j int) bool {
		return albums[i].Date > albums[j].Date
	})
	return albums, nil
}

func TestAlbum() (a Album) {
	a.FileName = "bjork"
	a.Title = "Björk Collection"
	a.Date = "2023-09-30"
	a.Description = "epic pics of björk"
	a.Thumbnail = Img{
		Src:   "https://www.newyorker.com/wp-content/uploads/2004/08/040823_r13315-967.jpg",
		Alt:   "Bjork with cool, asymmetric ponytail",
		Title: "bachelorette is her best song btw",
	}
	a.Photos = []Img{
		{
			Src:   "https://celebs-place.com/gallery/bjork/2-45.jpg",
			Alt:   "Björk forlornly dreaming",
			Title: "awesome hair",
		}, {
			Src:   "https://cdn.theatlantic.com/assets/media/img/mt/2013/11/bjork/lead_large.jpg",
			Alt:   "bjork and a tv",
			Title: "'You shouldn't let poets lie to you'",
		}, {
			Src: "https://townsquare.media/site/838/files/2016/01/bjork-boombox.jpg?w=1200&h=0&zc=1&s=0&a=t&q=89",
			Alt: "bjork with a boom box",
			// testing no title
		},
	}

	return a
}
