package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var uploadFormTmpl = []byte(`
<html>
<body>
<form action="/uploadCSV" method="post" enctype="multipart/form-data">
    File: <input type="file" name="csv_file">
    <input type="submit" value="UploadCSV">
</form>
<form action="/uploadCSVtemplate" method="post" enctype="multipart/form-data">
    File: <input type="file" name="csv_file">
    <input type="submit" value="uploadCSVtemplate">
</form>
<form action="/template" method="post" enctype="multipart/form-data">
    <input type="submit" value="template">
</form>
<i>Served by nginx</i>
</body>
</html>
`)

func mainPage(w http.ResponseWriter, r *http.Request) {
	w.Write(uploadFormTmpl)
}


type Params struct {
	ID   int
	User string
}


type ParamsCSV struct {
	ID   int
	User string
}

func uploadCSV(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("csv_file")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	var bytes []byte

	// read file content to byte buffer
	bytes, err = ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	// print file stats on page
	fmt.Fprintf(w, "<h3>Filename:</h3>%s\n", handler.Filename)
	fmt.Fprintf(w, "<br><h3>File content:</h3><xmp>\n%s</xmp>", string(bytes))


	// S3 section
	// Functions for content uploading and bucket items listing
	bucket := "test-k8s-kops-ha-storage"
	region := "eu-west-1"

	//select Region to use.
	conf := aws.Config{Region: aws.String(region)}
	uploaderSession := session.New(&conf)
	uploader := s3manager.NewUploader(uploaderSession)

	// file upload
	fmt.Println("Uploading file to S3...")

	// Setting file pointer to the beginning of the file
	// to enable reading it again for S3 upload
	file.Seek(0, io.SeekStart)
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(handler.Filename),
		Body:   file,
	})
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully uploaded to %s\n", result.Location)
	fmt.Fprintf(w, "<br><h3>Result:</h3>Successfully uploaded to <a href=%s>S3 bucket</a>\n", result.Location)

	// list previously uploaded files
	listerSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)

	// Create S3 service client
	lister := s3.New(listerSession)
	resp, err := lister.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(bucket)})
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	// print out the result
	fmt.Fprintf(w, "<br><h3>Uploaded files:</h3>\n")
	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)

		fmt.Fprintf(w, "Name:          %s\n", *item.Key)
		fmt.Fprintf(w, "<br>Last modified: %s\n", *item.LastModified)
		fmt.Fprintf(w, "<br>Size:          %v\n", *item.Size)
		fmt.Fprintf(w, "<br>Storage class: %s\n", *item.StorageClass)
		fmt.Fprintf(w, "<br><br>")
	}

}

func main() {
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/uploadCSV", uploadCSV)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
