package authprofile

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/RafayLabs/rcloud-cli/pkg/log"
	"github.com/RafayLabs/rcloud-cli/pkg/versioninfo"
	"github.com/ghodss/yaml"
	"github.com/levigross/grequests"
)

type HttpResp struct {
	StatusCode int       `json:"status_code"`
	Details    []Details `json:"details"`
}
type Details struct {
	ErrorCode string `json:"error_code"`
	Detail    string `json:"detail"`
	Info      string `json:"info"`
}

// error types
var (
	ResourceNotExists   = errors.New("resource does not exist")
	OperationNotAllowed = errors.New("operation not allowed")
	InvalidCredentials  = errors.New("invalid credentials")
)

func getSession(skipServerCertCheck bool) *grequests.Session {
	var sessionRequestOption *grequests.RequestOptions
	if skipServerCertCheck {
		sessionRequestOption = &grequests.RequestOptions{
			InsecureSkipVerify: true,
		}
	}
	return grequests.NewSession(sessionRequestOption)
}

// Common interface for the subclasses, key, session and null.
type AuthIf interface {
	Auth(s *grequests.Session) (map[string]string, error)
	SendRequest(s *grequests.Session, uri, method string, ro *grequests.RequestOptions) (*grequests.Response, error)
}

func (p *Profile) SubProfile() AuthIf {
	return &KeyProfile{Name: p.Name, URL: p.URL, Key: p.Key, Secret: p.Secret}
}

func (p *Profile) Auth(s *grequests.Session) (map[string]string, error) {
	sub := p.SubProfile()

	return sub.Auth(s)
}

func (p *Profile) SendRequest(s *grequests.Session, uri, method string, ro *grequests.RequestOptions) (*grequests.Response, error) {
	sub := p.SubProfile()

	version := versioninfo.Get()

	ro.Headers["cli-build-number"] = version.Version
	ro.Headers["cli-arch"] = version.Arch
	ro.Headers["cli-build-time"] = version.Time

	log.GetLogger().Debugf("Send request %s", uri)
	return sub.SendRequest(s, uri, method, ro)
}

func (p *Profile) AuthAndRequest(uri, method string, payload interface{}) (string, error) {
	response, err := p.AuthAndRequestFullResponse(uri, method, payload)
	if err != nil {
		return "", err
	}
	return response.String(), nil
}

func (p *Profile) AuthAndRequestFullResponse(uri, method string, payload interface{}) (*grequests.Response, error) {
	return p.AuthAndRequestWithHeadersFullResponse(uri, method, payload, nil)
}

func (p *Profile) AuthAndRequestWithHeadersFullResponse(uri, method string, payload interface{}, additionHeaders map[string]string) (*grequests.Response, error) {
	s := getSession(p.SkipServerCertValid)
	sub := p.SubProfile()

	headers, err := sub.Auth(s)
	if err != nil {
		return nil, err
	}
	headers["Content-Type"] = "application/json"
	if additionHeaders != nil {
		for h, v := range additionHeaders {
			headers[h] = v
		}
	}
	ro := &grequests.RequestOptions{
		Headers:   headers,
		JSON:      payload,
		UserAgent: getUserAgent(),
	}
	response, err := sub.SendRequest(s, uri, method, ro)
	if err != nil {
		log.GetLogger().Debugf("Error in response from core: %s\n", err)
		return response, fmt.Errorf("error in response from core: %s\n", err)
	}

	if response != nil && !response.Ok {
		log.GetLogger().Debugf("response not ok: %s", response.String())

		// check if error type is permission issue
		if strings.Contains(response.String(), "no or invalid credentials") {
			return response, InvalidCredentials
		}
		// check if error type is resource not found
		if strings.Contains(response.String(), "pg: no rows in result set") {
			return response, ResourceNotExists
		}
		// check if error type is permission issue
		if strings.Contains(response.String(), "method or route not allowed") {
			return response, OperationNotAllowed
		}
		// check if error type is permission issue
		if strings.Contains(response.String(), "You do not have enough privileges") {
			return response, OperationNotAllowed
		}

		var h HttpResp
		err = json.Unmarshal([]byte(response.String()), &h)
		if err == nil {
			log.GetLogger().Debugf("Struct values: %v", h)
			if len(h.Details) > 0 {
				log.GetLogger().Debugf("Returning error: server error: %s", h.Details[0].Detail)
				return response, fmt.Errorf("server error: %s", h.Details[0].Detail)
			}
		} else {
			log.GetLogger().Debugf("Response: %q", response.String())
		}
		return response, fmt.Errorf("server error [return code: %d]: %s", response.StatusCode, response.String())
	}

	return response, nil
}

func getUserAgent() string {
	version := versioninfo.Get()
	return fmt.Sprintf("RCTL/%s %s", version.Version, version.Arch)
}

