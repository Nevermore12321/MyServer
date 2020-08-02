package adminDB

type AdminDBOperation interface {
	Insert() error
	Delete(hardDel bool) error
	Update(modify map[string]interface{}, queryString interface{}, keyList ...interface{}) error
}
