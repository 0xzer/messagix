package modules

import (
	"log"
	"strconv"
	"strings"
)

type BootLoaderConfig struct {
	BtCutoffIndex              int   `json:"btCutoffIndex,omitempty"`
	DeferBootloads             bool  `json:"deferBootloads,omitempty"`
	EarlyRequireLazy           bool  `json:"earlyRequireLazy,omitempty"`
	FastPathForAlreadyRequired bool  `json:"fastPathForAlreadyRequired,omitempty"`
	HypStep4                   bool  `json:"hypStep4,omitempty"`
	JsRetries                  []int `json:"jsRetries,omitempty"`
	JsRetryAbortNum            int   `json:"jsRetryAbortNum,omitempty"`
	JsRetryAbortTime           int   `json:"jsRetryAbortTime,omitempty"`
	PhdOn                      bool  `json:"phdOn,omitempty"`
	SilentDups                 bool  `json:"silentDups,omitempty"`
	Timeout                    int   `json:"timeout,omitempty"`
	TranslationRetries         []int `json:"translationRetries,omitempty"`
	TranslationRetryAbortNum   int   `json:"translationRetryAbortNum,omitempty"`
	TranslationRetryAbortTime  int   `json:"translationRetryAbortTime,omitempty"`
}

type Bootloader_HandlePayload struct {
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

func HandlePayload(payload interface{}, bootloaderConfig *BootLoaderConfig) error {
	var data *Bootloader_HandlePayload
	err := InterfaceToStructJSON(&payload, &data)
	if err != nil {
		return err
	}

	if data.CsrUpgrade != "" {
		CsrBitmap = append(CsrBitmap, parseCSRBit(data.CsrUpgrade)...)
	}

	if len(data.RsrcMap) > 0 {
		for _, v := range data.RsrcMap {
			shouldAdd := (bootloaderConfig.PhdOn && v.C == 2) || (!bootloaderConfig.PhdOn && v.C != 0)
			if shouldAdd {
				CsrBitmap = append(CsrBitmap, parseCSRBit(v.P)...)
			}
		}
	}

	return nil
}

// s always start with :
func parseCSRBit(s string) []int {
	bits := make([]int, 0)
	splitUp := strings.Split(s[1:], ",")
	for _, b := range splitUp {
		conv, err := strconv.ParseInt(b, 10, 32)
		if err != nil {
			log.Fatalf("failed to parse csrbit: %e", err)
		}
		if conv == 0 {
			continue
		}
		bits = append(bits, int(conv))
	}
	return bits
}