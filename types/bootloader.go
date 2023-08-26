package types

type HandlePayload struct {
	Consistency Consistency            `json:"consistency,omitempty"`
	RsrcMap     map[string]RsrcDetails `json:"rsrcMap,omitempty"`
	CsrUpgrade  string                 `json:"csrUpgrade,omitempty"`
}

type Consistency struct {
	Rev int64 `json:"rev,omitempty"`
}

type RsrcDetails struct {
	Type string `json:"type,omitempty"`
	Src  string `json:"src,omitempty"`
	C    int64 `json:"c,omitempty"`
	Tsrc string `json:"tsrc,omitempty"`
	P    string `json:"p,omitempty"`
	M    string `json:"m,omitempty"`
}