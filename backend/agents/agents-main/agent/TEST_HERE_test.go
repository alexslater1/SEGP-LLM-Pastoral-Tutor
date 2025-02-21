package agent

import (
	"testing"
)

func Test_NewAcadmemicSupportStudyAgent(t *testing.T) {
	agent := NewAcadmemicSupportStudyAgent("test")
}

func Test_NewAdminUniServicesAgent(t *testing.T) {
	agent := NewAdminUniServicesAgent("test")
}

func Test_CareerProfessionalDevelopmentAgent(t *testing.T) {
	agent := NewCareerProfessionalDevelopmentAgent("test")
}

func Test_WellbeingMentalHealthPersonalDevelopmentAgent(t *testing.T) {
	agent := NewWellbeingMentalHealthPersonalDevelopmentAgent("test")
}

func Test_FinancialAccomodationResourceAgent(t *testing.T) {
	agent := NewFinancialAccomodationResourceAgent("test")
}

func Test_CampusLifeSocialAgent(t *testing.T) {
	agent := NewCampusLifeSocialAgent("test")
}

func Test_AccessibilityDisabilityAgent(t *testing.T) {
	agent := NewAccessibilityDisabilityAgent("test")
}

func Test_TransitionDiversityMiscAgent(t *testing.T) {
	agent := NewTransitionDiversityMiscAgent("test")
}
