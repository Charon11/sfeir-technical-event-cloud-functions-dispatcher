package subject

import (
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
	"log"
)

func InitializeAppDefault() (*firebase.App, error) {
	// [START initialize_app_default_golang]

	saJson, err := base64.StdEncoding.DecodeString("ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIsCiAgInByb2plY3RfaWQiOiAic2ZlaXJsdXh0ZWNobmljYWxldmVudCIsCiAgInByaXZhdGVfa2V5X2lkIjogIjE2NTliY2Q2NzdiYWIyZDEyZjY1NjVkMTkwNzk5ZjVhMzNjNzBmNDciLAogICJwcml2YXRlX2tleSI6ICItLS0tLUJFR0lOIFBSSVZBVEUgS0VZLS0tLS1cbk1JSUV2QUlCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktZd2dnU2lBZ0VBQW9JQkFRQ3BvTGhha3VJMjZLbE1cblVTbDUrdmpwWmdzUnlLc2UySksySjBwazNCd2dZZDc0endGc0pNM1M0cHV2RUw0bVI0bFhaMWpYYi9maUxFVXFcblV0S2QxalpGOXRXKzF6bkoxS1RwTWMydXNoLzlFMWtaeUsvejBjZ3JzNElrSUdILzFEeGJabXd6bWxBNVZnTi9cbjBmWktRY2p5bUlkeVBrZUdlQUNTdlFnVkhmc1ZFRXFMWTNSL1dvZ05qNDYrWjlnRnk0SW9lRXpRY2ZRRzhVaFpcbitFSHJTMGxYdUwxTVpIeVNGdXBLNUZTSjdjZDIrTmp3cEdFbUpQZ202cS9NRHd0enhvQXhMb3A1VXF1SDB4ZU9cbndLYzNFZjY5R3ZCZWxYcTd0eW5TUEpJOFZxUWFqVEFMWDBDZUVQL2VBdUR5YnU3a1NSam5hZEg4eEtoWTQ2ZWtcbksyWkxlUlJmQWdNQkFBRUNnZ0VBVHRVWVNGT3lHUldvdElTVnM1QmRSTTg5UHI4Vk0ra3Z5Y0xaMHFUdTEyZlVcbndhb212WlVmS0s1Uzd4bm5YUS9xOXJsYmN3Z2cyalo3MEc4Y1hla3pZUEdWcGJNTnNzeUY5YkJGS0RhMUloN3ZcblRrblQ2QXJyTGRBbm45V0M1ZXZEUHpFTEFUSnFyVExadm5vY0xhZUVGYWFvY0FJY3FKUTFvL3lBK0p0dDJIdnZcbkJLVXpYZVZvTFFaa2Z6MWVvQ1ZIMVRLcDIyMWpyY2YvOWdCMjAxZTFQdlpZSjVzZ2RsSFZCNVNHN3VreFI3UHJcbkZuNkQvejY5a0pNMU51UzFYTFF5NzJHQWxEMUd0czdSZnZuU2VnU0ZaZVVsWmZENzR4dEE3UjQvV1FTNm52ekVcblVMUVpPMHljQ2YzN2NweXlwcmdGNGFpT002ZkFXT21uVU1rUWpiRDVPUUtCZ1FEZ2xpL1d2citVZDZvRzI4dUdcblp5aW5hdUl3L2JaOFhEYVFFclVxeHFHNWMwQzNKWU1zcjFPekdsOU9YK0c0M0I4SGhveU1HMTdkUFpxelgwSUlcblU5NlkyaDN1czdKR3A1RklmMlRuSjB2VjhnZGVER0NRdWpiT0RMV3pzVVluTGRicmJmenkxeFVhb3FSU2ZQeUZcbmJpbUZ3VVNmcXRHMGZYVEtUVDFkaUFBSzF3S0JnUURCV3B5MzUzajlhRlBWN2g1UFZFSjlVbGhUZldBNWNrSUJcbjFnaGNYN2Q0ekVDNXZLRFJtRGI3SmVZMFVVWGNwVUU0VWNqV0g4amtuWm5JK2VnL09PUENQa2RTZWh5VW9pV2Jcbll3ZHZNRVhKd3RiS213SkpmZXJpbUhZTkJwUFhzdEYyN2FRdkdGbjFlVmxLaVRLanpGQUNONmRkT2dZdS9HMXdcbkljU3FXS2ZadVFLQmdHc3JINjdnblBqUzFWNnFlWWNzS0xrakJzYUwrdzJDejBLV3VyNnJ4RGFEYWNrN2JFbmhcbmJCWldLazV5OEhwUEI3dUdtQUN2ZXVnbDRuRmZybG5jODZhS3hxZkdOZmNETlErY1F1RU8zbUE4T3duRTdEdURcbnNHMUlvVmdhNnJmOVpzWTNXUEhrY3B5Z0tidDNDdVQ4K1hGckZUei92VXZjWmVPM1VlWVU1TGhIQW9HQUI1K2lcbi9qMUtML09sR3BKQTc2L0t6MVluMVdMa0lGKzQ2b21kMnNhTFhWL3dUV0o2bE1rRG9mTmttRHQ4SGE4R1cwb2Vcbk9STVh3S3ZXSEc1K0VjQVVldHdwdzc2ODBiOXk3Q1dEYlliaHVZck5IVC92WHladjFHOWVlRFVDTDBjRnEwTERcbmJuZnZWQlpzeHR1Uis1TmF0RHV0aGFkOEJ1NDlBbWJTY2tUOG1VRUNnWUJ0VGJCT2RBaEFnUjYwdUFPckxrdnhcbkNvRGZoMDNsUVdGdUw3eUV6ZS9MRXJ5UE1jQVdGQjU0TFRxSnRpL1NjZVFnNU9SU1VZNWxZby90TjBWeU0wRjhcbkx2Z0pzQjNobkJ4SUhtdzBFbmNpeEY3MCs1UDFqTXZ4blVYcXRnb2RuelE5U1l6UlhkS1VZamVFYy90U2dlTkZcbktRLzZQUDRyNmxVNDFYM2hBaWk0T2c9PVxuLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLVxuIiwKICAiY2xpZW50X2VtYWlsIjogImZpcmViYXNlLWFkbWluc2RrLW1pa3N0QHNmZWlybHV4dGVjaG5pY2FsZXZlbnQuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20iLAogICJjbGllbnRfaWQiOiAiMTExNTIwNjA0ODgxMDcwNDQzMzg4IiwKICAiYXV0aF91cmkiOiAiaHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29tL28vb2F1dGgyL2F1dGgiLAogICJ0b2tlbl91cmkiOiAiaHR0cHM6Ly9vYXV0aDIuZ29vZ2xlYXBpcy5jb20vdG9rZW4iLAogICJhdXRoX3Byb3ZpZGVyX3g1MDlfY2VydF91cmwiOiAiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vb2F1dGgyL3YxL2NlcnRzIiwKICAiY2xpZW50X3g1MDlfY2VydF91cmwiOiAiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vcm9ib3QvdjEvbWV0YWRhdGEveDUwOS9maXJlYmFzZS1hZG1pbnNkay1taWtzdCU0MHNmZWlybHV4dGVjaG5pY2FsZXZlbnQuaWFtLmdzZXJ2aWNlYWNjb3VudC5jb20iCn0K")
	if err != nil {
		log.Printf("decode error:", err)
		return nil, err
	}
	opt := option.WithCredentialsJSON(saJson)
	config := &firebase.Config{ProjectID: "sfeirluxtechnicalevent"}

	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Printf("error initializing app: %v\n", err)
		return nil, err
	}
	// [END initialize_app_default_golang]

	return app, nil
}

func VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {

	app, err := InitializeAppDefault()
	if err != nil {
		log.Printf("error init app: %v\n", err)
		return nil, err
	}
	// [START verify_id_token_golang]
	client, err := app.Auth(context.Background())
	if err != nil {
		log.Printf("error getting Auth client: %v\n", err)
		return nil, err
	}

	token, err := client.VerifyIDToken(ctx, idToken)
	if err != nil {
		log.Printf("error verifying ID token: %v\n", err)
		return nil, err
	}

	log.Printf("Verified ID token: %v\n", token)
	// [END verify_id_token_golang]

	return token, nil
}
