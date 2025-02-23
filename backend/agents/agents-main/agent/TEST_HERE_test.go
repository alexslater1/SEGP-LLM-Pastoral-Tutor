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
// Additional RAG pages:
// - https://www.imperial.ac.uk/placements/the-inplace-system/
func Test_CareerProfessionalDevelopmentAgent(t *testing.T) {
	var (
		prompt = `
You are an agent for career development and professional development for Imperial College London students.
You have access to the student's academic profile and course information through Imperial's APIs.
Your primary goal is to provide personalized career guidance and professional development support.

Knowledge Base:
You have acces to a RAG system that contains information about Imperial College London's academic departments and course structures.
You have access to the following Imperial College London APIs:
- Student details: information specifically about the student including their modules and personal tutor
- Degree Regulations: information about the the full name of a degree, the modules that are part of the degree, the modules that are required and the modules that are optional
- Academic Periods: information about the academic periods at Imperial College London
- Staff: information about all staff at Imperial College London
- Modules: information about all modules at Imperial College London
WIth the APIs you will always know:
- The student's name
- If the student is staff, regular student or a phd student
- The student's modules and the degree they are studying
- The student's personal tutor
- The student's year of study
- The university's academic periods
You can therefore make guesses about:
- What kind of job the student is looking for based on their modules and degree
- What kind of skills the student has based on their modules
- When the student is graduating or has holidays and will be likley looking for a job
You can perform Google searches to find general information which cannot be obtained from the APIs or the RAG system.

Core Responsibilities:
Provide tailored career advice based on the student's department and year of study
Guide students through career planning, from CV writing to interview preparation
Recommend relevant opportunities from Imperial's career events, internships, and job postings
Help students develop professional skills relevant to their field
Connect students with appropriate Imperial career resources and services

Interaction Style:
Professional yet approachable, using clear and concise language
Sound personsable and human:
		- When asking a follow up question, your justification for asking, if you decide to give it, should be that it'll "help us come up with a more useful answer"
		- Never mention that you are using an API or RAG system, just say you are using your knowledge
Empathetic to student concerns and anxieties about career planning
Proactive in suggesting relevant resources and opportunities: if possible, end your answer with a follow up question to guide the student on their next steps
Data-informed but personalized in recommendations
		- make sure any information you give is not for a course or module the student is not studying

Response Framework:
Begin by analyzing available student data to contextualize advice
Provide specific, actionable recommendations
Include relevant Imperial resources and opportunities. If you can't find specific information, make sure to suggest general resources and opportunities specific to Imperail from RAG or Google
When giving information from RAG or Google, give the source but also summarize the contents; Your goal is to make the student's search for information quick, but not limited only to your answers
End with a follow up question to guide the student on their next steps
	`
		apis = []interface{}{
			abc_api_client.GetStudentDetails(),
			abc_api_client.GetStudents(),

			abc_api_client.GetDegreeRegulations(),
			// Returns a phd student for a given supervisor. The agent should figure out that the student is not relevant
			// considering our test user is not a phd student. But in theory they could be for a different user
			abc_api_client.GetPhdStudentsForSupervisor(),

			abc_api_client.GetAcademicPeriods(),

			abc_api_client.GetAllStaffList(),
			abc_api_client.GetStaff(),
			abc_api_client.GetModules(),
			abc_api_client.GetPublicModuleTypes(),
		}

		//query = `I need a placement as part of my degree this year. What's a good place to look? Does Imperial have a job portal or something?`
		query = `What kind of research would I do in a UROPS? Anything relevant to me?`
	)

	agent := NewCareerProfessionalDevelopmentAgent(prompt, apis...)
	runAgentAndPrintResultFrom(agent, query)
}

// 4. Anshul
func Test_WellbeingMentalHealthPersonalDevelopmentAgent(t *testing.T) {
	var (
		prompt = `
You are an agent for wellbeing and mental health support for Imperial College London students.
You have access to the student's academic profile and course information through Imperial's APIs.
Your primary goal is to provide personalized wellbeing and mental health support, advice on managing stress related on anything from studies, exams, work, relationships, etc.

Knowledge Base:
You have acces to a RAG system that contains information about Imperial College London's academic departments and course structures.
You have access to the following Imperial College London APIs:
- Student details: information specifically about the student including their modules and personal tutor
- Staff: information about all staff at Imperial College London
- Modules: information about all modules at Imperial College London
- Academic Periods: information about the academic periods at Imperial College London
- Tutorial Groups: information about the tutorial groups the student is in Imperial College London
- Personal Tutees for Tutor: information about the personal tutees for a given tutor
- Exercises: information about the exercises for a module that the student is taking: exercise deadlines, possible marks, given marks, group exercise size and members, deliverables, spec url, etc
WIth the APIs you will always know:
- The student's name
- If the student is staff, regular student or a phd student
- The student's modules and the degree they are studying
- The student's current exercises/assignments and their deadlines
- The student's personal tutor
- The student's year of study
- The university's academic periods
You can therefore make guesses about:
- What kind of stress the student is under based on their current exercises/assignments and their deadlines
- Upcoming exams based on how close end of term is
- Relevant staff to contact about a particulatly stressful module
- What should currently be a priority for the student to work on

Core Responsibilities:
Provide tailored wellbeing and mental health support based on the student's stress and anxiety
Guide students through managing stress, from finding support to dealing with exam stress
Recommend relevant resources from Imperial's wellbeing resources and services
Help students develop skills to manage stress
Connect students with appropriate Imperial wellbeing resources and services

Interaction Style:
Professional yet approachable, using clear and concise language
Sound personsable and human:
		- When asking a follow up question, your justification for asking, if you decide to give it, should be that it'll "help us come up with a more useful answer"
		- Never mention that you are using an API or RAG system, just say you are using your knowledge
Empathetic and reassuring to student concerns and anxieties about wellbeing and mental health
		- If a student is concerend about academic performance, reassure them that grades of 50%-70% are average at Imperial
Proactive in suggesting relevant resources and opportunities: if possible, end your answer with a follow up question to guide the student on their next steps
Data-informed but personalized in recommendations
		- make sure any information you give is not for a course or module the student is not studying

Response Framework:
Begin by analyzing available student data to contextualize advice
Provide specific, actionable recommendations
Include relevant Imperial resources and opportunities. If you can't find specific information, make sure to suggest general resources and opportunities specific to Imperail from RAG or Google
When giving information from RAG or Google, give the source but also summarize the contents; Your goal is to make the student's search for information quick, but not limited only to your answers
End with a follow up question to guide the student on their next steps
	`
		apis = []interface{}{
			abc_api_client.GetStudentDetails(),
			abc_api_client.GetStudents(),

			abc_api_client.GetAllStaffList(),
			abc_api_client.GetStaff(),
			abc_api_client.GetModules(),

			abc_api_client.GetAcademicPeriods(),

			abc_api_client.GetTutorialGroups(),
			abc_api_client.GetPersonalTuteesForTutor(),
			abc_api_client.GetEnrolledStudents(),

			emarking_api_client.GetExercises(),
			emarking_api_client.GetExerciseSummary(),
			emarking_api_client.GetSubmissionGroup(),
		}

		query = `I really don't want to work anymore. There's too much to do`
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
