/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package image

import (
	"context"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/manifest"
	"strings"
)

//getLatestImageDigest - get digest of the latest version of a image tag
func (iu *ImageUpdater) getLatestImageDigest(ctx context.Context, imageName, imageId string) (string, error) {
	latestDigest, err := iu.getDockerDigest(ctx, imageName)
	if err != nil {
		iu.Log.Error(err, "Failed to get latest digest for image", "image", imageName)
		return "", err
	}
	latestImage := strings.Split(imageId, "@")[0] + "@" + latestDigest
	return latestImage, nil
}

//getDockerDigest - get digest of a given image
func (iu *ImageUpdater) getDockerDigest(ctx context.Context, imageName string) (string, error) {
	ref, err := docker.ParseReference("//" + imageName)
	if err != nil {
		return "", err
	}
	img, err := ref.NewImage(ctx, nil)
	if err != nil {
		return "", err
	}
	defer img.Close()
	b, _, err := img.Manifest(ctx)
	if err != nil {
		return "", err
	}
	digest, err := manifest.Digest(b)
	if err != nil {
		return "", err
	}
	digestStr := string(digest)
	return digestStr, nil
}
