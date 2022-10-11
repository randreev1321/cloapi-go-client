package servers

import (
	"context"
	"fmt"
	"github.com/clo-ru/cloapi-go-client/clo"
	"github.com/clo-ru/cloapi-go-client/clo/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"sync"
	"testing"
)

func TestServerListRequest_BuildRequest(t *testing.T) {
	ID := "id"
	req := ServerListRequest{
		ProjectID: ID,
	}
	h := http.Header{}
	h.Add("Authorization", fmt.Sprintf("Bearer %s", mocks.MockAuthKey))
	h.Add("Content-type", "application/json")
	h.Add("X-Add-Some", "SomeHeaderValue")
	rawReq, e := req.buildRequest(context.Background(), map[string]interface{}{
		"auth_key": mocks.MockAuthKey,
		"base_url": mocks.MockUrl,
	})
	rawReq.Header = h
	assert.Nil(t, e)
	expReq, _ := http.NewRequestWithContext(
		context.Background(), http.MethodGet, mocks.MockUrl+fmt.Sprintf(serverListEndpoint, ID), nil,
	)
	expReq.Header = h
	assert.Equal(t, expReq, rawReq)
}

func TestServerListRequest_Make(t *testing.T) {
	httpCli := mocks.MockClient{}
	cli := clo.ApiClient{
		HttpClient: &httpCli,
		Options: map[string]interface{}{
			"auth_key": "secret",
			"base_url": "https://clo.ru",
		},
	}
	var cases = []struct {
		Name           string
		ShouldFail     bool
		Req            ServerListRequest
		BodyStringFunc func() (string, int)
		Expected       ServerListResponse
	}{
		{
			Name: "Success",
			BodyStringFunc: func() (string, int) {
				return `{"count": 2, "results": [{"id": "first_item_id", "name": "host.com", "flavor":{"ram":3,"vcpus":2}},{"id": "second_item_id", "name": "host.com", "flavor":{"ram":2,"vcpus":3}}]}`,
					http.StatusOK
			},
			Req: ServerListRequest{
				ProjectID: "project_id",
			},
			Expected: ServerListResponse{
				Count: 2,
				Results: []ServerListItem{
					{
						ID:   "first_item_id",
						Name: "host.com",
						Flavor: ServerFlavorData{
							Ram:   3,
							Vcpus: 2,
						},
					},
					{
						ID:   "second_item_id",
						Name: "host.com",
						Flavor: ServerFlavorData{
							Ram:   2,
							Vcpus: 3,
						},
					},
				},
			},
		},
		{
			Name:       "WrongCount",
			ShouldFail: true,
			BodyStringFunc: func() (string, int) {
				return `{"count": 2, "results": [{"id": "first_item_id", "name": "host.com", "flavor":{"ram":3,"vcpus":2}},{"id": "second_item_id", "name": "host.com", "flavor":{"ram":2,"vcpus":3}}]}`,
					http.StatusOK
			},
			Req: ServerListRequest{
				ProjectID: "project_id",
			},
			Expected: ServerListResponse{
				Count: 1,
				Results: []ServerListItem{
					{
						ID:   "first_item_id",
						Name: "host.com",
						Flavor: ServerFlavorData{
							Ram:   3,
							Vcpus: 2,
						},
					},
					{
						ID:   "second_item_id",
						Name: "host.com",
						Flavor: ServerFlavorData{
							Ram:   2,
							Vcpus: 3,
						},
					},
				},
			},
		},
		{
			Name:       "WrongRamReturned",
			ShouldFail: true,
			BodyStringFunc: func() (string, int) {
				return `{"count": 2, "results": [{"id": "first_item_id", "name": "host.com", "flavor":{"ram":2,"vcpus":2}},{"id": "second_item_id", "name": "host.com", "flavor":{"ram":2,"vcpus":3}}]}`,
					http.StatusOK
			},
			Req: ServerListRequest{
				ProjectID: "project_id",
			},
			Expected: ServerListResponse{
				Count: 2,
				Results: []ServerListItem{
					{
						ID:   "first_item_id",
						Name: "host.com",
						Flavor: ServerFlavorData{
							Ram:   3,
							Vcpus: 2,
						},
					},
					{
						ID:   "second_item_id",
						Name: "host.com",
						Flavor: ServerFlavorData{
							Ram:   2,
							Vcpus: 3,
						},
					},
				},
			},
		},
		{
			Name:       "Error",
			ShouldFail: true,
			BodyStringFunc: func() (string, int) {
				return "", http.StatusInternalServerError
			},
			Req: ServerListRequest{
				ProjectID: "project_id",
			},
			Expected: ServerListResponse{Results: []ServerListItem{}},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			mocks.BodyStringFunc = c.BodyStringFunc
			res, e := c.Req.Make(context.Background(), &cli)
			if !c.ShouldFail {
				assert.Nil(t, e)
				assert.Equal(t, c.Expected, res)
			} else {
				assert.NotEqual(t, c.Expected, res)
			}
		})
	}
}

func TestServerListRequest_MakeRetry(t *testing.T) {
	retry := 5
	erCode := http.StatusInternalServerError
	httpCli := mocks.RequestDebugClient{}
	cli := clo.ApiClient{
		HttpClient: &httpCli,
		Options: map[string]interface{}{
			"auth_key": "secret",
			"base_url": "https://clo.ru",
		},
	}
	mocks.BodyStringFunc = func() (string, int) {
		return "", erCode
	}
	grNum := 4
	wg := sync.WaitGroup{}
	for n := 0; n < grNum; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := ServerListRequest{}
			req.WithRetry(retry, 0)
			_, _ = req.Make(context.Background(), &cli)
		}()
	}
	wg.Wait()
	assert.Equal(t, retry*grNum, httpCli.RequestCount)
}

