package ghwebhook

import (
	b64 "encoding/base64"
	"fmt"

	"github.com/go-playground/webhooks/v6/github"
)

func ProcessNewRelease(config *Config, release github.ReleasePayload) {
	fmt.Printf("Incoming release event from %s\n", release.Repository.FullName)
	for _, action := range config.Actions {
		if action.From == release.Repository.Name {
			commitComponent := release.Release.Name
			if action.Component != "" {
				commitComponent = &action.Component
			}
			for _, update := range action.Update {
				repoOwner := release.Repository.Owner.Login
				fmt.Printf("Fetching %s from %s/%s\n", update.File, repoOwner, update.Repo)
				content, resp := GetContentRest(repoOwner, update.Repo, update.File)
				fileContent, _ := b64.StdEncoding.DecodeString(*content.Content)
				fmt.Printf("Content: <%s>\n", fileContent)
				fmt.Printf("RegEx: <%s>\n", update.Regex)
				fmt.Printf("Resp: %+v\n", *resp)
				newVersion := ReplaceVersion(fileContent, update.Regex, release.Release.TagName)
				if string(newVersion) != string(fileContent) {
					fmt.Printf("Updating file %s\n", update.File)
					fmt.Printf("NewContent: <%s>\n", string(newVersion))

					message := fmt.Sprintf("pkg: Bump %s to %s", *commitComponent, release.Release.TagName)
					UpdateFileRest(repoOwner, update.Repo, update.File, string(newVersion), message, *content.SHA)
				} else {
					fmt.Println("No changes detected")
				}

			}
		}
	}
}
