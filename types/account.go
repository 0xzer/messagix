package types


type CurrentBusinessAccount struct {
	BusinessAccountName                              any    `json:"businessAccountName,omitempty"`
	BusinessID                                       any    `json:"business_id,omitempty"`
	BusinessPersonaID                                any    `json:"business_persona_id,omitempty"`
	BusinessProfilePicURL                            any    `json:"business_profile_pic_url,omitempty"`
	BusinessRole                                     any    `json:"business_role,omitempty"`
	BusinessUserID                                   any    `json:"business_user_id,omitempty"`
	Email                                            any    `json:"email,omitempty"`
	EnterpriseProfilePicURL                          any    `json:"enterprise_profile_pic_url,omitempty"`
	ExpiryTime                                       any    `json:"expiry_time,omitempty"`
	FirstName                                        any    `json:"first_name,omitempty"`
	HasVerifiedEmail                                 any    `json:"has_verified_email,omitempty"`
	IPPermission                                     any    `json:"ip_permission,omitempty"`
	IsBusinessPerson                                 bool   `json:"isBusinessPerson,omitempty"`
	IsEnterpriseBusiness                             bool   `json:"isEnterpriseBusiness,omitempty"`
	IsFacebookWorkAccount                            bool   `json:"isFacebookWorkAccount,omitempty"`
	IsInstagramBusinessPerson                        bool   `json:"isInstagramBusinessPerson,omitempty"`
	IsTwoFacNewFlow                                  bool   `json:"isTwoFacNewFlow,omitempty"`
	IsUserOptInAccountSwitchInfraUpgrade             bool   `json:"isUserOptInAccountSwitchInfraUpgrade,omitempty"`
	IsAdsFeatureLimited                              any    `json:"is_ads_feature_limited,omitempty"`
	IsBusinessBanhammered                            any    `json:"is_business_banhammered,omitempty"`
	LastName                                         any    `json:"last_name,omitempty"`
	PermittedBusinessAccountTaskIds                  []any  `json:"permitted_business_account_task_ids,omitempty"`
	PersonalUserID                                   string `json:"personal_user_id,omitempty"`
	ShouldHideComponentsByUnsupportedFirstPartyTools bool   `json:"shouldHideComponentsByUnsupportedFirstPartyTools,omitempty"`
	ShouldShowAccountSwitchComponents                bool   `json:"shouldShowAccountSwitchComponents,omitempty"`
}

type MessengerWebInitData struct {
	AccountKey      string          `json:"accountKey,omitempty"`
	AppID           int64           `json:"appId,omitempty"`
	CryptoAuthToken CryptoAuthToken `json:"cryptoAuthToken,omitempty"`
	LogoutToken     string          `json:"logoutToken,omitempty"`
	SessionID       string          `json:"sessionId,omitempty"`
}

type CryptoAuthToken struct {
	EncryptedSerializedCat  string `json:"encrypted_serialized_cat,omitempty"`
	ExpirationTimeInSeconds int    `json:"expiration_time_in_seconds,omitempty"`
}

type LSD struct {
	Token string `json:"token,omitempty"`
}

type IntlViewerContext struct {
	Gender         int `json:"GENDER,omitempty"`
	RegionalLocale any `json:"regionalLocale,omitempty"`
}

type IntlCurrentLocale struct {
	Code string `json:"code,omitempty"`
}

type DTSGInitData struct {
	AsyncGetToken string `json:"async_get_token,omitempty"`
	Token         string `json:"token,omitempty"`
}

type DTSGInitialData struct {
	Token string `json:"token,omitempty"`
}

type CurrentUserInitialData struct {
	AccountID                       string `json:"ACCOUNT_ID,omitempty"`
	AppID                           string `json:"APP_ID,omitempty"`
	HasSecondaryBusinessPerson      bool   `json:"HAS_SECONDARY_BUSINESS_PERSON,omitempty"`
	IsBusinessDomain                bool   `json:"IS_BUSINESS_DOMAIN,omitempty"`
	IsBusinessPersonAccount         bool   `json:"IS_BUSINESS_PERSON_ACCOUNT,omitempty"`
	IsDeactivatedAllowedOnMessenger bool   `json:"IS_DEACTIVATED_ALLOWED_ON_MESSENGER,omitempty"`
	IsFacebookWorkAccount           bool   `json:"IS_FACEBOOK_WORK_ACCOUNT,omitempty"`
	IsMessengerCallGuestUser        bool   `json:"IS_MESSENGER_CALL_GUEST_USER,omitempty"`
	IsMessengerOnlyUser             bool   `json:"IS_MESSENGER_ONLY_USER,omitempty"`
	IsWorkroomsUser                 bool   `json:"IS_WORKROOMS_USER,omitempty"`
	IsWorkMessengerCallGuestUser    bool   `json:"IS_WORK_MESSENGER_CALL_GUEST_USER,omitempty"`
	Name                            string `json:"NAME,omitempty"`
	ShortName                       string `json:"SHORT_NAME,omitempty"`
	UserID                          string `json:"USER_ID,omitempty"`
}