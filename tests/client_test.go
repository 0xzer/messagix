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

			contacts, err := cli.Account.GetContacts(100)
			if err != nil {
				log.Fatalf("failed to get contacts: %e", err)
			}

			cli.Logger.Info().Any("data", contacts).Msg("Contacts Response")
		case *messagix.Event_PublishResponse:
			cli.Logger.Info().Any("evtData", evtData).Msg("Got new event!")
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