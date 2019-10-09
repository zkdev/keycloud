package protocol_test

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"testing"

	"webauthn/protocol"
)

func TestIsValidAttestation(t *testing.T) {
	for i := range attestationRequests {
		t.Run(fmt.Sprintf("Run %d", i), func(t *testing.T) {
			r := protocol.CredentialCreationOptions{}
			if err := json.Unmarshal([]byte(attestationRequests[i]), &r); err != nil {
				t.Fatal(err)
			}

			b := protocol.AttestationResponse{}
			if err := json.Unmarshal([]byte(attestationResponses[i]), &b); err != nil {
				t.Fatal(err)
			}

			p, err := protocol.ParseAttestationResponse(b)
			if err != nil {
				t.Fatal(err)
			}

			d, err := protocol.IsValidAttestation(p, r.PublicKey.Challenge, "", "")
			if err != nil {
				e := protocol.ToWebAuthnError(err)
				t.Fatal(fmt.Sprintf("%s, %s: %s", e.Name, e.Description, e.Debug))
			}

			if !d {
				t.Fatal("is not valid")
			}
		})
	}
}

func TestIsValidAssertion(t *testing.T) {
	for i := range assertionRequests {
		t.Run(fmt.Sprintf("Run %d", i), func(t *testing.T) {
			rawAttestation := protocol.AttestationResponse{}
			if err := json.Unmarshal([]byte(attestationResponses[i]), &rawAttestation); err != nil {
				t.Fatal(err)
			}

			attestation, err := protocol.ParseAttestationResponse(rawAttestation)
			if err != nil {
				t.Fatal(err)
			}

			publicKey := attestation.Response.Attestation.AuthData.AttestedCredentialData.COSEKey

			cert := &x509.Certificate{
				PublicKey: publicKey,
			}

			r := protocol.CredentialCreationOptions{}
			if err := json.Unmarshal([]byte(assertionRequests[i]), &r); err != nil {
				t.Fatal(err)
			}

			b := protocol.AssertionResponse{}
			if err := json.Unmarshal([]byte(assertionResponses[i]), &b); err != nil {
				t.Fatal(err)
			}

			p, err := protocol.ParseAssertionResponse(b)
			if err != nil {
				t.Fatal(err)
			}

			d, err := protocol.IsValidAssertion(p, r.PublicKey.Challenge, "", "", cert)
			if err != nil {
				t.Fatal(err)
			}

			if !d {
				t.Fatal("is not valid")
			}
		})
	}
}

var attestationRequests = []string{
	`{"publicKey":{"rp":{"name":"accountsvc"},"user":{"id":"MTAwNjg1ODU4NDE3ODI5NDc4NA==","name":"Koen Vlaswinkel","displayName":"Koen Vlaswinkel"},"pubKeyCredParams":[{"type":"public-key","alg":-7}],"timeout":10000,"attestation":"direct","challenge":"+1jQysnwaIjNU+GrwRp4PWNBMlX0i9/caRkcKd7LPj8="}}`,
	`{"publicKey":{"rp":{"name":"webauthn-demo"},"user":{"name":"koen","id":"a29lbg==","displayName":"koen"},"challenge":"JUtlYcgpkSiFNzsThDYuOrtSVY1VeLofM+mWTRCCXqU=","pubKeyCredParams":[{"type":"public-key","alg":-7}],"timeout":30000,"authenticatorSelection":{"requireResidentKey":false},"attestation":"direct"}}`,
	`{"publicKey":{"rp":{"name":"webauthn-demo"},"user":{"name":"koen","id":"a29lbg==","displayName":"koen"},"challenge":"2HzAlPIGskbn53hBJZeH3kZ6XfcHWMnzbATVG/FSgkI=","pubKeyCredParams":[{"type":"public-key","alg":-7}],"timeout":30000,"authenticatorSelection":{"requireResidentKey":false},"attestation":"direct"}}`,
	// Android SafetyNet
	`{"publicKey":{"rp":{"name":"webauthn-demo"},"user":{"name":"Bewus","id":"QmV3dXM=","displayName":"koen"},"challenge":"d3cY1I6n1ar6gLpDEhTi5nBgP1xwIGsb6HM/NR8PK1o=","pubKeyCredParams":[{"type":"public-key","alg":-7}],"timeout":30000,"authenticatorSelection":{"requireResidentKey":false},"attestation":"direct"}}`,
}

