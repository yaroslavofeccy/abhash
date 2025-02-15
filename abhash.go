package abhash

// 🔥 ABHash реализует хеш-функцию, которая:
//   - Принимает:
//     – data: исходный срез байт ([]byte)
//     – sohp: ожидаемый размер части для хеша (по умолчанию 2)
//     – hpa: количество частей хеша (по умолчанию 2)
//   - Разбивает все входные данные на hpa частей так, что первые hpa-1 частей имеют длину sohp,
//     а последняя часть получает оставшиеся байты (она может быть больше sohp).
//   - Для каждой части вычисляется токен фиксированного размера sohp с использованием XOR‑свёртки
//     и побитового инвертирования. Таким образом, при изменении данных изменяется только токен той части,
//     в которой произошли изменения.
func ABHash(data []byte, sohp int, hpa int) []byte {
	// 💡 Устанавливаем значения по умолчанию, если заданы неверно
	if sohp <= 0 {
		sohp = 2
	}
	if hpa <= 0 {
		hpa = 2
	}

	// 📋 Разбиваем data на hpa частей.
	// Первые hpa-1 частей имеют ровно sohp байт (или меньше, если данных мало, с использованием оставшихся байт)
	// Последняя часть получает все оставшиеся данные.
	parts := make([][]byte, hpa)
	for i := 0; i < hpa-1; i++ {
		start := i * sohp
		end := start + sohp

		// Если начало выходит за границы data, то ограничиваем его до длины data
		if start > len(data) {
			start = len(data)
		}

		// Если конец выходит за границы data, то ограничиваем его до длины data
		if end > len(data) { // если данных меньше, чем требуется для очередной части
			end = len(data)
		}
		parts[i] = data[start:end]
	}
	// Последняя часть получает всё оставшееся
	if (hpa-1)*sohp < len(data) {
		parts[hpa-1] = data[(hpa-1)*sohp:]
	} else {
		// Добавляем паддинг к последней части, если она меньше sohp байт
		paddingSize := sohp - len(parts[hpa-1])
		padding := make([]byte, paddingSize)
		parts[hpa-1] = append(parts[hpa-1], padding...)
	}

	// 🤗 Генерируем токены для каждой части
	tokens := make([][]byte, hpa)
	for i, part := range parts {
		tokens[i] = generateToken(part, sohp)
	}

	// 🔄 Собираем итоговый хеш, конкатенируя токены в том же порядке, что и части данных
	result := make([]byte, 0, hpa*sohp)
	for _, token := range tokens {
		result = append(result, token...)
	}

	return result
}

// generateToken принимает произвольную часть данных и сворачивает её в токен фиксированного размера.
// Используется метод XOR‑свёртки: каждый байт части добавляется (XOR) в элемент токена по индексу i % size,
// а затем производится побитовое инвертирование для дополнительного преобразования.
func generateToken(part []byte, size int) []byte {
	token := make([]byte, size)
	// ⏱️ XOR-свёртка: складываем байты части по модулю размера токена
	for i, b := range part {
		token[i%size] ^= b
	}
	// Применяем побитовое инвертирование к каждому байту токена
	for i := 0; i < size; i++ {
		token[i] = ^token[i]
	}
	return token
}

// checkSimilarity checks for bandit bits (different bits) in both hashes and returns their positions
func checkSimilarity(hash1, hash2 []byte) []int {
	var banditsbits []int

	// Find the bandit bits (different bits) in both hashes and store their positions
	for i, h1 := range hash1 {
		for j, h2 := range hash2 {
			if h1 != h2 && i == j {
				banditsbits = append(banditsbits, i)
			}
		}
	}

	return banditsbits
}
