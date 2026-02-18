package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
	// коэффициент для расчета длины шага на основе роста.
	stepLengthCoefficient = 0.45
)

func parsePackage(data string) (int, time.Duration, error) {
	// разделение строки
	parts := strings.Split(data, ",")
	// проверка длины слайса
	if len(parts) != 2 {
		return 0, 0, errors.New("некорректный формат данных")
	}
	// приобразование количества шагов в int
	steps, err := strconv.Atoi(parts[0])
	if err != nil || steps <= 0 {
		return 0, 0, err
	}
	// преобразование продолжительности в time.Duration
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, err
	}
	// возврат результатов
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// получить данные с помощью parsePackage()
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	// проверить количество шагов
	if steps <= 0 {
		return ""
	}
	// вычислить дистанцию
	distanceInMeters := float64(steps) * (height * stepLengthCoefficient)
	distanceInKm := distanceInMeters / mInKm
	// вычислить калории
	calories, _ := WalkingSpentCalories(steps, weight, height, duration)
	// сформировать и вернуть результат
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.", steps, distanceInKm, calories)
}
