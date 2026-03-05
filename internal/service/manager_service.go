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

// ApproveUser подтверждает регистрацию клиента (is_active = true). Менеджер может подтверждать только пользователей с ролью client.
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

// GetUserHistory возвращает объединённую историю операций по всем счетам пользователя (то же, что /client/accounts/history по каждому счёту, объединённое и отсортированное по дате).
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

// BlockAccount блокирует счёт (менеджер). Операции по счёту будут запрещены.
func (s *Manager) BlockAccount(accountID int) error {
	_, err := s.accountRepo.GetAccountByID(accountID)
	if err != nil {
		return err
	}
	return s.accountRepo.SetAccountBlocked(accountID, true)
}

// UnblockAccount разблокирует счёт (менеджер).
func (s *Manager) UnblockAccount(accountID int) error {
	_, err := s.accountRepo.GetAccountByID(accountID)
	if err != nil {
		return err
	}
	return s.accountRepo.SetAccountBlocked(accountID, false)
}

// BlockDeposit блокирует вклад (менеджер). Может заблокировать в любой момент, без проверки на баланс.
func (s *Manager) BlockDeposit(depositID int) error {
	_, err := s.depositRepo.GetDepositByID(depositID)
	if err != nil {
		return err
	}
	return s.depositRepo.SetDepositBlocked(depositID, true)
}

// UnblockDeposit разблокирует вклад (менеджер).
func (s *Manager) UnblockDeposit(depositID int) error {
	_, err := s.depositRepo.GetDepositByID(depositID)
	if err != nil {
		return err
	}
	return s.depositRepo.SetDepositBlocked(depositID, false)
}

// GetEnterprisesWithEmployees возвращает список предприятий с привязкой ID сотрудников по каждому.
func (s *Manager) GetEnterprisesWithEmployees() ([]domain.EnterpriseWithEmployees, error) {
	return s.enterpriseRepo.GetEnterprisesWithEmployees()
}

// AddEmployee добавляет пользователя в предприятие как сотрудника.
func (s *Manager) AddEmployee(enterpriseID, userID int) error {
	_, err := s.enterpriseRepo.GetEnterpriseByID(enterpriseID)
	if err != nil {
		return err
	}
	return s.enterpriseRepo.AddEmployee(enterpriseID, userID)
}

// RemoveEmployee удаляет сотрудника из предприятия; pending заявки по этому предприятию отклоняются.
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

// ApproveSalaryApplication одобряет заявку на ЗП. Проверяет баланс предприятия >= сумма заявки.
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
