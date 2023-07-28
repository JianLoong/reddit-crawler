package main

import (
	"crawlers/models"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Wrapper for Storage client
type StorageService struct {
	db *gorm.DB
}

// Top level declarations for the storeService and sqlite
var (
	storeService = &StorageService{}
)

func initialiseDatabase(name string) *gorm.DB {

	var sqlite_name = "databases/" + name + ".sqlite"

	db, err := gorm.Open(sqlite.Open(sqlite_name), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func migrate() {

	// Migrate the schema
	storeService.db.AutoMigrate(&models.Comment{})
	storeService.db.AutoMigrate(&models.Submission{})

}

// Creates a request and returns the response
func request(url string) []byte {

	response, err := http.Get(url)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return responseData

}

func getComments(sub models.Submission) {

	var commentObject models.CommentsResponse

	var commentData = request("https://www.reddit.com/" + sub.Permalink + ".json")

	json.Unmarshal(commentData, &commentObject)

	fmt.Println("Obtained comments")

	if len(commentObject) == 0 {
		return
	}

	var comments []models.Comment

	for j := 0; j < len(commentObject[1].Data.Children); j++ {

		var comment models.Comment

		comment.Message = commentObject[1].Data.Children[j].Data.Body
		comment.CreatedUTC = commentObject[1].Data.Children[j].Data.CreatedUtc
		comment.Score = uint8(commentObject[1].Data.Children[j].Data.Score)
		// comment.SubmissionID = responseObject.Data.Children[i].Data.ID
		comment.SubmissionID = sub.SubmissionID
		comment.CommentID = commentObject[1].Data.Children[j].Data.ID

		comments = append(comments, comment)

		// storeService.db.Clauses(clause.OnConflict{
		// 	UpdateAll: true,
		// }).Create(&comment)

	}

	storeService.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&comments)

}

func crawl(subreddit_name string, no_of_post string) {

	var url = "https://www.reddit.com/r/" + subreddit_name + ".json?limit=" + no_of_post

	var responseObject models.SubmissionResponse

	var responseData = request(url)

	json.Unmarshal(responseData, &responseObject)

	fmt.Println("Processing")

	// var submissions []models.Submission

	for i := 0; i < len(responseObject.Data.Children); i++ {
		var sub models.Submission

		sub.SubmissionID = responseObject.Data.Children[i].Data.ID
		sub.Url = responseObject.Data.Children[i].Data.URL
		// sub.ID = responseObject.Data.Children[i].Data.ID
		sub.Title = responseObject.Data.Children[i].Data.Title
		sub.CreatedUTC = responseObject.Data.Children[i].Data.CreatedUtc
		sub.Selftext = responseObject.Data.Children[i].Data.Selftext
		sub.Permalink = responseObject.Data.Children[i].Data.Permalink
		sub.Score = uint8(responseObject.Data.Children[i].Data.Score)

		// db.Create((&sub))
		// Upsert
		storeService.db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&sub)

		getComments(sub)

	}
}

func build_indexes(name string) {

	var submissions []models.Submission

	var indexes []models.Index

	storeService.db.Find(&submissions)

	for i := 0; i < len(submissions); i++ {

		var index models.Index
		index.SubmissionID = submissions[i].SubmissionID
		index.CreatedUTC = submissions[i].CreatedUTC
		index.Score = submissions[i].Score

		indexes = append(indexes, index)
	}

	file, _ := json.MarshalIndent(indexes, "", " ")

	file_name := "./docs/api/" + name + "/indexes.json"

	write_err := os.WriteFile(file_name, file, 0777)

	if write_err != nil {
		panic(write_err)
	}

}

func create_end_points(name string) {

	var submissions []models.Submission

	storeService.db.Preload("Comments").Find(&submissions)

	if _, err := os.Stat("./docs/api/" + name); os.IsNotExist(err) {

		err := os.Mkdir("./docs/api/"+name, 0744)

		if err != nil {
			panic("Cant create directory")
		}
	}

	for i := 0; i < len(submissions); i++ {

		var id = ".docs/api/" + name + "/" + submissions[i].SubmissionID + ".json"

		file, _ := json.MarshalIndent(submissions[i], "", " ")

		err := os.WriteFile(id, file, 0744)

		if err != nil {
			panic(err)
		}

	}

}

func main() {

	subreddit_name := os.Args[1]

	no_of_post := os.Args[2]

	storeService.db = initialiseDatabase(subreddit_name)

	migrate()

	crawl(subreddit_name, no_of_post)

	create_end_points(subreddit_name)

	build_indexes(subreddit_name)

}
