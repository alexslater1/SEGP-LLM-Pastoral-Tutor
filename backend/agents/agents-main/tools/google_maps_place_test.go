package tools

import (
	"fmt"
	"testing"

	googleSearch "github.com/segp/agents-main/google_search"
)

func TestVerifyIsGoogleMapsPlaceURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "Valid Google Maps place URL",
			url:     "https://google.com/maps/place/Boucherie+Union+Square",
			wantErr: false,
		},
		{
			name:    "Valid Google Maps place URL with www",
			url:     "https://www.google.com/maps/place/Boucherie+Union+Square",
			wantErr: false,
		},
		{
			name:    "Empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "Invalid host",
			url:     "https://example.com/maps/place/Something",
			wantErr: true,
		},
		{
			name:    "Missing place in path",
			url:     "https://google.com/maps/something",
			wantErr: true,
		},
		{
			name:    "Invalid URL format",
			url:     "not-a-url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyIsGoogleMapsPlaceURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("verifyIsGoogleMapsPlaceURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParsePlaceDetails(t *testing.T) {
	tool := NewGoogleMapsPlaceTool(googleSearch.NewNonHeadlessRodClient())

	placeDetails, err := tool.PlaceDetailsFor("https://www.google.com/maps/place/Boucherie+Union+Square/data=!4m7!3m6!1s0x89c259a1ec5f5573:0x2fc6687f46f682d5!8m2!3d40.7372552!4d-73.9882246!16s%2Fg%2F11hbv5rh0_!19sChIJc1Vf7KFZwokR1YL2Rn9oxi8?authuser=0&hl=en&rclk=1")
	if err != nil {
		t.Errorf("PlaceDetailsFor() error =  %v", err)
	}

	fmt.Println(*placeDetails)
}

func TestPlaceAboutsFor(t *testing.T) {
	tool := NewGoogleMapsPlaceTool(googleSearch.NewNonHeadlessRodClient())

	placeAbout, err := tool.placeAboutFor("https://www.google.com/maps/place/Boucherie+Union+Square/data=!4m7!3m6!1s0x89c259a1ec5f5573:0x2fc6687f46f682d5!8m2!3d40.7372552!4d-73.9882246!16s%2Fg%2F11hbv5rh0_!19sChIJc1Vf7KFZwokR1YL2Rn9oxi8?authuser=0&hl=en&rclk=1")
	if err != nil {
		t.Errorf("PlaceDetailsFor() error = %v", err)
	}

	fmt.Printf("%+v\n", *placeAbout)
}
