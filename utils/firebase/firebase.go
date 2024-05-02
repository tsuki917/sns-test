package fire

import (
	"context"
	"fmt"
	"io"
	"strconv"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
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

	defer wc.Close() // 関数の最後にファイルを閉じる
	if err != nil {
		return err
	}

	return nil
}

func ImagefileUP(c *gin.Context, fileLen int64, seq int64) {

	for i := 0; i < int(fileLen); i++ {
		filedata, header, err := c.Request.FormFile("image[" + strconv.Itoa(i) + "]")
		if err != nil {
			print(err)
			return
		}
		binaryData, err := io.ReadAll(filedata)
		if err != nil {
			print(err)
			return
		}
		path := "postImg/" + strconv.Itoa(int(seq)) + "/" + header.Filename
		SavePostImg(binaryData, path)

	}

}
