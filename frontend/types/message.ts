export type JSONValue =
  | string
  | number
  | boolean
  | null
  | JSONValue[]
  | { [key: string]: JSONValue };

export enum Status {
  PENDING = "pending",
  COMPLETED = "completed",
  ERROR = "error",
}

export enum Role {
  SYSTEM = "system",
  USER = "user",
  ASSISTANT = "assistant",
  DATA = "data",
}

export type Message = {
  id: string;
  requestID: string;
  createdAt?: Date;
  content: string;
  role: Role;
  data?: JSONValue;
  actions: string[];
  status: Status;
};
