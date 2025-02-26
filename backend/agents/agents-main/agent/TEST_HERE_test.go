package agent

// import (
// 	"context"
// 	"log"
// 	"testing"
// )

// // 1. Angelo
// func Test_NewAcadmemicSupportStudyAgent(t *testing.T) {
// 	query := `I am struggling a bit with the Advanced Computer Architecture module. What methods are there to contact the lecturer for help?`
// 	agent := NewAcadmemicSupportStudyAgent()
// 	runAgentAndPrintResultFrom(agent, query)
// }

// // 2. Angelo
// func Test_NewAdminUniServicesAgent(t *testing.T) {
// 	query := `How do I get my student ID card replaced? I lost mine yesterday.`
// 	agent := NewAdminUniServicesAgent()
// 	runAgentAndPrintResultFrom(agent, query)
// }

// // 3. Anshul
// // Additional RAG pages:
// // - https://www.imperial.ac.uk/placements/the-inplace-system/
// func Test_CareerProfessionalDevelopmentAgent(t *testing.T) {
// 	query := `What kind of research would I do in a UROPS? Anything relevant to me?`
// 	agent := NewCareerProfessionalDevelopmentAgent()
// 	runAgentAndPrintResultFrom(agent, query)
// }

// // 4. Anshul
// func Test_WellbeingMentalHealthPersonalDevelopmentAgent(t *testing.T) {
// 	query := `I really don't want to work anymore. There's too much to do`
// 	agent := NewWellbeingMentalHealthPersonalDevelopmentAgent()
// 	runAgentAndPrintResultFrom(agent, query)
// }

// // 5. ALEX
// func Test_FinancialAccomodationResourceAgent(t *testing.T) {
// 	query := `I cant afford tuition fees`
// 	agent := NewFinancialAccomodationResourceAgent()
// 	runAgentAndPrintResultFrom(agent, query)
// }

// //TODO (problems):
// //TODO: Sometimes returns empty answer?

// // 6. ALEX
// func Test_CampusLifeSocialAgent(t *testing.T) {
// 	query := `How do i book a room in the huxley building?`
// 	agent := NewCampusLifeSocialAgent()
// 	runAgentAndPrintResultFrom(agent, query)
// }

// //TODO (for me - alex): see if can extract more from the apis

// // 7. TEO
// func Test_AccessibilityDisabilityAgent(t *testing.T) {
// 	query := `Can I get extra time in exams due to my ADHD?`
// 	agent := NewAccessibilityDisabilityAgent()
// 	runAgentAndPrintResultFrom(agent, query)
// }

// // 8. TEO
// func Test_TransitionDiversityMiscAgent(t *testing.T) {
// 	query := `How can I improve my public speaking skills for presentations?`
// 	agent := NewTransitionDiversityMiscAgent()
// 	runAgentAndPrintResultFrom(agent, query)
// }

// func runAgentAndPrintResultFrom(agent Agent, prompt string) {
// 	resp, err := agent.Run(context.Background(), prompt)
// 	if err != nil {
// 		log.Fatalf("Error running agent: %v", err)
// 	}

// 	log.Printf("Answer: %s\nReason: %s", *resp.Answer, *resp.Reason)
// }
