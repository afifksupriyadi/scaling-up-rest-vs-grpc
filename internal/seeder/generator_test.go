package seeder

import (
	"encoding/json"
	"fmt"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
)

func TestGenerateStudentsSample(t *testing.T) {
	resp, err := ToStudentResponse(1)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}

	b, err := protojson.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var pretty map[string]interface{}
	json.Unmarshal(b, &pretty)
	out, _ := json.MarshalIndent(pretty, "", "  ")
	fmt.Println(string(out))
}
