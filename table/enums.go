package table

type DisplayedContentTypes int64
const (
	TEXT DisplayedContentTypes = 1
)

type Gender int64
const (
	NOT_A_PERSON Gender = 0
	FEMALE_SINGULAR Gender = 1
	MALE_SINGULAR Gender = 2
	FEMALE_SINGULAR_GUESS Gender = 3
	MALE_SINGULAR_GUESS Gender = 4
	MIXED_UNKNOWN Gender = 5
	NEUTER_SINGULAR Gender = 6
	UNKNOWN_SINGULAR Gender = 7
	FEMALE_PLURAL Gender = 8
	MALE_PLURAL Gender = 9
	NEUTER_PLURAL Gender = 10
	UNKNOWN_PLURAL Gender = 11
)

type ContactViewerRelationship int64
const (
	UNKNOWN_RELATIONSHIP ContactViewerRelationship = 0
	NOT_CONTACT ContactViewerRelationship = 1
	CONTACT_OF_VIEWER ContactViewerRelationship = 2
	FACEBOOK_FRIEND ContactViewerRelationship = 3
	SOFT_CONTACT ContactViewerRelationship = 4
)

type ThreadSourceType int64

const (
   /* communityFriendsDialog, pagesHomeFriendsDialog, mutualFriendsDialog, birthday, groupMembers, fundraiserSupportersList, memories, feedPoll, reactorList, friendsList, pagesPrivateReply, timeline, feedOrganicPost */
   FB_FEED_ORGANIC_POST ThreadSourceType = 1572865
   /* inboxPendingRequests */
   MESSENGER_INBOX_PENDING_REQUESTS ThreadSourceType = 65546
   /* fullscreenChat */
   MESSENGER_COMMUNITY_MESSAGING_FULLSCREEN_CHAT ThreadSourceType = 2293762
   /* sidebarGroupsList */
   SIDEBAR_CONTACTS_GROUPS ThreadSourceType = 2228227
   /* jewel */
   JEWEL_THREAD_LIST ThreadSourceType = 2097153
   /* shop */
   MINI_SHOP_VIEW_MENU_BUTTON ThreadSourceType = 2818048
   /* chatheadsOverflow */
   CHATHEADS_OVERFLOW ThreadSourceType = 2162690
   /* hovercard, feedDynamicHoverCard */
   FB_FEED_DYNAMIC_HOVER_CARD ThreadSourceType = 1572868
   /* search, messengerUniversalSearch */
   MESSENGER_UNIVERSAL_SEARCH ThreadSourceType = 131072
   /* story, storyAggregatedUsers, storySeenByList */
   FB_STORY ThreadSourceType = 1310720
   /* pageAboutCard */
   FB_PAGE_ABOUT_CARD ThreadSourceType = 786433
   /* inboxInThread */
   MESSENGER_INBOX_IN_THREAD ThreadSourceType = 65537
   /* notificationInThreadReply */
   MESSENGER_NOTIFICATION_IN_THREAD_REPLY ThreadSourceType = 524289
   /* archieve */
   MESSENGER_ARCHIVED_THREADS ThreadSourceType = 2031616
   /* storyViewerSheetRow */
   FB_STORY_VIEWER_SHEET_ROW ThreadSourceType = 1310722
   /* chatheadsNewMessage */
   CHATHEADS_NEW_MESSAGE ThreadSourceType = 2162691
   /* event */
   FB_EVENT ThreadSourceType = 1703936
   /* jewelSearch */
   JEWEL_SEARCH ThreadSourceType = 2097154
   /* inboxSpam, inboxThreadList, inboxRestricted */
   MESSENGER_INBOX ThreadSourceType = 65536
   /* jewelNewMessage */
   JEWEL_NEW_MESSAGE ThreadSourceType = 2097155
   /* pendingRequests */
   MESSENGER_PENDING_REQUESTS_INBOX_THREAD_LIST ThreadSourceType = 327681
   /* sidebarSearch */
   SIDEBAR_CONTACTS_SEARCH ThreadSourceType = 2228226
   /* inboxRemainingThreads */
   MESSENGER_INBOX_REMAINING_THREADS ThreadSourceType = 65541
   /* pagesHeader */
   FB_PAGE_PROFILE_HEADER_MESSAGE_BUTTON ThreadSourceType = 786434
   /* inboxRecentThreads */
   MESSENGER_INBOX_RECENT_THREADS ThreadSourceType = 65540
   /* chatheads */
   CHATHEADS ThreadSourceType = 2162688
   /* pageResponsivenessCard */
   FB_PAGE_RESPONSIVENESS_CONTEXT_CARD ThreadSourceType = 786437
   /* inboxSearch */
   MESSENGER_INBOX_MESSAGE_SEARCH ThreadSourceType = 65542
   /* jewelNestedFolder */
   JEWEL_NESTED_FOLDER ThreadSourceType = 2097156
   /* marketplace */
   MARKETPLACE_SEND_MESSAGE ThreadSourceType = 1245186
   /* feedOrganicPostViewAndMessage */
   FB_FEED_ORGANIC_POST_VIEW_AND_MESSAGE ThreadSourceType = 1572866
   /* adsCta */
   CLICK_TO_MESSENGER_AD_SEND_MESSAGE_CTA ThreadSourceType = 589826
   /* chatInThread */
   MESSENGER_CHAT_IN_THREAD ThreadSourceType = 1966082
   /* payments */
   PAYMENTS ThreadSourceType = 655360
   /* inboxFolder */
   MESSENGER_INBOX_NESTED_FOLDER ThreadSourceType = 65539
   /* inboxArchived */
   MESSENGER_INBOX_ARCHIVED_THREADS ThreadSourceType = 65545
   /* inboxActiveContacts */
   MESSENGER_INBOX_ACTIVE_CONTACTS ThreadSourceType = 65547
   /* sidebarContactsList */
   SIDEBAR_CONTACTS_LIST ThreadSourceType = 2228225
   /* sidebarCommunityChatsList */
   SIDEBAR_CONTACTS_COMMUNITY_CHATS ThreadSourceType = 2228228
   /* None */
   UNKNOWN_THREAD_SOURCE_TYPE ThreadSourceType = 0
)

type SendType int64
const (
   UNKNOWN_SEND_TYPE SendType = 0
   DEFAULT SendType = 1
)

type InitiatingSource int64
const (
   FACEBOOK_CHAT InitiatingSource = 0
   FACEBOOK_INBOX InitiatingSource = 1
   ROOMS_SIDE_CHAT InitiatingSource = 2
   FACEBOOK_FULLSCREEN InitiatingSource = 3
)

type MessageUnsendabilityStatus int64
const (
   CAN_UNSEND MessageUnsendabilityStatus =  0
   DENY_LOG_MESSAGE MessageUnsendabilityStatus =  1
   DENY_TOMBSTONE_MESSAGE MessageUnsendabilityStatus =  2
   DENY_FOR_NON_SENDER MessageUnsendabilityStatus =  3
   DENY_P2P_PAYMENT MessageUnsendabilityStatus =  4
   DENY_STORY_REACTION MessageUnsendabilityStatus =  5
   DENY_BLOB_ATTACHMENT MessageUnsendabilityStatus =  6
   DENY_MESSAGE_NOT_FOUND MessageUnsendabilityStatus =  7
   DENY_MESSAGE_INSTAGRAM_DIRECT_WRITE_RESTRICTION MessageUnsendabilityStatus =  8
)

type AppState int64
const (
   BACKGROUND AppState = 0
   FOREGROUND AppState = 1
)