package sample

import "github.com/project-flogo/core/data/coerce"

type Settings struct {
	// BaseUrl: https://(login|test).salesforce.com
	BaseUrl      string `md:"BaseUrl,required"`
	ClientId     string `md:"ClientId,required"`
	ClientSecret string `md:"ClientSecret,required"`
	Username     string `md:"Username,required"`
	Password     string `md:"Password,required"`
}

type Input struct {
	AnInput string `md:"anInput,required"`
}

func (r *Input) FromMap(values map[string]interface{}) error {
	strVal, _ := coerce.ToString(values["anInput"])
	r.AnInput = strVal
	return nil
}

func (r *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"anInput": r.AnInput,
	}
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
