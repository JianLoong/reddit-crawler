// Reddit Crawler
package main

import (
	"crawlers/models"
	"io"
	"net/http"

	"encoding/json"
	"fmt"
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

	sqlite_path := fmt.Sprintf("databases/%v.sqlite", name)

	db, err := gorm.Open(sqlite.Open(sqlite_path), &gorm.Config{})
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

	// client := http.Client{}

	// request, err := http.NewRequest(http.MethodGet, "https://oauth.reddit.com/api/v1/scopes", nil)

	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }

	// response, err := client.Do(request)

	// respData, err := io.ReadAll(response.Body)

	// if err != nil {
	// 	fmt.Println(err)
	// 	panic(err)
	// }

	// fmt.Println(respData)

	// defer response.Body.Close()

	// return respData

	response, err := http.Get(url)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer response.Body.Close()

	return responseData

}

func getComments(sub models.Submission) {

	var commentObject models.CommentsResponse

	urlString := fmt.Sprintf("http://www.reddit.com/%v.json", sub.Permalink)

	var commentData = request(urlString)

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
		comment.SubmissionID = sub.SubmissionID
		comment.CommentID = commentObject[1].Data.Children[j].Data.ID

		comments = append(comments, comment)

	}
	// Batch insert comments
	storeService.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&comments)

}

func crawl(subreddit_name string, no_of_post string) {

	urlString := fmt.Sprintf("https://www.reddit.com/r/%v/.json?limit=%v", subreddit_name, no_of_post)

	var responseObject models.SubmissionsResponse

	var responseData = request(urlString)

	err := json.Unmarshal(responseData, &responseObject)

	fmt.Println(urlString)

	if err != nil {
		panic(err)
	}

	fmt.Println("Processing")

	var submissions []models.Submission

	for i := 0; i < len(responseObject.Data.Children); i++ {
		var sub models.Submission

		sub.SubmissionID = responseObject.Data.Children[i].Data.ID
		sub.Url = responseObject.Data.Children[i].Data.URL
		sub.Title = responseObject.Data.Children[i].Data.Title
		sub.CreatedUTC = responseObject.Data.Children[i].Data.CreatedUtc
		sub.Selftext = responseObject.Data.Children[i].Data.Selftext
		sub.Permalink = responseObject.Data.Children[i].Data.Permalink
		sub.Score = uint8(responseObject.Data.Children[i].Data.Score)

		//db.Create((&sub))
		//Upsert
		storeService.db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&sub)

		submissions = append(submissions, sub)

		getComments(sub)

	}

	//storeService.db.Clauses(clause.OnConflict{
	//	UpdateAll: true,
	//}).Create(&submissions)

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

	file_name := fmt.Sprintf("docs/api/%v/indexes.json", name)

	write_err := os.WriteFile(file_name, file, 0777)

	if write_err != nil {
		panic(write_err)
	}

}

func create_end_points(name string) {

	var submissions []models.Submission

	storeService.db.Preload("Comments").Find(&submissions)

	file_path := fmt.Sprintf("docs/api/%v", name)

	if _, err := os.Stat(file_path); os.IsNotExist(err) {

		err := os.Mkdir(file_path, 0744)

		if err != nil {
			panic("Cant create directory")
		}
	}

	for i := 0; i < len(submissions); i++ {

		var file_name = fmt.Sprintf("%v/%v.json", file_path, submissions[i].SubmissionID)

		file, _ := json.MarshalIndent(submissions[i], "", " ")

		err := os.WriteFile(file_name, file, 0744)

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
