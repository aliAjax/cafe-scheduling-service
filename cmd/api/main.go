package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"cafe-scheduling-api/internal/bootstrap"
	"cafe-scheduling-api/internal/businesshours"
	"cafe-scheduling-api/internal/employee"
	"cafe-scheduling-api/internal/schedule"
	"cafe-scheduling-api/internal/shift"
	"cafe-scheduling-api/internal/statistics"
	"cafe-scheduling-api/internal/user"
	"cafe-scheduling-api/migrations"
	"cafe-scheduling-api/pkg/auth"
	"cafe-scheduling-api/pkg/config"
	"cafe-scheduling-api/pkg/database"
	"cafe-scheduling-api/pkg/httputil"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	if err := database.WaitForDB(db, 2*time.Minute); err != nil {
		log.Fatalf("wait for database: %v", err)
	}
	if err := database.Migrate(db, migrations.FS); err != nil {
		log.Fatalf("run migrations: %v", err)
	}
	if err := bootstrap.Seed(context.Background(), db); err != nil {
		log.Fatalf("seed demo data: %v", err)
	}

	jwtManager := auth.NewManager(cfg.JWTSecret, cfg.JWTExpiration)

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, jwtManager)
	userHandler := user.NewHandler(userService)

	employeeRepo := employee.NewRepository(db)
	employeeService := employee.NewService(employeeRepo, userRepo)
	employeeHandler := employee.NewHandler(employeeService)

	shiftRepo := shift.NewRepository(db)
	shiftService := shift.NewService(shiftRepo)
	shiftHandler := shift.NewHandler(shiftService)

	businessHoursRepo := businesshours.NewRepository(db)
	businessHoursService := businesshours.NewService(businessHoursRepo)
	businessHoursHandler := businesshours.NewHandler(businessHoursService)

	scheduleRepo := schedule.NewRepository(db)
	scheduleService := schedule.NewService(scheduleRepo, employeeRepo, shiftRepo, businessHoursRepo)
	scheduleHandler := schedule.NewHandler(scheduleService)

	statisticsRepo := statistics.NewRepository(db)
	statisticsService := statistics.NewService(statisticsRepo)
	statisticsHandler := statistics.NewHandler(statisticsService)

	requireAuth := auth.Middleware(jwtManager)
	requireManager := auth.RequireRole("manager")
	requireEmployeeOrManager := auth.RequireRole("employee", "manager")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		httputil.JSON(w, http.StatusOK, httputil.MessageResponse{Message: "ok"})
	})

	mux.Handle("POST /api/v1/auth/login", http.HandlerFunc(userHandler.Login))
	mux.Handle("GET /api/v1/auth/me", requireAuth(http.HandlerFunc(userHandler.Me)))

	mux.Handle("GET /api/v1/employees", requireAuth(requireManager(http.HandlerFunc(employeeHandler.List))))
	mux.Handle("POST /api/v1/employees", requireAuth(requireManager(http.HandlerFunc(employeeHandler.Create))))
	mux.Handle("GET /api/v1/employees/{id}", requireAuth(requireManager(http.HandlerFunc(employeeHandler.Get))))
	mux.Handle("PUT /api/v1/employees/{id}", requireAuth(requireManager(http.HandlerFunc(employeeHandler.Update))))
	mux.Handle("DELETE /api/v1/employees/{id}", requireAuth(requireManager(http.HandlerFunc(employeeHandler.Delete))))

	mux.Handle("GET /api/v1/shift-types", requireAuth(requireManager(http.HandlerFunc(shiftHandler.List))))
	mux.Handle("POST /api/v1/shift-types", requireAuth(requireManager(http.HandlerFunc(shiftHandler.Create))))
	mux.Handle("GET /api/v1/shift-types/{id}", requireAuth(requireManager(http.HandlerFunc(shiftHandler.Get))))
	mux.Handle("PUT /api/v1/shift-types/{id}", requireAuth(requireManager(http.HandlerFunc(shiftHandler.Update))))
	mux.Handle("DELETE /api/v1/shift-types/{id}", requireAuth(requireManager(http.HandlerFunc(shiftHandler.Delete))))

	mux.Handle("GET /api/v1/business-hours", requireAuth(requireManager(http.HandlerFunc(businessHoursHandler.List))))
	mux.Handle("PUT /api/v1/business-hours", requireAuth(requireManager(http.HandlerFunc(businessHoursHandler.UpsertWeekly))))

	mux.Handle("GET /api/v1/schedules", requireAuth(requireEmployeeOrManager(http.HandlerFunc(scheduleHandler.List))))
	mux.Handle("GET /api/v1/schedules/me", requireAuth(http.HandlerFunc(scheduleHandler.Mine)))
	mux.Handle("POST /api/v1/schedules", requireAuth(requireManager(http.HandlerFunc(scheduleHandler.Create))))
	mux.Handle("PUT /api/v1/schedules/{id}", requireAuth(requireManager(http.HandlerFunc(scheduleHandler.Update))))
	mux.Handle("DELETE /api/v1/schedules/{id}", requireAuth(requireManager(http.HandlerFunc(scheduleHandler.Delete))))

	mux.Handle("GET /api/v1/statistics/weekly", requireAuth(requireManager(http.HandlerFunc(statisticsHandler.Weekly))))

	server := &http.Server{
		Addr:              ":" + cfg.ServerPort,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Printf("cafe scheduling API listening on :%s", cfg.ServerPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("serve HTTP: %v", err)
	}
}
