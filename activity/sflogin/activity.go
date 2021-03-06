package sflogin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/oliveagle/jsonpath"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New) //activity.Register(&Activity{}, New) to create instances using factory method 'New'
}

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

//New optional factory method, should be used if one activity instance per configuration is desired
func New(ctx activity.InitContext) (activity.Activity, error) {

	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	ctx.Logger().Debugf("Setting: %s", s.BaseUrl)

	fmt.Println("Setting baseUrl: ", s.BaseUrl)

	act := &Activity{
		BaseUrl: s.BaseUrl, ClientId: s.ClientId, ClientSecret: s.ClientSecret, Username: s.Username, Password: s.Password} //add aSetting to instance

	return act, nil
}

// Activity is an sflogin Activity that can be used as a base to create a custom activity
type Activity struct {
	BaseUrl      string
	ClientId     string
	ClientSecret string
	Username     string
	Password     string
}

// Metadata returns the activity's metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements api.Activity.Eval - Logs the Message
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	fmt.Println("Input baseUrl: ", a.BaseUrl)

	url := a.BaseUrl + "/services/oauth2/token"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("username", a.Username)
	_ = writer.WriteField("password", a.Password)
	_ = writer.WriteField("grant_type", "password")
	_ = writer.WriteField("client_id", a.ClientId)
	_ = writer.WriteField("client_secret", a.ClientSecret)

	writer.Close()

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "multipart/form-data;")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Response from Salesforce: ", string(body))

	var json_data interface{}

	json.Unmarshal([]byte(body), &json_data)

	val, err := jsonpath.JsonPathLookup(json_data, "$.access_token")

	if str, ok := val.(string); ok {

		/* act on str */
		output := &Output{AccessToken: str}
		err = ctx.SetOutputObject(output)
		if err != nil {
			return true, err
		}

		fmt.Println("==> Access Token: ", output.AccessToken)

		return true, nil
	} else {
		return true, errors.New("math: square root of negative number")
	}

}
