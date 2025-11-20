//go:build e2e
// +build e2e

package e2e

import (
	"encoding/json"
	"strconv"
	"testing"

	e2e "test-question/internal/tests/e2esuite"

	"github.com/stretchr/testify/suite"
)

type FullFlowResponse struct {
	ID int `json:"id"`
}

type FullE2ESuite struct {
	e2e.E2ESuite
}

func TestE2ESuite(t *testing.T) {
	s := &FullE2ESuite{}
	suite.Run(t, s)
}

func (f *FullE2ESuite) Test_FullFlow() {
	var qID int

	// ==== 1. Alice creates question ====
	{
		resp := f.IAmAlice().POST("/questions", map[string]any{
			"text": "hello world",
		})
		f.Require().Equal(201, resp.StatusCode)

		var out FullFlowResponse
		json.NewDecoder(resp.Body).Decode(&out)
		qID = out.ID
		f.True(qID > 0)
	}

	// ==== 2. Bob creates answer #1 ====
	var a1 int
	{
		resp := f.IAmBob().POST("/questions/"+strconv.Itoa(qID)+"/answers", map[string]any{
			"text": "first answer",
		})
		f.Require().Equal(201, resp.StatusCode)

		var out FullFlowResponse
		json.NewDecoder(resp.Body).Decode(&out)
		a1 = out.ID
		f.True(a1 > 0)
	}

	// ==== 3. Bob creates answer #2 ====
	var a2 int
	{
		resp := f.IAmBob().POST("/questions/"+strconv.Itoa(qID)+"/answers", map[string]any{
			"text": "second answer",
		})
		f.Require().Equal(201, resp.StatusCode)

		var out FullFlowResponse
		json.NewDecoder(resp.Body).Decode(&out)
		a2 = out.ID
		f.True(a2 > 0)
	}

	// ==== 4. GET question, must contain 2 answers ====
	{
		resp := f.IAmAlice().GET("/questions/" + strconv.Itoa(qID))
		f.Require().Equal(200, resp.StatusCode)

		var out map[string]any
		json.NewDecoder(resp.Body).Decode(&out)

		answers := out["answers"].([]any)
		f.Len(answers, 2)
	}

	// ==== 5. Alice tries delete answer #1 (must be 403) ====
	{
		resp := f.IAmAlice().DELETE("/answers/" + strconv.Itoa(a1))
		f.Require().Equal(403, resp.StatusCode)
	}

	// ==== 6. Bob deletes answer #1 ====
	{
		resp := f.IAmBob().DELETE("/answers/" + strconv.Itoa(a1))
		f.Require().Equal(204, resp.StatusCode)
	}

	// ==== 7. Bob tries delete answer #1 again (404) ====
	{
		resp := f.IAmBob().DELETE("/answers/" + strconv.Itoa(a1))
		f.Require().Equal(404, resp.StatusCode)

		resp = f.IAmBob().DELETE("/questions/" + strconv.Itoa(qID))
		f.Require().Equal(403, resp.StatusCode)
	}

	// ==== 8. Alice deletes question (answer #2 must be deleted CASCADE) ====
	{
		resp := f.IAmAlice().DELETE("/questions/" + strconv.Itoa(qID))
		f.Require().Equal(204, resp.StatusCode)
	}

	// ==== 9. GET question must be 404 ====
	{
		resp := f.IAmAlice().GET("/questions/" + strconv.Itoa(qID))
		f.Require().Equal(404, resp.StatusCode)
	}

	// ==== 10. GET answer #2 must be 404 (deleted by cascade) ====
	{
		resp := f.IAmBob().GET("/answers/" + strconv.Itoa(a2))
		f.Require().Equal(404, resp.StatusCode)
	}

	// ==== 11. DELETE answer #2 must be 404 ====
	{
		resp := f.IAmBob().DELETE("/answers/" + strconv.Itoa(a2))
		f.Require().Equal(404, resp.StatusCode)
	}

	// ==== 12. DELETE question again must be 404 ====
	{
		resp := f.IAmAlice().DELETE("/questions/" + strconv.Itoa(qID))
		f.Require().Equal(404, resp.StatusCode)
	}
}
