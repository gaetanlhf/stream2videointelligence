package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	videointelligence "cloud.google.com/go/videointelligence/apiv1p3beta1"
	videointelligencepb "cloud.google.com/go/videointelligence/apiv1p3beta1/videointelligencepb"
	"github.com/gogo/protobuf/jsonpb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

var (
	version                                string
	buildTime                              string
	serviceAccountCredentialsPath          *string
	videoPath                              *string
	CloudVideoIntelligenceStreamingFeature *string
	GoogleCloudStoragePath                 *string
	isGoogleCloudStorageEnabled            bool
	enableStdout                           *bool
	exportPath                             *string
	n                                      int
	err                                    error
	videoFile                              *os.File
	exportFile                             *os.File
)

func init() {
	serviceAccountCredentialsPath = flag.String("creds", "", "Service account JSON key file path")
	videoPath = flag.String("source", "", "Using a file as a source instead of a pipe (not mandatory)")
	CloudVideoIntelligenceStreamingFeature = flag.String("feature", "", "API Cloud Video Intelligence streaming feature")
	GoogleCloudStoragePath = flag.String("gcs", "", "GCS URI to store all annotation results (not mandatory)")
	enableStdout = flag.Bool("stdout", false, "Print in stdout results from the API (not mandatory)")
	exportPath = flag.String("export", "", "Export the annotation results from the API to a file (not mandatory)")
}

func main() {
	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})
	flag.Parse()

	if len(*serviceAccountCredentialsPath) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	serviceAccountCredentialsFile, err := os.OpenFile(*serviceAccountCredentialsPath, os.O_RDONLY, 0600)

	if os.IsNotExist(err) {
		log.Fatalf("Error while opening your service account JSON key file : file '%s' does not exist", *serviceAccountCredentialsPath)
	} else if os.IsPermission(err) {
		log.Fatalf("Error while opening your service account JSON key file : insufficient permissions to open file '%s' : %s", *serviceAccountCredentialsPath, err)
	} else if err != nil {
		log.Fatalf("Error while opening your service account JSON key file '%s' : %err", *serviceAccountCredentialsPath, err)
	}

	serviceAccountCredentialsFile.Close()

	if len(*videoPath) != 0 {
		videoFile, err = os.OpenFile(*videoPath, os.O_RDONLY, 0600)

		if os.IsNotExist(err) {
			log.Fatalf("Error while opening your video file : file '%s' does not exist", *videoFile)
		} else if os.IsPermission(err) {
			log.Fatalf("Error while opening your video file : insufficient permissions to open file '%s' : %s", *videoFile, err)
		} else if err != nil {
			log.Fatalf("Error while opening your video file '%s' : %s", *videoFile, err)
		}
	}

	if len(*exportPath) != 0 {
		exportFile, err := os.OpenFile(*exportPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)

		if os.IsNotExist(err) {
			log.Fatalf("File '%s' does not exist", *exportPath)
		} else if os.IsPermission(err) {
			log.Fatalf("Insufficient permissions to open/create file '%s' for appending: %s", *exportPath, err)
		} else if err != nil {
			log.Fatalf("Error while opening file '%s' for writing: %s", *exportPath, err)
		}

		defer exportFile.Close()
	}

	if len(*CloudVideoIntelligenceStreamingFeature) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if !*enableStdout && len(*exportPath) == 0 && len(*GoogleCloudStoragePath) == 0 {
		log.Fatal("Nothing to do!")
	}

	if videointelligencepb.StreamingFeature(videointelligencepb.StreamingFeature_value[*CloudVideoIntelligenceStreamingFeature]) == 0 {
		log.Fatal("Invalid feature!")
	}

	if len(*GoogleCloudStoragePath) != 0 {
		isGoogleCloudStorageEnabled = true
		log.Info("Data will be exported to : ", *GoogleCloudStoragePath)
	}

	log.Info("Starting Cloud Video Intelligence API Streaming " + version + " build on " + buildTime)

	stream := initStreaming(*serviceAccountCredentialsPath)
	sendConfiguration(stream, *CloudVideoIntelligenceStreamingFeature, *GoogleCloudStoragePath, isGoogleCloudStorageEnabled)
	streamVideoToGCP(stream, *videoPath)

	for {
		resp, err := stream.Recv()

		if err != nil {
			log.Fatal("An error occured : ", err)
		}

		m := jsonpb.Marshaler{}
		results, _ := m.MarshalToString(resp)

		if len(*exportPath) != 0 {
			_, err = exportFile.WriteString(results)

			if err != nil {
				log.Fatal("An error occured while exporting : ", err)
			}
		}

		if *enableStdout {
			fmt.Println(results)
		}
	}
}

func initStreaming(sa string) videointelligencepb.StreamingVideoIntelligenceService_StreamingAnnotateVideoClient {
	ctx := context.Background()

	log.Info("Connecting to Video Intelligence API...")

	client, err := videointelligence.NewStreamingVideoIntelligenceClient(ctx, option.WithCredentialsFile(sa))

	if err != nil {
		log.Fatal(err)
	}

	log.Info("Successfully connected!")

	stream, err := client.StreamingAnnotateVideo(ctx)

	if err != nil {
		log.Fatal(err)
	}

	return stream
}

func sendConfiguration(stream videointelligencepb.StreamingVideoIntelligenceService_StreamingAnnotateVideoClient, feature string, gcs string, gcsEnabled bool) {
	log.Info("Sending configuration...")

	streamConfig := videointelligencepb.StreamingAnnotateVideoRequest{
		StreamingRequest: &videointelligencepb.StreamingAnnotateVideoRequest_VideoConfig{
			VideoConfig: &videointelligencepb.StreamingVideoConfig{
				StreamingConfig: nil,
				Feature:         videointelligencepb.StreamingFeature(videointelligencepb.StreamingFeature_value[feature]),
				StorageConfig: &videointelligencepb.StreamingStorageConfig{
					EnableStorageAnnotationResult:    gcsEnabled,
					AnnotationResultStorageDirectory: gcs,
				},
			},
		},
	}

	stream.Send(&streamConfig)

	log.Info("Configuration sended!")
}

func streamVideoToGCP(stream videointelligencepb.StreamingVideoIntelligenceService_StreamingAnnotateVideoClient, videoPath string) {
	data := make([]byte, 0, 1*1024*1024)

	log.Info("Sending data...")

	go func() {
		for {
			data = data[:cap(data)]

			if len(videoPath) == 0 {
				n, err = os.Stdin.Read(data)
			} else {
				n, err = (videoFile.Read(data))
			}

			if n == 0 {
				if err == nil {
					continue
				}

				if err == io.EOF {
					continue
				}

				log.Fatal(err)
			}

			data := data[:n]

			streamData := videointelligencepb.StreamingAnnotateVideoRequest{
				StreamingRequest: &videointelligencepb.StreamingAnnotateVideoRequest_InputContent{
					InputContent: data,
				},
			}

			stream.Send(&streamData)
		}
	}()
}
