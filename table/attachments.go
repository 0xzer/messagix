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

type LSInsertXmaAttachment struct {
    Filename string `index:"1"`
    Filesize int64 `index:"2"`
    IsSharable bool `index:"3"`
    PlayableUrl string `index:"4"`
    PlayableUrlFallback string `index:"5"`
    PlayableUrlExpirationTimestampMs int64 `index:"6"`
    PlayableUrlMimeType string `index:"7"`
    PreviewUrl string `index:"8"`
    PreviewUrlFallback string `index:"9"`
    PreviewUrlExpirationTimestampMs int64 `index:"10"`
    PreviewUrlMimeType string `index:"11"`
    PreviewWidth int64 `index:"13"`
    PreviewHeight int64 `index:"14"`
    AttributionAppId int64 `index:"15"`
    AttributionAppName string `index:"16"`
    AttributionAppIcon int64 `index:"17"`
    AttributionAppIconFallback string `index:"18"`
    AttributionAppIconUrlExpirationTimestampMs int64 `index:"19"`
    AttachmentIndex int64 `index:"20"`
    AccessibilitySummaryText string `index:"21"`
    ShouldRespectServerPreviewSize bool `index:"22"`
    SubtitleIconUrl string `index:"23"`
    ShouldAutoplayVideo bool `index:"24"`
    ThreadKey int64 `index:"25"`
    AttachmentType int64 `index:"27"`
    TimestampMs int64 `index:"29"`
    MessageId string `index:"30"`
    OfflineAttachmentId int64 `index:"31"`
    AttachmentFbid string `index:"32"`
    XmaLayoutType int64 `index:"33"`
    XmasTemplateType int64 `index:"34"`
    CollapsibleId int64 `index:"35"`
    DefaultCtaId int64 `index:"36"`
    DefaultCtaTitle int64 `index:"37"`
    DefaultCtaType int64 `index:"38"`
    AttachmentCta1Id int64 `index:"40"`
    Cta1Title int64 `index:"41"`
    Cta1IconType int64 `index:"42"`
    Cta1Type int64 `index:"43"`
    AttachmentCta2Id int64 `index:"45"`
    Cta2Title int64 `index:"46"`
    Cta2IconType int64 `index:"47"`
    Cta2Type int64 `index:"48"`
    AttachmentCta3Id int64 `index:"50"`
    Cta3Title string `index:"51"`
    Cta3IconType int64 `index:"52"`
    Cta3Type int64 `index:"53"`
    ImageUrl int64 `string:"54"`
    ImageUrlFallback string `index:"55"`
    ImageUrlExpirationTimestampMs int64 `index:"56"`
    ActionUrl string `index:"57"`
    TitleText string `index:"58"`
    SubtitleText string `index:"59"`
    SubtitleDecorationType int64 `index:"60"`
    MaxTitleNumOfLines int64 `index:"61"`
    MaxSubtitleNumOfLines int64 `index:"62"`
    DescriptionText string `index:"63"`
    SourceText string `index:"64"`
    FaviconUrl string `index:"65"`
    FaviconUrlFallback string `index:"66"`
    FaviconUrlExpirationTimestampMs int64 `index:"67"`
    ListItemsId int64 `index:"69"`
    ListItemsDescriptionText string `index:"70"`
    ListItemsDescriptionSubtitleText string `index:"71"`
    ListItemsSecondaryDescriptionText string `index:"72"`
    ListItemId1 int64 `index:"73"`
    ListItemTitleText1 string `index:"74"`
    ListItemContactUrlList1 int64 `index:"75"`
    ListItemProgressBarFilledPercentage1 int64 `index:"76"`
    ListItemContactUrlExpirationTimestampList1 int64 `index:"77"`
    ListItemContactUrlFallbackList1 int64 `index:"78"`
    ListItemAccessibilityText1 string `index:"79"`
    ListItemTotalCount1 int64 `index:"80"`
    ListItemId2 int64 `index:"81"`
    ListItemTitleText2 int64 `index:"82"`
    ListItemContactUrlList2 int64 `index:"83"`
    ListItemProgressBarFilledPercentage2 int64 `index:"84"`
    ListItemContactUrlExpirationTimestampList2 int64 `index:"85"`
    ListItemContactUrlFallbackList2 int64 `index:"86"`
    ListItemAccessibilityText2 int64 `index:"87"`
    ListItemTotalCount2 int64 `index:"88"`
    ListItemId3 int64 `index:"89"`
    ListItemTitleText3 int64 `index:"90"`
    ListItemContactUrlList3 int64 `index:"91"`
    ListItemProgressBarFilledPercentage3 int64 `index:"92"`
    ListItemContactUrlExpirationTimestampList3 int64 `index:"93"`
    ListItemContactUrlFallbackList3 int64 `index:"94"`
    ListItemAccessibilityText3 int64 `index:"95"`
    ListItemTotalCount3 int64 `index:"96"`
    IsBorderless bool `index:"100"`
    HeaderImageUrlMimeType string `index:"101"`
    HeaderTitle string `index:"102"`
    HeaderSubtitleText string `index:"103"`
    HeaderImageUrl string `index:"104"`
    HeaderImageUrlFallback string `index:"105"`
    HeaderImageUrlExpirationTimestampMs int64 `index:"106"`
    PreviewImageDecorationType int64 `index:"107"`
    ShouldHighlightHeaderTitleInTitle bool `index:"108"`
    TargetId int64 `index:"109"`
    AttachmentLoggingType int64 `index:"112"`
    PreviewUrlLarge string `index:"114"`
    GatingType int64 `index:"115"`
    GatingTitle string `index:"116"`
    TargetExpiryTimestampMs int64 `index:"117"`
    CountdownTimestampMs int64 `index:"118"`
    ShouldBlurSubattachments int64 `index:"119"`
    VerifiedType int64 `index:"120"`
    CaptionBodyText int64 `index:"121"`
    IsPublicXma int64 `index:"122"`
    ReplyCount int64 `index:"123"`
    AuthorityLevel int64 `index:"124"`
}

