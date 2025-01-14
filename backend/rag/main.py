from chunker import document_chunker, print_chunks
from embedder import embed
from searcher import search
from pdf_to_text import pdf_to_text

#TODO: get document and convert to text files
pdf_path = "./pdfs/Student_Code_of_Conduct_2023_24.pdf"
txt_path = "./documents/Student_Code_of_Conduct_2023_24.txt"
pdf_to_text(pdf_path, txt_path)

#TODO: split text into chunks
chunks = document_chunker("./documents/", "BAAI/bge-small-en-v1.5")

text_chunks = []
for outer_key, outer_value in chunks.items():
    for inner_key, inner_value in outer_value.items():
        if 'text' in inner_value:
            # print("--------------------")
            # print(inner_value['text'])
            text_chunks.append(inner_value['text'])


#embed chunks as vectors
embedded_chunks = []
for chunk in text_chunks:
   embedded_chunks.append((chunk, embed(chunk)))

#TODO: store vector with document

#TODO: search chunks
query = "Hello World!"
return_chunk = search(query, embedded_chunks)
print(return_chunk)

