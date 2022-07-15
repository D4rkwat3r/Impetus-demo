package botnet

type BaseResponse struct {
	ApiStatusCode int    `json:"api:statuscode"`
	ApiDuration   string `json:"api:duration"`
	ApiMessage    string `json:"api:message"`
	ApiTimestamp  string `json:"api:timestamp"`
}

type BaseWebResponse struct {
	ApiStatusCode int    `json:"api_status_code"`
	Code          int    `json:"code"`
	ApiMessage    string `json:"message"`
}

type UserProfile struct {
	Status                         int    `json:"status"`
	ItemsCount                     int    `json:"itemsCount"`
	Uid                            string `json:"uid"`
	ModifiedTime                   string `json:"modifiedTime"`
	FollowingStatus                int    `json:"followingStatus"`
	OnlineStatus                   int    `json:"onlineStatus"`
	AccountMembershipStatus        int    `json:"accountMembershipStatus"`
	IsGlobal                       bool   `json:"isGlobal"`
	Reputation                     int    `json:"reputation"`
	PostsCount                     int    `json:"postsCount"`
	MembersCount                   int    `json:"membersCount"`
	Nickname                       string `json:"nickname"`
	Icon                           string `json:"icon"`
	IsNicknameVerified             bool   `json:"isNicknameVerified"`
	Level                          int    `json:"level"`
	NotificationSubscriptionStatus int    `json:"notificationSubscriptionStatus"`
	PushEnabled                    bool   `json:"pushEnabled"`
	MembershipStatus               int    `json:"membershipStatus"`
	JoinedCount                    int    `json:"joinedCount"`
	Role                           int    `json:"role"`
	CommentsCount                  int    `json:"commentsCount"`
	AminoId                        string `json:"aminoId"`
	NdcId                          int    `json:"ndcId"`
	CreatedTime                    string `json:"createdTime"`
	StoriesCount                   int    `json:"storiesCount"`
	BlogsCount                     int    `json:"blogsCount"`
}

type AccountResponse struct {
	BaseResponse
	AUID        string      `json:"auid"`
	Secret      string      `json:"secret"`
	SID         string      `json:"sid"`
	UserProfile UserProfile `json:"userProfile"`
}

type Auth struct {
	Email      string `json:"email"`
	V          int    `json:"v"`
	Secret     string `json:"secret"`
	DeviceID   string `json:"deviceID"`
	ClientType int    `json:"clientType"`
	Action     string `json:"action"`
	Timestamp  int64  `json:"timestamp"`
}

type ChangeNickname struct {
	Nickname  string `json:"nickname"`
	Timestamp int64  `json:"timestamp"`
}

type ChangeAvatar struct {
	Icon      string `json:"icon"`
	Timestamp int64  `json:"timestamp"`
}

type Empty struct {
	Timestamp int64 `json:"timestamp"`
}

type Community struct {
	JoinType     int    `json:"joinType"`
	Status       int    `json:"status"`
	ModifiedTime string `json:"modifiedTime"`
	NdcId        int    `json:"ndcId"`
	Link         string `json:"link"`
	Name         string `json:"name"`
	Path         string `json:"path"`
}

type CommunityResponse struct {
	BaseResponse
	Community Community `json:"community"`
}

type LinkInfoV2 struct {
	Path       string `json:"path"`
	Extensions struct {
		LinkInfo struct {
			ObjectId   string `json:"objectId"`
			TargetCode int    `json:"targetCode"`
			NdcId      int    `json:"ndcId"`
			FullPath   string `json:"fullPath"`
			ShortCode  string `json:"shortCode"`
			ObjectType int    `json:"objectType"`
		} `json:"linkInfo"`
	} `json:"extensions"`
}

type LinkInfoResponse struct {
	BaseResponse
	LinkInfoV2 LinkInfoV2 `json:"linkInfoV2"`
}

type CommunityMembershipChange struct {
	NdcId string `json:"ndcId"`
}

type ChatMembershipChange struct {
	NdcId    string `json:"ndcId"`
	ThreadId string `json:"threadId"`
}

type MessageSend struct {
	NdcId    string `json:"ndcId"`
	ThreadId string `json:"threadId"`
	Message  struct {
		Type    int    `json:"type"`
		Content string `json:"content"`
	} `json:"message"`
}
