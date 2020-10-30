package main

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
)

type Transaction struct {
	CCNum      string  `bson:"ccnum"`
	Date       string  `bson:"date"`
	Amount     float32 `bson:"amount"`
	Cvv        string  `bson:"cvv"`
	Expiration string  `bson:"exp"`
}

func main() {
	session, err := mgo.Dial("172.18.29.133")
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println("Connected")
	defer session.Close()

	results := make([]Transaction, 0)
	dbnames, err := session.DatabaseNames()
	if err != nil {
		log.Panicln(err)
	}
	for _, db := range dbnames {
		fmt.Printf("Database: %v\n", db)
	}

	iter := session.DB("store").C("transactions").Find(nil).Limit(100).Iter()
	if err := iter.All(&results); err != nil {
		log.Panicln(err)
	}
	defer iter.Close()

	fmt.Printf("\nResults\n")
	for _, txn := range results {
		fmt.Printf("%v\t%v\t%v\t%v\t%v\n", txn.CCNum, txn.Date, txn.Cvv, txn.Amount, txn.Expiration)
	}
}
