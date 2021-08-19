package client_test

import (
	"testing"

	client "github.com/jenkins-zh/jenkins-client/pkg"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	ghClient := client.GitHubReleaseClient{}

	assert.Nil(t, ghClient.Client)
	ghClient.Init()
	assert.NotNil(t, ghClient.Client)
}

func TestGetLatestReleaseAsset(t *testing.T) {
	c, teardown := client.PrepareForGetLatestReleaseAsset() //setup()
	defer teardown()

	ghClient := client.GitHubReleaseClient{
		Client: c,
	}
	asset, err := ghClient.GetLatestReleaseAsset("o", "r")

	assert.Nil(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, "tagName", asset.TagName)
	assert.Equal(t, "body", asset.Body)
}

func TestGetLatestJCLIAsset(t *testing.T) {
	c, teardown := client.PrepareForGetLatestJCLIAsset() //setup()
	defer teardown()

	ghClient := client.GitHubReleaseClient{
		Client: c,
	}
	asset, err := ghClient.GetLatestJCLIAsset()

	assert.Nil(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, "tagName", asset.TagName)
	assert.Equal(t, "body", asset.Body)
}

func TestGetJCLIAsset(t *testing.T) {
	c, teardown := client.PrepareForGetJCLIAsset("tagName") //setup()
	defer teardown()

	ghClient := client.GitHubReleaseClient{
		Client: c,
	}
	asset, err := ghClient.GetJCLIAsset("tagName")

	assert.Nil(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, "tagName", asset.TagName)
	assert.Equal(t, "body", asset.Body)
}

func TestGetReleaseAssetByTagName(t *testing.T) {
	c, teardown := client.PrepareForGetReleaseAssetByTagName() //setup()
	defer teardown()

	ghClient := client.GitHubReleaseClient{
		Client: c,
	}
	asset, err := ghClient.GetReleaseAssetByTagName("jenkins-zh", "jenkins-cli", "tagName")

	assert.Nil(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, "tagName", asset.TagName)
	assert.Equal(t, "body", asset.Body)
}
