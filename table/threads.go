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
    NeedsAdminApprovalForNewParticipant int64 `index:"5"`
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
    SnippetStringHash int64 `index:"21"`
    SnippetStringArgument1 int64 `index:"22"`
    SnippetAttribution int64 `index:"23"`
    SnippetAttributionStringHash int64 `index:"24"`
    DisappearingSettingTtl int64 `index:"25"`
    DisappearingSettingUpdatedTs int64 `index:"26"`
    DisappearingSettingUpdatedBy int64 `index:"27"`
    OngoingCallState int64 `index:"29"`
    CannotReplyReason int64 `index:"30"`
    CustomEmoji int64 `index:"31"`
    CustomEmojiImageUrl int64 `index:"32"`
    OutgoingBubbleColor int64 `index:"33"`
    ThemeFbid int64 `index:"34"`
    ParentThreadKey int64 `index:"35"`
    NullstateDescriptionText1 string `index:"36"`
    NullstateDescriptionType1 int64 `index:"37"`
    NullstateDescriptionText2 string `index:"38"`
    NullstateDescriptionType2 int64 `index:"39"`
    NullstateDescriptionText3 string `index:"40"`
    NullstateDescriptionType3 int64 `index:"41"`
    SnippetHasEmoji bool `index:"42"`
    HasPersistentMenu bool `index:"43"`
    DisableComposerInput bool `index:"44"`
    CannotUnsendReason int64 `index:"45"`
    ViewedPluginKey int64 `index:"46"`
    ViewedPluginContext int64 `index:"47"`
    ClientThreadKey int64 `index:"48"`
    Capabilities int64 `index:"49"`
    ShouldRoundThreadPicture int64 `index:"50"`
    ProactiveWarningDismissTime int64 `index:"51"`
    IsCustomThreadPicture bool `index:"52"`
    OtidOfFirstMessage int64 `index:"53"`
    NormalizedSearchTerms int64 `index:"54"`
    AdditionalThreadContext int64 `index:"55"`
    DisappearingThreadKey int64 `index:"56"`
    IsDisappearingMode bool `index:"57"`
    DisappearingModeInitiator int64 `index:"58"`
    UnreadDisappearingMessageCount int64 `index:"59"`
    LastMessageCtaId int64 `index:"61"`
    LastMessageCtaType int64 `index:"62"`
    ConsistentThreadFbid int64 `index:"63"`
    ThreadDescription int64 `index:"64"`
    UnsendLimitMs int64 `index:"65"`
    SyncGroup int64 `index:"66"`
    ThreadInvitesEnabled int64 `index:"67"`
    ThreadInviteLink int64 `index:"68"`
    NumUnreadSubthreads int64 `index:"69"`
    SubthreadCount int64 `index:"70"`
    ThreadInvitesEnabledV2 int64 `index:"71"`
    EventStartTimestampMs int64 `index:"72"`
    EventEndTimestampMs int64 `index:"73"`
    TakedownState int64 `index:"74"`
    MemberCount int64 `index:"75"`
    SecondaryParentThreadKey int64 `index:"76"`
    IgFolder int64 `index:"77"`
    InviterId int64 `index:"78"`
    ThreadTags int64 `index:"79"`
    ThreadStatus int64 `index:"80"`
    ThreadSubtype int64 `index:"81"`
    PauseThreadTimestamp int64 `index:"82"`
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

type LSUpdateTypingIndicator struct {
    ThreadKey int64 `index:"0"`
    SenderId int64 `index:"1"`
    IsTyping bool `index:"2"`
}

type LSMoveThreadToInboxAndUpdateParent struct {
    ThreadKey int64 `index:"0"`
    ParentThreadKey int64 `index:"1"`
}

type LSUpdateThreadSnippet struct {
    ThreadKey int64 `index:"0"`
    Snippet string `index:"1"`
    IsAdminSnippet bool `index:"2"`
    SnippetSenderContactId int64 `index:"3"`
    SnippetHasEmoji bool `index:"4"`
    ViewedPluginKey string `index:"5"`
    ViewedPluginContext string `index:"6"`
}

