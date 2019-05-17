package mgo

import (
	"fmt"
	"path"

	"github.com/madappgang/identifo/jwt"
	jwtService "github.com/madappgang/identifo/jwt/service"
	"github.com/madappgang/identifo/model"
	"github.com/madappgang/identifo/storage/mongo"
)

// NewComposer creates new database composer.
func NewComposer(settings model.ServerSettings, options ...func(*DatabaseComposer) error) (*DatabaseComposer, error) {
	c := DatabaseComposer{
		settings:        settings,
		newAppStorage:   mongo.NewAppStorage,
		newUserStorage:  mongo.NewUserStorage,
		newTokenStorage: mongo.NewTokenStorage,
	}

	for _, option := range options {
		if err := option(&c); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

// InitAppStorage returns an argument that sets the appStorage initialization function.
func InitAppStorage(initAS func(*mongo.DB) (model.AppStorage, error)) func(*DatabaseComposer) error {
	return func(dc *DatabaseComposer) error {
		dc.newAppStorage = initAS
		return nil
	}
}

// InitUserStorage returns an argument that sets the userStorage initialization function.
func InitUserStorage(initUS func(*mongo.DB) (model.UserStorage, error)) func(*DatabaseComposer) error {
	return func(dc *DatabaseComposer) error {
		dc.newUserStorage = initUS
		return nil
	}
}

// InitTokenStorage returns an argument that sets the tokenStorage initialization function.
func InitTokenStorage(initTS func(*mongo.DB) (model.TokenStorage, error)) func(*DatabaseComposer) error {
	return func(dc *DatabaseComposer) error {
		dc.newTokenStorage = initTS
		return nil
	}
}

// DatabaseComposer composes MongoDB services.
type DatabaseComposer struct {
	settings        model.ServerSettings
	newAppStorage   func(*mongo.DB) (model.AppStorage, error)
	newUserStorage  func(*mongo.DB) (model.UserStorage, error)
	newTokenStorage func(*mongo.DB) (model.TokenStorage, error)
}

// Compose composes all services with MongoDB support.
func (dc *DatabaseComposer) Compose() (
	model.AppStorage,
	model.UserStorage,
	model.TokenStorage,
	jwtService.TokenService,
	error,
) {

	db, err := mongo.NewDB(dc.settings.DBEndpoint, dc.settings.DBName)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	appStorage, err := dc.newAppStorage(db)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	userStorage, err := dc.newUserStorage(db)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tokenStorage, err := dc.newTokenStorage(db)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	tokenServiceAlg, ok := jwt.StrToTokenServiceAlg[dc.settings.Algorithm]
	if !ok {
		return nil, nil, nil, nil, fmt.Errorf("Unknow token service algoritm %s", dc.settings.Algorithm)
	}

	tokenService, err := jwtService.NewJWTokenService(
		path.Join(dc.settings.PEMFolderPath, dc.settings.PrivateKey),
		path.Join(dc.settings.PEMFolderPath, dc.settings.PublicKey),
		dc.settings.Issuer,
		tokenServiceAlg,
		tokenStorage,
		appStorage,
		userStorage,
	)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return appStorage, userStorage, tokenStorage, tokenService, nil
}
