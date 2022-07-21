package handlers

import (
	"bytes"
	"encoding/gob"
	"moonbrain/app/configs"
	"moonbrain/app/models"
	"moonbrain/app/services"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/rs/zerolog/log"
	"github.com/shareed2k/goth_fiber"
)

type OAuthRedirectData struct {
	RedirectURL string `json:"redirectUrl"`
}

func mapToUser(user goth.User) *models.User {
	return &models.User{
		Provider:            user.Provider,
		Email:               user.Email,
		Name:                user.Name,
		NickName:            user.NickName,
		AvatarURL:           user.AvatarURL,
		ExternalID:          user.UserID,
		FirstName:           user.FirstName,
		LastName:            user.LastName,
		Token:               user.AccessToken,
		RefreshToken:        &user.RefreshToken,
		TokenExpirationDate: user.ExpiresAt,
		ProfileURL:          user.RawData["html_url"].(string),
	}
}

// TODO: master refactor this code.
func RegisterAuthHandler(app fiber.Router, userService *services.UserService, config configs.Config) {
	goth.UseProviders(
		github.New(config.GithubID, config.GithubSecret, config.BackendHost+"/auth/github/callback"),
	)

	app.Get("/auth/github/login", func(c *fiber.Ctx) error {
		url, err := goth_fiber.GetAuthURL(c)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		// return c.Redirect(url, fiber.StatusTemporaryRedirect)
		data := Response[OAuthRedirectData]{
			Data: OAuthRedirectData{
				RedirectURL: url,
			},
		}
		return c.JSON(data)
	})

	app.Get("/auth/github/callback", func(c *fiber.Ctx) error {
		user, err := goth_fiber.CompleteUserAuth(c)
		if err != nil {
			log.Error().Err(err).Msgf("auth handlers: github auth handler: complete user auth")
			return c.Status(500).SendString("Internal server error")
		}
		var userBytes bytes.Buffer
		enc := gob.NewEncoder(&userBytes)
		err = enc.Encode(user)
		if err != nil {
			log.Error().Err(err).Msgf("auth handlers: github auth handler: encode user: %v", err)
			return c.Status(500).SendString("Internal server error")
		}
		u, err := userService.Login(*mapToUser(user))
		if err != nil {
			log.Error().Err(err).Msgf("auth handlers: github auth handler: login user %v", err)
			return c.Status(500).SendString("Internal server error")
		}
		redirectURL := config.ClientAddress + "/auth/login"
		url, err := url.Parse(redirectURL)
		if err != nil {
			log.Error().Err(err).Msgf("auth handlers: github auth handler: parse redirect url %v", err)
		}
		q := url.Query()
		q.Set("token", u.Token)
		q.Set("username", u.NickName)
		q.Set("avatarUrl", u.AvatarURL)
		q.Set("email", u.Email)
		q.Set("profileUrl", u.ProfileURL)
		url.RawQuery = q.Encode()
		log.Info().Msgf("auth handlers: github auth handler: redirect to %s, %v", url.String())
		return c.Redirect(url.String())

	})

	app.Get("/auth/logout", func(c *fiber.Ctx) error {
		if err := goth_fiber.Logout(c); err != nil {
			log.Fatal().Err(err).Msgf("auth handlers: github auth handler: logout")
			return c.Status(500).SendString("Internal server error")
		}

		c.SendString("logout")
		return c.Status(200).JSON(struct{}{})
	})
}