package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

// Методы логгера для записи логов различного уровня
//log.Debug("This is a debug log")
//log.Info("This is an info log")
//log.Warn("This is a warning log")
//log.Error("This is an error log")

func Logger() *logrus.Logger {
	log := logrus.New()

	// Установка уровня логгирования
	log.SetLevel(logrus.DebugLevel)

	// Установка форматирования логов
	log.SetFormatter(&logrus.TextFormatter{
		//DisableColors: false,
		FullTimestamp: false,
	})

	// Настройка вывода логов в файл
	file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Error("Failed to open log file:", err)
	}

	// Создание multi-writer для вывода в файл и консоль
	//mw := io.MultiWriter(os.Stdout, file)

	// Настройка вывода логов в multi-writer
	//log.SetOutput(mw)

	return log
}
