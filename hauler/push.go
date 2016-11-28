package hauler

import (
	"fmt"
	"github.com/docker/distribution/manifest/schema1"
	"github.com/docker/libtrust"
	"github.com/heroku/docker-registry-client/registry"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Push(target string, sourceTarget, destinationTarget RegistryTarget) error {

	imageJSON, err := ioutil.ReadFile(filepath.Join(target, "manifest.json"))

	if err != nil {
		fmt.Println("Error Reading Manifest File: ", err)
	}

	signedManifest := schema1.SignedManifest{}
	err = signedManifest.UnmarshalJSON(imageJSON)

	if err != nil {
		fmt.Println("Error Unmarshalling Manifest JSON: ", err)
	}

	// We only need to sign the manifest if the image name/tag changes
	var newSignedManifest *schema1.SignedManifest

	if sourceTarget.Image == destinationTarget.Image && sourceTarget.Tag == destinationTarget.Tag {
		newSignedManifest = &signedManifest
	} else {
		key, err := libtrust.GenerateECP256PrivateKey()
		if err != nil {
			return err
		}

		// Update Manifest With Image Name and Tag
		signedManifest.Manifest.Name = destinationTarget.Image
		signedManifest.Manifest.Tag = destinationTarget.Tag

		newSignedManifest, err = schema1.Sign(&signedManifest.Manifest, key)
		if err != nil {
			return fmt.Errorf("Error Signing Manifest: %v", err)
		}

	}

	hub, err := registry.New(destinationTarget.Url, destinationTarget.Username, destinationTarget.Password)

	if err != nil {
		return fmt.Errorf("Error Connecting to Registry to Push Image %v", err)
	}

	fmt.Println(fmt.Sprintf("\nPushing Image %s:%s to %s...\n", destinationTarget.Image, destinationTarget.Tag, hub.URL))

	// Push Layers to Registry
	for _, layer := range newSignedManifest.FSLayers {

		tarFileName := layer.BlobSum.String() + ".tar"
		sourceLayerPath := filepath.Join(target, tarFileName)

		_, err := os.Stat(sourceLayerPath)

		if err != nil {
			return fmt.Errorf("Error Locating Path", err)
		}

		uploadError := uploadLayer(hub, sourceLayerPath, layer, newSignedManifest)

		if uploadError != nil {
			return fmt.Errorf("Error Uploading Layer", uploadError)
		}

	}

	err = hub.PutManifest(destinationTarget.Image, destinationTarget.Tag, newSignedManifest)
	if err != nil {
		fmt.Println(err)
	}

	return nil

}

func uploadLayer(hub *registry.Registry, sourceLayerPath string, layer schema1.FSLayer, signedManifest *schema1.SignedManifest) error {
	reader, err := os.Open(sourceLayerPath)
	if err != nil {
		return fmt.Errorf("Failed to Open the Reader ", err)
	}

	defer reader.Close()

	exists, err := hub.HasLayer(signedManifest.Manifest.Name, layer.BlobSum)

	if err != nil {
		// Failed to check if layer exists in dest, fail!
		return fmt.Errorf("Unable to check if layer exists in destination registry. Error: ", err)
	}

	if !exists {
		fmt.Println("Pushing ", layer.BlobSum)

		uploadError := hub.UploadLayer(signedManifest.Manifest.Name, layer.BlobSum, reader)
		if uploadError != nil {
			return fmt.Errorf("Unable to upload layer, err: %s", err)
		}
	} else {
		fmt.Println("Skipping Layer ", layer.BlobSum)
	}

	return nil

}