var attestationResponses = []string{
	`{"id":"LOXI3xfiLvIP04MD_S2ZmJYwn3cvMX1FUXxiQO7xlfUvrfcj99UVO2aMrMAwsGvsujY7NHWiM6G3B6ryKJDBBdab-cl4tVZeOwOMhgvHLXk","rawId":"LOXI3xfiLvIP04MD/S2ZmJYwn3cvMX1FUXxiQO7xlfUvrfcj99UVO2aMrMAwsGvsujY7NHWiM6G3B6ryKJDBBdab+cl4tVZeOwOMhgvHLXk=","response":{"attestationObject":"o2dhdHRTdG10omNzaWdYRjBEAiAJ8Q7i8DQzKlb00g4Wby4PoEjlI+s3bS+kVKI3PKoyXQIgDzcP2c5vpplZdmftN+zUDNfXtG1TniWbJv2+6kGZ8bljeDVjgVkBKzCCAScwgc6gAwIBAgIBADAKBggqhkjOPQQDAjAWMRQwEgYDVQQDDAtLcnlwdG9uIEtleTAeFw0xODA5MTcxODQ3NDJaFw0yODA5MTcxODQ3NDJaMBYxFDASBgNVBAMMC0tyeXB0b24gS2V5MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEwzIpvM5A6mZQXYxRIhfp0sb/21yTcr/sp5Y5DU0IWODQf5ldS2rlDCl62yEaQDM9Akxbsay/vA/S5ut4VSsvoKMNMAswCQYDVR0TBAIwADAKBggqhkjOPQQDAgNIADBFAiA4Yx+5MtKVnjme6V3qXKQ2qcgaHfO6DMgXM9kwOCZcNAIhAJdNk5PPSA04ITfrX9HQy5azo8sH9yhkW7c6gLdb/Kz+aGF1dGhEYXRhWNRJlg3liA6MaHQ0Fw9kdmBbj+SuuaKGMseZXPO6gx2XY0EAAAAALOXI3xfiLvIP04MD/S2ZmABQLOXI3xfiLvIP04MD/S2ZmJYwn3cvMX1FUXxiQO7xlfUvrfcj99UVO2aMrMAwsGvsujY7NHWiM6G3B6ryKJDBBdab+cl4tVZeOwOMhgvHLXmlAQIDJiABIVggwzIpvM5A6mZQXYxRIhfp0sb/21yTcr/sp5Y5DU0IWOAiWCDQf5ldS2rlDCl62yEaQDM9Akxbsay/vA/S5ut4VSsvoGNmbXRoZmlkby11MmY=","clientDataJSON":"eyJjaGFsbGVuZ2UiOiItMWpReXNud2FJak5VLUdyd1JwNFBXTkJNbFgwaTlfY2FSa2NLZDdMUGo4IiwiY2xpZW50RXh0ZW5zaW9ucyI6e30sImhhc2hBbGdvcml0aG0iOiJTSEEtMjU2Iiwib3JpZ2luIjoiaHR0cDovL2xvY2FsaG9zdDo1Mzg3OSIsInRva2VuQmluZGluZyI6eyJzdGF0dXMiOiJub3Qtc3VwcG9ydGVkIn0sInR5cGUiOiJ3ZWJhdXRobi5jcmVhdGUifQ=="},"type":"public-key"}`,
	`{"id":"SNBSJTt1DHEuG9XBd6lfc4XXqxkppWfFbt4P5sRVQEPIPANIHHCmPo1AwY5pkUGcpVL3W-uHyWEn4vbgzp34Qw","rawId":"SNBSJTt1DHEuG9XBd6lfc4XXqxkppWfFbt4P5sRVQEPIPANIHHCmPo1AwY5pkUGcpVL3W+uHyWEn4vbgzp34Qw==","response":{"attestationObject":"o2NmbXRmcGFja2VkZ2F0dFN0bXSjY2FsZyZjc2lnWEcwRQIgFls/elhmdZmqEBEKafdcyvQPDrTdBRMW92v6RKJj1bACIQCZ+46sXn65dMEpPuGxvMUruV5i7XN25ctFV/iAi3wSomN4NWOBWQLCMIICvjCCAaagAwIBAgIEdIb9wjANBgkqhkiG9w0BAQsFADAuMSwwKgYDVQQDEyNZdWJpY28gVTJGIFJvb3QgQ0EgU2VyaWFsIDQ1NzIwMDYzMTAgFw0xNDA4MDEwMDAwMDBaGA8yMDUwMDkwNDAwMDAwMFowbzELMAkGA1UEBhMCU0UxEjAQBgNVBAoMCVl1YmljbyBBQjEiMCAGA1UECwwZQXV0aGVudGljYXRvciBBdHRlc3RhdGlvbjEoMCYGA1UEAwwfWXViaWNvIFUyRiBFRSBTZXJpYWwgMTk1NTAwMzg0MjBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABJVd8633JH0xde/9nMTzGk6HjrrhgQlWYVD7OIsuX2Unv1dAmqWBpQ0KxS8YRFwKE1SKE1PIpOWacE5SO8BN6+2jbDBqMCIGCSsGAQQBgsQKAgQVMS4zLjYuMS40LjEuNDE0ODIuMS4xMBMGCysGAQQBguUcAgEBBAQDAgUgMCEGCysGAQQBguUcAQEEBBIEEPigEfOMCk0VgAYXER+e3H0wDAYDVR0TAQH/BAIwADANBgkqhkiG9w0BAQsFAAOCAQEAMVxIgOaaUn44Zom9af0KqG9J655OhUVBVW+q0As6AIod3AH5bHb2aDYakeIyyBCnnGMHTJtuekbrHbXYXERIn4aKdkPSKlyGLsA/A+WEi+OAfXrNVfjhrh7iE6xzq0sg4/vVJoywe4eAJx0fS+Dl3axzTTpYl71Nc7p/NX6iCMmdik0pAuYJegBcTckE3AoYEg4K99AM/JaaKIblsbFh8+3LxnemeNf7UwOczaGGvjS6UzGVI0Odf9lKcPIwYhuTxM5CaNMXTZQ7xq4/yTfC3kPWtE4hFT34UJJflZBiLrxG4OsYxkHw/n5vKgmpspB3GfYuYTWhkDKiE8CYtyg87mhhdXRoRGF0YVjESZYN5YgOjGh0NBcPZHZgW4/krrmihjLHmVzzuoMdl2NBAAAAA/igEfOMCk0VgAYXER+e3H0AQEjQUiU7dQxxLhvVwXepX3OF16sZKaVnxW7eD+bEVUBDyDwDSBxwpj6NQMGOaZFBnKVS91vrh8lhJ+L24M6d+EOlAQIDJiABIVggLxxTguKmjCV4N5OMqd2Sl9AIxSltaPevmQxSqnyNlAciWCDEHOaQDaZ6pC2gC+Z0KS4Ln/XQiJp0X1BmTd+K+FdqSg==","clientDataJSON":"eyJjaGFsbGVuZ2UiOiJKVXRsWWNncGtTaUZOenNUaERZdU9ydFNWWTFWZUxvZk0tbVdUUkNDWHFVIiwibmV3X2tleXNfbWF5X2JlX2FkZGVkX2hlcmUiOiJkbyBub3QgY29tcGFyZSBjbGllbnREYXRhSlNPTiBhZ2FpbnN0IGEgdGVtcGxhdGUuIFNlZSBodHRwczovL2dvby5nbC95YWJQZXgiLCJvcmlnaW4iOiJodHRwOi8vbG9jYWxob3N0OjkwMDAiLCJ0eXBlIjoid2ViYXV0aG4uY3JlYXRlIn0="},"type":"public-key"}`,
	`{"id":"EBT1LOefp-8ID0n2jchlyaPrKcWZ6jdHH8nb0Z-hi9JHsOpTpCNUbJ7ijJOKdetLOy2cqdxNq8zkWYmCgpapKg","rawId":"EBT1LOefp+8ID0n2jchlyaPrKcWZ6jdHH8nb0Z+hi9JHsOpTpCNUbJ7ijJOKdetLOy2cqdxNq8zkWYmCgpapKg==","response":{"attestationObject":"o2NmbXRoZmlkby11MmZnYXR0U3RtdKJjc2lnWEgwRgIhAJkpVpWsMm/Z1OnF/+B/juq/IAlKqhakms5HkNf6ZKLWAiEAm2qNX/bHUkkdaJ0seanz5xxVDCn+bKGEPyQP3ZpPczNjeDVjgVkCUzCCAk8wggE3oAMCAQICBA0ACxYwDQYJKoZIhvcNAQELBQAwLjEsMCoGA1UEAxMjWXViaWNvIFUyRiBSb290IENBIFNlcmlhbCA0NTcyMDA2MzEwIBcNMTQwODAxMDAwMDAwWhgPMjA1MDA5MDQwMDAwMDBaMDExLzAtBgNVBAMMJll1YmljbyBVMkYgRUUgU2VyaWFsIDIzOTI1NzM0MDE1NzY1MjcwMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAETKz6btEEuhlL1uBm1+E/zGpgDxDSSFx+o9vUTNDVDbJROHujvR665t7mJQoFWMbpvmEYpEOOWkNfHtLrDOi7haM7MDkwIgYJKwYBBAGCxAoCBBUxLjMuNi4xLjQuMS40MTQ4Mi4xLjUwEwYLKwYBBAGC5RwCAQEEBAMCBSAwDQYJKoZIhvcNAQELBQADggEBAI7CTaiBlYLnMIQZnJ8UCvrqgFuin80CTT4UAiGWsBwh0eY+CRSwL4LEFZITkLlFYyOsfMDlI7oddSN/Jmn8HzrPWvzKVP/+mCuRMSdz735wFNYX5xle+NLkoctZjyHOCqdd4B8lgX0nzwNiPZuf+sdY5fhzhLRmtbpfBDToTP57tLR5WlIY6kJ6QKecpZ5sVNxCzSVxRncAptZV7YSsX2we05Kt5mHkBHqhi5CTPQQmOObHov7cB+4q5CpufDzEBFTKPL3tWxV6HvQr0J6Mp6bZFICq5nTP7VPatnnJelRA9VmPSpQuLjpRqpJFKRobj8eQ9yuveXG/7uutBOzBHW9oYXV0aERhdGFYxEmWDeWIDoxodDQXD2R2YFuP5K65ooYyx5lc87qDHZdjQQAAAAAAAAAAAAAAAAAAAAAAAAAAAEAQFPUs55+n7wgPSfaNyGXJo+spxZnqN0cfydvRn6GL0kew6lOkI1RsnuKMk4p160s7LZyp3E2rzORZiYKClqkqpQECAyYgASFYIF6oiA6H+mU150XH7WJ2vnzNmdzgr5YloPao7ePjNjlOIlggg0f3u4CtxsBkkKjo7v4luyJui9tJ1rGTBF3YkYlcADo=","clientDataJSON":"eyJjaGFsbGVuZ2UiOiIySHpBbFBJR3NrYm41M2hCSlplSDNrWjZYZmNIV01uemJBVFZHX0ZTZ2tJIiwib3JpZ2luIjoiaHR0cDovL2xvY2FsaG9zdDo5MDAwIiwidHlwZSI6IndlYmF1dGhuLmNyZWF0ZSJ9"},"type":"public-key"}`,
	// Android SafetyNet
	`{"id":"ARKBRFD84uLN6qG_rHsV0K2Bh9Lj3_HaJsXdC_DpPslKO6ZWmD38-hz90Lf_MzELErMa9AqR21Sr9brNzE2un1U","rawId":"ARKBRFD84uLN6qG/rHsV0K2Bh9Lj3/HaJsXdC/DpPslKO6ZWmD38+hz90Lf/MzELErMa9AqR21Sr9brNzE2un1U=","response":{"attestationObject":"o2NmbXRxYW5kcm9pZC1zYWZldHluZXRnYXR0U3RtdKJjdmVyaDE0MzY2MDE5aHJlc3BvbnNlWRS9ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbmcxWXlJNld5Sk5TVWxHYTJwRFEwSkljV2RCZDBsQ1FXZEpVVkpZY205T01GcFBaRkpyUWtGQlFVRkJRVkIxYm5wQlRrSm5hM0ZvYTJsSE9YY3dRa0ZSYzBaQlJFSkRUVkZ6ZDBOUldVUldVVkZIUlhkS1ZsVjZSV1ZOUW5kSFFURlZSVU5vVFZaU01qbDJXako0YkVsR1VubGtXRTR3U1VaT2JHTnVXbkJaTWxaNlRWSk5kMFZSV1VSV1VWRkVSWGR3U0ZaR1RXZFJNRVZuVFZVNGVFMUNORmhFVkVVMFRWUkJlRTFFUVROTlZHc3dUbFp2V0VSVVJUVk5WRUYzVDFSQk0wMVVhekJPVm05M1lrUkZURTFCYTBkQk1WVkZRbWhOUTFaV1RYaEZla0ZTUW1kT1ZrSkJaMVJEYTA1b1lrZHNiV0l6U25WaFYwVjRSbXBCVlVKblRsWkNRV05VUkZVeGRtUlhOVEJaVjJ4MVNVWmFjRnBZWTNoRmVrRlNRbWRPVmtKQmIxUkRhMlIyWWpKa2MxcFRRazFVUlUxNFIzcEJXa0puVGxaQ1FVMVVSVzFHTUdSSFZucGtRelZvWW0xU2VXSXliR3RNYlU1MllsUkRRMEZUU1hkRVVWbEtTMjlhU1doMlkwNUJVVVZDUWxGQlJHZG5SVkJCUkVORFFWRnZRMmRuUlVKQlRtcFlhM293WlVzeFUwVTBiU3N2UnpWM1QyOHJXRWRUUlVOeWNXUnVPRGh6UTNCU04yWnpNVFJtU3pCU2FETmFRMWxhVEVaSWNVSnJOa0Z0V2xaM01rczVSa2N3VHpseVVsQmxVVVJKVmxKNVJUTXdVWFZ1VXpsMVowaEROR1ZuT1c5MmRrOXRLMUZrV2pKd09UTllhSHAxYmxGRmFGVlhXRU40UVVSSlJVZEtTek5UTW1GQlpucGxPVGxRVEZNeU9XaE1ZMUYxV1ZoSVJHRkROMDlhY1U1dWIzTnBUMGRwWm5NNGRqRnFhVFpJTDNob2JIUkRXbVV5YkVvck4wZDFkSHBsZUV0d2VIWndSUzkwV2xObVlsazVNRFZ4VTJ4Q2FEbG1jR293TVRWamFtNVJSbXRWYzBGVmQyMUxWa0ZWZFdWVmVqUjBTMk5HU3pSd1pYWk9UR0Y0UlVGc0swOXJhV3hOZEVsWlJHRmpSRFZ1Wld3MGVFcHBlWE0wTVROb1lXZHhWekJYYUdnMVJsQXpPV2hIYXpsRkwwSjNVVlJxWVhwVGVFZGtkbGd3YlRaNFJsbG9hQzh5VmsxNVdtcFVORXQ2VUVwRlEwRjNSVUZCWVU5RFFXeG5kMmRuU2xWTlFUUkhRVEZWWkVSM1JVSXZkMUZGUVhkSlJtOUVRVlJDWjA1V1NGTlZSVVJFUVV0Q1oyZHlRbWRGUmtKUlkwUkJWRUZOUW1kT1ZraFNUVUpCWmpoRlFXcEJRVTFDTUVkQk1WVmtSR2RSVjBKQ1VYRkNVWGRIVjI5S1FtRXhiMVJMY1hWd2J6UlhObmhVTm1veVJFRm1RbWRPVmtoVFRVVkhSRUZYWjBKVFdUQm1hSFZGVDNaUWJTdDRaMjU0YVZGSE5rUnlabEZ1T1V0NlFtdENaMmR5UW1kRlJrSlJZMEpCVVZKWlRVWlpkMHAzV1VsTGQxbENRbEZWU0UxQlIwZEhNbWd3WkVoQk5reDVPWFpaTTA1M1RHNUNjbUZUTlc1aU1qbHVUREprTUdONlJuWk5WRUZ5UW1kbmNrSm5SVVpDVVdOM1FXOVpabUZJVWpCalJHOTJURE5DY21GVE5XNWlNamx1VERKa2VtTnFTWFpTTVZKVVRWVTRlRXh0VG5sa1JFRmtRbWRPVmtoU1JVVkdha0ZWWjJoS2FHUklVbXhqTTFGMVdWYzFhMk50T1hCYVF6VnFZakl3ZDBsUldVUldVakJuUWtKdmQwZEVRVWxDWjFwdVoxRjNRa0ZuU1hkRVFWbExTM2RaUWtKQlNGZGxVVWxHUVhwQmRrSm5UbFpJVWpoRlMwUkJiVTFEVTJkSmNVRm5hR2cxYjJSSVVuZFBhVGgyV1ROS2MweHVRbkpoVXpWdVlqSTVia3d3WkZWVmVrWlFUVk0xYW1OdGQzZG5aMFZGUW1kdmNrSm5SVVZCWkZvMVFXZFJRMEpKU0RGQ1NVaDVRVkJCUVdSM1EydDFVVzFSZEVKb1dVWkpaVGRGTmt4TldqTkJTMUJFVjFsQ1VHdGlNemRxYW1RNE1FOTVRVE5qUlVGQlFVRlhXbVJFTTFCTVFVRkJSVUYzUWtsTlJWbERTVkZEVTFwRFYyVk1Tblp6YVZaWE5rTm5LMmRxTHpsM1dWUktVbnAxTkVocGNXVTBaVmswWXk5dGVYcHFaMGxvUVV4VFlta3ZWR2g2WTNweGRHbHFNMlJyTTNaaVRHTkpWek5NYkRKQ01HODNOVWRSWkdoTmFXZGlRbWRCU0ZWQlZtaFJSMjFwTDFoM2RYcFVPV1ZIT1ZKTVNTdDRNRm95ZFdKNVdrVldla0UzTlZOWlZtUmhTakJPTUVGQlFVWnRXRkU1ZWpWQlFVRkNRVTFCVW1wQ1JVRnBRbU5EZDBFNWFqZE9WRWRZVURJM09IbzBhSEl2ZFVOSWFVRkdUSGx2UTNFeVN6QXJlVXhTZDBwVlltZEpaMlk0WjBocWRuQjNNbTFDTVVWVGFuRXlUMll6UVRCQlJVRjNRMnR1UTJGRlMwWlZlVm8zWmk5UmRFbDNSRkZaU2t0dldrbG9kbU5PUVZGRlRFSlJRVVJuWjBWQ1FVazVibFJtVWt0SlYyZDBiRmRzTTNkQ1REVTFSVlJXTm10aGVuTndhRmN4ZVVGak5VUjFiVFpZVHpReGExcDZkMG8yTVhkS2JXUlNVbFF2VlhORFNYa3hTMFYwTW1Nd1JXcG5iRzVLUTBZeVpXRjNZMFZYYkV4UldUSllVRXg1Um1wclYxRk9ZbE5vUWpGcE5GY3lUbEpIZWxCb2RETnRNV0kwT1doaWMzUjFXRTAyZEZnMVEzbEZTRzVVYURoQ2IyMDBMMWRzUm1sb2VtaG5iamd4Ukd4a2IyZDZMMHN5VlhkTk5sTTJRMEl2VTBWNGEybFdabllyZW1KS01ISnFkbWM1TkVGc1pHcFZabFYzYTBrNVZrNU5ha1ZRTldVNGVXUkNNMjlNYkRabmJIQkRaVVkxWkdkbVUxZzBWVGw0TXpWdmFpOUpTV1F6VlVVdlpGQndZaTl4WjBkMmMydG1aR1Y2ZEcxVmRHVXZTMU50Y21sM1kyZFZWMWRsV0daVVlra3plbk5wYTNkYVltdHdiVkpaUzIxcVVHMW9kalJ5YkdsNlIwTkhkRGhRYmpod2NUaE5Na3RFWmk5UU0ydFdiM1F6WlRFNFVUMGlMQ0pOU1VsRlUycERRMEY2UzJkQmQwbENRV2RKVGtGbFR6QnRjVWRPYVhGdFFrcFhiRkYxUkVGT1FtZHJjV2hyYVVjNWR6QkNRVkZ6UmtGRVFrMU5VMEYzU0dkWlJGWlJVVXhGZUdSSVlrYzVhVmxYZUZSaFYyUjFTVVpLZG1JelVXZFJNRVZuVEZOQ1UwMXFSVlJOUWtWSFFURlZSVU5vVFV0U01uaDJXVzFHYzFVeWJHNWlha1ZVVFVKRlIwRXhWVVZCZUUxTFVqSjRkbGx0Um5OVk1teHVZbXBCWlVaM01IaE9la0V5VFZSVmQwMUVRWGRPUkVwaFJuY3dlVTFVUlhsTlZGVjNUVVJCZDA1RVNtRk5SVWw0UTNwQlNrSm5UbFpDUVZsVVFXeFdWRTFTTkhkSVFWbEVWbEZSUzBWNFZraGlNamx1WWtkVloxWklTakZqTTFGblZUSldlV1J0YkdwYVdFMTRSWHBCVWtKblRsWkNRVTFVUTJ0a1ZWVjVRa1JSVTBGNFZIcEZkMmRuUldsTlFUQkhRMU54UjFOSllqTkVVVVZDUVZGVlFVRTBTVUpFZDBGM1oyZEZTMEZ2U1VKQlVVUlJSMDA1UmpGSmRrNHdOWHByVVU4NUszUk9NWEJKVW5aS2VucDVUMVJJVnpWRWVrVmFhRVF5WlZCRGJuWlZRVEJSYXpJNFJtZEpRMlpMY1VNNVJXdHpRelJVTW1aWFFsbHJMMnBEWmtNelVqTldXazFrVXk5a1RqUmFTME5GVUZwU2NrRjZSSE5wUzFWRWVsSnliVUpDU2pWM2RXUm5lbTVrU1UxWlkweGxMMUpIUjBac05YbFBSRWxMWjJwRmRpOVRTa2d2VlV3clpFVmhiSFJPTVRGQ2JYTkxLMlZSYlUxR0t5dEJZM2hIVG1oeU5UbHhUUzg1YVd3M01Va3laRTQ0UmtkbVkyUmtkM1ZoWldvMFlsaG9jREJNWTFGQ1ltcDRUV05KTjBwUU1HRk5NMVEwU1N0RWMyRjRiVXRHYzJKcWVtRlVUa001ZFhwd1JteG5UMGxuTjNKU01qVjRiM2x1VlhoMk9IWk9iV3R4TjNwa1VFZElXR3Q0VjFrM2IwYzVhaXRLYTFKNVFrRkNhemRZY2twbWIzVmpRbHBGY1VaS1NsTlFhemRZUVRCTVMxY3dXVE42Tlc5Nk1rUXdZekYwU2t0M1NFRm5UVUpCUVVkcVoyZEZlazFKU1VKTWVrRlBRbWRPVmtoUk9FSkJaamhGUWtGTlEwRlpXWGRJVVZsRVZsSXdiRUpDV1hkR1FWbEpTM2RaUWtKUlZVaEJkMFZIUTBOelIwRlJWVVpDZDAxRFRVSkpSMEV4VldSRmQwVkNMM2RSU1UxQldVSkJaamhEUVZGQmQwaFJXVVJXVWpCUFFrSlpSVVpLYWxJclJ6UlJOamdyWWpkSFEyWkhTa0ZpYjA5ME9VTm1NSEpOUWpoSFFURlZaRWwzVVZsTlFtRkJSa3AyYVVJeFpHNUlRamRCWVdkaVpWZGlVMkZNWkM5alIxbFpkVTFFVlVkRFEzTkhRVkZWUmtKM1JVSkNRMnQzU25wQmJFSm5aM0pDWjBWR1FsRmpkMEZaV1ZwaFNGSXdZMFJ2ZGt3eU9XcGpNMEYxWTBkMGNFeHRaSFppTW1OMldqTk9lVTFxUVhsQ1owNVdTRkk0UlV0NlFYQk5RMlZuU21GQmFtaHBSbTlrU0ZKM1QyazRkbGt6U25OTWJrSnlZVk0xYm1JeU9XNU1NbVI2WTJwSmRsb3pUbmxOYVRWcVkyMTNkMUIzV1VSV1VqQm5Ra1JuZDA1cVFUQkNaMXB1WjFGM1FrRm5TWGRMYWtGdlFtZG5ja0puUlVaQ1VXTkRRVkpaWTJGSVVqQmpTRTAyVEhrNWQyRXlhM1ZhTWpsMlduazVlVnBZUW5aak1td3dZak5LTlV4NlFVNUNaMnR4YUd0cFJ6bDNNRUpCVVhOR1FVRlBRMEZSUlVGSGIwRXJUbTV1TnpoNU5uQlNhbVE1V0d4UlYwNWhOMGhVWjJsYUwzSXpVazVIYTIxVmJWbElVRkZ4TmxOamRHazVVRVZoYW5aM1VsUXlhVmRVU0ZGeU1ESm1aWE54VDNGQ1dUSkZWRlYzWjFwUksyeHNkRzlPUm5ab2MwODVkSFpDUTA5SllYcHdjM2RYUXpsaFNqbDRhblUwZEZkRVVVZzRUbFpWTmxsYVdpOVlkR1ZFVTBkVk9WbDZTbkZRYWxrNGNUTk5SSGh5ZW0xeFpYQkNRMlkxYnpodGR5OTNTalJoTWtjMmVIcFZjalpHWWpaVU9FMWpSRTh5TWxCTVVrdzJkVE5OTkZSNmN6TkJNazB4YWpaaWVXdEtXV2s0ZDFkSlVtUkJka3RNVjFwMUwyRjRRbFppZWxsdGNXMTNhMjAxZWt4VFJGYzFia2xCU21KRlRFTlJRMXAzVFVnMU5uUXlSSFp4YjJaNGN6WkNRbU5EUmtsYVZWTndlSFUyZURaMFpEQldOMU4yU2tORGIzTnBjbE50U1dGMGFpODVaRk5UVmtSUmFXSmxkRGh4THpkVlN6UjJORnBWVGpnd1lYUnVXbm94ZVdjOVBTSmRmUS5leUp1YjI1alpTSTZJbEJQT1U0MWVrY3pZbTlOUms1Nk5UbHFWRU01ZG1vdlp6QkxlalExTjJoRVMyOTRTREZzZURSbVlWazlJaXdpZEdsdFpYTjBZVzF3VFhNaU9qRTFOREEwTURZeU16RTBOellzSW1Gd2ExQmhZMnRoWjJWT1lXMWxJam9pWTI5dExtZHZiMmRzWlM1aGJtUnliMmxrTG1kdGN5SXNJbUZ3YTBScFoyVnpkRk5vWVRJMU5pSTZJbVZSWXl0MmVsVmpaSGd3UmxaT1RIWllTSFZIY0VRd0sxSTRNRGR6VlVWMmNDdEtaV3hsV1ZwemFVRTlJaXdpWTNSelVISnZabWxzWlUxaGRHTm9JanAwY25WbExDSmhjR3REWlhKMGFXWnBZMkYwWlVScFoyVnpkRk5vWVRJMU5pSTZXeUk0VURGelZ6QkZVRXBqYzJ4M04xVjZVbk5wV0V3Mk5IY3JUelV3UldRclVrSkpRM1JoZVRGbk1qUk5QU0pkTENKaVlYTnBZMGx1ZEdWbmNtbDBlU0k2ZEhKMVpYMC5qaFE5a3RMLVVCdEJpX2NQNXd4SGRFejJZWXNGMXBLWXpqOEZUZXdXbVVvWE5CTWJSOWRBaEdDSnJCZ2w4RGZNSzFrMUJFQXdQUzRTMWJVczBXQ3haYmN4cWJtcS12UF9OWHI1QjlDQXkxUUpxdC1tRlRMUm5EZkZ1a2hfQjdZMUxEZUJaaGYtc1E4WnBfQUlncHRnYlBWa2Z6TE1PQVpkeE5xVk91dmU0YmJTSjQ5bWQwVklBbDkwc3h1YXVUT0x5bFpxN2ZhYXRZMGFqQ1VKZkNpRUJiLVZxSUhOZmQySExhaUZwcjVxRUFPeU8tQmx1V204TmVfSWZiMnFkTVZBa2p1a1YyVmZheElLbG93a05HZ2ZycjFjVVpJNE5oVTNfeDhLTjRGclpQd29tT2ZFdmlHWk9tZFRvNnNPSXpROTJVYkx2MXlGYUw1TDZIZ1I2Z1NnX0FoYXV0aERhdGFYxSpD77HzPafHbmULkXmwl2mw9P/lQRkLyNgWFO1qQIjVRQAAAAAAAAAAAAAAAAAAAAAAAAAAAEEBEoFEUPzi4s3qob+sexXQrYGH0uPf8domxd0L8Ok+yUo7plaYPfz6HP3Qt/8zMQsSsxr0CpHbVKv1us3MTa6fVaUBAgMmIAEhWCDavINo/+JM4T1eJeKjaZ+vGa2Do7YVh2EyD0vtmoZrrCJYIOtAindNbNXogAQxBAJii2Vd1Wl5rZb9KPak8J6iTKle","clientDataJSON":"eyJ0eXBlIjoid2ViYXV0aG4uY3JlYXRlIiwiY2hhbGxlbmdlIjoiZDNjWTFJNm4xYXI2Z0xwREVoVGk1bkJnUDF4d0lHc2I2SE1fTlI4UEsxbyIsIm9yaWdpbiI6Imh0dHBzOlwvXC9iMzk5ZmEwMC5uZ3Jvay5pbyIsImFuZHJvaWRQYWNrYWdlTmFtZSI6ImNvbS5hbmRyb2lkLmNocm9tZSJ9"},"type":"public-key"}`,
}

