// dharma is a supplementary corporation community tool for Eve Online.
// Copyright (C) 2021 Cory Slep
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package i18n

import (
	"github.com/cjslep/dharma/internal/data"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Messages struct {
	l *i18n.Localizer
}

func New(b *i18n.Bundle, langs ...string) *Messages {
	return &Messages{
		l: i18n.NewLocalizer(b, langs...),
	}
}

func (m *Messages) PageNotFound() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "pageNotFound",
			Description: "Notifies the user that the webpage was not found",
			Other:       "The requested page was not found.",
		},
	})
}

func (m *Messages) BadRequest() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "badRequest",
			Description: "Notifies the user that the browser issued a bad request",
			Other:       "Your browser issued a bad request.",
		},
	})
}

func (m *Messages) MethodNotAllowed() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "methodNotAllowed",
			Description: "Notifies the user that the browser issued an unsupported HTTP method",
			Other:       "Your browser issued a request with a method that is disallowed",
		},
	})
}

func (m *Messages) UhOh() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "uhOh",
			Description: "Exclaimation that a mistake or accident happened",
			Other:       "Uh oh!",
		},
	})
}

func (m *Messages) InternalServerError() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "internalServerError",
			Description: "Notifies the user that the server had an internal error",
			Other:       "An unrecoverable error occurred within the server.",
		},
	})
}

func (m *Messages) UedamaGateCampRequest() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "uedamaGateCampRequest",
			Description: "Jovial comment that perhaps a user's HTTP request was gate-camped and blown up in an EVE Online system named Uedama",
			Other:       "Perhaps your request was blown up in an Uedama gate camp?",
		},
	})
}

func (m *Messages) ConsiderFileABug() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "fileABug",
			Description: "Prompts the user to click on this text in order to file an issue against the software",
			Other:       "Found a problem? Please consider filing a bug report",
		},
	})
}

func (m *Messages) Email() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "email",
			Description: "E-Mail label for the login & register pages",
			Other:       "Email",
		},
	})
}

func (m *Messages) Password() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "password",
			Description: "Password label for the login & register pages",
			Other:       "Password",
		},
	})
}

func (m *Messages) ConfirmPassword() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "confirmPassword",
			Description: "Label for the box that the user must type a matching password on the register page",
			Other:       "Confirm Password",
		},
	})
}

func (m *Messages) Authorize() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "authorize",
			Description: "Authorize label for the OAuth2 authorization page",
			Other:       "Authorize",
		},
	})
}

func (m *Messages) Login() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "login",
			Description: "Login label for the button to submit the login page",
			Other:       "Login",
		},
	})
}

func (m *Messages) LoginError() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "loginError",
			Description: "Generic login error message when email, credentials, or something else is invalid",
			Other:       "Invalid email or password",
		},
	})
}

func (m *Messages) Forum() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "forum",
			Description: "Forum is the name of a bulletin-board like system of long-form posts and replies",
			Other:       "Forum",
		},
	})
}

func (m *Messages) Corporation() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "corporation",
			Description: "The \"corporation\" is the smallest guild-equivalent in EvE Online",
			Other:       "Corporation",
		},
	})
}

func (m *Messages) Activities() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "activities",
			Description: "The categorical word for generally doing things in the game",
			Other:       "Activities",
		},
	})
}

func (m *Messages) Social() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "social",
			Description: "The categorical word for generally interacting with other people",
			Other:       "Social",
		},
	})
}

func (m *Messages) Miscellaneous() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "miscellaneous",
			Description: "The categorical word for a collection of numerous other small arbitrary topics",
			Other:       "Miscellaneous",
		},
	})
}

func (m *Messages) TagName(id string) (string, error) {
	switch id {
	case data.Announce.ID:
		return m.l.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "announce",
				Description: "A word for announcements and broadcasts",
				Other:       "Announcements",
			},
		})
	case data.Events.ID:
		return m.l.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "events",
				Description: "A word for occurrences and events",
				Other:       "Events",
			},
		})
	case data.Discuss.ID:
		return m.l.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "discuss",
				Description: "A word for general discussions",
				Other:       "Discuss",
			},
		})
	case data.Fleet.ID:
		return m.l.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "fleet",
				Description: "A word for a fleet of ships in Eve Online",
				Other:       "Fleet",
			},
		})
	case data.Industry.ID:
		return m.l.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "industry",
				Description: "A word for the industry in EvE Online",
				Other:       "Industry",
			},
		})
	case data.Market.ID:
		return m.l.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "market",
				Description: "A word for the Market in EvE Online",
				Other:       "Market",
			},
		})
	case data.PVP.ID:
		return m.l.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "pvp",
				Description: "An acronym or word for player-versus-player gameplay",
				Other:       "PVP",
			},
		})
	case data.PVE.ID:
		return m.l.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "pve",
				Description: "An acronym or word for player-versus-environment gameplay",
				Other:       "PVE",
			},
		})
	case data.Relations.ID:
		return m.l.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "relations",
				Description: "A word for political relationships between corporations or alliances within EvE",
				Other:       "Relations",
			},
		})
	case data.Intel.ID:
		return m.l.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "intel",
				Description: "A word for dealing with secret or confidential intelligence",
				Other:       "Intel",
			},
		})
	case data.Justice.ID:
		return m.l.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "justice",
				Description: "A word for dealing with player actions that the corporation would deem criminal or unwanted",
				Other:       "Justice",
			},
		})
	case data.QNA.ID:
		return m.l.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "question-and-answer",
				Description: "A word for dedicating a space for asking questions and receiving answers",
				Other:       "Q&A",
			},
		})
	case data.OffTopic.ID:
		return m.l.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "offtopic",
				Description: "The categorical word for not being topical for any other topics",
				Other:       "Off-Topic",
			},
		})
	case data.Uncategorized.ID:
		fallthrough
	default:
		return m.l.Localize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "uncategorized",
				Description: "The word for not having a category",
				Other:       "Uncategorized",
			},
		})
	}
}

