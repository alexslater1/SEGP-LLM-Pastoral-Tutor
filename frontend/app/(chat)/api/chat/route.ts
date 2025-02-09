import {
  CoreUserMessage,
  CoreMessage,
  type Message,
  convertToCoreMessages,
  createDataStreamResponse,
  experimental_generateImage,
  streamObject,
  streamText,
} from 'ai';
import { z } from 'zod';

import {ResponseData, createResponse} from './response-generation';

import { auth } from '@/app/(auth)/auth';
import { customModel, imageGenerationModel } from '@/lib/ai';
import { models } from '@/lib/ai/models';
import {
  codePrompt,
  systemPrompt,
  updateDocumentPrompt,
} from '@/lib/ai/prompts';
import {
  deleteChatById,
  getChatById,
  getDocumentById,
  saveChat,
  saveDocument,
  saveMessages,
  saveSuggestions,
} from '@/lib/db/queries';
import type { Suggestion } from '@/lib/db/schema';
import {
  generateUUID,
  getMostRecentUserMessage,
  sanitizeResponseMessages,
} from '@/lib/utils';

import { generateTitleFromUserMessage } from '../../actions';

type AllowedTools =
  | 'createDocument'
  | 'updateDocument'
  | 'requestSuggestions'
  | 'getWeather';

const blocksTools: AllowedTools[] = [
  'createDocument',
  'updateDocument',
  'requestSuggestions',
];

const weatherTools: AllowedTools[] = ['getWeather'];

const API_FETCH_TIMEOUT_SECONDS: number = 120;

const allTools: AllowedTools[] = [...blocksTools, ...weatherTools];

export async function POST(request: Request) {
  const {
    id,
    messages,
    modelId,
  }: { id: string; messages: Array<Message>; modelId: string } =
    await request.json();

  const session = await auth();

  if (!session || !session.user || !session.user.id) {
    return new Response('Unauthorized', { status: 401 });
  }

  const model = models.find((model) => model.id === modelId);

  if (!model) {
    return new Response('Model not found', { status: 404 });
  }

  const coreMessages = convertToCoreMessages(messages);
  const userMessage = getMostRecentUserMessage(coreMessages);

  if (!userMessage) {
    return new Response('No user message found', { status: 400 });
  }

  const chat = await getChatById({ id });

  if (!chat) {
    const title = await generateTitleFromUserMessage({ message: userMessage });
    await saveChat({ id, userId: session.user.id, title });
  }

  const userMessageId = generateUUID();

  await saveMessages({
    messages: [
      { ...userMessage, id: userMessageId, createdAt: new Date(), chatId: id, annotations: [] },
    ],
  });

  async function makeCompletion(previousMessages: Message[], userMessage: CoreUserMessage): Promise<string> {
    type CompletionResponse = {
      response: string;
      reason: string;
    }
    try {
      let query = "Previous messages (oldest to most recent): " + previousMessages.slice(0, -1).map(message => {
        if (message.role == "user") {
          return "User: " + message.content;
        } else {
          return "Assistant: " + message.content;
        }
      }).join("\n");
      query += "\nMost recent user message to respond and answer to now: " + userMessage.content;

      const response = await fetch(process.env.BACKEND_AGENT_URL + "/completion" || "", {
        method: 'POST',
        signal: AbortSignal.timeout(1000 * API_FETCH_TIMEOUT_SECONDS),
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          query: query
        })
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      return (await response.json() as CompletionResponse).response;
    } catch (error) {
      console.error('Error:', error);
      return "Error fetching from backend";
    }
  }

  let result = createResponse({
    execute: async (response) => {
      response.addData ({
        type: 'user-message-id',
        content: userMessageId,
      })

      await response.setMessage(await makeCompletion(messages, userMessage));
    },
    onFinish: async (response) => {
      if (session.user?.id) {
        try {
          const messageId = generateUUID();
          if (response.getMessageRole() === 'assistant') {
            response.addAnnotation({
              messageIdFromServer: messageId,
            });
          }

          await saveMessages({
            messages: [{
              id: messageId,
              chatId: id,
              role: response.getMessageRole(),
              content: response.getMessage(),
              annotations: {},
              createdAt: new Date(),
            }]
          })
        } catch (error) {
          console.error('Failed to save chat');
        }
      }
    },
  })

  return result;
}

export async function DELETE(request: Request) {
  const { searchParams } = new URL(request.url);
  const id = searchParams.get('id');

  if (!id) {
    return new Response('Not Found', { status: 404 });
  }

  const session = await auth();

  if (!session || !session.user) {
    return new Response('Unauthorized', { status: 401 });
  }

  try {
    const chat = await getChatById({ id });

    if (chat.userId !== session.user.id) {
      return new Response('Unauthorized', { status: 401 });
    }

    await deleteChatById({ id });

    return new Response('Chat deleted', { status: 200 });
  } catch (error) {
    return new Response('An error occurred while processing your request', {
      status: 500,
    });
  }
}
