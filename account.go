package messagix

import (
	"fmt"
	"log"
	"reflect"
	"github.com/0xzer/messagix/crypto"
	"github.com/0xzer/messagix/modules"
	"github.com/0xzer/messagix/socket"
	"github.com/0xzer/messagix/table"
	"github.com/0xzer/messagix/types"
	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
)

var BASE_URL = "https://www.facebook.com"

type Account struct {
	client *Client
}

func (a *Account) Login(email, password string) (*types.Cookies, error) {
	moduleLoader := &ModuleParser{client: a.client}
	moduleLoader.load(BASE_URL+"/login")
	
	loginFormTags := moduleLoader.FormTags[0]
	loginPath, ok := loginFormTags.Attributes["action"]
	if !ok {
		return nil, fmt.Errorf("failed to resolve login path / action from html form tags")
	}

	err := a.client.configs.SetupConfigs()
	if err != nil {
		return nil, fmt.Errorf("failed to setup configs after moduleLoader finished: %e", err)
	}

	loginInputs := append(loginFormTags.Inputs, moduleLoader.LoginInputs...)
	loginForm := types.LoginForm{}
	v := reflect.ValueOf(&loginForm).Elem()
	a.client.configs.ParseFormInputs(loginInputs, v)

	a.client.configs.siteConfig.Jazoest = loginForm.Jazoest

	needsCookieConsent := len(modules.SchedulerJSDefined.InitialCookieConsent.InitialConsent) == 0
	if needsCookieConsent {
		err = a.client.sendCookieConsent(moduleLoader.JSDatr)
		if err != nil {
			return nil, err
		}
	}

	testDataSimulator := crypto.NewABTestData()
	data := testDataSimulator.GenerateAbTestData([]string{email, password})
	
	encryptedPW, err := crypto.EncryptPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt password: %e", err)
	}

	loginForm.Email = email
	loginForm.EncPass = encryptedPW
	loginForm.AbTestData = data
	loginForm.Lgndim = "eyJ3IjoyMjc1LCJoIjoxMjgwLCJhdyI6MjI3NiwiYWgiOjEyMzIsImMiOjI0fQ==" // irrelevant
	loginForm.Lgnjs = a.client.configs.siteConfig.SpinT
	loginForm.Timezone = "-120"

	form, err := query.Values(&loginForm)
	if err != nil {
		return nil, err
	}

	payloadBytes := []byte(form.Encode())
	h := a.client.buildHeaders()
	h.Add("sec-fetch-dest", "document")
	h.Add("sec-fetch-mode", "navigate")
	h.Add("sec-fetch-site", "same-origin") // header is required
	h.Add("sec-fetch-user", "?1")
	h.Add("upgrade-insecure-requests", "1")
	h.Add("origin", "https://www.facebook.com")
	h.Add("referer", "https://www.facebook.com/login")

	a.client.disableRedirects()
	resp, _, err := a.client.MakeRequest(BASE_URL + loginPath, "POST", h, payloadBytes, types.FORM)
	if err != nil {
		log.Fatalf("failed to send login request: %e", err)
	}
	a.client.enableRedirects()
	cookies := resp.Cookies()

	hasUserCookie := a.client.findCookie(cookies, "c_user")
	if hasUserCookie == nil {
		return nil, fmt.Errorf("failed to login")
	}
	
	session := types.NewCookiesFromResponse(cookies)

	a.client.cookies = session
	return session, nil
}

func (a *Account) GetContacts(limit int64) ([]table.LSVerifyContactRowExists, error) {
	tskm := a.client.NewTaskManager()
	tskm.AddNewTask(&socket.GetContactsTask{Limit: limit})

	payload, err := tskm.FinalizePayload()
	if err != nil {
		log.Fatal(err)
	}

	packetId, err := a.client.socket.makeLSRequest(payload, 3)
	if err != nil {
		log.Fatal(err)
	}

	resp := a.client.socket.responseHandler.waitForPubResponseDetails(packetId)
	if resp == nil {
		return nil, fmt.Errorf("failed to receive response from socket while trying to fetch contacts. packetId: %d", packetId)
	}

	return resp.Table.LSVerifyContactRowExists, nil
}

func (a *Account) GetContactsFull(contactIds []int64) ([]table.LSDeleteThenInsertContact, error) {
	tskm := a.client.NewTaskManager()
	for _, id := range contactIds {
		tskm.AddNewTask(&socket.GetContactsFullTask{
			ContactId: id,
		})
	}

	payload, err := tskm.FinalizePayload()
	if err != nil {
		log.Fatal(err)
	}

	packetId, err := a.client.socket.makeLSRequest(payload, 3)
	if err != nil {
		log.Fatal(err)
	}

	resp := a.client.socket.responseHandler.waitForPubResponseDetails(packetId)
	if resp == nil {
		return nil, fmt.Errorf("failed to receive response from socket while trying to fetch full contact information. packetId: %d", packetId)
	}

	return resp.Table.LSDeleteThenInsertContact, nil
}

func (a *Account) ReportAppState(state table.AppState) error {
	tskm := a.client.NewTaskManager()
	tskm.AddNewTask(&socket.ReportAppStateTask{AppState: state, RequestId: uuid.NewString()})

	payload, err := tskm.FinalizePayload()
	if err != nil {
		log.Fatal(err)
	}

	packetId, err := a.client.socket.makeLSRequest(payload, 3)
	if err != nil {
		log.Fatal(err)
	}


	resp := a.client.socket.responseHandler.waitForPubResponseDetails(packetId)
	if resp == nil {
		return fmt.Errorf("failed to receive response from socket while trying to report app state. packetId: %d", packetId)
	}

	a.client.Logger.Info().Any("data", resp).Msg("Got report app state resp app state")
	return nil
}