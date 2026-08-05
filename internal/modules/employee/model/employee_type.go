package model

type EmployeeType string

const (
	EmployeeTypeOwner     EmployeeType = "owner"
	EmployeeTypeManager   EmployeeType = "manager"
	EmployeeTypeForeman   EmployeeType = "foreman"
	EmployeeTypeLoader    EmployeeType = "loader"
	EmployeeTypeAssembler EmployeeType = "assembler"
)
