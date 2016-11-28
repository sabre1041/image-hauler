package main

import (
	"fmt"
	"github.com/pborman/uuid"
	"github.com/sabre1041/image-hauler/hauler"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	defaultTargetDir = "/tmp/image-hauler"
)

func run(config *hauler.Config) error {
	_, err := os.Stat(config.StorageDir)

	if err != nil {
		return fmt.Errorf("Error Verifying Storage Directory", err)
	}

	// Generate Target Directory Path
	uuid := uuid.NewRandom().String()

	target := filepath.Join(config.StorageDir, uuid)

	// Create Storage Objects
	sourceTarget := &hauler.RegistryTarget{
		Url:      config.SourceRegistry,
		Username: config.SourceUsername,
		Password: config.SourcePassword,
		Image:    config.SourceImage,
		Tag:      config.SourceTag,
	}

	destinationTarget := &hauler.RegistryTarget{
		Url:      config.DestinationRegistry,
		Username: config.DestinationUsername,
		Password: config.DestinationPassword,
		Image:    config.DestinationImage,
		Tag:      config.DestinationTag,
	}

	h := &hauler.Hauler{
		StorageDir:                target,
		SourceRegistryTarget:      *sourceTarget,
		DestinationRegistryTarget: *destinationTarget,
	}

	// Create Storage Directory
	os.MkdirAll(h.StorageDir, 0777)

	// Setup Cleanup of Storage Directory
	defer cleanupDirectory(h.StorageDir)

	// Pull the Source Image
	pull_err := hauler.Pull(h.StorageDir, h.SourceRegistryTarget)

	if pull_err != nil {
		return pull_err
	}

	// Push Image to Destination
	push_err := hauler.Push(h.StorageDir, h.SourceRegistryTarget, h.DestinationRegistryTarget)

	if push_err != nil {
		return push_err
	}

	return nil

}

func main() {

	// TODO: Remove this and use proper logging
	log.SetOutput(ioutil.Discard)

	config := &hauler.Config{}

	var rootCmd = &cobra.Command{
		Use:   "image-hauler",
		Short: "image-hauler is a Tool to Transfer Docker Images Between Registries",
		Run: func(cmd *cobra.Command, args []string) {

			if config.SourceRegistry == "" || config.SourceImage == "" || config.SourceTag == "" || config.DestinationRegistry == "" || config.DestinationImage == "" || config.DestinationTag == "" || config.StorageDir == "" {
				fmt.Println("ERROR: Source and Destination Components Must Be Provided\n")
				cmd.Help()
				os.Exit(1)
			}

			if err := run(config); err != nil {
				fmt.Println(err)
				os.Exit(-1)
			}

		},
	}

	rootCmd.Flags().StringVar(&(config.SourceRegistry), "source-registry", "", "Address of the Source Docker Registry")
	rootCmd.Flags().StringVar(&(config.SourceImage), "source-image", "", "Source Docker Image")
	rootCmd.Flags().StringVar(&(config.SourceTag), "source-tag", "latest", "Source Docker Tag")
	rootCmd.Flags().StringVar(&(config.SourceUsername), "source-username", "", "Username for the Source Registry")
	rootCmd.Flags().StringVar(&(config.SourcePassword), "source-password", "", "Password for the Source Registry")

	rootCmd.Flags().StringVar(&(config.DestinationRegistry), "destination-registry", "", "Address of the Destination Docker Registry")
	rootCmd.Flags().StringVar(&(config.DestinationImage), "destination-image", "", "Destination Docker Image")
	rootCmd.Flags().StringVar(&(config.DestinationTag), "destination-tag", "latest", "Destination Docker Tag")
	rootCmd.Flags().StringVar(&(config.SourceUsername), "destination-username", "", "Username for the Destination Registry")
	rootCmd.Flags().StringVar(&(config.DestinationPassword), "destination-password", "", "Password for the Destination Registry")

	rootCmd.Flags().StringVar(&(config.StorageDir), "storage-dir", defaultTargetDir, "Location to Store Docker Images During Execution")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}

func cleanupDirectory(d string) {
	fmt.Println("Removeing Storage directory", d)
	os.RemoveAll(d)
}
