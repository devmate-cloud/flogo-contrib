package s3uploader

import (
	"fmt"
	"os"
	"testing"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {

	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}

func TestEval(t *testing.T) {

	var awsAccessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	var awsAccessSecret = os.Getenv("AWS_SECRET_ACCESS_KEY")
	var bucket = os.Getenv("TEST_BUCKET")
	var key = os.Getenv("TEST_KEY")

	settings := &Settings{AwsAccessKey: awsAccessKey, AwsAccessSecret: awsAccessSecret, Region: "us-east-1", DefaultACL: "public-read"}

	act := &Activity{settings: settings}
	tc := test.NewActivityContext(act.Metadata())
	input := &Input{Bucket: bucket, Key: key, FileContent: "Hello S3 from Flogo"}
	err := tc.SetInputObject(input)
	assert.Nil(t, err)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.Nil(t, err)
	assert.Equal(t, fmt.Sprintf("https://%s.s3.amazonaws.com/%s", input.Bucket, input.Key), output.Location)

	fmt.Println("Output file: ", output.Location)
}