type LSInsertBlobAttachment struct {
    Filename string `index:"0"`
    Filesize int64 `index:"1"`
    HasMedia bool `index:"2"`
    PlayableUrl string `index:"3"`
    PlayableUrlFallback string `index:"4"`
    PlayableUrlExpirationTimestampMs int64 `index:"5"`
    PlayableUrlMimeType string `index:"6"`
    DashManifest string `index:"7"`
    PreviewUrl string `index:"8"`
    PreviewUrlFallback string `index:"9"`
    PreviewUrlExpirationTimestampMs int64 `index:"10"`
    PreviewUrlMimeType string `index:"11"`
    MiniPreview int64 `index:"13"`
    PreviewWidth int64 `index:"14"`
    PreviewHeight int64 `index:"15"`
    AttributionAppId int64 `index:"16"`
    AttributionAppName int64 `index:"17"`
    AttributionAppIcon int64 `index:"18"`
    AttributionAppIconFallback int64 `index:"19"`
    AttributionAppIconUrlExpirationTimestampMs int64 `index:"20"`
    LocalPlayableUrl int64 `index:"21"`
    PlayableDurationMs int64 `index:"22"`
    AttachmentIndex int64 `index:"23"`
    AccessibilitySummaryText int64 `index:"24"`
    IsPreviewImage bool `index:"25"`
    OriginalFileHash string `index:"26"`
    ThreadKey int64 `index:"27"`
    AttachmentType int64 `index:"29"`
    TimestampMs int64 `index:"31"`
    MessageId string `index:"32"`
    OfflineAttachmentId string `index:"33"`
    AttachmentFbid string `index:"34"`
    HasXma bool `index:"35"`
    XmaLayoutType int64 `index:"36"`
    XmasTemplateType int64 `index:"37"`
    TitleText string `index:"38"`
    SubtitleText string `index:"39"`
    DescriptionText string `index:"40"`
    SourceText string `index:"41"`
    FaviconUrlExpirationTimestampMs int64 `index:"42"`
    IsBorderless bool `index:"44"`
    PreviewUrlLarge string `index:"45"`
    SamplingFrequencyHz int64 `index:"46"`
    WaveformData string `index:"47"`
    AuthorityLevel int64 `index:"48"`
}

