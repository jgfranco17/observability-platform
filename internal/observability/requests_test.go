package observability

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJar_Cookies_SetAndGet(t *testing.T) {
	j := &jar{}
	testURL, err := url.Parse("http://example.com")
	require.NoError(t, err)

	cookies := []*http.Cookie{
		{
			Name:    "session",
			Value:   "abc123",
			Expires: time.Now().Add(1 * time.Hour),
		},
		{
			Name:    "user",
			Value:   "john",
			Expires: time.Now().Add(1 * time.Hour),
		},
	}

	j.SetCookies(testURL, cookies)
	retrievedCookies := j.Cookies(testURL)

	assert.Len(t, retrievedCookies, 2)
	assert.Equal(t, "session", retrievedCookies[0].Name)
	assert.Equal(t, "abc123", retrievedCookies[0].Value)
	assert.Equal(t, "user", retrievedCookies[1].Name)
	assert.Equal(t, "john", retrievedCookies[1].Value)
}

func TestJar_Cookies_Expired(t *testing.T) {
	j := &jar{}
	testURL, err := url.Parse("http://example.com")
	require.NoError(t, err)

	cookies := []*http.Cookie{
		{
			Name:    "expired",
			Value:   "old",
			Expires: time.Now().Add(-1 * time.Hour),
		},
		{
			Name:    "valid",
			Value:   "new",
			Expires: time.Now().Add(1 * time.Hour),
		},
	}

	j.SetCookies(testURL, cookies)
	retrievedCookies := j.Cookies(testURL)

	assert.Len(t, retrievedCookies, 1)
	assert.Equal(t, "valid", retrievedCookies[0].Name)
	assert.Equal(t, "new", retrievedCookies[0].Value)
}

func TestJar_DifferentHosts(t *testing.T) {
	j := &jar{}
	url1, err := url.Parse("http://example.com")
	require.NoError(t, err)

	url2, err := url.Parse("http://different.com")
	require.NoError(t, err)

	cookies1 := []*http.Cookie{
		{
			Name:    "site1",
			Value:   "value1",
			Expires: time.Now().Add(1 * time.Hour),
		},
	}

	cookies2 := []*http.Cookie{
		{
			Name:    "site2",
			Value:   "value2",
			Expires: time.Now().Add(1 * time.Hour),
		},
	}

	j.SetCookies(url1, cookies1)
	j.SetCookies(url2, cookies2)

	retrievedCookies1 := j.Cookies(url1)
	retrievedCookies2 := j.Cookies(url2)

	assert.Len(t, retrievedCookies1, 1)
	assert.Equal(t, "site1", retrievedCookies1[0].Name)

	assert.Len(t, retrievedCookies2, 1)
	assert.Equal(t, "site2", retrievedCookies2[0].Name)
}

func TestJar_EmptyForUnknownHost(t *testing.T) {
	j := &jar{}
	testURL, err := url.Parse("http://unknown.com")
	require.NoError(t, err)

	retrievedCookies := j.Cookies(testURL)
	assert.Empty(t, retrievedCookies)
}
