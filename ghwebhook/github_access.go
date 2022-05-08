package ghwebhook

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/go-github/v44/github"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var (
	v3Client *github.Client   = nil
	v4Client *githubv4.Client = nil
)

func QueryFile(repo string, fileRef string) string {
	client := getV4Client()
	var q struct {
		Repository struct {
			Object struct {
				Id   githubv4.ID
				Oid  githubv4.GitObjectID
				Blob struct {
					ByteSize githubv4.Int
					Text     githubv4.String
				} `graphql:"... on Blob"`
			} `graphql:"object(expression:$fileRef) "`
		} `graphql:"repository(owner:$owner, name:$repoName)"`
	}
	variables := map[string]interface{}{
		"owner":    githubv4.String("mpapenbr"),
		"repoName": githubv4.String(repo),
		"fileRef":  githubv4.String(fileRef),
	}
	err := client.Query(context.Background(), &q, variables)
	if err != nil {
		log.Fatalf("QueryError %v", err)
	}
	return string(q.Repository.Object.Blob.Text)
	// fmt.Printf("%+v\n", q)
}

// Note: This does not work atm. Not here, not in insomnia, not with curl.
// all yield: Something went wrong while executing your query. Please include `0F11:33C2:7D3186:7FEA60:6277DCCB` when reporting this issue.
func UpdateFileGraphQL(repo string, fileRef string, newContent string) string {
	client := getV4Client()
	var m struct {
		CreateCommitOnBranch struct {
			Commit struct {
				Url githubv4.URI
			}
		} `graphql:"createCommitOnBranch(input: $input)"`
	}

	input := githubv4.CreateCommitOnBranchInput{
		Branch: githubv4.CommittableBranch{
			RepositoryNameWithOwner: githubv4.NewString("mpapenbr/demo_app1"),
			BranchName:              githubv4.NewString("master"),
		},
		Message: githubv4.CommitMessage{
			Headline: *githubv4.NewString("Autocommit by ghwebhook"),
		},
		ExpectedHeadOid: *githubv4.NewGitObjectID("enterSomeGitObjectIdHere"),
		FileChanges: &githubv4.FileChanges{Additions: &[]githubv4.FileAddition{
			{Path: *githubv4.NewString(githubv4.String(fileRef))},
			{Contents: *githubv4.NewBase64String(githubv4.Base64String(newContent))},
		}},
	}

	err := client.Mutate(context.Background(), &m, input, nil)
	if err != nil {
		log.Fatalf("MutationError %v", err)
	}
	return ""
	// fmt.Printf("%+v\n", q)
}

func UpdateFileRest(repoOwner string, repo string, fileRef string, newContent string, message string, oid string) {
	client := getV3Client()
	_, resp, err := client.Repositories.UpdateFile(context.Background(), repoOwner, repo, fileRef, &github.RepositoryContentFileOptions{
		Content: []byte(newContent),
		Message: github.String(message),
		SHA:     github.String(oid),
	})
	if err != nil {
		log.Fatalf("UpdateFielRest %v", err)
	}
	fmt.Printf("%v", resp)
}

func GetContentRest(repoOwner string, repo string, fileRef string) (*github.RepositoryContent, *github.Response) {
	client := getV3Client()
	fileContent, _, resp, err := client.Repositories.GetContents(context.Background(), repoOwner, repo, fileRef, &github.RepositoryContentGetOptions{})
	if err != nil {
		log.Fatalf("GetContentRest %v", err)
	}
	fmt.Printf("%v\n", resp)
	return fileContent, resp
}

func getV4Client() *githubv4.Client {
	if v4Client == nil {
		src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
		httpClient := oauth2.NewClient(context.Background(), src)
		v4Client = githubv4.NewClient(httpClient)
	}
	return v4Client
}

func getV3Client() *github.Client {
	if v3Client == nil {
		src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})
		httpClient := oauth2.NewClient(context.Background(), src)
		v3Client = github.NewClient(httpClient)
	}
	return v3Client
}
