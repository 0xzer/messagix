package modules

type SiteData struct {
	SpinB                 string `json:"__spin_b,omitempty"`
	SpinR                 int    `json:"__spin_r,omitempty"`
	SpinT                 int    `json:"__spin_t,omitempty"`
	BeOneAhead            bool   `json:"be_one_ahead,omitempty"`
	BlHashVersion         int    `json:"bl_hash_version,omitempty"`
	ClientRevision        int    `json:"client_revision,omitempty"`
	CometEnv              int    `json:"comet_env,omitempty"`
	HasteSession          string `json:"haste_session,omitempty"`
	HasteSite             string `json:"haste_site,omitempty"`
	Hsi                   string `json:"hsi,omitempty"`
	IsComet               bool   `json:"is_comet,omitempty"`
	IsExperimentalTier    bool   `json:"is_experimental_tier,omitempty"`
	IsJitWarmedUp         bool   `json:"is_jit_warmed_up,omitempty"`
	IsRtl                 bool   `json:"is_rtl,omitempty"`
	ManifestBaseURI       string `json:"manifest_base_uri,omitempty"`
	ManifestOrigin        string `json:"manifest_origin,omitempty"`
	ManifestVersionPrefix string `json:"manifest_version_prefix,omitempty"`
	PkgCohort             string `json:"pkg_cohort,omitempty"`
	Pr                    int    `json:"pr,omitempty"`
	PushPhase             string `json:"push_phase,omitempty"`
	SemrHostBucket        string `json:"semr_host_bucket,omitempty"`
	ServerRevision        int    `json:"server_revision,omitempty"`
	SkipRdBl              bool   `json:"skip_rd_bl,omitempty"`
	Spin                  int    `json:"spin,omitempty"`
	Tier                  string `json:"tier,omitempty"`
	Vip                   string `json:"vip,omitempty"`
	WbloksEnv             bool   `json:"wbloks_env,omitempty"`
}