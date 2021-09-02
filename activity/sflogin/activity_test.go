package sflogin

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

	var baseUrl = os.Getenv("SF_BASE_URL")
	var clientId = os.Getenv("SF_CLIENT_ID")
	var clientSecret = os.Getenv("SF_CLIENT_SECRET")
	var username = os.Getenv("SF_USERNAME")
	var password = os.Getenv("SF_PASSWORD")

	act := &Activity{BaseUrl: baseUrl, ClientId: clientId, ClientSecret: clientSecret, Username: username, Password: password}
	tc := test.NewActivityContext(act.Metadata())
	input := &Input{AnInput: "test"}
	err := tc.SetInputObject(input)
	assert.Nil(t, err)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.Nil(t, err)
	assert.NotNil(t, output.AccessToken)

	fmt.Printf("Access token: %s", output.AccessToken)
	fmt.Println()
}
