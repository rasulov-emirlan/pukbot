package fs

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/rasulov-emirlan/pukbot/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type fileSystem struct {
	service *drive.Service
}

func NewFileSystem(cfg config.GoogleFS) (FileSystem, error) {

	ctx := context.Background()
	var config *oauth2.Config
	b, err := ioutil.ReadFile("googleapi_credentials.json")
	if err != nil {
		config = &oauth2.Config{
			ClientID:     cfg.Credentials.ClientID,
			ClientSecret: cfg.Credentials.ClientSecret,
			RedirectURL:  cfg.Credentials.RedirectURIs[0],
			Scopes:       []string{drive.DriveScope},
			Endpoint: oauth2.Endpoint{
				AuthURL:  cfg.Credentials.AuthURI,
				TokenURL: cfg.Credentials.TokenURI,
			},
		}
	}

	// If modifying these scopes, delete your previously saved token.json.
	if err == nil {
		config, err = google.ConfigFromJSON(b, drive.DriveScope)
		if err != nil {
			return nil, err
		}
	}

	client := getClient(config, cfg)
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	return &fileSystem{
		service: srv,
	}, nil
}

func (fs *fileSystem) UploadFile(name string, mimeType string, content io.Reader, folderID string) (string, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{folderID},
	}
	// mimtype shows the type of a file
	// TODO compress file in here
	// TODO if content is empty we will not upload anything
	// we will return a link to a random avatar or
	file, err := fs.service.Files.Create(f).Media(content).Do()
	if err != nil {
		return "", fmt.Errorf("FileServer: %v", err)
	}

	fileLink, err := fs.сreatePublicLink(file.Id)

	return fileLink, err
}

func (fs *fileSystem) сreatePublicLink(fileID string) (string, error) {
	_, err := fs.service.Permissions.Create(fileID, &drive.Permission{
		Role: "reader",
		Type: "anyone",
	}).Do()
	if err != nil {
		return "", err
	}
	return fileID, err
}

func (fs *fileSystem) DeleteFile(filename string) error {
	return fs.service.Files.Delete(filename).Do()
}

func getClient(config *oauth2.Config, cfg config.GoogleFS) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = &oauth2.Token{}
		tok.AccessToken = cfg.Token.AccessToken
		tok.RefreshToken = cfg.Token.RefreshToken
		tok.Expiry = cfg.Token.Expiry
		tok.TokenType = cfg.Token.TokenType

		// tok = getTokenFromWeb(config)
		// saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
