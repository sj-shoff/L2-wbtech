package main

import (
	"reflect"
	"sort"
	"testing"
)

// TestGroupAnagrams_BasicCase тестирует основной случай с несколькими множествами анаграмм
func TestGroupAnagrams_BasicCase(t *testing.T) {
	input := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}

	result := groupAnagrams(input)

	expected := map[string][]string{
		"пятак":  {"пятак", "пятка", "тяпка"},
		"листок": {"листок", "слиток", "столик"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Ожидалось %v, но получено %v", expected, result)
	}
}

// TestGroupAnagrams_SingleWords тестирует случай, когда нет анаграмм (слова без пар)
func TestGroupAnagrams_SingleWords(t *testing.T) {
	input := []string{"дом", "стол", "книга"}

	result := groupAnagrams(input)

	expected := map[string][]string{}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Ожидалось %v, но получено %v", expected, result)
	}
}

// TestGroupAnagrams_EmptyInput тестирует пустой входной массив
func TestGroupAnagrams_EmptyInput(t *testing.T) {
	input := []string{}

	result := groupAnagrams(input)

	expected := map[string][]string{}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Ожидалось %v, но получено %v", expected, result)
	}
}

// TestGroupAnagrams_CaseInsensitive тестирует регистронезависимость
func TestGroupAnagrams_CaseInsensitive(t *testing.T) {
	input := []string{"ПяТак", "пЯтКа", "ТяПкА", "ЛиСтОк", "СлИтОк", "СтОлИк"}

	result := groupAnagrams(input)

	expected := map[string][]string{
		"пятак":  {"пятак", "пятка", "тяпка"},
		"листок": {"листок", "слиток", "столик"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Ожидалось %v, но получено %v", expected, result)
	}
}

// TestGroupAnagrams_Duplicates тестирует обработку дубликатов
func TestGroupAnagrams_Duplicates(t *testing.T) {
	input := []string{"пятак", "пятка", "пятак", "тяпка", "пятка", "листок"}

	result := groupAnagrams(input)

	expected := map[string][]string{
		"пятак": {"пятак", "пятка", "тяпка"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Ожидалось %v, но получено %v", expected, result)
	}
}

// TestGroupAnagrams_MixedLength тестирует слова разной длины
func TestGroupAnagrams_MixedLength(t *testing.T) {
	input := []string{"кот", "ток", "окот", "окто"}

	result := groupAnagrams(input)

	expected := map[string][]string{
		"кот":  {"кот", "ток"},
		"окот": {"окот", "окто"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Ожидалось %v, но получено %v", expected, result)
	}
}

// TestGroupAnagrams_OneValidGroup тестирует только одну валидную группу
func TestGroupAnagrams_OneValidGroup(t *testing.T) {
	input := []string{"клоун", "колун", "кулон", "стол", "стул"}

	result := groupAnagrams(input)

	expected := map[string][]string{
		"клоун": {"клоун", "колун", "кулон"},
	}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Ожидалось %v, но получено %v", expected, result)
	}
}

// TestGroupAnagrams_DifferentLengths тестирует, что слова разной длины не считаются анаграммами
func TestGroupAnagrams_DifferentLengths(t *testing.T) {
	input := []string{"кот", "котик", "собака"}

	result := groupAnagrams(input)

	expected := map[string][]string{}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Ожидалось %v, но получено %v", expected, result)
	}
}

// TestGroupAnagrams_ActualAnagrams тестирует реальные анаграммы одинаковой длины
func TestGroupAnagrams_ActualAnagrams(t *testing.T) {
	input := []string{"котик", "тикок", "кокит"}

	result := groupAnagrams(input)

	if len(result) != 1 {
		t.Errorf("Ожидалась 1 группа анаграмм, но получено %d", len(result))
	}

	for _, group := range result {
		expectedWords := []string{"котик", "тикок", "кокит"}
		sort.Strings(expectedWords)
		sort.Strings(group)

		if !reflect.DeepEqual(group, expectedWords) {
			t.Errorf("Ожидалось %v, но получено %v", expectedWords, group)
		}
	}
}

// Вспомогательная функция для проверки отдельных свойств
func TestGroupAnagrams_Properties(t *testing.T) {
	input := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик"}

	result := groupAnagrams(input)

	if len(result) == 0 {
		t.Error("Результат не должен быть пустым")
	}

	for key, group := range result {
		for _, char := range key {
			if char < 'а' || char > 'я' {
				t.Errorf("Ключ %s содержит не-буквенные символы", key)
			}
		}

		if len(group) < 2 {
			t.Errorf("Группа для ключа %s содержит меньше 2 элементов: %v", key, group)
		}

		for _, word := range group {
			for _, char := range word {
				if char < 'а' || char > 'я' {
					t.Errorf("Слово %s в группе %s содержит не-буквенные символы", word, key)
				}
			}
		}

		if !sort.StringsAreSorted(group) {
			t.Errorf("Группа для ключа %s не отсортирована: %v", key, group)
		}
	}
}
