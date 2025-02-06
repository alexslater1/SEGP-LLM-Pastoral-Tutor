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

import { 
  StatusResponse,
  Status,
  makeV2InitialQuery, 
  makeV2StatusQuery 
} from './api-queries';

import { auth } from '@/app/(auth)/auth';
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

export const maxDuration = 60; // Setting timeout for Vercel serverless functions

const weatherTools: AllowedTools[] = ['getWeather'];

const allTools: AllowedTools[] = [...blocksTools, ...weatherTools];

const STATUS_QUERY_INTERVAL_SECONDS = 1

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

  let result = createResponse({
    execute: async (response) => {
      response.addData ({
        type: 'user-message-id',
        content: userMessageId,
      })

      let queryID = await makeV2InitialQuery(messages, userMessage);
      let status: StatusResponse = {type: Status.PENDING, current_action: "Thinking"};
      response.addAnnotation(status)
      while (status.type === Status.PENDING) {
        await new Promise(resolve => setTimeout(resolve, 1000 * STATUS_QUERY_INTERVAL_SECONDS));
        let newStatus = await makeV2StatusQuery(queryID);
        console.log(newStatus)
       
        if (newStatus.type === Status.FAILED || newStatus.type === Status.COMPLETED) {
          status = newStatus
          await response.setMessage(newStatus.type === Status.COMPLETED ?
            status.answer as string : status.error as string)
        } else if (newStatus.current_action && 
                   status.current_action !== newStatus.current_action) {
          status = newStatus
          response.addAnnotation(status)
          response.updateUIMessage()
        }
      }
      console.log("finished")
    },
    onFinish: async (response) => {
      console.log("starting")
      if (session.user?.id) {
        try {
          const messageId = generateUUID();
          if (response.getMessageRole() === 'assistant') {
            response.addAnnotation({
              messageIdFromServer: messageId,
            });
          }

          console.log("Saving:" + response.getMessage())

          await saveMessages({
            messages: [{
              id: messageId,
              chatId: id,
              role: response.getMessageRole(),
              content: response.getMessage(),
              annotations: response.getAnnotations(),
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
