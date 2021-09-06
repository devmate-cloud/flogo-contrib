package s3uploader

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	act := &Activity{settings: s}

	return act, nil
}

type Activity struct {
	settings *Settings
}

func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Upload file to S3 via AWS Go SDK
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	logger := ctx.Logger()

	if logger.DebugEnabled() {
		logger.Debugf("(debug) Upload file to S3: bucket=%s, key=%s", input.Bucket, input.Key)
		logger.Debugf("AWS Credential: key=%s, secret=%s", a.settings.AwsAccessKey, a.settings.AwsAccessSecret)
	}
	logger.Infof("Upload file to S3: bucket=%s, key=%s", input.Bucket, input.Key)
	fmt.Println("Default region=", a.settings.Region)

	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.

	creds := credentials.NewStaticCredentials(a.settings.AwsAccessKey, a.settings.AwsAccessSecret, "")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(a.settings.Region),
		Credentials: creds,
	})

	if err != nil {
		logger.Errorf("Can not create AWS session")
		return false, err
	}

	// Setup the S3 Upload Manager. Also see the SDK doc for the Upload Manager
	// for more information on configuring part size, and concurrency.
	//
	// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#NewUploader
	uploader := s3manager.NewUploader(sess)

	// Upload the file's body to S3 bucket as an object with the key being the
	// same as the filename.
	resp, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(input.Bucket),

		// Can also use the `filepath` standard library package to modify the
		// filename as need for an S3 object key. Such as turning absolute path
		// to a relative path.
		Key: aws.String(input.Key),

		ACL: aws.String(a.settings.DefaultACL),

		// The file to be uploaded. io.ReadSeeker is preferred as the Uploader
		// will be able to optimize memory when uploading large content. io.Reader
		// is supported, but will require buffering of the reader's bytes for
		// each part.
		Body: bytes.NewBufferString(input.FileContent),
	})

	if err != nil {
		logger.Errorf("Failed to upload S3 with error %s", err)
		return false, err
	}

	output := &Output{Location: resp.Location, UploadID: resp.UploadID}
	err = ctx.SetOutputObject(output)
	if err != nil {
		return false, err
	}

	logger.Infof("Successfully uploaded %q to %q\n", input.Bucket, input.Key)

	return true, err
}
