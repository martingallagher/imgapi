package main

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/martingallagher/imgapi/service"
	"github.com/spf13/cobra"
)

var (
	config *service.Config

	cmdImgAPI = &cobra.Command{
		Use:   "imgapi",
		Short: "The ImgAPI CLI.",
	}

	cmdServer = &cobra.Command{
		Use:   "server",
		Short: "ImgAPI server",
		Run: func(_ *cobra.Command, args []string) {
			startServer(config)
		},
	}

	cmdUpload = &cobra.Command{
		Use:   "upload",
		Short: "Upload an image.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			b, err := ioutil.ReadFile(args[0])

			if err != nil {
				log.Fatal(err)
			}

			conn, client, err := getClient(config)

			if err != nil {
				log.Fatal(err)
			}

			defer conn.Close() // nolint: errcheck

			resp, err := client.Upload(context.Background(), &service.UploadRequest{Data: b})

			if err != nil {
				log.Println(err)

				return
			}

			log.Printf("Image uploaded, ID: %s", resp.Id)
		},
	}

	cmdDownload = &cobra.Command{
		Use:   "download",
		Short: "Download an existing image, by ID",
		Args:  cobra.MinimumNArgs(1),
		Run: func(_ *cobra.Command, args []string) {
			conn, client, err := getClient(config)

			if err != nil {
				log.Fatal(err)
			}

			defer conn.Close() // nolint: errcheck

			format := ""

			if len(args) > 1 {
				format = args[1]
			}

			resp, err := client.Download(context.Background(), &service.DownloadRequest{
				Id:     args[0],
				Format: format,
			})

			if err != nil {
				log.Println(err)

				return
			}

			name := "./" + resp.Id + "." + resp.Format
			err = ioutil.WriteFile(name, resp.Data, 0700)

			if err != nil {
				log.Println(err)

				return
			}

			log.Printf("Image downloaded to %q", name)
		},
	}
)

func main() {
	log.SetFlags(0)

	var (
		err        error
		configFile string
	)

	cmdImgAPI.PersistentFlags().StringVarP(&configFile, "config", "c", "./config.yml", "Configuration file location")

	config, err = service.LoadConfig(configFile)

	if err != nil {
		log.Fatal(err)
	}

	cmdImgAPI.AddCommand(
		cmdServer,
		cmdUpload,
		cmdDownload,
	)

	if err := cmdImgAPI.Execute(); err != nil {
		log.Fatalf("Execution error: %s", err)
	}
}
