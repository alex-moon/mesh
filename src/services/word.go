package services

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"
)

type WordService struct {
	log       *slog.Logger
	blacklist []string
	filePath  string
}

// NewWordService creates a new word service and loads the blacklist from the specified file
func NewWordService(log *slog.Logger, blacklistFilePath string) (*WordService, error) {
	service := &WordService{
		log:      log,
		filePath: blacklistFilePath,
	}

	err := service.loadBlacklist()
	if err != nil {
		return nil, err
	}

	return service, nil
}

// loadBlacklist reads the newline-delimited blacklist file
func (w *WordService) loadBlacklist() error {
	file, err := os.Open(w.filePath)
	if err != nil {
		w.log.Error("Could not open blacklist file", "path", w.filePath, "error", err)
		return fmt.Errorf("could not open blacklist file: %w", err)
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			words = append(words, strings.ToLower(word))
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	w.blacklist = words
	w.log.Info("Loaded blacklist", "word_count", len(words))
	return nil
}

// Filter processes the input string and returns the first blacklisted word found, or empty string if none
func (w *WordService) Filter(input string) string {
	if input == "" {
		return ""
	}

	// Step 1: Translate non-ASCII characters to ASCII equivalents
	normalized := w.normalizeToASCII(input)

	// Step 2: Remove all non-alphanumeric characters
	alphanumeric := w.removeNonAlphanumeric(normalized)

	// Step 3: Compare against blacklist
	return w.checkAgainstBlacklist(alphanumeric)
}

// normalizeToASCII converts non-ASCII characters to their ASCII equivalents
func (w *WordService) normalizeToASCII(input string) string {
	var result strings.Builder

	for _, r := range input {
		// Simple ASCII conversion for common accented characters
		switch r {
		case 'á', 'à', 'â', 'ä', 'ã', 'å', 'ā', 'ă', 'ą':
			result.WriteRune('a')
		case 'Á', 'À', 'Â', 'Ä', 'Ã', 'Å', 'Ā', 'Ă', 'Ą':
			result.WriteRune('A')
		case 'é', 'è', 'ê', 'ë', 'ē', 'ĕ', 'ė', 'ę', 'ě':
			result.WriteRune('e')
		case 'É', 'È', 'Ê', 'Ë', 'Ē', 'Ĕ', 'Ė', 'Ę', 'Ě':
			result.WriteRune('E')
		case 'í', 'ì', 'î', 'ï', 'ĩ', 'ī', 'ĭ', 'į':
			result.WriteRune('i')
		case 'Í', 'Ì', 'Î', 'Ï', 'Ĩ', 'Ī', 'Ĭ', 'Į':
			result.WriteRune('I')
		case 'ó', 'ò', 'ô', 'ö', 'õ', 'ō', 'ŏ', 'ő':
			result.WriteRune('o')
		case 'Ó', 'Ò', 'Ô', 'Ö', 'Õ', 'Ō', 'Ŏ', 'Ő':
			result.WriteRune('O')
		case 'ú', 'ù', 'û', 'ü', 'ũ', 'ū', 'ŭ', 'ů', 'ű', 'ų':
			result.WriteRune('u')
		case 'Ú', 'Ù', 'Û', 'Ü', 'Ũ', 'Ū', 'Ŭ', 'Ů', 'Ű', 'Ų':
			result.WriteRune('U')
		case 'ç', 'ć', 'ĉ', 'ċ', 'č':
			result.WriteRune('c')
		case 'Ç', 'Ć', 'Ĉ', 'Ċ', 'Č':
			result.WriteRune('C')
		case 'ñ', 'ń', 'ņ', 'ň':
			result.WriteRune('n')
		case 'Ñ', 'Ń', 'Ņ', 'Ň':
			result.WriteRune('N')
		case 'ý', 'ÿ', 'ŷ':
			result.WriteRune('y')
		case 'Ý', 'Ÿ', 'Ŷ':
			result.WriteRune('Y')
		default:
			// If it's ASCII or not handled, keep as is
			if r <= 127 {
				result.WriteRune(r)
			}
			// Skip non-ASCII characters that aren't handled above
		}
	}

	return result.String()
}

// removeNonAlphanumeric removes all non-alphanumeric characters and converts to lowercase
func (w *WordService) removeNonAlphanumeric(input string) string {
	// Use regex to keep only alphanumeric characters
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	cleaned := reg.ReplaceAllString(input, "")
	return strings.ToLower(cleaned)
}

// checkAgainstBlacklist checks if any blacklisted word appears as a substring
func (w *WordService) checkAgainstBlacklist(input string) string {
	for _, blacklistedWord := range w.blacklist {
		if strings.Contains(input, blacklistedWord) {
			return blacklistedWord
		}
	}
	return ""
}

// ReloadBlacklist reloads the blacklist from the file
func (w *WordService) ReloadBlacklist() error {
	return w.loadBlacklist()
}
