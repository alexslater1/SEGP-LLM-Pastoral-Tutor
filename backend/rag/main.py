from transformers import BertTokenizer, BertModel
import torch
from chunker import document_chunker
from pdf_to_text import pdf_to_text


#other models can be used
embedding_model = 'bert-base-uncased'
tokenizer = BertTokenizer.from_pretrained(embedding_model)
model = BertModel.from_pretrained(embedding_model)

def embedder(document):
    #get tokens from input
    tokens = tokenizer(document, return_tensors='pt', padding=True, truncation=True)
    
    #embed tokens
    with torch.no_grad():  
        embeddings_output = model(**tokens)
    
    #get embeddings only
    embeddings = embeddings_output.last_hidden_state[:, 0, :]
    
    #numpy retains higher precision
    return embeddings[0].numpy()





#TODO: get document and convert to text files
pdf_path = "./pdfs/Student_Code_of_Conduct_2023_24.pdf"
txt_path = "./documents/Student_Code_of_Conduct_2023_24.txt"
pdf_to_text(pdf_path, txt_path)

#TODO: split text into chunks
chunks = document_chunker("./documents/", "BAAI/bge-small-en-v1.5")
print(chunks)

#embed chunks as vectors
embedded_chunks = []
for chunk in chunks:
   embedded_chunks.append(embedder(chunk))

#TODO: store vector with document
