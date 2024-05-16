package fire

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func SavePostImg(filedata []byte, filePath string) error {
	fmt.Println("Firebase setup")

	config := &firebase.Config{
		StorageBucket: "sns-image-storage.appspot.com",
	}
	ctx := context.Background()
	t := option.WithCredentialsFile("sns-image-storage.json")
	app, err := firebase.NewApp(ctx, config, t)
	if err != nil {
		fmt.Print(err)
		return err
	}
	client, err := app.Storage(ctx)
	if err != nil {
		fmt.Print(err)
		return err
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		fmt.Print(err)
		return err
	}
	wc := bucket.Object(filePath).NewWriter(ctx)

	_, err = wc.Write(filedata)

	if err != nil {
		defer wc.Close() // 失敗した場合もCloseを呼び出す
		return err
	}

	err = wc.Close() // 書き込みを終了し、ファイルを閉じる
	if err != nil {
		return err
	}

	// オブジェクトをパブリックに設定
	object := bucket.Object(filePath)
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return fmt.Errorf("error setting object ACL: %v", err)
	}

	fmt.Printf("File %s uploaded and set to public\n", filePath)
	return nil
}

// func GetPostImg(PostId int, FileName string) (baseURL string) {
// 	config := &firebase.Config{
// 		StorageBucket: "sns-image-storage.appspot.com",
// 	}
// 	ctx := context.Background()
// 	t := option.WithCredentialsFile("sns-image-storage.json")
// 	app, err := firebase.NewApp(ctx, config, t)
// 	if err != nil {
// 		fmt.Print(err)
// 		return
// 	}
// 	client, err := app.Storage(ctx)
// 	if err != nil {
// 		fmt.Print(err)
// 		return
// 	}

// 	bucket, err := client.DefaultBucket()
// 	if err != nil {
// 		fmt.Print(err)
// 		return
// 	}
// 	// object, err := bucket.Object("postImg/" + strconv.Itoa(PostId) + "/" + FileName).NewReader(ctx)
// 	objectName := "postImg/" + strconv.Itoa(PostId) + "/" + FileName
// 	object := bucket.Object(objectName)

// 	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	baseURL = fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucket.Name(), objectName)
// 	return baseURL
// 	// data, err := io.ReadAll(ob)
// 	// baseURL = base64.StdEncoding.EncodeToString(data)
// 	// if err != nil {
// 	// 	log.Fatalln(err)
// 	// 	return
// 	// }
// 	return
// }