type LSVerifyThreadExists struct {
    ThreadKey int64 `index:"0"`
    ThreadType int64 `index:"1"`
    FolderName string `index:"2"`
    ParentThreadKey int64 `index:"3"`
    AuthorityLevel int64 `index:"4"`
}

type LSBumpThread struct {
    LastReadWatermarkTimestampMs int64 `index:"0"`
    Unknown1 int64 `index:"1"`
    ThreadKey int64 `index:"2"`
}

// Idk which snippet is the correct, there's like 6 (snippet, snippetStringHash, snippetStringArgument1, snippetAttribution, snippetAttributionStringHash)
type LSUpdateThreadSnippetFromLastMessage struct {
    AccountId int64 `index:"0"`
    ThreadKey int64 `index:"1"`
    Snippet1 string `index:"2"`
    Snippet2 string `index:"3"`
    Snippet3 string `index:"4"`
    Snippet4 string `index:"5"`
    Snippet5 string `index:"6"`
    Snippet6 string `index:"7"`
    Snippet7 string `index:"8"`
    Snippet8 string `index:"9"`
    Snippet9 string `index:"10"`
    IsAdminSnippet bool `index:"11"`
}

type LSDeleteBannersByIds struct {
    ThreadKey int64 `index:"0"`
}

type LSUpdateDeliveryReceipt struct {
    DeliveredWatermarkTimestampMs int64 `index:"0"`
    ThreadKey int64 `index:"1"`
    ContactId int64 `index:"2"`
}

type LSUpdateOptimisticContextThreadKeys struct {
    ThreadKey1 int64 `index:"0"`
    ThreadKey2 int64 `index:"1"`
}

type LSReplaceOptimisticThread struct {
    ThreadKey1 int64 `index:"0"`
    ThreadKey2 int64 `index:"1"`
}

type LSApplyNewGroupThread struct {
    OtidOfFirstMessage string `index:"0"`
    ThreadKey int64 `index:"1"`
    ThreadType int64 `index:"2"`
    FolderName string `index:"3"`
    ParentThreadKey int64 `index:"4"`
    ThreadPictureUrlFallback string `index:"5"`
    LastActivityTimestampMs int64 `index:"6"`
    LastReadWatermarkTimestampMs int64 `index:"6"`
    NullstateDescriptionText1 string `index:"8"`
    NullstateDescriptionType1 int64 `index:"9"`
    NullstateDescriptionText2 string `index:"10"`
    NullstateDescriptionType2 int64 `index:"11"`
    CannotUnsendReason int64 `index:"12"`
    Capabilities int64 `index:"13"`
    InviterId int64 `index:"14"`
    IgFolder int64 `index:"15"`
    ThreadSubtype int64 `index:"16"`
}

type LSRemoveAllParticipantsForThread struct {
    ThreadKey int64 `index:"0"`
}

type LSUpdateThreadInviteLinksInfo struct {
    ThreadKey int64 `index:"0"`
    ThreadInvitesEnabled int64 `index:"1"` // 0 or 1
    ThreadInviteLink string `index:"2"`
}

type LSUpdateThreadParticipantAdminStatus struct {
    ThreadKey int64 `index:"0"`
    ContactId int64 `index:"1"`
    IsAdmin bool `index:"2"`
}

type LSUpdateParticipantSubscribeSourceText struct {
    ThreadKey int64 `index:"0"`
    ContactId int64 `index:"1"`
    SubscribeSource string `index:"2"`
}

type LSOverwriteAllThreadParticipantsAdminStatus struct {
    ThreadKey int64 `index:"0"`
    IsAdmin bool `index:"1"`
}

type LSUpdateParticipantCapabilities struct {
    ContactId int64 `index:"0"`
    ThreadKey int64 `index:"1"`
}

type LSChangeViewerStatus struct {
    ThreadKey int64 `index:"0"`
    CannotReplyReason string `index:"1"`
}

type LSSyncUpdateThreadName struct {
    ThreadName string `index:"0"`
    ThreadKey int64 `index:"1"`
    ThreadName1 string `index:"2"`
}