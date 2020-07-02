package spinup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestSpacesGetEndpoint(t *testing.T) {
	resource := Spaces{}
	expected := "http://localhost:8090/api/v2/spaces"

	if out := resource.GetEndpoint(map[string]string{}); out != expected {
		t.Errorf("expected %s, got %s", expected, out)
	}
}

func TestSpaceGetEndpoint(t *testing.T) {
	resource := Space{}
	expected := "http://localhost:8090/api/v2/spaces/123"

	if out := resource.GetEndpoint(map[string]string{"id": "123"}); out != expected {
		t.Errorf("expected %s, got %s", expected, out)
	}
}

func TestGetSpaceGetEndpoint(t *testing.T) {
	resource := GetSpace{}
	expected := "http://localhost:8090/api/v2/spaces/123"

	if out := resource.GetEndpoint(map[string]string{"id": "123"}); out != expected {
		t.Errorf("expected %s, got %s", expected, out)
	}
}

func TestSpaceCostGetEndpoint(t *testing.T) {
	resource := SpaceCost{}
	expected := "http://localhost:8090/api/v2/spaces/123/cost"

	if out := resource.GetEndpoint(map[string]string{"id": "123"}); out != expected {
		t.Errorf("expected %s, got %s", expected, out)
	}
}

type mockResourceOutput struct {
	Resources []*Resource `json:"resources"`
}

var mockResourceList map[string]*mockResourceOutput

func newMockResourceOutput(num int) *mockResourceOutput {
	resources := make([]*Resource, 0, num)
	for i := 0; i < num; i++ {
		fi := FlexInt(i)

		r := &Resource{
			ID:   &fi,
			Name: fmt.Sprintf("resource-%0.3d", i),
			Type: &Offering{
				Flavor: "linux",
			},
		}
		resources = append(resources, r)
	}
	return &mockResourceOutput{
		Resources: resources,
	}
}

func MockResourcesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	id := strings.TrimPrefix(r.URL.String(), "/api/v2/spaces/")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	mock, ok := mockResourceList[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}

	if id == "brokenJSON" {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{"))
		return
	}

	if id == "400error" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	out, err := json.Marshal(mock)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to marshall json: " + err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func TestResources(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(MockResourcesHandler))
	defer ts.Close()

	t.Logf("created server listening on %s", ts.URL)

	client, err := New(ts.URL, http.DefaultClient)
	if err != nil {
		t.Errorf("expected nil error, got %s", err)
	}

	mockResourceList = make(map[string]*mockResourceOutput)
	for i := 0; i < 100; i++ {
		spaceId := strconv.Itoa(i)
		mockResourceList[spaceId] = newMockResourceOutput(100)
	}
	mockResourceList["brokenJSON"] = nil
	mockResourceList["400error"] = nil

	for i := 0; i < 10; i++ {
		spaceId := strconv.Itoa(i)

		out, err := client.Resources(spaceId)
		if err != nil {
			t.Errorf("expected nil error, got %s", err)
		}

		if !reflect.DeepEqual(mockResourceList[spaceId].Resources, out) {
			t.Error("expected:")
			for _, e := range mockResourceList[spaceId].Resources {
				t.Errorf("%+v\n", *e)
			}

			t.Error("got:")
			for _, o := range out {
				t.Errorf("%+v\n", *o)
			}
		}
	}

	if _, err := client.Resources("brokenJSON"); err == nil {
		t.Error("expected error for broken JSON, got nil")
	}

	if _, err := client.Resources("400error"); err == nil {
		t.Error("expected error for 400 error, got nil")
	}
}
