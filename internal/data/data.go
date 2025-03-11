package data

import "SearchServices/internal/models"

func AllData() []models.Services {
	AllData := []models.Services{
		{
			Name: "Сброс пароля",
			Requests: []string{
				"Сброс пароля Active Directory",
				"Не могу войти в учетную запись",
				"Восстановление доступа к почте",
			},
		},
		{
			Name: "Установка ПО",
			Requests: []string{
				"Установка программного обеспечения",
				"Не запускается приложение",
			},
		},
	}
	return AllData
}