func TestServerListPaginator_NextPage(t *testing.T) {
	httpCli := mocks.RequestDebugClient{}
	cli := clo.ApiClient{
		HttpClient: &httpCli,
		Options: map[string]interface{}{
			"auth_key": "secret",
			"base_url": "https://clo.ru",
		},
	}
	mocks.BodyStringFunc = func() (string, int) {
		return "1", http.StatusOK
	}
	var cases = []struct {
		ShouldFail bool
		Name       string
		Expected   string
		PGOptions  PaginatorOptions
	}{
		{
			Name: "Success",
			PGOptions: PaginatorOptions{
				Limit:  3,
				Offset: 3,
			},
			Expected: "limit=3&offset=3",
		},
		{
			Name:       "WrongLimit",
			ShouldFail: true,
			PGOptions: PaginatorOptions{
				Limit:  2,
				Offset: 3,
			},
			Expected: "limit=3&offset=3",
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			req := ServerListRequest{
				ProjectID: "id",
			}
			pg, e := NewListPaginator(&cli, req, c.PGOptions)
			assert.Nil(t, e)
			_, e = pg.NextPage(context.Background())
			if !c.ShouldFail {
				assert.Equal(t, c.Expected, httpCli.URL.RawQuery)
			} else {
				assert.NotEqual(t, c.Expected, httpCli.URL.RawQuery)
			}
		})
	}
}

func TestServerListPaginator_lastPage(t *testing.T) {
	httpCli := mocks.RequestDebugClient{}
	cli := clo.ApiClient{
		HttpClient: &httpCli,
		Options: map[string]interface{}{
			"auth_key": "secret",
			"base_url": "https://clo.ru",
		},
	}
	req := ServerListRequest{
		ProjectID: "id",
	}
	mocks.BodyStringFunc = func() (string, int) {
		return `{"count": 2, "results": [{"id": "first_item_id", "ptr": "host.com", "attached_to_server":{"id":"server_id"}},{"id": "second_item_id", "ptr": "host.com", "attached_to_server":{"id":"server_id"}}]}`,
			http.StatusOK
	}
	pg, e := NewListPaginator(&cli, req, PaginatorOptions{
		Limit:  3,
		Offset: 3,
	})
	assert.Nil(t, e)
	assert.Equal(t, false, pg.lastPage)

	_, e = pg.NextPage(context.Background())
	assert.Nil(t, e)
	assert.Equal(t, true, pg.lastPage)

	_, e = pg.NextPage(context.Background())
	assert.Equal(t, "no more pages", e.Error())
}

func TestServerList_Filtering(t *testing.T) {
	httpCli := mocks.RequestDebugClient{}
	cli := clo.ApiClient{
		HttpClient: &httpCli,
		Options: map[string]interface{}{
			"auth_key": "secret",
			"base_url": "https://clo.ru",
		},
	}
	mocks.BodyStringFunc = func() (string, int) {
		return "1", http.StatusOK
	}
	var cases = []struct {
		ShouldFail   bool
		Name         string
		OrderFields  []string
		FilterFields []FilteringField
		RawExpected  map[string][]string
	}{
		{
			Name: "Success",
			FilterFields: []FilteringField{
				{
					FieldName: "field_gt",
					Condition: "gt",
					Value:     "3",
				},
				{
					FieldName: "field_in",
					Condition: "in",
					Value:     "2,3,4",
				},
				{
					FieldName: "field_range",
					Condition: "range",
					Value:     "2:3",
				},
			},
			OrderFields: []string{
				"field3", "-field4",
			},
			RawExpected: map[string][]string{
				"field_gt__gt":       {"3"},
				"field_in__in":       {"2,3,4"},
				"field_range__range": {"2:3"},
				"order":              {"field3", "-field4"},
			},
		},
		{
			Name:       "WrongCondition",
			ShouldFail: true,
			FilterFields: []FilteringField{
				{
					FieldName: "field_gt",
					Condition: "gt",
					Value:     "3",
				},
				{
					FieldName: "field_in",
					Condition: "in",
					Value:     "2,3,4",
				},
				{
					FieldName: "field_range",
					Condition: "range",
					Value:     "2:3",
				},
			},
			OrderFields: []string{
				"field3", "-field4",
			},
			RawExpected: map[string][]string{
				"field_gt__gt":       {"2"},
				"field_in__in":       {"2,3,4"},
				"field_range__range": {"2:3"},
				"order":              {"field3", "-field4"},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			var params url.Values
			params = c.RawExpected
			expected := params.Encode()
			req := ServerListRequest{
				ProjectID: "id",
			}
			for _, ff := range c.FilterFields {
				req.FilterBy(ff)
			}
			for _, of := range c.OrderFields {
				req.OrderBy(of)
			}
			_, _ = req.Make(context.Background(), &cli)
			if !c.ShouldFail {
				assert.Equal(t, expected, httpCli.URL.RawQuery)
			} else {
				assert.NotEqual(t, expected, httpCli.URL.RawQuery)
			}
		})
	}
}
