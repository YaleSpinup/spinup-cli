package spinup

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestURIVars(t *testing.T) {
	if ContainerURI != "/api/v3/containers" {
		t.Errorf("unexpected ContainerURI %s", ContainerURI)
	}

	if ResourceURI != "/api/v3/resources" {
		t.Errorf("unexpected ResourceURI %s", ResourceURI)
	}

	if SecretsURI != "/api/v3/spaces" {
		t.Errorf("unexpected SecretsURI %s", SecretsURI)
	}

	if SizeURI != "/api/v3/sizes" {
		t.Errorf("unexpected SizeURI %s", SizeURI)
	}

	if SpaceURI != "/api/v3/spaces" {
		t.Errorf("unexpected SpaceURI %s", SpaceURI)
	}

	if StorageURI != "/api/v3/storage" {
		t.Errorf("unexpected StorageURI %s", StorageURI)
	}
}

type MockResourceInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var (
	MockInfoURI   = "/api/v3/mock"
	testMockInfos = map[string]MockResourceInfo{
		"0": {
			ID:   "0",
			Name: "copy",
		},
		"1": {
			ID:   "1",
			Name: "forgery",
		},
		"2": {
			ID:   "2",
			Name: "sham",
		},
		"3": {
			ID:   "3",
			Name: "fraud",
		},
		"4": {
			ID:   "4",
			Name: "hoax",
		},
		"5": {
			ID:   "5",
			Name: "dummy",
		},
		"6": {
			ID:   "6",
			Name: "lookalike",
		},
		"7": {
			ID:   "7",
			Name: "conterfeit",
		},
		"8": {
			ID:   "8",
			Name: "phoney",
		},
		"9": {
			ID:   "9",
			Name: "reproduction",
		},
	}
)

// GetEndpoint returns the url for a mock resource
func (m *MockResourceInfo) GetEndpoint(params map[string]string) string {
	return BaseURL + MockInfoURI + "/" + params["id"]
}
func MockResourceGetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	id := strings.TrimPrefix(r.URL.String(), MockInfoURI+"/")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	mock, ok := testMockInfos[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}

	out, err := json.Marshal(mock)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to marshall json: " + err.Error()))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "XSRF-TOKEN",
		Value: "foobar",
	})

	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func TestGetResource(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(MockResourceGetHandler))
	defer ts.Close()

	t.Logf("created server listening on %s", ts.URL)

	client, err := New(ts.URL, http.DefaultClient, "token")
	if err != nil {
		t.Errorf("expected nil error, got %s", err)
	}

	for id, expected := range testMockInfos {
		output := MockResourceInfo{}
		if err := client.GetResource(map[string]string{"id": id}, &output); err != nil {
			t.Errorf("expected nil error, got %s", err)
		}

		if !reflect.DeepEqual(expected, output) {
			t.Errorf("expected '%+v', got '%+v'", expected, output)
		}

		if client.CSRFToken != "foobar" {
			t.Errorf("expected CSRF token to be set to 'foobar', got %s", client.CSRFToken)
		}
	}

	if err := client.GetResource(map[string]string{"id": "missing"}, &MockResourceInfo{}); err == nil {
		t.Error("expected error, got nil")
	}
}

func MockResourcePutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	id := strings.TrimPrefix(r.URL.String(), MockInfoURI+"/")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	mock, ok := testMockInfos[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
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

func TestPutResource(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(MockResourcePutHandler))
	defer ts.Close()

	t.Logf("created server listening on %s", ts.URL)

	client, err := New(ts.URL, http.DefaultClient, "token")
	if err != nil {
		t.Errorf("expected nil error, got %s", err)
	}
	client.CSRFToken = "foobar"

	for id, expected := range testMockInfos {
		input, _ := json.Marshal(expected)
		output := MockResourceInfo{}
		if err := client.PutResource(map[string]string{"id": id}, input, &output); err != nil {
			t.Errorf("expected nil error, got %s", err)
		}
	}

	if err := client.GetResource(map[string]string{"id": "missing"}, &MockResourceInfo{}); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestNew(t *testing.T) {
	expected := &Client{
		HTTPClient: http.DefaultClient,
		AuthToken:  "token",
	}
	spinupUrl := "https://spinup.example.com"

	output, err := New(spinupUrl, expected.HTTPClient, "token")
	if err != nil {
		t.Errorf("expected nil error, got %s", err)
	}

	if cType := reflect.TypeOf(output).String(); cType != "*spinup.Client" {
		t.Errorf("expected a new '*spinup.Client', got '%s'", cType)
	}

	if !reflect.DeepEqual(expected, output) {
		t.Errorf("expected '%+v', got '%+v'", expected, output)
	}

	if BaseURL != spinupUrl {
		t.Errorf("expected BaseURL to be set to '%s', got '%s'", spinupUrl, BaseURL)
	}

	// TODO find a URL that throws an error
	// if _, err := New("⌘", expected.HTTPClient); err == nil {
	// 	t.Error("expected error, got nil")
	// }
}

func TestFlexIntUnmarshallJSON(t *testing.T) {
	expectedInts := map[int]FlexInt{}
	for i := 0; i <= 100; i += 1 {
		expectedInts[i] = FlexInt(i)
	}

	for i, fi := range expectedInts {
		out := FlexInt(0)
		if err := out.UnmarshalJSON([]byte(strconv.Itoa(i))); err != nil {
			t.Errorf("expected nil error, got %s", err)
		}

		if out != fi {
			t.Errorf("expected %+v, got %+v", fi, out)
		}
	}

	expectedStrings := map[string]FlexInt{}
	for i, fi := range expectedInts {
		s := strconv.Itoa(i)
		expectedStrings[s] = fi
	}

	for s, fi := range expectedStrings {
		out := FlexInt(0)
		if err := out.UnmarshalJSON([]byte(s)); err != nil {
			t.Errorf("expected nil error, got %s", err)
		}

		if out != fi {
			t.Errorf("expected %+v, got %+v", fi, out)
		}
	}

	expectedErrs := []string{"1foo", "2foo", "3foo", "4foo", "5foo"}
	for _, s := range expectedErrs {
		out := FlexInt(0)
		if err := out.UnmarshalJSON([]byte(s)); err == nil {
			t.Errorf("expected error for input %s, got nil", s)
		}
	}
}

func TestFlexIntString(t *testing.T) {
	expectedStrings := map[string]FlexInt{}
	for i := 0; i <= 100; i += 1 {
		s := strconv.Itoa(i)
		expectedStrings[s] = FlexInt(i)
	}

	for s, fi := range expectedStrings {
		if fis := fi.String(); s != fis {
			t.Errorf("expected %s, got %s", s, fis)
		}
	}
}

func TestFlexBoolUnmarshall(t *testing.T) {
	ft := FlexBool(true)
	ff := FlexBool(false)

	var out FlexBool
	if err := out.UnmarshalJSON([]byte("false")); err != nil {
		t.Errorf("got unexpected error %s", err)
	}

	if out != ff {
		t.Error("expected false")
	}

	if err := out.UnmarshalJSON([]byte("true")); err != nil {
		t.Errorf("got unexpected error %s", err)
	}

	if out != ft {
		t.Error("expected true")
	}
}

func TestFlexBoolBool(t *testing.T) {
	ft := FlexBool(true)
	if ft.Bool() != true {
		t.Error("expected true")
	}

	ff := FlexBool(false)
	if ff.Bool() != false {
		t.Error("expected false")
	}
}

func TestFlexBoolString(t *testing.T) {
	ft := FlexBool(true)
	if ft.String() != "true" {
		t.Error("expected true")
	}

	ff := FlexBool(false)
	if ff.String() != "false" {
		t.Error("expected false")
	}
}

func TestNameValueString(t *testing.T) {
	tests := []NameValue{}
	for i := 0; i < 100; i++ {
		u, _ := uuid.NewRandom()
		name := u.String()

		s := sha256.Sum256([]byte(name))
		value := hex.EncodeToString(s[:])

		tests = append(tests, NameValue{Name: name, Value: value})
	}

	for _, test := range tests {
		expected := fmt.Sprintf("%s:%s", test.Name, test.Value)
		if out := test.String(); out != expected {
			t.Errorf("expected %s, got %s", expected, out)
		}
	}
}

func TestNameValueFromString(t *testing.T) {
	tests := []NameValueFrom{}
	for i := 0; i < 100; i++ {
		u, _ := uuid.NewRandom()
		name := u.String()

		s := sha256.Sum256([]byte(name))
		valueFrom := hex.EncodeToString(s[:])

		tests = append(tests, NameValueFrom{Name: name, ValueFrom: valueFrom})
	}

	for _, test := range tests {
		expected := fmt.Sprintf("%s:%s", test.Name, test.ValueFrom)
		if out := test.String(); out != expected {
			t.Errorf("expected %s, got %s", expected, out)
		}
	}
}
