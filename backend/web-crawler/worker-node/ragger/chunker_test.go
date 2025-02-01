package ragger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/sugarme/tokenizer"
	"github.com/sugarme/tokenizer/pretrained"
)

func TestChunk(t *testing.T) {
	tokenizerPath := "../model/tokenizer.json"
	tok, err := pretrained.FromFile(tokenizerPath)
	if err != nil {
		t.Fatalf("failed to load tokenizer: %v", err)
	}

	chunker := NewChunker(tok)

	text := `A paragraph (from Ancient Greek παράγραφος (parágraphos) 'to write beside') is a self-contained unit of discourse in writing dealing with a particular point or idea. Though not required by the orthographic conventions of any language with a writing system, paragraphs are a conventional means of organizing extended segments of prose.

History
The oldest classical British and Latin writings had little or no space between words and could be written in boustrophedon (alternating directions). Over time, text direction (left to right) became standardized. Word dividers and terminal punctuation became common. The first way to divide sentences into groups was the original paragraphos, similar to an underscore at the beginning of the new group.[1] The Greek parágraphos evolved into the pilcrow (¶), which in English manuscripts in the Middle Ages can be seen inserted inline between sentences.


Indented paragraphs demonstrated in the US Constitution
Ancient manuscripts also divided sentences into paragraphs with line breaks (newline) followed by an initial at the beginning of the next paragraph. An initial is an oversized capital letter, sometimes outdented beyond the margin of the text. This style can be seen, for example, in the original Old English manuscript of Beowulf. Outdenting is still used in English typography, though not commonly.[2] Modern English typography usually indicates a new paragraph by indenting the first line. This style can be seen in the (handwritten) United States Constitution from 1787. For additional ornamentation, a hedera leaf or other symbol can be added to the inter-paragraph white space, or put in the indentation space.

A second common modern English style is to use no indenting, but add vertical white space to create "block paragraphs." On a typewriter, a double carriage return produces a blank line for this purpose; professional typesetters (or word processing software) may put in an arbitrary vertical space by adjusting leading. This style is very common in electronic formats, such as on the World Wide Web and email. Wikipedia itself employs this format.

Typographical considerations
Professionally printed material in English typically does not indent the first paragraph, but indents those that follow. For example, Robert Bringhurst states that we should "Set opening paragraphs flush left."[2] Bringhurst explains as follows:

The function of a paragraph is to mark a pause, setting the paragraph apart from what precedes it. If a paragraph is preceded by a title or subhead, the indent is superfluous and can therefore be omitted.[2]

The Elements of Typographic Style states that "at least one en [space]" should be used to indent paragraphs after the first,[2] noting that that is the "practical minimum".[3] An em space is the most commonly used paragraph indent.[3] Miles Tinker, in his book Legibility of Print, concluded that indenting the first line of paragraphs increases readability by 7%, on average.[4]

When referencing a paragraph, typographic symbol U+00A7 § SECTION SIGN (&sect;) may be used: "See § Background".

In modern usage, paragraph initiation is typically indicated by one or more of a preceding blank line, indentation, an "Initial" ("drop cap") or other indication. Historically, the pilcrow symbol (¶) was used in Latin and western European languages. Other languages have their own marks with similar function.

Widows and orphans occur when the first line of a paragraph is the last in a column or page, or when the last line of a paragraph is the first line of a new column or page.

In computing
See also: Newline
In word processing and desktop publishing, a hard return or paragraph break indicates a new paragraph, to be distinguished from the soft return at the end of a line internal to a paragraph. This distinction allows word wrap to automatically re-flow text as it is edited, without losing paragraph breaks. The software may apply vertical white space or indenting at paragraph breaks, depending on the selected style.

How such documents are actually stored depends on the file format. For example, HTML uses the <p> tag as a paragraph container. In plaintext files, there are two common formats. The pre-formatted text will have a newline at the end of every physical line, and two newlines at the end of a paragraph, creating a blank line. An alternative is to only put newlines at the end of each paragraph, and leave word wrapping up to the application that displays or processes the text.

A line break that is inserted manually, and preserved when re-flowing, may still be distinct from a paragraph break, although this is typically not done in prose. HTML's <br /> tag produces a line break without ending the paragraph; the W3C recommends using it only to separate lines of verse (where each "paragraph" is a stanza), or in a street address.[5]

Numbering
Main article: Dot-decimal notation
See also: ISO 2145
Paragraphs are commonly numbered using the decimal system, where (in books) the integral part of the decimal represents the number of the chapter and the fractional parts are arranged in each chapter in order of magnitude. Thus in Whittaker and Watson's 1921 A Course of Modern Analysis, chapter 9 is devoted to Fourier Series; within that chapter §9.6 introduces Riemann's theory, the following section §9.61 treats an associated function, following §9.62 some properties of that function, following §9.621 a related lemma, while §9.63 introduces Riemann's main theorem, and so on. Whittaker and Watson attribute this system of numbering to Giuseppe Peano on their "Contents" page, although this attribution does not seem to be widely credited elsewhere.[6] Gradshteyn and Ryzhik is another book using this scheme since its third edition in 1951.

Section breaks
Main article: Section (typography)
Many published books use a device to separate certain paragraphs further when there is a change of scene or time. This extra space, especially when co-occurring at a page or section break, may contain a special symbol known as a dinkus, a fleuron, or a stylistic dingbat.

Style advice
The crafting of clear, coherent paragraphs is the subject of considerable stylistic debate. The form varies among different types of writing. For example, newspapers, scientific journals, and fictional essays have somewhat different conventions for the placement of paragraph breaks.

A common English usage misconception is that a paragraph has three to five sentences; single-word paragraphs can be seen in some professional writing, and journalists often use single-sentence paragraphs.[7]

English students are sometimes taught that a paragraph should have a topic sentence or "main idea", preferably first, and multiple "supporting" or "detail" sentences that explain or supply evidence. One technique of this type, intended for essay writing, is known as the Schaffer paragraph. Topic sentences are largely a phenomenon of school-based writing, and the convention does not necessarily obtain in other contexts.[8] This advice is also culturally specific, for example, it differs from stock advice for the construction of paragraphs in Japanese (translated as danraku 段落).[9]

See also
Inverted pyramid (journalism)
Notes
 Edwin Herbert Lewis (1894). The History of the English Paragraph. University of Chicago Press. p. 9.
 Bringhurst, Robert (2005). The Elements of Typographic Style. Vancouver: Hartley and Marks. p. 39. ISBN 0-88179-206-3.
 Bringhurst, Robert (2005). The Elements of Typographic Style. Vancouver: Hartley and Marks. p. 40. ISBN 0-88179-206-3.
 Tinker, Miles A. (1963). Legibility of Print. Iowa: Iowa State University Press. p. 127. ISBN 0-8138-2450-8.
 "<br>: The Line Break element". MDN Web Docs. Retrieved 15 March 2018.
 Kowwalski, E. (3 June 2008). "Peano paragraphing". blogs.ethz.ch.
 University of North Carolina at Chapel Hill. "Paragraph Development". The Writing Center. University of North Carolina at Chapel Hill. Retrieved 20 June 2018.
 Braddock, Richard (1974). "The Frequency and Placement of Topic Sentences in Expository Prose". Research in the Teaching of English. 8 (3): 287–302.
 com), Kazumi Kimura and Masako Kondo (timkondo *AT* nifty . com / Kazumikmr *AT* aol . "Effective writing instruction: From Japanese danraku to English paragraphs". jalt.org. Retrieved 15 March 2018.
References
The American Heritage Dictionary of the English Language. 4th ed. New York: Houghton Mifflin, 2000.
Johnson, Samuel. Lives of the Poets: Addison, Savage, etc.. Project Gutenberg, November 2003. E-Book, #4673.
Rozakis, Laurie E. Master the AP English Language and Composition Test. Lawrenceville, NJ: Peterson's, 2000. ISBN 0-7645-6184-7 (10). ISBN 978-0-7645-6184-9 (13).
External links
 The dictionary definition of paragraph at Wiktionary
vte
Typography
Page	
Canons of page constructionColumnEven workingMarginPage numberingPaper sizePaginationPull quoteRecto and versoIntentionally blank page
Paragraph	
AlignmentLeadingLine lengthRiverRunaroundWidows and orphans runt
Character	
Typeface anatomy	
CounterDiacriticsDingbatGlyphInk trapLigatureRotationSubscript and superscriptSwashText figuresTittle
Capitalization	
All capsCamel caseInitialLetter caseSmall capsSnake caseTitle case
Visual distinction	
Blackboard boldBoldColor printingItalicsObliqueUnderlineWhitespace
Horizontal aspects	
Figure spaceKerningLetter spacingPitchSentence spacingThin spaceWord spacing
Vertical aspects	
AscenderBaselineBody heightCap heightDescenderMean lineOvershootx-height
Typeface
classifications	
Roman type	
Serif AntiquaDidoneslab serifSans-serif
Blackletter type	
FrakturRotundaSchwabacher
Gaelic type	
InsularUncial
Specialist	
Record typeDisplay typeface scriptfat facereverse-contrast
Punctuation (List)	
BulletDashHanging punctuationHyphen minus signInterpunctSpaceVertical bar
Typesetting	
Etaoin shrdluFont computermonospacedFont catalogFor position onlyLetterpressLorem ipsumMicroprintingMicrotypographyMovable typePangramPhototypesettingPunchcuttingReversing typeSortType colorType designTypeface list
Typographic units	
AgateCiceroEmEnMetric unitsPicaPoint traditional point-size namesTwip
Digital typography	
Character encodingHintingText shapingRasterizationTypographic featuresWeb typographyBézier curvesDesktop publishing
Typography in other
writing systems	
ArabicCyrillic PT FontsEast AsianThai National Fonts
Related articles	
Penmanship HandwritingHandwriting scriptCalligraphyStyle guideType designType foundryHistory of Western typographyIntellectual property protection of typefacesTechnical letteringVox-ATypI classification
Related template	
Punctuation and other typographic symbols`

	chunks, err := chunker.Chunk(text)
	if err != nil {
		t.Fatalf("failed to chunk text: %v", err)
	}

	for _, chunk := range chunks {
		t.Log(chunk)

		tokens, err := tok.Encode(tokenizer.NewSingleEncodeInput(tokenizer.NewInputSequence(chunk)), true)
		if err != nil {
			t.Fatalf("failed to encode chunk: %v", err)
		}
		assert.LessOrEqual(t, len(tokens.Ids), maxTokenChunkSize)
	}

	assert.Equal(t, 16, len(chunks))
}

