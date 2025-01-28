package ragger

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContactsFor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Contact
	}{
		{
			name:  "should extract email",
			input: "You can reach me at john.doe@example.com for more information.",
			expected: []Contact{
				{
					Value: "john.doe@example.com",
					Type:  "email",
				},
			},
		},
		{
			name:  "should extract phone number",
			input: "Call me at +1 (555) 123-4567 anytime.",
			expected: []Contact{
				{
					Value: "+1 (555) 123-4567",
					Type:  "phone",
				},
			},
		},
		{
			name:  "should extract URL",
			input: "Visit our website at https://example.com/contact for details.",
			expected: []Contact{
				{
					Value: "https://example.com/contact",
					Type:  "url",
				},
			},
		},
	}

	client := &RagClient{} // No need to initialize with real model for contact extraction

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.ContactsFor(tt.input)
			if err != nil {
				t.Errorf("ContactsFor() error = %v", err)
				return
			}

			if len(got) != len(tt.expected) {
				t.Errorf("ContactsFor() got %d contacts, expected %d", len(got), len(tt.expected))
				return
			}

			for i, contact := range got {
				if contact.Value != tt.expected[i].Value {
					t.Errorf("ContactsFor() got value = %v, expected %v", contact.Value, tt.expected[i].Value)
				}
				if contact.Type != tt.expected[i].Type {
					t.Errorf("ContactsFor() got type = %v, expected %v", contact.Type, tt.expected[i].Type)
				}
			}
		})
	}
}

