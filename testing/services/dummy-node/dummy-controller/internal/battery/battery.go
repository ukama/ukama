package battery

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var provider BatteryProvider

type BatteryProvider interface {
    GetMetrics() (*BatteryMetrics, error)
}

type RealBatteryProvider struct {
    batteryPath string
}

type BatteryMetrics struct {
    Capacity    float64
    Voltage     float64
    Current     float64
    Temperature float64
    Status      string
    Health      string
}

func init() {
    if os.Getenv("MOCK_BATTERY") == "true" {
        provider = NewMockBatteryProvider()
    } else {
        provider = NewRealBatteryProvider()
    }
}

func GetBatteryMetrics() (*BatteryMetrics, error) {
    return provider.GetMetrics()
}

func NewRealBatteryProvider() *RealBatteryProvider {
    // Try to find the battery path
    batteryPath := findBatteryPath()
    return &RealBatteryProvider{
        batteryPath: batteryPath,
    }
}

func findBatteryPath() string {
    // Common battery paths to check
    basePaths := []string{
        "/sys/class/power_supply",
    }

    for _, basePath := range basePaths {
        entries, err := os.ReadDir(basePath)
        if (err != nil) {
            continue
        }

        // Look for any battery-like directory
        for _, entry := range entries {
            if !entry.IsDir() {
                continue
            }
            
            // Check if this is a battery directory
            fullPath := filepath.Join(basePath, entry.Name())
            typeFile := filepath.Join(fullPath, "type")
            typeContent, err := os.ReadFile(typeFile)
            if err == nil && strings.TrimSpace(string(typeContent)) == "Battery" {
                return fullPath
            }

            // Check for common battery names
            name := strings.ToLower(entry.Name())
            if strings.Contains(name, "bat") || strings.Contains(name, "battery") {
                return fullPath
            }
        }
    }
    
    return ""
}

func (r *RealBatteryProvider) GetMetrics() (*BatteryMetrics, error) {
    if r.batteryPath == "" {
        // If no battery found, try to find it again
        r.batteryPath = findBatteryPath()
        if r.batteryPath == "" {
            return nil, fmt.Errorf("no battery found in system")
        }
    }

    metrics := &BatteryMetrics{}
    var errs []string

    // Read capacity
    if capacity, err := r.readSysfs("capacity"); err == nil {
        if val, err := strconv.ParseFloat(capacity, 64); err == nil {
            metrics.Capacity = val
        } else {
            errs = append(errs, fmt.Sprintf("parse capacity error: %v", err))
        }
    } else {
        errs = append(errs, fmt.Sprintf("read capacity error: %v", err))
    }

    // Read voltage
    if voltage, err := r.readSysfs("voltage_now"); err == nil {
        if val, err := strconv.ParseFloat(voltage, 64); err == nil {
            metrics.Voltage = val / 1000000 // Convert to volts
        } else {
            errs = append(errs, fmt.Sprintf("parse voltage error: %v", err))
        }
    } else {
        errs = append(errs, fmt.Sprintf("read voltage error: %v", err))
    }

    // Read current
    if current, err := r.readSysfs("current_now"); err == nil {
        if val, err := strconv.ParseFloat(current, 64); err == nil {
            metrics.Current = val / 1000000 // Convert to amperes
        } else {
            errs = append(errs, fmt.Sprintf("parse current error: %v", err))
        }
    } else {
        errs = append(errs, fmt.Sprintf("read current error: %v", err))
    }

    // Read temperature
    if temp, err := r.readSysfs("temp"); err == nil {
        if val, err := strconv.ParseFloat(temp, 64); err == nil {
            metrics.Temperature = val / 10 // Convert to celsius
        } else {
            errs = append(errs, fmt.Sprintf("parse temperature error: %v", err))
        }
    } else {
        errs = append(errs, fmt.Sprintf("read temperature error: %v", err))
    }

    // Read status
    if status, err := r.readSysfs("status"); err == nil {
        metrics.Status = status
    } else {
        errs = append(errs, fmt.Sprintf("read status error: %v", err))
    }

    // Read health
    if health, err := r.readSysfs("health"); err == nil {
        metrics.Health = health
    } else {
        errs = append(errs, fmt.Sprintf("read health error: %v", err))
    }

    // Return error if no metrics were read successfully
    if len(errs) == 6 { // All reads failed
        return nil, fmt.Errorf("failed to read any battery metrics: %s", strings.Join(errs, "; "))
    }

    return metrics, nil
}

func (r *RealBatteryProvider) readSysfs(file string) (string, error) {
    path := filepath.Join(r.batteryPath, file)
    content, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return "", fmt.Errorf("battery file not found: %s", path)
        }
        if os.IsPermission(err) {
            return "", fmt.Errorf("permission denied accessing battery file: %s", path)
        }
        return "", fmt.Errorf("error reading battery file %s: %v", path, err)
    }
    return strings.TrimSpace(string(content)), nil
}