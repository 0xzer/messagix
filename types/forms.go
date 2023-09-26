package types

type LoginForm struct {
	Jazoest                 string `url:"jazoest" name:"jazoest"`
	Lsd                     string `url:"lsd" name:"lsd"`
	Display                 string `url:"display" name:"display"`
	IsPrivate               string `url:"isprivate" name:"isprivate"`
	ReturnSession           string `url:"return_session" name:"return_session"`
	SkipAPILogin            string `url:"skip_api_login" name:"skip_api_login"`
	SignedNext              string `url:"signed_next" name:"signed_next"`
	TryNum                  string `url:"trynum" name:"trynum"`
	Timezone                string `url:"timezone"`
	Lgndim                  string `url:"lgndim"`
	Lgnrnd                  string `url:"lgnrnd" name:"lgnrnd"`
	Lgnjs                   string `url:"lgnjs"`
	Email                   string `url:"email"`
	PrefillContactPoint     string `url:"prefill_contact_point" name:"prefill_contact_point"`
	PrefillSource           string `url:"prefill_source" name:"prefill_source"`
	PrefillType             string `url:"prefill_type" name:"prefill_type"`
	FirstPrefillSource      string `url:"first_prefill_source" name:"first_prefill_source"`
	FirstPrefillType        string `url:"first_prefill_type" name:"first_prefill_type"`
	HadCPPrefilled          string `url:"had_cp_prefilled" name:"had_cp_prefilled"`
	HadPasswordPrefilled    string `url:"had_password_prefilled" name:"had_password_prefilled"`
	AbTestData              string `url:"ab_test_data"`
	EncPass                 string `url:"encpass"`
}

type LgnDim struct {
	W  int `json:"w,omitempty"`
	H  int `json:"h,omitempty"`
	Aw int `json:"aw,omitempty"`
	Ah int `json:"ah,omitempty"`
	C  int `json:"c,omitempty"`
}