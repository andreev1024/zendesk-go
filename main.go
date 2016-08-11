package zendesk

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const ERR = "Error! Please, check status and body to get error details"

type API struct {
	email string
	token string
	host  string

	errorHandler func(err error)
}

func NewAPI(email, token, host string, errorHandler ...func(err error)) *API {
	var eh func(err error)
	if len(errorHandler) > 0 {
		eh = errorHandler[0]
	}

	return &API{
		token:        token,
		email:        email,
		host:         host,
		errorHandler: eh,
	}
}

func (a *API) Send(method, url string, reqData []byte) (body []byte, httpResp *http.Response, err error) {
	url = a.prepareUrl(url)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqData))
	if err != nil {
		a.HandleError(err)
		return
	}
	req.Header.Add("Accept", "application/json")
	if method == http.MethodPost || method == http.MethodPut {
		req.Header.Add("Content-Type", "application/json")
	}
	body, httpResp, err = a.sendRequest(req)
	return
}

func (a *API) SendFile(method string, url string, paramName, path string) (body []byte, httpResp *http.Response, err error) {
	url = a.prepareUrl(url)

	file, err := os.Open(path)
	if err != nil {
		a.HandleError(err)
		return
	}
	defer file.Close()

	b := &bytes.Buffer{}
	writer := multipart.NewWriter(b)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		a.HandleError(err)
		return
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		a.HandleError(err)
		return
	}

	req, err := http.NewRequest(method, url, b)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	body, httpResp, err = a.sendRequest(req)
	return
}

func (a *API) sendRequest(req *http.Request) (body []byte, httpResp *http.Response, err error) {
	a.useTokenAuth(req)

	client := http.DefaultClient
	httpResp, err = client.Do(req)

	if httpResp != nil {
		defer httpResp.Body.Close()
	}

	if err != nil {
		a.HandleError(err)
		return
	}

	body, err = ioutil.ReadAll(httpResp.Body)
	if err != nil {
		a.HandleError(err)
		return
	}

	match, err := regexp.MatchString("^[2|3].+$", httpResp.Status)
	if err != nil {
		a.HandleError(err)
		return
	}

	if !match {
		err = fmt.Errorf(ERR)
		return
	}

	return
}

func (a *API) useTokenAuth(req *http.Request) {
	req.SetBasicAuth(a.email+"/token", a.token)
}

func (a *API) prepareUrl(u string) string {
	return strings.Join([]string{a.host, "api/v2", u}, "/")
}

func (a *API) HandleError(err error) {
	if a.errorHandler != nil {
		a.errorHandler(err)
		return
	}
}