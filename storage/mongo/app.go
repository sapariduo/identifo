package mongo

import (
	"encoding/json"

	"github.com/madappgang/identifo/model"
	"gopkg.in/mgo.v2/bson"
)

// AppData is a MongoDb model that implements model.AppData.
type AppData struct {
	appData
}

type appData struct {
	ID                           bson.ObjectId          `bson:"_id,omitempty" json:"id,omitempty"`
	Secret                       string                 `bson:"secret,omitempty" json:"secret,omitempty"`
	Active                       bool                   `bson:"active" json:"active"`
	Name                         string                 `bson:"name,omitempty" json:"name,omitempty"`
	Description                  string                 `bson:"description,omitempty" json:"description,omitempty"`
	Scopes                       []string               `bson:"scopes,omitempty" json:"scopes,omitempty"`
	Offline                      bool                   `bson:"offline" json:"offline"`
	Type                         model.AppType          `bson:"type,omitempty" json:"type,omitempty"`
	RedirectURLs                 []string               `bson:"redirect_urls,omitempty" json:"redirect_urls,omitempty"`
	RefreshTokenLifespan         int64                  `bson:"refresh_token_lifespan,omitempty" json:"refresh_token_lifespan,omitempty"`
	InviteTokenLifespan          int64                  `bson:"invite_token_lifespan,omitempty" json:"invite_token_lifespan,omitempty"`
	TokenLifespan                int64                  `bson:"token_lifespan,omitempty" json:"token_lifespan,omitempty"`
	TokenPayload                 []string               `bson:"token_payload,omitempty" json:"token_payload,omitempty"`
	RegistrationForbidden        bool                   `bson:"registration_forbidden" json:"registration_forbidden"`
	AnonymousRegistrationAllowed bool                   `bson:"anonymous_registration_allowed" json:"anonymous_registration_allowed"`
	TFAStatus                    model.TFAStatus        `bson:"tfa_status" json:"tfa_status"`
	DebugTFACode                 string                 `bson:"debug_tfa_code,omitempty" json:"debug_tfa_code,omitempty"`
	AuthorizationWay             model.AuthorizationWay `bson:"authorization_way,omitempty" json:"authorization_way,omitempty"`
	AuthorizationModel           string                 `bson:"authorization_model,omitempty" json:"authorization_model,omitempty"`
	AuthorizationPolicy          string                 `bson:"authorization_policy,omitempty" json:"authorization_policy,omitempty"`
	RolesWhitelist               []string               `bson:"roles_whitelist,omitempty" json:"roles_whitelist,omitempty"`
	RolesBlacklist               []string               `bson:"roles_blacklist,omitempty" json:"roles_blacklist,omitempty"`
	NewUserDefaultRole           string                 `bson:"new_user_default_role,omitempty" json:"new_user_default_role,omitempty"`
	AppleInfo                    *model.AppleInfo       `bson:"apple_info,omitempty" json:"apple_info,omitempty"`
}

// NewAppData instantiates MongoDB app data model from the general one.
func NewAppData(data model.AppData) (AppData, error) {
	if !bson.IsObjectIdHex(data.ID()) {
		return AppData{}, model.ErrorWrongDataFormat
	}
	return AppData{appData: appData{
		ID:                           bson.ObjectIdHex(data.ID()),
		Secret:                       data.Secret(),
		Active:                       data.Active(),
		Name:                         data.Name(),
		Description:                  data.Description(),
		Scopes:                       data.Scopes(),
		Offline:                      data.Offline(),
		RedirectURLs:                 data.RedirectURLs(),
		RefreshTokenLifespan:         data.RefreshTokenLifespan(),
		InviteTokenLifespan:          data.InviteTokenLifespan(),
		TokenLifespan:                data.TokenLifespan(),
		TokenPayload:                 data.TokenPayload(),
		RegistrationForbidden:        data.RegistrationForbidden(),
		AnonymousRegistrationAllowed: data.AnonymousRegistrationAllowed(),
		TFAStatus:                    data.TFAStatus(),
		AuthorizationWay:             data.AuthzWay(),
		AuthorizationModel:           data.AuthzModel(),
		AuthorizationPolicy:          data.AuthzPolicy(),
		RolesWhitelist:               data.RolesWhitelist(),
		RolesBlacklist:               data.RolesBlacklist(),
		NewUserDefaultRole:           data.NewUserDefaultRole(),
	}}, nil
}

// AppDataFromJSON deserializes app data from JSON.
func AppDataFromJSON(d []byte) (AppData, error) {
	apd := appData{}
	if err := json.Unmarshal(d, &apd); err != nil {
		return AppData{}, err
	}
	return AppData{appData: apd}, nil
}

