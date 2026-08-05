package model

// LoadType Тип погрузки / выгрузки.
type LoadType string

const (
	// LoadTypeSide Боковая погрузка / выгрузка.
	LoadTypeSide LoadType = "side"

	// LoadTypeTop Верхняя погрузка / выгрузка.
	LoadTypeTop LoadType = "top"

	// LoadTypeBack Задняя погрузка / выгрузка.
	LoadTypeBack LoadType = "back"
)
