package main

import (
	"fleetify-backend/config"
	"fleetify-backend/routes"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect database
	config.ConnectDB()

	// === Manual Migration ===
	migrateTables()

	// Inisialisasi Gin
	app := gin.Default()

	// Middleware CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Routes
	routes.RegisterRoutes(app)

	// Jalankan server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("🚀 listening on :" + port)
	if err := app.Run(":" + port); err != nil {
		log.Fatal("server error: ", err)
	}
}

func migrateTables() {
	db := config.DB

	// ===========================
	// Department
	// ===========================
	deptSQL := `
	CREATE TABLE IF NOT EXISTS departments (
		id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		department_name VARCHAR(255) NOT NULL,
		max_clock_in_time TIME NOT NULL,
		max_clock_out_time TIME NOT NULL,
		created_at DATETIME(3),
		updated_at DATETIME(3)
	) ENGINE=InnoDB;
	`
	if err := db.Exec(deptSQL).Error; err != nil {
		log.Fatal("Failed to migrate departments:", err)
	}

	// ===========================
	// Employee
	// ===========================
	empSQL := `
	CREATE TABLE IF NOT EXISTS employees (
		id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		employee_id VARCHAR(50) UNIQUE NOT NULL,
		department_id BIGINT UNSIGNED NOT NULL,
		name VARCHAR(255),
		address TEXT,
		created_at DATETIME(3),
		updated_at DATETIME(3),
		FOREIGN KEY (department_id) REFERENCES departments(id)
		ON UPDATE CASCADE
		ON DELETE CASCADE
	) ENGINE=InnoDB;
	`
	if err := db.Exec(empSQL).Error; err != nil {
		log.Fatal("Failed to migrate employees:", err)
	}

	// ===========================
	// Attendance
	// ===========================
	attSQL := `
	CREATE TABLE IF NOT EXISTS attendances (
		id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		employee_id VARCHAR(50) NOT NULL,
		attendance_id VARCHAR(100) NOT NULL UNIQUE,
		clock_in DATETIME(3),
		clock_out DATETIME(3),
		created_at DATETIME(3),
		updated_at DATETIME(3),
		FOREIGN KEY (employee_id) REFERENCES employees(employee_id)
		ON UPDATE CASCADE
		ON DELETE CASCADE
	) ENGINE=InnoDB;
	`
	if err := db.Exec(attSQL).Error; err != nil {
		log.Fatal("Failed to migrate attendances:", err)
	}

	// ===========================
	// AttendanceHistory
	// ===========================
	histSQL := `
	CREATE TABLE IF NOT EXISTS attendance_histories (
		id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
		employee_id VARCHAR(50) NOT NULL,
		attendance_id VARCHAR(100) NOT NULL,
		date_attendance DATETIME(3),
		attendance_type TINYINT,
		description TEXT,
		created_at DATETIME(3),
		updated_at DATETIME(3),
		FOREIGN KEY (employee_id) REFERENCES employees(employee_id),
		FOREIGN KEY (attendance_id) REFERENCES attendances(attendance_id)
		ON UPDATE CASCADE
		ON DELETE CASCADE
	) ENGINE=InnoDB;
	`
	if err := db.Exec(histSQL).Error; err != nil {
		log.Fatal("Failed to migrate attendance_histories:", err)
	}

	log.Println("✅ Manual migration completed")
}
