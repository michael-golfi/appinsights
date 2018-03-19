package insights

import (
	"testing"
	"gitlab.com/michael.golfi/appinsights/constants"
	"github.com/stretchr/testify/require"
	"net/url"
)

func TestParseURL(t *testing.T) {
	var (
		endpoint string
		uri      *url.URL
		err      error
	)

	endpoint = constants.Endpoint
	uri, err = parseURL(endpoint)
	require.NoError(t, err)
	require.NotNil(t, uri)
	require.Equal(t, "/v2/track", uri.Path)
	require.Equal(t, "https", uri.Scheme)
	require.Equal(t, "dc.services.visualstudio.com", uri.Hostname())

	endpoints := []string{
		"not a url",
		"https://not a url",
		"https://dc.services.visualstudio.com",
		"https://dc.services.visualstudio.com/",
		"https://dc.services.visualstudio.com/?/",
		"https://dc.services.visualstudio.com/#/",
	}

	for _, endpoint = range endpoints {
		uri, err = parseURL(endpoint)
		require.Error(t, err)
		require.Nil(t, uri)
	}
}

func TestVerifyInsightsConnection(t *testing.T) {
	var (
		err error
	)

	allValues := []string{
		"",
		"https://somemalformedurl",
		"https://some malformed url",
		"http://getstatuscode.com/403",
	}

	err = verifyInsightsConnection(constants.Endpoint)
	require.NoError(t, err)

	for _, val := range allValues {
		err = verifyInsightsConnection(val)
		require.Error(t, err)
	}
}
