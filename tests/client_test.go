package client_test

import (
	"log"
	"os"
	"testing"

	"github.com/0xzer/messagix"
	"github.com/0xzer/messagix/debug"
	"github.com/0xzer/messagix/types"
)

var cli *messagix.Client

func TestClient(t *testing.T) {
	cookies := types.NewCookiesFromString("")

	cli = messagix.NewClient(cookies, debug.NewLogger(), "")
	cli.SetEventHandler(evHandler)

	err := cli.Connect()
	if err != nil {
		log.Fatal(err)
	}

	err = cli.SaveSession("session.json")
	if err != nil {
		log.Fatal(err)
	}

	// making sure the main program does not exit so that the socket can continue reading
	wait := make(chan struct{})
    <-wait
}

func evHandler(evt interface{}) {
	switch evtData := evt.(type) {
		case *messagix.Event_Ready:
			cli.Logger.Info().
			Any("connectionCode", evtData.ConnectionCode.ToString()).
			Any("isNewSession", evtData.IsNewSession).
			Any("total_messages", len(evtData.Messages)).
			Any("total_threads", len(evtData.Threads)).
			Any("total_contacts", len(evtData.Contacts)).
			Msg("Client is ready!")
			
		case *messagix.Event_PublishResponse:
			cli.Logger.Info().Any("tableData", evtData.Table).Msg("Received new event from socket")
		case *messagix.Event_Error:
			cli.Logger.Err(evtData.Err).Msg("The library encountered an error")
			os.Exit(1)
		case *messagix.Event_SocketClosed:
			cli.Logger.Info().Any("code", evtData.Code).Any("text", evtData.Text).Msg("Socket was closed.")
			os.Exit(1)
		default:
			cli.Logger.Info().Any("data", evtData).Interface("type", evt).Msg("Got unknown event!")
	}
}

func sendReaction() {
	resp, err := cli.Messages.SendReaction(123123123, "messageid", "")
	if err != nil {
		log.Fatalf("failed to send reaction: %e", err)
	}
	log.Println(resp.LSReplaceOptimisticReaction)
}

func sendMessageWithMedia() {
	mediaData, _ := os.ReadFile("test.jpeg")
	turtleData, _ := os.ReadFile("turtle.jpeg")

	medias := []*messagix.MercuryUploadMedia{
		{Filename: "test_image.jpg", MediaType: messagix.IMAGE_JPEG, MediaData: mediaData},
		{Filename: "turtle.jpg", MediaType: messagix.IMAGE_JPEG, MediaData: turtleData},
	}

	mediaUploads, err := cli.SendMercuryUploadRequest(medias)
	if err != nil {
		log.Fatalf("failed: %e", err)
	}

	mediaIds := make([]int64, 0)
	for _, m := range mediaUploads {
		switch data := m.Payload.Metadata.(type) {
		case *types.ImageMetadata:
			mediaIds = append(mediaIds, data.Fbid)
		case *types.VideoMetadata:
			mediaIds = append(mediaIds, data.VideoID)
		}
	}

	cli.Logger.Info().Any("mediaIds", mediaIds).Msg("Sending mediaIds")
	sendMsgBuilder := cli.Threads.NewMessageBuilder(1231231323)
	sendMsgBuilder.SetMediaIDs(mediaIds)
	sendMsgBuilder.SetText("media test")
	sendMsgBuilder.SetLastReadWatermarkTs(1695421957450)

	sendMsgBuilder.Execute() // make sure to execute to send the task
}