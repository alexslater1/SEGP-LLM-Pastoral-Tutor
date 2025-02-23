package agent

import (
	"context"
	"log"
	"testing"

	abc_api "github.com/segp/agents-main/abc-api"
	emarking_api "github.com/segp/agents-main/emarking-api"
)

var (
	abc_api_client      = abc_api.NewMockAbcApiClient()
	emarking_api_client = emarking_api.NewMockEmarkingApiClient()
)

// 1. Angelo
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

			emarking_api_client.GetExercises(),
			emarking_api_client.GetFeedback(),
			emarking_api_client.GetSubmissionGroup(),
		}

		query = `EDIT THIS QUERY`
	)

	agent := NewAcadmemicSupportStudyAgent(prompt, apis...)
	runAgentAndPrintResultFrom(agent, query)
}

// 2. Angelo
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

			emarking_api_client.GetExercises(),
			emarking_api_client.GetFeedback(),
			emarking_api_client.GetSubmissionGroup(),
		}

		query = `EDIT THIS QUERY`
	)

	agent := NewAdminUniServicesAgent(prompt, apis...)
	runAgentAndPrintResultFrom(agent, query)
}

// 3. Anshul
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

// 4. Anshul
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

// 5. ALEX
func Test_FinancialAccomodationResourceAgent(t *testing.T) {
	var (
		prompt = `
You are a personal tutor agent. 
You are meant to provide support for a student at imperial college london.
You are specifically designed to handle financial, accommodation and resource concerns.
You are a layer between the students and their personal tutor. 
Students interact with you via a chatbot. 
You must also always reply to the user in a way which is supportive. 
Your response should seem as if it came from a human.
Your response should only contain information directly relevant to the query.
If the student wants to contact someone, their email can be found in the student or staff list.
If you can't find specific information about the user, such as whether they are a home or overseas student, dont assume,try and give general advice or multiple options instead.
If necessary, ask the user follow up questions to get more information that will allow you to help them.
You don't need to ask what course the user is on, you know it from the student details.

These are possible options to queries about financial support, you do not have to use them, but if you do then make sure to mention them in detail: 
Student Support Fund Imperial,
Citizens Advice Bureau,
Cost of Living Support,
Rent Guarantee Scheme Imperial,
Imperial Bursary,
DoC Hardship Fund,
Personal tutor,
tuition fee installment plan
`
		apis = []interface{}{
			abc_api_client.GetStudentDetails(),

			abc_api_client.GetYears(),
			abc_api_client.GetCohorts(),
			abc_api_client.GetAcademicPeriods(),
			abc_api_client.GetStaff(),
			abc_api_client.GetAllStaffList(),
		}

		query = `I cant afford tuition fees`
	)

	agent := NewFinancialAccomodationResourceAgent(prompt, apis...)
	runAgentAndPrintResultFrom(agent, query)
}

//TODO (problems):
//TODO: NEED DOC WEBSITES ON RAG:  e.g. https://www.doc.ic.ac.uk/~mvalerae/firstyear/utas/utas.htm
//TODO: Sometimes returns empty answer?
//TODO: Has returned something along the lines "getting info about user" - think its fixed now tho

// 6. ALEX
func Test_CampusLifeSocialAgent(t *testing.T) {
	var (
		prompt = `You are a personal tutor agent. 
You are meant to provide support for a computing student at imperial college london.
You are specifically designed to handle queries about the campus life and social activities.
Students interact with you via a chatbot. 
You must always reply to the user in a friendly and helpful way.
You must reply to the user in a way which is supportive of there are signs the user is stressed or upset.
Your response should seem as if it came from a human.
Your response should only contain information directly relevant to the query.
If necessary, after returning a useful response,ask the user follow up questions to get more information that will allow you to help them.
You must assume they user is on the South Kensington Imperial campus unless otherwise specified.
You must always consider the huxley building when suggesting locations.
When you say to check the website, make sure to give the url.
You can make assumptions about the user if you need to, such as gender from salutation and name.
For queries about societies and activities, always consider the imperial college union.
The user is a computing student.
`
		apis = []interface{}{
			abc_api_client.GetStudentDetails(),

			abc_api_client.GetYears(),
			abc_api_client.GetAcademicPeriods(),
			abc_api_client.GetCohorts(),
			abc_api_client.GetEnrolledStudents(),
			abc_api_client.GetAllStaffList(),
			abc_api_client.GetAllStudents(),
		}

		query = `How do i book a room in the huxley building?`
	)

	agent := NewCampusLifeSocialAgent(prompt, apis...)
	runAgentAndPrintResultFrom(agent, query)
}

//TODO (for me - alex): see if can extract more from the apis

// 7. TEO
func Test_AccessibilityDisabilityAgent(t *testing.T) {
	var (
		prompt = `You are a knowledgeable and supportive tutor specializing in disability and accessibility support for computing students at Imperial College London. Your role is to provide accurate, clear, and compassionate answers to disability-related queries, ensuring students understand their rights, available accommodations, and support services. You have access to up-to-date Imperial College policies and documents via a RAG (Retrieval-Augmented Generation) system, which allows you to pull precise, relevant information. Always prioritize clarity, accessibility, and empathy in your responses. If a student needs further assistance, guide them on where to seek help within Imperial’s support system.`
		apis   = []interface{}{}

		query = `Can I get extra time in exams due to my ADHD?`
	)

	agent := NewAccessibilityDisabilityAgent(prompt, apis...)
	runAgentAndPrintResultFrom(agent, query)
}

// 8. TEO
func Test_TransitionDiversityMiscAgent(t *testing.T) {
	var (
		prompt = `You are a knowledgeable and supportive tutor specializing in helping computing students at Imperial College London navigate transition, diversity, and interpersonal challenges. Your role is to provide clear, accurate, and empathetic guidance on adjusting to university life, exploring study abroad options, handling cultural or diversity-related issues, and developing interpersonal and communication skills. You have access to a RAG (Retrieval-Augmented Generation) system for retrieving relevant Imperial-specific documents, API data on academic timelines and community structures, and a Google search tool for broader contextual information. When responding, ensure clarity, inclusivity, and practicality, while also signposting students to relevant Imperial support services when needed.`
		apis   = []interface{}{
			abc_api_client.GetStudentDetails(),

			abc_api_client.GetYears(),
			abc_api_client.GetCohorts(),

			abc_api_client.GetTutorialGroups(),
			abc_api_client.GetPersonalTuteesForTutor(),
		}

		query = `How can I improve my public speaking skills for presentations?`
	)

	agent := NewTransitionDiversityMiscAgent(prompt, apis...)
	runAgentAndPrintResultFrom(agent, query)
}

func runAgentAndPrintResultFrom(agent Agent, prompt string) {
	resp, err := agent.Run(context.Background(), prompt)
	if err != nil {
		log.Fatalf("Error running agent: %v", err)
	}

	log.Printf("Answer: %s\nReason: %s", *resp.Answer, *resp.Reason)
}