var assertionRequests = []string{
	`{"publicKey":{"allowCredentials":[{"id":"LOXI3xfiLvIP04MD/S2ZmJYwn3cvMX1FUXxiQO7xlfUvrfcj99UVO2aMrMAwsGvsujY7NHWiM6G3B6ryKJDBBdab+cl4tVZeOwOMhgvHLXk=","type":"public-key"}],"challenge":"+c0hMsULvTWp6ASl45YyOQRA/yVVK60XccCQ+Vui9j8=","timeout":10000}}`,
	`{"publicKey":{"challenge":"mcPXIDRHSPBF2gJWU58GPrR3TodLDXR1kHJhgVanYnU=","timeout":30000,"allowCredentials":[{"type":"public-key","id":"SNBSJTt1DHEuG9XBd6lfc4XXqxkppWfFbt4P5sRVQEPIPANIHHCmPo1AwY5pkUGcpVL3W+uHyWEn4vbgzp34Qw=="}]}}`,
	`{"publicKey":{"challenge":"/hXFS7WKYWTgqEx5AOG7SuGL3+6alkqi2TJkTu+MkBM=","timeout":30000,"allowCredentials":[{"type":"public-key","id":"EBT1LOefp+8ID0n2jchlyaPrKcWZ6jdHH8nb0Z+hi9JHsOpTpCNUbJ7ijJOKdetLOy2cqdxNq8zkWYmCgpapKg=="}]}}`,
	// Android SafetyNet
	`{"publicKey":{"challenge":"HC33hV7jFYx6m4hUkvNF0GLVn2WihTilaEtniha+Qvw=","timeout":30000,"allowCredentials":[{"type":"public-key","id":"ARKBRFD84uLN6qG/rHsV0K2Bh9Lj3/HaJsXdC/DpPslKO6ZWmD38+hz90Lf/MzELErMa9AqR21Sr9brNzE2un1U="}]}}`,
}

