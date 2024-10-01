package load_balancer

import (
	"fmt"
	intTesting "github.com/clo-ru/cloapi-go-client/v2/internal/testing"
	"github.com/clo-ru/cloapi-go-client/v2/internal/testing/mocks"
	"net/http"
	"testing"
)

func TestBalancerRulesListRequest_Build(t *testing.T) {
	ID := "id"
	tests := []struct {
		req *BalancerRulesListRequest
		url string
	}{
		{
			req: &BalancerRulesListRequest{ProjectID: ID},
			url: fmt.Sprintf(projectRulesListEndpoint, mocks.MockUrl, ID),
		},
		{
			req: &BalancerRulesListRequest{BalancerID: ID},
			url: fmt.Sprintf(balancerRulesListEndpoint, mocks.MockUrl, ID),
		},
	}
	for _, tt := range tests {
		intTesting.BuildTest(tt.req, http.MethodGet, tt.url, nil, t)
	}
}

func TestBalancerRulesListRequest_Make(t *testing.T) {
	intTesting.TestDoRequestCases(
		t,
		[]intTesting.DoTestCase{
			{
				Name: "Success",
				BodyStringFunc: func() (string, int) {
					return `{"count":2,"result":[{"id":"rule1","external_protocol_port":2,"internal_protocol_port":3, "server":"serv1"},{"id":"rule2","external_protocol_port":3,"internal_protocol_port":43,"loadbalancer":"lb1"}]}`, http.StatusOK
				},
				Req: &BalancerRulesListRequest{ProjectID: "project_id"},
				Expected: &BalancerRuleListResponse{
					Count: 2,
					Result: []BalancerRule{
						{
							ID:                   "rule1",
							ExternalProtocolPort: 2,
							InternalProtocolPort: 3,
							Server:               "serv1",
						},
						{
							ID:                   "rule2",
							ExternalProtocolPort: 3,
							InternalProtocolPort: 43,
							Loadbalancer:         "lb1",
						},
					},
				},
				Actual: &BalancerRuleListResponse{},
			},
			{
				Name:           "Error",
				ShouldFail:     true,
				BodyStringFunc: func() (string, int) { return "", http.StatusInternalServerError },
				Req:            &BalancerRulesListRequest{ProjectID: "project_id"},
			},
		},
	)
}

func TestBalancerRulesListRequest_MakeRetry(t *testing.T) {
	intTesting.ConcurrentRetryTest(&BalancerRulesListRequest{}, t)
}

func TestBalancerRulesListRequest_Filtering(t *testing.T) {
	intTesting.FilterTest(&BalancerRulesListRequest{}, t)
}
