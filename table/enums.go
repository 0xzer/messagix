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