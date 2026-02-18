package spentcalories

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// разделяем строку на части
	parts := strings.Split(data, ",")
	// проверяем что у нас ровно три элемента
	if len(parts) != 3 {
		return 0, "", 0, errors.New("некорректный формат данных")
	}
	// преобразуем количество шагов в int с возможными ошибками
	steps, err := strconv.Atoi(parts[0])
	if err != nil || steps <= 0 {
		return 0, "", 0, err
	}
	// получаем вид активности
	activity := parts[1]
	// парсим продолжительность в time.Duration
	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return 0, "", 0, err
	}
	// если всё прошло успешно возвращаем результат
	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	// расчитываем длину шага пользователя
	stepLength := height * stepLengthCoefficient
	// переводим количество шагов в метры
	totalDistanceMeters := float64(steps) * stepLength
	// переводим дистанцию из метров в км
	distanceKm := totalDistanceMeters / mInKm
	// возвращаем результат
	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// проверяем что продолжительность больше 0
	if duration <= 0 {
		return 0
	}
	// расчитывает дистанцию в км
	distanceKm := distance(steps, height)
	// переводим продолжительность в часы
	durationHours := duration.Hours()
	// вычисляем среднюю скорость
	speed := distanceKm / durationHours
	// возвращаем результат
	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// получаем данные из строки с помощью parseTraining
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	// рассчитываем дистанцию, среднюю скорость, и калории в зависимости от вида активности
	var distanceKm float64
	var calories float64
	var speedKmPerH float64

	switch activity {
	case "Ходьба":
		distanceKm = distance(steps, height)
		speedKmPerH = meanSpeed(steps, height, duration)
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	case "Бег":
		distanceKm = distance(steps, height)
		speedKmPerH = meanSpeed(steps, height, duration)
		calories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}
	// формируем строку с информацией о тренировке
	result := fmt.Srintf(`Тип тренировки: %s
	Длительность: %.2f ч.
	Дистанция: %.2f км.
	Скорость: %.2f км/ч.
	Сожгли калорий: %.2f`, activity, duration.Hours(), distanceKm, speedKmPerH, calories)
	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// проверяем входные параметры на корректность
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные параметры")
	}
	// расчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)
	// переводим продолжительность в минуты
	durationInMinutes := duration.Minutes()
	// расчитываем количество калорий
	calories := (weight * speed * durationInMinutes) / minInH
	// возвращаем результат
	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// проверяем входные параметры на корректность
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("некорректные входные параметры")
	}
	// расчитываем среднюю скорость
	speed := meanSpeed(steps, height, duration)
	// переводим продолжительность в минуты
	durationInMinutes := duration.Minutes()
	// рассчитываем количество калорий
	calories := (weight * speed * durationInMinutes) / minInH
	// умножаем на корректирующий коэффицент для ходьбы
	calories *= walkingCaloriesCoefficient
	// возвращаем результат
	return calories, nil
}
