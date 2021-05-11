package image

import (
	"context"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	multierror "github.com/hashicorp/go-multierror"
	"golang.org/x/xerrors"

	"github.com/wadeling/registry-client/types"
)

type extender interface {
	RepoTags() []string
	RepoDigests() []string
}

type Image struct {
	name   string
	client v1.Image
	extender
}

func (img Image) Name() string {
	return img.name
}

func (img Image) ID() (string, error) {
	h, err := img.client.ConfigName()
	if err != nil {
		return "", xerrors.Errorf("unable to get the image ID: %w", err)
	}
	return h.String(), nil
}

func (img Image) ConfigBlob() ([]byte, error) {
	return img.client.RawConfigFile()
}

func (img Image) LayerIDs() ([]string, error) {
	conf, err := img.client.ConfigFile()
	if err != nil {
		return nil, xerrors.Errorf("unable to get the config file: %w", err)
	}

	var layerIDs []string
	for _, d := range conf.RootFS.DiffIDs {
		layerIDs = append(layerIDs, d.String())
	}
	return layerIDs, nil
}

func (img Image) LayerByDiffID(h v1.Hash) (v1.Layer, error) {
	return img.client.LayerByDiffID(h)
}

func NewDockerImage(ctx context.Context, imageName string, option types.DockerOption) (Image, func(), error) {
	img, ext, cleanup, err := newDockerImage(ctx, imageName, option)
	if err != nil {
		return Image{}, func() {}, err
	}
	return Image{
		name:     imageName,
		client:   img,
		extender: ext,
	}, cleanup, nil
}

func newDockerImage(ctx context.Context, imageName string, option types.DockerOption) (v1.Image, extender, func(), error) {
	var errs error

	var nameOpts []name.Option
	if option.NonSSL {
		nameOpts = append(nameOpts, name.Insecure)
	}
	ref, err := name.ParseReference(imageName, nameOpts...)
	if err != nil {
		return nil, nil, func() {}, xerrors.Errorf("failed to parse the image name: %w", err)
	}


	// Try accessing Docker Registry
	img, ext, err := GetRemoteImage(ctx, ref, option)
	if err == nil {
		// Return v1.Image if the image is found in Docker Registry
		return img, ext, func() {}, nil
	}

	errs = multierror.Append(errs, err)
	return nil, nil, func() {}, errs
}


