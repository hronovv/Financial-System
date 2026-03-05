package service

import (
	"sort"

	"financial_system/internal/domain"
	"financial_system/internal/repository"
)

type Manager struct {
	userRepo    repository.UserRepository
	accountRepo repository.AccountRepository
}

func NewManagerService(userRepo repository.UserRepository, accountRepo repository.AccountRepository) *Manager {
	return &Manager{userRepo: userRepo, accountRepo: accountRepo}
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
