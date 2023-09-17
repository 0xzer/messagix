package table

type LSInsertStickerAttachment struct {
    PlayableUrl string `index:"0"`
    PlayableUrlFallback string `index:"1"`
    PlayableUrlExpirationTimestampMs int64 `index:"2"`
    PlayableUrlMimeType string `index:"3"`
    PreviewUrl string `index:"4"`
    PreviewUrlFallback string `index:"5"`
    PreviewUrlExpirationTimestampMs int64 `index:"6"`
    PreviewUrlMimeType string `index:"7"`
    PreviewWidth int64 `index:"9"`
    PreviewHeight int64 `index:"10"`
    ImageUrlMimeType string `index:"11"`
    AttachmentIndex int64 `index:"12"`
    AccessibilitySummaryText string `index:"13"`
    ThreadKey int64 `index:"14"`
    TimestampMs int64 `index:"17"`
    MessageId string `index:"18"`
    AttachmentFbid string `index:"19"`
    ImageUrl string `index:"20"`
    ImageUrlFallback string `index:"21"`
    ImageUrlExpirationTimestampMs int64 `index:"22"`
    FaviconUrlExpirationTimestampMs int64 `index:"23"`
    AvatarViewSize int64 `index:"25"`
    AvatarCount int64 `index:"26"`
    TargetId int64 `index:"27"`
    MustacheText string `index:"30"`
    AuthorityLevel int64 `index:"31"`
}