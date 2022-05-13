package omedb

import (
	"log"
	"encoding/json"
	"encoding/binary"

	bolt "go.etcd.io/bbolt"
	webpush "github.com/SherClockHolmes/webpush-go"
)

const dbPath = "../ome.db"


//func main () {
//	s := &webpush.Subscription{}
//      json.Unmarshal([]byte(tempSubscription), s)
//	addSubscription(*s)
//	allSubscriptions := getSubscriptions()
//	fmt.Println(allSubscriptions)
//}

func itob(v int) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, uint64(v))
    return b
}

func addSubscription (subscription webpush.Subscription) error {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("subscriptions"))

		id, _ := bucket.NextSequence()
		buf, err := json.Marshal(subscription)
		if (err != nil) {
			panic("Couldn't marshal the Subscription")
		}
		return bucket.Put(itob(int(id)), buf)
	})
}

func getSubscriptions () []webpush.Subscription {
	var subscriptions []webpush.Subscription

	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("subscriptions"))
		cursor := bucket.Cursor()

		for key, subscriptionData := cursor.First(); key != nil; key, subscriptionData = cursor.Next() {
			subscriptionObject := &webpush.Subscription{}
			json.Unmarshal(subscriptionData, subscriptionObject)
			subscriptions = append(subscriptions, *subscriptionObject)
		}
		return nil
	})
	return subscriptions
}

