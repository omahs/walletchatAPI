package entity

type Chatitem struct {
	Fromaddr  string `json:"fromaddr"`
	Toaddr    string `json:"toaddr"`
	Timestamp string `json:"timestamp"`
	Msgread   bool   `json:"read"`
	Mmkeyused bool   `json:"mmkeyused"`
	Message   string `json:"message"`
	Nftaddr   string `json:"nftaddr"`
	Nftid     int    `json:"nftid"`
}

type Chatiteminbox struct {
	Fromaddr  string `json:"fromaddr"`
	Toaddr    string `json:"toaddr"`
	Timestamp string `json:"timestamp"`
	Msgread   bool   `json:"read"`
	Message   string `json:"message"`
	Unreadcnt int    `json:"unread"`
}

// type ChatitemRsp struct {
// 	ID        int    `json:"id"`
// 	Fromaddr  string `json:"fromaddr"`
// 	Toaddr    string `json:"toaddr"`
// 	Timestamp string `json:"timestamp"`
// 	Read    string `json:"read"`
// 	Message   string `json:"message"`
// }
