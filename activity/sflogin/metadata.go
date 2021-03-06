package sflogin

import "github.com/project-flogo/core/data/coerce"

type Settings struct {
	// BaseUrl: https://(login|test).salesforce.com
	BaseUrl      string `md:"baseUrl,required"`
	ClientId     string `md:"clientId,required"`
	ClientSecret string `md:"clientSecret,required"`
	Username     string `md:"username,required"`
	Password     string `md:"password,required"`
}

type Input struct {
}

func (r *Input) FromMap(values map[string]interface{}) error {
	return nil
}

func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{}
}

type Output struct {
	AccessToken string `md:"AccessToken"`
}

func (o *Output) FromMap(values map[string]interface{}) error {
	strVal, _ := coerce.ToString(values["AccessToken"])
	o.AccessToken = strVal
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"AccessToken": o.AccessToken,
	}
}
