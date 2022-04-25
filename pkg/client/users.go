package client

import (
	"context"
	gen "github.com/hashicorp-dev-advocates/waypoint-client/pkg/waypoint"
)

type UserRef interface {
	Ref() string
}

type UserId string

func (u *UserId) Ref() string {
	return string(*u)
}

type Username string

func (u *Username) Ref() string {
	return string(*u)
}

// CreateToken returns a waypoint token for the current user context
func (c *waypointImpl) CreateToken(ctx context.Context, id UserRef) (string, error) {

	var user *gen.Ref_User

	switch id.(type) {
	case *UserId:
		user = &gen.Ref_User{
			Ref: &gen.Ref_User_Id{Id: &gen.Ref_UserId{Id: id.Ref()}},
		}

	case *Username:
		user = &gen.Ref_User{
			Ref: &gen.Ref_User_Username{Username: &gen.Ref_UserUsername{Username: id.Ref()}},
		}

	}
	if id != nil {
		user = &gen.Ref_User{
			Ref: &gen.Ref_User_Id{Id: &gen.Ref_UserId{Id: id.Ref()}},
		}
	}
	gtr := &gen.LoginTokenRequest{
		User:    user,
		Trigger: false,
	}

	token, err := c.client.GenerateLoginToken(ctx, gtr)
	if err != nil {
		return "", err
	}

	return token.Token, nil
}

// InviteUser returns a invitation token for a new user to the Waypoint server
func (c *waypointImpl) InviteUser(ctx context.Context, InitialUsername string, TokenTtl string) (string, error) {

	tis := &gen.Token_Invite_Signup{
		InitialUsername: InitialUsername,
	}

	uir := &gen.InviteTokenRequest{
		Duration:         TokenTtl,
		Signup:           tis,
		UnusedEntrypoint: nil,
	}

	inviteToken, err := c.client.GenerateInviteToken(ctx, uir)
	if err != nil {
		return "", err
	}

	return inviteToken.Token, nil

}

func (c *waypointImpl) AcceptInvitation(ctx context.Context, InviteToken string) (string, error) {
	citr := &gen.ConvertInviteTokenRequest{
		Token: InviteToken,
	}

	si, err := c.client.ConvertInviteToken(ctx, citr)
	if err != nil {
		return "", err
	}

	return si.Token, nil

}

func (c *waypointImpl) GetUser(ctx context.Context, username Username) (*gen.User, error) {

	gur := &gen.GetUserRequest{
		User: &gen.Ref_User{
			Ref: &gen.Ref_User_Username{Username: &gen.Ref_UserUsername{Username: username.Ref()}},
		},
	}

	gu, err := c.client.GetUser(ctx, gur)
	if err != nil {
		return nil, err
	}

	return gu.User, nil

}

func (c *waypointImpl) DeleteUser(ctx context.Context, id UserId) (string, error) {
	dur := &gen.DeleteUserRequest{User: &gen.Ref_User{
		Ref: &gen.Ref_User_Id{Id: &gen.Ref_UserId{Id: id.Ref()}},
	}}

	_, err := c.client.DeleteUser(ctx, dur)
	if err != nil {
		return "", err
	}

	return "User deleted", nil

}

