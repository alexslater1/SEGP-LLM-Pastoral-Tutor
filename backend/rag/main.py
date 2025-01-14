from transformers import BertTokenizer, BertModel
import torch
from chunker import document_chunker


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





#TODO: get document

#TODO: get text from document
text = "Hello World!" #placeholder text

#TODO: split text into chunks
chunks = document_chunker("./documents/", "BAAI/bge-small-en-v1.5")

#embed chunks as vectors
embedded_chunks = []
for chunk in chunks:
    embedded_chunks.append(embedder(chunk))

#TODO: store vector with document
