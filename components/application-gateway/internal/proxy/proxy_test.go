package proxy

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kyma-project/kyma/components/application-gateway/internal/apperrors"
	csrfMock "github.com/kyma-project/kyma/components/application-gateway/internal/authorization/csrf/mocks"
	authMock "github.com/kyma-project/kyma/components/application-gateway/internal/authorization/mocks"
	"github.com/kyma-project/kyma/components/application-gateway/internal/httperrors"
	metadataMock "github.com/kyma-project/kyma/components/application-gateway/internal/metadata/mocks"
	metadatamodel "github.com/kyma-project/kyma/components/application-gateway/internal/metadata/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProxy(t *testing.T) {

	proxyTimeout := 10

	t.Run("should proxy and use internal cache", func(t *testing.T) {
		// given
		ts := NewTestServer(func(req *http.Request) {
			assert.Equal(t, req.Method, http.MethodGet)
			assert.Equal(t, req.RequestURI, "/orders/123")
		})
		defer ts.Close()

		req, err := http.NewRequest(http.MethodGet, "/orders/123", nil)
		require.NoError(t, err)

		req.Host = "app-test-uuid-1.namespace.svc.cluster.local"

		authStrategyMock := &authMock.Strategy{}
		authStrategyMock.
			On("AddAuthorization", mock.AnythingOfType("*http.Request"), mock.AnythingOfType("TransportSetter")).
			Return(nil).Twice()

		csrfTokenStrategyMock := &csrfMock.TokenStrategy{}
		csrfTokenStrategyMock.On("AddCSRFToken", mock.AnythingOfType("*http.Request")).
			Return(nil)

		credentials := &metadatamodel.Credentials{}
		authStrategyFactoryMock := &authMock.StrategyFactory{}
		authStrategyFactoryMock.On("Create", credentials).Return(authStrategyMock).Once()

		csrfTokenStrategyFactoryMock := &csrfMock.TokenStrategyFactory{}
		csrfTokenStrategyFactoryMock.On("Create", authStrategyMock, credentials.CSRFTokenURL).Return(csrfTokenStrategyMock).Once()

		serviceDefServiceMock := &metadataMock.ServiceDefinitionService{}
		serviceDefServiceMock.On("GetAPI", "uuid-1").Return(&metadatamodel.API{
			TargetUrl:   ts.URL,
			Credentials: credentials,
		}, nil).Once()

		handler := New(serviceDefServiceMock, authStrategyFactoryMock, csrfTokenStrategyFactoryMock, createProxyConfig(proxyTimeout))
		rr := httptest.NewRecorder()

		// when
		handler.ServeHTTP(rr, req)

		// then
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "test", rr.Body.String())

		// given
		nextReq, _ := http.NewRequest(http.MethodGet, "/orders/123", nil)
		nextReq.Host = "app-test-uuid-1.namespace.svc.cluster.local"
		rr = httptest.NewRecorder()

		//when
		handler.ServeHTTP(rr, nextReq)

		//then
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "test", rr.Body.String())
		authStrategyFactoryMock.AssertExpectations(t)
		authStrategyMock.AssertExpectations(t)
	})

	t.Run("should proxy OAuth calls", func(t *testing.T) {
		// given
		ts := NewTestServer(func(req *http.Request) {
			assert.Equal(t, req.Method, http.MethodGet)
			assert.Equal(t, req.RequestURI, "/orders/123")
		})
		defer ts.Close()

		tsOAuth := NewTestServer(func(req *http.Request) {
			assert.Equal(t, req.Method, http.MethodPost)
			assert.Equal(t, req.RequestURI, "/token")
		})
		defer ts.Close()

		req, err := http.NewRequest(http.MethodGet, "/orders/123", nil)
		require.NoError(t, err)

		req.Host = "app-test-uuid-1.namespace.svc.cluster.local"

		authStrategyMock := &authMock.Strategy{}
		authStrategyMock.
			On("AddAuthorization", mock.AnythingOfType("*http.Request"), mock.AnythingOfType("TransportSetter")).
			Return(nil)

		csrfTokenStrategyMock := &csrfMock.TokenStrategy{}
		csrfTokenStrategyMock.On("AddCSRFToken", mock.AnythingOfType("*http.Request")).
			Return(nil)

		credentialsMatcher := createOAuthCredentialsMatcher("clientId", "clientSecret", tsOAuth.URL+"/token")

		authStrategyFactoryMock := &authMock.StrategyFactory{}
		authStrategyFactoryMock.On("Create", mock.MatchedBy(credentialsMatcher)).Return(authStrategyMock)

		csrfTokenStrategyFactoryMock := &csrfMock.TokenStrategyFactory{}
		csrfTokenStrategyFactoryMock.On("Create", authStrategyMock, "").Return(csrfTokenStrategyMock).Once()

		serviceDefServiceMock := &metadataMock.ServiceDefinitionService{}
		serviceDefServiceMock.On("GetAPI", "uuid-1").Return(&metadatamodel.API{
			TargetUrl: ts.URL,
			Credentials: &metadatamodel.Credentials{
				OAuth: &metadatamodel.OAuth{
					ClientID:     "clientId",
					ClientSecret: "clientSecret",
					URL:          tsOAuth.URL + "/token",
				},
			},
		}, nil)

		handler := New(serviceDefServiceMock, authStrategyFactoryMock, csrfTokenStrategyFactoryMock, createProxyConfig(proxyTimeout))
		rr := httptest.NewRecorder()

		// when
		handler.ServeHTTP(rr, req)

		// then
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "test", rr.Body.String())
	})

	t.Run("should proxy BasicAuth auth calls", func(t *testing.T) {
		// given
		ts := NewTestServer(func(req *http.Request) {
			assert.Equal(t, req.Method, http.MethodGet)
			assert.Equal(t, req.RequestURI, "/orders/123")
		})
		defer ts.Close()

		req, err := http.NewRequest(http.MethodGet, "/orders/123", nil)
		require.NoError(t, err)

		req.Host = "app-test-uuid-1.namespace.svc.cluster.local"

		authStrategyMock := &authMock.Strategy{}
		authStrategyMock.
			On("AddAuthorization", mock.AnythingOfType("*http.Request"), mock.AnythingOfType("TransportSetter")).
			Return(nil)

		csrfTokenStrategyMock := &csrfMock.TokenStrategy{}
		csrfTokenStrategyMock.On("AddCSRFToken", mock.AnythingOfType("*http.Request")).
			Return(nil)

		credentialsMatcher := createBasicCredentialsMatcher("username", "password")

		authStrategyFactoryMock := &authMock.StrategyFactory{}
		authStrategyFactoryMock.On("Create", mock.MatchedBy(credentialsMatcher)).Return(authStrategyMock)

		csrfTokenStrategyFactoryMock := &csrfMock.TokenStrategyFactory{}
		csrfTokenStrategyFactoryMock.On("Create", authStrategyMock, "").Return(csrfTokenStrategyMock).Once()

		serviceDefServiceMock := &metadataMock.ServiceDefinitionService{}
		serviceDefServiceMock.On("GetAPI", "uuid-1").Return(&metadatamodel.API{
			TargetUrl: ts.URL,
			Credentials: &metadatamodel.Credentials{
				BasicAuth: &metadatamodel.BasicAuth{
					Username: "username",
					Password: "password",
				},
			},
		}, nil)

		handler := New(serviceDefServiceMock, authStrategyFactoryMock, csrfTokenStrategyFactoryMock, createProxyConfig(proxyTimeout))
		rr := httptest.NewRecorder()

		// when
		handler.ServeHTTP(rr, req)

		// then
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "test", rr.Body.String())
	})

	t.Run("should fail with Bad Gateway error when failed to get OAuth token", func(t *testing.T) {
		// given
		ts := NewTestServer(func(req *http.Request) {
			assert.Equal(t, req.Method, http.MethodGet)
			assert.Equal(t, req.RequestURI, "/orders/123")
		})
		defer ts.Close()

		req, err := http.NewRequest(http.MethodGet, "/orders/123", nil)
		require.NoError(t, err)

		req.Host = "app-test-uuid-1.namespace.svc.cluster.local"

		authStrategyMock := &authMock.Strategy{}
		authStrategyMock.
			On("AddAuthorization", mock.AnythingOfType("*http.Request"), mock.AnythingOfType("TransportSetter")).
			Return(apperrors.UpstreamServerCallFailed("failed"))

		csrfTokenStrategyMock := &csrfMock.TokenStrategy{}
		csrfTokenStrategyMock.On("AddCSRFToken", mock.AnythingOfType("*http.Request")).
			Return(nil)

		credentialsMatcher := createOAuthCredentialsMatcher("clientId", "clientSecret", "www.example.com/token")

		authStrategyFactoryMock := &authMock.StrategyFactory{}
		authStrategyFactoryMock.On("Create", mock.MatchedBy(credentialsMatcher)).Return(authStrategyMock)

		csrfTokenStrategyFactoryMock := &csrfMock.TokenStrategyFactory{}
		csrfTokenStrategyFactoryMock.On("Create", authStrategyMock, "").Return(csrfTokenStrategyMock).Once()

		serviceDefServiceMock := &metadataMock.ServiceDefinitionService{}
		serviceDefServiceMock.On("GetAPI", "uuid-1").Return(&metadatamodel.API{
			TargetUrl: ts.URL,
			Credentials: &metadatamodel.Credentials{
				OAuth: &metadatamodel.OAuth{
					ClientID:     "clientId",
					ClientSecret: "clientSecret",
					URL:          "www.example.com/token",
				},
			},
		}, nil)

		handler := New(serviceDefServiceMock, authStrategyFactoryMock, csrfTokenStrategyFactoryMock, createProxyConfig(proxyTimeout))
		rr := httptest.NewRecorder()

		// when
		handler.ServeHTTP(rr, req)

		// then
		assert.Equal(t, http.StatusBadGateway, rr.Code)
	})

	t.Run("should return 500 if failed to get service definition", func(t *testing.T) {
		// given
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		require.NoError(t, err)
		req.Host = "app-test-uuid-1.namespace.svc.cluster.local"
		rr := httptest.NewRecorder()

		serviceDefServiceMock := &metadataMock.ServiceDefinitionService{}
		serviceDefServiceMock.On("GetAPI", "uuid-1").
			Return(&metadatamodel.API{}, apperrors.Internal("Failed to read services"))

		handler := New(serviceDefServiceMock, nil, nil, createProxyConfig(proxyTimeout))

		// when
		handler.ServeHTTP(rr, req)

		// then
		var errorResponse httperrors.ErrorResponse

		json.Unmarshal([]byte(rr.Body.String()), &errorResponse)

		serviceDefServiceMock.AssertExpectations(t)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, http.StatusInternalServerError, errorResponse.Code)
	})

	t.Run("should invalidate proxy and retry when 401 occurred", func(t *testing.T) {
		// given
		tsf := NewTestServerForRetryTest(http.StatusUnauthorized, func(req *http.Request) {
			assert.Equal(t, req.Method, http.MethodGet)
			assert.Equal(t, req.RequestURI, "/orders/123")
		})
		defer tsf.Close()

		req, _ := http.NewRequest(http.MethodGet, "/orders/123", nil)
		req.Host = "app-test-uuid-1.namespace.svc.cluster.local"

		serviceDefServiceMock := &metadataMock.ServiceDefinitionService{}
		serviceDefServiceMock.On("GetAPI", "uuid-1").Return(&metadatamodel.API{
			TargetUrl:   tsf.URL,
			Credentials: &metadatamodel.Credentials{},
		}, nil).Twice()

		authStrategyMock := &authMock.Strategy{}
		authStrategyMock.
			On("AddAuthorization", mock.Anything, mock.AnythingOfType("TransportSetter")).
			Return(nil).Twice()
		authStrategyMock.On("Invalidate").Return().Once()

		csrfTokenStrategyMock := &csrfMock.TokenStrategy{}
		csrfTokenStrategyMock.On("AddCSRFToken", mock.AnythingOfType("*http.Request")).
			Return(nil)
		csrfTokenStrategyMock.On("Invalidate").Return()

		authStrategyFactoryMock := &authMock.StrategyFactory{}
		authStrategyFactoryMock.On("Create", mock.Anything).Return(authStrategyMock).Twice()

		csrfTokenStrategyFactoryMock := &csrfMock.TokenStrategyFactory{}
		csrfTokenStrategyFactoryMock.On("Create", authStrategyMock, "").Return(csrfTokenStrategyMock)

		handler := New(serviceDefServiceMock, authStrategyFactoryMock, csrfTokenStrategyFactoryMock, createProxyConfig(proxyTimeout))
		rr := httptest.NewRecorder()

		// when
		handler.ServeHTTP(rr, req)

		// then
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "test", rr.Body.String())

		serviceDefServiceMock.AssertExpectations(t)
		authStrategyFactoryMock.AssertExpectations(t)
		authStrategyMock.AssertExpectations(t)
	})
}

