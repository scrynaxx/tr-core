package model

// EmployeeType определяет рабочую роль сотрудника.
type EmployeeType string

const (
	// EmployeeTypeOwner Владелец.
	EmployeeTypeOwner EmployeeType = "owner"

	// EmployeeTypeManager Менеджер.
	EmployeeTypeManager EmployeeType = "manager"

	// EmployeeTypeForeman Бригадир.
	EmployeeTypeForeman EmployeeType = "foreman"

	// EmployeeTypeLoader Грузчик.
	EmployeeTypeLoader EmployeeType = "loader"

	// EmployeeTypeAssembler Сборщик.
	EmployeeTypeAssembler EmployeeType = "assembler"
)
