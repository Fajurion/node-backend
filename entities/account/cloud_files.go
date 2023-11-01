package account

type CloudFile struct {
	Id        string `json:"id"`   // Format: a-[accountId]-[objectIdentifier]
	Name      string `json:"name"` // File name (encrypted with file key)
	Path      string `json:"path"` // Path to file on R2
	Type      string `json:"type"` // Mime type
	Key       string `json:"key"`  // Encryption key (encrypted with account public key)
	Account   string `json:"account"`
	Size      int64  `json:"size"` // In bytes
	Favorite  bool   `json:"favorite"`
	CreatedAt int64  `json:"-" gorm:"not null,autoCreateTime:milli"`
}
