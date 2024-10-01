package load_balancer

import (
	"fmt"
	intTesting "github.com/clo-ru/cloapi-go-client/v2/internal/testing"
	"github.com/clo-ru/cloapi-go-client/v2/internal/testing/mocks"
	"net/http"
	"testing"
)

func TestBalancerRulesDeleteRequest_BuildRequest(t *testing.T) {
	dID := "object_id"
	req := &BalancerRulesDeleteRequest{ObjectId: dID}
	intTesting.BuildTest(req, http.MethodDelete, fmt.Sprintf(balancerRulesDeleteEndpoint, mocks.MockUrl, dID), nil, t)
}

func TestBalancerRulesDeleteRequest_Make(t *testing.T) {
	intTesting.TestDoRequestCases(
		t,
		[]intTesting.DoTestCase{
			{
				Name: "Success",
				BodyStringFunc: func() (string, int) {
					return "", http.StatusOK
				},
				Req: &BalancerRulesDeleteRequest{ObjectId: "object_id"},
			},
			{
				Name:           "Error",
				ShouldFail:     true,
				CheckError:     true,
				BodyStringFunc: func() (string, int) { return "", http.StatusInternalServerError },
				Req:            &BalancerRulesDeleteRequest{ObjectId: "object_id"},
			},
		},
	)
}

func TestBalancerRulesDeleteRequest_MakeRetry(t *testing.T) {
	intTesting.ConcurrentRetryTest(&BalancerRulesDeleteRequest{}, t)
}
