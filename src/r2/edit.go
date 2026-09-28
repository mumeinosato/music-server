package r2

import (
	"context"
	"errors"
	"fmt"
	"music-server/src/config"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Upload は失敗した動画IDとエラーを返す
func Upload(files map[string]string) (failed []string, err error) {
	ctx := context.TODO()
	client := get_client()
	bucket_name := config.Get().R2BucketName

	var errs []error
	for id, file_path := range files {
		if err := upload_file(ctx, client, bucket_name, id+".opus", file_path); err != nil {
			failed = append(failed, id)
			errs = append(errs, err)
		}
	}
	return failed, errors.Join(errs...)
}

func upload_file(ctx context.Context, client *s3.Client, bucket_name string, key string, file_path string) error {
	file, err := os.Open(file_path)
	if err != nil {
		return fmt.Errorf("Cannot open file %s: %w", file_path, err)
	}
	defer file.Close()

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket_name),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("Failed to upload file %s: %w", file_path, err)
	}
	return nil
}

// Delete は失敗した動画IDとエラーを返す
func Delete(ids []string) (failed []string, err error) {
	ctx := context.TODO()
	client := get_client()
	bucket_name := config.Get().R2BucketName

	var errs []error
	for _, id := range ids {
		_, err = client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: aws.String(bucket_name),
			Key:    aws.String(id + ".opus"),
		})
		if err != nil {
			failed = append(failed, id)
			errs = append(errs, fmt.Errorf("Failed to delete file %s: %w", id+".opus", err))
		}
	}
	return failed, errors.Join(errs...)
}
