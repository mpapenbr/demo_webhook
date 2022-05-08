package ghwebhook

import (
	"testing"
)

// Note: These tests are no tests ;) More like do-i-use-the-API-correct methods

func TestQueryFile(t *testing.T) {
	type args struct {
		repo    string
		fileRef string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "dummy", args: args{repo: "demo_app1", fileRef: "HEAD:version.txt"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			QueryFile(tt.args.repo, tt.args.fileRef)
		})
	}
}

func TestUpdateFileGraphQL(t *testing.T) {
	type args struct {
		repo       string
		fileRef    string
		newContent string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "dummy", args: args{repo: "demo_app1", fileRef: "HEAD:version.txt"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := UpdateFileGraphQL(tt.args.repo, tt.args.fileRef, tt.args.newContent); got != tt.want {
				t.Errorf("UpdateFileGraphQL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateFileRest(t *testing.T) {
	type args struct {
		repoOwner  string
		repo       string
		fileRef    string
		newContent string
		oid        string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "dummyREST", args: args{
			repoOwner: "mpapenbr", repo: "demo_app1", fileRef: "version.txt", newContent: "version: v1.0", oid: "3bc4fa6046da515d14fdcedae44f43e55d6c0f97",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			UpdateFileRest(tt.args.repoOwner, tt.args.repo, tt.args.fileRef, tt.args.newContent, tt.args.oid)
		})
	}
}
