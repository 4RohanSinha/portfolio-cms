package firebase_service

import (
	"cms/version_control"
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func App() (*firestore.Client, error) {
	opt := option.WithCredentialsFile("personal-portfolio-57012-firebase-adminsdk-cogpk-791995c08a.json")

	app, err := firebase.NewApp(context.Background(), &firebase.Config{
		ProjectID: "personal-portfolio-57012",
	}, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Firestore(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error initializing firestore client: %v", err)
	}

	return client, nil
}

func SetHead(client *firestore.Client) error {
	head, err := version_control.Head()

	if err != nil {
		return err
	}

	ctx := context.Background()
	colRef := client.Collection("posts")

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
		_, err = client.Collection("posts-test").Doc(k).Set(context.Background(), map[string]interface{}{
			"document": head.Content[k],
		})

		if err != nil {
			return err
		}
	}

	return nil
}
