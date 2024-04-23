package models

type Role struct {
	Id   uint   `json:"id""`
	Name string `json:"name"`

	// Роль может иметь несколько разрешений.
	// Чтобы установить отношение  "многие ко многим",
	// необходимо создать таблицу role_permissions с помощью gorm с соответствующмим отношением (many2many):
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions"`
}
