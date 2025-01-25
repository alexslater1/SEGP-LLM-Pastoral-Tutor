import PyPDF2
import nltk
import re

nltk.download('punkt')

def extract_text_from_pdf(pdf_path):
    """Extract text from a PDF file."""
    text = ""
    with open(pdf_path, 'rb') as file:
        pdf_reader = PyPDF2.PdfReader(file)
        for page in pdf_reader.pages:
            text += page.extract_text()
    return text

def clean_text(text):
    """Clean the extracted text to remove unwanted artifacts."""
    # Normalize whitespace
    text = re.sub(r'\s+', ' ', text)
    # Remove headers, footers, or patterns (e.g., page numbers or URLs)
    text = re.sub(r'Page \d+|https?://\S+|\d+\s*$', '', text)
    # Remove placeholder tags like {insert here}
    text = re.sub(r'\{.*?\}', '', text)
    return text.strip()

def split_into_sentences(text):
    """Split text into sentences with additional handling for bullet points."""
    # Use nltk to split sentences initially
    sentences = nltk.sent_tokenize(text)
    # Further split sentences that contain bullet points or lists
    split_sentences = []
    for sentence in sentences:
        # Split by bullet points or semicolons
        parts = re.split(r'[•;]', sentence)
        split_sentences.extend([part.strip() for part in parts if part.strip()])
    return split_sentences

def document_sentence_chunker(pdf_path):

    # Extract and clean text
    pdf_text = extract_text_from_pdf(pdf_path)
    cleaned_text = clean_text(pdf_text)

    # Split text into sentences
    sentences = split_into_sentences(cleaned_text)

    return sentences
