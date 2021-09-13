package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	openapi3_routers "github.com/getkin/kin-openapi/routers"
	openapi3_legacy "github.com/getkin/kin-openapi/routers/legacy"
	"github.com/stretchr/testify/suite"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

//go:embed api.yaml
var apiSpec []byte

var ctx = context.Background()

func TestAPI(t *testing.T) {
	suite.Run(t, &APISuite{})
}

type APISuite struct {
	suite.Suite

	client        http.Client
	apiSpecRouter openapi3_routers.Router
}

func (s *APISuite) SetupSuite() {
	srv := NewServer()
	go func() {
		log.Printf("Start serving on %s", srv.Addr)
		log.Fatal(srv.ListenAndServe())
	}()

	spec, err := openapi3.NewLoader().LoadFromData(apiSpec)
	s.Require().NoError(err)
	s.Require().NoError(spec.Validate(ctx))
	router, err := openapi3_legacy.NewRouter(spec)
	s.Require().NoError(err)
	s.apiSpecRouter = router
	s.client.Transport = s.specValidating(http.DefaultTransport)
}

func (s *APISuite) TestNotFound() {
	// when:
	resp, err := s.client.Get("http://localhost:8080/bibab")

	// then:
	s.Require().NoError(err)
	s.Require().Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *APISuite) TestCreateAndGet() {
	// setup:
	targetContent := []byte("biba kuka")
	testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		_, err := rw.Write(targetContent)
		s.Require().NoError(err)
	}))

	var key string
	s.Run("CheckRedirectResponse", func() {
		// when:
		reqBody := io.NopCloser(strings.NewReader(fmt.Sprintf( /* language=json */ `{"url": "%s"}`, testServer.URL)))
		resp, err := s.client.Post("http://localhost:8080/api/urls", "application/json", reqBody)

		// then:
		s.Require().NoError(err)
		rawBody, err := ioutil.ReadAll(resp.Body)
		s.Require().NoError(err)
		var body map[string]string
		s.Require().NoError(json.Unmarshal(rawBody, &body))
		key = body["key"]
		s.Require().NotEmpty(key)
	})

	s.Run("CheckRedirectResponse", func() {
		// setup:
		s.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
		defer func() {
			s.client.CheckRedirect = nil
		}()

		// when:
		resp, err := s.client.Get(fmt.Sprintf("http://localhost:8080/%s", key))

		// then:
		s.Require().NoError(err)
		s.Require().Equal(http.StatusPermanentRedirect, resp.StatusCode)
	})

	s.Run("CheckFollowingRedirect", func() {
		// setup:
		var client http.Client // don't validate against spec and follow redirects

		// when:
		resp, err := client.Get(fmt.Sprintf("http://localhost:8080/%s", key))

		// then:
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		s.Require().NoError(err)
		s.Require().Equal(targetContent, body)
	})
}

func (s *APISuite) specValidating(transport http.RoundTripper) http.RoundTripper {
	return RoundTripperFunc(func(req *http.Request) (*http.Response, error) {
		reqBody := s.readAll(req.Body)

		// validate request
		route, params, err := s.apiSpecRouter.FindRoute(req)
		s.Require().NoError(err)
		reqDescriptor := &openapi3filter.RequestValidationInput{
			Request:     req,
			PathParams:  params,
			QueryParams: req.URL.Query(),
			Route:       route,
		}
		s.Require().NoError(openapi3filter.ValidateRequest(ctx, reqDescriptor))

		// do request
		req.Body = io.NopCloser(bytes.NewReader(reqBody))
		resp, err := transport.RoundTrip(req)
		if err != nil {
			return nil, err
		}
		respBody := s.readAll(resp.Body)

		// Validate response against OpenAPI spec
		s.Require().NoError(openapi3filter.ValidateResponse(ctx, &openapi3filter.ResponseValidationInput{
			RequestValidationInput: reqDescriptor,
			Status:                 resp.StatusCode,
			Header:                 resp.Header,
			Body:                   io.NopCloser(bytes.NewReader(respBody)),
		}))

		return resp, nil
	})
}

func (s *APISuite) readAll(in io.Reader) []byte {
	if in == nil {
		return nil
	}
	data, err := ioutil.ReadAll(in)
	s.Require().NoError(err)
	return data
}

type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (fn RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}
