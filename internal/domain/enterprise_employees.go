package domain

// EnterpriseWithEmployees — предприятие и список ID сотрудников (user_id).
type EnterpriseWithEmployees struct {
	Enterprise       Enterprise `json:"enterprise"`
	EmployeeUserIDs  []int      `json:"employee_user_ids"`
}
