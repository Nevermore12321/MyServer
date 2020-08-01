package adminDB

type AdminDBOperation interface {
	Insert([]*UserInfo) error
	Delete(users []*UserInfo) error
	Update(users []*UserInfo) error
	Query(users []*UserInfo) error
}
