package messagix_test

import (
	"log"
	"os"
	"testing"
	"github.com/0xzer/messagix"
	"github.com/0xzer/messagix/cookies"
	"github.com/0xzer/messagix/debug"
	"github.com/0xzer/messagix/table"
	"github.com/0xzer/messagix/types"
)

var cli *messagix.Client
func TestClient(t *testing.T) {
	session := &cookies.InstagramCookies{}
	err := cookies.NewCookiesFromString(`dpr=1.125; mid=ZYr4EQAEAAFtXZA_PWwg-dP-uF-y; ig_did=683F6ACB-8EFF-4282-80B0-53612B412C70; datr=UPiKZSanJbNB_rw-0iDBzL8d; csrftoken=VTMs6VGc9iS0lgQXoHEGtsaxfnhI4XNe; ds_user_id=62049897018; sessionid=62049897018%3AuRYpJ87BMghOr2%3A10%3AAYfDs2B0WgKmxOomlmCG-XMYdsnGH1uN1ybbAJVH7w; shbid="11985\05462049897018\0541735142373:01f77e690d77d926c0a4b80f5747a162216d54089c7653d41f155db36ab7b1d20db7d984"; shbts="1703606373\05462049897018\0541735142373:01f7a271224b9a153008890fbd292a2936aa87a8a21581b208d585d02acbf4538390e0e9"; fbm_124024574287414=base_domain=.instagram.com; fbsr_124024574287414=YkNRpB4pracPA0zmCcwK47TDPQvIQVqafE1i54l_zvk.eyJ1c2VyX2lkIjoiMTAwMDUzNDQwNTA3MDAzIiwiY29kZSI6IkFRQlRRNDFrOG9wR211SlgzLUVUYkdnNkZ0NUtxVWlpaUtQbFVMV190ZWRNQXltQWlLck5xMDFnUTl4UnVDYXozOVNzRnNITGxxNHpKMlhDcjNiYXdvNTNpN0dEMWdCR3VqS3NMaDFwS2RZWVZKRGRGejZoakRDNk9QTDVfR3JkOTM1WkRNaW15N0tfSEg5blJhdWZ3blhqNHVPcm1NSGd2ZEhsYnQ1MDFkckNFTDMteWkwbWJXZ1NYaEhZVU9SQVM1MGpCbHpuNjZScTQ4S2NVWmhpVUl5Ylp5QUMwaERfRlAxU3k4QnJ5Tzg3WHg1RHVlOVdxZDNXdGhBMktFNVRaX0NHYXBsZ2xKMWFTYjExcFB3cmNKZE5ZWEJfclN3N3pSak5TODloZm84MWFXS2pkZk5TWVV2ekdMR0Z2Q3ZyTC1tS0NaVDV2LVJ6b2ExSjJROFhZNm0yIiwib2F1dGhfdG9rZW4iOiJFQUFCd3pMaXhuallCT3kwNDl1RG51TkF4UnZ2WkNGSlpBUlVEQnJaQjhTSUNJMjBNeTk1TTBxRFZqT0xOUjA1QmdES0dXN1pDRlE2eEk3UjRzVTZCNXdWYjRhYWdMR2FUSnR4ODg2ODl6eFVvQ1haQVJpT2FqMzY5RG9rOGxRYU1sOWdBOWlkbEhMb0dXcmZSNlZ4SkhtcWZsbWlGblBBb0NWWHc0S1pDdEpZR0lFeWl2cFFHb3EzZWdUV1V3NXkxSkp2WkNCa1pCdmw2VnhBWkQiLCJhbGdvcml0aG0iOiJITUFDLVNIQTI1NiIsImlzc3VlZF9hdCI6MTcwMzYwNjcyMn0; fbsr_124024574287414=YkNRpB4pracPA0zmCcwK47TDPQvIQVqafE1i54l_zvk.eyJ1c2VyX2lkIjoiMTAwMDUzNDQwNTA3MDAzIiwiY29kZSI6IkFRQlRRNDFrOG9wR211SlgzLUVUYkdnNkZ0NUtxVWlpaUtQbFVMV190ZWRNQXltQWlLck5xMDFnUTl4UnVDYXozOVNzRnNITGxxNHpKMlhDcjNiYXdvNTNpN0dEMWdCR3VqS3NMaDFwS2RZWVZKRGRGejZoakRDNk9QTDVfR3JkOTM1WkRNaW15N0tfSEg5blJhdWZ3blhqNHVPcm1NSGd2ZEhsYnQ1MDFkckNFTDMteWkwbWJXZ1NYaEhZVU9SQVM1MGpCbHpuNjZScTQ4S2NVWmhpVUl5Ylp5QUMwaERfRlAxU3k4QnJ5Tzg3WHg1RHVlOVdxZDNXdGhBMktFNVRaX0NHYXBsZ2xKMWFTYjExcFB3cmNKZE5ZWEJfclN3N3pSak5TODloZm84MWFXS2pkZk5TWVV2ekdMR0Z2Q3ZyTC1tS0NaVDV2LVJ6b2ExSjJROFhZNm0yIiwib2F1dGhfdG9rZW4iOiJFQUFCd3pMaXhuallCT3kwNDl1RG51TkF4UnZ2WkNGSlpBUlVEQnJaQjhTSUNJMjBNeTk1TTBxRFZqT0xOUjA1QmdES0dXN1pDRlE2eEk3UjRzVTZCNXdWYjRhYWdMR2FUSnR4ODg2ODl6eFVvQ1haQVJpT2FqMzY5RG9rOGxRYU1sOWdBOWlkbEhMb0dXcmZSNlZ4SkhtcWZsbWlGblBBb0NWWHc0S1pDdEpZR0lFeWl2cFFHb3EzZWdUV1V3NXkxSkp2WkNCa1pCdmw2VnhBWkQiLCJhbGdvcml0aG0iOiJITUFDLVNIQTI1NiIsImlzc3VlZF9hdCI6MTcwMzYwNjcyMn0; rur="CLN\05462049897018\0541735143095:01f73951f072d104adde189ba2b9666b0d6097ca0161c11117e712b2a6c9367fbc0b446e"`, session)
	if err != nil {
		log.Fatalf("failed to create insta cookies: %e", err)
	}

	cli, err = messagix.NewClient(types.Instagram, session, debug.NewLogger(), "")
	if err != nil {
		log.Fatal(err)
	}
	cli.SetEventHandler(evHandler)

	err = cli.Connect()
	if err != nil {
		log.Fatalf("failed to connect to socket: %e", err)
	}

	cli.SaveSession("test_files/session.json")
	// making sure the main program does not exit so that the socket can continue reading
	wait := make(chan struct{})
    <-wait
}

