package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const (
	repoOwner  = "kanisterio"
	repoName   = "kanister"
	targetPath = "examples"
)

type GitHubContent struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	URL     string `json:"url"`
	HTMLURL string `json:"html_url"`
}

func main() {
	a := app.New()
	w := a.NewWindow("Kanister Blueprint Search")

	tokenInput := widget.NewPasswordEntry()
	tokenInput.SetPlaceHolder("Enter GitHub access token")

	searchInput := widget.NewEntry()
	searchInput.SetPlaceHolder("Enter search term")

	resultsLabel := widget.NewLabel("")

	resultsList := widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {},
	)

	resultsScroll := container.NewVScroll(resultsList)
	resultsScroll.SetMinSize(fyne.NewSize(0, 200))

	yamlViewer := widget.NewEntry()
	yamlViewer.MultiLine = true
	yamlViewer.Wrapping = fyne.TextWrapWord

	progressBar := widget.NewProgressBar()

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Selected Blueprint URL")

	openURLButton := widget.NewButton("Open in Browser", func() {
		err := openURL(urlEntry.Text)
		if err != nil {
			dialog.ShowError(err, w)
		}
	})

	openEditorButton := widget.NewButton("Open in Editor", func() {
		openInEditor(yamlViewer.Text)
	})

	viewButton := widget.NewButton("View", func() {
		displayYAMLFromURL(urlEntry.Text, tokenInput.Text, yamlViewer, w)
	})

	searchButton := widget.NewButton("Search", func() {
		accessToken := tokenInput.Text
		searchTerm := searchInput.Text

		if accessToken == "" {
			dialog.ShowInformation("Error", "Please enter a GitHub access token", w)
			return
		}

		resultsLabel.SetText("Searching...")
		progressBar.SetValue(0)
		go func() {
			blueprints, err := searchFiles(context.Background(), targetPath, searchTerm, accessToken, progressBar)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error searching files: %v", err), w)
				return
			}

			resultsLabel.SetText(fmt.Sprintf("%d files found containing '%s'", len(blueprints), searchTerm))
			resultsList.Length = func() int { return len(blueprints) }
			resultsList.UpdateItem = func(id widget.ListItemID, item fyne.CanvasObject) {
				item.(*widget.Label).SetText(blueprints[id])
			}
			resultsList.OnSelected = func(id widget.ListItemID) {
				selectedURL := blueprints[id]
				urlEntry.SetText(selectedURL)
			}
			resultsList.Refresh()
		}()
	})

	
	leftContainer := container.NewVBox(
		widget.NewLabel("GitHub Access Token:"),
		tokenInput,
		widget.NewLabel("Search Term:"),
		searchInput,
		searchButton,
		widget.NewLabel("Results:"),
		resultsLabel,
		progressBar,
		resultsScroll,
		widget.NewLabel("Selected Blueprint URL:"),
		urlEntry,
		container.NewHBox(openURLButton, openEditorButton, viewButton),
	)

	
	rightContainer := container.NewMax(
		yamlViewer,
	)

	
	content := container.NewHSplit(
		container.NewPadded(leftContainer),
		container.NewPadded(rightContainer),
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(1000, 600))
	w.ShowAndRun()
}

func searchFiles(ctx context.Context, path, searchTerm, accessToken string, progressBar *widget.ProgressBar) ([]string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", repoOwner, repoName, path), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("Authorization", "token "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request to GitHub API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var contents []GitHubContent
	err = json.Unmarshal(body, &contents)
	if err != nil {
		return nil, fmt.Errorf("error decoding GitHub API response: %v", err)
	}

	var blueprints []string
	totalItems := len(contents)
	for i, content := range contents {
		if content.Type == "file" && strings.HasSuffix(content.Name, ".yaml") && strings.Contains(strings.ToLower(content.Name), strings.ToLower(searchTerm)) {
			blueprints = append(blueprints, content.HTMLURL)
		} else if content.Type == "dir" {
			subBlueprints, err := searchFiles(ctx, path+"/"+content.Name, searchTerm, accessToken, progressBar)
			if err != nil {
				return nil, err
			}
			blueprints = append(blueprints, subBlueprints...)
		}
		progressBar.SetValue(float64(i+1) / float64(totalItems))
		time.Sleep(10 * time.Millisecond) 
	}
	return blueprints, nil
}

func fetchYAMLContent(url string) (string, error) {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return "", fmt.Errorf("error making request to fetch YAML content: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetch YAML content returned status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
}

func displayYAMLFromURL(url, accessToken string, viewer *widget.Entry, w fyne.Window) {
	rawURL := strings.Replace(url, "https://github.com", "https://raw.githubusercontent.com", 1)
	rawURL = strings.Replace(rawURL, "/blob/", "/", 1)

	content, err := fetchYAMLContent(rawURL)
	if err != nil {
		dialog.ShowError(fmt.Errorf("error fetching YAML content: %v", err), w)
		return
	}

	viewer.SetText(content)
}

func openURL(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default: 
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}

func openInEditor(content string) {
	file, err := ioutil.TempFile("", "*.yaml")
	if err != nil {
		dialog.ShowError(fmt.Errorf("error creating temp file: %v", err), nil)
		return
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		dialog.ShowError(fmt.Errorf("error writing to temp file: %v", err), nil)
		return
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("notepad.exe", file.Name())
	case "darwin":
		cmd = exec.Command("open", "-e", file.Name())
	default: 
		cmd = exec.Command("xdg-open", file.Name())
	}
	err = cmd.Start()
	if err != nil {
		dialog.ShowError(fmt.Errorf("error opening editor: %v", err), nil)
	}
}