func TestChunkerPanic(t *testing.T) {
	f := `In this section\n\n- [Study](/study/)\n\n- [Submit your offer conditions](/study/apply/undergraduate/offer-holders/next-steps/submit-offer-conditions/)\n- [Request a deferral](/study/apply/undergraduate/offer-holders/next-steps/deferred-entry/)\n\n[Study](/study/)\n\n- [Imperial Home](/)\n- [Study](/study/)\n- [Apply](/study/apply/)\n- [Undergraduate](/study/apply/undergraduate/)\n- [Offer holders](/study/apply/undergraduate/offer-holders/)\n- [Next steps](/study/apply/undergraduate/offer-holders/next-steps/)\n\n# Next steps\n\n[Offer holders](/study/apply/undergraduate/offer-holders/)\n\n- [Submit your offer conditions](/study/apply/undergraduate/offer-holders/next-steps/submit-offer-conditions/)\n- [Request a deferral](/study/apply/undergraduate/offer-holders/next-steps/deferred-entry/)\n\n![](https://pxl-imperialacuk.terminalfour.net/fit-in/1079x305/prod01/channel_3/media/study---restricted-images/offer-holder-hub/next-steps/nextsteps_content_header_3000X850.jpg)\n\nCongratulations on your offer! What's next?\n\n## Find out your next steps\n\nWe've put together key information to help you understand your offer, how to reply and all the key dates and deadlines you need to know.\n\nYou can also find information about accommodation, funding your studies and how to send us evidence of any additional qualifications you need to meet your offer conditions.\n\n## Next steps\n\n- [1 – Reply to your offer](#tabpanel_1456127)\n- [2 – Consider your accommodation options](#tabpanel_1456128)\n- [3 – Explore your funding options](#tabpanel_1456129)\n- [4 – Tell us about additional qualifications](#tabpanel_1456130)\n- [5 – Welcome Season](#tabpanel_1456131)\n\n1 – Reply to your offer\n\n\nAfter receiving decisions on all your applications, it’s time to reply to your offers by the deadline.\n\nIf you receive all of your offers by **Wednesday 14 May 2025**, you must reply to UCAS by **Wednesday** **4 June 2025**.\n\nIf you decide to accept your offer, you can choose to make Imperial either your firm or insurance choice.\n\n#### Firm\n\nThis means we’re your first choice. If your offer is unconditional, then congratulations – you’re in! If it’s conditional, your place will be confirmed after you’ve met the conditions of your offer.\n\n#### Insurance\n\nThis means we’re your backup choice. You’ll be offered a place here if you don’t meet the conditions of your firm choice offer but meet the conditions of your Imperial offer.\n\nFind out how to reply to your offer on [UCAS](https://www.ucas.com/undergraduate/after-you-apply/types-offer/replying-your-ucas-undergraduate-offers \"UCAS\").\n\n2 – Consider your accommodation options\n\n\nImperial accommodation is well-maintained with round-the-clock support, and the rent includes all bills. If you pick Imperial as your firm choice, we offer a guaranteed place in our accommodation if you want it.\\*\n\nIf Imperial is your firm choice, you can apply for accommodation in late May. You’ll be invited to apply in late July if we're your insurance choice. You can choose your preferences based on your preferred hall, room type, and price.\n\nWe aim to make our accommodation the best choice for as many students as possible, but we know it won’t be suitable for everyone. [Our Accommodation Office](/students/accommodation/private-accommodation/) can also offer advice about renting privately and how to find the right place for you.\n\n[Visit the Accommodation Office website to learn more.](/students/accommodation/)\n\n\\*You’re guaranteed a place in our accommodation if you make Imperial your firm choice, submit your accommodation application by **18 July 2025**, meet all conditions of your offer and hold an unconditional offer by (date to be confirmed; as a guide, the unconditional offer deadline for 2024-25 entry was 23 August 2024).\n\n3 – Explore your funding options\n\n\nLearn more about tuition fees and funding options at Imperial, including our Imperial Bursary, student loans and scholarships.\n\nFind out more on our Offer Holder [fees and funding page](/study/apply/undergraduate/offer-holders/fees-and-funding/).\n\n4 – Tell us about additional qualifications\n\n\nUCAS will send us your A-level, International Baccalaureate (IB), Irish Leaving Certificates, Pre-U Certificate, Scottish Advanced Highers and STEP grades on your results day. However, you may need to send us the results for other qualifications.\n\nFind out if that [includes your qualifications](/study/apply/undergraduate/offer-holders/next-steps/submit-offer-conditions/).\n\n5 – Welcome Season\n\n\nThe first month of our new term is Welcome Season, which **Welcome Week** is packed with opportunities to engage with the Imperial community for the first time.\n\nIt’s organised by Imperial and our Students’ Union and gives you the chance to make friends, get involved in loads of events and activities, and get to know your new home at Imperial.\n\n## Key dates\n\n- [May 2025](#tabpanel_1456133)\n- [June 2025](#tabpanel_1456134)\n- [July 2025](#tabpanel_1456135)\n- [August 2025](#tabpanel_1456136)\n- [September 2025](#tabpanel_1456137)\n\nMay 2025\n\n\n**Late May 2025 -** If you’ve made Imperial your firm choice on UCAS, you’ll be invited to apply for accommodation.\n\nJune 2025\n\n\n**5 June 2025 -** If you receive all of your offers by 14 May 2025, you must reply to UCAS by 5 June 2025.\n\nJuly 2025\n\n\n**5 July 2025 -** International Baccalaureate results day\n\n**18 July 2025** **-** Accommodation guarantee deadline\n\n**24 July 2025** **-** English Language and completed qualifications deadline\n\n**Late July** **-** If Imperial is your insurance choice on UCAS, you’ll be invited to apply for accommodation.\n\nAugust 2025\n\n\n**14 August 2025** **-** A-level results day\n\n**22 August 2025** **-** Deadline for receiving results that need additional, non-UCAS verification\n\n**31 August 2025** **-** Deadline to submit amended grades as a result of an appeal\n\nSeptember 2025\n\n\n**Early September 2025** **-** Accommodation offers sent out\n\n**27 September 2025** **-** Term starts\n\n[![](https://pxl-imperialacuk.terminalfour.net/fit-in/720x462/prod01/channel_3/media/study---restricted-images/offer-holder-hub/Contact-180227_southken_campus_snow_Pano.jpg)\\\n\\\n**Contact us** \\\n\\\nIf you have any questions about your offer, you can contact the relevant Admissions team.\\\n\\\nContact us](/study/apply/contact/)\n\n[![](https://pxl-imperialacuk.terminalfour.net/fit-in/720x462/prod01/channel_3/media/study---restricted-images/offer-holder-hub/Submit-offer-conditions.jpg)\\\n\\\n**Submit your offer conditions** \\\n\\\nWe may ask you to submit your results/certificates of qualifications that aren't sent to us via UCAS.\\\n\\\nFind out more](/study/apply/undergraduate/offer-holders/next-steps/submit-offer-conditions/)\n\n[![](https://pxl-imperialacuk.terminalfour.net/fit-in/720x462/prod01/channel_3/media/study---restricted-images/offer-holder-hub/Request-a-deferral.jpg)\\\n\\\n**Request a deferral** \\\n\\\nYou may apply to defer your admission to the next academic year, under exceptional circumstances. \\\n\\\nFind out more](/study/apply/undergraduate/offer-holders/next-steps/deferred-entry/)\n\n## Offer holders\n\n[Next steps](/study/apply/undergraduate/offer-holders/next-steps/)\n\n[Life on campus](/study/apply/undergraduate/offer-holders/life-on-campus/)\n\n[Accommodation](/study/apply/undergraduate/offer-holders/accommodation/)\n\n[Fees and funding](/study/apply/undergraduate/offer-holders/fees-and-funding/)\n\n[Support and FAQs](/study/apply/undergraduate/offer-holders/support/)\n\n[After you graduate](/study/apply/undergraduate/offer-holders/after-you-graduate/)\n\n[International students](/study/apply/undergraduate/offer-holders/international-students/)\n\n[Return to offer holder home](/study/apply/undergraduate/offer-holders/)\n\n[Accept your offer](https://www.ucas.com/undergraduate/after-you-apply/ucas-undergraduate-types-offer)\n\n[Submit your offer conditions](/study/apply/undergraduate/offer-holders/next-steps/submit-offer-conditions/)`

	tok, err := pretrained.FromFile("../model/tokenizer.json")
	if err != nil {
		t.Fatalf("failed to load tokenizer: %v", err)
	}
	chunker := NewChunker(tok)
	chunks, err := chunker.Chunk(f)
	if err != nil {
		t.Fatalf("failed to chunk text: %v", err)
	}

	for _, chunk := range chunks {
		t.Log(chunk)
	}
}