func evHandler(evt interface{}) {
	switch evtData := evt.(type) {
		case *messagix.Event_Ready:
			cli.Logger.Info().
			//Any("threads", evtData.Table.LSDeleteThenInsertThread).
			// Any("setcontentdisplay", evtData.Table.LSSetMessageDisplayedContentTypes).
			Msg("Client is ready!")
			threads := evtData.Table.LSDeleteThenInsertThread
			for _, thread := range threads {
				cli.Logger.Info().Any("thread_name", thread.ThreadName).Any("thread_key", thread.ThreadKey).Msg("[READY] Got thread info!")
			}
		case *messagix.Event_PublishResponse:
			// cli.Logger.Info().Any("tableData", evtData.Table).Msg("Received new event from socket")
			threads := evtData.Table.LSDeleteThenInsertThread
			for _, thread := range threads {
				cli.Logger.Info().Any("thread_name", thread.ThreadName).Any("thread_key", thread.ThreadKey).Msg("[PUBLISH] Got thread info!")
			}
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

func fetchMessages() {
	currentCursor := cli.SyncManager.GetCursor(1)
	resp, err := cli.Threads.FetchMessages(61550046156682, 1694956983720, "mid.$cAABtBT1ku_GQ1JRdqGKo08VYGlmT", currentCursor)
	if err != nil {
		log.Fatalf("failed to fetch messages: %e", err)
	}
	cli.Logger.Debug().Any("resp", resp).Msg("fetch messages")
}

func sendReaction() {
	resp, err := cli.Messages.SendReaction(11111, "mid", "👇")
	if err != nil {
		log.Fatalf("failed to send reaction: %e", err)
	}
	log.Println(resp.LSReplaceOptimisticReaction)
}

func sendMessageWithMedia() {
	mediaData, _ := os.ReadFile("test_files/test.jpeg")
	turtleData, _ := os.ReadFile("test_files/turtle.jpeg")

	medias := []*messagix.MercuryUploadMedia{
		{Filename: "test_image.jpg", MediaType: messagix.IMAGE_JPEG, MediaData: mediaData},
		{Filename: "turtle.jpg", MediaType: messagix.IMAGE_JPEG, MediaData: turtleData},
	}

	mediaUploads, err := cli.SendMercuryUploadRequest(medias)
	if err != nil {
		log.Fatalf("failed: %e", err)
	}

	cli.Logger.Info().Any("uploads", mediaUploads).Msg("Media uploads")

	sendMsgBuilder := cli.Threads.NewMessageBuilder(1111111)
	sendMsgBuilder.SetMedias(mediaUploads)
	sendMsgBuilder.SetText("media test :D")
	sendMsgBuilder.SetLastReadWatermarkTs(1696261558117)

	tableResp, err := sendMsgBuilder.Execute() // make sure to execute to send the task
	if err != nil {
		log.Fatalf("failed to send media: %e", err)
	}

	log.Println(tableResp)
}

func sendMessageText() {
	msgBuilder := cli.Threads.NewMessageBuilder(11111)
	msgBuilder.SetInitiatingSource(table.FACEBOOK_INBOX)
	msgBuilder.SetLastReadWatermarkTs(16962611558117)
	msgBuilder.SetSource(table.MESSENGER_INBOX_IN_THREAD)
	msgBuilder.SetText("hello there")
	tableResp, err := msgBuilder.Execute()
	if err != nil {
		log.Fatalf("failed to send text msg: %e", err)
	}
	log.Println(tableResp)
}