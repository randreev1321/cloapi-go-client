package snapshots

import (
	"github.com/clo-ru/cloapi-go-client/clo"
	"github.com/clo-ru/cloapi-go-client/clo/mocks"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"sync"
	"testing"
)

func TestSnapshotListRequest_BuildRequest(t *testing.T) {
	ID := "id"
	req := SnapshotListRequest{
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
		context.Background(), http.MethodGet, mocks.MockUrl+fmt.Sprintf(snapListEndpoint, ID), nil,
	)
	expReq.Header = h
	assert.Equal(t, expReq, rawReq)
}

func TestSnapshotListRequest_Make(t *testing.T) {
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
		Req            SnapshotListRequest
		BodyStringFunc func() (string, int)
		Expected       SnapshotListResponse
	}{
		{
			Name: "Success",
			BodyStringFunc: func() (string, int) {
				return `{"count": 2,"results": [{"name":"first_item_name","parent_server":{"id":"server1_id"}},{"name":"second_item_name","child_servers":[{"id":"server2_id"}]}]}`,
					http.StatusOK
			},
			Req: SnapshotListRequest{},
			Expected: SnapshotListResponse{
				Count: 2,
				Results: []SnapshotDetailItem{
					{
						Name:         "first_item_name",
						ParentServer: ResponseItem{ID: "server1_id"},
					},
					{
						Name: "second_item_name",
						ChildServers: []ResponseItem{
							{ID: "server2_id"},
						},
					},
				},
			},
		},
		{
			Name:       "WrongCount",
			ShouldFail: true,
			BodyStringFunc: func() (string, int) {
				return `{"count": 2,"results": [{"name":"first_item_name","parent_server":{"id":"server1_id"}},{"name":"second_item_name","child_servers":[{"id":"server2_id"}]}]}`,
					http.StatusOK
			},
			Req: SnapshotListRequest{},
			Expected: SnapshotListResponse{
				Count: 1,
				Results: []SnapshotDetailItem{
					{
						Name:         "first_item_name",
						ParentServer: ResponseItem{ID: "server1_id"},
					},
					{
						Name: "second_item_name",
						ChildServers: []ResponseItem{
							{ID: "server2_id"},
						},
					},
				},
			},
		},
		{
			Name:       "WrongParentServerIdReturned",
			ShouldFail: true,
			BodyStringFunc: func() (string, int) {
				return `{"count": 2,"results": [{"name":"first_item_name","parent_server":{"id":"server_id"}},{"name":"second_item_name","child_servers":[{"id":"server2_id"}]}]}`,
					http.StatusOK
			},
			Req: SnapshotListRequest{},
			Expected: SnapshotListResponse{
				Count: 2,
				Results: []SnapshotDetailItem{
					{
						Name:         "first_item_name",
						ParentServer: ResponseItem{ID: "server1_id"},
					},
					{
						Name: "second_item_name",
						ChildServers: []ResponseItem{
							{ID: "server2_id"},
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
			Req:      SnapshotListRequest{},
			Expected: SnapshotListResponse{Results: []SnapshotDetailItem{}},
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

func TestSnapshotListRequest_MakeRetry(t *testing.T) {
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
			req := SnapshotListRequest{}
			req.WithRetry(retry, 0)
			_, _ = req.Make(context.Background(), &cli)
		}()
	}
	wg.Wait()
	assert.Equal(t, retry*grNum, httpCli.RequestCount)
}

func TestSnapshotListPaginator_NextPage(t *testing.T) {
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
		PGOptions  SnapshotPaginatorOptions
	}{
		{
			Name: "Success",
			PGOptions: SnapshotPaginatorOptions{
				Limit:  3,
				Offset: 3,
			},
			Expected: "limit=3&offset=3",
		},
		{
			Name:       "WrongLimit",
			ShouldFail: true,
			PGOptions: SnapshotPaginatorOptions{
				Limit:  2,
				Offset: 3,
			},
			Expected: "limit=3&offset=3",
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			req := SnapshotListRequest{}
			pg, e := NewSnapshotListPaginator(&cli, req, c.PGOptions)
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

func TestSnapshotListPaginator_lastPage(t *testing.T) {
	httpCli := mocks.RequestDebugClient{}
	cli := clo.ApiClient{
		HttpClient: &httpCli,
		Options: map[string]interface{}{
			"auth_key": "secret",
			"base_url": "https://clo.ru",
		},
	}
	req := SnapshotListRequest{}
	mocks.BodyStringFunc = func() (string, int) {
		return `{"count": 2, "results": [{"id": "first_item_id", "ptr": "host.com", "attached_to_server":{"id":"server_id"}},{"id": "second_item_id", "ptr": "host.com", "attached_to_server":{"id":"server_id"}}]}`,
			http.StatusOK
	}
	pg, e := NewSnapshotListPaginator(&cli, req, SnapshotPaginatorOptions{
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

func TestSnapshotList_Filtering(t *testing.T) {
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
			req := SnapshotListRequest{}
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
