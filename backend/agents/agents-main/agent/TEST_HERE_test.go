package agent

import (
	"context"
	"log"
	"testing"

	abc_api "github.com/segp/agents-main/abc-api"
)

var (
	abc_api_client = abc_api.NewMockAbcApiClient()
)

func Test_NewAcadmemicSupportStudyAgent(t *testing.T) {
	var (
		prompt = `EDIT THIS PROMPT`
		apis   = []interface{}{
			abc_api_client.GetModules(),
			abc_api_client.GetModule(),
			abc_api_client.GetEnrolledStudents(),
			abc_api_client.GetExams(),
			abc_api_client.GetPublicModuleTypes(),
			abc_api_client.GetDegreeRegulations(),
			abc_api_client.GetStudentDetails(),
			abc_api_client.GetStaff(),
			abc_api_client.GetIdentity(),
			abc_api_client.GetAcademicPeriods(),
			abc_api_client.GetCohorts(),
		}

		query = `EDIT THIS QUERY`
	)

	agent := NewAcadmemicSupportStudyAgent(prompt, apis...)
	runAgentAndPrintResultFrom(agent, query)
}

func Test_NewAdminUniServicesAgent(t *testing.T) {
	var (
		prompt = `EDIT THIS PROMPT`
		apis   = []interface{}{
			abc_api_client.GetStudentDetails(),

			abc_api_client.GetModules(),
			abc_api_client.GetDegreeRegulations(),
			abc_api_client.GetAcademicPeriods(),
			abc_api_client.GetYears(),
			abc_api_client.GetCohorts(),

			abc_api_client.GetStudents(),
			abc_api_client.GetAllStaffList(),
			abc_api_client.GetStaff(),
			abc_api_client.GetPublicModuleTypes(),
			abc_api_client.GetAllStudents(),
			abc_api_client.GetTotalEnrolledStudents(),
		}

		query = `EDIT THIS QUERY`
	)

	agent := NewAdminUniServicesAgent(prompt, apis...)
	runAgentAndPrintResultFrom(agent, query)
}

func Test_CareerProfessionalDevelopmentAgent(t *testing.T) {
	var (
		prompt = `EDIT THIS PROMPT`
		apis   = []interface{}{
			abc_api_client.GetStudentDetails(),

			abc_api_client.GetDegreeRegulations(),
			abc_api_client.GetPhdStudentInformation(),
			abc_api_client.GetPhdStudentsForSupervisor(),

			abc_api_client.GetAllStaffList(),
			abc_api_client.GetStaff(),
			abc_api_client.GetModules(),
			abc_api_client.GetPublicModuleTypes(),
		}

		query = `EDIT THIS QUERY`
	)

	agent := NewCareerProfessionalDevelopmentAgent(prompt, apis...)
	runAgentAndPrintResultFrom(agent, query)
}

func Test_WellbeingMentalHealthPersonalDevelopmentAgent(t *testing.T) {
	var (
		prompt = `EDIT THIS PROMPT`
		apis   = []interface{}{
			abc_api_client.GetStudentDetails(),
			abc_api_client.GetTutorialGroups(),
			abc_api_client.GetPersonalTuteesForTutor(),
			abc_api_client.GetEnrolledStudents(),
		}

		query = `EDIT THIS QUERY`
	)

	agent := NewWellbeingMentalHealthPersonalDevelopmentAgent(prompt, apis...)
	runAgentAndPrintResultFrom(agent, query)
}

// ALEX
func Test_FinancialAccomodationResourceAgent(t *testing.T) {
	var (
		prompt = `EDIT THIS PROMPT`
		apis   = []interface{}{
			abc_api_client.GetStudentDetails(),

			abc_api_client.GetYears(),
			abc_api_client.GetCohorts(),
			abc_api_client.GetAcademicPeriods(),
			abc_api_client.GetStaff(),
		}

		query = `EDIT THIS QUERY`
	)

	agent := NewFinancialAccomodationResourceAgent(prompt, apis...)
	runAgentAndPrintResultFrom(agent, query)
}

// ALEX
func Test_CampusLifeSocialAgent(t *testing.T) {
	var (
		prompt = `EDIT THIS PROMPT`
		apis   = []interface{}{
			abc_api_client.GetStudentDetails(),

			abc_api_client.GetYears(),
			abc_api_client.GetAcademicPeriods(),
			abc_api_client.GetCohorts(),
			abc_api_client.GetEnrolledStudents(),
			abc_api_client.GetAllStaffList(),
			abc_api_client.GetAllStudents(),
		}

		query = `EDIT THIS QUERY`
	)

	agent := NewCampusLifeSocialAgent(prompt, apis...)
	runAgentAndPrintResultFrom(agent, query)
}

// TEO
func Test_AccessibilityDisabilityAgent(t *testing.T) {
	var (
		prompt = `EDIT THIS PROMPT`
		apis   = []interface{}{}

		query = `EDIT THIS QUERY`
	)

	agent := NewAccessibilityDisabilityAgent(prompt, apis...)
	runAgentAndPrintResultFrom(agent, query)
}

// TEO
func Test_TransitionDiversityMiscAgent(t *testing.T) {
	var (
		prompt = `EDIT THIS PROMPT`
		apis   = []interface{}{
			abc_api_client.GetStudentDetails(),

			abc_api_client.GetYears(),
			abc_api_client.GetCohorts(),

			abc_api_client.GetTutorialGroups(),
			abc_api_client.GetPersonalTuteesForTutor(),
		}

		query = `EDIT THIS QUERY`
	)

	agent := NewTransitionDiversityMiscAgent(prompt, apis...)
	runAgentAndPrintResultFrom(agent, query)
}

func runAgentAndPrintResultFrom(agent Agent, prompt string) {
	resp, err := agent.Run(context.Background(), prompt)
	if err != nil {
		log.Fatalf("Error running agent: %v", err)
	}

	log.Printf("Result: %+v", resp)
}
