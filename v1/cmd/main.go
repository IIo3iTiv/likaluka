package miniogo

import (
	"log"
	"miniogo/internal/common/config"
	"miniogo/internal/handler"
	"miniogo/pkg/minio"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загрузка конфигурации из файла .env
	config.LoadConfig()

	// Инициализация соединения с Minio
	minioClient := minio.NewMinioClient()
	err := minioClient.InitMinio()
	if err != nil {
		log.Fatalf("Ошибка инициализации Minio: %v", err)
	}

	_, s := handler.NewHandler(
		minioClient,
	)

	// Инициализация маршрутизатора Gin
	router := gin.Default()

	s.RegisterRoutes(router)

	// Запуск сервера Gin
	port := config.AppConfig.Port // Мы берем порт из конфига
	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("Ошибка запуска сервера Gin: %v", err)
	}
}
