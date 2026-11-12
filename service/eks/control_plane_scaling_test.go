package eks_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
)

func TestControlPlaneScalingConfigWireContract(t *testing.T) {
	var requests []map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decode request: %v", err)
				return
			}
			requests = append(requests, body)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/clusters":
			_, _ = w.Write([]byte(`{"cluster":{"name":"test","controlPlaneScalingConfig":{"tier":"tier-large"}}}`))
		case "/clusters/test/update-config":
			_, _ = w.Write([]byte(`{"update":{"id":"update-1","status":"InProgress"}}`))
		case "/clusters/test":
			_, _ = w.Write([]byte(`{"cluster":{"name":"test","controlPlaneScalingConfig":{"tier":"tier-2large"}}}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := eks.New(session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials("test", "test", ""),
		Endpoint:    aws.String(server.URL),
		Region:      aws.String("ru-msk"),
	})))
	created, err := client.CreateCluster(&eks.CreateClusterInput{
		Name:                      aws.String("test"),
		Version:                   aws.String("1.30"),
		PlatformVersion:           aws.String("eks.3"),
		ResourcesVpcConfig:        &eks.VpcConfigRequest{SubnetIds: aws.StringSlice([]string{"subnet-1"})},
		ControlPlaneScalingConfig: &eks.ControlPlaneScalingConfig{Tier: aws.String("tier-large")},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got := aws.StringValue(created.Cluster.ControlPlaneScalingConfig.Tier); got != "tier-large" {
		t.Fatalf("create response tier: %q", got)
	}
	if got := requests[0]["platformVersion"]; got != "eks.3" {
		t.Fatalf("create platform version: %v", got)
	}

	_, err = client.UpdateClusterConfig(&eks.UpdateClusterConfigInput{
		Name:                      aws.String("test"),
		ControlPlaneScalingConfig: &eks.ControlPlaneScalingConfig{Tier: aws.String("tier-2large")},
	})
	if err != nil {
		t.Fatal(err)
	}
	for i, want := range []string{"tier-large", "tier-2large"} {
		got := requests[i]["controlPlaneScalingConfig"].(map[string]interface{})["tier"]
		if got != want {
			t.Fatalf("request %d tier: got %v, want %s", i, got, want)
		}
	}
	described, err := client.DescribeCluster(&eks.DescribeClusterInput{Name: aws.String("test")})
	if err != nil {
		t.Fatal(err)
	}
	if got := aws.StringValue(described.Cluster.ControlPlaneScalingConfig.Tier); got != "tier-2large" {
		t.Fatalf("describe response tier: %q", got)
	}
}
