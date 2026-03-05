package service

import (
	"sort"

	"financial_system/internal/domain"
	"financial_system/internal/repository"
)

type Manager struct {
	userRepo       repository.UserRepository
	accountRepo    repository.AccountRepository
	depositRepo    repository.DepositRepository
	enterpriseRepo repository.EnterpriseRepository
	salaryRepo     repository.SalaryApplicationRepository
}

func NewManagerService(userRepo repository.UserRepository, accountRepo repository.AccountRepository, depositRepo repository.DepositRepository, enterpriseRepo repository.EnterpriseRepository, salaryRepo repository.SalaryApplicationRepository) *Manager {
	return &Manager{
		userRepo:       userRepo,
		accountRepo:    accountRepo,
		depositRepo:    depositRepo,
		enterpriseRepo: enterpriseRepo,
		salaryRepo:     salaryRepo,
	}
}

// ApproveUser устанавливает is_active для пользователя с ролью client.
func (s *Manager) ApproveUser(userID int) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user.Role != domain.RoleClient {
		return domain.ErrCanOnlyApproveClient
	}
	if user.IsActive {
		return domain.ErrUserAlreadyActive
	}
	return s.userRepo.SetUserActive(userID, true)
}

// GetUserHistory возвращает объединённую историю операций по всем счетам пользователя.
func (s *Manager) GetUserHistory(userID int) ([]domain.Transaction, error) {
	accounts, err := s.accountRepo.GetAccountsByUserID(userID)
	if err != nil {
		return nil, err
	}

	seen := make(map[int]struct{})
	var all []domain.Transaction

	for _, acc := range accounts {
		history, err := s.accountRepo.GetAccountHistory(acc.ID)
		if err != nil {
			return nil, err
		}
		for _, t := range history {
			if _, ok := seen[t.ID]; ok {
				continue
			}
			seen[t.ID] = struct{}{}
			all = append(all, t)
		}
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})
	return all, nil
}

// BlockAccount блокирует счёт.
func (s *Manager) BlockAccount(accountID int) error {
	_, err := s.accountRepo.GetAccountByID(accountID)
	if err != nil {
		return err
	}
	return s.accountRepo.SetAccountBlocked(accountID, true)
}

// UnblockAccount снимает блокировку счёта.
func (s *Manager) UnblockAccount(accountID int) error {
	_, err := s.accountRepo.GetAccountByID(accountID)
	if err != nil {
		return err
	}
	return s.accountRepo.SetAccountBlocked(accountID, false)
}

// BlockDeposit блокирует вклад.
func (s *Manager) BlockDeposit(depositID int) error {
	_, err := s.depositRepo.GetDepositByID(depositID)
	if err != nil {
		return err
	}
	return s.depositRepo.SetDepositBlocked(depositID, true)
}

// UnblockDeposit снимает блокировку вклада.
func (s *Manager) UnblockDeposit(depositID int) error {
	_, err := s.depositRepo.GetDepositByID(depositID)
	if err != nil {
		return err
	}
	return s.depositRepo.SetDepositBlocked(depositID, false)
}

// GetEnterprisesWithEmployees возвращает предприятия со списками ID сотрудников.
func (s *Manager) GetEnterprisesWithEmployees() ([]domain.EnterpriseWithEmployees, error) {
	return s.enterpriseRepo.GetEnterprisesWithEmployees()
}

// AddEmployee привязывает пользователя к предприятию как сотрудника.
func (s *Manager) AddEmployee(enterpriseID, userID int) error {
	_, err := s.enterpriseRepo.GetEnterpriseByID(enterpriseID)
	if err != nil {
		return err
	}
	return s.enterpriseRepo.AddEmployee(enterpriseID, userID)
}

// RemoveEmployee отвязывает сотрудника от предприятия; его pending-заявки отклоняются.
func (s *Manager) RemoveEmployee(enterpriseID, userID int) error {
	_, err := s.enterpriseRepo.GetEnterpriseByID(enterpriseID)
	if err != nil {
		return err
	}
	if err := s.salaryRepo.RejectPendingByUserAndEnterprise(userID, enterpriseID); err != nil {
		return err
	}
	return s.enterpriseRepo.RemoveEmployee(enterpriseID, userID)
}

// ApproveSalaryApplication одобряет заявку; при недостатке баланса предприятия возвращает ошибку.
func (s *Manager) ApproveSalaryApplication(applicationID int) error {
	app, err := s.salaryRepo.GetByID(applicationID)
	if err != nil {
		return err
	}
	if app.Status != domain.SalaryApplicationStatusPending {
		return domain.ErrApplicationNotPending
	}
	ent, err := s.enterpriseRepo.GetEnterpriseByID(app.EnterpriseID)
	if err != nil {
		return err
	}
	if ent.Balance < app.Amount {
		return domain.ErrInsufficientEnterpriseBalance
	}
	return s.salaryRepo.UpdateStatus(applicationID, domain.SalaryApplicationStatusApproved)
}