type LSInsertAttachmentItem struct {
    AttachmentFbid string `index:"0"`
    AttachmentIndex int64 `index:"1"`
    ThreadKey int64 `index:"2"`
    MessageId string `index:"4"`
    OriginalPageSenderId int64 `index:"7"`
    TitleText int64 `index:"8"`
    SubtitleText int64 `index:"9"`
    PlayableUrl int64 `index:"12"`
    PlayableUrlFallback int64 `index:"13"`
    PlayableUrlExpirationTimestampMs int64 `index:"14"`
    PlayableUrlMimeType int64 `index:"15"`
    DashManifest int64 `index:"16"`
    PreviewUrl string `index:"17"`
    PreviewUrlFallback string `index:"18"`
    PreviewUrlExpirationTimestampMs int64 `index:"19"`
    PreviewUrlMimeType string `index:"20"`
    PreviewWidth int64 `index:"21"`
    PreviewHeight int64 `index:"22"`
    ImageUrl int64 `index:"23"`
    DefaultCtaId int64 `index:"24"`
    DefaultCtaTitle int64 `index:"25"`
    DefaultCtaType int64 `index:"26"`
    DefaultButtonType int64 `index:"28"`
    DefaultActionUrl int64 `index:"29"`
    DefaultActionEnableExtensions bool `index:"30"`
    DefaultWebviewHeightRatio int64 `index:"32"`
    AttachmentCta1Id int64 `index:"34"`
    Cta1Title int64 `index:"35"`
    Cta1IconType int64 `index:"36"`
    Cta1Type int64 `index:"37"`
    AttachmentCta2Id int64 `index:"39"`
    Cta2Title int64 `index:"40"`
    Cta2IconType int64 `index:"41"`
    Cta2Type int64 `index:"42"`
    AttachmentCta3Id int64 `index:"44"`
    Cta3Title int64 `index:"45"`
    Cta3IconType int64 `index:"46"`
    Cta3Type int64 `index:"47"`
    FaviconUrl int64 `index:"48"`
    FaviconUrlFallback int64 `index:"49"`
    FaviconUrlExpirationTimestampMs int64 `index:"50"`
    PreviewUrlLarge int64 `index:"51"`
}

type LSGetFirstAvailableAttachmentCTAID struct {}

type LSInsertAttachmentCta struct {
    CtaId int64 `index:"0"`
    AttachmentFbid string `index:"1"`
    AttachmentIndex int64 `index:"2"`
    ThreadKey int64 `index:"3"`
    MessageId string `index:"5"`
    Title string `index:"6"`
    Type_ string `index:"7"`
    PlatformToken int64 `index:"8"`
    ActionUrl string `index:"9"`
    NativeUrl string `index:"10"`
    UrlWebviewType int64 `index:"11"`
    ActionContentBlob int64 `index:"12"`
    EnableExtensions bool `index:"13"`
    ExtensionHeightType int64 `index:"14"`
    TargetId int64 `index:"15"`
}

type LSUpdateAttachmentItemCtaAtIndex struct {
    AttachmentFbid string `index:"0"`
    Unknown int64 `index:"1"`
    AttachmentCtaId int64 `index:"2"`
    CtaTitle string `index:"3"`
    CtaType string `index:"4"`
    Index int64 `index:"5"`
}

type LSUpdateAttachmentCtaAtIndexIgnoringAuthority struct {
    ThreadKey int64 `index:"0"`
    MessageId string `index:"1"`
    AttachmentFbid string `index:"2"`
    AttachmentCtaId int64 `index:"3"`
    CtaTitle string `index:"4"`
    CtaType string `index:"5"`
    Index int64 `index:"6"`
}

type LSHasMatchingAttachmentCTA struct {
    ThreadKey int64 `index:"0"`
    AttachmentFbid string `index:"1"`
}