var assertionResponses = []string{
	`{"id":"LOXI3xfiLvIP04MD_S2ZmJYwn3cvMX1FUXxiQO7xlfUvrfcj99UVO2aMrMAwsGvsujY7NHWiM6G3B6ryKJDBBdab-cl4tVZeOwOMhgvHLXk","rawId":"LOXI3xfiLvIP04MD/S2ZmJYwn3cvMX1FUXxiQO7xlfUvrfcj99UVO2aMrMAwsGvsujY7NHWiM6G3B6ryKJDBBdab+cl4tVZeOwOMhgvHLXk=","response":{"clientDataJSON":"eyJjaGFsbGVuZ2UiOiItYzBoTXNVTHZUV3A2QVNsNDVZeU9RUkFfeVZWSzYwWGNjQ1EtVnVpOWo4IiwiaGFzaEFsZ29yaXRobSI6IlNIQS0yNTYiLCJvcmlnaW4iOiJodHRwOi8vbG9jYWxob3N0OjUzODc5IiwidHlwZSI6IndlYmF1dGhuLmdldCJ9","authenticatorData":"SZYN5YgOjGh0NBcPZHZgW4/krrmihjLHmVzzuoMdl2MBAAAAAQ==","signature":"MEYCIQD7W6TPIviP+BztYxEMsan/esy/O0S4pJO+9QxDaA0ehAIhANo5D+5UxwbtJGFcvSryl0+RdJd3j4lIKVhEe7WpvZeV","userHandle":""},"type":"public-key"}`,
	`{"id":"SNBSJTt1DHEuG9XBd6lfc4XXqxkppWfFbt4P5sRVQEPIPANIHHCmPo1AwY5pkUGcpVL3W-uHyWEn4vbgzp34Qw","rawId":"SNBSJTt1DHEuG9XBd6lfc4XXqxkppWfFbt4P5sRVQEPIPANIHHCmPo1AwY5pkUGcpVL3W+uHyWEn4vbgzp34Qw==","response":{"clientDataJSON":"eyJjaGFsbGVuZ2UiOiJtY1BYSURSSFNQQkYyZ0pXVTU4R1ByUjNUb2RMRFhSMWtISmhnVmFuWW5VIiwib3JpZ2luIjoiaHR0cDovL2xvY2FsaG9zdDo5MDAwIiwidHlwZSI6IndlYmF1dGhuLmdldCJ9","authenticatorData":"SZYN5YgOjGh0NBcPZHZgW4/krrmihjLHmVzzuoMdl2MBAAAABA==","signature":"MEUCIQCWGnyWIV4s13/9TRcLtDesxa0UJs+pwNaF3YDP/5RHDwIgIWlEiH74R7sPiyNffp8Tof3qo1s8jVvFDxCGejlICFI=","userHandle":""},"type":"public-key"}`,
	`{"id":"EBT1LOefp-8ID0n2jchlyaPrKcWZ6jdHH8nb0Z-hi9JHsOpTpCNUbJ7ijJOKdetLOy2cqdxNq8zkWYmCgpapKg","rawId":"EBT1LOefp+8ID0n2jchlyaPrKcWZ6jdHH8nb0Z+hi9JHsOpTpCNUbJ7ijJOKdetLOy2cqdxNq8zkWYmCgpapKg==","response":{"clientDataJSON":"eyJjaGFsbGVuZ2UiOiJfaFhGUzdXS1lXVGdxRXg1QU9HN1N1R0wzLTZhbGtxaTJUSmtUdS1Na0JNIiwib3JpZ2luIjoiaHR0cDovL2xvY2FsaG9zdDo5MDAwIiwidHlwZSI6IndlYmF1dGhuLmdldCJ9","authenticatorData":"SZYN5YgOjGh0NBcPZHZgW4/krrmihjLHmVzzuoMdl2MBAAAAAA==","signature":"MEUCIFGAxD82g/HBEQc2qblhIQsOCvMIuFzmiT54uMSCwYg6AiEAuuIUy6PyaW43xEpAnqrPcCPmUiJwpJ7IV/h6OGjqN2E=","userHandle":""},"type":"public-key"}`,
	// Android SafetyNet
	`{"id":"ARKBRFD84uLN6qG_rHsV0K2Bh9Lj3_HaJsXdC_DpPslKO6ZWmD38-hz90Lf_MzELErMa9AqR21Sr9brNzE2un1U","rawId":"ARKBRFD84uLN6qG/rHsV0K2Bh9Lj3/HaJsXdC/DpPslKO6ZWmD38+hz90Lf/MzELErMa9AqR21Sr9brNzE2un1U=","response":{"clientDataJSON":"eyJ0eXBlIjoid2ViYXV0aG4uZ2V0IiwiY2hhbGxlbmdlIjoiSEMzM2hWN2pGWXg2bTRoVWt2TkYwR0xWbjJXaWhUaWxhRXRuaWhhLVF2dyIsIm9yaWdpbiI6Imh0dHBzOlwvXC9iMzk5ZmEwMC5uZ3Jvay5pbyIsImFuZHJvaWRQYWNrYWdlTmFtZSI6ImNvbS5hbmRyb2lkLmNocm9tZSJ9","authenticatorData":"KkPvsfM9p8duZQuRebCXabD0/+VBGQvI2BYU7WpAiNUFAAAAAQ==","signature":"MEQCIBapcKD8L5Kp92QBr4XpHNwiRPjo/MGTEIEwCsklxfvAAiABn02rbcatTqHFtHwbnHwdNOLa5apxCBRuFPPwABBm3w==","userHandle":""},"type":"public-key"}`,
}
