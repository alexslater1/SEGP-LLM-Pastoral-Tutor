export interface RAGDocument {
  id: string;
  name: string;
  uploadedAt: Date;
  size: string;
  type: 'PDF' | 'DOCX' | 'TXT' | 'PPTX';
} 