func (p *Profile) AuthAndPostMultipartFile(uri, fileLocation, fileKeyName string) (string, error) {
	s := getSession(p.SkipServerCertValid)
	sub := p.SubProfile()
	headers, err := sub.Auth(s)
	if err != nil {
		return "", err
	}
	file, err := os.Open(fileLocation)
	if err != nil {
		log.GetLogger().Debugf("error in reading the helm file: %s\n", err)
		return "", fmt.Errorf("error in reading the helm file: %s", err)
	}
	ro := &grequests.RequestOptions{
		Headers:   headers,
		UserAgent: getUserAgent(),
		Files: []grequests.FileUpload{
			{
				FieldName:    fileKeyName,
				FileName:     filepath.Base(fileLocation),
				FileContents: file,
			},
		},
	}
	response, err := sub.SendRequest(s, uri, "POST", ro)
	if err != nil {
		log.GetLogger().Debugf("error in response from core: %s\n", err)
		return "", fmt.Errorf("error in response from core: %s", err)
	}

	if !response.Ok {
		log.GetLogger().Debugf("response not OK. [Status: %d] %s", response.StatusCode, response.String())
		return "", fmt.Errorf("server error [return code: %d]: %s", response.StatusCode, response.String())
	}

	return response.String(), nil
}

func (p *Profile) AuthAndPostMultipartFiles(uri string, files, fileKeys []string, valuesPath string) (string, error) {
	if len(files) == 0 {
		return "", fmt.Errorf("no files passed as argument to upload")
	}
	if len(files) != len(fileKeys) {
		return "", fmt.Errorf("number of file keys and files don't match")
	}
	s := getSession(p.SkipServerCertValid)
	sub := p.SubProfile()

	headers, err := sub.Auth(s)
	if err != nil {
		return "", err
	}
	fileUploads := make([]grequests.FileUpload, 0)

	for index, fileLocation := range files {
		file, err := os.Open(fileLocation)
		if err != nil {
			log.GetLogger().Debugf("error in reading the helm file: %s\n", err)
			return "", fmt.Errorf("error in reading the helm file: %s\n", err)
		}

		fileUploads = append(fileUploads, grequests.FileUpload{
			FieldName:    fileKeys[index],
			FileName:     filepath.Base(fileLocation),
			FileContents: file,
		})
	}
	ro := &grequests.RequestOptions{
		Headers:   headers,
		UserAgent: getUserAgent(),
		Files:     fileUploads,
		Data:      map[string]string{"valuesfileorder": valuesPath},
	}
	response, err := sub.SendRequest(s, uri, "POST", ro)
	if err != nil {
		log.GetLogger().Debugf("error in response from core: %s\n", err)
		return "", fmt.Errorf("error in response from core: %s\n", err)
	}

	if !response.Ok {
		log.GetLogger().Debugf("response not OK. [Status: %d] %s", response.StatusCode, response.String())
		return "", fmt.Errorf("server error [return code: %d]: %s", response.StatusCode, response.String())
	}

	return response.String(), nil
}

func (p *Profile) Request(uri, method string, payload interface{}, yamlResp bool) (string, error) {
	s := getSession(p.SkipServerCertValid)
	sub := p.SubProfile()

	headers, err := sub.Auth(s)
	if err != nil {
		return "", err
	}

	if yamlResp {
		headers["Content-Type"] = "application/yaml"
	} else {
		headers["Content-Type"] = "application/json"
	}

	ro := &grequests.RequestOptions{Headers: headers, JSON: payload}

	response, err := sub.SendRequest(s, uri, method, ro)
	if err != nil {
		log.GetLogger().Debugf("error in response from core: %s\n", err)
		return "", fmt.Errorf("error in response from core: %s\n", err)
	}

	if !response.Ok {
		log.GetLogger().Debugf("response not OK. [Status: %d] %s", response.StatusCode, response.String())
		return "", fmt.Errorf("server error [return code: %d]: %s", response.StatusCode, response.String())
	}

	return response.String(), nil
}

func (p *Profile) AuthAndPostData(uri, method, contentType string, params map[string]string, data []byte) (string, error) {
	s := getSession(p.SkipServerCertValid)
	sub := p.SubProfile()
	headers, err := sub.Auth(s)
	if err != nil {
		return "", err
	}
	headers["Content-Type"] = contentType

	ro := &grequests.RequestOptions{
		Params:      params,
		Headers:     headers,
		RequestBody: bytes.NewReader(data),
	}
	response, err := sub.SendRequest(s, uri, method, ro)
	if err != nil {
		log.GetLogger().Debugf("error in response from core: %s\n", err)
		return "", fmt.Errorf("error in response from core: %s\n", err)
	}

	if !response.Ok {
		log.GetLogger().Debugf("response not OK. [Status: %d] %s", response.StatusCode, response.String())
		return "", fmt.Errorf("server error [return code: %d]: %s", response.StatusCode, response.String())
	}

	return response.String(), nil
}

func (p *Profile) PostRequestFromFile(uri, filename string) (string, error) {
	return p.RequestFromFile(uri, "POST", filename)
}

func (p *Profile) RequestFromFile(uri, method, filename string) (string, error) {
	var payload interface{}

	r, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer r.Close()

	ext := filepath.Ext(filename)
	if ext == ".yml" || ext == ".yaml" {
		y, err := ioutil.ReadAll(r)
		if err != nil {
			return "", err
		}

		err = yaml.Unmarshal(y, &payload)
		if err != nil {
			return "", err
		}
	} else {
		j, err := ioutil.ReadAll(r)
		if err != nil {
			return "", err
		}

		err = json.Unmarshal(j, &payload)
		if err != nil {
			return "", err
		}
	}

	return p.AuthAndRequest(uri, method, payload)
}
