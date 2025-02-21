package agent

import (
	"testing"

	abc_api "github.com/segp/agents-main/abc-api"
)

var (
	abc_api_client = abc_api.NewMockAbcApiClient()
)

func Test_NewAcadmemicSupportStudyAgent(t *testing.T) {
	// agent := NewAcadmemicSupportStudyAgent("test")
}

func Test_NewAdminUniServicesAgent(t *testing.T) {
	// agent := NewAdminUniServicesAgent("test")
}

func Test_CareerProfessionalDevelopmentAgent(t *testing.T) {
	// agent := NewCareerProfessionalDevelopmentAgent("test")
}

func Test_WellbeingMentalHealthPersonalDevelopmentAgent(t *testing.T) {
	// agent := NewWellbeingMentalHealthPersonalDevelopmentAgent("test")
}

// ALEX
func Test_FinancialAccomodationResourceAgent(t *testing.T) {
	var (
		prompt = "TODO"
	)

	agent := NewFinancialAccomodationResourceAgent("test")
}

// ALEX
func Test_CampusLifeSocialAgent(t *testing.T) {
	agent := NewCampusLifeSocialAgent("test")
}

// TEO
func Test_AccessibilityDisabilityAgent(t *testing.T) {
	agent := NewAccessibilityDisabilityAgent("test")
}

// TEO
func Test_TransitionDiversityMiscAgent(t *testing.T) {
	var (
		prompt = "TODO"
		apis   = []interface{}{
			*abc_api_client.GetYears(),
			*abc_api_client.GetCohorts(),
		}
	)

	agent := NewTransitionDiversityMiscAgent(prompt, abc_api_client)
}
