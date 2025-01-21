import re
import os
import uuid
from transformers import AutoTokenizer

def document_chunker2(text, model_name,
                     paragraph_separator='\n\n',
                     chunk_size=250,
                     separator=' ',
                     secondary_chunking_regex=r'\S+?[\.,;!?]',
                     chunk_overlap=25):

    tokenizer = AutoTokenizer.from_pretrained(model_name)  # Load tokenizer for the specified model

    # Process each file using the existing chunking logic
    paragraphs = re.split(paragraph_separator, text)
    all_chunks = []
    for paragraph in paragraphs:
        words = paragraph.split(separator)
        current_chunk = ""
        chunks = []

        for word in words:
            new_chunk = current_chunk + (separator if current_chunk else '') + word
            if len(tokenizer.tokenize(new_chunk)) <= chunk_size:
                current_chunk = new_chunk
            else:
                if current_chunk:
                    chunks.append(current_chunk)
                current_chunk = word

        if current_chunk:
            chunks.append(current_chunk)

        refined_chunks = []
        for chunk in chunks:
            if len(tokenizer.tokenize(chunk)) > chunk_size:
                sub_chunks = re.split(secondary_chunking_regex, chunk)
                sub_chunk_accum = ""
                for sub_chunk in sub_chunks:
                    if sub_chunk_accum and len(tokenizer.tokenize(sub_chunk_accum + sub_chunk + ' ')) > chunk_size:
                        refined_chunks.append(sub_chunk_accum.strip())
                        sub_chunk_accum = sub_chunk
                    else:
                        sub_chunk_accum += (sub_chunk + ' ')
                if sub_chunk_accum:
                    refined_chunks.append(sub_chunk_accum.strip())
            else:
                refined_chunks.append(chunk)

        final_chunks = []
        if chunk_overlap > 0 and len(refined_chunks) > 1:
            for i in range(len(refined_chunks) - 1):
                final_chunks.append(refined_chunks[i])
                overlap_start = max(0, len(refined_chunks[i]) - chunk_overlap)
                overlap_end = min(chunk_overlap, len(refined_chunks[i+1]))
                overlap_chunk = refined_chunks[i][overlap_start:] + ' ' + refined_chunks[i+1][:overlap_end]
                final_chunks.append(overlap_chunk)
            final_chunks.append(refined_chunks[-1])
        else:
            final_chunks = refined_chunks
        
        for chunk in final_chunks:
            all_chunks.append(chunk)

    return all_chunks

def extract_emails_with_context(text, words_before=20, words_after=20):
    
    # Email regex pattern
    email_regex = r'[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}'
    
    results = []
    
    # Find all email occurrences in the text
    for match in re.finditer(email_regex, text):
        email = match.group()
        start, end = match.start(), match.end()
        
        # Get words before and after the email
        before_text = text[:start].split()[-words_before:]  # Last 20 words before email
        after_text = text[end:].split()[:words_after]  # First 20 words after email
        
        context = ' '.join(before_text + [email] + after_text)
        results.append((email, context))
    
    return results

def document_chunker(directory_path,
                     model_name,
                     paragraph_separator='\n\n',
                     chunk_size=250,
                     separator=' ',
                     secondary_chunking_regex=r'\S+?[\.,;!?]',
                     chunk_overlap=25):
    
    tokenizer = AutoTokenizer.from_pretrained(model_name)  # Load tokenizer for the specified model
    documents = {}  # Initialize dictionary to store results

    # Read each file in the specified directory
    for filename in os.listdir(directory_path):
        file_path = os.path.join(directory_path, filename)
        base = os.path.basename(file_path)
        sku = os.path.splitext(base)[0]
        if os.path.isfile(file_path):
            with open(file_path, 'r', encoding='utf-8') as file:
                text = file.read()

            # Generate a unique identifier for the document
            doc_id = str(uuid.uuid4())

            # Process each file using the existing chunking logic
            paragraphs = re.split(paragraph_separator, text)
            all_chunks = {}
            for paragraph in paragraphs:
                words = paragraph.split(separator)
                current_chunk = ""
                chunks = []

                for word in words:
                    new_chunk = current_chunk + (separator if current_chunk else '') + word
                    if len(tokenizer.tokenize(new_chunk)) <= chunk_size:
                        current_chunk = new_chunk
                    else:
                        if current_chunk:
                            chunks.append(current_chunk)
                        current_chunk = word

                if current_chunk:
                    chunks.append(current_chunk)

                refined_chunks = []
                for chunk in chunks:
                    if len(tokenizer.tokenize(chunk)) > chunk_size:
                        sub_chunks = re.split(secondary_chunking_regex, chunk)
                        sub_chunk_accum = ""
                        for sub_chunk in sub_chunks:
                            if sub_chunk_accum and len(tokenizer.tokenize(sub_chunk_accum + sub_chunk + ' ')) > chunk_size:
                                refined_chunks.append(sub_chunk_accum.strip())
                                sub_chunk_accum = sub_chunk
                            else:
                                sub_chunk_accum += (sub_chunk + ' ')
                        if sub_chunk_accum:
                            refined_chunks.append(sub_chunk_accum.strip())
                    else:
                        refined_chunks.append(chunk)

                final_chunks = []
                if chunk_overlap > 0 and len(refined_chunks) > 1:
                    for i in range(len(refined_chunks) - 1):
                        final_chunks.append(refined_chunks[i])
                        overlap_start = max(0, len(refined_chunks[i]) - chunk_overlap)
                        overlap_end = min(chunk_overlap, len(refined_chunks[i+1]))
                        overlap_chunk = refined_chunks[i][overlap_start:] + ' ' + refined_chunks[i+1][:overlap_end]
                        final_chunks.append(overlap_chunk)
                    final_chunks.append(refined_chunks[-1])
                else:
                    final_chunks = refined_chunks

                # Assign a UUID for each chunk and structure it with text and metadata
                for chunk in final_chunks:
                    chunk_id = str(uuid.uuid4())
                    all_chunks[chunk_id] = {"text": chunk, "metadata": {"file_name":sku}}  # Initialize metadata as dict

            # Map the document UUID to its chunk dictionary
            documents[doc_id] = all_chunks

    return documents

def print_chunks(chunks):
    for outer_key, outer_value in chunks.items():
        for inner_key, inner_value in outer_value.items():
            print("---------------------------------------------------")
            print(f"    Text: {inner_value.get('text', 'No text available')}")
