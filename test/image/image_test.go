package image

import (
	"context"
	"github.com/wadeling/registry-client/image"
	"github.com/wadeling/registry-client/types"
	"testing"
	"time"
)

const (
	registryImage    = "registry:2"
	registryPort     = "5443/tcp"
	registryUsername = "wadeling"
	registryPassword = "testpassword"
	imageName 		 = "wade23/dev:nginx-1.15"
)

func TestGetRemoteImage(t *testing.T) {
	t.Log("start get remote image")
	opt := types.DockerOption{
		Timeout:  600 * time.Second,
		SkipPing: true,
	}

	ctx := context.Background()
	img, _, err := image.NewDockerImage(ctx, imageName, opt)
	if err != nil {
		t.Fatalf("get remote image err %v",err)
	}
	imageId,err := img.ID()
	if err != nil {
		t.Fatalf("get img id err %v",err)
	}
	t.Logf("get image id %s,name %s",imageId,img.Name())

}
