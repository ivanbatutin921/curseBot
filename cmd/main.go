package main

import (
	"cursebot/internal/bot"
	"cursebot/internal/config"
	"cursebot/internal/models"
	"cursebot/internal/repository"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Определяем путь к .env файлу относительно текущей директории
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting working directory: %v", err)
	}
	envPath := filepath.Join(workDir, ".env")

	// Загрузка переменных окружения из .env файла
	if err := godotenv.Load(envPath); err != nil {
		log.Fatalf("Error loading .env file from %s: %v", envPath, err)
	}

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Подключение к базе данных
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{
		PrepareStmt: false,                                 // Отключаем подготовленные операторы
		Logger:      logger.Default.LogMode(logger.Silent), // Отключаем логирование SQL-запросов
	})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Получаем SQL-соединение для очистки кэша
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error getting sql.DB: %v", err)
	}

	// Очищаем кэш подготовленных запросов
	if _, err := sqlDB.Exec("DEALLOCATE ALL"); err != nil {
		log.Printf("Warning: could not deallocate prepared statements: %v", err)
	}

	// Миграция базы данных
	err = db.AutoMigrate(&models.User{}, &models.Course{}, &models.Access{})
	if err != nil {
		log.Printf("Warning: database migration error (tables might already exist): %v", err)
	}

	// Инициализация репозитория с отключенным кэшем
	repo := repository.NewRepository(db)

	// Создание тестового курса, если его нет
	var count int64
	if err := db.Model(&models.Course{}).Count(&count).Error; err != nil {
		log.Printf("Warning: could not count courses: %v", err)
	} else if count == 0 {
		testCourse := &models.Course{
			Title:       "Курс по программированию",
			Description: "Научитесь программировать с нуля за 3 месяца",
			Price:       15000,
		}
		if err := db.Create(testCourse).Error; err != nil {
			log.Printf("Warning: could not create test course: %v", err)
		}
	}

	// Инициализация бота
	b, err := gotgbot.NewBot(cfg.BotToken, nil)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}

	// Создание и запуск бота
	telegramBot := bot.NewBot(b, repo, cfg.ChannelID)
	if err := telegramBot.Start(); err != nil {
		log.Fatalf("Error starting bot: %v", err)
	}

	// Ожидание сигнала завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

// isTableExistsError проверяет, является ли ошибка ошибкой существования таблицы
func isTableExistsError(err error) bool {
	return err != nil && err.Error() == "ERROR: relation \"users\" already exists (SQLSTATE 42P07)"
}
