package main

import (
	"FileCrypt/crypt"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed frontend/dist
var assets embed.FS

// EncryptFile method that calls the function from aes.go
func (a *App) EncryptFile(filePath string, key string) (string, error) {
	return crypt.EncryptFile(filePath, key, a.updateProgress)
}

// DecryptFile method that calls the function from aes.go
func (a *App) DecryptFile(filePath string, key string) (string, error) {
	return crypt.DecryptFile(filePath, key, a.updateProgress)
}

// Update progress method (thread-safe)
func (a *App) updateProgress(percent int, status string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.progress = percent
	a.status = status
	runtime.EventsEmit(a.ctx, "progress", map[string]interface{}{
		"percentage": percent,
		"status":     status,
	})
}
func (a *App) GetFilePath() (string, error) {
	return runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择文件",
		Filters: []runtime.FileFilter{
			{DisplayName: "加密文件", Pattern: "*.enc"},
			{DisplayName: "所有文件", Pattern: "*.*"},
		},
	})
}
func main() {
	app := &App{}

	err := wails.Run(&options.App{
		Title:  "FileCrypt",
		Width:  800,
		Height: 600,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup: app.startup,
		Bind:      []interface{}{app},
	})

	if err != nil {
		panic(err)
	}
}
