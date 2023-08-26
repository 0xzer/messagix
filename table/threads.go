package table

type LSTruncateMetadataThreads struct {}

type LSUpsertInboxThreadsRange struct {
	SyncGroup int64 `index:"0"`
	MinLastActivityTimestampMs int64 `index:"1"`
	HasMoreBefore bool `index:"2"`
	IsLoadingBefore bool `index:"3"`
	MinThreadKey int64 `index:"4"`
}

type LSUpdateThreadsRangesV2 struct {
	FolderName string `index:"0"`
	ParentThreadKey int64 `index:"1"` /* not sure */
	MinLastActivityTimestampMs int64 `index:"2"`
	MinThreadKey int64 `index:"3"`
	IsLoadingBefore int64 `index:"4"` /* not sure */
}

type LSDeleteThenInsertThread struct {
    LastActivityTimestampMs int64 `index:"0"`
    LastReadWatermarkTimestampMs int64 `index:"1"`
    Snippet string `index:"2"`
    ThreadName string `index:"3"`
    ThreadPictureUrl string `index:"4"`
    NeedsAdminApprovalForNewParticipant bool `index:"5"`
    AuthorityLevel int64 `index:"6"`
    ThreadKey int64 `index:"7"`
    MailboxType int64 `index:"8"`
    ThreadType int64 `index:"9"`
    FolderName string `index:"10"`
    ThreadPictureUrlFallback string `index:"11"`
    ThreadPictureUrlExpirationTimestampMs int64 `index:"12"`
    RemoveWatermarkTimestampMs int64 `index:"13"`
    MuteExpireTimeMs int64 `index:"14"`
    MuteCallsExpireTimeMs int64 `index:"15"`
    GroupNotificationSettings int64 `index:"16"`
    IsAdminSnippet bool `index:"17"`
    SnippetSenderContactId int64 `index:"18"`
    SnippetStringHash string `index:"21"`
    SnippetStringArgument1 string `index:"22"`
    SnippetAttribution int64 `index:"23"`
    SnippetAttributionStringHash string `index:"24"`
    DisappearingSettingTtl int64 `index:"25"`
    DisappearingSettingUpdatedTs int64 `index:"26"`
    DisappearingSettingUpdatedBy int64 `index:"27"`
    OngoingCallState int64 `index:"29"`
    CannotReplyReason int64 `index:"30"`
    CustomEmoji int64 `index:"31"`
    CustomEmojiImageUrl string `index:"32"`
    OutgoingBubbleColor int64 `index:"33"`
    ThemeFbid int64 `index:"34"`
    ParentThreadKey int64 `index:"35"`
    NullstateDescriptionText1 string `index:"36"`
    NullstateDescriptionType1 int64 `index:"37"`
    NullstateDescriptionText2 string `index:"38"`
    NullstateDescriptionType2 int64 `index:"39"`
    NullstateDescriptionText3 string `index:"40"`
    NullstateDescriptionType3 int64 `index:"41"`
    DraftMessage string `index:"42"`
    SnippetHasEmoji bool `index:"43"`
    HasPersistentMenu bool `index:"44"`
    DisableComposerInput bool `index:"45"`
    CannotUnsendReason int64 `index:"46"`
    ViewedPluginKey int64 `index:"47"`
    ViewedPluginContext int64 `index:"48"`
    ClientThreadKey int64 `index:"49"`
    Capabilities int64 `index:"50"`
    ShouldRoundThreadPicture bool `index:"51"`
    ProactiveWarningDismissTime int64 `index:"52"`
    IsCustomThreadPicture bool `index:"53"`
    OtidOfFirstMessage int64 `index:"54"`
    NormalizedSearchTerms int64 `index:"55"`
    AdditionalThreadContext int64 `index:"56"`
    DisappearingThreadKey int64 `index:"57"`
    IsDisappearingMode bool `index:"58"`
    DisappearingModeInitiator int64 `index:"59"`
    UnreadDisappearingMessageCount int64 `index:"60"`
    LastMessageCtaId int64 `index:"62"`
    LastMessageCtaType int64 `index:"63"`
    ConsistentThreadFbid int64 `index:"64"`
    ThreadDescription int64 `index:"65"`
    UnsendLimitMs int64 `index:"66"`
    SyncGroup int64 `index:"67"`
    ThreadInvitesEnabled int64 `index:"68"`
    ThreadInviteLink string `index:"69"`
    NumUnreadSubthreads int64 `index:"70"`
    SubthreadCount int64 `index:"71"`
    ThreadInvitesEnabledV2 int64 `index:"72"`
    EventStartTimestampMs int64 `index:"73"`
    EventEndTimestampMs int64 `index:"74"`
    TakedownState int64 `index:"75"`
    MemberCount int64 `index:"76"`
    SecondaryParentThreadKey int64 `index:"77"`
    IgFolder int64 `index:"78"`
    InviterId int64 `index:"79"`
    ThreadTags int64 `index:"80"`
    ThreadStatus int64 `index:"81"`
    ThreadSubtype int64 `index:"82"`
    PauseThreadTimestamp int64 `index:"83"`
}

type LSAddParticipantIdToGroupThread struct {
    ThreadKey int64 `index:"0"`
    ContactId int64 `index:"1"`
    ReadWatermarkTimestampMs int64 `index:"2"`
    ReadActionTimestampMs int64 `index:"3"`
    DeliveredWatermarkTimestampMs int64 `index:"4"`
    Nickname string `index:"5"`
    IsAdmin bool `index:"6"`
    SubscribeSource int64 `index:"7"`
    AuthorityLevel int64 `index:"9"`
    NormalizedSearchTerms int64 `index:"10"`
    IsSuperAdmin bool `index:"11"`
    ThreadRoles int64 `index:"12"`
}

type LSWriteThreadCapabilities struct {
    ThreadKey int64 `index:"0"`
    Capabilities int64 `index:"1"`
    Capabilities2 int64 `index:"2"`
    Capabilities3 int64 `index:"3"`
}

type LSUpdateReadReceipt struct {
    ReadWatermarkTimestampMs int64 `index:"0"`
    ThreadKey int64 `index:"1"`
    ContactId int64 `index:"2"`
    ReadActionTimestampMs int64 `index:"3"`
}

type LSThreadsRangesQuery struct {
    ParentThreadKey int64 `index:"0"`
    Unknown1 bool `index:"1"`
    IsAfter bool `index:"2"`
    ReferenceThreadKey int64 `conditionField:"IsAfter" indexes:"4,3"`
    ReferenceActivityTimestamp int64 `conditionField:"IsAfter" indexes:"5,6"`
    AdditionalPagesToFetch int64 `index:"7"`
    Unknown8 bool `index:"8"`
}