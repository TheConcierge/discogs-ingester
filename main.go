package main

import (
	"discogs-adapter/go-discogs"
	"fmt"
)

const (
	MY_USERNAME = "PseudoRomulus"
	INGEST_QUEUE_NAME = "import-to-vinyl-tracker"
	POST_INGEST_QUEUE_NAME = "imported-to-vinyl-tracker"
	INVALID_COLLECTION_ID = -1
)
func main() {
	client, err := discogs.New(&discogs.Options{
		UserAgent: "LMAO U NO",
		Currency: "USD",
		Token: "",
	})
	if err != nil {
		fmt.Printf("error establishing discogs connection: %s", err.Error())
		return
	}

	collections, err := client.CollectionFolders(MY_USERNAME)
	if err != nil {
		fmt.Printf("could not search user collection: %s", err.Error())
		return
	}
	ingest_coll_id := INVALID_COLLECTION_ID
	post_coll_id := INVALID_COLLECTION_ID
	for _, folder := range collections.Folders {
		if folder.Name == INGEST_QUEUE_NAME {
			ingest_coll_id = folder.ID
		}
		if folder.Name == POST_INGEST_QUEUE_NAME {
			post_coll_id = folder.ID
		}
	}
	// hit a case where we didn't find a proper collection ID
	// this could indicate that the collection names have changed or something
	if ingest_coll_id == INVALID_COLLECTION_ID {
		fmt.Printf("could not find a valid collection id for collection %s", INGEST_QUEUE_NAME)
		return
	}
	if post_coll_id == INVALID_COLLECTION_ID {
		fmt.Printf("could not find a valid collection id for collection %s", POST_INGEST_QUEUE_NAME)
		return
	}

	ingest_items, err := client.CollectionItemsByFolder(MY_USERNAME, ingest_coll_id,  &discogs.Pagination{Sort: "artist", SortOrder: "desc", PerPage: 10})
	for _, item := range ingest_items.Items {

		err = client.AddToCollectionFolder(MY_USERNAME, post_coll_id, item.ID)
		if err != nil {
			fmt.Printf("could not add %s to new collection: %s\n", item.BasicInformation.Title, err.Error())
		} else {
			fmt.Printf("successfully added %s to collection\n", item.BasicInformation.Title)
		}
	}
}
