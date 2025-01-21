import fitz
from docx import Document
from pptx import Presentation

# def pdf_to_text(pdf_path, txt_path):

#     pdf_document = fitz.open(pdf_path)
#     with open(txt_path, 'w', encoding="utf-8") as text_file:
#         for page_number in range(len(pdf_document)):
#             page = pdf_document.load_page(page_number)
#             text = page.get_text()
#             text_file.write(text)
        
#     pdf_document.close()


def pdf_to_text(pdf_path):
    ret = ""

    pdf_document = fitz.open(pdf_path)
    for page_number in range(len(pdf_document)):
        page = pdf_document.load_page(page_number)
        text = page.get_text()
        ret += text
        
    pdf_document.close()
    return ret

def docx_to_text(file_path):
    doc = Document(file_path)
    full_text = []
    for paragraph in doc.paragraphs:
        full_text.append(paragraph.text)
    return '\n'.join(full_text)

def txt_to_text(file_path):
    with open(file_path, 'r', encoding='utf-8') as file:
        text = file.read()
    return text

def pptx_to_text(file_path):
    presentation = Presentation(file_path)
    text = []
    
    for slide in presentation.slides:
        for shape in slide.shapes:
            if hasattr(shape, "text"):
                text.append(shape.text)
    
    return '\n'.join(text)

def file_to_text(path):
    if path.endswith(".pdf"):
        return pdf_to_text(path)
    elif path.endswith(".docx"):
        return docx_to_text(path)
    elif path.endswith(".txt"):
        return txt_to_text(path)
    elif path.endswith(".pptx"):
        return pptx_to_text(path)
    else:
        print("Unsupported file type")