// MakeAppData creates new MongoDB app data instance.
func MakeAppData(id, secret string, active bool, name, description string, scopes []string, offline bool, redirectURLs []string,
	refreshTokenLifespan, inviteTokenLifespan, tokenLifespan int64, tokenPayload []string, registrationForbidden bool, anonymousRegistrationAllowed bool,
	tfaStatus model.TFAStatus, debugTFACode string, authzWay model.AuthorizationWay, authzModel, authzPolicy string, rolesWhitelist, rolesBlacklist []string, newUserDefaultRole string) (AppData, error) {

	if !bson.IsObjectIdHex(id) {
		return AppData{}, model.ErrorWrongDataFormat
	}
	return AppData{appData: appData{
		ID:                           bson.ObjectIdHex(id),
		Secret:                       secret,
		Active:                       active,
		Name:                         name,
		Description:                  description,
		Scopes:                       scopes,
		Offline:                      offline,
		RedirectURLs:                 redirectURLs,
		RefreshTokenLifespan:         refreshTokenLifespan,
		InviteTokenLifespan:          inviteTokenLifespan,
		TokenLifespan:                tokenLifespan,
		TokenPayload:                 tokenPayload,
		RegistrationForbidden:        registrationForbidden,
		AnonymousRegistrationAllowed: anonymousRegistrationAllowed,
		TFAStatus:                    tfaStatus,
		DebugTFACode:                 debugTFACode,
		AuthorizationWay:             authzWay,
		AuthorizationModel:           authzModel,
		AuthorizationPolicy:          authzPolicy,
		RolesWhitelist:               rolesWhitelist,
		RolesBlacklist:               rolesBlacklist,
		NewUserDefaultRole:           newUserDefaultRole,
	}}, nil
}

// Marshal serializes data to byte array.
func (ad AppData) Marshal() ([]byte, error) {
	return json.Marshal(ad.appData)
}

// ID implements model.AppData interface.
func (ad *AppData) ID() string { return ad.appData.ID.Hex() }

// Secret implements model.AppData interface.
func (ad *AppData) Secret() string { return ad.appData.Secret }

// Active implements model.AppData interface.
func (ad *AppData) Active() bool { return ad.appData.Active }

// Name implements model.AppData interface.
func (ad *AppData) Name() string { return ad.appData.Name }

// Description implements model.AppData interface.
func (ad *AppData) Description() string { return ad.appData.Description }

// Scopes implements model.AppData interface.
func (ad *AppData) Scopes() []string { return ad.appData.Scopes }

// Offline implements model.AppData interface.
func (ad *AppData) Offline() bool { return ad.appData.Offline }

// Type implements model.AppData interface.
func (ad *AppData) Type() model.AppType { return ad.appData.Type }

// RedirectURLs implements model.AppData interface.
func (ad *AppData) RedirectURLs() []string { return ad.appData.RedirectURLs }

// InviteTokenLifespan implements model.AppData interface.
func (ad *AppData) InviteTokenLifespan() int64 { return ad.appData.InviteTokenLifespan }

// RefreshTokenLifespan implements model.AppData interface.
func (ad *AppData) RefreshTokenLifespan() int64 { return ad.appData.RefreshTokenLifespan }

// TokenLifespan implements model.AppData interface.
func (ad *AppData) TokenLifespan() int64 { return ad.appData.TokenLifespan }

// TokenPayload implements model.AppData interface.
func (ad *AppData) TokenPayload() []string { return ad.appData.TokenPayload }

// RegistrationForbidden implements model.AppData interface.
func (ad *AppData) RegistrationForbidden() bool { return ad.appData.RegistrationForbidden }

// AnonymousRegistrationAllowed implements model.AppData interface.
func (ad *AppData) AnonymousRegistrationAllowed() bool { return ad.appData.AnonymousRegistrationAllowed }

// TFAStatus implements model.AppData interface.
func (ad *AppData) TFAStatus() model.TFAStatus { return ad.appData.TFAStatus }

// DebugTFACode implements model.AppData interface.
func (ad *AppData) DebugTFACode() string { return ad.appData.DebugTFACode }

// AuthzWay implements model.AppData interface.
func (ad *AppData) AuthzWay() model.AuthorizationWay { return ad.appData.AuthorizationWay }

// AuthzModel implements model.AppData interface.
func (ad *AppData) AuthzModel() string { return ad.appData.AuthorizationModel }

// AuthzPolicy implements model.AppData interface.
func (ad *AppData) AuthzPolicy() string { return ad.appData.AuthorizationPolicy }

// RolesWhitelist implements model.AppData interface.
func (ad *AppData) RolesWhitelist() []string { return ad.appData.RolesWhitelist }

// RolesBlacklist implements model.AppData interface.
func (ad *AppData) RolesBlacklist() []string { return ad.appData.RolesBlacklist }

// NewUserDefaultRole implements model.AppData interface.
func (ad *AppData) NewUserDefaultRole() string { return ad.appData.NewUserDefaultRole }

// AppleInfo implements model.AppData interface.
func (ad *AppData) AppleInfo() *model.AppleInfo { return ad.appData.AppleInfo }

// SetSecret implements model.AppData interface.
func (ad *AppData) SetSecret(secret string) {
	if ad == nil {
		return
	}
	ad.appData.Secret = secret
}

// Sanitize removes all sensitive data.
func (ad *AppData) Sanitize() {
	if ad == nil {
		return
	}
	ad.appData.Secret = ""
	if ad.appData.AppleInfo != nil {
		ad.appData.AppleInfo.ClientSecret = ""
	}

	ad.appData.AuthorizationWay = ""
	ad.appData.AuthorizationModel = ""
	ad.appData.AuthorizationPolicy = ""
}
