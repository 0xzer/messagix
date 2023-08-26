package table


type LSVerifyContactRowExists struct {
	ContactId int64 `index:"0"`
	Unknown1 int64 `index:"1"`
	ProfilePictureUrl string `index:"2"`
	Name string `index:"3"`
	ContactType int64 `index:"4"`
	ProfilePictureFallbackUrl string `index:"5"`
	Unknown6 int64 `index:"6"`
	Unknown7 int64 `index:"7"`
	IsMemorialized bool `index:"9"`
	BlockedByViewerStatus int64 `index:"11"`
	CanViewerMessage bool `index:"12"`
	AuthorityLevel int64 `index:"14"`
	Capabilities int64 `index:"15"`
	Capabilities2 int64 `index:"16"`
	Gender Gender `index:"18"`
	ContactViewerRelationship ContactViewerRelationship `index:"19"`
	SecondaryName string `index:"20"`
}