func TestInvalidStateHandler(t *testing.T) {
	t.Run("should always return Internal Server Error", func(t *testing.T) {
		// given
		req, err := http.NewRequest(http.MethodGet, "/test", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()

		handler := NewInvalidStateHandler("Application Gateway id not initialized properly")

		// when
		handler.ServeHTTP(rr, req)

		// then
		var errorResponse httperrors.ErrorResponse

		json.Unmarshal([]byte(rr.Body.String()), &errorResponse)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, http.StatusInternalServerError, errorResponse.Code)
	})
}

func NewTestServer(check func(req *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		check(r)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	}))
}

func NewTestServerForRetryTest(status int, check func(req *http.Request)) *httptest.Server {
	firstCall := true

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		check(r)
		if firstCall {
			w.WriteHeader(status)
			firstCall = false
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("test"))
		}
	}))
}

func createProxyConfig(proxyTimeout int) Config {
	return Config{
		SkipVerify:    true,
		ProxyTimeout:  proxyTimeout,
		Application:   "test",
		ProxyCacheTTL: proxyTimeout,
	}
}

func createOAuthCredentialsMatcher(clientId, clientSecret, url string) func(*metadatamodel.Credentials) bool {
	return func(c *metadatamodel.Credentials) bool {
		return c.OAuth != nil && c.OAuth.ClientID == clientId &&
			c.OAuth.ClientSecret == clientSecret &&
			c.OAuth.URL == url
	}
}

func createBasicCredentialsMatcher(username, password string) func(*metadatamodel.Credentials) bool {
	return func(c *metadatamodel.Credentials) bool {
		return c.BasicAuth != nil && c.BasicAuth.Username == username &&
			c.BasicAuth.Password == password
	}
}
