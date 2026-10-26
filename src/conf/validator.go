package conf

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"p2r/src/logger"
)

// Validator handles input validation and path checking.
type Validator struct {
	Lg *logger.Logger
}

// NewValidator creates a new Validator instance.
func NewValidator(lg *logger.Logger) *Validator {
	return &Validator{Lg: lg}
}

// ValidatePath validates and parses a path string (local or remote).
func (v *Validator) ValidatePath(pathStr string) (Path, error) {
	path := Path{
		FullPath: pathStr,
		Path:     pathStr,
	}

	// Check for obviously invalid patterns
	if v.isPathTooShort(pathStr) {
		return path, errors.New("path seems too short: " + pathStr)
	}

	// Detect if local or remote path
	path.IsLocal = !strings.Contains(pathStr, ":")

	if path.IsLocal {
		// Validate local path
		path.Exists, _ = v.pathExists(path.FullPath)
		path.IsFolder = v.isDirectory(path.FullPath)
		path.IsEmpty = v.isDirectoryEmpty(path.FullPath)
	} else {
		// Parse remote path (machine:path format)
		parts := strings.SplitN(pathStr, ":", 2)
		if len(parts) != 2 {
			return path, errors.New("invalid remote path format: " + pathStr)
		}
		path.Machine = parts[0]
		path.Path = parts[1]

		// Basic validation for remote host
		if err := v.validateHostname(path.Machine); err != nil {
			path.IsHealthy = false
			v.Lg.Warn("invalid hostname in remote path", map[string]interface{}{"machine": path.Machine, "error": err})
		}
	}

	// Check for dangerous patterns
	if err := v.checkDangerousPatterns(pathStr); err != nil {
		path.IsHealthy = false
		return path, err
	}

	path.IsHealthy = true
	return path, nil
}

// Path represents a validated path (local or remote).
type Path struct {
	IsLocal   bool
	Machine   string
	Path      string
	FullPath  string
	Exists    bool
	IsFolder  bool
	IsEmpty   bool
	IsHealthy bool
	Errors    []error
}

// validatePathList validates a list of paths.
func (v *Validator) validatePathList(paths []string) []Path {
	validated := make([]Path, 0, len(paths))
	for _, p := range paths {
		validatedPath, err := v.ValidatePath(p)
		if err != nil {
			validatedPath.Errors = append(validatedPath.Errors, err)
		}
		validated = append(validated, validatedPath)
	}
	return validated
}

// isPathTooShort checks if path is suspiciously short.
func (v *Validator) isPathTooShort(path string) bool {
	// Match paths that are just special characters
	matched, _ := regexp.MatchString(`^[:/-_]{1,2}$`, path)
	return matched
}

// pathExists checks if a local path exists.
func (v *Validator) pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// isDirectory checks if a path is a directory.
func (v *Validator) isDirectory(fn string) bool {
	fi, err := os.Stat(fn)
	if err != nil {
		v.Lg.Debug("path stat failed", map[string]interface{}{"path": fn, "error": err})
		return false
	}
	return fi.IsDir()
}

// isDirectoryEmpty checks if a directory is empty.
func (v *Validator) isDirectoryEmpty(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == io.EOF {
		return true
	}
	return false
}

// validateHostname performs basic hostname validation.
func (v *Validator) validateHostname(host string) error {
	if host == "" {
		return errors.New("empty hostname")
	}

	// Basic check: no spaces, no quotes, alphanumeric + dots + hyphens
	for _, ch := range host {
		if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '.' && ch != '-' {
			return fmt.Errorf("invalid character in hostname: %q", ch)
		}
	}

	return nil
}

// checkDangerousPatterns checks for injection attempts in paths.
func (v *Validator) checkDangerousPatterns(path string) error {
	dangerousPatterns := []string{
		";",  // Command separator
		"|",  // Pipe
		"&&", // Command chaining
		"||", // Or operator
		"``", // Backtick execution
		"$(", // Command substitution
		"$_", // Variable expansion
		"{",  // Shell pattern
		"}",  // Shell pattern
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(path, pattern) {
			return fmt.Errorf("dangerous character detected in path: %q contains %q", path, pattern)
		}
	}

	return nil
}

// ValidateSyncStep validates a sync step configuration.
func (v *Validator) ValidateSyncStep(step SyncStep) (validated Step, errArr []error) {
	validated.Cmd = step.Cmd

	// Validate command is not empty
	if len(step.Cmd) == 0 {
		errArr = append(errArr, errors.New("command array is empty"))
	}

	// Validate sources using validatePathList
	validatedSourcePaths := v.validatePathList(step.Sources)
	for i, src := range step.Sources {
		if err := v.checkDangerousPatterns(src); err != nil {
			errArr = append(errArr, fmt.Errorf("source path[%d] invalid: %w", i, err))
		}
		// Check for validation errors from validatePathList
		if len(validatedSourcePaths[i].Errors) > 0 {
			for _, err := range validatedSourcePaths[i].Errors {
				errArr = append(errArr, fmt.Errorf("source path[%d] validation error: %w", i, err))
			}
		}
		validated.SourcePaths = append(validated.SourcePaths, src)
	}

	// Validate targets using validatePathList
	validatedTargetPaths := v.validatePathList(step.Targets)
	for i, tgt := range step.Targets {
		if err := v.checkDangerousPatterns(tgt); err != nil {
			errArr = append(errArr, fmt.Errorf("target path[%d] invalid: %w", i, err))
		}
		// Check for validation errors from validatePathList
		if len(validatedTargetPaths[i].Errors) > 0 {
			for _, err := range validatedTargetPaths[i].Errors {
				errArr = append(errArr, fmt.Errorf("target path[%d] validation error: %w", i, err))
			}
		}
		validated.TargetPaths = append(validated.TargetPaths, tgt)
	}

	return validated, errArr
}

// Step represents a validated sync step.
type Step struct {
	Cmd         []string
	SourcePaths []string
	TargetPaths []string
	Errors      []error
}

// ValidateCommand validates a command array.
func (v *Validator) ValidateCommand(cmd []string) ([]string, error) {
	if len(cmd) == 0 {
		return nil, errors.New("command array is empty")
	}

	// Validate command arguments don't contain shell metacharacters
	for i, arg := range cmd {
		if err := v.checkDangerousPatterns(arg); err != nil {
			return cmd, fmt.Errorf("arg[%d]: %w", i, err)
		}
	}

	// Resolve command path
	absCmd, err := filepath.Abs(cmd[0])
	if err == nil {
		cmd[0] = absCmd
	}

	return cmd, nil
}