func TestChunksFor(t *testing.T) {
	client := NewRagClient("../model", "../libonnxruntime.so.1.20.1")

	text := `1 
 
	Student Code of Conduct 
	 
	The Student Code of Conduct is designed to supplement the College's Code of Ethics by providing 
	detailed explanations and examples of unacceptable behaviour.  
	The purpose of this Code is to encourage and maintain standards of conduct by students consistent 
	with the values and expectations of the Imperial community. 
	The Code applies to conduct both on and off the College's campuses and premises, as well as online. 
	It also applies to students engaged in Imperial College Union (ICU) activities. 
	The Code is complementary to, and does not replace, other standards, regulations, or professional 
	conduct requirements applying to students in the College. 
	The Code of conduct should be read in conjunction with the College's Ethics Code. 
	Respect for each other and our surroundings 
	• 
	Students should treat others as they themselves would like to be treated, with dignity and 
	due respect at all times. The College community is one in which discrimination, bullying, 
	harassment, and victimisation are never tolerated. 
	• 
	Students should treat others equitably and work to create an inclusive environment in which 
	everyone is safe to speak up and share their perspective. Students are encouraged to be 
	curious and seek to understand diverse perspectives. 
	• 
	Students should take responsibility for their behaviour and their impact on others. Students 
	should consider and respond to the needs of others, ensuring communications with others 
	are considerate and respectful. 
	• 
	Our campuses, property and facilities should be treated with respect. They are for the safe 
	and enjoyable use of all our community and should be used for their designated purposes 
	and not intentionally or recklessly damaged or defaced. 
	• 
	Whilst freedom of speech or expression is an important right for all in our community, it is 
	not an unqualified right. It is important to remember that a person's right to freedom of 
	speech means lawful freedom of speech. That means that speech (or other expressions of 
	views, such as slogans, tattoos, literature, emails, etc), that might be criminal in nature (for 
	example, inciting racial hatred; or displaying threatening, abusive or insulting writing likely 
	to cause harassment, alarm or distress), or otherwise in breach of civil law (for example, 
	causing a breach of the peace, or harassment), is not protected by a person's right to 
	freedom of speech or expression. It is very important therefore that members of our 
	community engage in debate in a way that is respectful: how something is said or expressed 
	needs to be given as much consideration as what is said or expressed.  
	• 
	Other examples of misconduct include but are not restricted to: 
	o disruption of, or improper interference with, the academic, administrative, sporting, 
	social or other activities of the College or Imperial College Union; 
	o disruption of College business, including teaching, research and studying that is not 
	authorised pursuant to a College recognised ballot or process;  
	o misuse or unauthorised use of College premises, facilities, or items of property; 
	2 
	 
	o disregard of the health and safety of self or others whilst on College campuses or 
	undertaking College activities; 
	o disregard for laboratory safety requirements following a warning from teaching staff 
	or laboratory technicians; 
	o wearing clothes, other items, or having visible tattoos with slogans or symbols that 
	might constitute a breach of the criminal or civil law of England and Wales; 
	o distributing material, including online material, which is intimidating, threatening, 
	indecent or illegal; and intentionally or recklessly harming other individuals or 
	putting others at risk of harm. 
	Honesty and integrity  
	• 
	Students should demonstrate a commitment to independence, honesty, and transparency. 
	They should be honest and truthful in their dealings with each other and with third parties.  
	• 
	Students engaged in research activities should conduct their research in a way that supports 
	public trust and confidence in the College's research methods and findings. They should 
	demonstrate rigour, honesty and integrity, and abide by relevant ethical and legal standards.  
	• 
	Students should not become complicit in any activities in which a student gains an unfair 
	advantage, through plagiarism, self-plagiarism, collusion, examination offences, dishonest 
	practice, or other means. Students should comply with the College's Policy on Academic 
	Misconduct.  
	• 
	Examples of misconduct include but are not restricted to: 
	o behaviour which brings the College into disrepute [this does not include 
	whistleblowing];  
	o academic misconduct: cheating, fabrication, plagiarism, collusion, facilitating 
	academic dishonesty, claiming authorship of others' work, or the submission of work 
	for assessment which has been generated through an artificial intelligence of 
	translation programme without acknowledgement or authorisation; 
	o research misconduct: fabrication or falsification of information or data, 
	misrepresentation of data and/or interests or involvement, plagiarism, and failure to 
	follow accepted procedures or to exercise due care in carrying out research; 
	o failure to disclose one's name or other relevant details to an officer or employee of 
	the College or ICU in circumstances where it is reasonable to require that such 
	information be given; 
	o making vexatious (such as frivolous allegations, or repeated allegations based on 
	substantially the same matter that has been dealt with), allegations against a 
	member of the College (allegations that are made in good faith are not vexatious, 
	even if they are not upheld after they have been investigated); and 
	o professional conduct violations. 
	Sexual misconduct and abuse 
	• 
	Sexual misconduct (i.e., sexual harassment, sexual violence) is never tolerated. All students 
	are expected to act to ensure a working and learning environment free from these 
	behaviours.   
	3 
	 
	• 
	Sexual misconduct is any act of violence or harassment which is sexual in nature or any kind 
	of unwanted, non-consensual1 sexual touching, or harassment, within or outside a 
	relationship. This may include rape, sexual assault, sexual exploitation or groping. It also 
	covers behaviours such as grooming, coercion, the promise of a reward for sexual access and 
	sexual demands or threats. 
	• 
	In addition to these, other examples, performed without consent, which might constitute 
	sexual misconduct include: 
	o sexually explicit remarks, innuendos or banter; 
	o sexual insults, jokes, teasing or songs; 
	o wolf whistling, cat calling or making other offensive sexual noises; 
	o offensive comments about someone's dress, appearance, or private life, including 
	their sexuality or gender identity; 
	o unwanted or inappropriate physical contact including touching, pinching, groping or 
	smacking; 
	o unwanted requests to engage in or discuss sexual activity; 
	o lifting or removing clothing;  
	o stalking;  
	o using humour to cover or deflect where sexual misconduct has occurred; 
	o display or distribution of pornographic or sexually explicit material. 
	 
	• 
	Stalking is persistent and unwanted conduct of one or more kinds of behaviours described 
	above. It can be physical or psychological and take place directly against a person, or by 
	approaching a third party about a person. The more common examples of stalking are 
	following a person home, following a person around, between or to/from campus, sending 
	or leaving them unwanted and repeated messages, bullying them on social media or making 
	intrusive or unwanted visits. Interpersonal relationships between individuals can also be 
	abusive without a sexual element to the behaviour, and can include emotional, financial or 
	physical abuse, threats, isolation or intimidation. They may involve bullying or coercive 
	behaviours used to maintain power or control.  
	Bullying, harassment, and discrimination 
	• 
	We do not tolerate bullying or harassment. 
	`

	chunks, err := client.ChunksFor(text)
	if err != nil {
		t.Fatalf("error getting chunks %v", err)
	}

	assert.Equal(t, 9, len(chunks))

	for _, chunk := range chunks {
		log.Printf("chunk: %v \n=============\n", chunk)
	}
}