func (m *Messages) Languages() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "languages",
			Description: "The word for multiple languages",
			Other:       "Languages",
		},
	})
}

func (m *Messages) NoPostsYet() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "noPostsYet",
			Description: "Describes the lack of posts yet in a forum category",
			Other:       "No posts yet.",
		},
	})
}

func (m *Messages) BeTheFirstToPost() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "beTheFirstToPost",
			Description: "Suggestion to the user to be the first one to post in an empty forum category",
			Other:       "Be the first one to create a post in this category!",
		},
	})
}

func (m *Messages) CreateNewPost() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "createNewPost",
			Description: "Button text prompting user to create a new forum post",
			Other:       "Create New Post",
		},
	})
}

func (m *Messages) ContentLabel() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "contentLabel",
			Description: "Label for the text box for entering the content of a new post",
			Other:       "Content",
		},
	})
}

func (m *Messages) TitleLabel() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "titleLabel",
			Description: "Label for the text box for entering the title of a new post",
			Other:       "Title",
		},
	})
}

func (m *Messages) TagLabel() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "tagLabel",
			Description: "Label for the select box for selecting a tag for a new post",
			Other:       "Tag",
		},
	})
}

func (m *Messages) RegisterNewAccount() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "registerNewAccount",
			Description: "Label for page registering new account",
			Other:       "Register New Account",
		},
	})
}

func (m *Messages) Register() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "register",
			Description: "Button for registering a new account",
			Other:       "Register",
		},
	})
}

func (m *Messages) Username() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "username",
			Description: "A user account's name",
			Other:       "Username",
		},
	})
}

func (m *Messages) PasswordsDoNotMatch() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "passwordsDoNotMatch",
			Description: "Error message shown when, upon registration, the password and confirm-password fields did not match",
			Other:       "Passwords do not match",
		},
	})
}

func (m *Messages) UsernameNotUnique() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "usernameNotUnique",
			Description: "Error message shown when, upon registration, the user's requested username is not unique",
			Other:       "Username is already taken",
		},
	})
}

func (m *Messages) EmailNotUnique() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "EmailNotUnique",
			Description: "Error message shown when, upon registration, the user's requested email is not unique",
			Other:       "Email is already taken",
		},
	})
}

func (m *Messages) UnknownRegistrationError() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "unknownRegistrationError",
			Description: "Error message shown when, upon registration, an error occurred but a programming error means it isn't being displayed properly. Should include a call to contact the site administrator and/or file a bug",
			Other:       "A registration error occurred, please contact the site administrator and/or file a bug against the software",
		},
	})
}

func (m *Messages) VerifyEmail() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "verifyEmail",
			Description: "Request to the user to verify their email",
			Other:       "Please verify your email address.",
		},
	})
}

func (m *Messages) ThanksForRegistering() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "thanksForRegistering",
			Description: "A title to tell the user \"thank-you for registering\" which is shown on the verify email page",
			Other:       "Thank You For Registering",
		},
	})
}

func (m *Messages) NotVerifiedEmailYet() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "notVerifiedEmailYet",
			Description: "A title to tell the user their email is not yet verified, which is shown on the verify email page when they get redirected",
			Other:       "Your Email Is Not Yet Verified",
		},
	})
}

func (m *Messages) CheckEmailForVerification() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "checkEmailForVerification",
			Description: "Instructions on how to verify the email, warning the user to check their spam folder",
			Other:       "Please check your email and its spam folder for a message containing a verification link. Clicking this link will verify your account.",
		},
	})
}

func (m *Messages) PleaseValidateEmailSubject() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "pleaseValidateEmailSubject",
			Description: "Email subject for the email sent for new registrants to verify their account",
			Other:       "Please Verify Your Account",
		},
	})
}

func (m *Messages) PleaseValidateEmailBodyText() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "pleaseValidateEmailSubject",
			Description: "Text email body for the email sent for new registrants to verify their account, must have a %s for the URL",
			Other:       "Please go to the following link to verify your account: %s",
		},
	})
}

func (m *Messages) VerifySuccess() (string, error) {
	return m.l.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:          "verifySuccess",
			Description: "Message displayed to a user after they have verified their email address by clicking a link in the email",
			Other:       "Your email has been successfully verified! You may now log in.",
		},
	})
}
