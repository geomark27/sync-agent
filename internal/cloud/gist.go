package cloud

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const gistAPIBase = "https://api.github.com"

// GistProvider implementa Provider sobre la API de GitHub Gists.
type GistProvider struct {
	token  string
	gistID string
	client *http.Client
}

// Verificación en tiempo de compilación de que cumple la interfaz.
var _ Provider = (*GistProvider)(nil)

// NewGistProvider crea un cliente para un Gist concreto.
func NewGistProvider(token, gistID string) *GistProvider {
	return &GistProvider{
		token:  token,
		gistID: gistID,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

type gistFile struct {
	Content   string `json:"content"`
	Truncated bool   `json:"truncated"`
	RawURL    string `json:"raw_url"`
}

type gistResponse struct {
	Files map[string]gistFile `json:"files"`
}

// Pull descarga todos los archivos del Gist.
func (g *GistProvider) Pull(ctx context.Context) (map[string]string, error) {
	url := fmt.Sprintf("%s/gists/%s", gistAPIBase, g.gistID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	g.setHeaders(req)

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("pull falló (%s): %s", resp.Status, string(body))
	}

	var gr gistResponse
	if err := json.NewDecoder(resp.Body).Decode(&gr); err != nil {
		return nil, err
	}

	files := make(map[string]string, len(gr.Files))
	for name, f := range gr.Files {
		content := f.Content
		// GitHub trunca el contenido de archivos grandes (>1MB); en ese caso
		// hay que recuperarlo desde la URL "raw".
		if f.Truncated && f.RawURL != "" {
			content, err = g.fetchRaw(ctx, f.RawURL)
			if err != nil {
				return nil, fmt.Errorf("no se pudo recuperar %s: %w", name, err)
			}
		}
		files[name] = content
	}
	return files, nil
}

type gistPushFile struct {
	Content string `json:"content"`
}

type gistPushBody struct {
	Files map[string]gistPushFile `json:"files"`
}

// Push crea/actualiza los archivos indicados mediante un PATCH al Gist.
func (g *GistProvider) Push(ctx context.Context, files map[string]string) error {
	if len(files) == 0 {
		return nil
	}

	payload := gistPushBody{Files: make(map[string]gistPushFile, len(files))}
	for name, content := range files {
		payload.Files[name] = gistPushFile{Content: content}
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/gists/%s", gistAPIBase, g.gistID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	g.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("push falló (%s): %s", resp.Status, string(body))
	}
	return nil
}

func (g *GistProvider) fetchRaw(ctx context.Context, rawURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+g.token)

	resp, err := g.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("estado inesperado: %s", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (g *GistProvider) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
}
