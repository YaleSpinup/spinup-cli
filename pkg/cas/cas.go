package cas

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

var casUrl = "https://secure.its.yale.edu/cas"

// Auth logs into CAS
func Auth(username, password, service string, client *http.Client) error {

	svc, err := url.Parse(service)
	if err != nil {
		return err
	}

	loginUrl, err := url.Parse(casUrl + "/login")
	if err != nil {
		return err
	}

	loginUrl.Query().Add("service", svc.String())

	log.Debugf("logging into cas (service: %s) with loginUrl %s, user: %s", service, loginUrl, username)

	formValues := url.Values{}
	formValues.Set("username", username)
	formValues.Set("password", password)
	formValues.Set("service", service)

	req, err := http.NewRequest("POST", loginUrl.String(), strings.NewReader(formValues.Encode())) // URL-encoded payload
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	log.Debugf("logging into cas with loginUrl %s and form values %+v: request %+v", loginUrl, formValues, req)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	log.Info("successfully logged into cas")
	log.Debugf("got response from POST form %+v", res)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	log.Debugf("got response body: %s", string(body))

	return nil
}
