package hauler

import (
	"fmt"
	"github.com/docker/distribution/digest"
	"github.com/heroku/docker-registry-client/registry"
	"io"
	"os"
	"path/filepath"
)

func Pull(target string, sourceTarget RegistryTarget) error {

	hub, err := registry.New(sourceTarget.Url, sourceTarget.Username, sourceTarget.Password)

	if err != nil {
		return fmt.Errorf("Error Connecting to Registry to Pull Image", err)
	}

	fmt.Println(fmt.Sprintf("Pulling Image %s:%s from %s...\n", sourceTarget.Image, sourceTarget.Tag, hub.URL))

	manifest, err := hub.Manifest(sourceTarget.Image, sourceTarget.Tag)

	if err != nil {
		return fmt.Errorf("Could Not Locate Manifest File")
	}
	// Write Manifiest
	bManifest, _ := manifest.MarshalJSON()
	manifestFile, err := os.Create(filepath.Join(target, "manifest.json"))
	if err != nil {
		return fmt.Errorf("Unable to create Manifest file")
	}

	manifestFile.Write(bManifest)

	for _, layer := range manifest.Manifest.FSLayers {
		// or obtain the digest from an existing manifest's FSLayer list
		digest, _ := digest.ParseDigest(layer.BlobSum.String())

		fmt.Println("Pulling image layer", digest)

		reader, err := hub.DownloadLayer(sourceTarget.Image, digest)

		if reader != nil {
			defer reader.Close()
		}
		if err != nil {
			return err
		}
		tarFile, err := os.Create(filepath.Join(target, layer.BlobSum.String()+".tar"))
		if err != nil {
			return fmt.Errorf("Unable to create tar file")
		}

		_, err = io.Copy(tarFile, reader)
		if err != nil {
			return fmt.Errorf("Unable to create tar file")
		}
	}

	return nil

}
