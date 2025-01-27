package codegenerator_test

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/grpc-gateway/v2/internal/codegenerator"
	"github.com/gogo/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/pluginpb"
)

var parseReqTests = []struct {
	name string
	in   io.Reader
	out  *pluginpb.CodeGeneratorRequest
	err  error
}{
	{
		"Empty input should produce empty output",
		mustGetReader(&pluginpb.CodeGeneratorRequest{}),
		&pluginpb.CodeGeneratorRequest{},
		nil,
	},
	{
		"Invalid reader should produce error",
		&invalidReader{},
		nil,
		fmt.Errorf("failed to read code generator request: invalid reader"),
	},
	{
		"Invalid proto message should produce error",
		strings.NewReader("{}"),
		nil,
		fmt.Errorf("failed to unmarshal code generator request: unexpected EOF"),
	},
}

func TestParseRequest(t *testing.T) {
	for _, tt := range parseReqTests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := codegenerator.ParseRequest(tt.in)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("got %v, want %v", err, tt.err)
			}
			if diff := cmp.Diff(out, tt.out, protocmp.Transform()); diff != "" {
				t.Errorf(diff)
			}
		})
	}
}

func mustGetReader(pb proto.Message) io.Reader {
	b, err := proto.Marshal(pb)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(b)
}

type invalidReader struct {
}

func (*invalidReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("invalid reader")
}
