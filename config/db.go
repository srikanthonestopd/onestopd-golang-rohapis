package config

import (
	"fmt"
	"log"
	"time"

	"github.com/couchbase/gocb/v2"
)

// Global variables for Couchbase
var Cluster *gocb.Cluster
var Collection *gocb.Collection

// Initialize Couchbase connection
func InitCouchbase() {
	var err error
	Cluster, err = gocb.Connect("couchbases://cb.tixyuflyocisfx7d.cloud.couchbase.com", gocb.ClusterOptions{
		Username: "Administrator",
		Password: "Admin123!",
	})
	if err != nil {
		log.Fatalf("Failed to connect to Couchbase: %v", err)
	}

	// Open bucket and collection
	bucket := Cluster.Bucket("roh-apis")
	err = bucket.WaitUntilReady(30*time.Second, nil)
	if err != nil {
		log.Fatalf("Bucket is not ready: %v", err)
	}

	Collection = bucket.Scope("myscope").Collection("mycollection")
	fmt.Println("✅ Connected to Couchbase and collection is ready!")
}
