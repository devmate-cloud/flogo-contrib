package s3uploader

import "github.com/project-flogo/core/data/coerce"

type Settings struct {
	AwsAccessKey    string `md:"awsAccessKey,required"`                 // AWS Credential Key
	AwsAccessSecret string `md:"awsAccessSecret,required"`              // AWS Credential Secret
	Region          string `md:"region,required" default:"us-region-1"` // AWS Region to upload
}

type Input struct {
	FileContent string `md:"fileContent,required"`
	Bucket      string `md:"bucket,required"`
	Key         string `md:"key,required"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"fileContent": i.FileContent,
		"bucket":      i.Bucket,
		"key":         i.Key,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {

	var err error
	i.FileContent, err = coerce.ToString(values["fileContent"])
	if err != nil {
		return err
	}

	i.Bucket, err = coerce.ToString(values["bucket"])
	if err != nil {
		return err
	}

	i.Key, err = coerce.ToString(values["key"])
	if err != nil {
		return err
	}

	return nil
}

type Output struct {
	UploadID  string `md:"uploadID"`  // The HTTP response data
	VersionID string `md:"versionID"` // The HTTP response data
	Location  string `md:"location"`  // The HTTP response data
}

func (o *Output) FromMap(values map[string]interface{}) error {
	var err error

	o.Location, err = coerce.ToString(values["location"])
	if err != nil {
		return err
	}
	o.UploadID, err = coerce.ToString(values["uploadID"])
	if err != nil {
		return err
	}

	o.VersionID, err = coerce.ToString(values["versionID"])
	if err != nil {
		return err
	}

	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"location":  o.Location,
		"versionID": o.VersionID,
		"uploadID":  o.UploadID,
	}
}
