from transformers import BertTokenizer, BertModel
import torch

#other models can be used
embedding_model = 'bert-base-uncased'
tokenizer = BertTokenizer.from_pretrained(embedding_model)
model = BertModel.from_pretrained(embedding_model)

def embed(text):
    #get tokens from input
    tokens = tokenizer(text, return_tensors='pt', padding=True, truncation=True)
    
    #embed tokens
    with torch.no_grad():  
        embeddings_output = model(**tokens)
    
    #get embeddings only
    embeddings = embeddings_output.last_hidden_state[:, 0, :]
    
    #numpy retains higher precision
    return embeddings[0].numpy()
