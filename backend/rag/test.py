import requests
import magic
import os

def upload_file(file_path: str, url: str):
    # Open the file in binary mode
    with open(file_path, 'rb') as file:
        # Get just the filename from the path
        file_name = os.path.basename(file_path)
        # Get the correct MIME type
        mime_type = magic.from_file(file_path, mime=True)
        # Create files dict with filename, not full path
        files = {'file': (file_name, file, mime_type)}
        
        # Send the POST request with the file
        response = requests.post(url, files=files)
        # Print the response status and content
        print("Status Code:", response.status_code)
        print("Response JSON:", response.json())

# Example usage
if __name__ == "__main__":
    file_path = '/home/alex/new/segp/rag/pdfs/testdoc.docx'  # Replace with the actual file path
    url = 'http://0.0.0.0:8000/rag-doc/'  # Replace with your FastAPI endpoint URL
    upload_file(file_path, url)
