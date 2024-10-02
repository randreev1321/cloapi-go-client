package load_balancer

import (
	"fmt"
	intTesting "github.com/clo-ru/cloapi-go-client/v2/internal/testing"
	"github.com/clo-ru/cloapi-go-client/v2/internal/testing/mocks"
	"net/http"
	"testing"
)

func TestBalancerRulesDetailRequest_Build(t *testing.T) {
	dID := "id"
	req := &BalancerRulesDetailRequest{ObjectId: dID}
	intTesting.BuildTest(req, http.MethodGet, fmt.Sprintf(balancerRulesDetailEndpoint, mocks.MockUrl, dID), nil, t)
}

func TestBalancerRulesDetailRequest_Make(t *testing.T) {
	intTesting.TestDoRequestCases(
		t,
		[]intTesting.DoTestCase{
			{
				Name: "Success",
				BodyStringFunc: func() (string, int) {
					return `{"result":{"id":"id1","external_protocol_port":80,"internal_protocol_port":80,"server":"id2","loadbalancer":"id3","status":"ACTIVE","address":"id4"}}`, http.StatusOK
				},
				Req: &BalancerRulesDetailRequest{ObjectId: "id"},
				Expected: &BalancerRuleDetailResponse{
					Result: BalancerRule{
						ID:                   "id1",
						ExternalProtocolPort: 80,
						InternalProtocolPort: 80,
						Server:               "id2",
						Loadbalancer:         "id3",
						Status:               "ACTIVE",
						Address:              "id4",
					},
				},
				Actual: &BalancerRuleDetailResponse{},
			},
			{
				Name:       "Error",
				ShouldFail: true,
				CheckError: true,
				BodyStringFunc: func() (string, int) {
					return "", http.StatusInternalServerError
				},
				Req: &BalancerRulesDetailRequest{ObjectId: "id"},
			},
		},
	)
}

func TestBalancerRulesDetailRequest_MakeRetry(t *testing.T) {
	intTesting.ConcurrentRetryTest(&BalancerRulesDetailRequest{}, t)
}
