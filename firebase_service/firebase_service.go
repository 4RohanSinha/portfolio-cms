package firebase_service

import (
	"cms/version_control"
	"context"
	"fmt"
	"io"
	"os"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func App() (*firebase.App, error) {

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "personal-portfolio-57012-firebase-adminsdk-cogpk-791995c08a.json")

	opt := option.WithCredentialsFile("personal-portfolio-57012-firebase-adminsdk-cogpk-791995c08a.json")

	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: "personal-portfolio-57012",
	}, opt)

	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}
	return app, nil
}

func FirestoreClient(app *firebase.App) (*firestore.Client, error) {

	client, err := app.Firestore(context.Background())

	if err != nil {
		return nil, err
	}

	return client, nil
}

func StorageClient(app *firebase.App) (*storage.Client, error) {
	client, err := storage.NewClient(context.Background())

	if err != nil {
		return nil, err
	}

	return client, nil
}

func UploadFile(fname string, app *firebase.App) error {
	client, err := StorageClient(app)

	if err != nil {
		return err
	}

	bucket := client.Bucket("personal-portfolio-57012.appspot.com")

	f, err := os.Open(".vc/content/" + fname)

	if err != nil {
		return err
	}

	defer f.Close()

	wc := bucket.Object(fname).NewWriter(context.Background())
	wc.ContentType = "text/markdown"

	if _, err := io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil

}

func SetHead(app *firebase.App) error {

	f_client, err := FirestoreClient(app)

	if err != nil {
		return err
	}

	head, err := version_control.Head()

	if err != nil {
		return err
	}

	ctx := context.Background()
	colRef := f_client.Collection("posts")

	docs, err := colRef.Documents(ctx).GetAll()

	if err != nil {
		return err
	}

	for _, doc := range docs {
		_, err := doc.Ref.Delete(ctx)

		if err != nil {
			return err
		}
	}

	for k := range head.Content {
		doc_id := k

		if len(k) > 3 {
			doc_id = k[:len(k)-3]
		}

		err := UploadFile(k, app)

		if err != nil {
			panic(err)
		}

		_, err = f_client.Collection("posts-dev").Doc(k).Set(context.Background(), map[string]interface{}{
			"document": k,
			"id":       doc_id,
			"title":    "New post",
		})

		if err != nil {
			return err
		}
	}

	return nil
}
