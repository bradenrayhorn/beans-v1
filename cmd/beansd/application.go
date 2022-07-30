package main

import (
	"fmt"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/http"
	"github.com/bradenrayhorn/beans/inmem"
	"github.com/bradenrayhorn/beans/logic"
	"github.com/bradenrayhorn/beans/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Application struct {
	httpServer *http.Server
	pool       *pgxpool.Pool

	config Config

	budgetRepository  beans.BudgetRepository
	budgetService     beans.BudgetService
	userRepository    beans.UserRepository
	userService       beans.UserService
	sessionRepository beans.SessionRepository
}

func NewApplication(c Config) *Application {
	return &Application{
		config: c,
	}
}

func (a *Application) Start() error {
	pool, err := postgres.CreatePool(
		fmt.Sprintf("postgres://%s:%s@%s/%s",
			a.config.Postgres.Username,
			a.config.Postgres.Password,
			a.config.Postgres.Addr,
			a.config.Postgres.Database,
		))

	if err != nil {
		panic(err)
	}
	a.pool = pool

	a.budgetRepository = postgres.NewBudgetRepository(pool)
	a.budgetService = logic.NewBudgetService(a.budgetRepository)
	a.userRepository = postgres.NewUserRepository(pool)
	a.userService = &logic.UserService{UserRepository: a.userRepository}
	a.sessionRepository = inmem.NewSessionRepository()

	a.httpServer = http.NewServer(
		a.budgetRepository,
		a.budgetService,
		a.userRepository,
		a.userService,
		a.sessionRepository,
	)
	if err := a.httpServer.Open(":" + a.config.Port); err != nil {
		panic(err)
	}

	return nil
}

func (a *Application) Stop() error {
	if err := a.httpServer.Close(); err != nil {
		return err
	}

	a.pool.Close()

	return nil
}

func (a *Application) HttpServer() *http.Server {
	return a.httpServer
}

func (a *Application) BudgetRepository() beans.BudgetRepository {
	return a.budgetRepository
}

func (a *Application) UserRepository() beans.UserRepository {
	return a.userRepository
}

func (a *Application) SessionRepository() beans.SessionRepository {
	return a.sessionRepository
}
