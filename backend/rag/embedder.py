import torch


def embed(text, tokenizer, model):
    #get tokens from input
    tokens = tokenizer(text, return_tensors='pt', padding=True, truncation=True)
    
    #embed tokens
    with torch.no_grad():  
        embeddings_output = model(**tokens)
    
    #get embeddings only
    embeddings = embeddings_output.last_hidden_state[:, 0, :]
    
    #numpy retains higher precision
    return embeddings[0].numpy()
